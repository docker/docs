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

	AddKey(name string, privKey data.PrivateKey) error
	GetKey(name string) (data.PrivateKey, error)
	ListKeys() []string
	RemoveKey(name string) error
}

// PassphraseRetriever is a callback function that should retrieve a passphrase
// for a given named key. If it should be treated as new passphrase (e.g. with
// confirmation), createNew will be true. Attempts is passed in so that implementers
// decide how many chances to give to a human, for example.
type PassphraseRetriever func(keyName string, createNew bool, attempts int) (passphrase string, giveup bool, err error)

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	SimpleFileStore
	PassphraseRetriever
}

// KeyMemoryStore manages private keys in memory
type KeyMemoryStore struct {
	MemoryFileStore
	PassphraseRetriever
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string, passphraseRetriever PassphraseRetriever) (*KeyFileStore, error) {
	fileStore, err := NewPrivateSimpleFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}

	return &KeyFileStore{*fileStore, passphraseRetriever}, nil
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyFileStore) AddKey(name string, privKey data.PrivateKey) error {
	return addKey(s, s.PassphraseRetriever, name, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (data.PrivateKey, error) {
	return getKey(s, s.PassphraseRetriever, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyFileStore) ListKeys() []string {
	return listKeys(s)
}

// RemoveKey removes the key from the keyfilestore
func (s *KeyFileStore) RemoveKey(name string) error {
	return remove(s, name)
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore(passphraseRetriever PassphraseRetriever) *KeyMemoryStore {
	memStore := NewMemoryFileStore()

	return &KeyMemoryStore{*memStore, passphraseRetriever}
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyMemoryStore) AddKey(name string, privKey data.PrivateKey) error {
	return addKey(s, s.PassphraseRetriever, name, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKey(name string) (data.PrivateKey, error) {
	return getKey(s, s.PassphraseRetriever, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyMemoryStore) ListKeys() []string {
	return listKeys(s)
}

// RemoveKey removes the key from the keystore
func (s *KeyMemoryStore) RemoveKey(name string) error {
	return remove(s, name)
}


func addKey(s LimitedFileStore, passphraseRetriever PassphraseRetriever, name string, privKey data.PrivateKey) error {
	pemPrivKey, err := KeyToPEM(privKey)
	if err != nil {
		return err
	}

	attempts := 0
	passphrase := ""
	giveup := false
	for (true) {
		passphrase, giveup, err = passphraseRetriever(name, true, attempts)
		if err != nil {
			attempts++
			continue
		} else if giveup {
			return err
		} else {
			break
		}
	}

	if passphrase != "" {
		encryptedPrivKey, err := EncryptPrivateKey(privKey, passphrase)
		if err != nil {
			return err
		}

		return s.Add(name, encryptedPrivKey)
	}

	return s.Add(name, pemPrivKey)
}

// GetKey returns the PrivateKey given a KeyID
func getKey(s LimitedFileStore, passphraseRetriever PassphraseRetriever, name string) (data.PrivateKey, error) {
	keyBytes, err := s.Get(name)
	if err != nil {
		return nil, err
	}

	// See if the key is encrypted. If its encrypted we'll fail to parse the private key
	privKey, err := ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		// We need to decrypt the key, lets get a passphrase
		attempts := 0
		for (true) {
			passphrase, giveup, err := passphraseRetriever(name, false, attempts)
			// Check if the passphrase retriever got an error or if it is telling us to give up
			if giveup || err != nil {
				return nil, err
			}

			// Try to convert PEM encoded bytes back to a PrivateKey using the passphrase
			privKey, err = ParsePEMPrivateKey(keyBytes, passphrase)
			if err == nil {
				// We managed to parse the PrivateKey. We've succeeded!
				break
			}
			attempts++
		}
	}
	return privKey, nil
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func listKeys(s LimitedFileStore) []string {
	var keyIDList []string
	for _, f := range s.ListFiles(false) {
		keyID := strings.TrimSpace(strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)))
		keyIDList = append(keyIDList, keyID)
	}
	return keyIDList
}

// RemoveKey removes the key from the keyfilestore
func remove(s LimitedFileStore, name string) error {
	return s.Remove(name)
}
