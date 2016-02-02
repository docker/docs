package cryptoservice

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

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
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

var oldPassphrase = "oldPassphrase"
var exportPassphrase = "exportPassphrase"
var oldPassphraseRetriever = func(string, string, bool, int) (string, bool, error) { return oldPassphrase, false, nil }
var newPassphraseRetriever = func(string, string, bool, int) (string, bool, error) { return exportPassphrase, false, nil }

func TestImportExportZip(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, newPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	pubKey, err := cs.Create(data.CanonicalRootRole, data.ECDSAKey)
	assert.NoError(t, err)

	rootKeyID := pubKey.ID()

	tempZipFile, err := ioutil.TempFile("", "notary-test-export-")
	tempZipFilePath := tempZipFile.Name()
	defer os.Remove(tempZipFilePath)

	err = cs.ExportAllKeys(tempZipFile, newPassphraseRetriever)
	tempZipFile.Close()
	assert.NoError(t, err)

	// Reopen the zip file for importing
	zipReader, err := zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Map of files to expect in the zip file, with the passphrases
	passphraseByFile := make(map[string]string)

	// Add non-root keys to the map. These should use the new passphrase
	// because the passwords were chosen by the newPassphraseRetriever.
	privKeyMap := cs.ListAllKeys()
	for privKeyName := range privKeyMap {
		_, alias, err := cs.GetPrivateKey(privKeyName)
		assert.NoError(t, err, "privKey %s has no alias", privKeyName)

		if alias == "root" {
			continue
		}
		relKeyPath := filepath.Join("tuf_keys", privKeyName+".key")
		passphraseByFile[relKeyPath] = exportPassphrase
	}

	// Add root key to the map. This will use the export passphrase because it
	// will be reencrypted.
	relRootKey := filepath.Join("root_keys", rootKeyID+".key")
	passphraseByFile[relRootKey] = exportPassphrase

	// Iterate through the files in the archive, checking that the files
	// exist and are encrypted with the expected passphrase.
	for _, f := range zipReader.File {
		expectedPassphrase, present := passphraseByFile[f.Name]
		assert.True(t, present, "unexpected file %s in zip file", f.Name)

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
	assert.Len(t, passphraseByFile, 0)

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, newPassphraseRetriever)
	assert.NoError(t, err)
	cs2 := NewCryptoService(gun, fileStore2)

	// Reopen the zip file for importing
	zipReader, err = zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Now try with a valid passphrase. This time it should succeed.
	err = cs2.ImportKeysZip(zipReader.Reader)
	assert.NoError(t, err)
	zipReader.Close()

	// Look for keys in private. The filenames should match the key IDs
	// in the repo's private key store.
	for privKeyName := range privKeyMap {
		_, alias, err := cs2.GetPrivateKey(privKeyName)
		assert.NoError(t, err, "privKey %s has no alias", privKeyName)

		if alias == "root" {
			continue
		}
		relKeyPath := filepath.Join("tuf_keys", privKeyName+".key")
		privKeyFileName := filepath.Join(tempBaseDir2, "private", relKeyPath)
		_, err = os.Stat(privKeyFileName)
		assert.NoError(t, err, "missing private key for role %s: %s", alias, privKeyName)
	}

	// Look for keys in root_keys
	// There should be a file named after the key ID of the root key we
	// passed in.
	rootKeyFilename := rootKeyID + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")
}

