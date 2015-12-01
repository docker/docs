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
	"crypto/rand"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/utils"
)

// Sign takes a data.Signed and replaces data.Signature with at least
// minPrimarySignature signatures using primaryKeys and any possible
// signatures using optionalKeys
func Sign(service CryptoService, s *data.Signed, primaryKeys []data.PublicKey,
	minSignatures int, optionalKeys []data.PublicKey) error {

	logrus.Debugf("sign called with %d/%d required + %d optional keys",
		minSignatures, len(primaryKeys), len(optionalKeys))

	privKeys := make(map[string]data.PrivateKey)

	// Get all the private key objects related to the public keys
	missingKeyIDs := []string{}
	for _, key := range primaryKeys {
		canonicalID, err := utils.CanonicalKeyID(key)
		if err != nil {
			return err
		}
		k, _, err := service.GetPrivateKey(canonicalID)
		if err != nil {
			if _, ok := err.(trustmanager.ErrKeyNotFound); ok {
				missingKeyIDs = append(missingKeyIDs, canonicalID)
				continue
			}
			return err
		}
		privKeys[key.ID()] = k
	}

	// Check to ensure we have enough signing keys
	if len(privKeys) < minSignatures {
		return ErrInsufficientSignatures{FoundKeys: len(privKeys),
			NeededKeys: minSignatures, MissingKeyIDs: missingKeyIDs}
	}

	for _, key := range optionalKeys {
		if _, ok := privKeys[key.ID()]; ok {
			continue
		}
		canonicalID, err := utils.CanonicalKeyID(key)
		if err != nil {
			return err
		}
		k, _, err := service.GetPrivateKey(canonicalID)
		if err != nil {
			continue
		}
		privKeys[key.ID()] = k
	}

	// Do signing and generate list of signatures
	signatures := make([]data.Signature, 0, len(privKeys))
	for keyID, pk := range privKeys {
		sig, err := pk.Sign(rand.Reader, *s.Signed, nil)
		if err != nil {
			logrus.Debugf("Failed to sign with key: %s. Reason: %v", keyID, err)
			return err
		}
		signatures = append(signatures, data.Signature{
			KeyID:     keyID,
			Method:    pk.SignatureAlgorithm(),
			Signature: sig[:],
		})
	}

	s.Signatures = signatures
	return nil
}
