package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"github.com/docker/notary/certs"
	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/server"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/store"
	"github.com/jfrazelle/go/canonical/json"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func simpleTestServer(t *testing.T) (
	*httptest.Server, *http.ServeMux, map[string]data.PrivateKey) {

	roles := []string{data.CanonicalTimestampRole, data.CanonicalSnapshotRole}
	keys := make(map[string]data.PrivateKey)
	mux := http.NewServeMux()

	for _, role := range roles {
		key, err := trustmanager.GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err)

		keys[role] = key
		pubKey := data.PublicKeyFromPrivate(key)
		jsonBytes, err := json.MarshalCanonical(&pubKey)
		assert.NoError(t, err)
		keyJSON := string(jsonBytes)

		// TUF will request /v2/docker.com/notary/_trust/tuf/<role>.key
		mux.HandleFunc(
			fmt.Sprintf("/v2/docker.com/notary/_trust/tuf/%s.key", role),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, keyJSON)
			})
	}

	ts := httptest.NewServer(mux)
	return ts, mux, keys
}

func fullTestServer(t *testing.T) *httptest.Server {
	// Set up server
	ctx := context.WithValue(
		context.Background(), "metaStore", storage.NewMemStorage())

	// Do not pass one of the const KeyAlgorithms here as the value! Passing a
	// string is in itself good test that we are handling it correctly as we
	// will be receiving a string from the configuration.
	ctx = context.WithValue(ctx, "keyAlgorithm", "ecdsa")

	// Eat the logs instead of spewing them out
	var b bytes.Buffer
	l := logrus.New()
	l.Out = &b
	ctx = ctxu.WithLogger(ctx, logrus.NewEntry(l))

	cryptoService := cryptoservice.NewCryptoService(
		"", trustmanager.NewKeyMemoryStore(passphraseRetriever))
	return httptest.NewServer(server.RootHandler(nil, ctx, cryptoService))
}

// server that returns some particular error code all the time
func errorTestServer(t *testing.T, errorCode int) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorCode)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	return server
}

