package trustmanager

import (
	"encoding/pem"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary"
	store "github.com/docker/notary/storage"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/utils"
)

type keyInfoMap map[string]KeyInfo

// KeyInfo stores the role, path, and gun for a corresponding private key ID
// It is assumed that each private key ID is unique
type KeyInfo struct {
	Gun  string
	Role string
}

// GenericKeyStore is a wrapper for Storage instances that provides
// translation between the []byte form and Public/PrivateKey objects
type GenericKeyStore struct {
	store Storage
	sync.Mutex
	notary.PassRetriever
	cachedKeys map[string]*cachedKey
	keyInfoMap
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string, p notary.PassRetriever) (*GenericKeyStore, error) {
	fileStore, err := store.NewPrivateKeyFileStorage(baseDir, notary.KeyExtension)
	if err != nil {
		return nil, err
	}
	return NewGenericKeyStore(fileStore, p), nil
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore(p notary.PassRetriever) *GenericKeyStore {
	memStore := store.NewMemoryStore(nil)
	return NewGenericKeyStore(memStore, p)
}

// NewGenericKeyStore creates a GenericKeyStore wrapping the provided
// Storage instance, using the PassRetriever to enc/decrypt keys
func NewGenericKeyStore(s Storage, p notary.PassRetriever) *GenericKeyStore {
	ks := GenericKeyStore{
		store:         s,
		PassRetriever: p,
		cachedKeys:    make(map[string]*cachedKey),
		keyInfoMap:    make(keyInfoMap),
	}
	ks.loadKeyInfo()
	return &ks
}

func generateKeyInfoMap(s Storage) map[string]KeyInfo {
	keyInfoMap := make(map[string]KeyInfo)
	for _, keyPath := range s.ListFiles() {
		d, err := s.Get(keyPath)
		if err != nil {
			logrus.Error(err)
			continue
		}
		keyID, keyInfo, err := KeyInfoFromPEM(d, keyPath)
		if err != nil {
			logrus.Error(err)
			continue
		}
		keyInfoMap[keyID] = keyInfo
	}
	return keyInfoMap
}

// Attempts to infer the keyID, role, and GUN from the specified key path.
// Note that non-root roles can only be inferred if this is a legacy style filename: KEYID_ROLE.key
func inferKeyInfoFromKeyPath(keyPath string) (string, string, string) {
	var keyID, role, gun string
	keyID = filepath.Base(keyPath)
	underscoreIndex := strings.LastIndex(keyID, "_")

	// This is the legacy KEYID_ROLE filename
	// The keyID is the first part of the keyname
	// The keyRole is the second part of the keyname
	// in a key named abcde_root, abcde is the keyID and root is the KeyAlias
	if underscoreIndex != -1 {
		role = keyID[underscoreIndex+1:]
		keyID = keyID[:underscoreIndex]
	}

	if filepath.HasPrefix(keyPath, notary.RootKeysSubdir+"/") {
		return keyID, data.CanonicalRootRole, ""
	}

	keyPath = strings.TrimPrefix(keyPath, notary.NonRootKeysSubdir+"/")
	gun = getGunFromFullID(keyPath)
	return keyID, role, gun
}

func getGunFromFullID(fullKeyID string) string {
	keyGun := filepath.Dir(fullKeyID)
	// If the gun is empty, Dir will return .
	if keyGun == "." {
		keyGun = ""
	}
	return keyGun
}

func (s *GenericKeyStore) loadKeyInfo() {
	s.keyInfoMap = generateKeyInfoMap(s.store)
}

// GetKeyInfo returns the corresponding gun and role key info for a keyID
func (s *GenericKeyStore) GetKeyInfo(keyID string) (KeyInfo, error) {
	if info, ok := s.keyInfoMap[keyID]; ok {
		return info, nil
	}
	return KeyInfo{}, fmt.Errorf("Could not find info for keyID %s", keyID)
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *GenericKeyStore) AddKey(keyInfo KeyInfo, privKey data.PrivateKey) error {
	var (
		chosenPassphrase string
		giveup           bool
		err              error
		pemPrivKey       []byte
	)
	s.Lock()
	defer s.Unlock()
	if keyInfo.Role == data.CanonicalRootRole || data.IsDelegation(keyInfo.Role) || !data.ValidRole(keyInfo.Role) {
		keyInfo.Gun = ""
	}
	name := filepath.Join(keyInfo.Gun, privKey.ID())
	for attempts := 0; ; attempts++ {
		chosenPassphrase, giveup, err = s.PassRetriever(name, keyInfo.Role, true, attempts)
		if err != nil {
			continue
		}
		if giveup || attempts > 10 {
			return ErrAttemptsExceeded{}
		}
		break
	}

	if chosenPassphrase != "" {
		pemPrivKey, err = utils.EncryptPrivateKey(privKey, keyInfo.Role, chosenPassphrase)
	} else {
		pemPrivKey, err = utils.KeyToPEM(privKey, keyInfo.Role)
	}

	if err != nil {
		return err
	}

	s.cachedKeys[name] = &cachedKey{alias: keyInfo.Role, key: privKey}
	err = s.store.Set(filepath.Join(getSubdir(keyInfo.Role), name), pemPrivKey)
	if err != nil {
		return err
	}
	s.keyInfoMap[privKey.ID()] = keyInfo
	return nil
}

// GetKey returns the PrivateKey given a KeyID
func (s *GenericKeyStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[name]; ok {
		name = filepath.Join(keyInfo.Gun, name)
	}

	cachedKeyEntry, ok := s.cachedKeys[name]
	if ok {
		return cachedKeyEntry.key, cachedKeyEntry.alias, nil
	}

	keyAlias, legacy, err := getKeyRole(s.store, name)
	if err != nil {
		return nil, "", err
	}

	if legacy {
		name = name + "_" + keyAlias
	}

	keyBytes, err := s.store.Get(filepath.Join(getSubdir(keyAlias), name))
	if err != nil {
		return nil, "", err
	}

	// See if the key is encrypted. If its encrypted we'll fail to parse the private key
	privKey, err := utils.ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		privKey, _, err = GetPasswdDecryptBytes(s.PassRetriever, keyBytes, name, string(keyAlias))
		if err != nil {
			return nil, "", err
		}
	}
	s.cachedKeys[name] = &cachedKey{alias: keyAlias, key: privKey}
	return privKey, keyAlias, nil
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore, by returning a copy of the keyInfoMap
func (s *GenericKeyStore) ListKeys() map[string]KeyInfo {
	return copyKeyInfoMap(s.keyInfoMap)
}

// RemoveKey removes the key from the keyfilestore
func (s *GenericKeyStore) RemoveKey(keyID string) error {
	s.Lock()
	defer s.Unlock()
	// If this is a bare key ID without the gun, prepend the gun so the filestore lookup succeeds
	if keyInfo, ok := s.keyInfoMap[keyID]; ok {
		keyID = filepath.Join(keyInfo.Gun, keyID)
	}

	role, legacy, err := getKeyRole(s.store, keyID)
	if err != nil {
		return err
	}

	delete(s.cachedKeys, keyID)

	name := keyID
	if legacy {
		name = keyID + "_" + role
	}

	// being in a subdirectory is for backwards compatibliity
	err = s.store.Remove(filepath.Join(getSubdir(role), name))
	if err != nil {
		return err
	}

	// Remove this key from our keyInfo map if we removed from our filesystem
	delete(s.keyInfoMap, filepath.Base(keyID))
	return nil
}

// Name returns a user friendly name for the location this store
// keeps its data
func (s *GenericKeyStore) Name() string {
	return s.store.Location()
}

// copyKeyInfoMap returns a deep copy of the passed-in keyInfoMap
func copyKeyInfoMap(keyInfoMap map[string]KeyInfo) map[string]KeyInfo {
	copyMap := make(map[string]KeyInfo)
	for keyID, keyInfo := range keyInfoMap {
		copyMap[keyID] = KeyInfo{Role: keyInfo.Role, Gun: keyInfo.Gun}
	}
	return copyMap
}

// KeyInfoFromPEM attempts to get a keyID and KeyInfo from the filename and PEM bytes of a key
func KeyInfoFromPEM(pemBytes []byte, filename string) (string, KeyInfo, error) {
	keyID, role, gun := inferKeyInfoFromKeyPath(filename)
	if role == "" {
		block, _ := pem.Decode(pemBytes)
		if block == nil {
			return "", KeyInfo{}, fmt.Errorf("could not decode PEM block for key %s", filename)
		}
		if keyRole, ok := block.Headers["role"]; ok {
			role = keyRole
		}
	}
	return keyID, KeyInfo{Gun: gun, Role: role}, nil
}

// getKeyRole finds the role for the given keyID. It attempts to look
// both in the newer format PEM headers, and also in the legacy filename
// format. It returns: the role, whether it was found in the legacy format
// (true == legacy), and an error
func getKeyRole(s Storage, keyID string) (string, bool, error) {
	name := strings.TrimSpace(strings.TrimSuffix(filepath.Base(keyID), filepath.Ext(keyID)))

	for _, file := range s.ListFiles() {
		filename := filepath.Base(file)

		if strings.HasPrefix(filename, name) {
			d, err := s.Get(file)
			if err != nil {
				return "", false, err
			}
			block, _ := pem.Decode(d)
			if block != nil {
				if role, ok := block.Headers["role"]; ok {
					return role, false, nil
				}
			}

			role := strings.TrimPrefix(filename, name+"_")
			return role, true, nil
		}
	}

	return "", false, ErrKeyNotFound{KeyID: keyID}
}

// Assumes 2 subdirectories, 1 containing root keys and 1 containing TUF keys
func getSubdir(alias string) string {
	if alias == data.CanonicalRootRole {
		return notary.RootKeysSubdir
	}
	return notary.NonRootKeysSubdir
}

// GetPasswdDecryptBytes gets the password to decrypt the given pem bytes.
// Returns the password and private key
func GetPasswdDecryptBytes(passphraseRetriever notary.PassRetriever, pemBytes []byte, name, alias string) (data.PrivateKey, string, error) {
	var (
		passwd  string
		retErr  error
		privKey data.PrivateKey
	)
	for attempts := 0; ; attempts++ {
		var (
			giveup bool
			err    error
		)
		passwd, giveup, err = passphraseRetriever(name, alias, false, attempts)
		// Check if the passphrase retriever got an error or if it is telling us to give up
		if giveup || err != nil {
			return nil, "", ErrPasswordInvalid{}
		}
		if attempts > 10 {
			return nil, "", ErrAttemptsExceeded{}
		}

		// Try to convert PEM encoded bytes back to a PrivateKey using the passphrase
		privKey, err = utils.ParsePEMPrivateKey(pemBytes, passwd)
		if err != nil {
			retErr = ErrPasswordInvalid{}
		} else {
			// We managed to parse the PrivateKey. We've succeeded!
			retErr = nil
			break
		}
	}
	if retErr != nil {
		return nil, "", retErr
	}
	return privKey, passwd, nil
}
