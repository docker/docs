// +build pkcs11

package client

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/docker/notary"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
)

// NewNotaryRepository is a helper method that returns a new notary repository.
// It takes the base directory under where all the trust files will be stored
// (usually ~/.docker/trust/).
func NewNotaryRepository(baseDir, gun, baseURL string, rt http.RoundTripper,
	retriever passphrase.Retriever) (*NotaryRepository, error) {

	keysPath := filepath.Join(baseDir, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, retriever)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key store in directory: %s", keysPath)
	}

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(baseDir, fileKeyStore)
	yubiKeyStore, _ := trustmanager.NewYubiKeyStore(fileKeyStore, retriever)
	var cryptoService signed.CryptoService
	if yubiKeyStore == nil {
		cryptoService = cryptoservice.NewCryptoService(gun, keyStoreManager.KeyStore)
	} else {
		cryptoService = cryptoservice.NewCryptoService(gun, yubiKeyStore, keyStoreManager.KeyStore)
	}

	nRepo := &NotaryRepository{
		gun:             gun,
		baseDir:         baseDir,
		baseURL:         baseURL,
		tufRepoPath:     filepath.Join(baseDir, tufDir, filepath.FromSlash(gun)),
		CryptoService:   cryptoService,
		roundTrip:       rt,
		KeyStoreManager: keyStoreManager,
	}

	fileStore, err := store.NewFilesystemStore(
		nRepo.tufRepoPath,
		"metadata",
		"json",
		"",
	)
	if err != nil {
		return nil, err
	}
	nRepo.fileStore = fileStore

	return nRepo, nil
}