func initializeRepo(t *testing.T, rootType, tempBaseDir, gun, url string) (
	*NotaryRepository, string) {

	repo, err := NewNotaryRepository(
		tempBaseDir, gun, url, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	assert.NoError(t, err, "error generating root key: %s", err)

	err = repo.Initialize(rootPubKey.ID())
	assert.NoError(t, err, "error creating repository: %s", err)

	return repo, rootPubKey.ID()
}

// TestInitRepo runs through the process of initializing a repository and makes
// sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key
func TestInitRepo(t *testing.T) {
	testInitRepo(t, data.ECDSAKey)
	if !testing.Short() {
		testInitRepo(t, data.RSAKey)
	}
}

// This creates a new KeyFileStore in the repo's base directory and makes sure
// the repo has the right number of keys
func assertRepoHasExpectedKeys(t *testing.T, repo *NotaryRepository,
	rootKeyID string) {

	// The repo should have a keyFileStore and have created keys using it,
	// so create a new KeyFileStore, and check that the keys do exist and are
	// valid
	ks, err := trustmanager.NewKeyFileStore(repo.baseDir, passphraseRetriever)
	assert.NoError(t, err)

	roles := make(map[string]bool)
	for keyID, role := range ks.ListKeys() {
		if role == data.CanonicalRootRole {
			assert.Equal(t, rootKeyID, keyID, "Unexpected root key ID")
		}
		// just to ensure the content of the key files created are valid
		_, r, err := ks.GetKey(keyID)
		assert.NoError(t, err)
		assert.Equal(t, role, r)
		roles[role] = true
	}
	// there is a root key and a targets key
	for _, role := range data.ValidRoles {
		if role != data.CanonicalTimestampRole {
			_, ok := roles[role]
			assert.True(t, ok, fmt.Sprintf("missing %s key", role))
		}
	}

	// The server manages the timestamp key - there should not be a timestamp
	// key
	_, ok := roles[data.CanonicalTimestampRole]
	assert.False(t, ok)
}

// This creates a new certificate manager in the repo's base directory and
// makes sure the repo has the right certificates
func assertRepoHasExpectedCerts(t *testing.T, repo *NotaryRepository) {
	// The repo should have a certificate manager and have created certs using
	// it, so create a new manager, and check that the certs do exist and
	// are valid
	certManager, err := certs.NewManager(repo.baseDir)
	assert.NoError(t, err)
	certificates := certManager.TrustedCertificateStore().GetCertificates()
	assert.Len(t, certificates, 1, "unexpected number of trusted certificates")

	certID, err := trustmanager.FingerprintCert(certificates[0])
	assert.NoError(t, err, "unable to fingerprint the trusted certificate")
	assert.NotEqual(t, certID, "")
}

// Sanity check the TUF metadata files. Verify that they exist, the JSON is
// well-formed, and the signatures exist. For the root.json file, also check
// that the root, snapshot, and targets key IDs are present.
func assertRepoHasExpectedMetadata(t *testing.T, repo *NotaryRepository) {
	expectedTUFMetadataFiles := []string{
		filepath.Join(tufDir, filepath.FromSlash(repo.gun), "metadata", "root.json"),
		filepath.Join(tufDir, filepath.FromSlash(repo.gun), "metadata", "snapshot.json"),
		filepath.Join(tufDir, filepath.FromSlash(repo.gun), "metadata", "targets.json"),
	}
	for _, filename := range expectedTUFMetadataFiles {
		fullPath := filepath.Join(repo.baseDir, filename)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "missing TUF metadata file: %s", filename)

		jsonBytes, err := ioutil.ReadFile(fullPath)
		assert.NoError(t, err, "error reading TUF metadata file %s: %s", filename, err)

		var decoded data.Signed
		err = json.Unmarshal(jsonBytes, &decoded)
		assert.NoError(t, err, "error parsing TUF metadata file %s: %s", filename, err)

		assert.Len(t, decoded.Signatures, 1,
			"incorrect number of signatures in TUF metadata file %s", filename)

		assert.NotEmpty(t, decoded.Signatures[0].KeyID,
			"empty key ID field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Method,
			"empty method field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Signature,
			"empty signature in TUF metadata file %s", filename)

		// Special case for root.json: also check that the signed
		// content for keys and roles
		if strings.HasSuffix(filename, "root.json") {
			var decodedRoot data.Root
			err := json.Unmarshal(decoded.Signed, &decodedRoot)
			assert.NoError(t, err, "error parsing root.json signed section: %s", err)

			assert.Equal(t, "Root", decodedRoot.Type, "_type mismatch in root.json")

			// Expect 1 key for each valid role in the Keys map - one for
			// each of root, targets, snapshot, timestamp
			assert.Len(t, decodedRoot.Keys, len(data.ValidRoles),
				"wrong number of keys in root.json")
			assert.Len(t, decodedRoot.Roles, len(data.ValidRoles),
				"wrong number of roles in root.json")

			for role := range data.ValidRoles {
				_, ok := decodedRoot.Roles[role]
				assert.True(t, ok, "Missing role %s in root.json", role)
			}
		}
	}
}

func testInitRepo(t *testing.T, rootType string) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootKeyID := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	assertRepoHasExpectedKeys(t, repo, rootKeyID)
	assertRepoHasExpectedCerts(t, repo)
	assertRepoHasExpectedMetadata(t, repo)
}

// TestAddTarget adds a target to the repo and confirms that the changelist
// is updated correctly.
// We test this with both an RSA and ECDSA root key
func TestAddTarget(t *testing.T) {
	testAddTarget(t, data.ECDSAKey)
	if !testing.Short() {
		testAddTarget(t, data.RSAKey)
	}
}

func addTarget(t *testing.T, repo *NotaryRepository, targetName, targetFile string) *Target {
	target, err := NewTarget(targetName, targetFile)
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(target)
	assert.NoError(t, err, "error adding target")
	return target
}

// calls GetChangelist and gets the actual changes out
func getChanges(t *testing.T, repo *NotaryRepository) []changelist.Change {
	changeList, err := repo.GetChangelist()
	assert.NoError(t, err)
	return changeList.List()
}

