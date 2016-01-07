package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"github.com/jfrazelle/go/canonical/json"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

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
	"github.com/docker/notary/tuf/validation"
)

func simpleTestServer(t *testing.T, roles ...string) (
	*httptest.Server, *http.ServeMux, map[string]data.PrivateKey) {

	if len(roles) == 0 {
		roles = []string{data.CanonicalTimestampRole, data.CanonicalSnapshotRole}
	}
	keys := make(map[string]data.PrivateKey)
	mux := http.NewServeMux()

	for _, role := range roles {
		key, err := trustmanager.GenerateECDSAKey(rand.Reader)
		require.NoError(t, err)

		keys[role] = key
		pubKey := data.PublicKeyFromPrivate(key)
		jsonBytes, err := json.MarshalCanonical(&pubKey)
		require.NoError(t, err)
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

// initializes a repository in a temporary directory
func initializeRepo(t *testing.T, rootType, gun, url string,
	serverManagesSnapshot bool) (*NotaryRepository, string) {

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory: %s", err)

	serverManagedRoles := []string{}
	if serverManagesSnapshot {
		serverManagedRoles = []string{data.CanonicalSnapshotRole}
	}

	repo, rootPubKeyID := createRepoAndKey(t, rootType, tempBaseDir, gun, url)

	err = repo.Initialize(rootPubKeyID, serverManagedRoles...)
	require.NoError(t, err, "error creating repository: %s", err)
	if err != nil {
		os.RemoveAll(tempBaseDir)
	}

	return repo, rootPubKeyID
}

// Creates a new repository and adds a root key.  Returns the repo and key ID.
func createRepoAndKey(t *testing.T, rootType, tempBaseDir, gun, url string) (
	*NotaryRepository, string) {

	repo, err := NewNotaryRepository(
		tempBaseDir, gun, url, http.DefaultTransport, passphraseRetriever)
	require.NoError(t, err, "error creating repo: %s", err)

	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	require.NoError(t, err, "error generating root key: %s", err)

	return repo, rootPubKey.ID()
}

// creates a new notary repository with the same gun and url as the previous
// repo, so that it can be used to pull the trust data the original repo pushed
func newRepoToTestRepo(t *testing.T, existingRepo *NotaryRepository) *NotaryRepository {
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")

	repo, err := NewNotaryRepository(
		tempBaseDir, existingRepo.gun, existingRepo.baseURL,
		http.DefaultTransport, passphraseRetriever)
	require.NoError(t, err, "error creating repository: %s", err)
	if err != nil {
		defer os.RemoveAll(tempBaseDir)
	}

	return repo
}

// Initializing a new repo while specifying that the server should manage the root
// role will fail.
func TestInitRepositoryManagedRolesIncludingRoot(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, data.CanonicalRootRole)
	require.Error(t, err)
	require.IsType(t, ErrInvalidRemoteRole{}, err)
	// Just testing the error message here in this one case
	require.Equal(t, err.Error(),
		"notary does not support the server managing the root key")
}

// Initializing a new repo while specifying that the server should manage some
// invalid role will fail.
func TestInitRepositoryManagedRolesInvalidRole(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, "randomrole")
	require.Error(t, err)
	require.IsType(t, ErrInvalidRemoteRole{}, err)
}

// Initializing a new repo while specifying that the server should manage the
// targets role will fail.
func TestInitRepositoryManagedRolesIncludingTargets(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", "http://localhost")
	err = repo.Initialize(rootPubKeyID, data.CanonicalTargetsRole)
	require.Error(t, err)
	require.IsType(t, ErrInvalidRemoteRole{}, err)
}

// Initializing a new repo while specifying that the server should manage the
// timestamp key is fine - that's what it already does, so no error.
func TestInitRepositoryManagedRolesIncludingTimestamp(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", ts.URL)
	err = repo.Initialize(rootPubKeyID, data.CanonicalTimestampRole)
	require.NoError(t, err)
}

// Initializing a new repo fails if unable to get the timestamp key, even if
// the snapshot key is available
func TestInitRepositoryNeedsRemoteTimestampKey(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t, data.CanonicalSnapshotRole)
	defer ts.Close()

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", ts.URL)
	err = repo.Initialize(rootPubKeyID, data.CanonicalTimestampRole)
	require.Error(t, err)
	require.IsType(t, store.ErrMetaNotFound{}, err)
}

// Initializing a new repo with remote server signing fails if unable to get
// the snapshot key, even if the timestamp key is available
func TestInitRepositoryNeedsRemoteSnapshotKey(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t, data.CanonicalTimestampRole)
	defer ts.Close()

	repo, rootPubKeyID := createRepoAndKey(
		t, data.ECDSAKey, tempBaseDir, "docker.com/notary", ts.URL)
	err = repo.Initialize(rootPubKeyID, data.CanonicalSnapshotRole)
	require.Error(t, err)
	require.IsType(t, store.ErrMetaNotFound{}, err)
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
func requireRepoHasExpectedKeys(t *testing.T, repo *NotaryRepository,
	rootKeyID string, expectedSnapshotKey bool) {

	// The repo should have a keyFileStore and have created keys using it,
	// so create a new KeyFileStore, and check that the keys do exist and are
	// valid
	ks, err := trustmanager.NewKeyFileStore(repo.baseDir, passphraseRetriever)
	require.NoError(t, err)

	roles := make(map[string]bool)
	for keyID, role := range ks.ListKeys() {
		if role == data.CanonicalRootRole {
			require.Equal(t, rootKeyID, keyID, "Unexpected root key ID")
		}
		// just to ensure the content of the key files created are valid
		_, r, err := ks.GetKey(keyID)
		require.NoError(t, err)
		require.Equal(t, role, r)
		roles[role] = true
	}
	// there is a root key and a targets key
	alwaysThere := []string{data.CanonicalRootRole, data.CanonicalTargetsRole}
	for _, role := range alwaysThere {
		_, ok := roles[role]
		require.True(t, ok, "missing %s key", role)
	}

	// there may be a snapshots key, depending on whether the server is managing
	// the snapshots key
	_, ok := roles[data.CanonicalSnapshotRole]
	if expectedSnapshotKey {
		require.True(t, ok, "missing snapshot key")
	} else {
		require.False(t, ok,
			"there should be no snapshot key because the server manages it")
	}

	// The server manages the timestamp key - there should not be a timestamp
	// key
	_, ok = roles[data.CanonicalTimestampRole]
	require.False(t, ok,
		"there should be no timestamp key because the server manages it")
}

// This creates a new certificate manager in the repo's base directory and
// makes sure the repo has the right certificates
func requireRepoHasExpectedCerts(t *testing.T, repo *NotaryRepository) {
	// The repo should have a certificate manager and have created certs using
	// it, so create a new manager, and check that the certs do exist and
	// are valid
	certManager, err := certs.NewManager(repo.baseDir)
	require.NoError(t, err)
	certificates := certManager.TrustedCertificateStore().GetCertificates()
	require.Len(t, certificates, 1, "unexpected number of trusted certificates")

	certID, err := trustmanager.FingerprintCert(certificates[0])
	require.NoError(t, err, "unable to fingerprint the trusted certificate")
	require.NotEqual(t, certID, "")
}

// Sanity check the TUF metadata files. Verify that it exists for a particular
// role, the JSON is well-formed, and the signatures exist.
// For the root.json file, also check that the root, snapshot, and
// targets key IDs are present.
func requireRepoHasExpectedMetadata(t *testing.T, repo *NotaryRepository,
	role string, expected bool) {

	filename := filepath.Join(tufDir, filepath.FromSlash(repo.gun),
		"metadata", role+".json")
	fullPath := filepath.Join(repo.baseDir, filename)
	_, err := os.Stat(fullPath)

	if expected {
		require.NoError(t, err, "missing TUF metadata file: %s", filename)
	} else {
		require.Error(t, err,
			"%s metadata should not exist, but does: %s", role, filename)
		return
	}

	jsonBytes, err := ioutil.ReadFile(fullPath)
	require.NoError(t, err, "error reading TUF metadata file %s: %s", filename, err)

	var decoded data.Signed
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err, "error parsing TUF metadata file %s: %s", filename, err)

	require.Len(t, decoded.Signatures, 1,
		"incorrect number of signatures in TUF metadata file %s", filename)

	require.NotEmpty(t, decoded.Signatures[0].KeyID,
		"empty key ID field in TUF metadata file %s", filename)
	require.NotEmpty(t, decoded.Signatures[0].Method,
		"empty method field in TUF metadata file %s", filename)
	require.NotEmpty(t, decoded.Signatures[0].Signature,
		"empty signature in TUF metadata file %s", filename)

	// Special case for root.json: also check that the signed
	// content for keys and roles
	if role == data.CanonicalRootRole {
		var decodedRoot data.Root
		err := json.Unmarshal(decoded.Signed, &decodedRoot)
		require.NoError(t, err, "error parsing root.json signed section: %s", err)

		require.Equal(t, "Root", decodedRoot.Type, "_type mismatch in root.json")

		// Expect 1 key for each valid role in the Keys map - one for
		// each of root, targets, snapshot, timestamp
		require.Len(t, decodedRoot.Keys, len(data.ValidRoles),
			"wrong number of keys in root.json")
		require.Len(t, decodedRoot.Roles, len(data.ValidRoles),
			"wrong number of roles in root.json")

		for role := range data.ValidRoles {
			_, ok := decodedRoot.Roles[role]
			require.True(t, ok, "Missing role %s in root.json", role)
		}
	}
}

