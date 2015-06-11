package data

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Sirupsen/logrus"
	cjson "github.com/tent/canonical-json-go"
)

type Key interface {
	ID() string
	Cipher() string
	Public() string
	Private() string
}

type KeyPair struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

type TUFKey struct {
	id    string  `json:"-"`
	Type  string  `json:"keytype"`
	Value KeyPair `json:"keyval"`
}

func NewTUFKey(cipher, public, private string) *TUFKey {
	return &TUFKey{
		Type: cipher,
		Value: KeyPair{
			Public:  public,
			Private: private,
		},
	}
}

func (k TUFKey) Cipher() string {
	return k.Type
}

func (k *TUFKey) ID() string {
	if k.id == "" {
		logrus.Debug("Generating Key ID")
		pubK := NewTUFKey(k.Cipher(), k.Public(), "")
		data, err := cjson.Marshal(&pubK)
		if err != nil {
			logrus.Error("Error generating key ID:", err)
		}
		digest := sha256.Sum256(data)
		k.id = hex.EncodeToString(digest[:])
	}
	return k.id
}

func (k TUFKey) Public() string {
	return k.Value.Public
}

type PublicKey struct {
	TUFKey
}

func (k PublicKey) Private() string {
	return ""
}

func NewPublicKey(cipher, public string) *PublicKey {
	return &PublicKey{
		TUFKey{
			Type: cipher,
			Value: KeyPair{
				Public:  public,
				Private: "",
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

func NewPrivateKey(cipher, public, private string) *PrivateKey {
	return &PrivateKey{
		TUFKey{
			Type: cipher,
			Value: KeyPair{
				Public:  public,
				Private: private,
			},
		},
	}
}

func (k PrivateKey) Private() string {
	return k.Value.Private
}
