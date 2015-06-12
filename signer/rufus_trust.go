package signer

import (
	"errors"
	"net"

	"github.com/Sirupsen/logrus"
	pb "github.com/docker/rufus/proto"
	"github.com/endophage/gotuf/data"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// RufusSigner implements a RPC based Trust service that calls the Rufus Service
type RufusSigner struct {
	kmClient pb.KeyManagementClient
	sClient  pb.SignerClient
}

func NewRufusSigner(hostname string, port string, tlscafile string) *RufusSigner {
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
	return &RufusSigner{
		kmClient: kmClient,
		sClient:  sClient,
	}
}

// addKey allows you to add a private key to the trust service
func (trust *RufusSigner) addKey(k *data.PrivateKey) error {
	return errors.New("Not implemented: RufusSigner.addKey")
}

// RemoveKey allows you to remove a private key from the trust service
func (trust *RufusSigner) RemoveKey(keyID string) error {
	toBeDeletedKeyID := &pb.KeyID{ID: keyID}
	_, err := trust.kmClient.DeleteKey(context.Background(), toBeDeletedKeyID)
	return err
}

// Sign signs a byte string with a number of KeyIDs
func (trust *RufusSigner) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
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
			KeyID:     sig.KeyID.ID,
			Method:    "TODOALGORITHM",
			Signature: sig.Content,
		})
	}
	return signatures, nil
}

// Create creates a remote key and returns the PublicKey associated with the remote private key
func (trust *RufusSigner) Create() (*data.PublicKey, error) {
	publicKey, err := trust.kmClient.CreateKey(context.Background(), &pb.Void{})
	if err != nil {
		return nil, err
	}
	//TODO(mccauley): Update API to return algorithm and/or take it as a param
	public := data.NewPublicKey("TODOALGORITHM", string(publicKey.PublicKey))
	return public, nil
}

// PublicKeys returns the public key(s) associated with the passed in keyIDs
func (trust *RufusSigner) PublicKeys(keyIDs ...string) (map[string]*data.PublicKey, error) {
	publicKeys := make(map[string]*data.PublicKey)
	for _, ID := range keyIDs {
		keyID := pb.KeyID{ID: ID}
		sig, err := trust.kmClient.GetKeyInfo(context.Background(), &keyID)
		if err != nil {
			return nil, err
		}
		publicKeys[sig.KeyID.ID] =
			data.NewPublicKey("TODOALGORITHM", string(sig.PublicKey))
	}
	return publicKeys, nil
}

func (trust *RufusSigner) CanSign(kID string) bool {
	return true
}