func testInitRepo(t *testing.T, rootType string, serverManagesSnapshot bool) {
	gun := "docker.com/notary"

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootKeyID := initializeRepo(t, rootType, gun, ts.URL, serverManagesSnapshot)
	defer os.RemoveAll(repo.baseDir)

	requireRepoHasExpectedKeys(t, repo, rootKeyID, !serverManagesSnapshot)
	requireRepoHasExpectedCerts(t, repo)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalRootRole, true)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTargetsRole, true)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole,
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
	require.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	retriever := passphrase.ConstantRetriever("password")
	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	require.NoError(t, err, "error creating repo: %s", err)
	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	require.NoError(t, err, "error generating root key: %s", err)

	retriever = passphrase.ConstantRetriever("incorrect password")
	// repo.CryptoService’s FileKeyStore caches the unlocked private key, so to test
	// private key unlocking we need a new repo instance.
	repo, err = NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	require.NoError(t, err, "error creating repo: %s", err)
	err = repo.Initialize(rootPubKey.ID())
	require.EqualError(t, err, trustmanager.ErrAttemptsExceeded{}.Error())
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
	require.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	retriever := passphrase.ConstantRetriever("password")
	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, retriever)
	require.NoError(t, err, "error creating repo: %s", err)
	rootPubKey, err := repo.CryptoService.Create("root", rootType)
	require.NoError(t, err, "error generating root key: %s", err)

	// repo.CryptoService’s FileKeyStore caches the unlocked private key, so to test
	// private key unlocking we need a new repo instance.
	repo, err = NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, giveUpPassphraseRetriever)
	require.NoError(t, err, "error creating repo: %s", err)
	err = repo.Initialize(rootPubKey.ID())
	require.EqualError(t, err, trustmanager.ErrPasswordInvalid{}.Error())
}

func addTarget(t *testing.T, repo *NotaryRepository, targetName, targetFile string,
	roles ...string) *Target {
	target, err := NewTarget(targetName, targetFile)
	require.NoError(t, err, "error creating target")
	err = repo.AddTarget(target, roles...)
	require.NoError(t, err, "error adding target")
	return target
}

// calls GetChangelist and gets the actual changes out
func getChanges(t *testing.T, repo *NotaryRepository) []changelist.Change {
	changeList, err := repo.GetChangelist()
	require.NoError(t, err)
	return changeList.List()
}

// TestAddTargetToTargetRoleByDefault adds a target without specifying a role
// to a repo without delegations.  Confirms that the changelist is created
// correctly, for the targets scope.
func TestAddTargetToTargetRoleByDefault(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	testAddOrDeleteTarget(t, repo, changelist.ActionCreate, nil,
		[]string{data.CanonicalTargetsRole})
}

// Tests that adding a target to a repo or deleting a target from a repo,
// with the given roles, makes a change to the expected scopes
func testAddOrDeleteTarget(t *testing.T, repo *NotaryRepository, action string,
	rolesToChange []string, expectedScopes []string) {

	require.Len(t, getChanges(t, repo), 0, "should start with zero changes")

	if action == changelist.ActionCreate {
		// Add fixtures/intermediate-ca.crt as a target. There's no particular
		// reason for using this file except that it happens to be available as
		// a fixture.
		addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt", rolesToChange...)
	} else {
		err := repo.RemoveTarget("latest", rolesToChange...)
		require.NoError(t, err, "error removing target")
	}

	changes := getChanges(t, repo)
	require.Len(t, changes, len(expectedScopes), "wrong number of changes files found")

	foundScopes := make(map[string]bool)
	for _, c := range changes { // there is only one
		require.EqualValues(t, action, c.Action())
		foundScopes[c.Scope()] = true
		require.Equal(t, "target", c.Type())
		require.Equal(t, "latest", c.Path())
		if action == changelist.ActionCreate {
			require.NotEmpty(t, c.Content())
		} else {
			require.Empty(t, c.Content())
		}
	}
	require.Len(t, foundScopes, len(expectedScopes))
	for _, expectedScope := range expectedScopes {
		_, ok := foundScopes[expectedScope]
		require.True(t, ok, "Target was not added/removed from %s", expectedScope)
	}

	// add/delete a second time
	if action == changelist.ActionCreate {
		addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt", rolesToChange...)
	} else {
		err := repo.RemoveTarget("current", rolesToChange...)
		require.NoError(t, err, "error removing target")
	}

	changes = getChanges(t, repo)
	require.Len(t, changes, 2*len(expectedScopes),
		"wrong number of changelist files found")

	newFileFound := false
	foundScopes = make(map[string]bool)
	for _, c := range changes {
		if c.Path() != "latest" {
			require.EqualValues(t, action, c.Action())
			foundScopes[c.Scope()] = true
			require.Equal(t, "target", c.Type())
			require.Equal(t, "current", c.Path())
			if action == changelist.ActionCreate {
				require.NotEmpty(t, c.Content())
			} else {
				require.Empty(t, c.Content())
			}

			newFileFound = true
		}
	}
	require.True(t, newFileFound, "second changelist file not found")
	require.Len(t, foundScopes, len(expectedScopes))
	for _, expectedScope := range expectedScopes {
		_, ok := foundScopes[expectedScope]
		require.True(t, ok, "Target was not added/removed from %s", expectedScope)
	}
}

// TestAddTargetToSpecifiedValidRoles adds a target to the specified roles.
// Confirms that the changelist is created correctly, one for each of the
// the specified roles as scopes.
func TestAddTargetToSpecifiedValidRoles(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	roleName := filepath.Join(data.CanonicalTargetsRole, "a")
	testAddOrDeleteTarget(t, repo, changelist.ActionCreate,
		[]string{
			strings.ToUpper(data.CanonicalTargetsRole),
			strings.ToUpper(roleName),
		},
		[]string{data.CanonicalTargetsRole, roleName})
}

// TestAddTargetToSpecifiedInvalidRoles expects errors to be returned if
// adding a target to an invalid role.  If any of the roles are invalid,
// no targets are added to any roles.
func TestAddTargetToSpecifiedInvalidRoles(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	invalidRoles := []string{
		data.CanonicalRootRole,
		data.CanonicalSnapshotRole,
		data.CanonicalTimestampRole,
		"target/otherrole",
		"otherrole",
	}

	for _, invalidRole := range invalidRoles {
		target, err := NewTarget("latest", "../fixtures/intermediate-ca.crt")
		require.NoError(t, err, "error creating target")

		err = repo.AddTarget(target, data.CanonicalTargetsRole, invalidRole)
		require.Error(t, err, "Expected an ErrInvalidRole error")
		require.IsType(t, data.ErrInvalidRole{}, err)

		changes := getChanges(t, repo)
		require.Len(t, changes, 0)
	}
}

// General way to require that errors writing a changefile are propagated up
func testErrorWritingChangefiles(t *testing.T, writeChangeFile func(*NotaryRepository) error) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// first, make the actual changefile unwritable by making the changelist
	// directory unwritable
	changelistPath := filepath.Join(repo.tufRepoPath, "changelist")
	err := os.MkdirAll(changelistPath, 0744)
	require.NoError(t, err, "could not create changelist dir")
	err = os.Chmod(changelistPath, 0600)
	require.NoError(t, err, "could not change permission of changelist dir")

	err = writeChangeFile(repo)
	require.Error(t, err, "Expected an error writing the change")
	require.IsType(t, &os.PathError{}, err)

	// then break prevent the changlist directory from being able to be created
	err = os.Chmod(changelistPath, 0744)
	require.NoError(t, err, "could not change permission of temp dir")
	err = os.RemoveAll(changelistPath)
	require.NoError(t, err, "could not remove changelist dir")
	// creating a changelist file so the directory can't be created
	err = ioutil.WriteFile(changelistPath, []byte("hi"), 0644)
	require.NoError(t, err, "could not write temporary file")

	err = writeChangeFile(repo)
	require.Error(t, err, "Expected an error writing the change")
	require.IsType(t, &os.PathError{}, err)
}

// TestAddTargetErrorWritingChanges expects errors writing a change to file
// to be propagated.
func TestAddTargetErrorWritingChanges(t *testing.T) {
	testErrorWritingChangefiles(t, func(repo *NotaryRepository) error {
		target, err := NewTarget("latest", "../fixtures/intermediate-ca.crt")
		require.NoError(t, err, "error creating target")
		return repo.AddTarget(target, data.CanonicalTargetsRole)
	})
}

// TestRemoveTargetToTargetRoleByDefault removes a target without specifying a
// role from a repo.  Confirms that the changelist is created correctly for
// the targets scope.
func TestRemoveTargetToTargetRoleByDefault(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	testAddOrDeleteTarget(t, repo, changelist.ActionDelete, nil,
		[]string{data.CanonicalTargetsRole})
}

// TestRemoveTargetFromSpecifiedValidRoles removes a target from the specified
// roles. Confirms that the changelist is created correctly, one for each of
// the the specified roles as scopes.
func TestRemoveTargetFromSpecifiedValidRoles(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	roleName := filepath.Join(data.CanonicalTargetsRole, "a")
	testAddOrDeleteTarget(t, repo, changelist.ActionDelete,
		[]string{
			strings.ToUpper(data.CanonicalTargetsRole),
			strings.ToUpper(roleName),
		},
		[]string{data.CanonicalTargetsRole, roleName})
}

