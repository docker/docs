// +build pkcs11

package yubikey

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

var ret = passphrase.ConstantRetriever("passphrase")

// create a new store for clearing out keys, because we don't want to pollute
// any cache
func clearAllKeys(t *testing.T) {
	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	for k := range store.ListKeys() {
		err := store.RemoveKey(k)
		assert.NoError(t, err)
	}
}

func TestEnsurePrivateKeySizePassesThroughRightSizeArrays(t *testing.T) {
	fullByteArray := make([]byte, ecdsaPrivateKeySize)
	for i := range fullByteArray {
		fullByteArray[i] = byte(1)
	}

	result := ensurePrivateKeySize(fullByteArray)
	assert.True(t, reflect.DeepEqual(fullByteArray, result))
}

// The pad32Byte helper function left zero-pads byte arrays that are less than
// ecdsaPrivateKeySize bytes
func TestEnsurePrivateKeySizePadsLessThanRequiredSizeArrays(t *testing.T) {
	shortByteArray := make([]byte, ecdsaPrivateKeySize/2)
	for i := range shortByteArray {
		shortByteArray[i] = byte(1)
	}

	expected := append(
		make([]byte, ecdsaPrivateKeySize-ecdsaPrivateKeySize/2),
		shortByteArray...)

	result := ensurePrivateKeySize(shortByteArray)
	assert.True(t, reflect.DeepEqual(expected, result))
}

func testAddKey(t *testing.T, store trustmanager.KeyStore) (data.PrivateKey, error) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	err = store.AddKey(privKey.ID(), data.CanonicalRootRole, privKey)
	return privKey, err
}

func addMaxKeys(t *testing.T, store trustmanager.KeyStore) []string {
	keys := make([]string, 0, numSlots)
	// create the maximum number of keys
	for i := 0; i < numSlots; i++ {
		privKey, err := testAddKey(t, store)
		assert.NoError(t, err)
		keys = append(keys, privKey.ID())
	}
	return keys
}

// We can add keys enough times to fill up all the slots in the Yubikey.
// They are backed up, and we can then list them and get the keys.
func TestYubiAddKeysAndRetrieve(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	keys := addMaxKeys(t, origStore)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// All 4 keys should be in the original store, in the clean store (which
	// makes sure the keys are actually on the Yubikey and not on the original
	// store's cache, and on the backup store)
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		listedKeys := store.ListKeys()
		assert.Len(t, listedKeys, numSlots)
		for _, k := range keys {
			r, ok := listedKeys[k]
			assert.True(t, ok)
			assert.Equal(t, data.CanonicalRootRole, r)

			_, _, err := store.GetKey(k)
			assert.NoError(t, err)
		}
	}
}

// We can't add a key if there are no more slots
func TestYubiAddKeyFailureIfNoMoreSlots(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	addMaxKeys(t, origStore)

	// add another key - should fail because there are no more slots
	badKey, err := testAddKey(t, origStore)
	assert.Error(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// The key should not be in the original store, in the new clean store, or
	// in teh backup store.
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		// the key that wasn't created should not appear in ListKeys or GetKey
		_, _, err := store.GetKey(badKey.ID())
		assert.Error(t, err)
		for k := range store.ListKeys() {
			assert.NotEqual(t, badKey, k)
		}
	}
}

// If some random key in the middle was removed, adding a key will work (keys
// do not have to be deleted/added in order)
func TestYubiAddKeyCanAddToMiddleSlot(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	// create 4 keys on the original store
	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)
	keys := addMaxKeys(t, origStore)

	// delete one of the middle keys, and assert we can still create a new key
	keyIDToDelete := keys[numSlots/2]
	err = origStore.RemoveKey(keyIDToDelete)
	assert.NoError(t, err)

	newKey, err := testAddKey(t, origStore)
	assert.NoError(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// The new key should be in the original store, in the new clean store, and
	// in the backup store.  The old key should not be in the original store,
	// or the new clean store.
	for _, store := range []trustmanager.KeyStore{origStore, cleanStore, backup} {
		// new key should appear in all stores
		gottenKey, _, err := store.GetKey(newKey.ID())
		assert.NoError(t, err)
		assert.Equal(t, gottenKey.ID(), newKey.ID())

		listedKeys := store.ListKeys()
		_, ok := listedKeys[newKey.ID()]
		assert.True(t, ok)

		// old key should not be in the non-backup stores
		if store != backup {
			_, _, err := store.GetKey(keyIDToDelete)
			assert.Error(t, err)
			_, ok = listedKeys[keyIDToDelete]
			assert.False(t, ok)
		}
	}
}

