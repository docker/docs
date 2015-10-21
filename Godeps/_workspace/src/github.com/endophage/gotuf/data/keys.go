package data

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Sirupsen/logrus"
	"github.com/jfrazelle/go/canonical/json"
)

// Key is the minimal interface for a public key. It is declared
// independently of PublicKey for composability.
type Key interface {
	ID() string
	Algorithm() KeyAlgorithm
	Public() []byte
}

// PublicKey is the necessary interface for public keys
type PublicKey interface {
	Key
}

// PrivateKey adds the ability to access the private key
type PrivateKey interface {
	Key

	Private() []byte
}

// KeyPair holds the public and private key bytes
type KeyPair struct {
	Public  []byte `json:"public"`
	Private []byte `json:"private"`
}

// TUFKey is the structure used for both public and private keys in TUF.
// Normally it would make sense to use a different structures for public and
// private keys, but that would change the key ID algorithm (since the canonical
// JSON would be different). This structure should normally be accessed through
// the PublicKey or PrivateKey interfaces.
type TUFKey struct {
	id    string
	Type  KeyAlgorithm `json:"keytype"`
	Value KeyPair      `json:"keyval"`
}

// NewPrivateKey instantiates a new TUFKey with the private key component
// populated
func NewPrivateKey(algorithm KeyAlgorithm, public, private []byte) *TUFKey {
	return &TUFKey{
		Type: algorithm,
		Value: KeyPair{
			Public:  public,
			Private: private,
		},
	}
}

// Algorithm returns the algorithm of the key
func (k TUFKey) Algorithm() KeyAlgorithm {
	return k.Type
}

// ID efficiently generates if necessary, and caches the ID of the key
func (k *TUFKey) ID() string {
	if k.id == "" {
		pubK := NewPublicKey(k.Algorithm(), k.Public())
		data, err := json.MarshalCanonical(&pubK)
		if err != nil {
			logrus.Error("Error generating key ID:", err)
		}
		digest := sha256.Sum256(data)
		k.id = hex.EncodeToString(digest[:])
	}
	return k.id
}

// Public returns the public bytes
func (k TUFKey) Public() []byte {
	return k.Value.Public
}

// Private returns the private bytes
func (k TUFKey) Private() []byte {
	return k.Value.Private
}

// NewPublicKey instantiates a new TUFKey where the private bytes are
// guaranteed to be nil
func NewPublicKey(algorithm KeyAlgorithm, public []byte) PublicKey {
	return &TUFKey{
		Type: algorithm,
		Value: KeyPair{
			Public:  public,
			Private: nil,
		},
	}
}

// PublicKeyFromPrivate returns a new TUFKey based on a private key, with
// the private key bytes guaranteed to be nil.
func PublicKeyFromPrivate(pk PrivateKey) PublicKey {
	return &TUFKey{
		Type: pk.Algorithm(),
		Value: KeyPair{
			Public:  pk.Public(),
			Private: nil,
		},
	}
}
