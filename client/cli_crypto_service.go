package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"path/filepath"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

type CliCryptoService struct {
	gun      string
	keyStore trustmanager.FileStore
}

// NewCryptoService returns an instance ofS cliCryptoService
func NewCryptoService(gun string, keyStore trustmanager.FileStore) *CliCryptoService {
	return &CliCryptoService{gun: gun, keyStore: keyStore}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *CliCryptoService) Create(role string) (*data.PublicKey, error) {
	keyData, pemCert, err := GenerateKeyAndCert(ccs.gun)
	if err != nil {
		return nil, err
	}

	fingerprint, err := trustmanager.FingerprintPEMCert(pemCert)
	if err != nil {
		return nil, err
	}

	// The key is going to be stored in the private directory, using the GUN and
	// the filename will be the TUF-compliant ID. The Store takes care of extensions.
	privKeyFilename := filepath.Join(ccs.gun, fingerprint)

	// Store this private key
	ccs.keyStore.Add(privKeyFilename, keyData)

	return data.NewPublicKey("RSA", pemCert), nil
}

// Sign returns the signatures for data with the given keyIDs
func (ccs *CliCryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
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
			KeyID:     fingerprint,
			Method:    "RSASSA-PKCS1-V1_5-SIGN",
			Signature: sig[:],
		})
	}
	return signatures, nil
}

// generateKeyAndCert deals with the creation and storage of a key and returns a
// PEM encoded cert
func GenerateKeyAndCert(gun string) ([]byte, []byte, error) {
	// Generates a new RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate private key: %v", err)
	}

	// Creates a new Certificate template. We need the certificate to calculate the
	// TUF-compliant keyID
	//TODO (diogo): We're hardcoding the Organization to be the GUN. Probably want to
	// change it
	template := trustmanager.NewCertificate(gun, gun)
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate the certificate for key: %v", err)
	}

	// Encode the new certificate into PEM
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate the certificate for key: %v", err)
	}

	pemKey, err := trustmanager.KeyToPEM(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate the certificate for key: %v", err)
	}

	return pemKey, trustmanager.CertToPEM(cert), nil
}