// TestRemoveTargetFromSpecifiedInvalidRoles expects errors to be returned if
// removing a target to an invalid role.  If any of the roles are invalid,
// no targets are removed from any roles.
func TestRemoveTargetToSpecifiedInvalidRoles(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	invalidRoles := []string{
		data.CanonicalRootRole,
		data.CanonicalSnapshotRole,
		data.CanonicalTimestampRole,
		"target/otherrole",
		"otherrole",
	}

	for _, invalidRole := range invalidRoles {
		err := repo.RemoveTarget("latest", data.CanonicalTargetsRole, invalidRole)
		require.Error(t, err, "Expected an ErrInvalidRole error")
		require.IsType(t, data.ErrInvalidRole{}, err)

		changes := getChanges(t, repo)
		require.Len(t, changes, 0)
	}
}

// TestRemoveTargetErrorWritingChanges expects errors writing a change to file
// to be propagated.
func TestRemoveTargetErrorWritingChanges(t *testing.T) {
	testErrorWritingChangefiles(t, func(repo *NotaryRepository) error {
		return repo.RemoveTarget("latest", data.CanonicalTargetsRole)
	})
}

// TestListTarget fakes serving signed metadata files over the test's
// internal HTTP server to ensure that ListTargets returns the correct number
// of listed targets.
// We test this with both an RSA and ECDSA root key
func TestListTarget(t *testing.T) {
	testListEmptyTargets(t, data.ECDSAKey)
	testListTarget(t, data.ECDSAKey)
	testListTargetWithDelegates(t, data.ECDSAKey)
	if !testing.Short() {
		testListEmptyTargets(t, data.RSAKey)
		testListTarget(t, data.RSAKey)
		testListTargetWithDelegates(t, data.RSAKey)
	}
}

func testListEmptyTargets(t *testing.T, rootType string) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	_, err := repo.ListTargets(data.CanonicalTargetsRole)
	require.Error(t, err) // no trust data
}

// reads data from the repository in order to fake data being served via
// the ServeMux.
func fakeServerData(t *testing.T, repo *NotaryRepository, mux *http.ServeMux,
	keys map[string]data.PrivateKey) {

	timestampKey, ok := keys[data.CanonicalTimestampRole]
	require.True(t, ok)
	savedTUFRepo := repo.tufRepo // in case this is overwritten

	fileStore, err := trustmanager.NewKeyFileStore(repo.baseDir, passphraseRetriever)
	require.NoError(t, err)
	fileStore.AddKey(
		filepath.Join(filepath.FromSlash(repo.gun), timestampKey.ID()),
		"nonroot", timestampKey)

	rootJSONFile := filepath.Join(repo.baseDir, "tuf",
		filepath.FromSlash(repo.gun), "metadata", "root.json")
	rootFileBytes, err := ioutil.ReadFile(rootJSONFile)

	signedTargets, err := savedTUFRepo.SignTargets(
		"targets", data.DefaultExpires("targets"))
	require.NoError(t, err)

	signedLevel1, err := savedTUFRepo.SignTargets(
		"targets/level1",
		data.DefaultExpires(data.CanonicalTargetsRole),
	)
	if _, ok := savedTUFRepo.Targets["targets/level1"]; ok {
		require.NoError(t, err)
	}

	signedLevel2, err := savedTUFRepo.SignTargets(
		"targets/level2",
		data.DefaultExpires(data.CanonicalTargetsRole),
	)
	if _, ok := savedTUFRepo.Targets["targets/level2"]; ok {
		require.NoError(t, err)
	}

	signedSnapshot, err := savedTUFRepo.SignSnapshot(
		data.DefaultExpires("snapshot"))
	require.NoError(t, err)

	signedTimestamp, err := savedTUFRepo.SignTimestamp(
		data.DefaultExpires("timestamp"))
	require.NoError(t, err)

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/root.json",
		func(w http.ResponseWriter, r *http.Request) {
			require.NoError(t, err)
			fmt.Fprint(w, string(rootFileBytes))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.json",
		func(w http.ResponseWriter, r *http.Request) {
			timestampJSON, _ := json.MarshalCanonical(signedTimestamp)
			fmt.Fprint(w, string(timestampJSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/snapshot.json",
		func(w http.ResponseWriter, r *http.Request) {
			snapshotJSON, _ := json.MarshalCanonical(signedSnapshot)
			fmt.Fprint(w, string(snapshotJSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets.json",
		func(w http.ResponseWriter, r *http.Request) {
			targetsJSON, _ := json.MarshalCanonical(signedTargets)
			fmt.Fprint(w, string(targetsJSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets/level1.json",
		func(w http.ResponseWriter, r *http.Request) {
			level1JSON, err := json.MarshalCanonical(signedLevel1)
			require.NoError(t, err)
			fmt.Fprint(w, string(level1JSON))
		})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets/level2.json",
		func(w http.ResponseWriter, r *http.Request) {
			level2JSON, err := json.MarshalCanonical(signedLevel2)
			require.NoError(t, err)
			fmt.Fprint(w, string(level2JSON))
		})
}

// We want to sort by name, so we can guarantee ordering.
type targetSorter []*TargetWithRole

func (k targetSorter) Len() int           { return len(k) }
func (k targetSorter) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k targetSorter) Less(i, j int) bool { return k[i].Name < k[j].Name }

func testListTarget(t *testing.T, rootType string) {
	ts, mux, keys := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err := repo.tufRepo.InitTimestamp()
	require.NoError(t, err, "error creating repository: %s", err)

	latestTarget := addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	currentTarget := addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")

	// Apply the changelist. Normally, this would be done by Publish

	// load the changelist for this repo
	cl, err := changelist.NewFileChangelist(
		filepath.Join(repo.baseDir, "tuf", filepath.FromSlash(repo.gun), "changelist"))
	require.NoError(t, err, "could not open changelist")

	// apply the changelist to the repo
	err = applyChangelist(repo.tufRepo, cl)
	require.NoError(t, err, "could not apply changelist")

	fakeServerData(t, repo, mux, keys)

	targets, err := repo.ListTargets(data.CanonicalTargetsRole)
	require.NoError(t, err)

	// Should be two targets
	require.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

	sort.Stable(targetSorter(targets))

	// the targets should both be found in the targets role
	for _, foundTarget := range targets {
		require.Equal(t, data.CanonicalTargetsRole, foundTarget.Role)
	}

	// current should be first
	require.True(t, reflect.DeepEqual(*currentTarget, targets[0].Target), "current target does not match")
	require.True(t, reflect.DeepEqual(*latestTarget, targets[1].Target), "latest target does not match")

	// Also test GetTargetByName
	newLatestTarget, err := repo.GetTargetByName("latest")
	require.NoError(t, err)
	require.Equal(t, data.CanonicalTargetsRole, newLatestTarget.Role)
	require.True(t, reflect.DeepEqual(*latestTarget, newLatestTarget.Target), "latest target does not match")

	newCurrentTarget, err := repo.GetTargetByName("current")
	require.NoError(t, err)
	require.Equal(t, data.CanonicalTargetsRole, newCurrentTarget.Role)
	require.True(t, reflect.DeepEqual(*currentTarget, newCurrentTarget.Target), "current target does not match")
}

func testListTargetWithDelegates(t *testing.T, rootType string) {
	ts, mux, keys := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err := repo.tufRepo.InitTimestamp()
	require.NoError(t, err, "error creating repository: %s", err)

	latestTarget := addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	currentTarget := addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")

	// setup delegated targets/level1 role
	k, err := repo.CryptoService.Create("targets/level1", rootType)
	require.NoError(t, err)
	r, err := data.NewRole("targets/level1", 1, []string{k.ID()}, []string{""}, nil)
	require.NoError(t, err)
	repo.tufRepo.UpdateDelegations(r, []data.PublicKey{k})
	delegatedTarget := addTarget(t, repo, "current", "../fixtures/root-ca.crt", "targets/level1")
	otherTarget := addTarget(t, repo, "other", "../fixtures/root-ca.crt", "targets/level1")

	// setup delegated targets/level2 role
	k, err = repo.CryptoService.Create("targets/level2", rootType)
	require.NoError(t, err)
	r, err = data.NewRole("targets/level2", 1, []string{k.ID()}, []string{""}, nil)
	require.NoError(t, err)
	repo.tufRepo.UpdateDelegations(r, []data.PublicKey{k})
	// this target should not show up as the one in targets/level1 takes higher priority
	_ = addTarget(t, repo, "current", "../fixtures/notary-server.crt", "targets/level2")
	// this target should show up as the name doesn't exist elsewhere
	level2Target := addTarget(t, repo, "level2", "../fixtures/notary-server.crt", "targets/level2")

	// Apply the changelist. Normally, this would be done by Publish

	// load the changelist for this repo
	cl, err := changelist.NewFileChangelist(
		filepath.Join(repo.baseDir, "tuf", filepath.FromSlash(repo.gun), "changelist"))
	require.NoError(t, err, "could not open changelist")

	// apply the changelist to the repo
	err = applyChangelist(repo.tufRepo, cl)
	require.NoError(t, err, "could not apply changelist")

	_, ok := repo.tufRepo.Targets["targets/level1"].Signed.Targets["current"]
	require.True(t, ok)
	_, ok = repo.tufRepo.Targets["targets/level1"].Signed.Targets["other"]
	require.True(t, ok)
	_, ok = repo.tufRepo.Targets["targets/level2"].Signed.Targets["level2"]
	require.True(t, ok)

	fakeServerData(t, repo, mux, keys)

	// test default listing
	targets, err := repo.ListTargets()
	require.NoError(t, err)

	// Should be four targets
	require.Len(t, targets, 4, "unexpected number of targets returned by ListTargets")

	sort.Stable(targetSorter(targets))

	// current should be first.
	require.True(t, reflect.DeepEqual(*currentTarget, targets[0].Target), "current target does not match")
	require.Equal(t, data.CanonicalTargetsRole, targets[0].Role)

	require.True(t, reflect.DeepEqual(*latestTarget, targets[1].Target), "latest target does not match")
	require.Equal(t, data.CanonicalTargetsRole, targets[1].Role)

	require.True(t, reflect.DeepEqual(*level2Target, targets[2].Target), "level2 target does not match")
	require.Equal(t, "targets/level2", targets[2].Role)

	require.True(t, reflect.DeepEqual(*otherTarget, targets[3].Target), "other target does not match")
	require.Equal(t, "targets/level1", targets[3].Role)

	// test listing with priority specified
	targets, err = repo.ListTargets("targets/level1", data.CanonicalTargetsRole)
	require.NoError(t, err)

	// Should be four targets
	require.Len(t, targets, 4, "unexpected number of targets returned by ListTargets")

	sort.Stable(targetSorter(targets))

	// current (in delegated role) should be first
	require.True(t, reflect.DeepEqual(*delegatedTarget, targets[0].Target), "current target does not match")
	require.Equal(t, "targets/level1", targets[0].Role)

	require.True(t, reflect.DeepEqual(*latestTarget, targets[1].Target), "latest target does not match")
	require.Equal(t, data.CanonicalTargetsRole, targets[1].Role)

	require.True(t, reflect.DeepEqual(*level2Target, targets[2].Target), "level2 target does not match")
	require.Equal(t, "targets/level2", targets[2].Role)

	require.True(t, reflect.DeepEqual(*otherTarget, targets[3].Target), "other target does not match")
	require.Equal(t, "targets/level1", targets[3].Role)

	// Also test GetTargetByName
	newLatestTarget, err := repo.GetTargetByName("latest")
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(*latestTarget, newLatestTarget.Target), "latest target does not match")
	require.Equal(t, data.CanonicalTargetsRole, newLatestTarget.Role)

	newCurrentTarget, err := repo.GetTargetByName("current", "targets/level1", "targets")
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(*delegatedTarget, newCurrentTarget.Target), "current target does not match")
	require.Equal(t, "targets/level1", newCurrentTarget.Role)

	newOtherTarget, err := repo.GetTargetByName("other")
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(*otherTarget, newOtherTarget.Target), "other target does not match")
	require.Equal(t, "targets/level1", newOtherTarget.Role)

	newLevel2Target, err := repo.GetTargetByName("level2")
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(*level2Target, newLevel2Target.Target), "level2 target does not match")
	require.Equal(t, "targets/level2", newLevel2Target.Role)
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
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	rootJSONFile := filepath.Join(repo.baseDir, "tuf", filepath.FromSlash(repo.gun),
		"metadata", "root.json")

	jsonBytes, err := ioutil.ReadFile(rootJSONFile)
	require.NoError(t, err, "error reading TUF metadata file %s: %s", rootJSONFile, err)

	var decoded data.Signed
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err, "error parsing TUF metadata file %s: %s", rootJSONFile, err)

	var decodedRoot data.Root
	err = json.Unmarshal(decoded.Signed, &decodedRoot)
	require.NoError(t, err, "error parsing root.json signed section: %s", err)

	keyids := []string{}
	for role, roleData := range decodedRoot.Roles {
		if role == "root" {
			keyids = append(keyids, roleData.KeyIDs...)
		}
	}
	require.NotEmpty(t, keyids)

	for _, keyid := range keyids {
		key, ok := decodedRoot.Keys[keyid]
		require.True(t, ok, "key id not found in keys")
		_, err := trustmanager.LoadCertFromPEM(key.Public())
		require.NoError(t, err, "key is not a valid cert")
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
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	require.Len(t, getChanges(t, repo), 0, "No changes should be in changelist yet")

	// Create 2 targets
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	addTarget(t, repo, "current", "../fixtures/intermediate-ca.crt")

	// Test loading changelist
	chgs := getChanges(t, repo)
	require.Len(t, chgs, 2, "Wrong number of changes returned from changelist")

	changes := make(map[string]changelist.Change)
	for _, ch := range chgs {
		changes[ch.Path()] = ch
	}

	currentChange := changes["current"]
	require.NotNil(t, currentChange, "Expected changelist to contain a change for path 'current'")
	require.EqualValues(t, changelist.ActionCreate, currentChange.Action())
	require.Equal(t, "targets", currentChange.Scope())
	require.Equal(t, "target", currentChange.Type())
	require.Equal(t, "current", currentChange.Path())

	latestChange := changes["latest"]
	require.NotNil(t, latestChange, "Expected changelist to contain a change for path 'latest'")
	require.EqualValues(t, changelist.ActionCreate, latestChange.Action())
	require.Equal(t, "targets", latestChange.Scope())
	require.Equal(t, "target", latestChange.Type())
	require.Equal(t, "latest", latestChange.Path())
}

// Create a repo, instantiate a notary server, and publish the bare repo to the
// server, signing all the non-timestamp metadata.  Root, targets, and snapshots
// (if locally signing) should be sent.
func TestPublishBareRepo(t *testing.T) {
	testPublishNoData(t, data.ECDSAKey, true)
	testPublishNoData(t, data.ECDSAKey, false)
	if !testing.Short() {
		testPublishNoData(t, data.RSAKey, true)
		testPublishNoData(t, data.RSAKey, false)
	}
}

func testPublishNoData(t *testing.T, rootType string, serverManagesSnapshot bool) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo1, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL,
		serverManagesSnapshot)
	defer os.RemoveAll(repo1.baseDir)
	require.NoError(t, repo1.Publish())

	// use another repo to check metadata
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	targets, err := repo2.ListTargets()
	require.NoError(t, err)
	require.Empty(t, targets)

	for role := range data.ValidRoles {
		// we don't cache timstamp metadata
		if role != data.CanonicalTimestampRole {
			requireRepoHasExpectedMetadata(t, repo2, role, true)
		}
	}
}

// Publishing an uninitialized repo will fail, but initializing and republishing
// after should succeed
func TestPublishUninitializedRepo(t *testing.T) {
	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	// uninitialized repo should fail to publish
	tempBaseDir, err := ioutil.TempDir("", "notary-tests")
	require.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL,
		http.DefaultTransport, passphraseRetriever)
	require.NoError(t, err, "error creating repository: %s", err)
	err = repo.Publish()
	require.Error(t, err)

	// no metadata created
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalRootRole, false)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, false)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTargetsRole, false)

	// now, initialize and republish in the same directory
	rootPubKey, err := repo.CryptoService.Create("root", data.ECDSAKey)
	require.NoError(t, err, "error generating root key: %s", err)

	require.NoError(t, repo.Initialize(rootPubKey.ID()))

	// now metadata is created
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalRootRole, true)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, true)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTargetsRole, true)

	require.NoError(t, repo.Publish())
}

