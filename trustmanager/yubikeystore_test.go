// +build pkcs11

package trustmanager

import (
	"crypto/rand"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

func clearAllKeys(t *testing.T) {
	// TODO(cyli): this is creating a new yubikey store because for some reason,
	// removing and then adding with the same YubiKeyStore causes
	// non-deterministic failures at least on Mac OS
	ret := passphrase.ConstantRetriever("passphrase")
	store, err := NewYubiKeyStore(NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)

	for k := range store.ListKeys() {
		err := store.RemoveKey(k)
		assert.NoError(t, err)
	}
}

func TestAddKeyToNextEmptyYubikeySlot(t *testing.T) {
	if !YubikeyAccessible() {
		t.Skip("Must have Yubikey access.")
	}
	clearAllKeys(t)

	ret := passphrase.ConstantRetriever("passphrase")
	store, err := NewYubiKeyStore(NewKeyMemoryStore(ret), ret)
	assert.NoError(t, err)
	SetYubikeyKeyMode(KeymodeNone)
	defer func() {
		SetYubikeyKeyMode(KeymodeTouch | KeymodePinOnce)
	}()

	keys := make([]string, 0, numSlots)

	// create the maximum number of keys
	for i := 0; i < numSlots; i++ {
		privKey, err := GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err)

		err = store.AddKey(privKey.ID(), data.CanonicalRootRole, privKey)
		assert.NoError(t, err)

		keys = append(keys, privKey.ID())
	}

	listedKeys := store.ListKeys()
	assert.Len(t, listedKeys, numSlots)
	for _, k := range keys {
		r, ok := listedKeys[k]
		assert.True(t, ok)
		assert.Equal(t, data.CanonicalRootRole, r)
	}

	// numSlots is not actually the max - some keys might have more, so do not
	// test that adding more keys will fail.
}
