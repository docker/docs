package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
)

// TestInitRepo runs through the process of initializing a repository and makes
// sure the repository looks correct on disk.
func TestInitRepo(t *testing.T) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo, err := NewNotaryRepository(tempBaseDir, gun, "", nil)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID, err := repo.GenRootKey("passphrase")
	assert.NoError(t, err, "error generating root key: %s", err)

	rootSigner, err := repo.GetRootSigner(rootKeyID, "passphrase")
	assert.NoError(t, err, "error retreiving root key: %s", err)

	err = repo.Initialize(rootSigner)
	assert.NoError(t, err, "error creating repository: %s", err)

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
		assert.NoError(t, err, "missing directory in base directory: %s", dir)
		assert.True(t, fi.Mode().IsDir(), "%s is not a directory", dir)
	}

	// Look for keys in private. The filenames should match the key IDs
	// in the private key store.
	privKeyList := repo.privKeyStore.ListAll()
	for _, privKeyName := range privKeyList {
		_, err := os.Stat(privKeyName)
		assert.NoError(t, err, "missing private key: %s", privKeyName)
	}

	// Look for keys in root_keys
	// There should be a file named after the key ID of the root key we
	// passed in.
	rootKeyFilename := rootSigner.ID() + ".key"
	_, err = os.Stat(filepath.Join(tempBaseDir, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")

	// Also expect a symlink from the key ID of the certificate key to this
	// root key
	certificates := repo.certificateStore.GetCertificates()
	assert.Len(t, certificates, 1, "unexpected number of certificates")

	certID := trustmanager.FingerprintCert(certificates[0])

	actualDest, err := os.Readlink(filepath.Join(tempBaseDir, "private", "root_keys", certID+".key"))
	assert.NoError(t, err, "missing symlink to root key")

	assert.Equal(t, rootKeyFilename, actualDest, "symlink to root key has wrong destination")

	// There should be a trusted certificate
	_, err = os.Stat(filepath.Join(tempBaseDir, "trusted_certificates", gun, certID+".crt"))
	assert.NoError(t, err, "missing trusted certificate")

	// Sanity check the TUF metadata files. Verify that they exist, the JSON is
	// well-formed, and the signatures exist. For the root.json file, also check
	// that the root, snapshot, and targets key IDs are present.
	expectedTUFMetadataFiles := []string{
		filepath.Join("tuf", gun, "metadata", "root.json"),
		filepath.Join("tuf", gun, "metadata", "snapshot.json"),
		filepath.Join("tuf", gun, "metadata", "targets.json"),
	}
	for _, filename := range expectedTUFMetadataFiles {
		fullPath := filepath.Join(tempBaseDir, filename)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "missing TUF metadata file: %s", filename)

		jsonBytes, err := ioutil.ReadFile(fullPath)
		assert.NoError(t, err, "error reading TUF metadata file %s: %s", filename, err)

		var decoded data.Signed
		err = json.Unmarshal(jsonBytes, &decoded)
		assert.NoError(t, err, "error parsing TUF metadata file %s: %s", filename, err)

		assert.Len(t, decoded.Signatures, 1, "incorrect number of signatures in TUF metadata file %s", filename)

		assert.NotEmpty(t, decoded.Signatures[0].KeyID, "empty key ID field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Method, "empty method field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Signature, "empty signature in TUF metadata file %s", filename)

		// Special case for root.json: also check that the signed
		// content for keys and roles
		if strings.HasSuffix(filename, "root.json") {
			var decodedRoot data.Root
			err := json.Unmarshal(decoded.Signed, &decodedRoot)
			assert.NoError(t, err, "error parsing root.json signed section: %s", err)

			assert.Equal(t, "Root", decodedRoot.Type, "_type mismatch in root.json")

			// Expect 4 keys in the Keys map: root, targets, snapshot, timestamp
			assert.Len(t, decodedRoot.Keys, 4, "wrong number of keys in root.json")

			roleCount := 0
			for role := range decodedRoot.Roles {
				roleCount++
				if role != "root" && role != "snapshot" && role != "targets" && role != "timestamp" {
					t.Fatalf("unexpected role %s in root.json", role)
				}
			}
			assert.Equal(t, 4, roleCount, "wrong number of roles (%d) in root.json", roleCount)
		}
	}
}
