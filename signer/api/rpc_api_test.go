package api_test

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/pkg/passphrase"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/docker/notary/proto"
)

var (
	kmClient   pb.KeyManagementClient
	sClient    pb.SignerClient
	grpcServer *grpc.Server
	void       *pb.Void
	pr         passphrase.Retriever
    health=    map[string]string {
        "db": "ok",
        "other": "not ok",
    }
)

func init() {
	pr = func(string, string, bool, int) (string, bool, error) { return "passphrase", false, nil }
	keyStore := trustmanager.NewKeyMemoryStore(pr)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	cryptoServices := signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService}
	void = &pb.Void{}
	//server setup
	kms := &api.KeyManagementServer{CryptoServices: cryptoServices}
	ss := &api.SignerServer{CryptoServices: cryptoServices}
	grpcServer = grpc.NewServer()
	pb.RegisterKeyManagementServer(grpcServer, kms)
	pb.RegisterSignerServer(grpcServer, ss)
	lis, err := net.Listen("tcp", "127.0.0.1:7899")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	go grpcServer.Serve(lis)

	//client setup
	conn, err := grpc.Dial("127.0.0.1:7899")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

    fakeHealth := func() map[string]string {
        return health
    }

	kmClient = pb.NewKeyManagementClient(conn, fakeHealth)
	sClient = pb.NewSignerClient(conn, fakeHealth)
}

func TestDeleteKeyHandlerReturnsNotFoundWithNonexistentKey(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}

	ret, err := kmClient.DeleteKey(context.Background(), keyID)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.NotFound)
	assert.Nil(t, ret)
}

func TestCreateKeyHandlerCreatesKey(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	assert.NotNil(t, publicKey)
	assert.NotEmpty(t, publicKey.PublicKey)
	assert.NotEmpty(t, publicKey.KeyInfo)
	assert.Nil(t, err)
	assert.Equal(t, grpc.Code(err), codes.OK)
}

func TestDeleteKeyHandlerDeletesCreatedKey(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	ret, err := kmClient.DeleteKey(context.Background(), publicKey.KeyInfo.KeyID)
	assert.Nil(t, err)
	assert.Equal(t, ret, void)
}

func TestKeyInfoReturnsCreatedKeys(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	fmt.Println("Pubkey ID: " + publicKey.GetKeyInfo().KeyID.ID)
	returnedPublicKey, err := kmClient.GetKeyInfo(context.Background(), publicKey.KeyInfo.KeyID)
	fmt.Println("returnedPublicKey ID: " + returnedPublicKey.GetKeyInfo().KeyID.ID)

	assert.Nil(t, err)
	assert.Equal(t, publicKey.KeyInfo, returnedPublicKey.KeyInfo)
	assert.Equal(t, publicKey.PublicKey, returnedPublicKey.PublicKey)
}

func TestCreateKeyCreatesNewKeys(t *testing.T) {
	publicKey1, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	assert.Nil(t, err)
	publicKey2, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	assert.Nil(t, err)
	assert.NotEqual(t, publicKey1, publicKey2)
	assert.NotEqual(t, publicKey1.KeyInfo, publicKey2.KeyInfo)
	assert.NotEqual(t, publicKey1.PublicKey, publicKey2.PublicKey)
}

func TestGetKeyInfoReturnsNotFoundOnNonexistKeys(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}

	ret, err := kmClient.GetKeyInfo(context.Background(), keyID)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.NotFound)
	assert.Nil(t, ret)
}

func TestCreatedKeysCanBeUsedToSign(t *testing.T) {
	message := []byte{0, 0, 0, 0}

	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key.String()})
	assert.Nil(t, err)
	assert.NotNil(t, publicKey)

	sr := &pb.SignatureRequest{Content: message, KeyID: publicKey.KeyInfo.KeyID}
	assert.NotNil(t, sr)
	signature, err := sClient.Sign(context.Background(), sr)
	assert.Nil(t, err)
	assert.NotNil(t, signature)
	assert.NotEmpty(t, signature.Content)
	assert.Equal(t, publicKey.KeyInfo, signature.KeyInfo)
}

func TestSignReturnsNotFoundOnNonexistKeys(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}
	message := []byte{0, 0, 0, 0}
	sr := &pb.SignatureRequest{Content: message, KeyID: keyID}

	ret, err := sClient.Sign(context.Background(), sr)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.NotFound)
	assert.Nil(t, ret)
}
