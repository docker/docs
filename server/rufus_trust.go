package server

import (
	"errors"
	"log"

	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/keys"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/docker/rufus/proto"
)

// RufusSigner implements a simple in memory keystore and trust service
type RufusSigner struct {
	kmClient pb.KeyManagementClient
	sClient  pb.SignerClient
}

func newRufusSigner(hostNameAndPort string) *RufusSigner {
	conn, err := grpc.Dial(hostNameAndPort)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	kmClient := pb.NewKeyManagementClient(conn)
	sClient := pb.NewSignerClient(conn)
	return &RufusSigner{
		kmClient: kmClient,
		sClient:  sClient,
	}
}

// addKey allows you to add a private key to the trust service
func (trust *RufusSigner) addKey(k *keys.PrivateKey) error {
	return errors.New("Not implemented: RufusSigner.addKey")
}

func (trust *RufusSigner) RemoveKey(keyID string) error {
	toBeDeletedKeyID := &pb.KeyID{ID: keyID}
	_, err := trust.kmClient.DeleteKey(context.Background(), toBeDeletedKeyID)
	return err
}

func (trust *RufusSigner) Sign(keyIDs []string, toSign []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, kID := range keyIDs {
		keyID := pb.KeyID{ID: kID}
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

func (trust *RufusSigner) Create() (*keys.PublicKey, error) {
	publicKey, err := trust.kmClient.CreateKey(context.Background(), &pb.Void{})
	if err != nil {
		return nil, err
	}
	//TODO(mccauley): Update API to return algorithm and/or take it as a param
	public := keys.NewPublicKey("TODOALGORITHM", publicKey.PublicKey)
	return public, nil
}

func (trust *RufusSigner) PublicKeys(keyIDs ...string) (map[string]*keys.PublicKey, error) {
	publicKeys := make(map[string]*keys.PublicKey)
	for _, kID := range keyIDs {
		keyID := pb.KeyID{ID: kID}
		sig, err := trust.kmClient.GetKeyInfo(context.Background(), &keyID)
		if err != nil {
			return nil, err
		}
		publicKeys[sig.KeyID.ID] =
			keys.NewPublicKey("TODOALGORITHM", sig.PublicKey)
	}
	return publicKeys, nil
}