// Create a repo, instantiate a notary server, and publish the repo with
// some targets to the server, signing all the non-timestamp metadata.
// We test this with both an RSA and ECDSA root key
func TestPublishClientHasSnapshotKey(t *testing.T) {
	testPublishWithData(t, data.ECDSAKey, false)
	if !testing.Short() {
		testPublishWithData(t, data.RSAKey, false)
	}
}

// Create a repo, instantiate a notary server (designating the server as the
// snapshot signer) , and publish the repo with some targets to the server,
// signing the root and targets metadata only.  The server should sign just fine.
// We test this with both an RSA and ECDSA root key
func TestPublishAfterInitServerHasSnapshotKey(t *testing.T) {
	testPublishWithData(t, data.ECDSAKey, true)
	if !testing.Short() {
		testPublishWithData(t, data.RSAKey, true)
	}
}

func testPublishWithData(t *testing.T, rootType string, serverManagesSnapshot bool) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL,
		serverManagesSnapshot)
	defer os.RemoveAll(repo.baseDir)
	requirePublishSucceeds(t, repo)
}

// requires that publish succeeds by adding to the default only and publishing;
// the targets should appear in targets
func requirePublishSucceeds(t *testing.T, repo1 *NotaryRepository) {
	requirePublishToRolesSucceeds(t, repo1, nil, []string{data.CanonicalTargetsRole})
}

