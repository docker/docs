package keystoremanager

import (
	"path/filepath"

	"github.com/docker/notary/trustmanager"
)

// KeyStoreManager is an abstraction around the root and non-root key stores
type KeyStoreManager struct {
	rootKeyStore    *trustmanager.KeyFileStore
	nonRootKeyStore *trustmanager.KeyFileStore
}

const (
	privDir        = "private"
	rootKeysSubdir = "root_keys"
)

// NewKeyStoreManager returns an initialized KeyStoreManager, or an error
// if it fails to create the KeyFileStores
func NewKeyStoreManager(baseDir string) (*KeyStoreManager, error) {
	nonRootKeyStore, err := trustmanager.NewKeyFileStore(filepath.Join(baseDir, privDir))
	if err != nil {
		return nil, err
	}

	// Load the keystore that will hold all of our encrypted Root Private Keys
	rootKeysPath := filepath.Join(baseDir, privDir, rootKeysSubdir)
	rootKeyStore, err := trustmanager.NewKeyFileStore(rootKeysPath)
	if err != nil {
		return nil, err
	}

	return &KeyStoreManager{
		rootKeyStore:    rootKeyStore,
		nonRootKeyStore: nonRootKeyStore,
	}, nil
}

// RootKeyStore returns the root key store being managed by this
// KeyStoreManager
func (km *KeyStoreManager) RootKeyStore() *trustmanager.KeyFileStore {
	return km.rootKeyStore
}

// NonRootKeyStore returns the non-root key store being managed by this
// KeyStoreManager
func (km *KeyStoreManager) NonRootKeyStore() *trustmanager.KeyFileStore {
	return km.nonRootKeyStore
}
