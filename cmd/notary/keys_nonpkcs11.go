// +build !pkcs11

package main

import (
	"errors"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
)

func getYubiKeyStore(fileKeyStore trustmanager.KeyStore, ret passphrase.Retriever) (trustmanager.KeyStore, error) {
	return nil, errors.New("Not built with hardware support")
}
