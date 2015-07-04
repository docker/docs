package trustmanager

import (
	"crypto"
	"errors"
	"fmt"
)

const (
	keyExtension = "key"
)

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	fingerprintMap map[string]string
	fileStore      FileStore
}

// NewKeyFileStore returns a new KeyFileStore.
func NewKeyFileStore(directory string) (*KeyFileStore, error) {
	fileStore, err := NewPrivateFileStore(directory, keyExtension)
	if err != nil {
		return nil, err
	}

	return &KeyFileStore{
		fingerprintMap: make(map[string]string),
		fileStore:      fileStore,
	}, nil
}

// Add stores both the PrivateKey bytes in a file
func (s *KeyFileStore) Add(fileName string, privKey crypto.PrivateKey) error {
	if privKey == nil {
		return errors.New("adding nil key to keyFileStore")
	}

	pemKey, err := KeyToPEM(privKey)
	if err != nil {
		return err
	}

	return s.fileStore.Add(fileName, pemKey)
}

// Get returns a PrivateKey given a filename
func (s *KeyFileStore) Get(fileName string) (crypto.PrivateKey, error) {
	keyBytes, err := s.fileStore.GetData(fileName)
	if err != nil {
		return nil, errors.New("Could not retrieve private key material")
	}

	return ParseRawPrivateKey(keyBytes)
}

// AddEncrypted stores the contents of the private key as an encrypted PEM block
func (s *KeyFileStore) AddEncrypted(fileName string, privKey crypto.PrivateKey, passphrase string) error {
	if privKey == nil {
		return errors.New("adding nil key to keyFileStore")
	}

	encryptedKey, err := KeyToEncryptedPEM(privKey, passphrase)
	if err != nil {
		return err
	}

	fmt.Println(string(encryptedKey))
	return s.fileStore.Add(fileName, encryptedKey)
}

// GetDecrypted decrypts and returns the private key
func (s *KeyFileStore) GetDecrypted(fileName string, passphrase string) (crypto.PrivateKey, error) {
	keyBytes, err := s.fileStore.GetData(fileName)
	if err != nil {
		return nil, errors.New("could not retrieve private key material")
	}

	return ParseRawEncryptedPrivateKey(keyBytes, passphrase)
}

// Remove removes a key from a store
func (s *KeyFileStore) Remove(fileName string) error {
	return s.fileStore.Remove(fileName)
}

// RemoveAll removes all the keys under a directory
func (s *KeyFileStore) RemoveAll(directoryName string) error {
	return s.fileStore.RemoveDir(directoryName)
}

// List returns a list of all the keys the store is currently managing
func (s *KeyFileStore) ListAll() []string {
	return s.fileStore.ListAll()
}

// List returns a list of all the keys the store is currently managing
func (s *KeyFileStore) ListDir(directoryName string) []string {
	return s.fileStore.ListDir(directoryName)
}
