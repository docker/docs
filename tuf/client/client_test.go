package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	tuf "github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/require"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
)

func TestRotation(t *testing.T) {
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(signer)
	remote := store.NewMemoryStore(nil)
	cache := store.NewMemoryStore(nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil)
	require.NoError(t, err, "Error creating root role")

	originalRoot, err := data.NewRoot(
		map[string]data.PublicKey{rootKey.ID(): rootKey},
		map[string]*data.RootRole{"root": &rootRole.RootRole},
		false,
	)

	repo.Root = originalRoot

	// Generate new key and role.
	replacementKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil)
	require.NoError(t, err, "Error creating replacement root role")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{
			data.CanonicalRootRole:      &replacementRole.RootRole,
			data.CanonicalSnapshotRole:  &replacementRole.RootRole,
			data.CanonicalTargetsRole:   &replacementRole.RootRole,
			data.CanonicalTimestampRole: &replacementRole.RootRole,
		},
		false,
	)
	require.NoError(t, err, "Failed to create new root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, []data.PublicKey{rootKey, replacementKey})
	require.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	require.True(t, origKeySig, "Original root key signature not present")
	require.True(t, replKeySig, "Replacement root key signature not present")

	client := NewClient(repo, remote, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	require.NoError(t, err, "Failed to verify key rotated root")
}

func TestRotationNewSigMissing(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(signer)
	remote := store.NewMemoryStore(nil)
	cache := store.NewMemoryStore(nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil)
	require.NoError(t, err, "Error creating root role")

	originalRoot, err := data.NewRoot(
		map[string]data.PublicKey{rootKey.ID(): rootKey},
		map[string]*data.RootRole{"root": &rootRole.RootRole},
		false,
	)

	repo.Root = originalRoot

	// Generate new key and role.
	replacementKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil)
	require.NoError(t, err, "Error creating replacement root role")

	require.NotEqual(t, rootKey.ID(), replacementKey.ID(), "Key IDs are the same")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{"root": &replacementRole.RootRole},
		false,
	)
	require.NoError(t, err, "Failed to create new root")

	_, ok := testRoot.Signed.Keys[rootKey.ID()]
	require.False(t, ok, "Old root key appeared in test root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, []data.PublicKey{rootKey})
	require.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	require.True(t, origKeySig, "Original root key signature not present")
	require.False(t, replKeySig, "Replacement root key signature was present and shouldn't be")

	client := NewClient(repo, remote, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	require.Error(t, err, "Should have errored on verify as replacement signature was missing.")

}

func TestRotationOldSigMissing(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(signer)
	remote := store.NewMemoryStore(nil)
	cache := store.NewMemoryStore(nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil)
	require.NoError(t, err, "Error creating root role")

	originalRoot, err := data.NewRoot(
		map[string]data.PublicKey{rootKey.ID(): rootKey},
		map[string]*data.RootRole{"root": &rootRole.RootRole},
		false,
	)

	repo.Root = originalRoot

	// Generate new key and role.
	replacementKey, err := signer.Create("root", "", data.ED25519Key)
	require.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil)
	require.NoError(t, err, "Error creating replacement root role")

	require.NotEqual(t, rootKey.ID(), replacementKey.ID(), "Key IDs are the same")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{"root": &replacementRole.RootRole},
		false,
	)
	require.NoError(t, err, "Failed to create new root")

	_, ok := testRoot.Signed.Keys[rootKey.ID()]
	require.False(t, ok, "Old root key appeared in test root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, []data.PublicKey{replacementKey})
	require.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	require.False(t, origKeySig, "Original root key signature was present and shouldn't be")
	require.True(t, replKeySig, "Replacement root key signature was not present")

	client := NewClient(repo, remote, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	require.Error(t, err, "Should have errored on verify as replacement signature was missing.")

}

func TestCheckRootExpired(t *testing.T) {
	repo := tuf.NewRepo(nil)
	storage := store.NewMemoryStore(nil)
	client := NewClient(repo, storage, storage)

	root := &data.SignedRoot{}
	root.Signed.Expires = time.Now().AddDate(-1, 0, 0)

	signedRoot, err := root.ToSigned()
	require.NoError(t, err)
	rootJSON, err := json.Marshal(signedRoot)
	require.NoError(t, err)

	rootHash := sha256.Sum256(rootJSON)

	testSnap := &data.SignedSnapshot{
		Signed: data.Snapshot{
			Meta: map[string]data.FileMeta{
				"root": {
					Length: int64(len(rootJSON)),
					Hashes: map[string][]byte{
						"sha256": rootHash[:],
					},
				},
			},
		},
	}
	repo.SetRoot(root)
	repo.SetSnapshot(testSnap)

	storage.SetMeta("root", rootJSON)

	err = client.checkRoot()
	require.Error(t, err)
	require.IsType(t, tuf.ErrLocalRootExpired{}, err)
}

func TestChecksumMismatch(t *testing.T) {
	repo := tuf.NewRepo(nil)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := testutils.NewCorruptingMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	require.NoError(t, err)

	origHashes, err := GetSupportedHashes(orig)
	require.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", int64(len(orig)), origHashes)
	require.IsType(t, ErrChecksumMismatch{}, err)
}

func TestChecksumMatch(t *testing.T) {
	repo := tuf.NewRepo(nil)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	require.NoError(t, err)

	origHashes, err := GetSupportedHashes(orig)
	require.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", int64(len(orig)), origHashes)
	require.NoError(t, err)
}

func TestSizeMismatchLong(t *testing.T) {
	repo := tuf.NewRepo(nil)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := testutils.NewLongMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	require.NoError(t, err)
	l := int64(len(orig))

	origHashes, err := GetSupportedHashes(orig)
	require.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", l, origHashes)
	// size just limits the data received, the error is caught
	// either during checksum verification or during json deserialization
	require.IsType(t, ErrChecksumMismatch{}, err)
}

func TestSizeMismatchShort(t *testing.T) {
	repo := tuf.NewRepo(nil)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := testutils.NewShortMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	require.NoError(t, err)
	l := int64(len(orig))

	origHashes, err := GetSupportedHashes(orig)
	require.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", l, origHashes)
	// size just limits the data received, the error is caught
	// either during checksum verification or during json deserialization
	require.IsType(t, ErrChecksumMismatch{}, err)
}

func TestDownloadTargetsHappy(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	err = client.downloadTargets("targets")
	require.NoError(t, err)
}

// TestDownloadTargetsLarge: Check that we can download very large targets metadata files,
// which may be caused by adding a large number of targets.
// This test is slow, so it will not run in short mode.
func TestDownloadTargetsLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	// Add a ton of target files to the targets role to make this targets metadata huge
	// 75,000 targets results in > 5MB (~6.5MB on recent runs)
	for i := 0; i < 75000; i++ {
		_, err = repo.AddTargets(data.CanonicalTargetsRole, data.Files{strconv.Itoa(i): f})
		require.NoError(t, err)
	}

	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	// Clear the cache to force an online download
	client.cache.RemoveAll()

	err = client.downloadTargets("targets")
	require.NoError(t, err)
}

