package api

import (
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"

	pb "github.com/docker/notary/proto"
)

// FindKeyByID looks for the key with the given ID in each of the
// signing services in sigServices. It returns the first matching key it finds,
// or ErrInvalidKeyID if the key is not found in any of the signing services.
func FindKeyByID(sigServices signer.SigningServiceIndex, keyID *pb.KeyID) (*pb.PublicKey, signer.SigningService, error) {
	for _, service := range sigServices {
		key, err := service.KeyInfo(keyID)
		if err == nil {
			return key, service, nil
		}
	}

	return nil, nil, keys.ErrInvalidKeyID
}
