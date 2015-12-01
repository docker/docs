package api

import (
	"crypto/rand"
	"fmt"

	ctxu "github.com/docker/distribution/context"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/docker/notary/proto"
)

//KeyManagementServer implements the KeyManagementServer grpc interface
type KeyManagementServer struct {
	CryptoServices signer.CryptoServiceIndex
	HealthChecker  func() map[string]string
}

//SignerServer implements the SignerServer grpc interface
type SignerServer struct {
	CryptoServices signer.CryptoServiceIndex
	HealthChecker  func() map[string]string
}

//CreateKey returns a PublicKey created using KeyManagementServer's SigningService
func (s *KeyManagementServer) CreateKey(ctx context.Context, algorithm *pb.Algorithm) (*pb.PublicKey, error) {
	keyAlgo := algorithm.Algorithm

	service := s.CryptoServices[keyAlgo]

	logger := ctxu.GetLogger(ctx)

	if service == nil {
		logger.Error("CreateKey: unsupported algorithm: ", algorithm.Algorithm)
		return nil, fmt.Errorf("algorithm %s not supported for create key", algorithm.Algorithm)
	}

	tufKey, err := service.Create("", keyAlgo)
	if err != nil {
		logger.Error("CreateKey: failed to create key: ", err)
		return nil, grpc.Errorf(codes.Internal, "Key creation failed")
	}
	logger.Info("CreateKey: Created KeyID ", tufKey.ID())
	return &pb.PublicKey{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
		},
		PublicKey: tufKey.Public(),
	}, nil
}

//DeleteKey deletes they key associated with a KeyID
func (s *KeyManagementServer) DeleteKey(ctx context.Context, keyID *pb.KeyID) (*pb.Void, error) {
	_, service, err := FindKeyByID(s.CryptoServices, keyID)

	logger := ctxu.GetLogger(ctx)

	if err != nil {
		logger.Errorf("DeleteKey: key %s not found", keyID.ID)
		return nil, grpc.Errorf(codes.NotFound, "key %s not found", keyID.ID)
	}

	err = service.RemoveKey(keyID.ID)
	logger.Info("DeleteKey: Deleted KeyID ", keyID.ID)
	if err != nil {
		switch err {
		case keys.ErrInvalidKeyID:
			logger.Errorf("DeleteKey: key %s not found", keyID.ID)
			return nil, grpc.Errorf(codes.NotFound, "key %s not found", keyID.ID)
		default:
			logger.Error("DeleteKey: deleted key ", keyID.ID)
			return nil, grpc.Errorf(codes.Internal, "Key deletion for KeyID %s failed", keyID.ID)
		}
	}

	return &pb.Void{}, nil
}

//GetKeyInfo returns they PublicKey associated with a KeyID
func (s *KeyManagementServer) GetKeyInfo(ctx context.Context, keyID *pb.KeyID) (*pb.PublicKey, error) {
	_, service, err := FindKeyByID(s.CryptoServices, keyID)

	logger := ctxu.GetLogger(ctx)

	if err != nil {
		logger.Errorf("GetKeyInfo: key %s not found", keyID.ID)
		return nil, grpc.Errorf(codes.NotFound, "key %s not found", keyID.ID)
	}

	tufKey := service.GetKey(keyID.ID)
	if tufKey == nil {
		logger.Errorf("GetKeyInfo: key %s not found", keyID.ID)
		return nil, grpc.Errorf(codes.NotFound, "key %s not found", keyID.ID)
	}
	logger.Debug("GetKeyInfo: Returning PublicKey for KeyID ", keyID.ID)
	return &pb.PublicKey{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
		},
		PublicKey: tufKey.Public(),
	}, nil
}

//CheckHealth returns the HealthStatus with the service
func (s *KeyManagementServer) CheckHealth(ctx context.Context, v *pb.Void) (*pb.HealthStatus, error) {
	return &pb.HealthStatus{
		Status: s.HealthChecker(),
	}, nil
}

//Sign signs a message and returns the signature using a private key associate with the KeyID from the SignatureRequest
func (s *SignerServer) Sign(ctx context.Context, sr *pb.SignatureRequest) (*pb.Signature, error) {
	tufKey, service, err := FindKeyByID(s.CryptoServices, sr.KeyID)

	logger := ctxu.GetLogger(ctx)

	if err != nil {
		logger.Errorf("Sign: key %s not found", sr.KeyID.ID)
		return nil, grpc.Errorf(codes.NotFound, "key %s not found", sr.KeyID.ID)
	}

	privKey, _, err := service.GetPrivateKey(tufKey.ID())
	if err != nil {
		logger.Errorf("Sign: key %s not found", sr.KeyID.ID)
		return nil, grpc.Errorf(codes.NotFound, "key %s not found", sr.KeyID.ID)
	}
	sig, err := privKey.Sign(rand.Reader, sr.Content, nil)
	if err != nil {
		logger.Errorf("Sign: signing failed for KeyID %s on hash %s", sr.KeyID.ID, sr.Content)
		return nil, grpc.Errorf(codes.Internal, "Signing failed for KeyID %s on hash %s", sr.KeyID.ID, sr.Content)
	}

	logger.Info("Sign: Signed ", string(sr.Content), " with KeyID ", sr.KeyID.ID)

	signature := &pb.Signature{
		KeyInfo: &pb.KeyInfo{
			KeyID:     &pb.KeyID{ID: tufKey.ID()},
			Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
		},
		Algorithm: &pb.Algorithm{Algorithm: privKey.SignatureAlgorithm().String()},
		Content:   sig,
	}

	return signature, nil
}

//CheckHealth returns the HealthStatus with the service
func (s *SignerServer) CheckHealth(ctx context.Context, v *pb.Void) (*pb.HealthStatus, error) {
	return &pb.HealthStatus{
		Status: s.HealthChecker(),
	}, nil
}
