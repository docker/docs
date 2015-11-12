// +build pkcs11

package main

import (
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/trustmanager/yubikey"
)

func getYubiKeyStore(fileKeyStore trustmanager.KeyStore, ret passphrase.Retriever) (trustmanager.KeyStore, error) {
	return yubikey.NewYubiKeyStore(fileKeyStore, ret)
}
