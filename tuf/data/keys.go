package data

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/Sirupsen/logrus"
	"github.com/jfrazelle/go/canonical/json"
)

// PublicKey is the necessary interface for public keys
type PublicKey interface {
	ID() string
	Algorithm() string
	Public() []byte
}

// PrivateKey adds the ability to access the private key
type PrivateKey interface {
	PublicKey

	Private() []byte
}

// KeyPair holds the public and private key bytes
type KeyPair struct {
	Public  []byte `json:"public"`
	Private []byte `json:"private"`
}

type Keys map[string]PublicKey

func (ks *Keys) UnmarshalJSON(data []byte) error {
	parsed := make(map[string]tufKey)
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	final := make(map[string]PublicKey)
	for k, tk := range parsed {
		final[k] = typedPublicKey(tk)
	}
	*ks = final
	return nil
}

type KeyList []PublicKey

func (ks *KeyList) UnmarshalJSON(data []byte) error {
	parsed := make([]tufKey, 0, 1)
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	final := make([]PublicKey, 0, len(parsed))
	for _, tk := range parsed {
		final = append(final, typedPublicKey(tk))
	}
	*ks = final
	return nil
}

func typedPublicKey(tk tufKey) PublicKey {
	switch tk.Algorithm() {
	case ECDSAKey:
		return &ECDSAPublicKey{tufKey: tk}
	case ECDSAx509Key:
		return &ECDSAx509PublicKey{tufKey: tk}
	case RSAKey:
		return &RSAPublicKey{tufKey: tk}
	case RSAx509Key:
		return &RSAx509PublicKey{tufKey: tk}
	case ED25519Key:
		return &ED25519PublicKey{tufKey: tk}
	}
	return &UnknownPublicKey{tufKey: tk}
}

func typedPrivateKey(tk tufKey) PrivateKey {
	private := tk.Value.Private
	tk.Value.Private = nil
	switch tk.Algorithm() {
	case ECDSAKey:
		return &ECDSAPrivateKey{
			ECDSAPublicKey: ECDSAPublicKey{
				tufKey: tk,
			},
			private: private,
		}
	case ECDSAx509Key:
		return &ECDSAx509PrivateKey{
			ECDSAx509PublicKey: ECDSAx509PublicKey{
				tufKey: tk,
			},
			private: private,
		}
	case RSAKey:
		return &RSAPrivateKey{
			RSAPublicKey: RSAPublicKey{
				tufKey: tk,
			},
			private: private,
		}
	case RSAx509Key:
		return &RSAx509PrivateKey{
			RSAx509PublicKey: RSAx509PublicKey{
				tufKey: tk,
			},
			private: private,
		}
	case ED25519Key:
		return &ED25519PrivateKey{
			ED25519PublicKey: ED25519PublicKey{
				tufKey: tk,
			},
			private: private,
		}
	}
	return &UnknownPrivateKey{
		tufKey:  tk,
		private: private,
	}
}

func NewPublicKey(alg string, public []byte) PublicKey {
	tk := tufKey{
		Type: alg,
		Value: KeyPair{
			Public: public,
		},
	}
	return typedPublicKey(tk)
}

func NewPrivateKey(pubKey PublicKey, private []byte) PrivateKey {
	tk := tufKey{
		Type: pubKey.Algorithm(),
		Value: KeyPair{
			Public:  pubKey.Public(),
			Private: private, // typedPrivateKey moves this value
		},
	}
	return typedPrivateKey(tk)
}

func UnmarshalPublicKey(data []byte) (PublicKey, error) {
	var parsed tufKey
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return typedPublicKey(parsed), nil
}

func UnmarshalPrivateKey(data []byte) (PrivateKey, error) {
	var parsed tufKey
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return typedPrivateKey(parsed), nil
}

// tufKey is the structure used for both public and private keys in TUF.
// Normally it would make sense to use a different structures for public and
// private keys, but that would change the key ID algorithm (since the canonical
// JSON would be different). This structure should normally be accessed through
// the PublicKey or PrivateKey interfaces.
type tufKey struct {
	id    string
	Type  string  `json:"keytype"`
	Value KeyPair `json:"keyval"`
}

// Algorithm returns the algorithm of the key
func (k tufKey) Algorithm() string {
	return k.Type
}

