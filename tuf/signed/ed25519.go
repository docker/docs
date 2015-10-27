package signed

import (
	"crypto/rand"
	"errors"

	"github.com/agl/ed25519"
	"github.com/docker/notary/tuf/data"
)

// Ed25519 implements a simple in memory cryptosystem for ED25519 keys
type Ed25519 struct {
	keys map[string]data.PrivateKey
}

// NewEd25519 initializes a new empty Ed25519 CryptoService that operates
// entirely in memory
func NewEd25519() *Ed25519 {
	return &Ed25519{
		make(map[string]data.PrivateKey),
	}
}

// addKey allows you to add a private key
func (e *Ed25519) addKey(k data.PrivateKey) {
	e.keys[k.ID()] = k
}

// RemoveKey deletes a key from the signer
func (e *Ed25519) RemoveKey(keyID string) error {
	delete(e.keys, keyID)
	return nil
}

// Sign generates an Ed25519 signature over the data
func (e *Ed25519) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		priv := [ed25519.PrivateKeySize]byte{}
		copy(priv[:], e.keys[keyID].Private())
		sig := ed25519.Sign(&priv, toSign)
		signatures = append(signatures, data.Signature{
			KeyID:     keyID,
			Method:    data.EDDSASignature,
			Signature: sig[:],
		})
	}
	return signatures, nil

}

// Create generates a new key and returns the public part
func (e *Ed25519) Create(role string, algorithm data.KeyAlgorithm) (data.PublicKey, error) {
	if algorithm != data.ED25519Key {
		return nil, errors.New("only ED25519 supported by this cryptoservice")
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	public := data.NewPublicKey(data.ED25519Key, pub[:])
	private := data.NewPrivateKey(data.ED25519Key, pub[:], priv[:])
	e.addKey(private)
	return public, nil
}

// PublicKeys returns a map of public keys for the ids provided, when those IDs are found
// in the store.
func (e *Ed25519) PublicKeys(keyIDs ...string) (map[string]data.PublicKey, error) {
	k := make(map[string]data.PublicKey)
	for _, keyID := range keyIDs {
		if key, ok := e.keys[keyID]; ok {
			k[keyID] = data.PublicKeyFromPrivate(key)
		}
	}
	return k, nil
}

// GetKey returns a single public key based on the ID
func (e *Ed25519) GetKey(keyID string) data.PublicKey {
	return data.PublicKeyFromPrivate(e.keys[keyID])
}