func testAddTarget(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err = repo.tufRepo.InitTimestamp()
	assert.NoError(t, err, "error creating repository: %s", err)
	assert.Len(t, getChanges(t, repo), 0, "should start with zero changes")

	// Add fixtures/intermediate-ca.crt as a target. There's no particular
	// reason for using this file except that it happens to be available as
	// a fixture.
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	changes := getChanges(t, repo)
	assert.Len(t, changes, 1, "wrong number of changes files found")

	for _, c := range changes { // there is only one
		assert.EqualValues(t, changelist.ActionCreate, c.Action())
		assert.Equal(t, "targets", c.Scope())
		assert.Equal(t, "target", c.Type())
		assert.Equal(t, "latest", c.Path())
		assert.NotEmpty(t, c.Content())
	}

	// Create a second target
	addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")
	changes = getChanges(t, repo)
	assert.Len(t, changes, 2, "wrong number of changelist files found")

	newFileFound := false
	for _, c := range changes {
		if c.Path() != "latest" {
			assert.EqualValues(t, changelist.ActionCreate, c.Action())
			assert.Equal(t, "targets", c.Scope())
			assert.Equal(t, "target", c.Type())
			assert.Equal(t, "current", c.Path())
			assert.NotEmpty(t, c.Content())

			newFileFound = true
		}
	}
	assert.True(t, newFileFound, "second changelist file not found")
}

// TestListTarget fakes serving signed metadata files over the test's
// internal HTTP server to ensure that ListTargets returns the correct number
// of listed targets.
// We test this with both an RSA and ECDSA root key
func TestListTarget(t *testing.T) {
	testListEmptyTargets(t, data.ECDSAKey)
	testListTarget(t, data.ECDSAKey)
	if !testing.Short() {
		testListEmptyTargets(t, data.RSAKey)
		testListTarget(t, data.RSAKey)
	}
}

func testListEmptyTargets(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err = repo.tufRepo.InitTimestamp()
	assert.NoError(t, err, "error creating repository: %s", err)

	_, err = repo.ListTargets()
	assert.Error(t, err) // no trust data
}