// ID efficiently generates if necessary, and caches the ID of the key
func (k *tufKey) ID() string {
	if k.id == "" {
		pubK := tufKey{
			Type: k.Algorithm(),
			Value: KeyPair{
				Public:  k.Public(),
				Private: nil,
			},
		}
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
func (k tufKey) Public() []byte {
	return k.Value.Public
}

// Public key types

type ECDSAPublicKey struct {
	tufKey
}

type ECDSAx509PublicKey struct {
	tufKey
}

type RSAPublicKey struct {
	tufKey
}

type RSAx509PublicKey struct {
	tufKey
}

type ED25519PublicKey struct {
	tufKey
}

type UnknownPublicKey struct {
	tufKey
}

func NewECDSAPublicKey(public []byte) *ECDSAPublicKey {
	return &ECDSAPublicKey{
		tufKey: tufKey{
			Type: ECDSAKey,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

func NewECDSAx509PublicKey(public []byte) *ECDSAx509PublicKey {
	return &ECDSAx509PublicKey{
		tufKey: tufKey{
			Type: ECDSAx509Key,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

func NewRSAPublicKey(public []byte) *RSAPublicKey {
	return &RSAPublicKey{
		tufKey: tufKey{
			Type: RSAKey,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

func NewRSAx509PublicKey(public []byte) *RSAx509PublicKey {
	return &RSAx509PublicKey{
		tufKey: tufKey{
			Type: RSAx509Key,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

func NewED25519PublicKey(public []byte) *ED25519PublicKey {
	return &ED25519PublicKey{
		tufKey: tufKey{
			Type: ED25519Key,
			Value: KeyPair{
				Public:  public,
				Private: nil,
			},
		},
	}
}

// Private key types

// ECDSAPrivateKey
type ECDSAPrivateKey struct {
	ECDSAPublicKey
	private []byte `json:"-"`
}

type ECDSAx509PrivateKey struct {
	ECDSAx509PublicKey
	private []byte `json:"-"`
}

type RSAPrivateKey struct {
	RSAPublicKey
	private []byte `json:"-"`
}

type RSAx509PrivateKey struct {
	RSAx509PublicKey
	private []byte `json:"-"`
}

type ED25519PrivateKey struct {
	ED25519PublicKey
	private []byte `json:"-"`
}

type UnknownPrivateKey struct {
	tufKey
	private []byte `json:"-"`
}

func NewECDSAPrivateKey(public ECDSAPublicKey, private []byte) *ECDSAPrivateKey {
	return &ECDSAPrivateKey{
		ECDSAPublicKey: public,
		private:        private,
	}
}

func NewECDSAx509PrivateKey(public ECDSAx509PublicKey, private []byte) *ECDSAx509PrivateKey {
	return &ECDSAx509PrivateKey{
		ECDSAx509PublicKey: public,
		private:            private,
	}
}

func NewRSAPrivateKey(public RSAPublicKey, private []byte) *RSAPrivateKey {
	return &RSAPrivateKey{
		RSAPublicKey: public,
		private:      private,
	}
}

func NewRSAx509PrivateKey(public RSAx509PublicKey, private []byte) *RSAx509PrivateKey {
	return &RSAx509PrivateKey{
		RSAx509PublicKey: public,
		private:          private,
	}
}

func NewED25519PrivateKey(public ED25519PublicKey, private []byte) *ED25519PrivateKey {
	return &ED25519PrivateKey{
		ED25519PublicKey: public,
		private:          private,
	}
}

func (k ECDSAPrivateKey) Private() []byte {
	return k.private
}

func (k ECDSAx509PrivateKey) Private() []byte {
	return k.private
}

func (k RSAPrivateKey) Private() []byte {
	return k.private
}

func (k RSAx509PrivateKey) Private() []byte {
	return k.private
}

func (k ED25519PrivateKey) Private() []byte {
	return k.private
}

func (k UnknownPrivateKey) Private() []byte {
	return k.private
}

// PublicKeyFromPrivate returns a new tufKey based on a private key, with
// the private key bytes guaranteed to be nil.
func PublicKeyFromPrivate(pk PrivateKey) PublicKey {
	return typedPublicKey(tufKey{
		Type: pk.Algorithm(),
		Value: KeyPair{
			Public:  pk.Public(),
			Private: nil,
		},
	})
}

type ed25519PrivateKey []byte
