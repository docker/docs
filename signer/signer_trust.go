package signer

import (
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	pb "github.com/docker/notary/proto"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	tlsConfig, err := utils.ConfigureClientTLS(&utils.ClientTLSOpts{
		RootCAFile: tlscafile,
		ServerName: hostname,
	})
	if err != nil {
		logrus.Fatal("Unable to set up TLS: ", err)
	}
	creds := credentials.NewTLS(tlsConfig)
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

// CheckHealth checks the health of one of the clients, since both clients run
// from the same GRPC server.
func (trust *NotarySigner) CheckHealth(timeout time.Duration) error {

	// Do not bother starting checking at all if the connection is broken.
	if trust.clientConn.State() != grpc.Idle &&
		trust.clientConn.State() != grpc.Ready {
		return fmt.Errorf("Not currently connected to trust server.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	status, err := trust.kmClient.CheckHealth(ctx, &pb.Void{})
	defer cancel()
	if err == nil && len(status.Status) > 0 {
		return fmt.Errorf("Trust is not healthy")
	} else if err != nil && grpc.Code(err) == codes.DeadlineExceeded {
		return fmt.Errorf(
			"Timed out reaching trust service after %s.", timeout)
	}
	return err
}