func TestImportExportGUN(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, newPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	_, err = cs.Create(data.CanonicalRootRole, data.ECDSAKey)
	_, err = cs.Create(data.CanonicalTargetsRole, data.ECDSAKey)
	_, err = cs.Create(data.CanonicalSnapshotRole, data.ECDSAKey)
	assert.NoError(t, err)

	tempZipFile, err := ioutil.TempFile("", "notary-test-export-")
	tempZipFilePath := tempZipFile.Name()
	defer os.Remove(tempZipFilePath)

	err = cs.ExportKeysByGUN(tempZipFile, gun, newPassphraseRetriever)
	assert.NoError(t, err)

	// With an invalid GUN, this should return an error
	err = cs.ExportKeysByGUN(tempZipFile, "does.not.exist/in/repository", newPassphraseRetriever)
	assert.EqualError(t, err, ErrNoKeysFoundForGUN.Error())

	tempZipFile.Close()

	// Reopen the zip file for importing
	zipReader, err := zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Map of files to expect in the zip file, with the passphrases
	passphraseByFile := make(map[string]string)

	// Add keys non-root keys to the map. These should use the new passphrase
	// because they were formerly unencrypted.
	privKeyMap := cs.ListAllKeys()
	for privKeyName := range privKeyMap {
		_, alias, err := cs.GetPrivateKey(privKeyName)
		assert.NoError(t, err, "privKey %s has no alias", privKeyName)
		if alias == "root" {
			continue
		}
		relKeyPath := filepath.Join("tuf_keys", privKeyName+".key")

		passphraseByFile[relKeyPath] = exportPassphrase
	}

	// Iterate through the files in the archive, checking that the files
	// exist and are encrypted with the expected passphrase.
	for _, f := range zipReader.File {

		expectedPassphrase, present := passphraseByFile[f.Name]
		assert.True(t, present, "unexpected file %s in zip file", f.Name)

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
	assert.Len(t, passphraseByFile, 0)

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, newPassphraseRetriever)
	cs2 := NewCryptoService(gun, fileStore2)

	// Reopen the zip file for importing
	zipReader, err = zip.OpenReader(tempZipFilePath)
	assert.NoError(t, err, "could not open zip file")

	// Now try with a valid passphrase. This time it should succeed.
	err = cs2.ImportKeysZip(zipReader.Reader)
	assert.NoError(t, err)
	zipReader.Close()

	// Look for keys in private. The filenames should match the key IDs
	// in the repo's private key store.
	for privKeyName, role := range privKeyMap {
		if role == "root" {
			continue
		}
		_, alias, err := cs2.GetPrivateKey(privKeyName)
		assert.NoError(t, err, "privKey %s has no alias", privKeyName)
		if alias == "root" {
			continue
		}
		relKeyPath := filepath.Join("tuf_keys", privKeyName+".key")
		privKeyFileName := filepath.Join(tempBaseDir2, "private", relKeyPath)
		_, err = os.Stat(privKeyFileName)
		assert.NoError(t, err)
	}
}

func TestImportExportRootKey(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, oldPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	pubKey, err := cs.Create(data.CanonicalRootRole, data.ECDSAKey)
	assert.NoError(t, err)

	rootKeyID := pubKey.ID()

	tempKeyFile, err := ioutil.TempFile("", "notary-test-export-")
	tempKeyFilePath := tempKeyFile.Name()
	defer os.Remove(tempKeyFilePath)

	err = cs.ExportKey(tempKeyFile, rootKeyID, data.CanonicalRootRole)
	assert.NoError(t, err)
	tempKeyFile.Close()

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, oldPassphraseRetriever)
	cs2 := NewCryptoService(gun, fileStore2)

	keyReader, err := os.Open(tempKeyFilePath)
	assert.NoError(t, err, "could not open key file")

	err = cs2.ImportRootKey(keyReader)
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
	decryptedPEMBytes, err := trustmanager.KeyToPEM(privKey, "root")
	assert.NoError(t, err, "could not convert key to PEM")

	err = cs2.ImportRootKey(bytes.NewReader(decryptedPEMBytes))
	assert.EqualError(t, err, ErrRootKeyNotEncrypted.Error())

	// Try to import garbage and make sure it doesn't succeed
	err = cs2.ImportRootKey(strings.NewReader("this is not PEM"))
	assert.EqualError(t, err, ErrNoValidPrivateKey.Error())

	// Should be able to unlock the root key with the old password
	key, alias, err := cs2.GetPrivateKey(rootKeyID)
	assert.NoError(t, err, "could not unlock root key")
	assert.Equal(t, "root", alias)
	assert.Equal(t, rootKeyID, key.ID())
}

