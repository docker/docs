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

	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

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

	gun := "docker.io/notary0.1/samplerepo"
	passwd := "randompass"

	ts := fullTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tmpDir, gun, ts.URL, http.DefaultTransport,
		passphrase.ConstantRetriever(passwd))
	require.NoError(t, err, "error creating repo: %s", err)

	// rotate the timestamp key, since the server doesn't have that one
	timestampPubKey, err := getRemoteKey(ts.URL, gun, data.CanonicalTimestampRole, http.DefaultTransport)
	require.NoError(t, err)
	require.NoError(
		t, repo.rootFileKeyChange(data.CanonicalTimestampRole, changelist.ActionCreate, timestampPubKey))

	require.NoError(t, repo.Publish())

	targets, err := repo.ListTargets()
	require.NoError(t, err)
	require.Len(t, targets, 1)
	require.Equal(t, "v1", targets[0].Name)

	// Also check that we can add/remove keys by rotating keys
	oldTargetsKeys := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)
	require.NoError(t, repo.RotateKey(data.CanonicalTargetsRole, false))
	require.NoError(t, repo.Publish())
	newTargetsKeys := repo.CryptoService.ListKeys(data.CanonicalTargetsRole)

	require.Len(t, oldTargetsKeys, 1)
	require.Len(t, newTargetsKeys, 1)
	require.NotEqual(t, oldTargetsKeys[0], newTargetsKeys[0])
}
