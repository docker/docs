package cryptoservice

import (
	"crypto/rand"
	"crypto/x509"
	"testing"
	"time"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCertificate(t *testing.T) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err, "could not generate key")

	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)

	err = keyStore.AddKey(trustmanager.KeyInfo{Role: data.CanonicalRootRole, Gun: ""}, privKey)
	assert.NoError(t, err, "could not add key to store")

	// Check GenerateCertificate method
	gun := "docker.com/notary"
	startTime := time.Now()
	cert, err := GenerateCertificate(privKey, gun, startTime, startTime.AddDate(10, 0, 0))
	assert.NoError(t, err, "could not generate certificate")

	// Check public key
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(privKey.Private())
	assert.NoError(t, err)
	ecdsaPublicKey := ecdsaPrivateKey.Public()
	assert.Equal(t, ecdsaPublicKey, cert.PublicKey)

	// Check CommonName
	assert.Equal(t, cert.Subject.CommonName, gun)
}
