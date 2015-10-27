package cryptoservice

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

func TestCryptoService(t *testing.T) {
	testCryptoService(t, data.ECDSAKey, signed.ECDSAVerifier{})
	testCryptoService(t, data.ED25519Key, signed.Ed25519Verifier{})
	if !testing.Short() {
		testCryptoService(t, data.RSAKey, signed.RSAPSSVerifier{})
	}
}

var passphraseRetriever = func(string, string, bool, int) (string, bool, error) { return "", false, nil }

func testCryptoService(t *testing.T, keyAlgo data.KeyAlgorithm, verifier signed.Verifier) {
	content := []byte("this is a secret")

	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := NewCryptoService("", keyStore)

	// Test Create
	tufKey, err := cryptoService.Create("", keyAlgo)
	assert.NoError(t, err, "error creating key")

	// Test Sign
	signatures, err := cryptoService.Sign([]string{tufKey.ID()}, content)
	assert.NoError(t, err, "signing failed")
	assert.Len(t, signatures, 1, "wrong number of signatures")

	err = verifier.Verify(tufKey, signatures[0].Signature, content)
	assert.NoError(t, err, "verification failed")

	// Test GetKey
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.Equal(t, tufKey.Public(), retrievedKey.Public(), "retrieved key didn't match")

	assert.Nil(t, cryptoService.GetKey("boguskeyid"), "non-nil result for bogus keyid")

	// Test RemoveKey
	err = cryptoService.RemoveKey(tufKey.ID())
	assert.NoError(t, err, "could not remove key")
	retrievedKey = cryptoService.GetKey(tufKey.ID())
	assert.Nil(t, retrievedKey, "remove didn't work")
}
