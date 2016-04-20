package client

import (
	"encoding/json"
	"testing"

	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/store"
)

// TestDownloadTargetsNoChecksum: it's never valid to download any targets
// role (incl. delegations) when a checksum is not available.
func TestDownloadTargetsNoChecksum(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	delete(repo.Snapshot.Signed.Meta["targets"].Hashes, "sha256")
	delete(repo.Snapshot.Signed.Meta["targets"].Hashes, "sha512")

	err = client.downloadTargets("targets")
	require.IsType(t, data.ErrMissingMeta{}, err)
}

// TestDownloadTargetsNoSnapshot: it's never valid to download any targets
// role (incl. delegations) when a checksum is not available.
func TestDownloadTargetsNoSnapshot(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	repo.Snapshot = nil

	err = client.downloadTargets("targets")
	require.IsType(t, tuf.ErrNotLoaded{}, err)
}

func TestUpdateDownloadRootChecksumNotFound(t *testing.T) {
	remoteStore := store.NewMemoryStore(nil)
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStore, localStorage)

	// sign snapshot to make sure we have current checksum for root
	_, err = repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)

	// sign and "upload" sample root
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStore.SetMeta("root", orig)
	require.NoError(t, err)

	// don't sign snapshot again to ensure checksum is out of date (bad)

	err = client.downloadRoot()
	require.IsType(t, store.ErrMetaNotFound{}, err)
}

func TestDownloadSnapshotNoTimestamp(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	require.NoError(t, err)

	repo.Timestamp = nil

	err = client.downloadSnapshot()
	require.IsType(t, tuf.ErrNotLoaded{}, err)
}

// TestDownloadSnapshotNoChecksum: It should never be valid to download a
// snapshot if we don't have a checksum
func TestDownloadSnapshotNoChecksum(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	require.NoError(t, err)

	delete(repo.Timestamp.Signed.Meta["snapshot"].Hashes, "sha256")
	delete(repo.Timestamp.Signed.Meta["snapshot"].Hashes, "sha512")

	err = client.downloadSnapshot()
	require.IsType(t, data.ErrMissingMeta{}, err)
}