// reads data from the repository in order to fake data being served via
// the ServeMux.
func fakeServerData(t *testing.T, repo *NotaryRepository, mux *http.ServeMux,
	keys map[string]data.PrivateKey) {

	timestampKey, ok := keys[data.CanonicalTimestampRole]
	assert.True(t, ok)
	savedTUFRepo := repo.tufRepo // in case this is overwritten

	fileStore, err := trustmanager.NewKeyFileStore(repo.baseDir, passphraseRetriever)
	assert.NoError(t, err)
	fileStore.AddKey(
		filepath.Join(filepath.FromSlash(repo.gun), timestampKey.ID()),
		"nonroot", timestampKey)

	rootJSONFile := filepath.Join(repo.baseDir, "tuf",
		filepath.FromSlash(repo.gun), "metadata", "root.json")
	rootFileBytes, err := ioutil.ReadFile(rootJSONFile)

	signedTargets, err := savedTUFRepo.SignTargets(
		"targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)

	signedSnapshot, err := savedTUFRepo.SignSnapshot(
		data.DefaultExpires("snapshot"))
	assert.NoError(t, err)

	signedTimestamp, err := savedTUFRepo.SignTimestamp(
		data.DefaultExpires("timestamp"))
	assert.NoError(t, err)

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/root.json",
		func(w http.ResponseWriter, r *http.Request) {
			assert.NoError(t, err)
			fmt.Fprint(w, string(rootFileBytes))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.json",
		func(w http.ResponseWriter, r *http.Request) {
			timestampJSON, _ := json.Marshal(signedTimestamp)
			fmt.Fprint(w, string(timestampJSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/snapshot.json",
		func(w http.ResponseWriter, r *http.Request) {
			snapshotJSON, _ := json.Marshal(signedSnapshot)
			fmt.Fprint(w, string(snapshotJSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets.json",
		func(w http.ResponseWriter, r *http.Request) {
			targetsJSON, _ := json.Marshal(signedTargets)
			fmt.Fprint(w, string(targetsJSON))
		})
}

// We want to sort by name, so we can guarantee ordering.
type targetSorter []*Target

func (k targetSorter) Len() int           { return len(k) }
func (k targetSorter) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k targetSorter) Less(i, j int) bool { return k[i].Name < k[j].Name }

func testListTarget(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, mux, keys := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err = repo.tufRepo.InitTimestamp()
	assert.NoError(t, err, "error creating repository: %s", err)

	latestTarget := addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	currentTarget := addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")

	// Apply the changelist. Normally, this would be done by Publish

	// load the changelist for this repo
	cl, err := changelist.NewFileChangelist(
		filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "changelist"))
	assert.NoError(t, err, "could not open changelist")

	// apply the changelist to the repo
	err = applyChangelist(repo.tufRepo, cl)
	assert.NoError(t, err, "could not apply changelist")

	fakeServerData(t, repo, mux, keys)

	targets, err := repo.ListTargets()
	assert.NoError(t, err)

	// Should be two targets
	assert.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

	sort.Stable(targetSorter(targets))

	// current should be first
	assert.Equal(t, currentTarget, targets[0], "current target does not match")
	assert.Equal(t, latestTarget, targets[1], "latest target does not match")

	// Also test GetTargetByName
	newLatestTarget, err := repo.GetTargetByName("latest")
	assert.NoError(t, err)
	assert.Equal(t, latestTarget, newLatestTarget, "latest target does not match")

	newCurrentTarget, err := repo.GetTargetByName("current")
	assert.NoError(t, err)
	assert.Equal(t, currentTarget, newCurrentTarget, "current target does not match")
}

// TestValidateRootKey verifies that the public data in root.json for the root
// key is a valid x509 certificate.
func TestValidateRootKey(t *testing.T) {
	testValidateRootKey(t, data.ECDSAKey)
	if !testing.Short() {
		testValidateRootKey(t, data.RSAKey)
	}
}

func testValidateRootKey(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	rootJSONFile := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "metadata", "root.json")

	jsonBytes, err := ioutil.ReadFile(rootJSONFile)
	assert.NoError(t, err, "error reading TUF metadata file %s: %s", rootJSONFile, err)

	var decoded data.Signed
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err, "error parsing TUF metadata file %s: %s", rootJSONFile, err)

	var decodedRoot data.Root
	err = json.Unmarshal(decoded.Signed, &decodedRoot)
	assert.NoError(t, err, "error parsing root.json signed section: %s", err)

	keyids := []string{}
	for role, roleData := range decodedRoot.Roles {
		if role == "root" {
			keyids = append(keyids, roleData.KeyIDs...)
		}
	}
	assert.NotEmpty(t, keyids)

	for _, keyid := range keyids {
		if key, ok := decodedRoot.Keys[keyid]; !ok {
			t.Fatal("key id not found in keys")
		} else {
			_, err := trustmanager.LoadCertFromPEM(key.Public())
			assert.NoError(t, err, "key is not a valid cert")
		}
	}
}

// TestGetChangelist ensures that the changelist returned matches the changes
// added.
// We test this with both an RSA and ECDSA root key
func TestGetChangelist(t *testing.T) {
	testGetChangelist(t, data.ECDSAKey)
	if !testing.Short() {
		testGetChangelist(t, data.RSAKey)
	}
}

func testGetChangelist(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)
	assert.Len(t, getChanges(t, repo), 0, "No changes should be in changelist yet")

	// Create 2 targets
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")

	// Test loading changelist
	chgs := getChanges(t, repo)
	assert.Len(t, chgs, 2, "Wrong number of changes returned from changelist")

	changes := make(map[string]changelist.Change)
	for _, ch := range chgs {
		changes[ch.Path()] = ch
	}

	currentChange := changes["current"]
	assert.NotNil(t, currentChange, "Expected changelist to contain a change for path 'current'")
	assert.EqualValues(t, changelist.ActionCreate, currentChange.Action())
	assert.Equal(t, "targets", currentChange.Scope())
	assert.Equal(t, "target", currentChange.Type())
	assert.Equal(t, "current", currentChange.Path())

	latestChange := changes["latest"]
	assert.NotNil(t, latestChange, "Expected changelist to contain a change for path 'latest'")
	assert.EqualValues(t, changelist.ActionCreate, latestChange.Action())
	assert.Equal(t, "targets", latestChange.Scope())
	assert.Equal(t, "target", latestChange.Type())
	assert.Equal(t, "latest", latestChange.Path())
}

// TestPublish creates a repo, instantiates a notary server, and publishes
// the repo to the server.
// We test this with both an RSA and ECDSA root key
func TestPublish(t *testing.T) {
	testPublish(t, data.ECDSAKey)
	if !testing.Short() {
		testPublish(t, data.RSAKey)
	}
}

func testPublish(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()
	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL)

	// Create 2 targets
	latestTarget := addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	currentTarget := addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")
	assert.Len(t, getChanges(t, repo), 2, "wrong number of changelist files found")

	// Now test Publish
	err = repo.Publish()
	assert.NoError(t, err)
	assert.Len(t, getChanges(t, repo), 0, "wrong number of changelist files found")

	// Create a new repo and pull from the server
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo2, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	targets, err := repo2.ListTargets()
	assert.NoError(t, err)

	// Should be two targets
	assert.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

	sort.Stable(targetSorter(targets))

	assert.Equal(t, currentTarget, targets[0], "current target does not match")
	assert.Equal(t, latestTarget, targets[1], "latest target does not match")

	// Also test GetTargetByName
	newLatestTarget, err := repo2.GetTargetByName("latest")
	assert.NoError(t, err)
	assert.Equal(t, latestTarget, newLatestTarget, "latest target does not match")

	newCurrentTarget, err := repo2.GetTargetByName("current")
	assert.NoError(t, err)
	assert.Equal(t, currentTarget, newCurrentTarget, "current target does not match")
}

func TestRotate(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, tempBaseDir, gun, ts.URL)

	// Adding a target will allow us to confirm the repository is still valid after
	// rotating the keys.
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")

	// Publish
	err = repo.Publish()
	assert.NoError(t, err)

	// Get root.json and capture targets + snapshot key IDs
	repo.GetTargetByName("latest") // force a pull
	targetsKeyIDs := repo.tufRepo.Root.Signed.Roles["targets"].KeyIDs
	snapshotKeyIDs := repo.tufRepo.Root.Signed.Roles["snapshot"].KeyIDs
	assert.Len(t, targetsKeyIDs, 1)
	assert.Len(t, snapshotKeyIDs, 1)

	// Do rotation
	repo.RotateKeys()

	// Publish
	err = repo.Publish()
	assert.NoError(t, err)

	// Get root.json. Check targets + snapshot keys have changed
	// and that they match those found in the changelist.
	_, err = repo.GetTargetByName("latest") // force a pull
	assert.NoError(t, err)
	newTargetsKeyIDs := repo.tufRepo.Root.Signed.Roles["targets"].KeyIDs
	newSnapshotKeyIDs := repo.tufRepo.Root.Signed.Roles["snapshot"].KeyIDs
	assert.Len(t, newTargetsKeyIDs, 1)
	assert.Len(t, newSnapshotKeyIDs, 1)
	assert.NotEqual(t, targetsKeyIDs[0], newTargetsKeyIDs[0])
	assert.NotEqual(t, snapshotKeyIDs[0], newSnapshotKeyIDs[0])

	// Confirm changelist dir empty after publishing changes
	changes := getChanges(t, repo)
	assert.Len(t, changes, 0, "wrong number of changelist files found")
}

// If there is no local cache, notary operations return the remote error code
func TestRemoteServerUnavailableNoLocalCache(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts := errorTestServer(t, 500)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, "docker.com/notary",
		ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repo: %s", err)

	_, err = repo.ListTargets()
	assert.Error(t, err)
	assert.IsType(t, store.ErrServerUnavailable{}, err)

	_, err = repo.GetTargetByName("targetName")
	assert.Error(t, err)
	assert.IsType(t, store.ErrServerUnavailable{}, err)

	err = repo.Publish()
	assert.Error(t, err)
	assert.IsType(t, store.ErrServerUnavailable{}, err)
}
