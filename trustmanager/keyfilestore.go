package trustmanager

import (
	"path/filepath"
	"strings"

	"github.com/endophage/gotuf/data"
	"errors"
	"fmt"
)

const (
	keyExtension = "key"
	aliasExtension = "alias"
)

// KeyStore is a generic interface for private key storage
type KeyStore interface {
	LimitedFileStore

	AddKey(name, alias string, privKey data.PrivateKey) error
	GetKey(name string) (data.PrivateKey, error)
	GetKeyAlias(name string) (string, error)
	ListKeys() []string
	RemoveKey(name string) error
}

// PassphraseRetriever is a callback function that should retrieve a passphrase
// for a given named key. If it should be treated as new passphrase (e.g. with
// confirmation), createNew will be true. Attempts is passed in so that implementers
// decide how many chances to give to a human, for example.
type PassphraseRetriever func(keyId, alias string, createNew bool, attempts int) (passphrase string, giveup bool, err error)

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
func (s *KeyFileStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	return addKey(s, s.PassphraseRetriever, name, alias, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (data.PrivateKey, error) {
	return getKey(s, s.PassphraseRetriever, name)
}

// GetKeyAlias returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKeyAlias(name string) (string, error) {
	return getKeyAlias(s, name)
}

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyFileStore) ListKeys() []string {
	return listKeys(s)
}

// RemoveKey removes the key from the keyfilestore
func (s *KeyFileStore) RemoveKey(name string) error {
	return removeKey(s, name)
}

// NewKeyMemoryStore returns a new KeyMemoryStore which holds keys in memory
func NewKeyMemoryStore(passphraseRetriever PassphraseRetriever) *KeyMemoryStore {
	memStore := NewMemoryFileStore()

	return &KeyMemoryStore{*memStore, passphraseRetriever}
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyMemoryStore) AddKey(name, alias string, privKey data.PrivateKey) error {
	return addKey(s, s.PassphraseRetriever, name, alias, privKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKey(name string) (data.PrivateKey, error) {
	return getKey(s, s.PassphraseRetriever, name)
}

// GetKeyAlias returns the PrivateKey given a KeyID
func (s *KeyMemoryStore) GetKeyAlias(name string) (string, error) {
	return getKeyAlias(s, name)
}


// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyMemoryStore) ListKeys() []string {
	return listKeys(s)
}

// RemoveKey removes the key from the keystore
func (s *KeyMemoryStore) RemoveKey(name string) error {
	return removeKey(s, name)
}


func addKey(s LimitedFileStore, passphraseRetriever PassphraseRetriever, name, alias string, privKey data.PrivateKey) error {
	pemPrivKey, err := KeyToPEM(privKey)
	if err != nil {
		return err
	}

	attempts := 0
	passphrase := ""
	giveup := false
	for {
		passphrase, giveup, err = passphraseRetriever(name, alias, true, attempts)
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
		pemPrivKey, err = EncryptPrivateKey(privKey, passphrase)
		if err != nil {
			return err
		}
	}

	return s.Add(name + "_" + alias, pemPrivKey)
}

func getKeyAlias(s LimitedFileStore, name string) (string, error) {
	files := s.ListFiles(true)

	fmt.Println(name)
	name = name[strings.LastIndexAny(name, "/\\")+1:]
	//name = strings.TrimSpace(strings.TrimSuffix(filepath.Base(name), filepath.Ext(name)))

	fmt.Println(name)

	for _, file := range files {
		fmt.Println(file, " ======= ", name)
		if strings.HasSuffix(file, keyExtension) {
			lastPathSeparator := strings.LastIndexAny(file, "/\\")
			filename := file[lastPathSeparator+1:]
			//filename := strings.TrimSpace(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))

			fmt.Println(filename, " : ", name)

			if strings.HasPrefix(filename, name) {
				fmt.Println("filename:", filename)
				fmt.Println("name:", name)
				aliasPlusDotKey := strings.TrimPrefix(filename, name + "_")
				fmt.Println("aliasPlusDotKey:", aliasPlusDotKey)

				retVal := strings.TrimSuffix(aliasPlusDotKey, "." + keyExtension)
				fmt.Println("retVal:", retVal)

				return retVal, nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("keyId %s has no alias", name))
}

// GetKey returns the PrivateKey given a KeyID
func getKey(s LimitedFileStore, passphraseRetriever PassphraseRetriever, name string) (data.PrivateKey, error) {

	keyAlias, err := getKeyAlias(s, name)
	if err != nil {
		return nil, err
	}

	keyBytes, err := s.Get(name + "_" + keyAlias)
	if err != nil {
		return nil, err
	}

	// See if the key is encrypted. If its encrypted we'll fail to parse the private key
	privKey, err := ParsePEMPrivateKey(keyBytes, "")
	if err != nil {
		// We need to decrypt the key, lets get a passphrase
		attempts := 0
		for {
			passphrase, giveup, err := passphraseRetriever(name, string(keyAlias), false, attempts)
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
		keyID = keyID[:strings.LastIndex(keyID,"_")]
		keyIDList = append(keyIDList, keyID)
	}
	return keyIDList
}

// RemoveKey removes the key from the keyfilestore
func removeKey(s LimitedFileStore, name string) error {
	return s.Remove(name)
}
