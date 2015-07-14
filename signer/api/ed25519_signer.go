package api

import (
	"github.com/agl/ed25519"
	"github.com/docker/rufus/keys"

	pb "github.com/docker/rufus/proto"
)

// ED25519 represents an ed25519 algorithm
const ED25519 string = "ed25519"

// Ed25519Signer implements the Signer interface for Ed25519 keys
type Ed25519Signer struct {
	privateKey *keys.Key
}

// Sign returns a signature for a given blob
func (s *Ed25519Signer) Sign(request *pb.SignatureRequest) (*pb.Signature, error) {
	priv := [ed25519.PrivateKeySize]byte{}
	copy(priv[:], s.privateKey.Private[:])
	sig := ed25519.Sign(&priv, request.Content)

	return &pb.Signature{KeyInfo: &pb.KeyInfo{ID: s.privateKey.ID, Algorithm: &pb.Algorithm{Algorithm: ED25519}}, Content: sig[:]}, nil
}

// NewEd25519Signer returns a Ed25519Signer, given a private key
func NewEd25519Signer(key *keys.Key) *Ed25519Signer {
	return &Ed25519Signer{privateKey: key}
}
