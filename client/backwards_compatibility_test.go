// The client can read and operate on older repository formats

package client

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/store"
	"github.com/stretchr/testify/require"
)

// Once a fixture is read in, ensure that it's valid by making sure the expiry
// times of all the metadata and certificates is > 10 years ahead
func requireValidFixture(t *testing.T, notaryRepo *NotaryRepository) {
	tenYearsInFuture := time.Now().AddDate(10, 0, 0)
	require.True(t, notaryRepo.tufRepo.Root.Signed.Expires.After(tenYearsInFuture))
	require.True(t, notaryRepo.tufRepo.Snapshot.Signed.Expires.After(tenYearsInFuture))
	require.True(t, notaryRepo.tufRepo.Timestamp.Signed.Expires.After(tenYearsInFuture))
	for _, targetObj := range notaryRepo.tufRepo.Targets {
		require.True(t, targetObj.Signed.Expires.After(tenYearsInFuture))
	}
}

// recursively copies the contents of one directory into another - ignores
// symlinks
func recursiveCopy(sourceDir, targetDir string) error {
	return filepath.Walk(sourceDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		targetFP := filepath.Join(targetDir, strings.TrimPrefix(fp, sourceDir+"/"))

		if fi.IsDir() {
			return os.MkdirAll(targetFP, fi.Mode())
		}

		// Ignore symlinks
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		// copy the file
		in, err := os.Open(fp)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(targetFP)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return nil
	})
}

// We can read and publish from notary0.1 repos
func Test0Dot1RepoFormat(t *testing.T) {
	// make a temporary directory and copy the fixture into it, since updating
	// and publishing will modify the files
	tmpDir, err := ioutil.TempDir("", "notary-backwards-compat-test")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)
	require.NoError(t, recursiveCopy("../fixtures/compatibility/notary0.1", tmpDir))

	gun := "docker.com/notary0.1/samplerepo"
	passwd := "randompass"

	ts := fullTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tmpDir, gun, ts.URL, http.DefaultTransport,
		passphrase.ConstantRetriever(passwd), trustpinning.TrustPinConfig{})
	require.NoError(t, err, "error creating repo: %s", err)

	// targets should have 1 target, and it should be readable offline
	targets, err := repo.ListTargets()
	require.NoError(t, err)
	require.Len(t, targets, 1)
	require.Equal(t, "LICENSE", targets[0].Name)

	// ok, now that everything has been loaded, verify that the fixture is valid
	requireValidFixture(t, repo)

	// delete the timestamp metadata, since the server will ignore the uploaded
	// one and try to create a new one from scratch, which will be the wrong version
	require.NoError(t, repo.fileStore.RemoveMeta(data.CanonicalTimestampRole))

	// rotate the timestamp key, since the server doesn't have that one
	err = repo.RotateKey(data.CanonicalTimestampRole, true)
	require.NoError(t, err)

	require.NoError(t, repo.Publish())

	targets, err = repo.ListTargets()
	require.NoError(t, err)
	require.Len(t, targets, 2)

	// Also check that we can add/remove keys by rotating keys
	oldTargetsKeys := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
	require.NoError(t, repo.RotateKey(data.CanonicalTargetsRole, false))
	require.NoError(t, repo.Publish())
	newTargetsKeys := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)

	require.Len(t, oldTargetsKeys, 1)
	require.Len(t, newTargetsKeys, 1)
	require.NotEqual(t, oldTargetsKeys[0], newTargetsKeys[0])

	// rotate the snapshot key to the server and ensure that the server can re-generate the snapshot
	// and we can download the snapshot
	require.NoError(t, repo.RotateKey(data.CanonicalSnapshotRole, true))
	require.NoError(t, repo.Publish())
	err = repo.Update(false)
	require.NoError(t, err)
}

// Ensures that the current client can download metadata that is published from notary 0.1 repos
func TestDownloading0Dot1RepoFormat(t *testing.T) {
	gun := "docker.com/notary0.1/samplerepo"
	passwd := "randompass"

	metaCache, err := store.NewFilesystemStore(
		filepath.Join("../fixtures/compatibility/notary0.1/tuf", filepath.FromSlash(gun)),
		"metadata", "json")
	require.NoError(t, err)

	ts := readOnlyServer(t, metaCache, http.StatusNotFound, gun)
	defer ts.Close()

	repoDir, err := ioutil.TempDir("", "notary-backwards-compat-test")
	require.NoError(t, err)
	defer os.RemoveAll(repoDir)

	repo, err := NewNotaryRepository(repoDir, gun, ts.URL, http.DefaultTransport,
		passphrase.ConstantRetriever(passwd), trustpinning.TrustPinConfig{})
	require.NoError(t, err, "error creating repo: %s", err)

	err = repo.Update(true)
	require.NoError(t, err, "error updating repo: %s", err)
}
