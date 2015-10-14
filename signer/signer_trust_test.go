package signer

import (
	"errors"
	"testing"

	"google.golang.org/grpc"

	pb "github.com/docker/notary/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type StubKeyManagementClient struct {
	pb.KeyManagementClient
	status map[string]string
	err    error
}

func (c StubKeyManagementClient) CheckHealth(x context.Context, v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &pb.HealthStatus{c.status}, nil
}

type StubSignerClient struct {
	pb.SignerClient
	status map[string]string
	err    error
}

func (c StubSignerClient) CheckHealth(x context.Context, v *pb.Void, o ...grpc.CallOption) (*pb.HealthStatus, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &pb.HealthStatus{c.status}, nil
}

// KMHealth only succeeds if the KeyManagement service is healthy.
func TestKMHealthFailure(t *testing.T) {
	errs := []error{errors.New("Connection failure!"), nil}
	status := map[string]string{"db": "bad"}

	for _, err := range errs {
		signer := NotarySigner{
			StubKeyManagementClient{
				pb.NewKeyManagementClient(nil),
				status,
				err,
			},
			StubSignerClient{
				pb.NewSignerClient(nil),
				make(map[string]string),
				nil,
			},
		}
		err := signer.KMHealth()
		assert.Error(t, err)
	}
}

// SHealth only succeeds if the Signer service is healthy.
func TestSignerHealthFailure(t *testing.T) {
	errs := []error{errors.New("Connection failure!"), nil}
	status := map[string]string{"db": "bad"}

	for _, err := range errs {
		signer := NotarySigner{
			StubKeyManagementClient{
				pb.NewKeyManagementClient(nil),
				make(map[string]string),
				nil,
			},
			StubSignerClient{
				pb.NewSignerClient(nil),
				status,
				err,
			},
		}
		err := signer.SHealth()
		assert.Error(t, err)
	}
}

// SHealth succeeds if the signer service is healthy, regardless of the
// key management service.
func TestSignerHealthGood(t *testing.T) {
	signer := NotarySigner{
		StubKeyManagementClient{
			pb.NewKeyManagementClient(nil),
			nil,
			errors.New("Connection failure!"),
		},
		StubSignerClient{
			pb.NewSignerClient(nil),
			make(map[string]string),
			nil,
		},
	}
	err := signer.SHealth()
	assert.NoError(t, err)
}

// KMHealth succeeds if the key management service is healthy, regardless of
// the signer service.
func TestKMHealthGood(t *testing.T) {
	signer := NotarySigner{
		StubKeyManagementClient{
			pb.NewKeyManagementClient(nil),
			make(map[string]string),
			nil,
		},
		StubSignerClient{
			pb.NewSignerClient(nil),
			nil,
			errors.New("Connection failure!"),
		},
	}
	err := signer.KMHealth()
	assert.NoError(t, err)
}
