package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/agl/ed25519"
	"github.com/docker/rufus"
	"github.com/docker/rufus/keys"

	pb "github.com/docker/rufus/proto"
)

// EdDSASigningService is an implementation of SigningService
type EdDSASigningService struct {
	KeyDB rufus.KeyDatabase
}

// CreateKey creates a key and returns its public components
func (s EdDSASigningService) CreateKey() (*pb.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	k := &keys.Key{
		Algorithm: ED25519,
		Public:    *pub,
		Private:   priv,
	}
	digest := sha256.Sum256(k.Public[:])
	k.ID = hex.EncodeToString(digest[:])

	err = s.KeyDB.AddKey(k)
	if err != nil {
		return nil, err
	}

	pubKey := &pb.PublicKey{KeyInfo: &pb.KeyInfo{ID: k.ID, Algorithm: &pb.Algorithm{Algorithm: k.Algorithm}}, PublicKey: k.Public[:]}

	return pubKey, nil
}

// DeleteKey removes a key from the key database
func (s EdDSASigningService) DeleteKey(keyInfo *pb.KeyInfo) (*pb.Void, error) {
	return s.KeyDB.DeleteKey(keyInfo)
}

// KeyInfo returns the public components of a particular key
func (s EdDSASigningService) KeyInfo(keyInfo *pb.KeyInfo) (*pb.PublicKey, error) {
	return s.KeyDB.KeyInfo(keyInfo)
}

// Signer returns a Signer for a specific KeyID
func (s EdDSASigningService) Signer(keyInfo *pb.KeyInfo) (rufus.Signer, error) {
	key, err := s.KeyDB.GetKey(keyInfo)
	if err != nil {
		return nil, keys.ErrInvalidKeyID
	}
	return &Ed25519Signer{privateKey: key}, nil
}

// NewEdDSASigningService returns an instance of KeyDB
func NewEdDSASigningService(keyDB rufus.KeyDatabase) *EdDSASigningService {
	return &EdDSASigningService{
		KeyDB: keyDB,
	}
}
