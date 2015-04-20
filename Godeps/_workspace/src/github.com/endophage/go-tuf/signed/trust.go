package signed

import (
	"crypto/rand"

	"github.com/agl/ed25519"
	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/keys"
)

// Ed25519 implements a simple in memory keystore and trust service
type Ed25519 struct {
	keys map[string]*keys.PrivateKey
}

var _ TrustService = &Ed25519{}

func NewEd25519() *Ed25519 {
	return &Ed25519{
		make(map[string]*keys.PrivateKey),
	}
}

// addKey allows you to add a private key to the trust service
func (trust *Ed25519) addKey(k *keys.PrivateKey) {
	key := keys.PrivateKey{
		PublicKey: keys.PublicKey{
			Key: data.Key{
				Type: k.Type,
				Value: data.KeyValue{
					Public: make([]byte, len(k.Value.Public)),
				},
			},
			ID: k.ID,
		},
		Private: make([]byte, len(k.Private)),
	}

	copy(key.Value.Public, k.Value.Public)
	copy(key.Private, k.Private)
	trust.keys[k.ID] = &key
}

func (trust *Ed25519) RemoveKey(keyID string) {
	delete(trust.keys, keyID)
}

func (trust *Ed25519) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, kID := range keyIDs {
		priv := [ed25519.PrivateKeySize]byte{}
		pub := [ed25519.PublicKeySize]byte{}
		copy(priv[:], trust.keys[kID].Private)
		copy(pub[:], trust.keys[kID].Value.Public)
		sig := ed25519.Sign(&priv, toSign)
		signatures = append(signatures, data.Signature{
			KeyID:     kID,
			Method:    "ed25519",
			Signature: sig[:],
		})
	}
	return signatures, nil

}

func (trust *Ed25519) Create() (*keys.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pubBytes := make([]byte, ed25519.PublicKeySize)
	copy(pubBytes, pub[:])
	privBytes := make([]byte, ed25519.PrivateKeySize)
	copy(privBytes, priv[:])
	public := keys.NewPublicKey("ed25519", pubBytes)
	private := keys.PrivateKey{*public, privBytes}
	trust.addKey(&private)
	return public, nil
}

func (trust *Ed25519) PublicKeys(keyIDs ...string) (map[string]*keys.PublicKey, error) {
	k := make(map[string]*keys.PublicKey)
	for _, kID := range keyIDs {
		if key, ok := trust.keys[kID]; ok {
			k[kID] = &key.PublicKey
		}
	}
	return k, nil
}
