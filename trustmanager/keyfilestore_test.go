package trustmanager

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/notary/pkg/passphrase"
	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
)

const cannedPassphrase = "passphrase"

var passphraseRetriever = func(keyID string, alias string, createNew bool, numAttempts int) (string, bool, error) {
	if numAttempts > 5 {
		giveup := true
		return "", giveup, errors.New("passPhraseRetriever failed after too many requests")
	}
	return cannedPassphrase, false, nil
}

func TestAddKey(t *testing.T) {
	testName := "docker.com/notary/root"
	testExt := "key"
	testAlias := "root"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, rootKeysSubdir, testName+"_"+testAlias+"."+testExt)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(testName, "root", privKey)
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	assert.NoError(t, err, "expected file not found")
	assert.Contains(t, string(b), "-----BEGIN EC PRIVATE KEY-----")
}

func TestGet(t *testing.T) {
	testData := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAyUIXjsrWRrvPa4Bzp3VJ6uOUGPay2fUpSV8XzNxZxIG/Opdr
+k3EQi1im6WOqF3Y5AS1UjYRxNuRN+cAZeo3uS1pOTuoSupBXuchVw8s4hZJ5vXn
TRmGb+xY7tZ1ZVgPfAZDib9sRSUsL/gC+aSyprAjG/YBdbF06qKbfOfsoCEYW1OQ
82JqHzQH514RFYPTnEGpvfxWaqmFQLmv0uMxV/cAYvqtrGkXuP0+a8PknlD2obw5
0rHE56Su1c3Q42S7L51K38tpbgWOSRcTfDUWEj5v9wokkNQvyKBwbS996s4EJaZd
7r6M0h1pHnuRxcSaZLYRwgOe1VNGg2VfWzgd5QIDAQABAoIBAF9LGwpygmj1jm3R
YXGd+ITugvYbAW5wRb9G9mb6wspnwNsGTYsz/UR0ZudZyaVw4jx8+jnV/i3e5PC6
QRcAgqf8l4EQ/UuThaZg/AlT1yWp9g4UyxNXja87EpTsGKQGwTYxZRM4/xPyWOzR
mt8Hm8uPROB9aA2JG9npaoQG8KSUj25G2Qot3ukw/IOtqwN/Sx1EqF0EfCH1K4KU
a5TrqlYDFmHbqT1zTRec/BTtVXNsg8xmF94U1HpWf3Lpg0BPYT7JiN2DPoLelRDy
a/A+a3ZMRNISL5wbq/jyALLOOyOkIqa+KEOeW3USuePd6RhDMzMm/0ocp5FCwYfo
k4DDeaECgYEA0eSMD1dPGo+u8UTD8i7ZsZCS5lmXLNuuAg5f5B/FGghD8ymPROIb
dnJL5QSbUpmBsYJ+nnO8RiLrICGBe7BehOitCKi/iiZKJO6edrfNKzhf4XlU0HFl
jAOMa975pHjeCoZ1cXJOEO9oW4SWTCyBDBSqH3/ZMgIOiIEk896lSmkCgYEA9Xf5
Jqv3HtQVvjugV/axAh9aI8LMjlfFr9SK7iXpY53UdcylOSWKrrDok3UnrSEykjm7
UL3eCU5jwtkVnEXesNn6DdYo3r43E6iAiph7IBkB5dh0yv3vhIXPgYqyTnpdz4pg
3yPGBHMPnJUBThg1qM7k6a2BKHWySxEgC1DTMB0CgYAGvdmF0J8Y0k6jLzs/9yNE
4cjmHzCM3016gW2xDRgumt9b2xTf+Ic7SbaIV5qJj6arxe49NqhwdESrFohrKaIP
kM2l/o2QaWRuRT/Pvl2Xqsrhmh0QSOQjGCYVfOb10nAHVIRHLY22W4o1jk+piLBo
a+1+74NRaOGAnu1J6/fRKQKBgAF180+dmlzemjqFlFCxsR/4G8s2r4zxTMXdF+6O
3zKuj8MbsqgCZy7e8qNeARxwpCJmoYy7dITNqJ5SOGSzrb2Trn9ClP+uVhmR2SH6
AlGQlIhPn3JNzI0XVsLIloMNC13ezvDE/7qrDJ677EQQtNEKWiZh1/DrsmHr+irX
EkqpAoGAJWe8PC0XK2RE9VkbSPg9Ehr939mOLWiHGYTVWPttUcum/rTKu73/X/mj
WxnPWGtzM1pHWypSokW90SP4/xedMxludvBvmz+CTYkNJcBGCrJumy11qJhii9xp
EMl3eFOJXjIch/wIesRSN+2dGOsl7neercjMh1i9RvpCwHDx/E0=
-----END RSA PRIVATE KEY-----
`)
	testName := "docker.com/notary/root"
	testExt := "key"
	testAlias := "root"
	perms := os.FileMode(0755)

	emptyPassphraseRetriever := func(string, string, bool, int) (string, bool, error) { return "", false, nil }

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	filePath := filepath.Join(tempBaseDir, rootKeysSubdir, testName+"_"+testAlias+"."+testExt)

	os.MkdirAll(filepath.Dir(filePath), perms)
	err = ioutil.WriteFile(filePath, testData, perms)
	assert.NoError(t, err, "failed to write test file")

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, emptyPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Call the GetKey function
	privKey, _, err := store.GetKey(testName)
	assert.NoError(t, err, "failed to get key from store")

	pemPrivKey, err := KeyToPEM(privKey)
	assert.NoError(t, err, "failed to convert key to PEM")
	assert.Equal(t, testData, pemPrivKey)
}

func TestAddGetKeyMemStore(t *testing.T) {
	testName := "docker.com/notary/root"
	testAlias := "root"

	// Create our store
	store := NewKeyMemoryStore(passphraseRetriever)

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(testName, testAlias, privKey)
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	retrievedKey, retrievedAlias, err := store.GetKey(testName)
	assert.NoError(t, err, "failed to get key from store")

	assert.Equal(t, retrievedAlias, testAlias)
	assert.Equal(t, retrievedKey.Public(), privKey.Public())
	assert.Equal(t, retrievedKey.Private(), privKey.Private())
}
func TestGetDecryptedWithTamperedCipherText(t *testing.T) {
	testExt := "key"
	testAlias := "root"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Generate a new Private Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddEncryptedKey function
	err = store.AddKey(privKey.ID(), testAlias, privKey)
	assert.NoError(t, err, "failed to add key to store")

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, rootKeysSubdir, privKey.ID()+"_"+testAlias+"."+testExt)

	// Get file description, open file
	fp, err := os.OpenFile(expectedFilePath, os.O_WRONLY, 0600)
	assert.NoError(t, err, "expected file not found")

	// Tamper the file
	fp.WriteAt([]byte("a"), int64(1))

	// Recreate the KeyFileStore to avoid caching
	store, err = NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Try to decrypt the file
	_, _, err = store.GetKey(privKey.ID())
	assert.Error(t, err, "expected error while decrypting the content due to invalid cipher text")
}

func TestGetDecryptedWithInvalidPassphrase(t *testing.T) {

	// Make a passphraseRetriever that always returns a different passphrase in order to test
	// decryption failure
	a := "a"
	var invalidPassphraseRetriever = func(keyId string, alias string, createNew bool, numAttempts int) (string, bool, error) {
		if numAttempts > 5 {
			giveup := true
			return "", giveup, nil
		}
		a = a + a
		return a, false, nil
	}

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Test with KeyFileStore
	fileStore, err := NewKeyFileStore(tempBaseDir, invalidPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	newFileStore, err := NewKeyFileStore(tempBaseDir, invalidPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	testGetDecryptedWithInvalidPassphrase(t, fileStore, newFileStore, ErrPasswordInvalid{})

	// Can't test with KeyMemoryStore because we cache the decrypted version of
	// the key forever
}

func TestGetDecryptedWithConsistentlyInvalidPassphrase(t *testing.T) {
	// Make a passphraseRetriever that always returns a different passphrase in order to test
	// decryption failure
	a := "aaaaaaaaaaaaa"
	var consistentlyInvalidPassphraseRetriever = func(keyID string, alias string, createNew bool, numAttempts int) (string, bool, error) {
		a = a + "a"
		return a, false, nil
	}

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Test with KeyFileStore
	fileStore, err := NewKeyFileStore(tempBaseDir, consistentlyInvalidPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	newFileStore, err := NewKeyFileStore(tempBaseDir, consistentlyInvalidPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	testGetDecryptedWithInvalidPassphrase(t, fileStore, newFileStore, ErrAttemptsExceeded{})

	// Can't test with KeyMemoryStore because we cache the decrypted version of
	// the key forever
}

// testGetDecryptedWithInvalidPassphrase takes two keystores so it can add to
// one and get from the other (to work around caching)
func testGetDecryptedWithInvalidPassphrase(t *testing.T, store KeyStore, newStore KeyStore, expectedFailureType interface{}) {
	testAlias := "root"

	// Generate a new random RSA Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(privKey.ID(), testAlias, privKey)
	assert.NoError(t, err, "failed to add key to store")

	// Try to decrypt the file with an invalid passphrase
	_, _, err = newStore.GetKey(privKey.ID())
	assert.Error(t, err, "expected error while decrypting the content due to invalid passphrase")
	assert.IsType(t, err, expectedFailureType)
}

func TestRemoveKey(t *testing.T) {
	testName := "docker.com/notary/root"
	testExt := "key"
	testAlias := "alias"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, nonRootKeysSubdir, testName+"_"+testAlias+"."+testExt)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(testName, testAlias, privKey)
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	_, err = ioutil.ReadFile(expectedFilePath)
	assert.NoError(t, err, "expected file not found")

	// Call remove key
	err = store.RemoveKey(testName)
	assert.NoError(t, err, "unable to remove key")

	// Check to see if file still exists
	_, err = ioutil.ReadFile(expectedFilePath)
	assert.Error(t, err, "file should not exist")
}

func TestKeysAreCached(t *testing.T) {
	testName := "docker.com/notary/root"
	testAlias := "alias"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	var countingPassphraseRetriever passphrase.Retriever

	numTimesCalled := 0
	countingPassphraseRetriever = func(keyId, alias string, createNew bool, attempts int) (passphrase string, giveup bool, err error) {
		numTimesCalled++
		return "password", false, nil
	}

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, countingPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(testName, testAlias, privKey)
	assert.NoError(t, err, "failed to add key to store")

	assert.Equal(t, 1, numTimesCalled, "numTimesCalled should have been 1")

	// Call the AddKey function
	privKey2, _, err := store.GetKey(testName)
	assert.NoError(t, err, "failed to add key to store")

	assert.Equal(t, privKey.Public(), privKey2.Public(), "cachedPrivKey should be the same as the added privKey")
	assert.Equal(t, privKey.Private(), privKey2.Private(), "cachedPrivKey should be the same as the added privKey")
	assert.Equal(t, 1, numTimesCalled, "numTimesCalled should be 1 -- no additional call to passphraseRetriever")

	// Create a new store
	store2, err := NewKeyFileStore(tempBaseDir, countingPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Call the GetKey function
	privKey3, _, err := store2.GetKey(testName)
	assert.NoError(t, err, "failed to get key from store")

	assert.Equal(t, privKey2.Private(), privKey3.Private(), "privkey from store1 should be the same as privkey from store2")
	assert.Equal(t, privKey2.Public(), privKey3.Public(), "privkey from store1 should be the same as privkey from store2")
	assert.Equal(t, 2, numTimesCalled, "numTimesCalled should be 2 -- one additional call to passphraseRetriever")

	// Call the GetKey function a bunch of times
	for i := 0; i < 10; i++ {
		_, _, err := store2.GetKey(testName)
		assert.NoError(t, err, "failed to get key from store")
	}
	assert.Equal(t, 2, numTimesCalled, "numTimesCalled should be 2 -- no additional call to passphraseRetriever")
}

// Exporting a key is successful (it is a valid key)
func TestKeyFileStoreExportSuccess(t *testing.T) {
	// Generate a new Private Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore and add the key
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err)
	err = store.AddKey(privKey.ID(), "root", privKey)
	assert.NoError(t, err)

	assertExportKeySuccess(t, store, privKey)
}

// Exporting a key that doesn't exist fails (it is a valid key)
func TestKeyFileStoreExportNonExistantFailure(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Create empty FileStore
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err)

	_, err = store.ExportKey("12345")
	assert.Error(t, err)
}

// Exporting a key is successful (it is a valid key)
func TestKeyMemoryStoreExportSuccess(t *testing.T) {
	// Generate a new Private Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// Create our MemoryStore and add key to it
	store := NewKeyMemoryStore(passphraseRetriever)
	assert.NoError(t, err)
	err = store.AddKey(privKey.ID(), "root", privKey)
	assert.NoError(t, err)

	assertExportKeySuccess(t, store, privKey)
}

// Exporting a key that doesn't exist fails (it is a valid key)
func TestKeyMemoryStoreExportNonExistantFailure(t *testing.T) {
	store := NewKeyMemoryStore(passphraseRetriever)
	_, err := store.ExportKey("12345")
	assert.Error(t, err)
}

// Importing a key is successful
func TestKeyFileStoreImportSuccess(t *testing.T) {
	// Generate a new Private Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err)

	assertImportKeySuccess(t, store, privKey)
}

// Importing a key is successful
func TestKeyMemoryStoreImportSuccess(t *testing.T) {
	// Generate a new Private Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	// Create our MemoryStore
	store := NewKeyMemoryStore(passphraseRetriever)
	assert.NoError(t, err)

	assertImportKeySuccess(t, store, privKey)
}

// Given a keystore and expected key that is in the store, export the key
// and assert that the exported key is the same and encrypted with the right
// password.
func assertExportKeySuccess(
	t *testing.T, s KeyStore, expectedKey data.PrivateKey) {

	pemBytes, err := s.ExportKey(expectedKey.ID())
	assert.NoError(t, err)

	reparsedKey, err := ParsePEMPrivateKey(pemBytes, cannedPassphrase)
	assert.NoError(t, err)
	assert.Equal(t, expectedKey.Private(), reparsedKey.Private())
	assert.Equal(t, expectedKey.Public(), reparsedKey.Public())
}

// Given a keystore and expected key, generate an encrypted PEM of the key
// and assert that the then imported key is the same and encrypted with the
// right password.
func assertImportKeySuccess(
	t *testing.T, s KeyStore, expectedKey data.PrivateKey) {

	pemBytes, err := EncryptPrivateKey(expectedKey, cannedPassphrase)
	assert.NoError(t, err)

	err = s.ImportKey(pemBytes, "root")
	assert.NoError(t, err)

	reimportedKey, reimportedAlias, err := s.GetKey(expectedKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, "root", reimportedAlias)
	assert.Equal(t, expectedKey.Private(), reimportedKey.Private())
	assert.Equal(t, expectedKey.Public(), reimportedKey.Public())
}
