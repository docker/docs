package trustmanager

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/notary"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/tuf/data"
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
	gun := "docker.com/notary"
	testAddKeyWithRole(t, data.CanonicalRootRole, notary.RootKeysSubdir)
	testAddKeyWithRole(t, data.CanonicalTargetsRole, filepath.Join(notary.NonRootKeysSubdir, gun))
	testAddKeyWithRole(t, data.CanonicalSnapshotRole, filepath.Join(notary.NonRootKeysSubdir, gun))
	testAddKeyWithRole(t, "targets/a/b/c", notary.NonRootKeysSubdir)
	testAddKeyWithRole(t, "invalidRole", notary.NonRootKeysSubdir)
}

func testAddKeyWithRole(t *testing.T, role, expectedSubdir string) {
	gun := "docker.com/notary"
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)
	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, notary.PrivDir, expectedSubdir, privKey.ID()+"."+testExt)

	// Call the AddKey function
	err = store.AddKey(privKey, KeyInfo{Role: role, Gun: gun})
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	b, err := ioutil.ReadFile(expectedFilePath)
	assert.NoError(t, err, "expected file not found")
	assert.Contains(t, string(b), "-----BEGIN EC PRIVATE KEY-----")

	// Check that we have the role and gun info for this key's ID
	keyInfo, ok := store.keyInfoMap[privKey.ID()]
	assert.True(t, ok)
	assert.Equal(t, role, keyInfo.Role)
	if role == data.CanonicalRootRole || data.IsDelegation(role) || !data.ValidRole(role) {
		assert.Empty(t, keyInfo.Gun)
	} else {
		assert.Equal(t, gun, keyInfo.Gun)
	}
}

func TestKeyStoreInternalState(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	gun := "docker.com/notary"

	// Mimic a notary repo setup, and test that bringing up a keyfilestore creates the correct keyInfoMap
	roles := []string{data.CanonicalRootRole, data.CanonicalTargetsRole, data.CanonicalSnapshotRole, "targets/delegation"}
	// Keep track of the key IDs for each role, so we can validate later against the keystore state
	roleToID := make(map[string]string)
	for _, role := range roles {
		// generate a key for the role
		privKey, err := GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err, "could not generate private key")

		// generate the correct PEM role header
		privKeyPEM, err := KeyToPEM(privKey, role)
		assert.NoError(t, err, "could not generate PEM")

		// write the key file to the correct location
		// Pretend our GUN is docker.com/notary
		keyPath := filepath.Join(tempBaseDir, "private", getSubdir(role))
		if role == data.CanonicalTargetsRole || role == data.CanonicalSnapshotRole {
			keyPath = filepath.Join(keyPath, gun)
		}
		keyPath = filepath.Join(keyPath, privKey.ID())
		assert.NoError(t, os.MkdirAll(filepath.Dir(keyPath), 0755))
		assert.NoError(t, ioutil.WriteFile(keyPath+".key", privKeyPEM, 0755))

		roleToID[role] = privKey.ID()
	}

	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err)
	assert.Len(t, store.keyInfoMap, 4)
	for _, role := range roles {
		keyID, _ := roleToID[role]
		// make sure this keyID is the right length
		assert.Len(t, keyID, notary.Sha256HexSize)
		assert.Equal(t, role, store.keyInfoMap[keyID].Role)
		// targets and snapshot keys should have a gun set, root and delegation keys should not
		if role == data.CanonicalTargetsRole || role == data.CanonicalSnapshotRole {
			assert.Equal(t, gun, store.keyInfoMap[keyID].Gun)
		} else {
			assert.Empty(t, store.keyInfoMap[keyID].Gun)
		}
	}

	// Try removing the targets key only by ID (no gun provided)
	assert.NoError(t, store.RemoveKey(roleToID[data.CanonicalTargetsRole]))
	// The key file itself should have been removed
	_, err = os.Stat(filepath.Join(tempBaseDir, "private", "tuf_keys", gun, roleToID[data.CanonicalTargetsRole]+".key"))
	assert.Error(t, err)
	// The keyInfoMap should have also updated by deleting the key
	_, ok := store.keyInfoMap[roleToID[data.CanonicalTargetsRole]]
	assert.False(t, ok)

	// Try removing the delegation key only by ID (no gun provided)
	assert.NoError(t, store.RemoveKey(roleToID["targets/delegation"]))
	// The key file itself should have been removed
	_, err = os.Stat(filepath.Join(tempBaseDir, "private", "tuf_keys", roleToID["targets/delegation"]+".key"))
	assert.Error(t, err)
	// The keyInfoMap should have also updated
	_, ok = store.keyInfoMap[roleToID["targets/delegation"]]
	assert.False(t, ok)

	// Try removing the root key only by ID (no gun provided)
	assert.NoError(t, store.RemoveKey(roleToID[data.CanonicalRootRole]))
	// The key file itself should have been removed
	_, err = os.Stat(filepath.Join(tempBaseDir, "private", "root_keys", roleToID[data.CanonicalRootRole]+".key"))
	assert.Error(t, err)
	// The keyInfoMap should have also updated_
	_, ok = store.keyInfoMap[roleToID[data.CanonicalRootRole]]
	assert.False(t, ok)

	// Generate a new targets key and add it with its gun, check that the map gets updated back
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")
	assert.NoError(t, store.AddKey(privKey, KeyInfo{Role: data.CanonicalTargetsRole, Gun: gun}))
	assert.Equal(t, gun, store.keyInfoMap[privKey.ID()].Gun)
	assert.Equal(t, data.CanonicalTargetsRole, store.keyInfoMap[privKey.ID()].Role)
}