func TestImportExportRootKeyReencrypt(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, oldPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	pubKey, err := cs.Create(data.CanonicalRootRole, data.ECDSAKey)
	assert.NoError(t, err)

	rootKeyID := pubKey.ID()

	tempKeyFile, err := ioutil.TempFile("", "notary-test-export-")
	tempKeyFilePath := tempKeyFile.Name()
	defer os.Remove(tempKeyFilePath)

	err = cs.ExportKeyReencrypt(tempKeyFile, rootKeyID, newPassphraseRetriever)
	assert.NoError(t, err)
	tempKeyFile.Close()

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, newPassphraseRetriever)
	cs2 := NewCryptoService(gun, fileStore2)

	keyReader, err := os.Open(tempKeyFilePath)
	assert.NoError(t, err, "could not open key file")

	err = cs2.ImportRootKey(keyReader)
	assert.NoError(t, err)
	keyReader.Close()

	// Look for repo's root key in repo2
	// There should be a file named after the key ID of the root key we
	// imported.
	rootKeyFilename := rootKeyID + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")

	// Should be able to unlock the root key with the new password
	key, alias, err := cs2.GetPrivateKey(rootKeyID)
	assert.NoError(t, err, "could not unlock root key")
	assert.Equal(t, "root", alias)
	assert.Equal(t, rootKeyID, key.ID())
}

func TestImportExportNonRootKey(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	fmt.Println(tempBaseDir)
	//defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, oldPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	pubKey, err := cs.Create(data.CanonicalTargetsRole, data.ECDSAKey)
	assert.NoError(t, err)

	targetsKeyID := pubKey.ID()

	tempKeyFile, err := ioutil.TempFile("", "notary-test-export-")
	tempKeyFilePath := tempKeyFile.Name()
	defer os.Remove(tempKeyFilePath)

	err = cs.ExportKey(tempKeyFile, targetsKeyID, data.CanonicalTargetsRole)
	assert.NoError(t, err)
	tempKeyFile.Close()

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	fmt.Println(tempBaseDir2)
	//defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, oldPassphraseRetriever)
	cs2 := NewCryptoService(gun, fileStore2)

	keyReader, err := os.Open(tempKeyFilePath)
	assert.NoError(t, err, "could not open key file")

	err = cs2.ImportRoleKey(keyReader, data.CanonicalTargetsRole, oldPassphraseRetriever)
	assert.NoError(t, err)
	keyReader.Close()

	// Look for repo's targets key in repo2
	// There should be a file named after the key ID of the targets key we
	// imported.
	targetsKeyFilename := targetsKeyID + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "tuf_keys", "docker.com/notary", targetsKeyFilename))
	assert.NoError(t, err, "missing targets key")

	// Check that the key is the same
	key, alias, err := cs2.GetPrivateKey(targetsKeyID)
	assert.NoError(t, err, "could not unlock targets key")
	assert.Equal(t, "targets", alias)
	assert.Equal(t, targetsKeyID, key.ID())
}

func TestImportExportNonRootKeyReencrypt(t *testing.T) {
	gun := "docker.com/notary"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore, err := trustmanager.NewKeyFileStore(tempBaseDir, oldPassphraseRetriever)
	cs := NewCryptoService(gun, fileStore)
	pubKey, err := cs.Create(data.CanonicalSnapshotRole, data.ECDSAKey)
	assert.NoError(t, err)

	snapshotKeyID := pubKey.ID()

	tempKeyFile, err := ioutil.TempFile("", "notary-test-export-")
	tempKeyFilePath := tempKeyFile.Name()
	defer os.Remove(tempKeyFilePath)

	err = cs.ExportKeyReencrypt(tempKeyFile, snapshotKeyID, newPassphraseRetriever)
	assert.NoError(t, err)
	tempKeyFile.Close()

	// Create new repo to test import
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	fileStore2, err := trustmanager.NewKeyFileStore(tempBaseDir2, newPassphraseRetriever)
	cs2 := NewCryptoService(gun, fileStore2)

	keyReader, err := os.Open(tempKeyFilePath)
	assert.NoError(t, err, "could not open key file")

	err = cs2.ImportRoleKey(keyReader, "snapshot", newPassphraseRetriever)
	assert.NoError(t, err)
	keyReader.Close()

	// Look for repo's snapshot key in repo2
	// There should be a file named after the key ID of the snapshot key we
	// imported.
	snapshotKeyFilename := snapshotKeyID + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir2, "private", "tuf_keys", "docker.com/notary", snapshotKeyFilename))
	assert.NoError(t, err, "missing snapshot key")

	// Should be able to unlock the root key with the new password
	key, alias, err := cs2.GetPrivateKey(snapshotKeyID)
	assert.NoError(t, err, "could not unlock snapshot key")
	assert.Equal(t, "snapshot", alias)
	assert.Equal(t, snapshotKeyID, key.ID())
}
