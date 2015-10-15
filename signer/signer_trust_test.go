package signer

import (
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "github.com/docker/notary/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type StubKeyManagementClient struct {
	pb.KeyManagementClient
	healthCheck func() (map[string]string, error)
}

func (c StubKeyManagementClient) CheckHealth(x context.Context, v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	status, err := c.healthCheck()
	if err != nil {
		return nil, err
	}
	return &pb.HealthStatus{status}, nil
}

type StubSignerClient struct {
	pb.SignerClient
	healthCheck func() (map[string]string, error)
}

func (c StubSignerClient) CheckHealth(x context.Context, v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	status, err := c.healthCheck()
	if err != nil {
		return nil, err
	}
	return &pb.HealthStatus{status}, nil
}

type StubGRPCConnection struct {
	fakeConnStatus grpc.ConnectivityState
}

func (c StubGRPCConnection) State() grpc.ConnectivityState {
	return c.fakeConnStatus
}

type healthSideEffect func() (map[string]string, error)

func healthOk() (map[string]string, error) {
	return make(map[string]string), nil
}

func healthBad() (map[string]string, error) {
	return map[string]string{"health": "not good"}, nil
}

func healthError() (map[string]string, error) {
	return nil, errors.New("Something's wrong")
}

func healthTimeout() (map[string]string, error) {
	time.Sleep(time.Second * 10)
	return healthOk()
}

func makeSigner(kmFunc healthSideEffect, sFunc healthSideEffect, conn StubGRPCConnection) NotarySigner {
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
	signer := makeSigner(healthBad, healthOk, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the health check to the KM server errors
func TestHealthCheckKMError(t *testing.T) {
	signer := makeSigner(healthBad, healthOk, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the health check to the KM server times out
func TestHealthCheckKMTimeout(t *testing.T) {
	signer := makeSigner(healthTimeout, healthOk, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the signer is unhealthy
func TestHealthCheckSignerUnhealthy(t *testing.T) {
	signer := makeSigner(healthOk, healthBad, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the health check to the signer errors
func TestHealthCheckSignerError(t *testing.T) {
	signer := makeSigner(healthOk, healthBad, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth does not succeed if the health check to the signer times out
func TestHealthCheckSignerTimeout(t *testing.T) {
	signer := makeSigner(healthOk, healthTimeout, StubGRPCConnection{})
	assert.Error(t, signer.CheckHealth(1))
}

// CheckHealth succeeds if both services are healthy and reachable
func TestHealthCheckBothHealthy(t *testing.T) {
	signer := makeSigner(healthOk, healthOk, StubGRPCConnection{})
	assert.NoError(t, signer.CheckHealth(1))
}

// CheckHealth fails immediately if not connected to the server.
func TestHealthCheckConnectionDied(t *testing.T) {
	signer := makeSigner(healthTimeout, healthTimeout,
		StubGRPCConnection{grpc.Connecting})
	assert.Error(t, signer.CheckHealth(30))
}