func TestGet(t *testing.T) {
	nonRootRolesToTest := []string{
		data.CanonicalTargetsRole,
		data.CanonicalSnapshotRole,
		"targets/a/b/c",
		"invalidRole",
	}

	gun := "docker.io/notary"

	// Root role needs to go in the rootKeySubdir to be read.
	// All other roles can go in the nonRootKeysSubdir, possibly under a GUN
	rootKeysSubdirWithGUN := filepath.Clean(filepath.Join(notary.RootKeysSubdir, gun))
	nonRootKeysSubdirWithGUN := filepath.Clean(filepath.Join(notary.NonRootKeysSubdir, gun))

	testGetKeyWithRole(t, "", data.CanonicalRootRole, notary.RootKeysSubdir, true)
	testGetKeyWithRole(t, gun, data.CanonicalRootRole, rootKeysSubdirWithGUN, true)
	for _, role := range nonRootRolesToTest {
		testGetKeyWithRole(t, "", role, notary.NonRootKeysSubdir, true)
		testGetKeyWithRole(t, gun, role, nonRootKeysSubdirWithGUN, true)
	}

	// Root cannot go in the nonRootKeysSubdir, or it won't be able to be read,
	// and vice versa
	testGetKeyWithRole(t, "", data.CanonicalRootRole, notary.NonRootKeysSubdir, false)
	testGetKeyWithRole(t, gun, data.CanonicalRootRole, nonRootKeysSubdirWithGUN, false)
	for _, role := range nonRootRolesToTest {
		testGetKeyWithRole(t, "", role, notary.RootKeysSubdir, false)
		testGetKeyWithRole(t, gun, role, rootKeysSubdirWithGUN, false)
	}
}

func testGetKeyWithRole(t *testing.T, gun, role, expectedSubdir string, success bool) {
	testData := []byte(fmt.Sprintf(`-----BEGIN RSA PRIVATE KEY-----
role: %s

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
`, role))
	testName := "keyID"
	testExt := "key"
	perms := os.FileMode(0755)

	emptyPassphraseRetriever := func(string, string, bool, int) (string, bool, error) { return "", false, nil }

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Since we're generating this manually we need to add the extension '.'
	filePath := filepath.Join(tempBaseDir, notary.PrivDir, expectedSubdir, testName+"."+testExt)
	os.MkdirAll(filepath.Dir(filePath), perms)
	err = ioutil.WriteFile(filePath, testData, perms)
	assert.NoError(t, err, "failed to write test file")

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, emptyPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Call the GetKey function
	if gun != "" {
		testName = gun + "/keyID"
	}
	privKey, _, err := store.GetKey(testName)
	if success {
		assert.NoError(t, err, "failed to get %s key from store (it's in %s)", role, expectedSubdir)

		pemPrivKey, err := KeyToPEM(privKey, role)
		assert.NoError(t, err, "failed to convert key to PEM")
		assert.Equal(t, testData, pemPrivKey)

		// Test that we can get purely by the ID we provided to AddKey (without gun)
		privKeyByID, _, err := store.GetKey("keyID")
		assert.NoError(t, err)
		assert.Equal(t, privKey, privKeyByID)
	} else {
		assert.Error(t, err, "should not have succeeded getting key from store")
		assert.Nil(t, privKey)
	}
}