func TestDownloadTargetsDeepHappy(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	delegations := []string{
		// left subtree
		"targets/level1",
		"targets/level1/a",
		"targets/level1/a/i",
		"targets/level1/a/ii",
		"targets/level1/a/iii",
		// right subtree
		"targets/level2",
		"targets/level2/b",
		"targets/level2/b/i",
		"targets/level2/b/i/0",
		"targets/level2/b/i/1",
	}

	for _, r := range delegations {
		// create role
		k, err := cs.Create(r, "docker.com/notary", data.ED25519Key)
		require.NoError(t, err)

		// add role to repo
		err = repo.UpdateDelegationKeys(r, []data.PublicKey{k}, []string{}, 1)
		require.NoError(t, err)
		err = repo.UpdateDelegationPaths(r, []string{""}, []string{}, false)
		require.NoError(t, err)
		repo.InitTargets(r)
	}

	// can only sign after adding all delegations
	for _, r := range delegations {
		// serialize and store role
		signedOrig, err := repo.SignTargets(r, data.DefaultExpires("targets"))
		require.NoError(t, err)
		orig, err := json.Marshal(signedOrig)
		require.NoError(t, err)
		err = remoteStorage.SetMeta(r, orig)
		require.NoError(t, err)
	}

	// serialize and store targets after adding all delegations
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	delete(repo.Targets, "targets")
	for _, r := range delegations {
		delete(repo.Targets, r)
		_, ok := repo.Targets[r]
		require.False(t, ok)
	}

	err = client.downloadTargets("targets")
	require.NoError(t, err)

	_, ok := repo.Targets["targets"]
	require.True(t, ok)

	for _, r := range delegations {
		_, ok = repo.Targets[r]
		require.True(t, ok)
	}
}

func TestDownloadTargetChecksumMismatch(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := testutils.NewCorruptingMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	origSha256 := sha256.Sum256(orig)
	err = remoteStorage.SetMeta("targets", orig)
	require.NoError(t, err)

	// create local snapshot with targets file
	snap := data.SignedSnapshot{
		Signed: data.Snapshot{
			Meta: data.Files{
				"targets": data.FileMeta{
					Length: int64(len(orig)),
					Hashes: data.Hashes{
						"sha256": origSha256[:],
					},
				},
			},
		},
	}

	repo.Snapshot = &snap

	err = client.downloadTargets("targets")
	require.IsType(t, ErrChecksumMismatch{}, err)
}

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

func TestBootstrapDownloadRootHappy(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample root
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	require.NoError(t, err)

	// unset snapshot as if we're bootstrapping from nothing
	repo.Snapshot = nil

	err = client.downloadRoot()
	require.NoError(t, err)
}

