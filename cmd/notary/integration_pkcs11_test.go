// +build pkcs11

package main

import (
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

var rootOnHardware = trustmanager.YubikeyAccessible

// Per-test set up that returns a cleanup function.  This set up:
// - changes the passphrase retriever to always produce a constant passphrase
// - disables touch on yubikeys
// - deletes all keys on the yubikey
func setUp(t *testing.T) func() {
	oldRetriever := retriever

	var fake = func(k, a string, c bool, n int) (string, bool, error) {
		if k == "Yubikey" {
			return oldRetriever(k, a, c, n)
		}
		return testPassphrase, false, nil
	}

	retriever = fake
	getRetriever = func() passphrase.Retriever { return fake }
	trustmanager.SetYubikeyKeyMode(trustmanager.KeymodeNone)

	// //we're just removing keys here, so nil is fine
	s, err := trustmanager.NewYubiKeyStore(nil, retriever)
	assert.NoError(t, err)
	for k := range s.ListKeys() {
		err := s.RemoveKey(k)
		assert.NoError(t, err)
	}

	return func() {
		retriever = oldRetriever
		getRetriever = getPassphraseRetriever
		trustmanager.SetYubikeyKeyMode(trustmanager.KeymodeTouch | trustmanager.KeymodePinOnce)
	}
}

// ensures that the root is actually on the yubikey - this makes sure the
// commands are hooked up to interact with the yubikey, rather than right files
// on disk
func verifyRootKeyOnHardware(t *testing.T, rootKeyID string) {
	// do not bother verifying if there is no yubikey available
	if trustmanager.YubikeyAccessible() {
		// //we're just getting keys here, so nil is fine
		s, err := trustmanager.NewYubiKeyStore(nil, retriever)
		assert.NoError(t, err)
		privKey, role, err := s.GetKey(rootKeyID)
		assert.NoError(t, err)
		assert.NotNil(t, privKey)
		assert.Equal(t, data.CanonicalRootRole, role)
	}
}
