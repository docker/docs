package trustmanager

import "github.com/endophage/gotuf/data"

const (
	keyExtension = "key"
)

// KeyFileStore persists and manages private keys on disk
type KeyFileStore struct {
	FileStore
}

// NewKeyFileStore returns a new KeyFileStore creating a private directory to
// hold the keys.
func NewKeyFileStore(baseDir string) (*KeyFileStore, error) {
	fileStore, err := NewPrivateSimpleFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}

	return &KeyFileStore{fileStore}, nil
}

// AddKey stores the contents of a PEM-encoded private key as a PEM block
func (s *KeyFileStore) AddKey(name string, privKey *data.PrivateKey) error {
	pemPrivKey, err := KeyToPEM(privKey)
	if err != nil {
		return err
	}

	return s.Add(name, pemPrivKey)
}

// GetKey returns the PrivateKey given a KeyID
func (s *KeyFileStore) GetKey(name string) (*data.PrivateKey, error) {
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

// AddEncrypted stores the contents of a PEM-encoded private key as an encrypted PEM block
func (s *KeyFileStore) AddEncryptedKey(name string, privKey *data.PrivateKey, passphrase string) error {
	encryptedPrivKey, err := EncryptPrivateKey(privKey, passphrase)
	if err != nil {
		return err
	}

	return s.Add(name, encryptedPrivKey)
}

// GetDecrypted decrypts and returns the PEM Encoded private key given a flename
// and a passphrase
func (s *KeyFileStore) GetDecryptedKey(name string, passphrase string) (*data.PrivateKey, error) {
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

// ListKeys returns a list of unique PublicKeys present on the KeyFileStore.
// There might be symlinks associating Certificate IDs to Public Keys, so this
// method only returns the IDs that aren't symlinks
func (s *KeyFileStore) ListKeys() []string {
	return s.ListFiles(false)
}
