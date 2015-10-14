package signer

import (
	"errors"
	"net"

	"github.com/Sirupsen/logrus"
	pb "github.com/docker/notary/proto"
	"github.com/endophage/gotuf/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NotarySigner implements a RPC based Trust service that calls the Notary-signer Service
type NotarySigner struct {
	kmClient pb.KeyManagementClient
	sClient  pb.SignerClient
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
		kmClient: kmClient,
		sClient:  sClient,
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

// KMHealth returns the health of the KeyManagement RPC service.
func (trust *NotarySigner) KMHealth() error {
	status, err := trust.kmClient.CheckHealth(context.Background(), &pb.Void{})
	if err != nil {
		return err
	}
	if len(status.Status) > 0 {
		return errors.New("KeyManagement not healthy")
	}
	return nil
}

// SHealth returns the health of the Signer RPC service.
func (trust *NotarySigner) SHealth() error {
	status, err := trust.sClient.CheckHealth(context.Background(), &pb.Void{})
	if err != nil {
		return err
	}
	if len(status.Status) > 0 {
		return errors.New("Signer not healthy")
	}
	return nil
}
