package signed

// The Sign function is a choke point for all code paths that do signing.
// We use this fact to do key ID translation. There are 2 types of key ID:
//   - Scoped: the key ID based purely on the data that appears in the TUF
//             files. This may be wrapped by a certificate that scopes the
//             key to be used in a specific context.
//   - Canonical: the key ID based purely on the public key bytes. This is
//             used by keystores to easily identify keys that may be reused
//             in many scoped locations.
// Currently these types only differ in the context of Root Keys in Notary
// for which the root key is wrapped using an x509 certificate.

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/errors"
	"github.com/endophage/gotuf/utils"
)

type idPair struct {
	scopedKeyID    string
	canonicalKeyID string
}

// Sign takes a data.Signed and a key, calculated and adds the signature
// to the data.Signed
func Sign(service CryptoService, s *data.Signed, keys ...data.PublicKey) error {
	logrus.Debugf("sign called with %d keys", len(keys))
	signatures := make([]data.Signature, 0, len(s.Signatures)+1)
	keyIDMemb := make(map[string]struct{})
	keyIDs := make(
		[]idPair,
		0,
		len(keys),
	)

	for _, key := range keys {
		keyID, err := utils.CanonicalKeyID(key)
		if err != nil {
			continue
		}
		keyIDMemb[key.ID()] = struct{}{}
		keyIDs = append(keyIDs, idPair{
			scopedKeyID:    key.ID(),
			canonicalKeyID: keyID,
		})
	}

	// we need to ask the signer to sign with the canonical key ID
	// (ID of the TUF key's public key bytes only), but
	// we need to translate back to the scoped key ID (the hash of the TUF key
	// with the full PEM bytes) before giving the
	// signature back to TUF.
	for _, pair := range keyIDs {
		newSigs, err := service.Sign([]string{pair.canonicalKeyID}, s.Signed)
		if err != nil {
			return err
		}
		// we only asked to sign with 1 key ID, so there will either be 1
		// or zero signatures
		if len(newSigs) == 1 {
			newSig := newSigs[0]
			newSig.KeyID = pair.scopedKeyID
			signatures = append(signatures, newSig)
		}
	}
	if len(signatures) < 1 {
		return errors.ErrInsufficientSignatures{
			Name: fmt.Sprintf("Cryptoservice failed to produce any signatures for keys with IDs: %v", keyIDs),
			Err:  nil,
		}
	}
	for _, sig := range s.Signatures {
		if _, ok := keyIDMemb[sig.KeyID]; ok {
			continue
		}
		signatures = append(signatures, sig)
	}
	s.Signatures = signatures
	return nil
}