// TestGetLegacyKey ensures we can still load keys where the role
// is stored as part of the filename (i.e. <hexID>_<role>.key
func TestGetLegacyKey(t *testing.T) {
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
	filePath := filepath.Join(tempBaseDir, notary.PrivDir, notary.RootKeysSubdir, testName+"_"+testAlias+"."+testExt)

	os.MkdirAll(filepath.Dir(filePath), perms)
	err = ioutil.WriteFile(filePath, testData, perms)
	assert.NoError(t, err, "failed to write test file")

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, emptyPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Call the GetKey function
	_, role, err := store.GetKey(testName)
	assert.NoError(t, err, "failed to get key from store")
	assert.Equal(t, testAlias, role)
}

func TestListKeys(t *testing.T) {
	testName := "docker.com/notary/root"
	perms := os.FileMode(0755)

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	roles := append(data.BaseRoles, "targets/a", "invalidRoleName")

	for i, role := range roles {
		// Make a new key for each role
		privKey, err := GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err, "could not generate private key")

		// Call the AddKey function
		gun := filepath.Dir(testName)
		err = store.AddKey(privKey, KeyInfo{Role: role, Gun: gun})
		assert.NoError(t, err, "failed to add key to store")

		// Check to see if the keystore lists this key
		keyMap := store.ListKeys()

		// Expect to see exactly one key in the map
		assert.Len(t, keyMap, i+1)
		// Expect to see privKeyID inside of the map
		listedInfo, ok := keyMap[privKey.ID()]
		assert.True(t, ok)
		assert.Equal(t, role, listedInfo.Role)
	}

	// Write an invalid filename to the directory
	filePath := filepath.Join(tempBaseDir, notary.PrivDir, notary.RootKeysSubdir, "fakekeyname.key")
	err = ioutil.WriteFile(filePath, []byte("data"), perms)
	assert.NoError(t, err, "failed to write test file")

	// Check to see if the keystore still lists two keys
	keyMap := store.ListKeys()
	assert.Len(t, keyMap, len(roles))
}

func TestAddGetKeyMemStore(t *testing.T) {
	testAlias := data.CanonicalRootRole

	// Create our store
	store := NewKeyMemoryStore(passphraseRetriever)

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(privKey, KeyInfo{Role: testAlias, Gun: ""})
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	retrievedKey, retrievedAlias, err := store.GetKey(privKey.ID())
	assert.NoError(t, err, "failed to get key from store")

	assert.Equal(t, retrievedAlias, testAlias)
	assert.Equal(t, retrievedKey.Public(), privKey.Public())
	assert.Equal(t, retrievedKey.Private(), privKey.Private())
}

func TestAddGetKeyInfoMemStore(t *testing.T) {
	gun := "docker.com/notary"

	// Create our store
	store := NewKeyMemoryStore(passphraseRetriever)

	rootKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(rootKey, KeyInfo{Role: data.CanonicalRootRole, Gun: ""})
	assert.NoError(t, err, "failed to add key to store")

	// Get and validate key info
	rootInfo, err := store.GetKeyInfo(rootKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, data.CanonicalRootRole, rootInfo.Role)
	assert.Equal(t, "", rootInfo.Gun)

	targetsKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(targetsKey, KeyInfo{Role: data.CanonicalTargetsRole, Gun: gun})
	assert.NoError(t, err, "failed to add key to store")

	// Get and validate key info
	targetsInfo, err := store.GetKeyInfo(targetsKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, data.CanonicalTargetsRole, targetsInfo.Role)
	assert.Equal(t, gun, targetsInfo.Gun)

	delgKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(delgKey, KeyInfo{Role: "targets/delegation", Gun: gun})
	assert.NoError(t, err, "failed to add key to store")

	// Get and validate key info
	delgInfo, err := store.GetKeyInfo(delgKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, "targets/delegation", delgInfo.Role)
	assert.Equal(t, "", delgInfo.Gun)
}

func TestGetDecryptedWithTamperedCipherText(t *testing.T) {
	testExt := "key"
	testAlias := data.CanonicalRootRole

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
	err = store.AddKey(privKey, KeyInfo{Role: testAlias, Gun: ""})
	assert.NoError(t, err, "failed to add key to store")

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, notary.PrivDir, notary.RootKeysSubdir, privKey.ID()+"."+testExt)

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
	testAlias := data.CanonicalRootRole

	// Generate a new random RSA Key
	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Call the AddKey function
	err = store.AddKey(privKey, KeyInfo{Role: testAlias, Gun: ""})
	assert.NoError(t, err, "failed to add key to store")

	// Try to decrypt the file with an invalid passphrase
	_, _, err = newStore.GetKey(privKey.ID())
	assert.Error(t, err, "expected error while decrypting the content due to invalid passphrase")
	assert.IsType(t, err, expectedFailureType)
}

