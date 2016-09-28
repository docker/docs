package enzi

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"path"
	"strings"

	libkv "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/utils"
)

// Users' Public Keys are organized in the folowing directory
// structure:
//     orca/v1/auth/
//         public_keys/
//             by_username/
//                 {username}/
//                     {key_id} -> JSON({key_pem, label, key_id, username})
//                     ...
//             by_key_id/
//                 {key_id} -> JSON({username})
//                 ...
//
// Primary lookup/storage is (username, key_id).
// Secondary index on key_id only, enforce uniqueness.
//
// Prefix for user public keys which are used to authenticate using
// TLS.
const kvStorePublicKeysPrefix = "orca/v1/auth/public_keys"

func userPublicKeysDir(username string) string {
	return path.Join(kvStorePublicKeysPrefix, "by_username", username)
}

func userPubliKeyPath(username, keyID string) string {
	return path.Join(userPublicKeysDir(username), keyID)
}

func publicKeyUsernamePath(keyID string) string {
	return path.Join(kvStorePublicKeysPrefix, "by_key_id", keyID)
}

func generateKeyID(publicKey crypto.PublicKey) (string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("unable to encode public key to DER bytes: %s", err)
	}

	return keyIDFromDer(derBytes), nil
}

func keyIDFromDer(derBytes []byte) string {
	hasher := sha256.New()
	hasher.Write(derBytes)

	return hex.EncodeToString(hasher.Sum(nil))
}

type keyUsernameWrapper struct {
	Username string `json:"username"`
}

// AddUserPublicKey adds the given key to the user's list of public keys. When
// adding or updating a key:
//     1. add to primary storage
//     2. add to secondary, fail if already exists with a different username
//     3. if step 2 failed, remove from primary storage
func (a *Authenticator) AddUserPublicKey(user *auth.Account, label string, publicKey crypto.PublicKey) error {
	// Encode the public key into DER bytes.
	derBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return auth.ErrInvalidPublicKey
	}

	// Generate the key ID as a hash of these bytes.
	keyID := keyIDFromDer(derBytes)

	// Encode to PEM.
	pemKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}))

	accountKey := &auth.AccountKey{
		UserID:    user.ID,
		KeyID:     keyID,
		Label:     label,
		PublicKey: strings.TrimSpace(pemKey),
	}

	return a.storeAccountKey(user.Username, accountKey)
}

// storeAccountKey is used by both SetUserPublicKey and the bulk
// UpdateUserPublicKeys functions.
func (a *Authenticator) storeAccountKey(username string, accountKey *auth.AccountKey) error {
	accountKeyJSON, err := json.Marshal(accountKey)
	if err != nil {
		return fmt.Errorf("unable to encode account key to JSON: %s", err)
	}

	if err := a.kvStore.Put(userPubliKeyPath(username, accountKey.KeyID), accountKeyJSON, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// If there's an error in the following steps, remove the key.
	defer func() {
		if err != nil {
			// Best-effort delete, don't check errors.
			a.kvStore.Delete(userPubliKeyPath(username, accountKey.KeyID))
		}
	}()

	// Need to add secondary lookup. Check existence first.
	exists, err := a.kvStore.Exists(publicKeyUsernamePath(accountKey.KeyID))
	if err != nil {
		return err
	}

	if exists {
		// Check if the key already belongs to another user.
		kvPair, err := a.kvStore.Get(publicKeyUsernamePath(accountKey.KeyID))
		if err != nil {
			return utils.MaybeWrapEtcdClusterErr(err)
		}

		var wrapper keyUsernameWrapper
		if err := json.Unmarshal(kvPair.Value, &wrapper); err != nil {
			return fmt.Errorf("unable to decode user public key from JSON: %s", err)
		}

		if wrapper.Username != username {
			return fmt.Errorf("key with ID %s is already used by another account", accountKey.KeyID)
		}
	}

	// Finally, add the keyID -> username lookup.
	wrapper := keyUsernameWrapper{
		Username: username,
	}

	wrapperJSON, err := json.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("unable to encode username lookup to JSON: %s", err)
	}

	return utils.MaybeWrapEtcdClusterErr(a.kvStore.Put(publicKeyUsernamePath(accountKey.KeyID), wrapperJSON, nil))
}

func (a *Authenticator) deleteUserPublicKey(username, keyID string) error {
	// First, delete the reverse lookup key.
	if err := a.kvStore.Delete(publicKeyUsernamePath(keyID)); err != nil {
		return err
	}

	return a.kvStore.Delete(userPubliKeyPath(username, keyID))
}