// requires that adding to the given roles results in the targets actually being
// added only to the expected roles and no others
func requirePublishToRolesSucceeds(t *testing.T, repo1 *NotaryRepository,
	publishToRoles []string, expectedPublishedRoles []string) {

	// were there unpublished changes before?
	changesOffset := len(getChanges(t, repo1))

	// Create 2 targets - (actually 3, but we delete 1)
	addTarget(t, repo1, "toDelete", "../fixtures/intermediate-ca.crt", publishToRoles...)
	latestTarget := addTarget(
		t, repo1, "latest", "../fixtures/intermediate-ca.crt", publishToRoles...)
	currentTarget := addTarget(
		t, repo1, "current", "../fixtures/intermediate-ca.crt", publishToRoles...)
	repo1.RemoveTarget("toDelete", publishToRoles...)

	// if no roles are provided, the default role is target
	numRoles := int(math.Max(1, float64(len(publishToRoles))))
	require.Len(t, getChanges(t, repo1), changesOffset+4*numRoles,
		"wrong number of changelist files found")

	// Now test Publish
	err := repo1.Publish()
	require.NoError(t, err)
	require.Len(t, getChanges(t, repo1), 0, "wrong number of changelist files found")

	// use another repo to check metadata
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	// Should be two targets per role
	for _, role := range expectedPublishedRoles {
		for _, repo := range []*NotaryRepository{repo1, repo2} {
			targets, err := repo.ListTargets(role)
			require.NoError(t, err)

			require.Len(t, targets, 2,
				"unexpected number of targets returned by ListTargets(%s)", role)

			sort.Stable(targetSorter(targets))

			require.True(t, reflect.DeepEqual(*currentTarget, targets[0].Target), "current target does not match")
			require.Equal(t, role, targets[0].Role)
			require.True(t, reflect.DeepEqual(*latestTarget, targets[1].Target), "latest target does not match")
			require.Equal(t, role, targets[1].Role)

			// Also test GetTargetByName
			newLatestTarget, err := repo.GetTargetByName("latest", role)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(*latestTarget, newLatestTarget.Target), "latest target does not match")
			require.Equal(t, role, newLatestTarget.Role)

			newCurrentTarget, err := repo.GetTargetByName("current", role)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(*currentTarget, newCurrentTarget.Target), "current target does not match")
			require.Equal(t, role, newCurrentTarget.Role)
		}
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
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, true)
	defer os.RemoveAll(repo.baseDir)
	// no timestamp metadata because that comes from the server
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, false)
	// no snapshot metadata because that comes from the server
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, false)

	// Publish something
	published := addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")
	require.NoError(t, repo.Publish())

	// still no timestamp or snapshot metadata info
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, false)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, false)

	// list, so that the snapshot metadata is pulled from server
	targets, err := repo.ListTargets(data.CanonicalTargetsRole)
	require.NoError(t, err)
	require.Equal(t, []*TargetWithRole{{Target: *published, Role: data.CanonicalTargetsRole}}, targets)
	// listing downloaded the timestamp and snapshot metadata info
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalTimestampRole, true)
	requireRepoHasExpectedMetadata(t, repo, data.CanonicalSnapshotRole, true)

	// Publish again should succeed
	addTarget(t, repo, "v2", "../fixtures/intermediate-ca.crt")
	err = repo.Publish()
	require.NoError(t, err)
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
	ts := fullTestServer(t)
	defer ts.Close()

	// create repo and delete the snapshot key and metadata
	repo, _ := initializeRepo(t, rootType, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	snapshotRole, ok := repo.tufRepo.Root.Signed.Roles[data.CanonicalSnapshotRole]
	require.True(t, ok)
	for _, keyID := range snapshotRole.KeyIDs {
		repo.CryptoService.RemoveKey(keyID)
	}

	// Publish something
	addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")
	err := repo.Publish()
	require.Error(t, err)
	require.IsType(t, validation.ErrBadHierarchy{}, err)
}

// If the snapshot metadata is corrupt or the snapshot metadata is unreadable,
// we can't publish for the first time (whether the client or server has the
// snapshot key), because there is no existing data for us to download. If the
// repo has already been published, it doesn't matter if the metadata is corrupt
// because we can just redownload if it is.
func TestPublishSnapshotCorrupt(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	// do not publish first - publish should fail with corrupt snapshot data even with server signing snapshot
	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary1", ts.URL, true)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalSnapshotRole, repo, false, false)

	// do not publish first - publish should fail with corrupt snapshot data with local snapshot signing
	repo, _ = initializeRepo(t, data.ECDSAKey, "docker.com/notary2", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalSnapshotRole, repo, false, false)

	// publish first - publish again should succeed despite corrupt snapshot data (server signing snapshot)
	repo, _ = initializeRepo(t, data.ECDSAKey, "docker.com/notary3", ts.URL, true)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalSnapshotRole, repo, true, true)

	// publish first - publish again should succeed despite corrupt snapshot data (local snapshot signing)
	repo, _ = initializeRepo(t, data.ECDSAKey, "docker.com/notary4", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalSnapshotRole, repo, true, true)
}

// If the targets metadata is corrupt or the targets metadata is unreadable,
// we can't publish for the first time, because there is no existing data for.
// us to download. If the repo has already been published, it doesn't matter
// if the metadata is corrupt because we can just redownload if it is.
func TestPublishTargetsCorrupt(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	// do not publish first - publish should fail with corrupt snapshot data
	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary1", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalTargetsRole, repo, false, false)

	// publish first - publish again should succeed despite corrupt snapshot data
	repo, _ = initializeRepo(t, data.ECDSAKey, "docker.com/notary2", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalTargetsRole, repo, true, true)
}

// If the root metadata is corrupt or the root metadata is unreadable,
// we can't publish for the first time.  If there is already a remote root,
// we just download that and verify (using our trusted certificate trust
// anchors) that it is signed with the same keys, and if so, we just use the
// remote root.
func TestPublishRootCorrupt(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	// do not publish first - publish should fail with corrupt snapshot data
	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary1", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalRootRole, repo, false, false)

	// publish first - publish should still succeed if root corrupt since the
	// remote root is signed with the same key.
	repo, _ = initializeRepo(t, data.ECDSAKey, "docker.com/notary2", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	testPublishBadMetadata(t, data.CanonicalRootRole, repo, true, true)
}

// When publishing snapshot, root, or target, if the repo hasn't been published
// before, if the metadata is corrupt, it can't be published.  If it has been
// published already, then the corrupt metadata can just be re-downloaded, so
// publishing is successful.
func testPublishBadMetadata(t *testing.T, roleName string, repo *NotaryRepository,
	publishFirst, succeeds bool) {

	if publishFirst {
		require.NoError(t, repo.Publish())
	}

	addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")

	// readable, but corrupt file
	repo.fileStore.SetMeta(roleName, []byte("this isn't JSON"))
	err := repo.Publish()
	if succeeds {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.IsType(t, &json.SyntaxError{}, err)
	}

	// make an unreadable file by creating a directory instead of a file
	path := fmt.Sprintf("%s.%s",
		filepath.Join(repo.baseDir, tufDir, filepath.FromSlash(repo.gun),
			"metadata", roleName), "json")
	os.RemoveAll(path)
	require.NoError(t, os.Mkdir(path, 0755))
	defer os.RemoveAll(path)

	err = repo.Publish()
	if succeeds {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.IsType(t, &os.PathError{}, err)
	}
}

// If the repo is not initialized, calling repo.Publish() should return ErrRepoNotInitialized
func TestNotInitializedOnPublish(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	require.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := createRepoAndKey(t, data.ECDSAKey, tempBaseDir, gun, ts.URL)

	addTarget(t, repo, "v1", "../fixtures/intermediate-ca.crt")

	err = repo.Publish()
	require.Error(t, err)
	require.IsType(t, ErrRepoNotInitialized{}, err)
}

type cannotCreateKeys struct {
	signed.CryptoService
}

func (cs cannotCreateKeys) Create(_, _ string) (data.PublicKey, error) {
	return nil, fmt.Errorf("Oh no I cannot create keys")
}

// If there is an error creating the local keys, no call is made to get a
// remote key.
func TestPublishSnapshotLocalKeysCreatedFirst(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	require.NoError(t, err, "failed to create a temporary directory: %s", err)
	gun := "docker.com/notary"

	requestMade := false
	ts := httptest.NewServer(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) { requestMade = true }))
	defer ts.Close()

	repo, err := NewNotaryRepository(
		tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	require.NoError(t, err, "error creating repo: %s", err)

	cs := cryptoservice.NewCryptoService(gun,
		trustmanager.NewKeyMemoryStore(passphraseRetriever))

	rootPubKey, err := cs.Create(data.CanonicalRootRole, data.ECDSAKey)
	require.NoError(t, err, "error generating root key: %s", err)

	repo.CryptoService = cannotCreateKeys{CryptoService: cs}

	err = repo.Initialize(rootPubKey.ID(), data.CanonicalSnapshotRole)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Oh no I cannot create keys")
	require.False(t, requestMade)
}

// Publishing delegations works so long as the delegation parent exists by the
// time that delegation addition change is applied.  Most of the tests for
// applying delegation changes in in helpers_test.go (applyTargets tests), so
// this is just a sanity test to make sure Publish calls it correctly and
// no fallback happens.
func TestPublishDelegations(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo1, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo1.baseDir)

	delgKey, err := repo1.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// This should publish fine, even though targets/a/b is dependent upon
	// targets/a, because these should execute in order
	for _, delgName := range []string{"targets/a", "targets/a/b", "targets/c"} {
		require.NoError(t,
			repo1.AddDelegation(delgName, 1, []data.PublicKey{delgKey}, []string{""}),
			"error creating delegation")
	}
	require.Len(t, getChanges(t, repo1), 3, "wrong number of changelist files found")
	require.NoError(t, repo1.Publish())
	require.Len(t, getChanges(t, repo1), 0, "wrong number of changelist files found")

	// this should not publish, because targets/z doesn't exist
	require.NoError(t,
		repo1.AddDelegation("targets/z/y", 1, []data.PublicKey{delgKey}, []string{""}),
		"error creating delegation")
	require.Len(t, getChanges(t, repo1), 1, "wrong number of changelist files found")
	require.Error(t, repo1.Publish())
	require.Len(t, getChanges(t, repo1), 1, "wrong number of changelist files found")

	// use another repo to check metadata
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	// pull
	_, err = repo2.ListTargets()
	require.NoError(t, err, "unable to pull repo")

	for _, repo := range []*NotaryRepository{repo1, repo2} {
		// targets should have delegations targets/a and targets/c
		targets := repo.tufRepo.Targets[data.CanonicalTargetsRole]
		require.Len(t, targets.Signed.Delegations.Roles, 2)
		require.Len(t, targets.Signed.Delegations.Keys, 1)

		_, ok := targets.Signed.Delegations.Keys[delgKey.ID()]
		require.True(t, ok)

		foundRoleNames := make(map[string]bool)
		for _, r := range targets.Signed.Delegations.Roles {
			foundRoleNames[r.Name] = true
		}
		require.True(t, foundRoleNames["targets/a"])
		require.True(t, foundRoleNames["targets/c"])

		// targets/a should have delegation targets/a/b only
		a := repo.tufRepo.Targets["targets/a"]
		require.Len(t, a.Signed.Delegations.Roles, 1)
		require.Len(t, a.Signed.Delegations.Keys, 1)

		_, ok = a.Signed.Delegations.Keys[delgKey.ID()]
		require.True(t, ok)

		require.Equal(t, "targets/a/b", a.Signed.Delegations.Roles[0].Name)
	}
}

