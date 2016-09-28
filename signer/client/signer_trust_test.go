package client

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	pb "github.com/docker/notary/proto"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils/interfaces"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type rpcHealthCheck func(
	context.Context, *pb.Void, ...grpc.CallOption) (*pb.HealthStatus, error)

type StubKeyManagementClient struct {
	pb.KeyManagementClient
	healthCheck rpcHealthCheck
}

func (c StubKeyManagementClient) CheckHealth(x context.Context,
	v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	return c.healthCheck(x, v, o...)
}

type StubGRPCConnection struct {
	fakeConnStatus grpc.ConnectivityState
}

func (c StubGRPCConnection) State() grpc.ConnectivityState {
	return c.fakeConnStatus
}

func stubHealthFunction(t *testing.T, status map[string]string, err error) rpcHealthCheck {
	return func(ctx context.Context, v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
		_, withDeadline := ctx.Deadline()
		require.True(t, withDeadline)

		return &pb.HealthStatus{Status: status}, err
	}
}

func makeSigner(kmFunc rpcHealthCheck, conn StubGRPCConnection) NotarySigner {
	return NotarySigner{
		StubKeyManagementClient{
			pb.NewKeyManagementClient(nil),
			kmFunc,
		},
		pb.NewSignerClient(nil),
		conn,
	}
}

// CheckHealth does not succeed if the KM server is unhealthy
func TestHealthCheckKMUnhealthy(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, map[string]string{"health": "not good"}, nil),
		StubGRPCConnection{})
	require.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server errors
func TestHealthCheckKMError(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, nil, errors.New("Something's wrong")),
		StubGRPCConnection{})
	require.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server times out
func TestHealthCheckKMTimeout(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, nil, grpc.Errorf(codes.DeadlineExceeded, "")),
		StubGRPCConnection{})
	err := signer.CheckHealth(1 * time.Second)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Timed out"))
}

// CheckHealth succeeds if KM is healthy and reachable.
func TestHealthCheckKMHealthy(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, make(map[string]string), nil),
		StubGRPCConnection{})
	require.NoError(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth fails immediately if not connected to the server.
func TestHealthCheckConnectionDied(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, make(map[string]string), nil),
		StubGRPCConnection{grpc.Connecting})
	require.Error(t, signer.CheckHealth(1*time.Second))
}

var ret = passphrase.ConstantRetriever("pass")

func TestGetPrivateKeyAndSignWithExistingKey(t *testing.T) {
	key, err := trustmanager.GenerateECDSAKey(rand.Reader)
	require.NoError(t, err, "could not generate key")

	store := trustmanager.NewKeyMemoryStore(ret)

	err = store.AddKey(trustmanager.KeyInfo{Role: data.CanonicalTimestampRole, Gun: "gun"}, key)
	require.NoError(t, err, "could not add key to store")

	signer := setUpSigner(t, store)

	privKey, _, err := signer.GetPrivateKey(key.ID())
	require.NoError(t, err)
	require.NotNil(t, privKey)

	msg := []byte("message!")
	sig, err := privKey.Sign(rand.Reader, msg, nil)
	require.NoError(t, err)

	err = signed.Verifiers[data.ECDSASignature].Verify(
		data.PublicKeyFromPrivate(key), sig, msg)
	require.NoError(t, err)
}

// Signer conforms to the signed.CryptoService interface behavior
func TestCryptoSignerInterfaceBehavior(t *testing.T) {
	signer := setUpSigner(t, trustmanager.NewKeyMemoryStore(ret))
	interfaces.EmptyCryptoServiceInterfaceBehaviorTests(t, &signer)
	interfaces.CreateGetKeyCryptoServiceInterfaceBehaviorTests(t, &signer, data.ECDSAKey, false)
	// can't test AddKey, because the signer does not support adding keys, and can't test listing
	// keys because the signer doesn't support listing keys.  Signer also doesn't support tracking
	// roles
}

type StubClientFromServers struct {
	api.KeyManagementServer
	api.SignerServer
}

func (c *StubClientFromServers) CreateKey(ctx context.Context,
	algorithm *pb.Algorithm, _ ...grpc.CallOption) (*pb.PublicKey, error) {
	return c.KeyManagementServer.CreateKey(ctx, algorithm)
}

func (c *StubClientFromServers) DeleteKey(ctx context.Context, keyID *pb.KeyID,
	_ ...grpc.CallOption) (*pb.Void, error) {
	return c.KeyManagementServer.DeleteKey(ctx, keyID)
}

func (c *StubClientFromServers) GetKeyInfo(ctx context.Context, keyID *pb.KeyID,
	_ ...grpc.CallOption) (*pb.PublicKey, error) {
	return c.KeyManagementServer.GetKeyInfo(ctx, keyID)
}

func (c *StubClientFromServers) Sign(ctx context.Context,
	sr *pb.SignatureRequest, _ ...grpc.CallOption) (*pb.Signature, error) {
	return c.SignerServer.Sign(ctx, sr)
}

func (c *StubClientFromServers) CheckHealth(ctx context.Context, v *pb.Void,
	_ ...grpc.CallOption) (*pb.HealthStatus, error) {
	return c.KeyManagementServer.CheckHealth(ctx, v)
}

func setUpSigner(t *testing.T, store trustmanager.KeyStore) NotarySigner {
	cryptoService := cryptoservice.NewCryptoService(store)
	cryptoServices := signer.CryptoServiceIndex{
		data.ED25519Key: cryptoService,
		data.RSAKey:     cryptoService,
		data.ECDSAKey:   cryptoService,
	}

	fakeHealth := func() map[string]string { return map[string]string{} }

	client := StubClientFromServers{
		KeyManagementServer: api.KeyManagementServer{CryptoServices: cryptoServices,
			HealthChecker: fakeHealth},
		SignerServer: api.SignerServer{CryptoServices: cryptoServices,
			HealthChecker: fakeHealth},
	}

	return NotarySigner{kmClient: &client, sClient: &client}
}
