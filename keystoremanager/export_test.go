package keystoremanager_test

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/notary/client"
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

func TestExportKeys(t *testing.T) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _ := createTestServer(t)
	defer ts.Close()

	repo, err := client.NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(data.ECDSAKey.String(), "oldPassphrase")
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID, "oldPassphrase")
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	tempZipFile, err := ioutil.TempFile("", "notary-test-export-")
	tempZipFilePath := tempZipFile.Name()
	defer os.Remove(tempZipFilePath)

	err = repo.KeyStoreManager.ExportAllKeys(tempZipFile, "exportPassphrase")
	tempZipFile.Close()
	assert.NoError(t, err)

	// Check the contents of the zip file
	zipReader, err := zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")
	defer zipReader.Close()

	// Map of files to expect in the zip file, with the passphrases
	passphraseByFile := make(map[string]string)

	// Add keys in private to the map. These should use the new passphrase
	// because they were formerly unencrypted.
	privKeyList := repo.KeyStoreManager.NonRootKeyStore().ListFiles(false)
	for _, privKeyName := range privKeyList {
		relName := strings.TrimPrefix(privKeyName, tempBaseDir+string(filepath.Separator))
		passphraseByFile[relName] = "exportPassphrase"
	}

	// Add root key to the map. This will use the old passphrase because it
	// won't be reencrypted.
	relRootKey := filepath.Join("private", "root_keys", rootCryptoService.ID()+".key")
	passphraseByFile[relRootKey] = "oldPassphrase"

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

	// Are there any keys that didn't make it to the zip?
	for fileNotFound, _ := range passphraseByFile {
		t.Fatalf("%s not found in zip", fileNotFound)
	}
}
