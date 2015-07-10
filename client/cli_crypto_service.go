package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

// CryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type CryptoService struct {
	gun        string
	passphrase string
	keyStore   *trustmanager.KeyFileStore
}

// NewCryptoService returns an instance of CryptoService
func NewCryptoService(gun string, keyStore *trustmanager.KeyFileStore) *CryptoService {
	return &CryptoService{gun: gun, keyStore: keyStore}
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

// SetPassphrase tells the cryptoservice the passphrase. Use only if the key needs
// to be decrypted.
func (ccs *CryptoService) SetPassphrase(passphrase string) {
	ccs.passphrase = passphrase
}

// Sign returns the signatures for data with the given root Key ID, falling back
// if not rootKeyID is found
func (ccs *CryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, fingerprint)

		var privKey *data.PrivateKey
		var err error
		var method string

		// Read PrivateKey from file
		if ccs.passphrase != "" {
			// This is a root key
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
			method = "RSASSA-PSS-X509"
		} else {
			privKey, err = ccs.keyStore.GetKey(keyName)
			method = "RSASSA-PSS"
		}
		if err != nil {
			// Note that GetDecryptedKey always fails on InitRepo.
			// InitRepo gets a signer that doesn't have access to
			// the root keys. Continuing here is safe because we
			// end up not returning any signatures.
			continue
		}

		sig, err := sign(privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    method,
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
