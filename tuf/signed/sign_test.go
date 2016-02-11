package signed

import (
	"crypto/rand"
	"encoding/pem"
	"io"
	"testing"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

const (
	testKeyPEM1 = "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAnKuXZeefa2LmgxaL5NsM\nzKOHNe+x/nL6ik+lDBCTV6OdcwAhHQS+PONGhrChIUVR6Vth3hUCrreLzPO73Oo5\nVSCuRJ53UronENl6lsa5mFKP8StYLvIDITNvkoT3j52BJIjyNUK9UKY9As2TNqDf\nBEPIRp28ev/NViwGOEkBu2UAbwCIdnDXm8JQErCZA0Ydm7PKGgjLbFsFGrVzqXHK\n6pdzJXlhr9yap3UpgQ/iO9JtoEYB2EXsnSrPc9JRjR30bNHHtnVql3fvinXrAEwq\n3xmN4p+R4VGzfdQN+8Kl/IPjqWB535twhFYEG/B7Ze8IwbygBjK3co/KnOPqMUrM\nBI8ztvPiogz+MvXb8WvarZ6TMTh8ifZI96r7zzqyzjR1hJulEy3IsMGvz8XS2J0X\n7sXoaqszEtXdq5ef5zKVxkiyIQZcbPgmpHLq4MgfdryuVVc/RPASoRIXG4lKaTJj\n1ANMFPxDQpHudCLxwCzjCb+sVa20HBRPTnzo8LSZkI6jAgMBAAE=\n-----END PUBLIC KEY-----"
	testKeyID1  = "51324b59d4888faa91219ebbe5a3876bb4efb21f0602ddf363cd4c3996ded3d4"
)

type FailingCryptoService struct {
	testKey data.PrivateKey
}

func (mts *FailingCryptoService) Create(_, _ string) (data.PublicKey, error) {
	return mts.testKey, nil
}

func (mts *FailingCryptoService) ListKeys(role string) []string {
	return []string{mts.testKey.ID()}
}

func (mts *FailingCryptoService) AddKey(role string, key data.PrivateKey) error {
	return nil
}

func (mts *FailingCryptoService) ListAllKeys() map[string]string {
	return map[string]string{
		mts.testKey.ID(): data.CanonicalRootRole,
		mts.testKey.ID(): data.CanonicalTargetsRole,
		mts.testKey.ID(): data.CanonicalSnapshotRole,
		mts.testKey.ID(): data.CanonicalTimestampRole,
	}
}

func (mts *FailingCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == "testID" {
		return data.PublicKeyFromPrivate(mts.testKey)
	}
	return nil
}

func (mts *FailingCryptoService) GetPrivateKey(keyID string) (data.PrivateKey, string, error) {
	if mts.testKey != nil {
		return mts.testKey, "testRole", nil
	}
	return nil, "", trustmanager.ErrKeyNotFound{KeyID: keyID}
}

func (mts *FailingCryptoService) RemoveKey(keyID string) error {
	return nil
}

func (mts *FailingCryptoService) ImportRootKey(r io.Reader) error {
	return nil
}

type MockCryptoService struct {
	testKey data.PrivateKey
}

func (mts *MockCryptoService) Create(_ string, _ string) (data.PublicKey, error) {
	return mts.testKey, nil
}

func (mts *MockCryptoService) AddKey(role string, key data.PrivateKey) error {
	return nil
}

func (mts *MockCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == "testID" {
		return data.PublicKeyFromPrivate(mts.testKey)
	}
	return nil
}

func (mts *MockCryptoService) ListKeys(role string) []string {
	return []string{mts.testKey.ID()}
}

func (mts *MockCryptoService) ListAllKeys() map[string]string {
	return map[string]string{
		mts.testKey.ID(): data.CanonicalRootRole,
		mts.testKey.ID(): data.CanonicalTargetsRole,
		mts.testKey.ID(): data.CanonicalSnapshotRole,
		mts.testKey.ID(): data.CanonicalTimestampRole,
	}
}

func (mts *MockCryptoService) GetPrivateKey(keyID string) (data.PrivateKey, string, error) {
	return mts.testKey, "testRole", nil
}

func (mts *MockCryptoService) RemoveKey(keyID string) error {
	return nil
}

func (mts *MockCryptoService) ImportRootKey(r io.Reader) error {
	return nil
}

var _ CryptoService = &MockCryptoService{}

type StrictMockCryptoService struct {
	MockCryptoService
}

func (mts *StrictMockCryptoService) GetKey(keyID string) data.PublicKey {
	if keyID == mts.testKey.ID() {
		return data.PublicKeyFromPrivate(mts.testKey)
	}
	return nil
}

func (mts *StrictMockCryptoService) ListKeys(role string) []string {
	return []string{mts.testKey.ID()}
}

func (mts *StrictMockCryptoService) ListAllKeys() map[string]string {
	return map[string]string{
		mts.testKey.ID(): data.CanonicalRootRole,
		mts.testKey.ID(): data.CanonicalTargetsRole,
		mts.testKey.ID(): data.CanonicalSnapshotRole,
		mts.testKey.ID(): data.CanonicalTimestampRole,
	}
}

func (mts *StrictMockCryptoService) AddKey(role string, key data.PrivateKey) error {
	return nil
}

func (mts *StrictMockCryptoService) ImportRootKey(r io.Reader) error {
	return nil
}

// Test signing and ensure the expected signature is added
func TestBasicSign(t *testing.T) {
	cs := NewEd25519()
	key, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)
	testData := data.Signed{}

	err = Sign(cs, &testData, key)
	require.NoError(t, err)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != key.ID() {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}
}

