// +build !pkcs11

package main

import (
	"testing"

	"github.com/docker/notary/passphrase"
)

func init() {
	retriever = passphrase.ConstantRetriever("pass")
	getRetriever = func() passphrase.Retriever { return retriever }
}

func rootOnHardware() bool {
	return false
}

// Per-test set up that is a no-op
func setUp(t *testing.T) {}

// no-op
func verifyRootKeyOnHardware(t *testing.T, rootKeyID string) {}
