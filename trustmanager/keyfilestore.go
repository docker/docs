package trustmanager

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
	fileStore, err := NewFileStore(baseDir, keyExtension)
	if err != nil {
		return nil, err
	}

	return &KeyFileStore{fileStore}, nil
}

// AddEncrypted stores the contents of a PEM-encoded private key as an encrypted PEM block
func (s *KeyFileStore) AddEncrypted(fileName string, pemKey []byte, passphrase string) error {

	privKey, err := ParsePEMPrivateKey(pemKey)
	if err != nil {
		return err
	}

	encryptedKey, err := EncryptPrivateKey(privKey, passphrase)
	if err != nil {
		return err
	}

	return s.Add(fileName, encryptedKey)
}

// GetDecrypted decrypts and returns the PEM Encoded private key given a flename
// and a passphrase
func (s *KeyFileStore) GetDecrypted(fileName string, passphrase string) ([]byte, error) {
	keyBytes, err := s.Get(fileName)
	if err != nil {
		return nil, err
	}

	// Gets an unencrypted PrivateKey.
	privKey, err := ParsePEMEncryptedPrivateKey(keyBytes, passphrase)
	if err != nil {
		return nil, err
	}

	return KeyToPEM(privKey)
}

func (s *KeyFileStore) Link(src, dst string) error {
	return s.FileStore.Link(src, dst)
}
