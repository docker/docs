package signer

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	pb "github.com/docker/notary/proto"
	"github.com/endophage/gotuf/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// The only thing needed from grpc.ClientConn is it's state.
type checkableConnectionState interface {
	State() grpc.ConnectivityState
}

// NotarySigner implements a RPC based Trust service that calls the Notary-signer Service
type NotarySigner struct {
	kmClient   pb.KeyManagementClient
	sClient    pb.SignerClient
	clientConn checkableConnectionState
}

// NewNotarySigner is a convinience method that returns NotarySigner
func NewNotarySigner(hostname string, port string, tlscafile string) *NotarySigner {
	var opts []grpc.DialOption
	netAddr := net.JoinHostPort(hostname, port)
	creds, err := credentials.NewClientTLSFromFile(tlscafile, hostname)
	if err != nil {
		logrus.Fatal("fail to read: ", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(netAddr, opts...)

	if err != nil {
		logrus.Fatal("fail to dial: ", err)
	}
	kmClient := pb.NewKeyManagementClient(conn)
	sClient := pb.NewSignerClient(conn)
	return &NotarySigner{
		kmClient:   kmClient,
		sClient:    sClient,
		clientConn: conn,
	}
}

// Sign signs a byte string with a number of KeyIDs
func (trust *NotarySigner) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, ID := range keyIDs {
		keyID := pb.KeyID{ID: ID}
		sr := &pb.SignatureRequest{
			Content: toSign,
			KeyID:   &keyID,
		}
		sig, err := trust.sClient.Sign(context.Background(), sr)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, data.Signature{
			KeyID:     sig.KeyInfo.KeyID.ID,
			Method:    data.SigAlgorithm(sig.Algorithm.Algorithm),
			Signature: sig.Content,
		})
	}
	return signatures, nil
}

// Create creates a remote key and returns the PublicKey associated with the remote private key
func (trust *NotarySigner) Create(role string, algorithm data.KeyAlgorithm) (data.PublicKey, error) {
	publicKey, err := trust.kmClient.CreateKey(context.Background(), &pb.Algorithm{Algorithm: algorithm.String()})
	if err != nil {
		return nil, err
	}
	public := data.NewPublicKey(data.KeyAlgorithm(publicKey.KeyInfo.Algorithm.Algorithm), publicKey.PublicKey)
	return public, nil
}

// RemoveKey deletes a key
func (trust *NotarySigner) RemoveKey(keyid string) error {
	_, err := trust.kmClient.DeleteKey(context.Background(), &pb.KeyID{ID: keyid})
	return err
}

// GetKey retrieves a key
func (trust *NotarySigner) GetKey(keyid string) data.PublicKey {
	publicKey, err := trust.kmClient.GetKeyInfo(context.Background(), &pb.KeyID{ID: keyid})
	if err != nil {
		return nil
	}
	return data.NewPublicKey(data.KeyAlgorithm(publicKey.KeyInfo.Algorithm.Algorithm), publicKey.PublicKey)
}

// CheckHealth returns true if trust is healthy - that is if it can connect
// to the trust server, and both the key management and signer services are
// healthy
func (trust *NotarySigner) CheckHealth(timeout int) error {
	if e := trust.checkServiceHealth("Key manager", trust.kmClient,
		timeout); e != nil {
		return e
	}
	return trust.checkServiceHealth("Signer", trust.sClient, timeout)
}

// Generalized function that can check the health of both the signer client
// and the key management client.
func (trust *NotarySigner) checkServiceHealth(
	serviceName string,
	client interface {
		CheckHealth(context.Context, *pb.Void, ...grpc.CallOption) (*pb.HealthStatus, error)
	},
	timeout int) error {

	// Do not bother starting goroutine if the connection is broken.
	if trust.clientConn.State() != grpc.Idle &&
		trust.clientConn.State() != grpc.Ready {
		return errors.New("Not currently connected to trust server.")
	}

	// We still want to time out getting health, because the connection could
	// have disconnected sometime between when we checked the connection and
	// when we try to make an RPC call.
	channel := make(chan error)
	go func() {
		status, err := client.CheckHealth(context.Background(), &pb.Void{})
		// if this function gets timed out, it might panic when writing to a
		// closed channel
		defer func() { recover() }()
		if err != nil {
			channel <- err
		}
		if len(status.Status) > 0 {
			channel <- fmt.Errorf("%s not healthy", serviceName)
		}
		channel <- nil
	}()
	select {
	case err := <-channel:
		close(channel)
		return err
	case <-time.After(time.Second * time.Duration(timeout)):
		close(channel)
		return fmt.Errorf("Timed out connecting to %s after 60 seconds", serviceName)
	}
}