// Signing with the same key multiple times should not produce multiple sigs
// with the same key ID
func TestReSign(t *testing.T) {
	cs := NewEd25519()
	key, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)
	testData := data.Signed{}

	Sign(cs, &testData, key)
	Sign(cs, &testData, key)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != key.ID() {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}

}

// Should not remove signatures for valid keys that were not resigned with
func TestMultiSign(t *testing.T) {
	cs := NewEd25519()
	testData := data.Signed{}

	key1, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	Sign(cs, &testData, key1)

	// reinitializing cs means it won't know about key1. We want
	// to attempt to sign passing both key1 and key2, while expecting
	// that the signature for key1 is left intact and the signature
	// for key2 is added
	cs = NewEd25519()
	key2, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	Sign(
		cs,
		&testData,
		key1,
		key2,
	)

	if len(testData.Signatures) != 2 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	keyIDs := map[string]struct{}{key1.ID(): {}, key2.ID(): {}}
	count := 0
	for _, sig := range testData.Signatures {
		count++
		if _, ok := keyIDs[sig.KeyID]; !ok {
			t.Fatalf("Got a signature we didn't expect: %s", sig.KeyID)
		}
	}
	require.Equal(t, 2, count)
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

func TestSignWithX509(t *testing.T) {
	// generate a key becase we need a cert
	privKey, err := trustmanager.GenerateRSAKey(rand.Reader, 1024)
	require.NoError(t, err)

	// make a RSA x509 key
	cert, err := cryptoservice.GenerateTestingCertificate(privKey.CryptoSigner(), "test")
	require.NoError(t, err)

	tufRSAx509Key := trustmanager.CertToKey(cert)
	require.NoError(t, err)

	// test signing against a service that only recognizes a RSAKey (not
	// RSAx509 key)
	mockCryptoService := &StrictMockCryptoService{MockCryptoService{privKey}}
	testData := data.Signed{}
	err = Sign(mockCryptoService, &testData, tufRSAx509Key)
	require.NoError(t, err)

	require.Len(t, testData.Signatures, 1)
	require.Equal(t, tufRSAx509Key.ID(), testData.Signatures[0].KeyID)
}

func TestSignRemovesValidSigByInvalidKey(t *testing.T) {
	cs := NewEd25519()
	testData := data.Signed{}

	key1, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	Sign(cs, &testData, key1)
	require.Len(t, testData.Signatures, 1)
	require.Equal(t, key1.ID(), testData.Signatures[0].KeyID)

	key2, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	// should remove key1 sig even though it's valid. It no longer appears
	// in the list of signing keys for the role
	Sign(
		cs,
		&testData,
		key2,
	)

	require.Len(t, testData.Signatures, 1)
	require.Equal(t, key2.ID(), testData.Signatures[0].KeyID)
}

func TestSignRemovesInvalidSig(t *testing.T) {
	cs := NewEd25519()
	testData := data.Signed{}

	key1, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	Sign(cs, &testData, key1)
	require.Len(t, testData.Signatures, 1)
	require.Equal(t, key1.ID(), testData.Signatures[0].KeyID)

	// we need cs to "forget" key1 so we can't sign with it
	cs = NewEd25519()
	key2, err := cs.Create(data.CanonicalRootRole, data.ED25519Key)
	require.NoError(t, err)

	// modify test data to invalidate key1 sig
	testData.Signed = []byte{0xff}
	// should remove key1 sig because it's out of date
	Sign(
		cs,
		&testData,
		key1,
		key2,
	)

	require.Len(t, testData.Signatures, 1)
	require.Equal(t, key2.ID(), testData.Signatures[0].KeyID)
}
