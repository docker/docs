package cryptoservice

import (
	"crypto/rand"
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
)

const (
	rsaKeySize = 2048 // Used for snapshots and targets keys
)

// CryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type CryptoService struct {
	gun      string
	keyStore trustmanager.KeyStore
}

// NewCryptoService returns an instance of CryptoService
func NewCryptoService(gun string, keyStore trustmanager.KeyStore) *CryptoService {
	return &CryptoService{gun: gun, keyStore: keyStore}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *CryptoService) Create(role, algorithm string) (data.PublicKey, error) {
	var privKey data.PrivateKey
	var err error

	switch algorithm {
	case data.RSAKey:
		privKey, err = trustmanager.GenerateRSAKey(rand.Reader, rsaKeySize)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %v", err)
		}
	case data.ECDSAKey:
		privKey, err = trustmanager.GenerateECDSAKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate EC key: %v", err)
		}
	case data.ED25519Key:
		privKey, err = trustmanager.GenerateED25519Key(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ED25519 key: %v", err)
		}
	default:
		return nil, fmt.Errorf("private key type not supported for key generation: %s", algorithm)
	}
	logrus.Debugf("generated new %s key for role: %s and keyID: %s", algorithm, role, privKey.ID())

	// Store the private key into our keystore with the name being: /GUN/ID.key with an alias of role
	err = ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), role, privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to add key to filestore: %v", err)
	}
	return data.PublicKeyFromPrivate(privKey), nil
}

// GetPrivateKey returns a private key by ID
func (ccs *CryptoService) GetPrivateKey(keyID string) (data.PrivateKey, string, error) {
	return ccs.keyStore.GetKey(keyID)
}

// GetKey returns a key by ID
func (ccs *CryptoService) GetKey(keyID string) data.PublicKey {
	key, _, err := ccs.keyStore.GetKey(keyID)
	if err != nil {
		return nil
	}
	return data.PublicKeyFromPrivate(key)
}

// RemoveKey deletes a key by ID
func (ccs *CryptoService) RemoveKey(keyID string) error {
	return ccs.keyStore.RemoveKey(keyID)
}

// Sign returns the signatures for the payload with a set of keyIDs. It ignores
// errors to sign and expects the called to validate if the number of returned
// signatures is adequate.
func (ccs *CryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		keyName := keyid

		// Try to get the key first without a GUN (in which case it's a root
		// key).  If that fails, try to get the key with the GUN (non-root
		// key).  If that fails, then we don't have the key.
		privKey, _, err := ccs.keyStore.GetKey(keyName)
		if err != nil {
			keyName = filepath.Join(ccs.gun, keyid)
			privKey, _, err = ccs.keyStore.GetKey(keyName)
			if err != nil {
				logrus.Debugf("error attempting to retrieve key ID: %s, %v", keyid, err)
				return nil, err
			}
		}

		var sigAlgorithm data.SigAlgorithm

		switch privKey.(type) {
		case *data.RSAPrivateKey:
			sigAlgorithm = data.RSAPSSSignature
		case *data.ECDSAPrivateKey:
			sigAlgorithm = data.ECDSASignature
		case *data.ED25519PrivateKey:
			sigAlgorithm = data.EDDSASignature
		}
		sig, err := privKey.Sign(rand.Reader, payload, nil)
		if err != nil {
			logrus.Debugf("ignoring error attempting to %s sign with keyID: %s, %v", privKey.Algorithm(), keyid, err)
			return nil, err
		}

		logrus.Debugf("appending %s signature with Key ID: %s", privKey.Algorithm(), keyid)

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    sigAlgorithm,
			Signature: sig[:],
		})
	}

	return signatures, nil
}
