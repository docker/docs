package signed

import (
	cjson "github.com/endophage/go-tuf/Godeps/_workspace/src/github.com/tent/canonical-json-go"

	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/keys"
)

// Sign takes a data.Signed and a key, calculated and adds the signature
// to the data.Signed
//func Sign(s *data.Signed, k *data.Key) {
//	id := k.ID()
//	signatures := make([]data.Signature, 0, len(s.Signatures)+1)
//	for _, sig := range s.Signatures {
//		if sig.KeyID == id {
//			continue
//		}
//		signatures = append(signatures, sig)
//	}
//	priv := [ed25519.PrivateKeySize]byte{}
//	copy(priv[:], k.Value.Private)
//	sig := ed25519.Sign(&priv, s.Signed)
//	s.Signatures = append(signatures, data.Signature{
//		KeyID:     id,
//		Method:    "ed25519",
//		Signature: sig[:],
//	})
//}

// Signer encapsulates a signing service with some convenience methods to
// interface between TUF keys and the generic service interface
type Signer struct {
	service TrustService
}

func NewSigner(service TrustService) *Signer {
	return &Signer{service}
}

// Sign takes a data.Signed and a key, calculated and adds the signature
// to the data.Signed
func (signer *Signer) Sign(s *data.Signed, keys ...*keys.PublicKey) error {
	signatures := make([]data.Signature, 0, len(s.Signatures)+1)
	keyIDMemb := make(map[string]struct{})
	keyIDs := make([]string, 0, len(keys))
	for _, key := range keys {
		keyIDMemb[key.ID] = struct{}{}
		keyIDs = append(keyIDs, key.ID)
	}
	for _, sig := range s.Signatures {
		if _, ok := keyIDMemb[sig.KeyID]; ok {
			continue
		}
		signatures = append(signatures, sig)
	}
	newSigs, err := signer.service.Sign(keyIDs, s.Signed)

	if err != nil {
		return err
	}
	s.Signatures = append(signatures, newSigs...)
	return nil
}

func (signer *Signer) Marshal(v interface{}, keys ...*keys.PublicKey) (*data.Signed, error) {
	b, err := cjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	s := &data.Signed{Signed: b}
	err = signer.Sign(s, keys...)
	return s, err // err may be nil but there's no point in checking, just return it
}

func (signer *Signer) NewKey() (*keys.PublicKey, error) {
	return signer.service.Create()
}
