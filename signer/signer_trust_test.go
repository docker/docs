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
		assert.True(t, withDeadline)

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
	assert.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server errors
func TestHealthCheckKMError(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, nil, errors.New("Something's wrong")),
		StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth does not succeed if the health check to the KM server times out
func TestHealthCheckKMTimeout(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, nil, grpc.Errorf(codes.DeadlineExceeded, "")),
		StubGRPCConnection{})
	err := signer.CheckHealth(1 * time.Second)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Timed out"))
}

// CheckHealth succeeds if KM is healthy and reachable.
func TestHealthCheckKMHealthy(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, make(map[string]string), nil),
		StubGRPCConnection{})
	assert.NoError(t, signer.CheckHealth(1*time.Second))
}

// CheckHealth fails immediately if not connected to the server.
func TestHealthCheckConnectionDied(t *testing.T) {
	signer := makeSigner(
		stubHealthFunction(t, make(map[string]string), nil),
		StubGRPCConnection{grpc.Connecting})
	assert.Error(t, signer.CheckHealth(1*time.Second))
}
