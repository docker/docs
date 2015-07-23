package api

import (
	"fmt"
	"log"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"github.com/endophage/gotuf/data"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/docker/notary/proto"
)

//KeyManagementServer implements the KeyManagementServer grpc interface
type KeyManagementServer struct {
	CryptoServices signer.CryptoServiceIndex
}

//SignerServer implements the SignerServer grpc interface
type SignerServer struct {
	CryptoServices signer.CryptoServiceIndex
}

//CreateKey returns a PublicKey created using KeyManagementServer's SigningService
func (s *KeyManagementServer) CreateKey(ctx context.Context, algorithm *pb.Algorithm) (*pb.PublicKey, error) {
	keyAlgo := data.KeyAlgorithm(algorithm.Algorithm)

	service := s.CryptoServices[keyAlgo]

	if service == nil {
		return nil, fmt.Errorf("algorithm %s not supported for create key", algorithm.Algorithm)
	}

	tufKey, err := service.Create("", keyAlgo)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Key creation failed")
	}
	log.Println("[Notary-signer CreateKey] : Created KeyID ", tufKey.ID())
	return &pb.PublicKey{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
		},
		PublicKey: tufKey.Public(),
	}, nil
}

//DeleteKey deletes they key associated with a KeyID
func (s *KeyManagementServer) DeleteKey(ctx context.Context, keyID *pb.KeyID) (*pb.Void, error) {
	_, service, err := FindKeyByID(s.CryptoServices, keyID)

	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyID.ID)
	}

	err = service.RemoveKey(keyID.ID)
	log.Println("[Notary-signer DeleteKey] : Deleted KeyID ", keyID.ID)
	if err != nil {
		switch err {
		case keys.ErrInvalidKeyID:
			return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyID.ID)
		default:
			return nil, grpc.Errorf(codes.Internal, "Key deletion for keyID %s failed", keyID.ID)
		}
	}

	return &pb.Void{}, nil
}

//GetKeyInfo returns they PublicKey associated with a KeyID
func (s *KeyManagementServer) GetKeyInfo(ctx context.Context, keyID *pb.KeyID) (*pb.PublicKey, error) {
	_, service, err := FindKeyByID(s.CryptoServices, keyID)

	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyID.ID)
	}

	tufKey := service.GetKey(keyID.ID)
	if tufKey == nil {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyID.ID)
	}
	log.Println("[Notary-signer GetKeyInfo] : Returning PublicKey for KeyID ", keyID.ID)
	return &pb.PublicKey{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
		},
		PublicKey: tufKey.Public(),
	}, nil
}

//Sign signs a message and returns the signature using a private key associate with the KeyID from the SignatureRequest
func (s *SignerServer) Sign(ctx context.Context, sr *pb.SignatureRequest) (*pb.Signature, error) {
	tufKey, service, err := FindKeyByID(s.CryptoServices, sr.KeyID)

	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", sr.KeyID.ID)
	}

	log.Println("[Notary-signer Sign] : Signing ", string(sr.Content), " with KeyID ", sr.KeyID.ID)

	signatures, err := service.Sign([]string{sr.KeyID.ID}, sr.Content)
	if err != nil || len(signatures) != 1 {
		return nil, grpc.Errorf(codes.Internal, "Signing failed for keyID %s on hash %s", sr.KeyID.ID, sr.Content)
	}

	signature := &pb.Signature{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
		},
		Algorithm: &pb.Algorithm{Algorithm: signatures[0].Method.String()},
		Content:   signatures[0].Signature,
	}

	return signature, nil
}
