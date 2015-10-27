package trustmanager

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/notary/pkg/passphrase"
	"github.com/endophage/gotuf/data"
)

const (
	rootKeysSubdir    = "root_keys"
	nonRootKeysSubdir = "tuf_keys"
)

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	sync.Mutex
	SimpleFileStore
	passphrase.Retriever
	cachedKeys map[string]*cachedKey
}

// KeyMemoryStore manages private keys in memory
type KeyMemoryStore struct {
	sync.Mutex
	MemoryFileStore
	passphrase.Retriever
	cachedKeys map[string]*cachedKey
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string, passphraseRetriever passphrase.Retriever) (*KeyFileStore, error) {
	fileStore, err := NewPrivateSimpleFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}
	cachedKeys := make(map[string]*cachedKey)

	return &KeyFileStore{SimpleFileStore: *fileStore,
		Retriever:  passphraseRetriever,
		cachedKeys: cachedKeys}, nil
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyFileStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	s.Lock()
	defer s.Unlock()
	return addKey(s, s.Retriever, s.cachedKeys, name, alias, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	return getKey(s, s.Retriever, s.cachedKeys, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
func (s *KeyFileStore) ListKeys() map[string]string {
	return listKeys(s)
}

// RemoveKey removes the key from the keyfilestore
func (s *KeyFileStore) RemoveKey(name string) error {
	s.Lock()
	defer s.Unlock()
	return removeKey(s, s.cachedKeys, name)
}

// ExportKey exportes the encrypted bytes from the keystore and writes it to
// dest.
func (s *KeyFileStore) ExportKey(name string) ([]byte, error) {
	keyBytes, _, err := getRawKey(s, name)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// ImportKey imports the private key in the encrypted bytes into the keystore
// with the given key ID and alias.
func (s *KeyFileStore) ImportKey(pemBytes []byte, alias string) error {
	return importKey(s, s.Retriever, s.cachedKeys, alias, pemBytes)
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore(passphraseRetriever passphrase.Retriever) *KeyMemoryStore {
	memStore := NewMemoryFileStore()
	cachedKeys := make(map[string]*cachedKey)

	return &KeyMemoryStore{MemoryFileStore: *memStore,
		Retriever:  passphraseRetriever,
		cachedKeys: cachedKeys}
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyMemoryStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	s.Lock()
	defer s.Unlock()
	return addKey(s, s.Retriever, s.cachedKeys, name, alias, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKey(name string) (data.PrivateKey, string, error) {
	s.Lock()
	defer s.Unlock()
	return getKey(s, s.Retriever, s.cachedKeys, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
func (s *KeyMemoryStore) ListKeys() map[string]string {
	return listKeys(s)
}

// RemoveKey removes the key from the keystore
func (s *KeyMemoryStore) RemoveKey(name string) error {
	s.Lock()
	defer s.Unlock()
	return removeKey(s, s.cachedKeys, name)
}

// ExportKey exportes the encrypted bytes from the keystore and writes it to
// dest.
func (s *KeyMemoryStore) ExportKey(name string) ([]byte, error) {
	keyBytes, _, err := getRawKey(s, name)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// ImportKey imports the private key in the encrypted bytes into the keystore
// with the given key ID and alias.
func (s *KeyMemoryStore) ImportKey(pemBytes []byte, alias string) error {
	return importKey(s, s.Retriever, s.cachedKeys, alias, pemBytes)
}

func addKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, name, alias string, privKey data.PrivateKey) error {

	var (
		chosenPassphrase string
		giveup           bool
		err              error
	)

	for attempts := 0; ; attempts++ {
		chosenPassphrase, giveup, err = passphraseRetriever(name, alias, true, attempts)
		if err != nil {
			continue
		}
		if giveup {
			return ErrAttemptsExceeded{}
		}
		if attempts > 10 {
			return ErrAttemptsExceeded{}
		}
		break
	}

	return encryptAndAddKey(s, chosenPassphrase, cachedKeys, name, alias, privKey)
}

func getKeyAlias(s LimitedFileStore, keyID string) (string, error) {
	files := s.ListFiles()

	name := strings.TrimSpace(strings.TrimSuffix(filepath.Base(keyID), filepath.Ext(keyID)))

	for _, file := range files {
		filename := filepath.Base(file)

		if strings.HasPrefix(filename, name) {
			aliasPlusDotKey := strings.TrimPrefix(filename, name+"_")
			retVal := strings.TrimSuffix(aliasPlusDotKey, "."+keyExtension)
			return retVal, nil
		}
	}

	return "", &ErrKeyNotFound{KeyID: keyID}
}

// GetKey returns the PrivateKey given a KeyID
func getKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, name string) (data.PrivateKey, string, error) {
	cachedKeyEntry, ok := cachedKeys[name]
	if ok {
		return cachedKeyEntry.key, cachedKeyEntry.alias, nil
	}

	keyBytes, keyAlias, err := getRawKey(s, name)
	if err != nil {
		return nil, "", err
	}

	var retErr error
	// See if the key is encrypted. If its encrypted we'll fail to parse the private key
	privKey, err := ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		privKey, _, retErr = getPasswdDecryptBytes(s, passphraseRetriever, keyBytes, name, string(keyAlias))
	}
	if retErr != nil {
		return nil, "", retErr
	}
	cachedKeys[name] = &cachedKey{alias: keyAlias, key: privKey}
	return privKey, keyAlias, nil
}

// ListKeys returns a map of unique PublicKeys present on the KeyFileStore and
// their corresponding aliases.
func listKeys(s LimitedFileStore) map[string]string {
	keyIDMap := make(map[string]string)

	for _, f := range s.ListFiles() {
		if f[:len(rootKeysSubdir)] == rootKeysSubdir {
			f = strings.TrimPrefix(f, rootKeysSubdir+"/")
		} else {
			f = strings.TrimPrefix(f, nonRootKeysSubdir+"/")
		}
		keyIDFull := strings.TrimSpace(strings.TrimSuffix(f, filepath.Ext(f)))
		keyID := keyIDFull[:strings.LastIndex(keyIDFull, "_")]
		keyAlias := keyIDFull[strings.LastIndex(keyIDFull, "_")+1:]
		keyIDMap[keyID] = keyAlias
	}
	return keyIDMap
}

// RemoveKey removes the key from the keyfilestore
func removeKey(s LimitedFileStore, cachedKeys map[string]*cachedKey, name string) error {
	keyAlias, err := getKeyAlias(s, name)
	if err != nil {
		return err
	}

	delete(cachedKeys, name)

	// being in a subdirectory is for backwards compatibliity
	filename := name + "_" + keyAlias
	err = s.Remove(filepath.Join(getSubdir(keyAlias), filename))
	if err != nil {
		return err
	}
	return nil
}

// Assumes 2 subdirectories, 1 containing root keys and 1 containing tuf keys
func getSubdir(alias string) string {
	if alias == "root" {
		return rootKeysSubdir
	}
	return nonRootKeysSubdir
}

// Given a key ID, gets the bytes and alias belonging to that key if the key
// exists
func getRawKey(s LimitedFileStore, name string) ([]byte, string, error) {
	keyAlias, err := getKeyAlias(s, name)
	if err != nil {
		return nil, "", err
	}

	filename := name + "_" + keyAlias
	var keyBytes []byte
	keyBytes, err = s.Get(filepath.Join(getSubdir(keyAlias), filename))
	if err != nil {
		return nil, "", err
	}
	return keyBytes, keyAlias, nil
}

// Get the password to decript the given pem bytes.  Return the password,
// because it is useful for importing
func getPasswdDecryptBytes(s LimitedFileStore, passphraseRetriever passphrase.Retriever, pemBytes []byte, name, alias string) (data.PrivateKey, string, error) {
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
		privKey, err = ParsePEMPrivateKey(pemBytes, passwd)
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

func encryptAndAddKey(s LimitedFileStore, passwd string, cachedKeys map[string]*cachedKey, name, alias string, privKey data.PrivateKey) error {

	var (
		pemPrivKey []byte
		err        error
	)

	if passwd != "" {
		pemPrivKey, err = EncryptPrivateKey(privKey, passwd)
	} else {
		pemPrivKey, err = KeyToPEM(privKey)
	}

	if err != nil {
		return err
	}

	cachedKeys[name] = &cachedKey{alias: alias, key: privKey}
	return s.Add(filepath.Join(getSubdir(alias), name+"_"+alias), pemPrivKey)
}

func importKey(s LimitedFileStore, passphraseRetriever passphrase.Retriever, cachedKeys map[string]*cachedKey, alias string, pemBytes []byte) error {

	privKey, passphrase, err := getPasswdDecryptBytes(s, passphraseRetriever, pemBytes, "imported", alias)

	if err != nil {
		return err
	}

	return encryptAndAddKey(
		s, passphrase, cachedKeys, privKey.ID(), alias, privKey)
}
