// +build !pkcs11

package main

import (
	"path/filepath"

	"github.com/docker/notary"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/spf13/cobra"
)

func getCryptoService(cmd *cobra.Command, directory string,
	ret passphrase.Retriever, _ bool) *cryptoservice.CryptoService {

	keysPath := filepath.Join(directory, notary.PrivDir)
	fileKeyStore, err := trustmanager.NewKeyFileStore(keysPath, ret)
	if err != nil {
		fatalf("Failed to create private key store in directory: %s", keysPath)
	}
	return cryptoservice.NewCryptoService("", fileKeyStore)
}
