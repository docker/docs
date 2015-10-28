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

// Keys represents a map of key ID to PublicKey object. It's necessary
// to allow us to unmarshal into an interface via the json.Unmarshaller
// interface
type Keys map[string]PublicKey

// UnmarshalJSON implements the json.Unmarshaller interface
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

// KeyList represents a list of keys
type KeyList []PublicKey

// UnmarshalJSON implements the json.Unmarshaller interface
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

// NewPublicKey creates a new, correctly typed PublicKey, using the
// UnknownPublicKey catchall for unsupported ciphers
func NewPublicKey(alg string, public []byte) PublicKey {
	tk := tufKey{
		Type: alg,
		Value: KeyPair{
			Public: public,
		},
	}
	return typedPublicKey(tk)
}

// NewPrivateKey creates a new, correctly typed PrivateKey, using the
// UnknownPrivateKey catchall for unsupported ciphers
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

// UnmarshalPublicKey is used to parse individual public keys in JSON
func UnmarshalPublicKey(data []byte) (PublicKey, error) {
	var parsed tufKey
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return typedPublicKey(parsed), nil
}

// UnmarshalPrivateKey is used to parse individual private keys in JSON
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

// ECDSAPublicKey represents an ECDSA key using a raw serialization
// of the public key
type ECDSAPublicKey struct {
	tufKey
}

// ECDSAx509PublicKey represents an ECDSA key using an x509 cert
// as the serialized format of the public key
type ECDSAx509PublicKey struct {
	tufKey
}

// RSAPublicKey represents an RSA key using a raw serialization
// of the public key
type RSAPublicKey struct {
	tufKey
}

// RSAx509PublicKey represents an RSA key using an x509 cert
// as the serialized format of the public key
type RSAx509PublicKey struct {
	tufKey
}

// ED25519PublicKey represents an ED25519 key using a raw serialization
// of the public key
type ED25519PublicKey struct {
	tufKey
}

// UnknownPublicKey is a catchall for key types that are not supported
type UnknownPublicKey struct {
	tufKey
}

// NewECDSAPublicKey initializes a new public key with the ECDSAKey type
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

// NewECDSAx509PublicKey initializes a new public key with the ECDSAx509Key type
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

// NewRSAPublicKey initializes a new public key with the RSA type
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

// NewRSAx509PublicKey initializes a new public key with the RSAx509Key type
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

// NewED25519PublicKey initializes a new public key with the ED25519Key type
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

// ECDSAPrivateKey represents a private ECDSA key
type ECDSAPrivateKey struct {
	ECDSAPublicKey
	private []byte
}

// ECDSAx509PrivateKey represents a private ECDSA key where the public
// component is serialized into an x509 cert
type ECDSAx509PrivateKey struct {
	ECDSAx509PublicKey
	private []byte
}

// RSAPrivateKey represents a private RSA key
type RSAPrivateKey struct {
	RSAPublicKey
	private []byte
}

// RSAx509PrivateKey represents a private RSA key where the public
// component is serialized into an x509 cert
type RSAx509PrivateKey struct {
	RSAx509PublicKey
	private []byte
}

// ED25519PrivateKey represents a private ED25519 key
type ED25519PrivateKey struct {
	ED25519PublicKey
	private []byte
}

// UnknownPrivateKey is a catchall for unsupported key types
type UnknownPrivateKey struct {
	tufKey
	private []byte
}

// NewECDSAPrivateKey initializes a new ECDSA private key
func NewECDSAPrivateKey(public ECDSAPublicKey, private []byte) *ECDSAPrivateKey {
	return &ECDSAPrivateKey{
		ECDSAPublicKey: public,
		private:        private,
	}
}

// NewECDSAx509PrivateKey initializes a new ECDSA private key
func NewECDSAx509PrivateKey(public ECDSAx509PublicKey, private []byte) *ECDSAx509PrivateKey {
	return &ECDSAx509PrivateKey{
		ECDSAx509PublicKey: public,
		private:            private,
	}
}

// NewRSAPrivateKey initialized a new RSA private key
func NewRSAPrivateKey(public RSAPublicKey, private []byte) *RSAPrivateKey {
	return &RSAPrivateKey{
		RSAPublicKey: public,
		private:      private,
	}
}

// NewRSAx509PrivateKey initialized a new RSA private key
func NewRSAx509PrivateKey(public RSAx509PublicKey, private []byte) *RSAx509PrivateKey {
	return &RSAx509PrivateKey{
		RSAx509PublicKey: public,
		private:          private,
	}
}

// NewED25519PrivateKey initialized a new ED25519 private key
func NewED25519PrivateKey(public ED25519PublicKey, private []byte) *ED25519PrivateKey {
	return &ED25519PrivateKey{
		ED25519PublicKey: public,
		private:          private,
	}
}

// Private return the serialized private bytes of the key
func (k ECDSAPrivateKey) Private() []byte {
	return k.private
}

// Private return the serialized private bytes of the key
func (k ECDSAx509PrivateKey) Private() []byte {
	return k.private
}

// Private return the serialized private bytes of the key
func (k RSAPrivateKey) Private() []byte {
	return k.private
}

// Private return the serialized private bytes of the key
func (k RSAx509PrivateKey) Private() []byte {
	return k.private
}

// Private return the serialized private bytes of the key
func (k ED25519PrivateKey) Private() []byte {
	return k.private
}

// Private return the serialized private bytes of the key
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
