// +build pkcs11

package main

import (
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
)

func getYubiKeyStore(fileKeyStore trustmanager.KeyStore, ret passphrase.Retriever) (trustmanager.KeyStore, error) {
	yubiStore, err := trustmanager.NewYubiKeyStore(fileKeyStore, ret)
	if err != nil {
		return nil, err
	}
	return yubiStore, nil
}
