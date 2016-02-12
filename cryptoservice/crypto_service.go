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
	gun       string
	keyStores []trustmanager.KeyStore
}

// NewCryptoService returns an instance of CryptoService
func NewCryptoService(gun string, keyStores ...trustmanager.KeyStore) *CryptoService {
	return &CryptoService{gun: gun, keyStores: keyStores}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (cs *CryptoService) Create(role, algorithm string) (data.PublicKey, error) {
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

	// Store the private key into our keystore
	for _, ks := range cs.keyStores {
		err = ks.AddKey(privKey, trustmanager.KeyInfo{Role: role, Gun: cs.gun})
		if err == nil {
			return data.PublicKeyFromPrivate(privKey), nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to add key to filestore: %v", err)
	}

	return nil, fmt.Errorf("keystores would not accept new private keys for unknown reasons")
}

// GetPrivateKey returns a private key and role if present by ID.
func (cs *CryptoService) GetPrivateKey(keyID string) (k data.PrivateKey, role string, err error) {
	for _, ks := range cs.keyStores {
		k, role, err = ks.GetKey(keyID)
		if err == nil {
			return
		}
		switch err.(type) {
		case trustmanager.ErrPasswordInvalid, trustmanager.ErrAttemptsExceeded:
			return
		default:
			continue
		}
	}
	return // returns whatever the final values were
}

// GetKey returns a key by ID
func (cs *CryptoService) GetKey(keyID string) data.PublicKey {
	privKey, _, err := cs.GetPrivateKey(keyID)
	if err != nil {
		return nil
	}
	return data.PublicKeyFromPrivate(privKey)
}

// GetKeyInfo returns role and GUN info of a key by ID
func (cs *CryptoService) GetKeyInfo(keyID string) (trustmanager.KeyInfo, error) {
	for _, store := range cs.keyStores {
		if info, err := store.GetKeyInfo(keyID); err == nil {
			return info, nil
		}
	}
	return trustmanager.KeyInfo{}, fmt.Errorf("Could not find info for keyID %s", keyID)
}

// RemoveKey deletes a key by ID
func (cs *CryptoService) RemoveKey(keyID string) (err error) {
	for _, ks := range cs.keyStores {
		ks.RemoveKey(keyID)
	}
	return // returns whatever the final values were
}

// AddKey adds a private key to a specified role.
// The GUN is inferred from the cryptoservice itself for non-root roles
func (cs *CryptoService) AddKey(key data.PrivateKey, role string) (err error) {
	keyID := key.ID()
	if role != data.CanonicalRootRole {
		keyID = filepath.Join(cs.gun, key.ID())
	}
	for _, ks := range cs.keyStores {
		if keyInfo, err := ks.GetKeyInfo(keyID); err == nil {
			if keyInfo.Role != role {
				return fmt.Errorf("key with same ID already exists for role: %s", keyInfo.Role)
			}
			continue
		}
		ks.AddKey(key, trustmanager.KeyInfo{Role: role, Gun: cs.gun})
	}
	return // returns whatever the final values were
}

// ListKeys returns a list of key IDs valid for the given role
func (cs *CryptoService) ListKeys(role string) []string {
	var res []string
	for _, ks := range cs.keyStores {
		for k, r := range ks.ListKeys() {
			if r.Role == role {
				res = append(res, k)
			}
		}
	}
	return res
}

// ListAllKeys returns a map of key IDs to role
func (cs *CryptoService) ListAllKeys() map[string]string {
	res := make(map[string]string)
	for _, ks := range cs.keyStores {
		for k, r := range ks.ListKeys() {
			res[k] = r.Role // keys are content addressed so don't care about overwrites
		}
	}
	return res
}
