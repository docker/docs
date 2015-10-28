package api

import (
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"

	pb "github.com/docker/notary/proto"
)

// FindKeyByID looks for the key with the given ID in each of the
// signing services in sigServices. It returns the first matching key it finds,
// or ErrInvalidKeyID if the key is not found in any of the signing services.
// It also returns the CryptoService associated with the key, so the caller
// can perform operations with the key (such as signing).
func FindKeyByID(cryptoServices signer.CryptoServiceIndex, keyID *pb.KeyID) (data.PublicKey, signed.CryptoService, error) {
	for _, service := range cryptoServices {
		key := service.GetKey(keyID.ID)
		if key != nil {
			return key, service, nil
		}
	}

	return nil, nil, keys.ErrInvalidKeyID
}