func TestRemoveKey(t *testing.T) {
	gun := "docker.com/notary"
	testRemoveKeyWithRole(t, data.CanonicalRootRole, notary.RootKeysSubdir)
	testRemoveKeyWithRole(t, data.CanonicalTargetsRole, filepath.Join(notary.NonRootKeysSubdir, gun))
	testRemoveKeyWithRole(t, data.CanonicalSnapshotRole, filepath.Join(notary.NonRootKeysSubdir, gun))
	testRemoveKeyWithRole(t, "targets/a/b/c", notary.NonRootKeysSubdir)
	testRemoveKeyWithRole(t, "invalidRole", notary.NonRootKeysSubdir)
}

func testRemoveKeyWithRole(t *testing.T, role, expectedSubdir string) {
	gun := "docker.com/notary"
	testExt := "key"

	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory")
	defer os.RemoveAll(tempBaseDir)

	// Create our store
	store, err := NewKeyFileStore(tempBaseDir, passphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	privKey, err := GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate private key")

	// Since we're generating this manually we need to add the extension '.'
	expectedFilePath := filepath.Join(tempBaseDir, notary.PrivDir, expectedSubdir, privKey.ID()+"."+testExt)

	err = store.AddKey(privKey, KeyInfo{Role: role, Gun: gun})
	assert.NoError(t, err, "failed to add key to store")

	// Check to see if file exists
	_, err = ioutil.ReadFile(expectedFilePath)
	assert.NoError(t, err, "expected file not found")

	// Call remove key
	err = store.RemoveKey(privKey.ID())
	assert.NoError(t, err, "unable to remove key")

	// Check to see if file still exists
	_, err = ioutil.ReadFile(expectedFilePath)
	assert.Error(t, err, "file should not exist")
}

func TestKeysAreCached(t *testing.T) {
	gun := "docker.com/notary"
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
	err = store.AddKey(privKey, KeyInfo{Role: testAlias, Gun: gun})
	assert.NoError(t, err, "failed to add key to store")

	assert.Equal(t, 1, numTimesCalled, "numTimesCalled should have been 1")

	// Call the AddKey function
	privKey2, _, err := store.GetKey(privKey.ID())
	assert.NoError(t, err, "failed to add key to store")

	assert.Equal(t, privKey.Public(), privKey2.Public(), "cachedPrivKey should be the same as the added privKey")
	assert.Equal(t, privKey.Private(), privKey2.Private(), "cachedPrivKey should be the same as the added privKey")
	assert.Equal(t, 1, numTimesCalled, "numTimesCalled should be 1 -- no additional call to passphraseRetriever")

	// Create a new store
	store2, err := NewKeyFileStore(tempBaseDir, countingPassphraseRetriever)
	assert.NoError(t, err, "failed to create new key filestore")

	// Call the GetKey function
	privKey3, _, err := store2.GetKey(privKey.ID())
	assert.NoError(t, err, "failed to get key from store")

	assert.Equal(t, privKey2.Private(), privKey3.Private(), "privkey from store1 should be the same as privkey from store2")
	assert.Equal(t, privKey2.Public(), privKey3.Public(), "privkey from store1 should be the same as privkey from store2")
	assert.Equal(t, 2, numTimesCalled, "numTimesCalled should be 2 -- one additional call to passphraseRetriever")

	// Call the GetKey function a bunch of times
	for i := 0; i < 10; i++ {
		_, _, err := store2.GetKey(privKey.ID())
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
	err = store.AddKey(privKey, KeyInfo{Role: data.CanonicalRootRole, Gun: ""})
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
	err = store.AddKey(privKey, KeyInfo{Role: data.CanonicalRootRole, Gun: ""})
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

	pemBytes, err := EncryptPrivateKey(expectedKey, "root", cannedPassphrase)
	assert.NoError(t, err)

	err = s.ImportKey(pemBytes, "root")
	assert.NoError(t, err)

	reimportedKey, reimportedAlias, err := s.GetKey(expectedKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, "root", reimportedAlias)
	assert.Equal(t, expectedKey.Private(), reimportedKey.Private())
	assert.Equal(t, expectedKey.Public(), reimportedKey.Public())
}
