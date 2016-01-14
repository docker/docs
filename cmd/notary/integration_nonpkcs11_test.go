// +build !pkcs11

package main

import (
	"testing"

	"github.com/docker/notary/passphrase"
)

func init() {
	fake := passphrase.ConstantRetriever("pass")
	retriever = fake
	getRetriever = func() passphrase.Retriever { return fake }
}

func rootOnHardware() bool {
	return false
}

// Per-test set up that is a no-op
func setUp(t *testing.T) {}

// no-op
func verifyRootKeyOnHardware(t *testing.T, rootKeyID string) {}
