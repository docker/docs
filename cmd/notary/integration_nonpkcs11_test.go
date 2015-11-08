// +build !pkcs11

package main

import (
	"testing"

	"github.com/docker/notary/passphrase"
)

func rootOnHardware() bool {
	return false
}

// Per-test set up that returns a cleanup function.  This set up changes the
// passphrase retriever to always produce a constant passphrase
func setUp(t *testing.T) func() {
	oldRetriever := retriever

	var fake = func(k, a string, c bool, n int) (string, bool, error) {
		return testPassphrase, false, nil
	}

	retriever = fake
	getRetriever = func() passphrase.Retriever { return fake }

	return func() {
		retriever = oldRetriever
		getRetriever = getPassphraseRetriever
	}
}

// no-op
func verifyRootKeyOnHardware(t *testing.T, rootKeyID string) {}
