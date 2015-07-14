package api_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/agl/ed25519"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/signer/keys"
	"github.com/stretchr/testify/assert"

	pb "github.com/docker/notary/proto"
)

type zeroReader struct{}

func (zeroReader) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

func TestSign(t *testing.T) {
	var zero zeroReader
	public, private, _ := ed25519.GenerateKey(zero)

	blob := []byte("test message")

	directSig := ed25519.Sign(private, blob)
	directSigHex := hex.EncodeToString(directSig[:])

	signer := api.NewEd25519Signer(&keys.Key{Private: private, Public: *public, ID: "fake_id"})

	sigRequest := &pb.SignatureRequest{KeyID: &pb.KeyID{ID: "fake_id"}, Content: blob}

	sig, err := signer.Sign(sigRequest)
	assert.Nil(t, err)
	signatureHex := fmt.Sprintf("%x", sig.Content)

	assert.Equal(t, directSigHex, signatureHex)
	assert.Equal(t, sig.KeyInfo.KeyID.ID, "fake_id")
}

func BenchmarkSign(b *testing.B) {
	blob := []byte("7d16f1d0b95310a7bc557747fc4f20fcd41c1c5095ae42f189df")

	keyDB := keys.NewKeyDB()
	var sigService = api.NewEdDSASigningService(keyDB)

	key, _ := sigService.CreateKey()
	privkey, _ := keyDB.GetKey(key.KeyInfo.KeyID)

	signer := api.NewEd25519Signer(privkey)
	sigRequest := &pb.SignatureRequest{KeyID: key.KeyInfo.KeyID, Content: blob}

	for n := 0; n < b.N; n++ {
		_, _ = signer.Sign(sigRequest)
	}
}
