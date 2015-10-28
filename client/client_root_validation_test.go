package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

var passphraseRetriever = func(string, string, bool, int) (string, bool, error) { return "passphrase", false, nil }

// TestValidateRoot through the process of initializing a repository and makes
// sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key
func TestValidateRoot(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	validateRootSuccessfully(t, data.ECDSAKey)
	if !testing.Short() {
		validateRootSuccessfully(t, data.RSAKey)
	}
}

func validateRootSuccessfully(t *testing.T, rootType data.KeyAlgorithm) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, mux := createTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err = repo.tufRepo.InitTimestamp()
	assert.NoError(t, err, "error creating repository: %s", err)

	// Initialize is supposed to have created new certificate for this repository
	// Lets check for it and store it for later use
	allCerts := repo.KeyStoreManager.TrustedCertificateStore().GetCertificates()
	assert.Len(t, allCerts, 1)

	// Now test ListTargets. In preparation, we need to expose some signed
	// metadata files on the internal HTTP server.
	var tempKey data.TUFKey
	json.Unmarshal([]byte(timestampECDSAKeyJSON), &tempKey)

	repo.KeyStoreManager.NonRootKeyStore().AddKey(filepath.Join(filepath.FromSlash(gun), tempKey.ID()), "root", &tempKey)

	// Because ListTargets will clear this
	savedTUFRepo := repo.tufRepo

	rootJSONFile := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "metadata", "root.json")
	rootFileBytes, err := ioutil.ReadFile(rootJSONFile)

	signedTargets, err := savedTUFRepo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	assert.NoError(t, err)

	signedSnapshot, err := savedTUFRepo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)

	signedTimestamp, err := savedTUFRepo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/root.json", func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, err)
		fmt.Fprint(w, string(rootFileBytes))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.json", func(w http.ResponseWriter, r *http.Request) {
		timestampJSON, _ := json.Marshal(signedTimestamp)
		fmt.Fprint(w, string(timestampJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/snapshot.json", func(w http.ResponseWriter, r *http.Request) {
		snapshotJSON, _ := json.Marshal(signedSnapshot)
		fmt.Fprint(w, string(snapshotJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets.json", func(w http.ResponseWriter, r *http.Request) {
		targetsJSON, _ := json.Marshal(signedTargets)
		fmt.Fprint(w, string(targetsJSON))
	})

	_, err = repo.ListTargets()
	assert.NoError(t, err)

	//
	// Test TOFUS logic. We remove all certs and expect a new one to be added after ListTargets
	//
	err = repo.KeyStoreManager.TrustedCertificateStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 0)

	// This list targets is expected to succeed and the certificate store to have the new certificate
	_, err = repo.ListTargets()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 1)

	//
	// Test certificate mismatch logic. We remove all certs, add a different cert to the
	// same CN, and expect ValidateRoot to fail
	//

	// First, remove all certs
	err = repo.KeyStoreManager.TrustedCertificateStore().RemoveAll()
	assert.NoError(t, err)
	assert.Len(t, repo.KeyStoreManager.TrustedCertificateStore().GetCertificates(), 0)

	// Add a previously generated certificate with CN=docker.com/notary
	err = repo.KeyStoreManager.TrustedCertificateStore().AddCertFromFile("../fixtures/self-signed_docker.com-notary.crt")
	assert.NoError(t, err)

	// This list targets is expected to fail, since there already exists a certificate
	// in the store for the dnsName docker.com/notary, so TOFUS doesn't apply
	_, err = repo.ListTargets()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, &keystoremanager.ErrValidationFail{Reason: "failed to validate data with current trusted certificates"})
	}
}
