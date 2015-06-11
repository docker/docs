package signed

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/agl/ed25519"
	"github.com/endophage/gotuf/data"
)

// Ed25519 implements a simple in memory keystore and trust service
type Ed25519 struct {
	keys map[string]*data.PrivateKey
}

func NewEd25519() *Ed25519 {
	return &Ed25519{
		make(map[string]*data.PrivateKey),
	}
}

// addKey allows you to add a private key to the trust service
func (trust *Ed25519) addKey(k *data.PrivateKey) {
	trust.keys[k.ID()] = k
}

func (trust *Ed25519) RemoveKey(keyID string) {
	delete(trust.keys, keyID)
}

func (trust *Ed25519) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, kID := range keyIDs {
		priv := [ed25519.PrivateKeySize]byte{}
		pub := [ed25519.PublicKeySize]byte{}
		copy(priv[:], trust.keys[kID].Private())
		copy(pub[:], trust.keys[kID].Public())
		sig := ed25519.Sign(&priv, toSign)
		signatures = append(signatures, data.Signature{
			KeyID:     kID,
			Method:    "ed25519",
			Signature: sig[:],
		})
	}
	return signatures, nil

}

func (trust *Ed25519) Create() (*data.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pubStr := hex.EncodeToString(pub[:])
	privStr := hex.EncodeToString(priv[:])
	public := data.NewPublicKey("ed25519", pubStr)
	private := data.NewPrivateKey("ed25519", pubStr, privStr)
	trust.addKey(private)
	return public, nil
}

func (trust *Ed25519) PublicKeys(keyIDs ...string) (map[string]*data.PublicKey, error) {
	k := make(map[string]*data.PublicKey)
	for _, kID := range keyIDs {
		if key, ok := trust.keys[kID]; ok {
			k[kID] = data.PublicKeyFromPrivate(*key)
		}
	}
	return k, nil
}

func (trust *Ed25519) CanSign(keyID string) bool {
	_, ok := trust.keys[keyID]
	return ok
}
