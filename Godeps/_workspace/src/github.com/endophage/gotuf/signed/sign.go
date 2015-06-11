package signed

import (
	"github.com/endophage/gotuf/data"
)

// Signer encapsulates a signing service with some convenience methods to
// interface between TUF keys and the generic service interface
type Signer struct {
	service CryptoService
}

func NewSigner(service CryptoService) *Signer {
	return &Signer{service}
}

// Sign takes a data.Signed and a key, calculated and adds the signature
// to the data.Signed
func (signer *Signer) Sign(s *data.Signed, keys ...*data.PublicKey) error {
	signatures := make([]data.Signature, 0, len(s.Signatures)+1)
	keyIDMemb := make(map[string]struct{})
	keyIDs := make([]string, 0, len(keys))
	for _, key := range keys {
		keyIDMemb[key.ID()] = struct{}{}
		keyIDs = append(keyIDs, key.ID())
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

func (signer *Signer) Create() (*data.PublicKey, error) {
	return signer.service.Create()
}

func (signer *Signer) PublicKeys(keyIDs ...string) (map[string]*data.PublicKey, error) {
	return signer.service.PublicKeys(keyIDs...)
}
