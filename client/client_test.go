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
	"testing"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"github.com/docker/notary/certs"
	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/server"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
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

func initializeRepo(t *testing.T, rootType, tempBaseDir, gun, url string,
	serverManagesSnapshot bool) (*NotaryRepository, string) {

	serverManagedRoles := []string{}
	if serverManagesSnapshot {
		serverManagedRoles = []string{data.CanonicalSnapshotRole}
	}

	repo, rootPubKeyID := createRepoAndKey(t, rootType, tempBaseDir, gun, url)

	err := repo.Initialize(rootPubKeyID, serverManagedRoles...)
	assert.NoError(t, err, "error creating repository: %s", err)

	return repo, rootPubKeyID
}

// Creates a new repository and adds a root key.  Returns the repo and key ID.
func createRepoAndKey(t *testing.T, rootType, tempBaseDir, gun, url string) (
	*NotaryRepository, string) {

	repo, err := NewNotaryRepository(
		tempBaseDir, gun, url, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	assert.NoError(t, err, "error generating root key: %s", err)

	return repo, rootPubKey.ID()
}

// Initializing a new repo while specifying that the server should manage the root
// role will fail.
func TestInitRepositoryManagedRolesIncludingRoot(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, data.CanonicalRootRole)
	assert.Error(t, err)
	assert.Contains(t, err.Error(),
		"Notary does not support the server managing the root key")
}

// Initializing a new repo while specifying that the server should manage some
// invalid role will fail.
func TestInitRepositoryManagedRolesInvalidRole(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, "randomrole")
	assert.Error(t, err)
	assert.Contains(t, err.Error(),
		"Notary does not support the server managing the randomrole key")
}

// Initializing a new repo while specifying that the server should manage the
// targets role will fail.
func TestInitRepositoryManagedRolesIncludingTargets(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, data.CanonicalTargetsRole)
	assert.Error(t, err)
	assert.Contains(t, err.Error(),
		"Notary does not support the server managing the targets key")
}

// Initializing a new repo while specifying that the server should manage the
// timestamp key is fine - that's what it already does, so no error.
func TestInitRepositoryManagedRolesIncludingTimestamp(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", ts.URL)
	err = repo.Initialize(rootPubKeyID, data.CanonicalTimestampRole)
	assert.NoError(t, err)
}

// passing timestamp + snapshot, or just snapshot, is tested in the next two
// test cases.

// TestInitRepoServerOnlyManagesTimestampKey runs through the process of
// initializing a repository and makes sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key.
// This test case covers the default case where the server only manages the
// timestamp key.
func TestInitRepoServerOnlyManagesTimestampKey(t *testing.T) {
	testInitRepo(t, data.ECDSAKey, false)
	if !testing.Short() {
		testInitRepo(t, data.RSAKey, false)
	}
}

// TestInitRepoServerManagesTimestampAndSnapshotKeys runs through the process of
// initializing a repository and makes sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key.
// This test case covers the server managing both the timestap and snapshot keys.
func TestInitRepoServerManagesTimestampAndSnapshotKeys(t *testing.T) {
	testInitRepo(t, data.ECDSAKey, true)
	if !testing.Short() {
		testInitRepo(t, data.RSAKey, true)
	}
}