func (a *Authenticator) listUserPublicKeys(username string) ([]auth.AccountKey, error) {
	kvPairs, err := a.kvStore.List(userPublicKeysDir(username))
	if err != nil {
		if err == libkv.ErrKeyNotFound {
			return nil, nil
		}

		return nil, err
	}

	accountKeys := make([]auth.AccountKey, len(kvPairs))
	for i, kvPair := range kvPairs {
		if err := json.Unmarshal(kvPair.Value, &accountKeys[i]); err != nil {
			return nil, fmt.Errorf("unable to decode account key from JSON: %s", err)
		}
	}

	return accountKeys, nil
}

func (a *Authenticator) getPublicKey(publicKey crypto.PublicKey) (*auth.AccountKey, error) {
	keyID, err := generateKeyID(publicKey)
	if err != nil {
		return nil, auth.ErrInvalidPublicKey
	}

	// Lookup the username in the KV store.
	kvPair, err := a.kvStore.Get(publicKeyUsernamePath(keyID))
	if err != nil {
		if err == libkv.ErrKeyNotFound {
			return nil, auth.ErrInvalidPublicKey
		}

		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	var wrapper keyUsernameWrapper
	if err := json.Unmarshal(kvPair.Value, &wrapper); err != nil {
		return nil, fmt.Errorf("unable to decode user public key from JSON: %s", err)
	}

	// We have the username for the key now, so we can retrieve the key
	// from the primary location in the KV store.
	kvPair, err = a.kvStore.Get(userPubliKeyPath(wrapper.Username, keyID))
	if err != nil {
		if err == libkv.ErrKeyNotFound {
			return nil, auth.ErrInvalidPublicKey
		}

		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	var accountKey auth.AccountKey
	if err := json.Unmarshal(kvPair.Value, &accountKey); err != nil {
		return nil, fmt.Errorf("unable to decode user public key from JSON: %s", err)
	}

	return &accountKey, nil
}

func parsePublicKeyPEM(publicKeyPEM string) (pubKey crypto.PublicKey, keyID string, err error) {
	pemBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pemBlock == nil {
		return nil, "", auth.ErrInvalidPublicKey
	}

	if pemBlock.Type != "PUBLIC KEY" {
		return nil, "", auth.ErrInvalidPublicKey
	}

	pubKey, err = x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, "", auth.ErrInvalidPublicKey
	}

	return pubKey, keyIDFromDer(pemBlock.Bytes), nil
}

// newAccountKey is a type used for validating new account keys.
type newAccountKey struct {
	keyID        string
	label        string
	publicKey    crypto.PublicKey
	publicKeyPEM string
}

// UpdateUserPublicKeys makes the stored set of public keys for the given user
// equal to the list of keys on the account object.
func (a *Authenticator) updateUserPublicKeys(user *auth.Account) error {
	// Get the current list of keys for this user.
	currentKeys, err := a.listUserPublicKeys(user.Username)
	if err != nil {
		return fmt.Errorf("unable to list current public keys for user: %s", err)
	}

	// Make a set of the current keys by key ID. We can assume that the
	// public keys from the KV store are valid and consistent.
	currentKeysSet := make(map[string]*auth.AccountKey, len(currentKeys))
	for i := range currentKeys {
		currentKey := &currentKeys[i]
		currentKeysSet[currentKey.KeyID] = &currentKeys[i]
	}

	// Detect which keys are actually new and those for which we need to
	// update the label.
	for i := range user.PublicKeys {
		// Get pointer ref so we can modify the value on the user
		// object.
		newKey := &user.PublicKeys[i]

		_, keyID, err := parsePublicKeyPEM(newKey.PublicKey)
		if err != nil {
			// The new key is invalid.
			return err
		}

		if existingKey, ok := currentKeysSet[keyID]; ok {
			// The user already had this key. Remove the "new" key
			// from the currentKeysSet now so that by the end of
			// this loop, whatever is left in this set should be
			// deleted.
			delete(currentKeysSet, keyID)

			if existingKey.Label == newKey.Label {
				// The label did not change so we can skip the
				// update.
				continue
			}
		}

		// Ensure that these fields are set.
		newKey.UserID = user.ID
		newKey.KeyID = keyID

		// Update the key in the kvStore.
		if err := a.storeAccountKey(user.Username, newKey); err != nil {
			return err
		}
	}

	// Now that all of the new/updated keys have been stored, any keys that
	// are still in the currentKeysSet should be deleted.
	for _, keyToDelete := range currentKeysSet {
		if err := a.deleteUserPublicKey(user.Username, keyToDelete.KeyID); err != nil {
			return err
		}
	}

	return nil
}
