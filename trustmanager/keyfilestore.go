package trustmanager

import (
	"path/filepath"
	"strings"

	"github.com/endophage/gotuf/data"
)

const (
	keyExtension = "key"
)

// KeyStore is a generic interface for private key storage
type KeyStore interface {
	LimitedFileStore

	AddKey(name string, privKey *data.PrivateKey) error
	GetKey(name string) (*data.PrivateKey, error)
	AddEncryptedKey(name string, privKey *data.PrivateKey, passphrase string) error
	GetDecryptedKey(name string, passphrase string) (*data.PrivateKey, error)
	ListKeys() []string
}

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	SimpleFileStore
}

// KeyMemoryStore manages private keys in memory
type KeyMemoryStore struct {
	MemoryFileStore
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string) (*KeyFileStore, error) {
	fileStore, err := NewPrivateSimpleFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}

	return &KeyFileStore{*fileStore}, nil
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyFileStore) AddKey(name string, privKey *data.PrivateKey) error {
	return addKey(s, name, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (*data.PrivateKey, error) {
	return getKey(s, name)
}

// AddEncryptedKey stores the contents of a PEM-encoded private key as an encrypted PEM block
func (s *KeyFileStore) AddEncryptedKey(name string, privKey *data.PrivateKey, passphrase string) error {
	return addEncryptedKey(s, name, privKey, passphrase)
}

// GetDecryptedKey decrypts and returns the PEM Encoded private key given a flename
// and a passphrase
func (s *KeyFileStore) GetDecryptedKey(name string, passphrase string) (*data.PrivateKey, error) {
	return getDecryptedKey(s, name, passphrase)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyFileStore) ListKeys() []string {
	return listKeys(s)
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore() *KeyMemoryStore {
	memStore := NewMemoryFileStore()

	return &KeyMemoryStore{*memStore}
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyMemoryStore) AddKey(name string, privKey *data.PrivateKey) error {
	return addKey(s, name, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKey(name string) (*data.PrivateKey, error) {
	return getKey(s, name)
}

// AddEncryptedKey stores the contents of a PEM-encoded private key as an encrypted PEM block
func (s *KeyMemoryStore) AddEncryptedKey(name string, privKey *data.PrivateKey, passphrase string) error {
	return addEncryptedKey(s, name, privKey, passphrase)
}

// GetDecryptedKey decrypts and returns the PEM Encoded private key given a flename
// and a passphrase
func (s *KeyMemoryStore) GetDecryptedKey(name string, passphrase string) (*data.PrivateKey, error) {
	return getDecryptedKey(s, name, passphrase)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyMemoryStore) ListKeys() []string {
	return listKeys(s)
}

func addKey(s LimitedFileStore, name string, privKey *data.PrivateKey) error {
	pemPrivKey, err := KeyToPEM(privKey)
	if err != nil {
		return err
	}

	return s.Add(name, pemPrivKey)
}

func getKey(s LimitedFileStore, name string) (*data.PrivateKey, error) {
	keyBytes, err := s.Get(name)
	if err != nil {
		return nil, err
	}

	// Convert PEM encoded bytes back to a PrivateKey
	privKey, err := ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func addEncryptedKey(s LimitedFileStore, name string, privKey *data.PrivateKey, passphrase string) error {
	encryptedPrivKey, err := EncryptPrivateKey(privKey, passphrase)
	if err != nil {
		return err
	}

	return s.Add(name, encryptedPrivKey)
}

func getDecryptedKey(s LimitedFileStore, name string, passphrase string) (*data.PrivateKey, error) {
	keyBytes, err := s.Get(name)
	if err != nil {
		return nil, err
	}

	// Gets an unencrypted PrivateKey.
	privKey, err := ParsePEMPrivateKey(keyBytes, passphrase)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func listKeys(s LimitedFileStore) []string {
	var keyIDList []string
	for _, f := range s.ListFiles(false) {
		keyID := strings.TrimSpace(strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)))
		keyIDList = append(keyIDList, keyID)
	}
	return keyIDList
}
