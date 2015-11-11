// +build pkcs11

package main

import (
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
)

func getYubiKeyStore(fileKeyStore trustmanager.KeyStore, ret passphrase.Retriever) (trustmanager.KeyStore, error) {
	return trustmanager.NewYubiKeyStore(fileKeyStore, ret)
}
