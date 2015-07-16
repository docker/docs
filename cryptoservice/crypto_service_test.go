package cryptoservice

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
)

func TestED25519(t *testing.T) {
	content := []byte("this is a secret")

	keyStore := trustmanager.NewKeyMemoryStore()
	cryptoService := NewCryptoService("", keyStore, "")

	tufKey, err := cryptoService.Create("", data.ED25519Key)
	assert.NoError(t, err, "error creating key")

	signatures, err := cryptoService.Sign([]string{tufKey.ID()}, content)
	assert.NoError(t, err, "signing failed")
	assert.Len(t, signatures, 1, "wrong number of signatures")

	verifier := &signed.Ed25519Verifier{}
	err = verifier.Verify(tufKey, signatures[0].Signature, content)
	assert.NoError(t, err, "verification failed")
}

func TestRSA(t *testing.T) {
	content := []byte("this is a secret")

	keyStore := trustmanager.NewKeyMemoryStore()
	cryptoService := NewCryptoService("", keyStore, "")

	tufKey, err := cryptoService.Create("", data.RSAKey)
	assert.NoError(t, err, "error creating key")

	signatures, err := cryptoService.Sign([]string{tufKey.ID()}, content)
	assert.NoError(t, err, "signing failed")
	assert.Len(t, signatures, 1, "wrong number of signatures")

	verifier := &signed.RSAPSSVerifier{}
	err = verifier.Verify(tufKey, signatures[0].Signature, content)
	assert.NoError(t, err, "verification failed")
}

func TestECDSA(t *testing.T) {
	content := []byte("this is a secret")

	keyStore := trustmanager.NewKeyMemoryStore()
	cryptoService := NewCryptoService("", keyStore, "")

	tufKey, err := cryptoService.Create("", data.ECDSAKey)
	assert.NoError(t, err, "error creating key")

	signatures, err := cryptoService.Sign([]string{tufKey.ID()}, content)
	assert.NoError(t, err, "signing failed")
	assert.Len(t, signatures, 1, "wrong number of signatures")

	verifier := &signed.ECDSAVerifier{}
	err = verifier.Verify(tufKey, signatures[0].Signature, content)
	assert.NoError(t, err, "verification failed")
}
