package api

import (
	"fmt"
	"log"

	"github.com/docker/rufus"
	"github.com/docker/rufus/keys"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/docker/rufus/proto"
)

//KeyManagementServer implements the KeyManagementServer grpc interface
type KeyManagementServer struct {
	SigServices rufus.SigningServiceIndex
}

//SignerServer implements the SignerServer grpc interface
type SignerServer struct {
	SigServices rufus.SigningServiceIndex
}

//CreateKey returns a PublicKey created using KeyManagementServer's SigningService
func (s *KeyManagementServer) CreateKey(ctx context.Context, algorithm *pb.Algorithm) (*pb.PublicKey, error) {
	service := s.SigServices[algorithm.Algorithm]

	if service == nil {
		return nil, fmt.Errorf("algorithm %s not supported for create key", algorithm.Algorithm)
	}

	key, err := service.CreateKey()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Key creation failed")
	}
	log.Println("[Rufus CreateKey] : Created KeyID ", key.KeyInfo.ID)
	return key, nil
}

//DeleteKey deletes they key associated with a KeyID
func (s *KeyManagementServer) DeleteKey(ctx context.Context, keyInfo *pb.KeyInfo) (*pb.Void, error) {
	service := s.SigServices[keyInfo.Algorithm.Algorithm]

	if service == nil {
		return nil, fmt.Errorf("algorithm %s not supported for delete key", keyInfo.Algorithm.Algorithm)
	}

	_, err := service.DeleteKey(keyInfo)
	log.Println("[Rufus DeleteKey] : Deleted KeyID ", keyInfo.ID)
	if err != nil {
		switch err {
		case keys.ErrInvalidKeyID:
			return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyInfo.ID)
		default:
			return nil, grpc.Errorf(codes.Internal, "Key deletion for keyID %s failed", keyInfo.ID)
		}
	}

	return &pb.Void{}, nil
}

//GetKeyInfo returns they PublicKey associated with a KeyID
func (s *KeyManagementServer) GetKeyInfo(ctx context.Context, keyInfo *pb.KeyInfo) (*pb.PublicKey, error) {
	service := s.SigServices[keyInfo.Algorithm.Algorithm]

	if service == nil {
		return nil, fmt.Errorf("algorithm %s not supported for get key info", keyInfo.Algorithm.Algorithm)
	}

	key, err := service.KeyInfo(keyInfo)
	if err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key %s not found", keyInfo.ID)
	}
	log.Println("[Rufus GetKeyInfo] : Returning PublicKey for KeyID ", keyInfo.ID)
	return key, nil
}

//Sign signs a message and returns the signature using a private key associate with the KeyID from the SignatureRequest
func (s *SignerServer) Sign(ctx context.Context, sr *pb.SignatureRequest) (*pb.Signature, error) {
	service := s.SigServices[sr.KeyInfo.Algorithm.Algorithm]

	if service == nil {
		return nil, fmt.Errorf("algorithm %s not supported for sign", sr.KeyInfo.Algorithm.Algorithm)
	}

	log.Println("[Rufus Sign] : Signing ", string(sr.Content), " with KeyID ", sr.KeyInfo.ID)
	signer, err := service.Signer(sr.KeyInfo)
	if err == keys.ErrInvalidKeyID {
		return nil, grpc.Errorf(codes.NotFound, "Invalid keyID: key not found")
	}

	signature, err := signer.Sign(sr)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Signing failed for keyID %s on hash %s", sr.KeyInfo.ID, sr.Content)
	}

	return signature, nil
}
