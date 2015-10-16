package signer

import (
	"errors"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/docker/notary/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type StubKeyManagementClient struct {
	pb.KeyManagementClient
	healthCheck rpcHealthCheck
}

func (c StubKeyManagementClient) CheckHealth(x context.Context,
	v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	return c.healthCheck(x, v, o...)
}

type StubSignerClient struct {
	pb.SignerClient
	healthCheck rpcHealthCheck
}

func (c StubSignerClient) CheckHealth(x context.Context, v *pb.Void,
	o ...grpc.CallOption) (*pb.HealthStatus, error) {
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
		assert.True(t, withDeadline)

		return &pb.HealthStatus{status}, err
	}
}

func healthOk(t *testing.T) rpcHealthCheck {
	return stubHealthFunction(t, make(map[string]string), nil)
}

func healthBad(t *testing.T) rpcHealthCheck {
	return stubHealthFunction(t, map[string]string{"health": "not good"}, nil)
}

func healthError(t *testing.T) rpcHealthCheck {
	return stubHealthFunction(t, nil, errors.New("Something's wrong"))
}

func healthTimeout(t *testing.T) rpcHealthCheck {
	return stubHealthFunction(
		t, nil, grpc.Errorf(codes.DeadlineExceeded, ""))
}

func makeSigner(kmFunc rpcHealthCheck, sFunc rpcHealthCheck, conn StubGRPCConnection) NotarySigner {
	return NotarySigner{
		StubKeyManagementClient{
			pb.NewKeyManagementClient(nil),
			kmFunc,
		},
		StubSignerClient{
			pb.NewSignerClient(nil),
			sFunc,
		},
		conn,
	}
}

// CheckHealth does not succeed if the KM server is unhealthy
func TestHealthCheckKMUnhealthy(t *testing.T) {
	signer := makeSigner(healthBad(t), healthOk(t), StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server errors
func TestHealthCheckKMError(t *testing.T) {
	signer := makeSigner(healthBad(t), healthOk(t), StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server times out
func TestHealthCheckKMTimeout(t *testing.T) {
	signer := makeSigner(healthTimeout(t), healthOk(t), StubGRPCConnection{})
	err := signer.CheckHealth(1 * time.Second)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Timed out"))
}

// CheckHealth does not succeed if the signer is unhealthy
func TestHealthCheckSignerUnhealthy(t *testing.T) {
	signer := makeSigner(healthOk(t), healthBad(t), StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the signer errors
func TestHealthCheckSignerError(t *testing.T) {
	signer := makeSigner(healthOk(t), healthBad(t), StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the health check to the signer times out
func TestHealthCheckSignerTimeout(t *testing.T) {
	signer := makeSigner(healthOk(t), healthTimeout(t), StubGRPCConnection{})
	err := signer.CheckHealth(1 * time.Second)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Timed out"))
}

// CheckHealth succeeds if both services are healthy and reachable
func TestHealthCheckBothHealthy(t *testing.T) {
	signer := makeSigner(healthOk(t), healthOk(t), StubGRPCConnection{})
	assert.NoError(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth fails immediately if not connected to the server.
func TestHealthCheckConnectionDied(t *testing.T) {
	signer := makeSigner(healthOk(t), healthOk(t),
		StubGRPCConnection{grpc.Connecting})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}
