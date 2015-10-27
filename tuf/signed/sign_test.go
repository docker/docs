package signed

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

const (
	testKeyPEM1 = "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAnKuXZeefa2LmgxaL5NsM\nzKOHNe+x/nL6ik+lDBCTV6OdcwAhHQS+PONGhrChIUVR6Vth3hUCrreLzPO73Oo5\nVSCuRJ53UronENl6lsa5mFKP8StYLvIDITNvkoT3j52BJIjyNUK9UKY9As2TNqDf\nBEPIRp28ev/NViwGOEkBu2UAbwCIdnDXm8JQErCZA0Ydm7PKGgjLbFsFGrVzqXHK\n6pdzJXlhr9yap3UpgQ/iO9JtoEYB2EXsnSrPc9JRjR30bNHHtnVql3fvinXrAEwq\n3xmN4p+R4VGzfdQN+8Kl/IPjqWB535twhFYEG/B7Ze8IwbygBjK3co/KnOPqMUrM\nBI8ztvPiogz+MvXb8WvarZ6TMTh8ifZI96r7zzqyzjR1hJulEy3IsMGvz8XS2J0X\n7sXoaqszEtXdq5ef5zKVxkiyIQZcbPgmpHLq4MgfdryuVVc/RPASoRIXG4lKaTJj\n1ANMFPxDQpHudCLxwCzjCb+sVa20HBRPTnzo8LSZkI6jAgMBAAE=\n-----END PUBLIC KEY-----"
	testKeyID1  = "51324b59d4888faa91219ebbe5a3876bb4efb21f0602ddf363cd4c3996ded3d4"

	testKeyPEM2 = "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEArvqUPYb6JJROPJQglPTj\n5uDrsxQKl34Mo+3pSlBVuD6puE4lDnG649a2YksJy+C8ZIPJgokn5w+C3alh+dMe\nzbdWHHxrY1h9CLpYz5cbMlE16303ubkt1rvwDqEezG0HDBzPaKj4oP9YJ9x7wbsq\ndvFcy+Qc3wWd7UWcieo6E0ihbJkYcY8chRXVLg1rL7EfZ+e3bq5+ojA2ECM5JqzZ\nzgDpqCv5hTCYYZp72MZcG7dfSPAHrcSGIrwg7whzz2UsEtCOpsJTuCl96FPN7kAu\n4w/WyM3+SPzzr4/RQXuY1SrLCFD8ebM2zHt/3ATLhPnGmyG5I0RGYoegFaZ2AViw\nlqZDOYnBtgDvKP0zakMtFMbkh2XuNBUBO7Sjs0YcZMjLkh9gYUHL1yWS3Aqus1Lw\nlI0gHS22oyGObVBWkZEgk/Foy08sECLGao+5VvhmGpfVuiz9OKFUmtPVjWzRE4ng\niekEu4drSxpH41inLGSvdByDWLpcTvWQI9nkgclh3AT/AgMBAAE=\n-----END PUBLIC KEY-----"
	testKeyID2  = "26f2f5c0fbfa98823bf1ad39d5f3b32575895793baf80f1df675597d5b95dba8"
)

type FailingCryptoService struct {
	testKey data.PublicKey
}

func (mts *FailingCryptoService) Sign(keyIDs []string, _ []byte) ([]data.Signature, error) {
	sigs := make([]data.Signature, 0, len(keyIDs))
	return sigs, nil
}

func (mts *FailingCryptoService) Create(_ string, _ data.KeyAlgorithm) (data.PublicKey, error) {
	return mts.testKey, nil
}

func (mts *FailingCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == "testID" {
		return mts.testKey
	}
	return nil
}

func (mts *FailingCryptoService) RemoveKey(keyID string) error {
	return nil
}

type MockCryptoService struct {
	testKey data.PublicKey
}

func (mts *MockCryptoService) Sign(keyIDs []string, _ []byte) ([]data.Signature, error) {
	sigs := make([]data.Signature, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		sigs = append(sigs, data.Signature{KeyID: keyID})
	}
	return sigs, nil
}

func (mts *MockCryptoService) Create(_ string, _ data.KeyAlgorithm) (data.PublicKey, error) {
	return mts.testKey, nil
}

func (mts *MockCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == "testID" {
		return mts.testKey
	}
	return nil
}

func (mts *MockCryptoService) RemoveKey(keyID string) error {
	return nil
}

var _ CryptoService = &MockCryptoService{}

type StrictMockCryptoService struct {
	MockCryptoService
}

