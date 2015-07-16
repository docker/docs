package api

import (
	"github.com/agl/ed25519"
	"github.com/endophage/gotuf/data"

	pb "github.com/docker/notary/proto"
)

// Ed25519Signer implements the Signer interface for Ed25519 keys
type Ed25519Signer struct {
	privateKey data.Key
}

// Sign returns a signature for a given blob
func (s *Ed25519Signer) Sign(request *pb.SignatureRequest) (*pb.Signature, error) {
	priv := [ed25519.PrivateKeySize]byte{}
	copy(priv[:], s.privateKey.Private())
	sig := ed25519.Sign(&priv, request.Content)

	return &pb.Signature{KeyInfo: &pb.KeyInfo{KeyID: &pb.KeyID{ID: s.privateKey.ID()}, Algorithm: &pb.Algorithm{Algorithm: data.ED25519Key.String()}}, Content: sig[:]}, nil
}

// NewEd25519Signer returns a Ed25519Signer, given a private key
func NewEd25519Signer(key data.Key) *Ed25519Signer {
	return &Ed25519Signer{privateKey: key}
}
