package client

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

// TestValidateRoot through the process of initializing a repository and makes
// sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key
func TestValidateRoot(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	validateRootSuccessfully(t, data.ECDSAKey)
	if !testing.Short() {
		validateRootSuccessfully(t, data.RSAKey)
	}
}

func validateRootSuccessfully(t *testing.T, rootType string) {
	gun := "docker.com/notary"

	ts, mux, keys := simpleTestServer(t)
	defer ts.Close()

	repo, _ := initializeRepo(t, rootType, gun, ts.URL, false)
	defer os.RemoveAll(repo.baseDir)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err := repo.tufRepo.InitTimestamp()
	require.NoError(t, err, "error creating repository: %s", err)

	// Initialize is supposed to have created new certificate for this repository
	// Lets check for it and store it for later use
	allCerts := repo.CertStore.GetCertificates()
	require.Len(t, allCerts, 1)

	fakeServerData(t, repo, mux, keys)

	//
	// Test TOFUS logic. We remove all certs and expect a new one to be added after ListTargets
	//
	err = repo.CertStore.RemoveAll()
	require.NoError(t, err)
	require.Len(t, repo.CertStore.GetCertificates(), 0)

	// This list targets is expected to succeed and the certificate store to have the new certificate
	_, err = repo.ListTargets(data.CanonicalTargetsRole)
	require.NoError(t, err)
	require.Len(t, repo.CertStore.GetCertificates(), 1)

	//
	// Test certificate mismatch logic. We remove all certs, add a different cert to the
	// same CN, and expect ValidateRoot to fail
	//

	// First, remove all certs
	err = repo.CertStore.RemoveAll()
	require.NoError(t, err)
	require.Len(t, repo.CertStore.GetCertificates(), 0)

	// Add a previously generated certificate with CN=docker.com/notary
	err = repo.CertStore.AddCertFromFile(
		"../fixtures/self-signed_docker.com-notary.crt")
	require.NoError(t, err)

	// This list targets is expected to fail, since there already exists a certificate
	// in the store for the dnsName docker.com/notary, so TOFUS doesn't apply
	_, err = repo.ListTargets(data.CanonicalTargetsRole)
	require.Error(t, err, "An error was expected")
	require.Equal(t, err, &trustpinning.ErrValidationFail{
		Reason: "failed to validate data with current trusted certificates",
	})
}