// Publishing delegations works so long as the delegation parent exists by the
// time that delegation addition change is applied.  Most of the tests for
// applying delegation changes in in helpers_test.go (applyTargets tests), so
// this is just a sanity test to make sure Publish calls it correctly and
// no fallback happens.
func TestPublishDelegationsX509(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo1, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo1.baseDir)

	delgKey, err := repo1.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	start := time.Now()
	privKey, _, err := repo1.CryptoService.GetPrivateKey(delgKey.ID())
	require.NoError(t, err)
	cert, err := cryptoservice.GenerateCertificate(
		privKey, "targets/a", start, start.AddDate(1, 0, 0),
	)
	require.NoError(t, err)
	delgCert := data.NewECDSAx509PublicKey(trustmanager.CertToPEM(cert))

	// This should publish fine, even though targets/a/b is dependent upon
	// targets/a, because these should execute in order
	for _, delgName := range []string{"targets/a", "targets/a/b", "targets/c"} {
		require.NoError(t,
			repo1.AddDelegation(delgName, 1, []data.PublicKey{delgCert}, []string{""}),
			"error creating delegation")
	}
	require.Len(t, getChanges(t, repo1), 3, "wrong number of changelist files found")
	require.NoError(t, repo1.Publish())
	require.Len(t, getChanges(t, repo1), 0, "wrong number of changelist files found")

	// this should not publish, because targets/z doesn't exist
	require.NoError(t,
		repo1.AddDelegation("targets/z/y", 1, []data.PublicKey{delgCert}, []string{""}),
		"error creating delegation")
	require.Len(t, getChanges(t, repo1), 1, "wrong number of changelist files found")
	require.Error(t, repo1.Publish())
	require.Len(t, getChanges(t, repo1), 1, "wrong number of changelist files found")

	// Create a new repo and pull from the server
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	// pull
	_, err = repo2.ListTargets()
	require.NoError(t, err, "unable to pull repo")

	for _, repo := range []*NotaryRepository{repo1, repo2} {
		// targets should have delegations targets/a and targets/c
		targets := repo.tufRepo.Targets[data.CanonicalTargetsRole]
		require.Len(t, targets.Signed.Delegations.Roles, 2)
		require.Len(t, targets.Signed.Delegations.Keys, 1)

		_, ok := targets.Signed.Delegations.Keys[delgCert.ID()]
		require.True(t, ok)

		foundRoleNames := make(map[string]bool)
		for _, r := range targets.Signed.Delegations.Roles {
			foundRoleNames[r.Name] = true
		}
		require.True(t, foundRoleNames["targets/a"])
		require.True(t, foundRoleNames["targets/c"])

		// targets/a should have delegation targets/a/b only
		a := repo.tufRepo.Targets["targets/a"]
		require.Len(t, a.Signed.Delegations.Roles, 1)
		require.Len(t, a.Signed.Delegations.Keys, 1)

		_, ok = a.Signed.Delegations.Keys[delgCert.ID()]
		require.True(t, ok)

		require.Equal(t, "targets/a/b", a.Signed.Delegations.Roles[0].Name)
	}
}

// If a changelist specifies a particular role to push targets to, and there
// is no such role, publish will try to publish to its parent.  If the parent
// doesn't work, it falls back on its parent, and so forth, and eventually
// falls back on publishing to "target".  This *only* falls back if the role
// doesn't exist, not if the user doesn't have a key.  (different test)
func TestPublishTargetsDelgationScopeFallback(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	requirePublishToRolesSucceeds(t, repo, []string{"targets/a/b", "targets/b/c"},
		[]string{data.CanonicalTargetsRole})
}

// If a changelist specifies a particular role to push targets to, and there
// is a role but no key, publish not fall back and just fail.
func TestPublishTargetsDelgationScopeNoFallbackIfNoKeys(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// generate a key that isn't in the cryptoservice, so we can't sign this
	// one
	aPrivKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err, "error generating key that is not in our cryptoservice")
	aPubKey := data.PublicKeyFromPrivate(aPrivKey)

	// ensure that the role exists
	require.NoError(t, repo.AddDelegation("targets/a", 1, []data.PublicKey{aPubKey}, []string{""}))
	require.NoError(t, repo.Publish())

	// add a target to targets/a/b - no role b, so it falls back on a, which
	// exists but there is no signing key for
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt", "targets/a/b")
	require.Len(t, getChanges(t, repo), 1, "wrong number of changelist files found")

	// Now Publish should fail
	require.Error(t, repo.Publish())
	require.Len(t, getChanges(t, repo), 1, "wrong number of changelist files found")

	targets, err := repo.ListTargets("targets", "targets/a", "targets/a/b")
	require.NoError(t, err)
	require.Empty(t, targets)
}

// If a changelist specifies a particular role to push targets to, and such
// a role and the keys are present, publish will write to that role only, and
// not its parents.  This tests the case where the local machine knows about
// all the roles (in fact, the role creations will be applied before the
// targets)
func TestPublishTargetsDelgationSuccessLocallyHasRoles(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	delgKey, err := repo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	for _, delgName := range []string{"targets/a", "targets/a/b"} {
		require.NoError(t,
			repo.AddDelegation(delgName, 1, []data.PublicKey{delgKey}, []string{""}),
			"error creating delegation")
	}

	requirePublishToRolesSucceeds(t, repo, []string{"targets/a/b"},
		[]string{"targets/a/b"})
}

// If a changelist specifies a particular role to push targets to, and the role
// is present, publish will write to that role only.  The targets keys are not
// needed.
func TestPublishTargetsDelgationNoTargetsKeyNeeded(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	delgKey, err := repo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	for _, delgName := range []string{"targets/a", "targets/a/b"} {
		require.NoError(t,
			repo.AddDelegation(delgName, 1, []data.PublicKey{delgKey}, []string{""}),
			"error creating delegation")
	}
	require.NoError(t, repo.Publish())

	// remove targets key - it is not even needed
	targetsKeys := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
	require.Len(t, targetsKeys, 1)
	require.NoError(t, repo.CryptoService.RemoveKey(targetsKeys[0]))

	requirePublishToRolesSucceeds(t, repo, []string{"targets/a/b"},
		[]string{"targets/a/b"})
}

// If a changelist specifies a particular role to push targets to, and is such
// a role and the keys are present, publish will write to that role only, and
// not its parents.  Tests:
// - case where the local doesn't know about all the roles, and has to download
//   them before publish.
// - owner of a repo may not have the delegated keys, so can't sign a delegated
//   role
func TestPublishTargetsDelgationSuccessNeedsToDownloadRoles(t *testing.T) {
	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	// this is the original repo - it owns the root/targets keys and creates
	// the delegation to which it doesn't have the key (so server snapshot
	// signing would be required)
	ownerRepo, _ := initializeRepo(t, data.ECDSAKey, gun, ts.URL, true)
	defer os.RemoveAll(ownerRepo.baseDir)

	// this is a user, or otherwise a repo that only has access to the delegation
	// key so it can publish targets to the delegated role
	delgRepo := newRepoToTestRepo(t, ownerRepo)
	defer os.RemoveAll(delgRepo.baseDir)

	// create a key on the owner repo
	aKey, err := ownerRepo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// create a key on the delegated repo
	bKey, err := delgRepo.CryptoService.Create("targets/a/b", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// owner creates delegations, adds the delegated key to them, and publishes them
	require.NoError(t,
		ownerRepo.AddDelegation("targets/a", 1, []data.PublicKey{aKey}, []string{""}),
		"error creating delegation")
	require.NoError(t,
		ownerRepo.AddDelegation("targets/a/b", 1, []data.PublicKey{bKey}, []string{""}),
		"error creating delegation")

	require.NoError(t, ownerRepo.Publish())

	// delegated repo now publishes to delegated roles, but it will need
	// to download those roles first, since it doesn't know about them
	requirePublishToRolesSucceeds(t, delgRepo, []string{"targets/a/b"},
		[]string{"targets/a/b"})
}

// Ensure that two clients can publish delegations with two different keys and
// the changes will not clobber each other.
func TestPublishTargetsDelgationFromTwoRepos(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	// this happens to be the client that creates the repo, but can also
	// write a delegation
	repo1, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, true)
	defer os.RemoveAll(repo1.baseDir)

	// this is the second writable repo
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	// create keys for each repo
	key1, err := repo1.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// create a key on the delegated repo
	key2, err := repo2.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// delegation includes both keys
	require.NoError(t,
		repo1.AddDelegation("targets/a", 1, []data.PublicKey{key1, key2}, []string{""}),
		"error creating delegation")

	require.NoError(t, repo1.Publish())

	// both repos add targets and publish
	addTarget(t, repo1, "first", "../fixtures/root-ca.crt", "targets/a")
	require.NoError(t, repo1.Publish())
	addTarget(t, repo2, "second", "../fixtures/root-ca.crt", "targets/a")
	require.NoError(t, repo2.Publish())

	// first repo can publish again
	addTarget(t, repo1, "third", "../fixtures/root-ca.crt", "targets/a")
	require.NoError(t, repo1.Publish())

	// both repos should be able to see all targets
	for _, repo := range []*NotaryRepository{repo1, repo2} {
		targets, err := repo.ListTargets()
		require.NoError(t, err)
		require.Len(t, targets, 3)

		found := make(map[string]bool)
		for _, t := range targets {
			found[t.Name] = true
		}

		for _, targetName := range []string{"first", "second", "third"} {
			_, ok := found[targetName]
			require.True(t, ok)
		}
	}
}

