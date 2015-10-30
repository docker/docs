// +build pkcs11

package client

import (
	"net/http"
	"path/filepath"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/keystoremanager"
	"github.com/docker/notary/pkg/passphrase"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/tuf/store"
)

// NewNotaryRepository is a helper method that returns a new notary repository.
// It takes the base directory under where all the trust files will be stored
// (usually ~/.docker/trust/).
func NewNotaryRepository(baseDir, gun, baseURL string, rt http.RoundTripper,
	passphraseRetriever passphrase.Retriever) (*NotaryRepository, error) {

	keyStoreManager, err := keystoremanager.NewKeyStoreManager(baseDir, passphraseRetriever)
	if err != nil {
		return nil, err
	}

	yubiKeyStore := api.NewYubiKeyStore()
	cryptoService := cryptoservice.NewCryptoService(gun, yubiKeyStore, keyStoreManager.KeyStore)

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
