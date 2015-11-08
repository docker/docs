// +build pkcs11

package main

import (
	"path/filepath"

	"github.com/docker/notary"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/spf13/cobra"
)

// Build a CryptoService, optionally including a hardware keystore.  Returns
// the CryptoService and whether or not a hardware keystore was included.
func getCryptoService(cmd *cobra.Command, directory string,
	ret passphrase.Retriever, withHardware bool) *cryptoservice.CryptoService {

	keysPath := filepath.Join(directory, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, ret)
	if err != nil {
		fatalf("Failed to create private key store in directory: %s", keysPath)
	}

	ks := []trustmanager.KeyStore{fileKeyStore}

	if withHardware {
		yubiStore, err := api.NewYubiKeyStore(fileKeyStore, ret)
		if err != nil {
			cmd.Println("No YubiKey detected - using local filesystem only.")
		} else {
			// Note that the order is important, since we want to prioritize
			// the yubikey store
			ks = []trustmanager.KeyStore{yubiStore, fileKeyStore}
		}
	}

	return cryptoservice.NewCryptoService("", ks...)
}
