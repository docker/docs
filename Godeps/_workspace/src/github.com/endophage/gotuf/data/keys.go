package data

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Sirupsen/logrus"
	cjson "github.com/tent/canonical-json-go"
)

type Key interface {
	ID() string
	Algorithm() KeyAlgorithm
	Public() []byte
	Private() []byte
}

type KeyPair struct {
	Public  []byte `json:"public"`
	Private []byte `json:"private"`
}

type TUFKey struct {
	id    string       `json:"-"`
	Type  KeyAlgorithm `json:"keytype"`
	Value KeyPair      `json:"keyval"`
}

func NewTUFKey(algorithm KeyAlgorithm, public, private []byte) *TUFKey {
	return &TUFKey{
		Type: algorithm,
		Value: KeyPair{
			Public:  public,
			Private: private,
		},
	}
}

func (k TUFKey) Algorithm() KeyAlgorithm {
	return k.Type
}

func (k *TUFKey) ID() string {
	if k.id == "" {
		pubK := NewTUFKey(k.Algorithm(), k.Public(), nil)
		data, err := cjson.Marshal(&pubK)
		if err != nil {
			logrus.Error("Error generating key ID:", err)
		}
		digest := sha256.Sum256(data)
		k.id = hex.EncodeToString(digest[:])
	}
	return k.id
}

func (k TUFKey) Public() []byte {
	return k.Value.Public
}

type PublicKey struct {
	TUFKey
}

func (k PublicKey) Private() []byte {
	return nil
}

func NewPublicKey(algorithm KeyAlgorithm, public []byte) *PublicKey {
	return &PublicKey{
		TUFKey{
			Type: algorithm,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

func PublicKeyFromPrivate(pk PrivateKey) *PublicKey {
	return &PublicKey{
		pk.TUFKey,
	}
}

type PrivateKey struct {
	TUFKey
}

func NewPrivateKey(algorithm KeyAlgorithm, public, private []byte) *PrivateKey {
	return &PrivateKey{
		TUFKey{
			Type: algorithm,
			Value: KeyPair{
				Public:  []byte(public),
				Private: []byte(private),
			},
		},
	}
}

func (k PrivateKey) Private() []byte {
	return k.Value.Private
}
