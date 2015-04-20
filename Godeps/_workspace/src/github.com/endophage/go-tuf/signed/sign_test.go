package signed

import (
	"testing"

	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/keys"
)

type MockTrustService struct {
	testKey keys.PublicKey
}

func (mts *MockTrustService) Sign(keyIDs []string, _ []byte) ([]data.Signature, error) {
	sigs := make([]data.Signature, 0, len(keyIDs))
	for _, keyID := range keyIDs {
		sigs = append(sigs, data.Signature{KeyID: keyID})
	}
	return sigs, nil
}

func (mts *MockTrustService) Create() (*keys.PublicKey, error) {
	return &mts.testKey, nil
}

func (mts *MockTrustService) PublicKeys(keyIDs ...string) (map[string]*keys.PublicKey, error) {
	keys := map[string]*keys.PublicKey{"testID": &mts.testKey}
	return keys, nil
}

var _ TrustService = &MockTrustService{}

// Test signing and ensure the expected signature is added
func TestBasicSign(t *testing.T) {
	signer := Signer{&MockTrustService{
		testKey: keys.PublicKey{ID: "testID"},
	}}
	key, err := signer.Create()
	if err != nil {
		t.Fatal(err)
	}
	testData := data.Signed{}

	signer.Sign(&testData, key)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != "testID" {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}

}

// Test signing with the same key multiple times only registers a single signature
// for the key (N.B. MockTrustService.Sign will still be called again, but Signer.Sign
// should be cleaning previous signatures by the KeyID when asked to sign again)
func TestReSign(t *testing.T) {
	signer := Signer{&MockTrustService{
		testKey: keys.PublicKey{},
	}}
	key := keys.PublicKey{ID: "testID"}
	testData := data.Signed{}

	signer.Sign(&testData, &key)
	signer.Sign(&testData, &key)

	if len(testData.Signatures) != 1 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	if testData.Signatures[0].KeyID != "testID" {
		t.Fatalf("Wrong signature ID returned: %s", testData.Signatures[0].KeyID)
	}

}

func TestMultiSign(t *testing.T) {
	signer := Signer{&MockTrustService{}}
	key := keys.PublicKey{ID: "testID1"}
	testData := data.Signed{}

	signer.Sign(&testData, &key)

	key = keys.PublicKey{ID: "testID2"}
	signer.Sign(&testData, &key)

	if len(testData.Signatures) != 2 {
		t.Fatalf("Incorrect number of signatures: %d", len(testData.Signatures))
	}

	keyIDs := map[string]struct{}{"testID1": struct{}{}, "testID2": struct{}{}}
	for _, sig := range testData.Signatures {
		if _, ok := keyIDs[sig.KeyID]; !ok {
			t.Fatalf("Got a signature we didn't expect: %s", sig.KeyID)
		}
	}

}

func TestCreate(t *testing.T) {
	signer := Signer{&MockTrustService{
		testKey: keys.PublicKey{ID: "testID"},
	}}

	key, err := signer.Create()

	if err != nil {
		t.Fatal(err)
	}
	if key.ID != "testID" {
		t.Fatalf("Expected key ID not found: %s", key.ID)
	}
}
