package client

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

// TestInitRepo runs through the process of initializing a repository and makes
// sure the repository looks correct on disk.
func TestInitRepo(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	if err != nil {
		t.Fatalf("failed to create a temporary directory: %s", err)
	}

	client, err := NewClient(tempBaseDir)
	if err != nil {
		t.Fatalf("error creating client: %s", err)
	}

	rootKeyID, err := client.GenRootKey("passphrase")
	if err != nil {
		t.Fatalf("error generating root key: %s", err)
	}

	rootKey, err := client.GetRootKey(rootKeyID, "passphrase")
	if err != nil {
		t.Fatalf("error retreiving root key: %s", err)
	}

	gun := "docker.com/notary"
	repo, err := client.InitRepository(gun, "", nil, rootKey)
	if err != nil {
		t.Fatalf("error creating repository: %s", err)
	}

	// Inspect contents of the temporary directory
	expectedDirs := []string{
		"private",
		filepath.Join("private", gun),
		filepath.Join("private", "root_keys"),
		"trusted_certificates",
		filepath.Join("trusted_certificates", gun),
		"tuf",
		filepath.Join("tuf", gun, "metadata"),
		filepath.Join("tuf", gun, "targets"),
	}
	for _, dir := range expectedDirs {
		fi, err := os.Stat(filepath.Join(tempBaseDir, dir))
		if err != nil {
			t.Fatalf("missing directory in base directory: %s", dir)
		}
		if !fi.Mode().IsDir() {
			t.Fatalf("%s is not a directory", dir)
		}
	}

	// Look for keys in private. The filenames should match the key IDs
	// in the private key store.
	privKeyList := repo.privKeyStore.ListAll()
	for _, privKeyName := range privKeyList {
		if _, err := os.Stat(privKeyName); err != nil {
			t.Fatalf("missing private key: %s", privKeyName)
		}
	}

	// Look for keys in root_keys
	// There should be a file named after the key ID of the root key we
	// passed in.
	rootKeyFilename := rootKey.ID() + ".key"
	if _, err := os.Stat(filepath.Join(tempBaseDir, "private", "root_keys", rootKeyFilename)); err != nil {
		t.Fatal("missing root key")
	}

	// Also expect a symlink from the key ID of the certificate key to this
	// root key
	certificates := client.certificateStore.GetCertificates()
	if len(certificates) != 1 {
		t.Fatalf("unexpected number of certificates (%d)", len(certificates))
	}

	certID := data.NewPublicKey("RSA", trustmanager.CertToPEM(certificates[0])).ID()

	actualDest, err := os.Readlink(filepath.Join(tempBaseDir, "private", "root_keys", certID+".key"))
	if err != nil {
		t.Fatal("missing symlink to root key")
	}

	if actualDest != rootKeyFilename {
		t.Fatalf("symlink to root key has wrong destination (got: %s, expected: %s)", actualDest, rootKeyFilename)
	}

	// There should be a trusted certificate
	if _, err := os.Stat(filepath.Join(tempBaseDir, "trusted_certificates", gun, certID+".crt")); err != nil {
		t.Fatal("missing trusted certificate")
	}

	// Check that tuf metadata files were created
	expectedTUFMetadataFiles := []string{
		filepath.Join("tuf", gun, "metadata", "root.json"),
		filepath.Join("tuf", gun, "metadata", "snapshot.json"),
		filepath.Join("tuf", gun, "metadata", "targets.json"),
	}
	for _, filename := range expectedTUFMetadataFiles {
		_, err := os.Stat(filepath.Join(tempBaseDir, filename))
		if err != nil {
			t.Fatalf("missing TUF metadata file: %s", filename)
		}
	}
}
