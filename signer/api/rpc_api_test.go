package api_test

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
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
	health     = map[string]string{
		"db":    "ok",
		"other": "not ok",
	}
)

func init() {
	pr = func(string, string, bool, int) (string, bool, error) { return "passphrase", false, nil }
	keyStore := trustmanager.NewKeyMemoryStore(pr)
	cryptoService := cryptoservice.NewCryptoService(keyStore)
	cryptoServices := signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService}
	void = &pb.Void{}

	fakeHealth := func() map[string]string {
		return health
	}

	//server setup
	kms := &api.KeyManagementServer{CryptoServices: cryptoServices,
		HealthChecker: fakeHealth}
	ss := &api.SignerServer{CryptoServices: cryptoServices,
		HealthChecker: fakeHealth}
	grpcServer = grpc.NewServer()
	pb.RegisterKeyManagementServer(grpcServer, kms)
	pb.RegisterSignerServer(grpcServer, ss)
	lis, err := net.Listen("tcp", "127.0.0.1:7899")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	go grpcServer.Serve(lis)

	//client setup
	conn, err := grpc.Dial("127.0.0.1:7899", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	kmClient = pb.NewKeyManagementClient(conn)
	sClient = pb.NewSignerClient(conn)
}

func TestDeleteKeyHandlerReturnsNilWithNonexistentKey(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}

	_, err := kmClient.DeleteKey(context.Background(), keyID)
	require.NoError(t, err)
}

func TestCreateKeyHandlerCreatesKey(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	require.NotNil(t, publicKey)
	require.NotEmpty(t, publicKey.PublicKey)
	require.NotEmpty(t, publicKey.KeyInfo)
	require.Nil(t, err)
	require.Equal(t, grpc.Code(err), codes.OK)
}

func TestDeleteKeyHandlerDeletesCreatedKey(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	ret, err := kmClient.DeleteKey(context.Background(), publicKey.KeyInfo.KeyID)
	require.Nil(t, err)
	require.Equal(t, ret, void)
}

func TestKeyInfoReturnsCreatedKeys(t *testing.T) {
	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	fmt.Println("Pubkey ID: " + publicKey.GetKeyInfo().KeyID.ID)
	returnedPublicKey, err := kmClient.GetKeyInfo(context.Background(), publicKey.KeyInfo.KeyID)
	fmt.Println("returnedPublicKey ID: " + returnedPublicKey.GetKeyInfo().KeyID.ID)

	require.Nil(t, err)
	require.Equal(t, publicKey.KeyInfo, returnedPublicKey.KeyInfo)
	require.Equal(t, publicKey.PublicKey, returnedPublicKey.PublicKey)
}

func TestCreateKeyCreatesNewKeys(t *testing.T) {
	publicKey1, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	require.Nil(t, err)
	publicKey2, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	require.Nil(t, err)
	require.NotEqual(t, publicKey1, publicKey2)
	require.NotEqual(t, publicKey1.KeyInfo, publicKey2.KeyInfo)
	require.NotEqual(t, publicKey1.PublicKey, publicKey2.PublicKey)
}

func TestGetKeyInfoReturnsNotFoundOnNonexistKeys(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}

	ret, err := kmClient.GetKeyInfo(context.Background(), keyID)
	require.NotNil(t, err)
	require.Equal(t, grpc.Code(err), codes.NotFound)
	require.Nil(t, ret)
}

func TestCreatedKeysCanBeUsedToSign(t *testing.T) {
	message := []byte{0, 0, 0, 0}

	publicKey, err := kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: data.ED25519Key})
	require.Nil(t, err)
	require.NotNil(t, publicKey)

	sr := &pb.SignatureRequest{Content: message, KeyID: publicKey.KeyInfo.KeyID}
	require.NotNil(t, sr)
	signature, err := sClient.Sign(context.Background(), sr)
	require.Nil(t, err)
	require.NotNil(t, signature)
	require.NotEmpty(t, signature.Content)
	require.Equal(t, publicKey.KeyInfo, signature.KeyInfo)
}

func TestSignReturnsNotFoundOnNonexistKeys(t *testing.T) {
	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyID := &pb.KeyID{ID: fakeID}
	message := []byte{0, 0, 0, 0}
	sr := &pb.SignatureRequest{Content: message, KeyID: keyID}

	ret, err := sClient.Sign(context.Background(), sr)
	require.NotNil(t, err)
	require.Equal(t, grpc.Code(err), codes.NotFound)
	require.Nil(t, ret)
}

func TestHealthChecksForServices(t *testing.T) {
	sHealthStatus, err := sClient.CheckHealth(context.Background(), void)
	require.Nil(t, err)
	require.Equal(t, health, sHealthStatus.Status)

	kmHealthStatus, err := kmClient.CheckHealth(context.Background(), void)
	require.Nil(t, err)
	require.Equal(t, health, kmHealthStatus.Status)
}