// A client who could publish before can no longer publish once the owner
// removes their delegation key from the delegation role.
func TestPublishRemoveDelgationKeyFromDelegationRole(t *testing.T) {
	gun := "docker.com/notary"
	ts := fullTestServer(t)
	defer ts.Close()

	// this is the original repo - it owns the root/targets keys and creates
	// the delegation to which it doesn't have the key (so server snapshot
	// signing would be required)
	ownerRepo, _ := initializeRepo(t, data.ECDSAKey, gun, ts.URL, true)
	defer os.RemoveAll(ownerRepo.baseDir)

	// this is a user, or otherwise a repo that only has access to the delegation
	// key so it can publish targets to the delegated role
	delgRepo := newRepoToTestRepo(t, ownerRepo)
	defer os.RemoveAll(delgRepo.baseDir)

	// create a key on the delegated repo
	aKey, err := delgRepo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// owner creates delegation, adds the delegated key to it, and publishes it
	require.NoError(t,
		ownerRepo.AddDelegation("targets/a", 1, []data.PublicKey{aKey}, []string{""}),
		"error creating delegation")
	require.NoError(t, ownerRepo.Publish())

	// delegated repo can now publish to delegated role
	addTarget(t, delgRepo, "v1", "../fixtures/root-ca.crt", "targets/a")
	require.NoError(t, delgRepo.Publish())

	// owner revokes delegation
	// note there is no removekeyfromdelegation yet, so here's a hack to do so
	newKey, err := ownerRepo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err)
	tdJSON, err := json.Marshal(&changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      data.KeyList([]data.PublicKey{newKey}),
		RemoveKeys:   []string{aKey.ID()},
	})
	require.NoError(t, err)

	cl, err := changelist.NewFileChangelist(filepath.Join(ownerRepo.tufRepoPath, "changelist"))
	require.NoError(t, cl.Add(changelist.NewTufChange(
		changelist.ActionUpdate,
		"targets/a",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)))
	cl.Close()
	require.NoError(t, ownerRepo.Publish())

	// delegated repo can now no longer publish to delegated role
	addTarget(t, delgRepo, "v2", "../fixtures/root-ca.crt", "targets/a")
	require.Error(t, delgRepo.Publish())
}

// A client who could publish before can no longer publish once the owner
// deletes the delegation
func TestPublishRemoveDelgation(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	// this is the original repo - it owns the root/targets keys and creates
	// the delegation to which it doesn't have the key (so server snapshot
	// signing would be required)
	ownerRepo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, true)
	defer os.RemoveAll(ownerRepo.baseDir)

	// this is a user, or otherwise a repo that only has access to the delegation
	// key so it can publish targets to the delegated role
	delgRepo := newRepoToTestRepo(t, ownerRepo)
	defer os.RemoveAll(delgRepo.baseDir)

	// create a key on the delegated repo
	aKey, err := delgRepo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	// owner creates delegation, adds the delegated key to it, and publishes it
	require.NoError(t,
		ownerRepo.AddDelegation("targets/a", 1, []data.PublicKey{aKey}, []string{""}),
		"error creating delegation")
	require.NoError(t, ownerRepo.Publish())

	// delegated repo can now publish to delegated role
	addTarget(t, delgRepo, "v1", "../fixtures/root-ca.crt", "targets/a")
	require.NoError(t, delgRepo.Publish())

	// owner removes delegation
	require.NoError(t, ownerRepo.RemoveDelegation("targets/a"))
	require.NoError(t, ownerRepo.Publish())

	// delegated repo can now no longer publish to delegated role
	addTarget(t, delgRepo, "v2", "../fixtures/root-ca.crt", "targets/a")
	require.Error(t, delgRepo.Publish())
}

// If the delegation data is corrupt or unreadable, it doesn't matter because
// all the delegation information is just re-downloaded.  When bootstrapping
// the repository from disk, we just don't load the data from disk because
// there should not be anything there yet.
func TestPublishSucceedsDespiteDelegationCorrupt(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	delgKey, err := repo.CryptoService.Create("targets/a", data.ECDSAKey)
	require.NoError(t, err, "error creating delegation key")

	require.NoError(t,
		repo.AddDelegation("targets/a", 1, []data.PublicKey{delgKey}, []string{""}),
		"error creating delegation")

	testPublishBadMetadata(t, "targets/a", repo, false, true)

	// publish again, now that it has already been published, and again there
	// is no error.
	testPublishBadMetadata(t, "targets/a", repo, true, true)
}

// Rotate invalid roles, or attempt to delegate target signing to the server
func TestRotateKeyInvalidRole(t *testing.T) {
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// the equivalent of: (root, true), (root, false), (timestamp, true),
	// (timestamp, false), (targets, true)
	for role := range data.ValidRoles {
		if role == data.CanonicalSnapshotRole {
			continue
		}
		for _, serverManagesKey := range []bool{true, false} {
			if role == data.CanonicalTargetsRole && !serverManagesKey {
				continue
			}
			err := repo.RotateKey(role, serverManagesKey)
			require.Error(t, err,
				"Rotating a %s key with server-managing the key as %v should fail",
				role, serverManagesKey)
		}
	}
}

// Rotates the keys.  After the rotation, downloading the latest metadata
// and require that the keys have changed
func requireRotationSuccessful(t *testing.T, repo1 *NotaryRepository,
	keysToRotate map[string]bool, alreadyPublished bool) {
	// Create 2 new repos:  1 will download repo data before the publish,
	// and one only downloads after the publish. This reflects a client
	// that has some previous trust data (but is not the publisher), and a
	// completely new client being able to read the rotated trust data.
	repo2 := newRepoToTestRepo(t, repo1)
	defer os.RemoveAll(repo2.baseDir)

	repos := []*NotaryRepository{repo1, repo2}

	if alreadyPublished {
		repo3 := newRepoToTestRepo(t, repo1)
		defer os.RemoveAll(repo2.baseDir)

		// force a pull on repo3
		_, err := repo3.GetTargetByName("latest")
		require.NoError(t, err)

		repos = append(repos, repo3)
	}

	oldKeyIDs := make(map[string][]string)
	for role := range keysToRotate {
		keyIDs := repo1.tufRepo.Root.Signed.Roles[role].KeyIDs
		oldKeyIDs[role] = keyIDs
	}

	// Do rotation
	for role, serverManaged := range keysToRotate {
		require.NoError(t, repo1.RotateKey(role, serverManaged))
	}

	// Publish
	err := repo1.Publish()
	require.NoError(t, err)

	// Download data from remote and check that keys have changed
	for _, repo := range repos {
		_, err := repo.GetTargetByName("latest") // force a pull
		require.NoError(t, err)

		for role, isRemoteKey := range keysToRotate {
			keyIDs := repo.tufRepo.Root.Signed.Roles[role].KeyIDs
			require.Len(t, keyIDs, 1)

			// the new key is not the same as any of the old keys, and the
			// old keys have been removed not just from the TUF file, but
			// from the cryptoservice
			for _, oldKeyID := range oldKeyIDs[role] {
				require.NotEqual(t, oldKeyID, keyIDs[0])
				_, _, err := repo.CryptoService.GetPrivateKey(oldKeyID)
				require.Error(t, err)
			}

			// On the old repo, the new key is present in the cryptoservice, or
			// not present if remote.  On the new repo, no keys are ever in the
			// cryptoservice
			key, _, err := repo.CryptoService.GetPrivateKey(keyIDs[0])
			if repo != repo1 || isRemoteKey {
				require.Error(t, err)
				require.Nil(t, key)
			} else {
				require.NoError(t, err)
				require.NotNil(t, key)
			}
		}

		// Confirm changelist dir empty (on repo1, it should be empty after
		// after publishing changes, on repo2, there should never have been
		// any changelists)
		changes := getChanges(t, repo)
		require.Len(t, changes, 0, "wrong number of changelist files found")
	}
}

// Initialize repo to have the server sign snapshots (remote snapshot key)
// Without downloading a server-signed snapshot file, rotate keys so that
//    snapshots are locally signed (local snapshot key)
// require that we can publish.
func TestRotateBeforePublishFromRemoteKeyToLocalKey(t *testing.T) {
	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL, true)
	defer os.RemoveAll(repo.baseDir)

	// Adding a target will allow us to confirm the repository is still valid
	// after rotating the keys.
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")
	requireRotationSuccessful(t, repo, map[string]bool{
		data.CanonicalTargetsRole:  false,
		data.CanonicalSnapshotRole: false}, false)
}

