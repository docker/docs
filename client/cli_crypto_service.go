package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

type CryptoService struct {
	gun      string
	keyStore *trustmanager.KeyFileStore
}

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
	// Generates a new RSA key
	key, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, fmt.Errorf("could not generate private key: %v", err)
	}

	pemKey, err := trustmanager.KeyToPEM(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the certificate for key: %v (%s)", role, err)
	}
	rsaPublicKey := key.PublicKey
	// Using x509 to Marshal the Public key into der encoding
	pubBytes, err := x509.MarshalPKIXPublicKey(&rsaPublicKey)
	if err != nil {
		return nil, errors.New("Failed to Marshal public key.")
	}
	tufKey := data.NewPublicKey("RSA", pubBytes)

	// Passing in the the GUN + keyID as the name for the private key and adding it
	// to our KeyFileStore. Final storage will be under $BASE_PATH/GUN/keyID.key
	privKeyFilename := filepath.Join(ccs.gun, tufKey.ID())
	ccs.keyStore.Add(privKeyFilename, pemKey)

	return tufKey, nil
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
		privPEMBytes, err := ccs.keyStore.Get(privKeyFilename)
		if err != nil {
			continue
		}

		// Parse PrivateKey
		privKeyBytes, _ := pem.Decode(privPEMBytes)
		privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes.Bytes)
		if err != nil {
			return nil, err
		}

		// Sign the data
		sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:  fingerprint,
			Method: "RSA",
			//Method:    "RSASSA-PKCS1-V1_5-SIGN",
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
func (ccs *RootCryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// Read PrivateKey from file
		privPEMBytes, err := ccs.rootKeyStore.GetDecrypted(fingerprint, ccs.passphrase)
		if err != nil {
			// TODO(diogo): This error should be returned to the user in someway
			continue
		}

		// Parse PrivateKey
		privKeyBytes, _ := pem.Decode(privPEMBytes)
		privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes.Bytes)
		if err != nil {
			return nil, err
		}

		// Sign the data
		sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    "RSASSA-PKCS1-V1_5-SIGN",
			Signature: sig[:],
		})
	}
	return signatures, nil
}
