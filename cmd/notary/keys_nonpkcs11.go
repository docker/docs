// +build !pkcs11

package main

import (
	"errors"

	"github.com/docker/notary"
	"github.com/docker/notary/trustmanager"
)

func getYubiStore(fileKeyStore trustmanager.KeyStore, ret notary.PassRetriever) (trustmanager.KeyStore, error) {
	return nil, errors.New("Not built with hardware support")
}
