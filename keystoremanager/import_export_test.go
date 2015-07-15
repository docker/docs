package keystoremanager_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/notary/client"
	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
)

const timestampECDSAKeyJSON = `
{"keytype":"ecdsa","keyval":{"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw==","private":"MHcCAQEEIDqtcdzU7H3AbIPSQaxHl9+xYECt7NpK7B1+6ep5cv9CoAoGCCqGSM49AwEHoUQDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw=="}}`

func createTestServer(t *testing.T) (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	// TUF will request /v2/docker.com/notary/_trust/tuf/timestamp.key
	// Return a canned timestamp.key
	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.key", func(w http.ResponseWriter, r *http.Request) {
		// Also contains the private key, but for the purpose of this
		// test, we don't care
		fmt.Fprint(w, timestampECDSAKeyJSON)
	})

	ts := httptest.NewServer(mux)

	return ts, mux
}

func TestImportExportZip(t *testing.T) {
	gun := "docker.com/notary"
	oldPassphrase := "oldPassphrase"
	exportPassphrase := "exportPassphrase"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _ := createTestServer(t)
	defer ts.Close()

	repo, err := client.NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(data.ECDSAKey.String(), oldPassphrase)
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID, oldPassphrase)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	tempZipFile, err := ioutil.TempFile("", "notary-test-export-")
	tempZipFilePath := tempZipFile.Name()
	defer os.Remove(tempZipFilePath)

	err = repo.KeyStoreManager.ExportAllKeys(tempZipFile, exportPassphrase)
	tempZipFile.Close()
	assert.NoError(t, err)

	// Reopen the zip file for importing
	zipReader, err := zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Map of files to expect in the zip file, with the passphrases
	passphraseByFile := make(map[string]string)

	// Add keys in private to the map. These should use the new passphrase
	// because they were formerly unencrypted.
	privKeyList := repo.KeyStoreManager.NonRootKeyStore().ListFiles(false)
	for _, privKeyName := range privKeyList {
		relName := strings.TrimPrefix(privKeyName, tempBaseDir+string(filepath.Separator))
		passphraseByFile[relName] = exportPassphrase
	}

	// Add root key to the map. This will use the old passphrase because it
	// won't be reencrypted.
	relRootKey := filepath.Join("private", "root_keys", rootCryptoService.ID()+".key")
	passphraseByFile[relRootKey] = oldPassphrase

	// Iterate through the files in the archive, checking that the files
	// exist and are encrypted with the expected passphrase.
	for _, f := range zipReader.File {
		expectedPassphrase, present := passphraseByFile[f.Name]
		if !present {
			t.Fatalf("unexpected file %s in zip file", f.Name)
		}

		delete(passphraseByFile, f.Name)

		rc, err := f.Open()
		assert.NoError(t, err, "could not open file inside zip archive")

		pemBytes, err := ioutil.ReadAll(rc)
		assert.NoError(t, err, "could not read file from zip")

		_, err = trustmanager.ParsePEMPrivateKey(pemBytes, expectedPassphrase)
		assert.NoError(t, err, "PEM not encrypted with the expected passphrase")

		rc.Close()
	}

	zipReader.Close()

	// Are there any keys that didn't make it to the zip?
	for fileNotFound := range passphraseByFile {
		t.Fatalf("%s not found in zip", fileNotFound)
	}

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo2, err := client.NewNotaryRepository(tempBaseDir2, gun, ts.URL, http.DefaultTransport)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID2, err := repo2.KeyStoreManager.GenRootKey(data.ECDSAKey.String(), "oldPassphrase")
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService2, err := repo2.KeyStoreManager.GetRootCryptoService(rootKeyID2, "oldPassphrase")
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo2.Initialize(rootCryptoService2)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Reopen the zip file for importing
	zipReader, err = zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// First try with an incorrect passphrase
	err = repo2.KeyStoreManager.ImportKeysZip(zipReader.Reader, "wrongpassphrase")
	// Don't use EqualError here because occasionally decrypting with the
	// wrong passphrase returns a parse error
	assert.Error(t, err)
	zipReader.Close()

	// Reopen the zip file for importing
	zipReader, err = zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Now try with a valid passphrase. This time it should succeed.
	err = repo2.KeyStoreManager.ImportKeysZip(zipReader.Reader, exportPassphrase)
	assert.NoError(t, err)
	zipReader.Close()

	// Look for repo's keys in repo2

	// Look for keys in private. The filenames should match the key IDs
	// in the repo's private key store.
	for _, privKeyName := range privKeyList {
		privKeyRel := strings.TrimPrefix(privKeyName, tempBaseDir)
		_, err := os.Stat(filepath.Join(tempBaseDir2, privKeyRel))
		assert.NoError(t, err, "missing private key: %s", privKeyName)
	}

	// Look for keys in root_keys
	// There should be a file named after the key ID of the root key we
	// passed in.
	rootKeyFilename := rootCryptoService.ID() + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")
}

func TestImportExportRootKey(t *testing.T) {
	gun := "docker.com/notary"
	oldPassphrase := "oldPassphrase"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _ := createTestServer(t)
	defer ts.Close()

	repo, err := client.NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(data.ECDSAKey.String(), oldPassphrase)
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID, oldPassphrase)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	tempKeyFile, err := ioutil.TempFile("", "notary-test-export-")
	tempKeyFilePath := tempKeyFile.Name()
	defer os.Remove(tempKeyFilePath)

	err = repo.KeyStoreManager.ExportRootKey(tempKeyFile, rootKeyID)
	assert.NoError(t, err)
	tempKeyFile.Close()

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo2, err := client.NewNotaryRepository(tempBaseDir2, gun, ts.URL, http.DefaultTransport)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID2, err := repo2.KeyStoreManager.GenRootKey(data.ECDSAKey.String(), oldPassphrase)
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService2, err := repo2.KeyStoreManager.GetRootCryptoService(rootKeyID2, oldPassphrase)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo2.Initialize(rootCryptoService2)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Check the contents of the zip file
	keyReader, err := os.Open(tempKeyFilePath)
	assert.NoError(t, err, "could not open key file")

	err = repo2.KeyStoreManager.ImportRootKey(keyReader, rootKeyID)
	assert.NoError(t, err)
	keyReader.Close()

	// Look for repo's root key in repo2
	// There should be a file named after the key ID of the root key we
	// imported.
	rootKeyFilename := rootKeyID + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")

	// Try to import a decrypted version of the root key and make sure it
	// doesn't succeed
	pemBytes, err := ioutil.ReadFile(tempKeyFilePath)
	assert.NoError(t, err, "could not read key file")
	privKey, err := trustmanager.ParsePEMPrivateKey(pemBytes, oldPassphrase)
	assert.NoError(t, err, "could not decrypt key file")
	decryptedPEMBytes, err := trustmanager.KeyToPEM(privKey)
	assert.NoError(t, err, "could not convert key to PEM")

	err = repo2.KeyStoreManager.ImportRootKey(bytes.NewReader(decryptedPEMBytes), rootKeyID)
	assert.EqualError(t, err, keystoremanager.ErrRootKeyNotEncrypted.Error())

	// Try to import garbage and make sure it doesn't succeed
	err = repo2.KeyStoreManager.ImportRootKey(strings.NewReader("this is not PEM"), rootKeyID)
	assert.EqualError(t, err, keystoremanager.ErrNoValidPrivateKey.Error())
}
