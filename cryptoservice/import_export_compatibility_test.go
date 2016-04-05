// Ensures we can import/export old-style repos

package cryptoservice

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

// Zips up the keys in the old repo, and assert that we can import it and use
// said keys.  The 0.1 exported format is just a zip file of all the keys
func TestImport0Dot1Zip(t *testing.T) {
	ks, ret, _ := get0Dot1(t)

	zipFile, err := ioutil.TempFile("", "notary-test-zipFile")
	defer os.RemoveAll(zipFile.Name())
	zipWriter := zip.NewWriter(zipFile)
	require.NoError(t, err)
	require.NoError(t, addKeysToArchive(zipWriter, ks))
	zipWriter.Close()
	zipFile.Close()

	origKeys := make(map[string]string)
	for keyID, keyInfo := range ks.ListKeys() {
		origKeys[keyID] = keyInfo.Role
	}
	require.Len(t, origKeys, 3)

	// now import the zip file into a new cryptoservice

	tempDir, err := ioutil.TempDir("", "notary-test-import")
	defer os.RemoveAll(tempDir)
	require.NoError(t, err)

	ks, err = trustmanager.NewKeyFileStore(tempDir, ret)
	require.NoError(t, err)
	cs := NewCryptoService(ks)

	zipReader, err := zip.OpenReader(zipFile.Name())
	require.NoError(t, err)
	defer zipReader.Close()

	require.NoError(t, cs.ImportKeysZip(zipReader.Reader, passphrase.ConstantRetriever("randompass")))
	assertHasKeys(t, cs, origKeys)
}

func get0Dot1(t *testing.T) (*trustmanager.KeyFileStore, passphrase.Retriever, string) {
	gun := "docker.com/notary0.1/samplerepo"
	ret := passphrase.ConstantRetriever("randompass")

	// produce the zip file
	ks, err := trustmanager.NewKeyFileStore("../fixtures/compatibility/notary0.1", ret)
	require.NoError(t, err)

	return ks, ret, gun
}

// Given a map of key IDs to roles, asserts that the cryptoService has all and
// only those keys
func assertHasKeys(t *testing.T, cs *CryptoService, expectedKeys map[string]string) {
	keys := cs.ListAllKeys()
	require.Len(t, keys, len(expectedKeys))

	for keyID, role := range keys {
		expectedRole, ok := expectedKeys[keyID]
		require.True(t, ok)
		require.Equal(t, expectedRole, role)
	}
}

// Export all the keys of a cryptoservice to a zipfile, and import it into a
// new cryptoService, and return that new cryptoService
func importExportedZip(t *testing.T, original *CryptoService,
	ret passphrase.Retriever, gun string) (*CryptoService, string) {

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	require.NoError(t, err, "failed to create a temporary directory: %s", err)

	ks, err := trustmanager.NewKeyFileStore(tempBaseDir, ret)
	require.NoError(t, err)
	var cs *CryptoService

	// export keys
	zipFile, err := ioutil.TempFile("", "notary-test-zipFile")
	defer os.RemoveAll(zipFile.Name())
	if gun != "" {
		err = original.ExportKeysByGUN(zipFile, gun, ret)
		require.NoError(t, err)
		cs = NewCryptoService(ks)
	} else {
		err = original.ExportAllKeys(zipFile, ret)
		require.NoError(t, err)
		cs = NewCryptoService(ks)
	}
	zipFile.Close()

	// import keys into the cryptoservice now
	zipReader, err := zip.OpenReader(zipFile.Name())
	require.NoError(t, err)
	defer zipReader.Close()

	require.NoError(t, cs.ImportKeysZip(zipReader.Reader, passphrase.ConstantRetriever("randompass")))
	return cs, tempBaseDir
}

func TestImportExport0Dot1AllKeys(t *testing.T) {
	ks, ret, _ := get0Dot1(t)
	cs := NewCryptoService(ks)

	newCS, tempDir := importExportedZip(t, cs, ret, "")
	defer os.RemoveAll(tempDir)

	assertHasKeys(t, newCS, cs.ListAllKeys())
}

func TestImportExport0Dot1GUNKeys(t *testing.T) {
	ks, ret, gun := get0Dot1(t)

	// remove root from expected key list, because root is not exported when
	// we export by gun
	expectedKeys := make(map[string]string)
	for keyID, keyInfo := range ks.ListKeys() {
		if keyInfo.Role != data.CanonicalRootRole {
			expectedKeys[keyID] = keyInfo.Role
		}
	}

	// make some other temp directory to create new keys in
	tempDir, err := ioutil.TempDir("", "notary-tests-keystore")
	defer os.RemoveAll(tempDir)
	require.NoError(t, err)

	otherKS, err := trustmanager.NewKeyFileStore(tempDir, ret)
	require.NoError(t, err)
	cs := NewCryptoService(otherKS, ks)

	// create a keys that is not of the same GUN, and be sure it's in this
	// CryptoService
	otherPubKey, err := cs.Create(data.CanonicalTargetsRole, "some/other/gun", data.ECDSAKey)
	require.NoError(t, err)

	k, _, err := cs.GetPrivateKey(otherPubKey.ID())
	require.NoError(t, err)
	require.NotNil(t, k)

	// export/import, and ensure that the other-gun key is not in the new
	// CryptoService
	newCS, tempDir := importExportedZip(t, cs, ret, gun)
	defer os.RemoveAll(tempDir)

	assertHasKeys(t, newCS, expectedKeys)

	_, _, err = newCS.GetPrivateKey(otherPubKey.ID())
	require.Error(t, err)
}
