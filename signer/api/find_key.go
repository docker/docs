package api

import (
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"

	pb "github.com/docker/notary/proto"
)

func FindKeyByID(sigServices signer.SigningServiceIndex, keyID *pb.KeyID) (*pb.PublicKey, signer.SigningService, error) {
	for _, service := range sigServices {
		key, err := service.KeyInfo(keyID)
		if err == nil {
			return key, service, nil
		}
	}

	return nil, nil, keys.ErrInvalidKeyID
}