// RemoveKey removes a key from the yubikey, but not from the backup store.
func TestYubiRemoveKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	key, err := testAddKey(t, origStore)
	assert.NoError(t, err)
	err = origStore.RemoveKey(key.ID())
	assert.NoError(t, err)

	// key remains in the backup store
	backupKey, role, err := backup.GetKey(key.ID())
	assert.NoError(t, err)
	assert.Equal(t, data.CanonicalRootRole, role)
	assert.Equal(t, key.ID(), backupKey.ID())

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	// key is not in either the original store or the clean store
	for _, store := range []*YubiKeyStore{origStore, cleanStore} {
		_, _, err := store.GetKey(key.ID())
		assert.Error(t, err)
	}
}

// ImportKey imports a key as root without adding it to the backup store
func TestYubiImportKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	origStore, err := NewYubiKeyStore(backup, ret)
	assert.NoError(t, err)

	// generate key and import it
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)

	pemBytes, err := trustmanager.EncryptPrivateKey(privKey, "passphrase")
	assert.NoError(t, err)

	err = origStore.ImportKey(pemBytes, "root")
	assert.NoError(t, err)

	// key is not in backup store
	_, _, err = backup.GetKey(privKey.ID())
	assert.Error(t, err)

	// create a new store, since we want to be sure the original store's cache
	// is not masking any issues
	cleanStore, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	for _, store := range []*YubiKeyStore{origStore, cleanStore} {
		gottenKey, role, err := store.GetKey(privKey.ID())
		assert.NoError(t, err)
		assert.Equal(t, data.CanonicalRootRole, role)
		assert.Equal(t, privKey.Public(), gottenKey.Public())
	}
}

// One cannot export from hardware - it will not export from the backup
func TestYubiExportKeyFails(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	key, err := testAddKey(t, store)
	assert.NoError(t, err)

	_, err = store.ExportKey(key.ID())
	assert.Error(t, err)
}

// If there are keys in the backup store but no keys in the Yubikey,
// listing and getting cannot access the keys in the backup store
func TestYubiListAndGetKeysIgnoresBackup(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	backup := trustmanager.NewKeyMemoryStore(ret)
	key, err := testAddKey(t, backup)
	assert.NoError(t, err)

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.Len(t, store.ListKeys(), 0)
	_, _, err = store.GetKey(key.ID())
	assert.Error(t, err)
}

// Get a YubiPrivateKey.  Check that it has the right algorithm, etc, and
// specifically that you cannot get the private bytes out.
func TestYubiKey(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	ecdsaPrivateKey, err := testAddKey(t, store)
	assert.NoError(t, err)

	yubiPrivateKey, _, err := store.GetKey(ecdsaPrivateKey.ID())
	assert.NoError(t, err)

	assert.Equal(t, data.ECDSAKey, yubiPrivateKey.Algorithm())
	assert.Equal(t, data.ECDSASignature, yubiPrivateKey.SignatureAlgorithm())
	assert.Equal(t, ecdsaPrivateKey.Public(), yubiPrivateKey.Public())
	assert.Nil(t, yubiPrivateKey.Private())
}

// Get a YubiPrivateKey.  Sign something with it.
func TestYubiSigning(t *testing.T) {
	// TODO(cyli): the signature should be verified, but the importing the
	// verifiers causes an import cycle.  A bigger refactor needs to be done
	// to fix it.
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	store, err := NewYubiKeyStore(trustmanager.NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	ecdsaPrivateKey, err := testAddKey(t, store)
	assert.NoError(t, err)

	yubiPrivateKey, _, err := store.GetKey(ecdsaPrivateKey.ID())
	assert.NoError(t, err)

	msg := []byte("Hello there")
	_, err = yubiPrivateKey.Sign(bytes.NewBuffer(msg), msg, nil)
	assert.NoError(t, err)
}
