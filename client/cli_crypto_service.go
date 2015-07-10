package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

// CryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type CryptoService struct {
	gun      string
	keyStore *trustmanager.KeyFileStore
}

// RootCryptoService implements Sign and Create and operates on a rootKeyStore,
// taking in a passphrase and calling decrypt when signing.
type RootCryptoService struct {
	// TODO(diogo): support multiple passphrases per key
	passphrase   string
	rootKeyStore *trustmanager.KeyFileStore
}

// NewCryptoService returns an instance of CryptoService
func NewCryptoService(gun string, keyStore *trustmanager.KeyFileStore) *CryptoService {
	return &CryptoService{gun: gun, keyStore: keyStore}
}

// NewRootCryptoService returns an instance of CryptoService
func NewRootCryptoService(rootKeyStore *trustmanager.KeyFileStore, passphrase string) *RootCryptoService {
	return &RootCryptoService{rootKeyStore: rootKeyStore, passphrase: passphrase}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *CryptoService) Create(role string) (*data.PublicKey, error) {
	privKey, err := trustmanager.GenerateRSAKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Store the private key into our keystore with the name being: /GUN/ID.key
	ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)

	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for data with the given keyIDs
func (ccs *CryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// Get the PrivateKey filename
		privKeyFilename := filepath.Join(ccs.gun, fingerprint)
		// Read PrivateKey from file
		privKey, err := ccs.keyStore.GetKey(privKeyFilename)
		if err != nil {
			continue
		}

		sig, err := sign(privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    "RSASSA-PSS",
			Signature: sig[:],
		})
	}
	return signatures, nil
}

// Create in a root crypto service is not implemented
func (rcs *RootCryptoService) Create(role string) (*data.PublicKey, error) {
	return nil, errors.New("create on a root key filestore is not implemented")
}

// Sign returns the signatures for data with the given root Key ID, falling back
// if not rootKeyID is found
// TODO(diogo): This code has 1 line change from the Sign from Crypto service. DRY it up.
func (rcs *RootCryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// Read PrivateKey from file
		privKey, err := rcs.rootKeyStore.GetDecryptedKey(fingerprint, rcs.passphrase)
		if err != nil {
			// TODO(diogo): This error should be returned to the user in someway
			continue
		}

		sig, err := sign(privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    "RSASSA-PSS-X509",
			Signature: sig[:],
		})
	}

	return signatures, nil
}

func sign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	// TODO(diogo): Implement support for ECDSA.
	if strings.ToLower(privKey.Cipher()) != "rsa" {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an rsa.PrivateKey out of the private key bytes
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the RSA key to sign the data
	sig, err := rsa.SignPSS(rand.Reader, rsaPrivKey, hash, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		return nil, err
	}

	return sig, nil
}
