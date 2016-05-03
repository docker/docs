package interfaces

import (
	"crypto/rand"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/stretchr/testify/require"
)

// These are tests that can be used to test a cryptoservice

// EmptyCryptoServiceInterfaceNormalBehaviorTests tests expected behavior for
// an empty signed.CryptoService:
// 1.  Getting the public key of a key that doesn't exist should fail
// 2.  Listing an empty cryptoservice returns no keys
// 3.  Removing a non-existent key succeeds (no-op)
func EmptyCryptoServiceInterfaceNormalBehaviorTests(t *testing.T, empty signed.CryptoService) {
	for _, role := range append(data.BaseRoles, "targets/delegation", "invalid") {
		keys := empty.ListKeys(role)
		require.Len(t, keys, 0)
	}
	keys := empty.ListAllKeys()
	require.Len(t, keys, 0)

	require.NoError(t, empty.RemoveKey("nonexistent"))

	require.Nil(t, empty.GetKey("nonexistent"))

	k, _, err := empty.GetPrivateKey("nonexistent")
	require.Error(t, err)
	require.Nil(t, k)
}

// CreateKeyCryptoServiceInterfaceNormalBehaviorTests tests expected behavior for
// creating keys in a signed.CryptoService and other read operations on the
// crypto service after keys are present
// 1.  Creating a key succeeds and returns a non-nil public key
// 2.  Getting the key should return the same key, without error
// 3.  Listing returns the correct number of keys and right roles
// 4.  Removing the key succeeds
func CreateKeyCryptoServiceInterfaceNormalBehaviorTests(t *testing.T, cs signed.CryptoService, algo string) {
	expectedRolesToKeys := make(map[string]string)
	for i := 0; i < 2; i++ {
		role := data.BaseRoles[i+1]
		createdPubKey, err := cs.Create(role, "docker.io/notary", algo)
		require.NoError(t, err)
		require.NotNil(t, createdPubKey)
		expectedRolesToKeys[role] = createdPubKey.ID()
	}

	testReadNonEmptyRepo(t, cs, expectedRolesToKeys)
}

// AddKeyCryptoServiceInterfaceNormalBehaviorTests tests expected behavior for
// adding keys in a signed.CryptoService and other read operations on the
// crypto service after keys are present
// 1.  Adding a key succeeds
// 3.  Getting the key should return the same key, without error
// 4.  Listing returns the correct number of keys and right roles
// 5.  Removing the key succeeds
func AddKeyCryptoServiceInterfaceNormalBehaviorTests(t *testing.T, cs signed.CryptoService, algo string) {
	expectedRolesToKeys := make(map[string]string)
	for i := 0; i < 2; i++ {
		var (
			addedPrivKey data.PrivateKey
			err          error
		)
		role := data.BaseRoles[i+1]
		switch algo {
		case data.RSAKey:
			addedPrivKey, err = trustmanager.GenerateRSAKey(rand.Reader, 2048)
		case data.ECDSAKey:
			addedPrivKey, err = trustmanager.GenerateECDSAKey(rand.Reader)
		case data.ED25519Key:
			addedPrivKey, err = trustmanager.GenerateED25519Key(rand.Reader)
		default:
			require.FailNow(t, "invalid algorithm %s", algo)
		}
		require.NoError(t, err)
		require.NotNil(t, addedPrivKey)
		require.NoError(t, cs.AddKey(role, "docker.io/notary", addedPrivKey))
		expectedRolesToKeys[role] = addedPrivKey.ID()
	}

	testReadNonEmptyRepo(t, cs, expectedRolesToKeys)
}

func testReadNonEmptyRepo(t *testing.T, cs signed.CryptoService, expectedRolesToKeys map[string]string) {
	for _, role := range append(data.BaseRoles, "targets/delegation", "invalid") {
		keys := cs.ListKeys(role)

		if keyID, ok := expectedRolesToKeys[role]; ok {
			require.Len(t, keys, 1)
			require.Equal(t, keyID, keys[0])
		} else {
			require.Len(t, keys, 0)
		}
	}

	keys := cs.ListAllKeys()
	require.Len(t, keys, len(expectedRolesToKeys))
	for role, keyID := range expectedRolesToKeys {
		require.Equal(t, role, keys[keyID])

		pubKey := cs.GetKey(keyID)
		require.NotNil(t, pubKey)
		require.Equal(t, keyID, pubKey.ID())

		privKey, gotRole, err := cs.GetPrivateKey(keyID)
		require.NoError(t, err)
		require.NotNil(t, privKey)
		require.Equal(t, keyID, privKey.ID())
		require.Equal(t, role, gotRole)

		require.NoError(t, cs.RemoveKey(keyID))
		require.Nil(t, cs.GetKey(keyID))
	}
}
