package cryptoservice

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
)

var algoToSigType = map[string]data.SigAlgorithm{
	data.ECDSAKey:   data.ECDSASignature,
	data.ED25519Key: data.EDDSASignature,
	data.RSAKey:     data.RSAPSSSignature,
}

var passphraseRetriever = func(string, string, bool, int) (string, bool, error) { return "", false, nil }

type CryptoServiceTester struct {
	cryptoServiceFactory func() *CryptoService
	role                 string
	keyAlgo              string
}

// asserts that created key exists
func (c CryptoServiceTester) TestCreateAndGetKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	// Test Create
	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Test GetKey
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.NotNil(t, retrievedKey,
		c.errorMsg("Could not find key ID %s", tufKey.ID()))
	assert.Equal(t, tufKey.Public(), retrievedKey.Public(),
		c.errorMsg("retrieved public key didn't match"))

	// Test GetPrivateKey
	retrievedKey, alias, err := cryptoService.GetPrivateKey(tufKey.ID())
	assert.NoError(t, err)
	assert.Equal(t, tufKey.ID(), retrievedKey.ID(),
		c.errorMsg("retrieved private key didn't have the right ID"))
	assert.Equal(t, c.role, alias)
}

// asserts that getting key fails for a non-existent key
func (c CryptoServiceTester) TestGetNonexistentKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	assert.Nil(t, cryptoService.GetKey("boguskeyid"),
		c.errorMsg("non-nil result for bogus keyid"))

	_, _, err := cryptoService.GetPrivateKey("boguskeyid")
	assert.NotNil(t, err)
}

// asserts that signing with a created key creates a valid signature
func (c CryptoServiceTester) TestSignWithKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()
	content := []byte("this is a secret")

	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))

	// Test Sign
	signatures, err := cryptoService.Sign([]string{tufKey.ID()}, content)
	assert.NoError(t, err, c.errorMsg("signing failed"))
	assert.Len(t, signatures, 1, c.errorMsg("wrong number of signatures"))

	verifier, ok := signed.Verifiers[algoToSigType[c.keyAlgo]]
	assert.True(t, ok, c.errorMsg("Unknown verifier for algorithm"))

	err = verifier.Verify(tufKey, signatures[0].Signature, content)
	assert.NoError(t, err,
		c.errorMsg("verification failed for %s key type", c.keyAlgo))
}

// asserts that removing key that exists succeeds
func (c CryptoServiceTester) TestRemoveCreatedKey(t *testing.T) {
	cryptoService := c.cryptoServiceFactory()

	tufKey, err := cryptoService.Create(c.role, c.keyAlgo)
	assert.NoError(t, err, c.errorMsg("error creating key"))
	assert.NotNil(t, cryptoService.GetKey(tufKey.ID()))

	// Test RemoveKey
	err = cryptoService.RemoveKey(tufKey.ID())
	assert.NoError(t, err, c.errorMsg("could not remove key"))
	retrievedKey := cryptoService.GetKey(tufKey.ID())
	assert.Nil(t, retrievedKey, c.errorMsg("remove didn't work"))
}

// Prints out an error message with information about the key algorithm,
// role, and test name. Ideally we could generate different tests given
// data, without having to put for loops in one giant test function, but
// that involves a lot of boilerplate.  So as a compromise, everything will
// still be run in for loops in one giant test function, but we can at
// least provide an error message stating what data/helper test function
// failed.
func (c CryptoServiceTester) errorMsg(message string, args ...interface{}) string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)    // the caller of errorMsg
	f := runtime.FuncForPC(pc[0])
	return fmt.Sprintf("%s (role: %s, keyAlgo: %s): %s", f.Name(), c.role,
		c.keyAlgo, fmt.Sprintf(message, args...))
}

func testCryptoService(t *testing.T, gun string) {
	getCryptoService := func() *CryptoService {
		return NewCryptoService(
			gun, trustmanager.NewKeyMemoryStore(passphraseRetriever))
	}
	roles := []string{
		data.CanonicalRootRole,
		data.CanonicalTargetsRole,
		data.CanonicalSnapshotRole,
		data.CanonicalTimestampRole,
	}

	for _, role := range roles {
		for algo := range algoToSigType {
			cst := CryptoServiceTester{
				cryptoServiceFactory: getCryptoService,
				role:                 role,
				keyAlgo:              algo,
			}
			cst.TestCreateAndGetKey(t)
			cst.TestGetNonexistentKey(t)
			cst.TestSignWithKey(t)
			cst.TestRemoveCreatedKey(t)
		}
	}
}

func TestCryptoServiceWithNonEmptyGUN(t *testing.T) {
	testCryptoService(t, "org/repo")
}

func TestCryptoServiceWithEmptyGUN(t *testing.T) {
	testCryptoService(t, "")
}