// Initialize a repo, locally signed snapshots
// Publish some content (so that the server has a root.json), and download root.json
// Rotate keys
// Download the latest metadata and require that the keys have changed.
func TestRotateKeyAfterPublishNoServerManagementChange(t *testing.T) {
	// rotate a single target key
	testRotateKeySuccess(t, false, map[string]bool{data.CanonicalTargetsRole: false})
	testRotateKeySuccess(t, false, map[string]bool{data.CanonicalSnapshotRole: false})
	// rotate two at once before publishing
	testRotateKeySuccess(t, false, map[string]bool{
		data.CanonicalSnapshotRole: false,
		data.CanonicalTargetsRole:  false})
}

// Tests rotating keys when there's a change from locally managed keys to
// remotely managed keys and vice versa
// Before rotating, publish some content (so that the server has a root.json),
// and download root.json
func TestRotateKeyAfterPublishServerManagementChange(t *testing.T) {
	// delegate snapshot key management to the server
	testRotateKeySuccess(t, false, map[string]bool{
		data.CanonicalSnapshotRole: true,
		data.CanonicalTargetsRole:  false,
	})
	// reclaim snapshot key management from the server
	testRotateKeySuccess(t, true, map[string]bool{
		data.CanonicalSnapshotRole: false,
		data.CanonicalTargetsRole:  false,
	})
}

func testRotateKeySuccess(t *testing.T, serverManagesSnapshotInit bool,
	keysToRotate map[string]bool) {

	ts := fullTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, "docker.com/notary", ts.URL,
		serverManagesSnapshotInit)
	defer os.RemoveAll(repo.baseDir)

	// Adding a target will allow us to confirm the repository is still valid after
	// rotating the keys.
	addTarget(t, repo, "latest", "../fixtures/intermediate-ca.crt")

	// Publish
	require.NoError(t, repo.Publish())

	// Get root.json and capture targets + snapshot key IDs
	repo.GetTargetByName("latest") // force a pull
	requireRotationSuccessful(t, repo, keysToRotate, true)
}

// If there is no local cache, notary operations return the remote error code
func TestRemoteServerUnavailableNoLocalCache(t *testing.T) {
	tempBaseDir, err := ioutil.TempDir("/tmp", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory: %s", err)
	defer os.RemoveAll(tempBaseDir)

	ts := errorTestServer(t, 500)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, "docker.com/notary",
		ts.URL, http.DefaultTransport, passphraseRetriever)
	require.NoError(t, err, "error creating repo: %s", err)

	_, err = repo.ListTargets(data.CanonicalTargetsRole)
	require.Error(t, err)
	require.IsType(t, store.ErrServerUnavailable{}, err)

	_, err = repo.GetTargetByName("targetName")
	require.Error(t, err)
	require.IsType(t, store.ErrServerUnavailable{}, err)

	err = repo.Publish()
	require.Error(t, err)
	require.IsType(t, store.ErrServerUnavailable{}, err)
}

// AddDelegation creates a valid changefile (rejects invalid delegation names,
// but does not check the delegation hierarchy).  When applied, the change adds
// a new delegation role with the correct keys.
func TestAddDelegationChangefileValid(t *testing.T) {
	gun := "docker.com/notary"
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, gun, ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	targetKeyIds := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
	require.NotEmpty(t, targetKeyIds)
	targetPubKey := repo.CryptoService.GetKey(targetKeyIds[0])
	require.NotNil(t, targetPubKey)

	err := repo.AddDelegation("root", 1, []data.PublicKey{targetPubKey}, []string{""})
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
	require.Empty(t, getChanges(t, repo))

	// to show that adding does not care about the hierarchy
	err = repo.AddDelegation("targets/a/b/c", 1, []data.PublicKey{targetPubKey}, []string{""})
	require.NoError(t, err)

	// ensure that the changefiles is correct
	changes := getChanges(t, repo)
	require.Len(t, changes, 1)
	require.Equal(t, changelist.ActionCreate, changes[0].Action())
	require.Equal(t, "targets/a/b/c", changes[0].Scope())
	require.Equal(t, changelist.TypeTargetsDelegation, changes[0].Type())
	require.Equal(t, "", changes[0].Path())
	require.NotEmpty(t, changes[0].Content())
}

// The changefile produced by AddDelegation, when applied, actually adds
// the delegation to the repo (assuming the delegation hierarchy is correct -
// tests for change application validation are in helpers_test.go)
func TestAddDelegationChangefileApplicable(t *testing.T) {
	gun := "docker.com/notary"
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, data.ECDSAKey, gun, ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	targetKeyIds := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
	require.NotEmpty(t, targetKeyIds)
	targetPubKey := repo.CryptoService.GetKey(targetKeyIds[0])
	require.NotNil(t, targetPubKey)

	// this hierarchy has to be right to be applied
	err := repo.AddDelegation("targets/a", 1, []data.PublicKey{targetPubKey}, []string{""})
	require.NoError(t, err)
	changes := getChanges(t, repo)
	require.Len(t, changes, 1)

	// ensure that it can be applied correctly
	err = applyTargetsChange(repo.tufRepo, changes[0])
	require.NoError(t, err)

	targetRole := repo.tufRepo.Targets[data.CanonicalTargetsRole]
	require.Len(t, targetRole.Signed.Delegations.Roles, 1)
	require.Len(t, targetRole.Signed.Delegations.Keys, 1)

	_, ok := targetRole.Signed.Delegations.Keys[targetPubKey.ID()]
	require.True(t, ok)

	newDelegationRole := targetRole.Signed.Delegations.Roles[0]
	require.Len(t, newDelegationRole.KeyIDs, 1)
	require.Equal(t, targetPubKey.ID(), newDelegationRole.KeyIDs[0])
	require.Equal(t, "targets/a", newDelegationRole.Name)
}

// TestAddDelegationErrorWritingChanges expects errors writing a change to file
// to be propagated.
func TestAddDelegationErrorWritingChanges(t *testing.T) {
	testErrorWritingChangefiles(t, func(repo *NotaryRepository) error {
		targetKeyIds := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
		require.NotEmpty(t, targetKeyIds)
		targetPubKey := repo.CryptoService.GetKey(targetKeyIds[0])
		require.NotNil(t, targetPubKey)

		return repo.AddDelegation("targets/a", 1, []data.PublicKey{targetPubKey}, []string{""})
	})
}

// RemoveDelegation rejects attempts to remove invalidly-named delegations,
// but otherwise does not validate the name of the delegation to remove.  This
// test ensures that the changefile generated by RemoveDelegation is correct.
func TestRemoveDelegationChangefileValid(t *testing.T) {
	gun := "docker.com/notary"
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootKeyID := initializeRepo(t, data.ECDSAKey, gun, ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	rootPubKey := repo.CryptoService.GetKey(rootKeyID)
	require.NotNil(t, rootPubKey)

	err := repo.RemoveDelegation("root")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
	require.Empty(t, getChanges(t, repo))

	// to demonstrate that so long as the delegation name is valid, the
	// existence of the delegation doesn't matter
	require.NoError(t, repo.RemoveDelegation("targets/a/b/c"))

	// ensure that the changefile is correct
	changes := getChanges(t, repo)
	require.Len(t, changes, 1)
	require.Equal(t, changelist.ActionDelete, changes[0].Action())
	require.Equal(t, "targets/a/b/c", changes[0].Scope())
	require.Equal(t, changelist.TypeTargetsDelegation, changes[0].Type())
	require.Equal(t, "", changes[0].Path())
	require.Empty(t, changes[0].Content())
}

// The changefile produced by RemoveDelegation, when applied, actually removes
// the delegation from the repo (assuming the repo exists - tests for
// change application validation are in helpers_test.go)
func TestRemoveDelegationChangefileApplicable(t *testing.T) {
	gun := "docker.com/notary"
	ts, _, _ := simpleTestServer(t)
	defer ts.Close()

	repo, rootKeyID := initializeRepo(t, data.ECDSAKey, gun, ts.URL, false)
	defer os.RemoveAll(repo.baseDir)
	rootPubKey := repo.CryptoService.GetKey(rootKeyID)
	require.NotNil(t, rootPubKey)

	// add a delegation first so it can be removed
	require.NoError(t, repo.AddDelegation("targets/a", 1, []data.PublicKey{rootPubKey}, []string{""}))
	changes := getChanges(t, repo)
	require.Len(t, changes, 1)
	require.NoError(t, applyTargetsChange(repo.tufRepo, changes[0]))

	targetRole := repo.tufRepo.Targets[data.CanonicalTargetsRole]
	require.Len(t, targetRole.Signed.Delegations.Roles, 1)
	require.Len(t, targetRole.Signed.Delegations.Keys, 1)

	// now remove it
	require.NoError(t, repo.RemoveDelegation("targets/a"))
	changes = getChanges(t, repo)
	require.Len(t, changes, 2)
	require.NoError(t, applyTargetsChange(repo.tufRepo, changes[1]))

	targetRole = repo.tufRepo.Targets[data.CanonicalTargetsRole]
	require.Empty(t, targetRole.Signed.Delegations.Roles)
	require.Empty(t, targetRole.Signed.Delegations.Keys)
}

// TestRemoveDelegationErrorWritingChanges expects errors writing a change to
// file to be propagated.
func TestRemoveDelegationErrorWritingChanges(t *testing.T) {
	testErrorWritingChangefiles(t, func(repo *NotaryRepository) error {
		return repo.RemoveDelegation("targets/a")
	})
}
