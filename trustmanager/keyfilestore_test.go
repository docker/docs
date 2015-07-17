package trustmanager

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddKey(t *testing.T) {
	testName := "docker.com/notary/root"
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	privKey, err := GenerateRSAKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("could not generate private key: %v", err)
	}

	// Call the AddKey function
	err = store.AddKey(testName, privKey)
	if err != nil {
		t.Fatalf("failed to add file to store: %v", err)
	}

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected file not found: %v", err)
	}

	if !strings.Contains(string(b), "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatalf("expected private key content in the file: %s", expectedFilePath)
	}
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
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	filePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	os.MkdirAll(filepath.Dir(filePath), perms)
	if err = ioutil.WriteFile(filePath, testData, perms); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	// Call the GetKey function
	privKey, err := store.GetKey(testName)
	if err != nil {
		t.Fatalf("failed to get file from store: %v", err)
	}

	pemPrivKey, err := KeyToPEM(privKey)
	if err != nil {
		t.Fatalf("failed to convert key to PEM: %v", err)
	}

	if !bytes.Equal(testData, pemPrivKey) {
		t.Fatalf("unexpected content in the file: %s", filePath)
	}
}

func TestAddEncryptedAndGetDecrypted(t *testing.T) {
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	// Generate new PrivateKey
	privKey, err := GenerateRSAKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("could not generate private key: %v", err)
	}

	// Call the AddEncryptedKey function
	err = store.AddEncryptedKey(privKey.ID(), privKey, "diogomonica")
	if err != nil {
		t.Fatalf("failed to add file to store: %v", err)
	}

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, privKey.ID()+"."+testExt)

	// Check to see if file exists
	_, err = ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected file not found: %v", err)
	}

	// Call the GetDecryptedKey function
	readPrivKey, err := store.GetDecryptedKey(privKey.ID(), "diogomonica")
	if err != nil {
		t.Fatalf("could not decrypt private key: %v", err)
	}

	if !bytes.Equal(privKey.Private(), readPrivKey.Private()) {
		t.Fatalf("written key and loaded key do not match")
	}
}

func TestGetDecryptedWithTamperedCipherText(t *testing.T) {
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	// Generate a new Private Key
	privKey, err := GenerateRSAKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("could not generate private key: %v", err)
	}

	// Call the AddEncryptedKey function
	err = store.AddEncryptedKey(privKey.ID(), privKey, "diogomonica")
	if err != nil {
		t.Fatalf("failed to add file to store: %v", err)
	}

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, privKey.ID()+"."+testExt)

	// Get file description, open file
	fp, err := os.OpenFile(expectedFilePath, os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("expected file not found: %v", err)
	}

	// Tamper the file
	fp.WriteAt([]byte("a"), int64(1))

	// Try to decrypt the file
	_, err = store.GetDecryptedKey(privKey.ID(), "diogomonica")
	if err == nil {
		t.Fatalf("expected error while decrypting the content due to invalid cipher text")
	}
}

func TestGetDecryptedWithInvalidPassphrase(t *testing.T) {
	testName := "docker.com/notary/root"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Create our FileStore
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	// Generate a new random RSA Key
	privKey, err := GenerateRSAKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("could not generate private key: %v", err)
	}

	// Call the AddEncryptedKey function
	err = store.AddEncryptedKey(privKey.ID(), privKey, "diogomonica")
	if err != nil {
		t.Fatalf("failed to add file to stoAFre: %v", err)
	}

	// Try to decrypt the file with an invalid passphrase
	_, err = store.GetDecryptedKey(testName, "diegomonica")
	if err == nil {
		t.Fatalf("expected error while decrypting the content due to invalid passphrase")
	}
}

func TestRemoveKey(t *testing.T) {
	testName := "docker.com/notary/root"
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, testName+"."+testExt)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir)
	if err != nil {
		t.Fatalf("failed to create new key filestore: %v", err)
	}

	privKey, err := GenerateRSAKey(rand.Reader, 512)
	if err != nil {
		t.Fatalf("could not generate private key: %v", err)
	}

	// Call the AddKey function
	err = store.AddKey(testName, privKey)
	if err != nil {
		t.Fatalf("failed to add file to store: %v", err)
	}

	// Check to see if file exists
	_, err = ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("expected file not found: %v", err)
	}

	// Call remove key
	err = store.RemoveKey(testName)
	if err != nil {
		t.Fatalf("unable to remove key: %v", err)
	}

	// Check to see if file still exists
	_, err = ioutil.ReadFile(expectedFilePath)
	if err == nil {
		t.Fatalf("file should not exist %s", expectedFilePath)
	}
}