func TestUpdateDownloadRootHappy(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample root, snapshot, and timestamp
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	require.NoError(t, err)

	// sign snapshot to make root meta in snapshot get updated
	signedOrig, err = repo.SignSnapshot(data.DefaultExpires("snapshot"))

	err = client.downloadRoot()
	require.NoError(t, err)
}

func TestUpdateDownloadRootBadChecksum(t *testing.T) {
	remoteStore := testutils.NewCorruptingMemoryStore(nil)

	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStore, localStorage)

	// sign and "upload" sample root
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStore.SetMeta("root", orig)
	require.NoError(t, err)

	// sign snapshot to make sure we have current checksum for root
	_, err = repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)

	err = client.downloadRoot()
	require.IsType(t, ErrChecksumMismatch{}, err)
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

func TestDownloadTimestampHappy(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// create and "upload" sample timestamp
	signedOrig, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	require.NoError(t, err)

	err = client.downloadTimestamp()
	require.NoError(t, err)
}

func TestDownloadSnapshotHappy(t *testing.T) {
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

	signedOrig, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	orig, err = json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	require.NoError(t, err)

	err = client.downloadSnapshot()
	require.NoError(t, err)
}

// TestDownloadSnapshotLarge: Check that we can download very large snapshot metadata files,
// which may be caused by adding a large number of delegations.
// This test is slow, so it will not run in short mode.
func TestDownloadSnapshotLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// Add a ton of empty delegation roles to targets to make snapshot data huge
	// This can also be done by adding legitimate delegations but it will be much slower
	// 75,000 delegation roles results in > 5MB (~7.3MB on recent runs)
	for i := 0; i < 75000; i++ {
		newRole := &data.SignedTargets{}
		repo.Targets[fmt.Sprintf("targets/%d", i)] = newRole
	}

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	require.NoError(t, err)

	signedOrig, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	orig, err = json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	require.NoError(t, err)

	// Clear the cache to force an online download
	client.cache.RemoveAll()

	err = client.downloadSnapshot()
	require.NoError(t, err)
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

func TestDownloadSnapshotChecksumNotFound(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	// sign timestamp to ensure it has a checksum for snapshot
	_, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	require.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	require.NoError(t, err)

	// by not signing timestamp again we ensure it has the wrong checksum

	err = client.downloadSnapshot()
	require.IsType(t, store.ErrMetaNotFound{}, err)
}

// If there is no local cache and also no remote timestamp, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoTimestamps(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	err = client.downloadTimestamp()
	require.Error(t, err)
	notFoundErr, ok := err.(store.ErrMetaNotFound)
	require.True(t, ok)
	require.Equal(t, data.CanonicalTimestampRole, notFoundErr.Resource)
}

// If there is no local cache and the remote timestamp is empty, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoLocalTimestampRemoteTimestampEmpty(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: {}})
	client := NewClient(repo, remoteStorage, localStorage)

	err = client.downloadTimestamp()
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
}

// If there is no local cache and the remote timestamp is invalid, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoLocalTimestampRemoteTimestampInvalid(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(nil)

	// add a timestamp to the remote cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	tsSigned.Signatures[0].Signature = []byte("12345") // invalidate the signature
	ts, err := json.Marshal(tsSigned)
	require.NoError(t, err)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts})

	client := NewClient(repo, remoteStorage, localStorage)
	err = client.downloadTimestamp()
	require.Error(t, err)
	require.IsType(t, signed.ErrRoleThreshold{}, err)
}

// If there is is a local cache and no remote timestamp, we fall back on the cached timestamp
func TestDownloadTimestampLocalTimestampNoRemoteTimestamp(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	// add a timestamp to the local cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	ts, err := json.Marshal(tsSigned)
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts})

	remoteStorage := store.NewMemoryStore(nil)
	client := NewClient(repo, remoteStorage, localStorage)

	err = client.downloadTimestamp()
	require.NoError(t, err)
}

// If there is is a local cache and the remote timestamp is invalid, we fall back on the cached timestamp
func TestDownloadTimestampLocalTimestampInvalidRemoteTimestamp(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	require.NoError(t, err)

	// add a timestamp to the local cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	ts, err := json.Marshal(tsSigned)
	require.NoError(t, err)
	localStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts})

	// add a timestamp to the remote cache
	tsSigned.Signatures[0].Signature = []byte("12345") // invalidate the signature
	ts, err = json.Marshal(tsSigned)
	require.NoError(t, err)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts})

	client := NewClient(repo, remoteStorage, localStorage)
	err = client.downloadTimestamp()
	require.NoError(t, err)
}

// GetSupportedHashes is a helper function that returns
// the checksums of all the supported hash algorithms
// of the given payload.
func GetSupportedHashes(payload []byte) (data.Hashes, error) {
	meta, err := data.NewFileMeta(bytes.NewReader(payload), data.NotaryDefaultHashes...)
	if err != nil {
		return nil, err
	}

	return meta.Hashes, nil
}
