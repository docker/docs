package client

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
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
	// Test TOFUS logic.
	//

	_, err = repo.ListTargets(data.CanonicalTargetsRole)
	require.NoError(t, err)
}