// This creates a new KeyFileStore in the repo's base directory and makes sure
// the repo has the right number of keys
func assertRepoHasExpectedKeys(t *testing.T, repo *NotaryRepository,
	rootKeyID string, expectedSnapshotKey bool) {

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
	alwaysThere := []string{data.CanonicalRootRole, data.CanonicalTargetsRole}
	for _, role := range alwaysThere {
		_, ok := roles[role]
		assert.True(t, ok, "missing %s key", role)
	}

	// there may be a snapshots key, depending on whether the server is managing
	// the snapshots key
	_, ok := roles[data.CanonicalSnapshotRole]
	if expectedSnapshotKey {
		assert.True(t, ok, "missing snapshot key")
	} else {
		assert.False(t, ok,
			"there should be no snapshot key because the server manages it")
	}

	// The server manages the timestamp key - there should not be a timestamp
	// key
	_, ok = roles[data.CanonicalTimestampRole]
	assert.False(t, ok,
		"there should be no timestamp key because the server manages it")
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

// Sanity check the TUF metadata files. Verify that it exists for a particular
// role, the JSON is well-formed, and the signatures exist.
// For the root.json file, also check that the root, snapshot, and
// targets key IDs are present.
func assertRepoHasExpectedMetadata(t *testing.T, repo *NotaryRepository,
	role string, expected bool) {

	filename := filepath.Join(tufDir, filepath.FromSlash(repo.gun),
		"metadata", role+".json")
	fullPath := filepath.Join(repo.baseDir, filename)
	_, err := os.Stat(fullPath)

	if expected {
		assert.NoError(t, err, "missing TUF metadata file: %s", filename)
	} else {
		assert.Error(t, err,
			"%s metadata should not exist, but does: %s", role, filename)
		return
	}

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
	if role == data.CanonicalRootRole {
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

func testInitRepo(t *testing.T, rootType string, serverManagesSnapshot bool) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootKeyID := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL,
		serverManagesSnapshot)

	assertRepoHasExpectedKeys(t, repo, rootKeyID, !serverManagesSnapshot)
	assertRepoHasExpectedCerts(t, repo)
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalRootRole, true)
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalTargetsRole, true)
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole,
		!serverManagesSnapshot)
}

// TestInitRepoAttemptsExceeded tests error handling when passphrase.Retriever
// (or rather the user) insists on an incorrect password.
func TestInitRepoAttemptsExceeded(t *testing.T) {
	testInitRepoAttemptsExceeded(t, data.ECDSAKey)
	if !testing.Short() {
		testInitRepoAttemptsExceeded(t, data.RSAKey)
	}
}

func testInitRepoAttemptsExceeded(t *testing.T, rootType string) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	retriever := passphrase.ConstantRetriever("password")
	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	assert.NoError(t, err, "error creating repo: %s", err)
	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	assert.NoError(t, err, "error generating root key: %s", err)

	retriever = passphrase.ConstantRetriever("incorrect password")
	// repo.CryptoService’s FileKeyStore caches the unlocked private key, so to test
	// private key unlocking we need a new repo instance.
	repo, err = NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	assert.NoError(t, err, "error creating repo: %s", err)
	err = repo.Initialize(rootPubKey.ID())
	assert.EqualError(t, err, trustmanager.ErrAttemptsExceeded{}.Error())
}

// TestInitRepoPasswordInvalid tests error handling when passphrase.Retriever
// (or rather the user) fails to provide a correct password.
func TestInitRepoPasswordInvalid(t *testing.T) {
	testInitRepoPasswordInvalid(t, data.ECDSAKey)
	if !testing.Short() {
		testInitRepoPasswordInvalid(t, data.RSAKey)
	}
}

func giveUpPassphraseRetriever(_, _ string, _ bool, _ int) (string, bool, error) {
	return "", true, nil
}

func testInitRepoPasswordInvalid(t *testing.T, rootType string) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	retriever := passphrase.ConstantRetriever("password")
	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	assert.NoError(t, err, "error creating repo: %s", err)
	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	assert.NoError(t, err, "error generating root key: %s", err)

	// repo.CryptoService’s FileKeyStore caches the unlocked private key, so to test
	// private key unlocking we need a new repo instance.
	repo, err = NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, giveUpPassphraseRetriever)
	assert.NoError(t, err, "error creating repo: %s", err)
	err = repo.Initialize(rootPubKey.ID())
	assert.EqualError(t, err, trustmanager.ErrPasswordInvalid{}.Error())
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

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)

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

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)

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

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)

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

	initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)

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

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)
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

// Create a repo, instantiate a notary server, and publish the repo to the
// server, signing all the non-timestamp metadata.
// We test this with both an RSA and ECDSA root key
func TestPublishClientHasSnapshotKey(t *testing.T) {
	testPublish(t, data.ECDSAKey, false)
	if !testing.Short() {
		testPublish(t, data.RSAKey, false)
	}
}

// Create a repo, instantiate a notary server (designating the server as the
// snapshot signer) , and publish the repo to the server, signing the root and
// targets metadata only.  The server should sign just fine.
// We test this with both an RSA and ECDSA root key
func TestPublishAfterInitServerHasSnapshotKey(t *testing.T) {
	testPublish(t, data.ECDSAKey, true)
	if !testing.Short() {
		testPublish(t, data.RSAKey, true)
	}
}