func (mts *StrictMockCryptoService) Sign(keyIDs []string, _ []byte) ([]data.Signature, error) {
	sigs := make([]data.Signature, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		if keyID == mts.testKey.ID() {
			sigs = append(sigs, data.Signature{KeyID: keyID})
		}
	}
	return sigs, nil
}

func (mts *StrictMockCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == mts.testKey.ID() {
		return mts.testKey
	}
	return nil
}

// Test signing and ensure the expected signature is added
func TestBasicSign(t *testing.T) {
	testKey, _ := pem.Decode([]byte(testKeyPEM1))
	k := data.NewPublicKey(data.RSAKey, testKey.Bytes)
	mockCryptoService := &MockCryptoService{testKey: k}
	key, err := mockCryptoService.Create("root", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	testData := data.Signed{}

	Sign(mockCryptoService, &testData, key)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != testKeyID1 {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}
}

// Test signing with the same key multiple times only registers a single signature
// for the key (N.B. MockCryptoService.Sign will still be called again, but Signer.Sign
// should be cleaning previous signatures by the KeyID when asked to sign again)
func TestReSign(t *testing.T) {
	testKey, _ := pem.Decode([]byte(testKeyPEM1))
	k := data.NewPublicKey(data.RSAKey, testKey.Bytes)
	mockCryptoService := &MockCryptoService{testKey: k}
	testData := data.Signed{}

	Sign(mockCryptoService, &testData, k)
	Sign(mockCryptoService, &testData, k)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != testKeyID1 {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}

}

func TestMultiSign(t *testing.T) {
	mockCryptoService := &MockCryptoService{}
	testData := data.Signed{}

	testKey, _ := pem.Decode([]byte(testKeyPEM1))
	key := data.NewPublicKey(data.RSAKey, testKey.Bytes)
	Sign(mockCryptoService, &testData, key)

	testKey, _ = pem.Decode([]byte(testKeyPEM2))
	key = data.NewPublicKey(data.RSAKey, testKey.Bytes)
	Sign(mockCryptoService, &testData, key)

	if len(testData.Signatures) != 2 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	keyIDs := map[string]struct{}{testKeyID1: struct{}{}, testKeyID2: struct{}{}}
	for _, sig := range testData.Signatures {
		if _, ok := keyIDs[sig.KeyID]; !ok {
			t.Fatalf("Got a signature we didn't expect: %s", sig.KeyID)
		}
	}

}

func TestSignReturnsNoSigs(t *testing.T) {
	failingCryptoService := &FailingCryptoService{}
	testData := data.Signed{}

	testKey, _ := pem.Decode([]byte(testKeyPEM1))
	key := data.NewPublicKey(data.RSAKey, testKey.Bytes)
	err := Sign(failingCryptoService, &testData, key)

	if err == nil {
		t.Fatalf("Expected failure due to no signature being returned by the crypto service")
	}
	if len(testData.Signatures) != 0 {
		t.Fatalf("Incorrect number of signatures, expected 0: %d", len(testData.Signatures))
	}
}

func TestCreate(t *testing.T) {
	testKey, _ := pem.Decode([]byte(testKeyPEM1))
	k := data.NewPublicKey(data.RSAKey, testKey.Bytes)
	mockCryptoService := &MockCryptoService{testKey: k}

	key, err := mockCryptoService.Create("root", data.ED25519Key)

	if err != nil {
		t.Fatal(err)
	}
	if key.ID() != testKeyID1 {
		t.Fatalf("Expected key ID not found: %s", key.ID())
	}
}

func TestSignWithX509(t *testing.T) {
	// generate a key becase we need a cert
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	// make a RSA x509 key
	template, err := trustmanager.NewCertificate("test")
	assert.NoError(t, err)

	derBytes, err := x509.CreateCertificate(
		rand.Reader, template, template, &privKey.PublicKey, privKey)
	assert.NoError(t, err)

	cert, err := x509.ParseCertificate(derBytes)
	assert.NoError(t, err)

	tufRSAx509Key := trustmanager.CertToKey(cert)
	assert.NoError(t, err)

	// make a data.PublicKey from the generated private key
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	assert.NoError(t, err)
	tufRSAKey := data.NewPublicKey(data.RSAKey, pubBytes)

	// test signing against a service that only recognizes a RSAKey (not
	// RSAx509 key)
	mockCryptoService := &StrictMockCryptoService{MockCryptoService{tufRSAKey}}
	testData := data.Signed{}
	err = Sign(mockCryptoService, &testData, tufRSAx509Key)
	assert.NoError(t, err)

	assert.Len(t, testData.Signatures, 1)
	assert.Equal(t, tufRSAx509Key.ID(), testData.Signatures[0].KeyID)
}