func testPublish(t *testing.T, rootType string, serverManagesSnapshot bool) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL,
		serverManagesSnapshot)
	assertPublishSucceeds(t, repo)
}

func assertPublishSucceeds(t *testing.T, repo1 *NotaryRepository) {
	// Create 2 targets
	latestTarget := addTarget(t, repo1, "latest", "../fixtures/intermediate-ca.crt")
	currentTarget := addTarget(t, repo1, "current", "../fixtures/intermediate-ca.crt")
	assert.Len(t, getChanges(t, repo1), 2, "wrong number of changelist files found")

	// Now test Publish
	err := repo1.Publish()
	assert.NoError(t, err)
	assert.Len(t, getChanges(t, repo1), 0, "wrong number of changelist files found")

	// Create a new repo and pull from the server
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo2, err := NewNotaryRepository(tempBaseDir, repo1.gun, repo1.baseURL,
		http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Should be two targets
	for _, repo := range []*NotaryRepository{repo1, repo2} {
		targets, err := repo.ListTargets()
		assert.NoError(t, err)

		assert.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

		sort.Stable(targetSorter(targets))

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
}

// After pulling a repo from the server, so there is a snapshots metadata file,
// push a different target to the server (the server is still the snapshot
// signer).  The server should sign just fine.
// We test this with both an RSA and ECDSA root key
func TestPublishAfterPullServerHasSnapshotKey(t *testing.T) {
	testPublishAfterPullServerHasSnapshotKey(t, data.ECDSAKey)
	if !testing.Short() {
		testPublishAfterPullServerHasSnapshotKey(t, data.RSAKey)
	}
}

func testPublishAfterPullServerHasSnapshotKey(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, true)
	// no timestamp metadata because that comes from the server
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, false)
	// no snapshot metadata because that comes from the server
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, false)

	// Publish something
	published := addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")
	err = repo.Publish()
	assert.NoError(t, err)
	// still no timestamp or snapshot metadata info
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, false)
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, false)

	// list, so that the snapshot metadata is pulled from server
	targets, err := repo.ListTargets()
	assert.NoError(t, err)
	assert.Equal(t, []*Target{published}, targets)
	// listing downloaded the timestamp and snapshot metadata info
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, true)
	assertRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, true)

	// Publish again should succeed
	addTarget(t, repo, "v2", "../fixtures/intermediate-ca.crt")
	err = repo.Publish()
	assert.NoError(t, err)
}

// If neither the client nor the server has the snapshot key, signing will fail
// with an ErrNoKeys error.
// We test this with both an RSA and ECDSA root key
func TestPublishNoOneHasSnapshotKey(t *testing.T) {
	testPublishNoOneHasSnapshotKey(t, data.ECDSAKey)
	if !testing.Short() {
		testPublishNoOneHasSnapshotKey(t, data.RSAKey)
	}
}

func testPublishNoOneHasSnapshotKey(t *testing.T, rootType string) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	// create repo and delete the snapshot key and metadata
	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, false)
	snapshotRole, ok := repo.tufRepo.Root.Signed.Roles[data.CanonicalSnapshotRole]
	assert.True(t, ok)
	for _, keyID := range snapshotRole.KeyIDs {
		repo.CryptoService.RemoveKey(keyID)
	}

	// Publish something
	addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")
	err = repo.Publish()
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)
}

// If the snapshot metadata is corrupt, whether the client or server has the
// snapshot key, we can't publish.
// We test this with both an RSA and ECDSA root key
func TestPublishSnapshotCorrupt(t *testing.T) {
	testPublishSnapshotCorrupt(t, data.ECDSAKey, true)
	testPublishSnapshotCorrupt(t, data.ECDSAKey, false)
	if !testing.Short() {
		testPublishSnapshotCorrupt(t, data.RSAKey, true)
		testPublishSnapshotCorrupt(t, data.RSAKey, false)
	}
}

func testPublishSnapshotCorrupt(t *testing.T, rootType string, serverManagesSnapshot bool) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, tempBaseDir, gun, ts.URL, serverManagesSnapshot)
	// write a corrupt snapshots file
	repo.fileStore.SetMeta(data.CanonicalSnapshotRole, []byte("this isn't JSON"))

	addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")
	err = repo.Publish()
	assert.Error(t, err)
}

func TestRotate(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, tempBaseDir, gun, ts.URL, false)

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
