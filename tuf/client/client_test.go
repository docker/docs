package client

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	tuf "github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
)

func TestRotation(t *testing.T) {
	kdb := keys.NewDB()
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(kdb, signer)
	remote := store.NewMemoryStore(nil, nil)
	cache := store.NewMemoryStore(nil, nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating root role")

	kdb.AddKey(rootKey)
	err = kdb.AddRole(rootRole)
	assert.NoError(t, err, "Error adding root role to db")

	// Generate new key and role. These will appear in the root.json
	// but will not be added to the keyDB.
	replacementKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating replacement root role")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{"root": &replacementRole.RootRole},
		false,
	)
	assert.NoError(t, err, "Failed to create new root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, rootKey, replacementKey)
	assert.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	assert.True(t, origKeySig, "Original root key signature not present")
	assert.True(t, replKeySig, "Replacement root key signature not present")

	client := NewClient(repo, remote, kdb, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	assert.NoError(t, err, "Failed to verify key rotated root")
}

func TestRotationNewSigMissing(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	kdb := keys.NewDB()
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(kdb, signer)
	remote := store.NewMemoryStore(nil, nil)
	cache := store.NewMemoryStore(nil, nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating root role")

	kdb.AddKey(rootKey)
	err = kdb.AddRole(rootRole)
	assert.NoError(t, err, "Error adding root role to db")

	// Generate new key and role. These will appear in the root.json
	// but will not be added to the keyDB.
	replacementKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating replacement root role")

	assert.NotEqual(t, rootKey.ID(), replacementKey.ID(), "Key IDs are the same")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{"root": &replacementRole.RootRole},
		false,
	)
	assert.NoError(t, err, "Failed to create new root")

	_, ok := testRoot.Signed.Keys[rootKey.ID()]
	assert.False(t, ok, "Old root key appeared in test root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, rootKey)
	assert.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	assert.True(t, origKeySig, "Original root key signature not present")
	assert.False(t, replKeySig, "Replacement root key signature was present and shouldn't be")

	client := NewClient(repo, remote, kdb, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	assert.Error(t, err, "Should have errored on verify as replacement signature was missing.")

}

func TestRotationOldSigMissing(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	kdb := keys.NewDB()
	signer := signed.NewEd25519()
	repo := tuf.NewRepo(kdb, signer)
	remote := store.NewMemoryStore(nil, nil)
	cache := store.NewMemoryStore(nil, nil)

	// Generate initial root key and role and add to key DB
	rootKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating root key")
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating root role")

	kdb.AddKey(rootKey)
	err = kdb.AddRole(rootRole)
	assert.NoError(t, err, "Error adding root role to db")

	// Generate new key and role. These will appear in the root.json
	// but will not be added to the keyDB.
	replacementKey, err := signer.Create("root", data.ED25519Key)
	assert.NoError(t, err, "Error creating replacement root key")
	replacementRole, err := data.NewRole("root", 1, []string{replacementKey.ID()}, nil, nil)
	assert.NoError(t, err, "Error creating replacement root role")

	assert.NotEqual(t, rootKey.ID(), replacementKey.ID(), "Key IDs are the same")

	// Generate a new root with the replacement key and role
	testRoot, err := data.NewRoot(
		map[string]data.PublicKey{replacementKey.ID(): replacementKey},
		map[string]*data.RootRole{"root": &replacementRole.RootRole},
		false,
	)
	assert.NoError(t, err, "Failed to create new root")

	_, ok := testRoot.Signed.Keys[rootKey.ID()]
	assert.False(t, ok, "Old root key appeared in test root")

	// Sign testRoot with both old and new keys
	signedRoot, err := testRoot.ToSigned()
	err = signed.Sign(signer, signedRoot, replacementKey)
	assert.NoError(t, err, "Failed to sign root")
	var origKeySig bool
	var replKeySig bool
	for _, sig := range signedRoot.Signatures {
		if sig.KeyID == rootKey.ID() {
			origKeySig = true
		} else if sig.KeyID == replacementKey.ID() {
			replKeySig = true
		}
	}
	assert.False(t, origKeySig, "Original root key signature was present and shouldn't be")
	assert.True(t, replKeySig, "Replacement root key signature was not present")

	client := NewClient(repo, remote, kdb, cache)

	err = client.verifyRoot("root", signedRoot, 0)
	assert.Error(t, err, "Should have errored on verify as replacement signature was missing.")

}

func TestCheckRootExpired(t *testing.T) {
	repo := tuf.NewRepo(nil, nil)
	storage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, storage, nil, storage)

	root := &data.SignedRoot{}
	root.Signed.Expires = time.Now().AddDate(-1, 0, 0)

	signedRoot, err := root.ToSigned()
	assert.NoError(t, err)
	rootJSON, err := json.Marshal(signedRoot)
	assert.NoError(t, err)

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
	assert.Error(t, err)
	assert.IsType(t, tuf.ErrLocalRootExpired{}, err)
}

func TestChecksumMismatch(t *testing.T) {
	repo := tuf.NewRepo(nil, nil)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, nil, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	origSha256 := sha256.Sum256(orig)
	orig[0] = '}' // corrupt data, should be a {
	assert.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", int64(len(orig)), origSha256[:])
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

func TestChecksumMatch(t *testing.T) {
	repo := tuf.NewRepo(nil, nil)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, nil, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	origSha256 := sha256.Sum256(orig)
	assert.NoError(t, err)

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", int64(len(orig)), origSha256[:])
	assert.NoError(t, err)
}

func TestSizeMismatchLong(t *testing.T) {
	repo := tuf.NewRepo(nil, nil)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, nil, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	origSha256 := sha256.Sum256(orig)
	assert.NoError(t, err)
	l := int64(len(orig))

	orig = append([]byte(" "), orig...)
	assert.Equal(t, l+1, int64(len(orig)))

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", l, origSha256[:])
	// size just limits the data received, the error is caught
	// either during checksum verification or during json deserialization
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

func TestSizeMismatchShort(t *testing.T) {
	repo := tuf.NewRepo(nil, nil)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, nil, localStorage)

	sampleTargets := data.NewTargets()
	orig, err := json.Marshal(sampleTargets)
	origSha256 := sha256.Sum256(orig)
	assert.NoError(t, err)
	l := int64(len(orig))

	orig = orig[1:]

	remoteStorage.SetMeta("targets", orig)

	_, _, err = client.downloadSigned("targets", l, origSha256[:])
	// size just limits the data received, the error is caught
	// either during checksum verification or during json deserialization
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

func TestDownloadTargetsHappy(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	err = client.downloadTargets("targets")
	assert.NoError(t, err)
}

// TestDownloadTargetsLarge: Check that we can download very large targets metadata files,
// which may be caused by adding a large number of targets.
// This test is slow, so it will not run in short mode.
func TestDownloadTargetsLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

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
		assert.NoError(t, err)
	}

	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	// Clear the cache to force an online download
	client.cache.RemoveAll()

	err = client.downloadTargets("targets")
	assert.NoError(t, err)
}

func TestDownloadTargetsDeepHappy(t *testing.T) {
	kdb, repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

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
		k, err := cs.Create(r, data.ED25519Key)
		assert.NoError(t, err)
		role, err := data.NewRole(r, 1, []string{k.ID()}, []string{""}, nil)
		assert.NoError(t, err)

		// add role to repo
		repo.UpdateDelegations(role, []data.PublicKey{k})
		repo.InitTargets(r)
	}

	// can only sign after adding all delegations
	for _, r := range delegations {
		// serialize and store role
		signedOrig, err := repo.SignTargets(r, data.DefaultExpires("targets"))
		assert.NoError(t, err)
		orig, err := json.Marshal(signedOrig)
		assert.NoError(t, err)
		err = remoteStorage.SetMeta(r, orig)
		assert.NoError(t, err)
	}

	// serialize and store targets after adding all delegations
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"))

	delete(repo.Targets, "targets")
	for _, r := range delegations {
		delete(repo.Targets, r)
		_, ok := repo.Targets[r]
		assert.False(t, ok)
	}

	err = client.downloadTargets("targets")
	assert.NoError(t, err)

	_, ok := repo.Targets["targets"]
	assert.True(t, ok)

	for _, r := range delegations {
		_, ok = repo.Targets[r]
		assert.True(t, ok)
	}
}

func TestDownloadTargetChecksumMismatch(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	origSha256 := sha256.Sum256(orig)
	orig[0] = '}' // corrupt data, should be a {
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	// create local snapshot with targets file
	// It's necessary to do it this way rather than calling repo.SignSnapshot
	// so that we have the wrong sha256 in the snapshot.
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
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

// TestDownloadTargetsNoChecksum: it's never valid to download any targets
// role (incl. delegations) when a checksum is not available.
func TestDownloadTargetsNoChecksum(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	delete(repo.Snapshot.Signed.Meta["targets"].Hashes, "sha256")

	err = client.downloadTargets("targets")
	assert.IsType(t, ErrMissingMeta{}, err)
}

// TestDownloadTargetsNoSnapshot: it's never valid to download any targets
// role (incl. delegations) when a checksum is not available.
func TestDownloadTargetsNoSnapshot(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	repo.Snapshot = nil

	err = client.downloadTargets("targets")
	assert.IsType(t, ErrMissingMeta{}, err)
}

func TestBootstrapDownloadRootHappy(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample root
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	assert.NoError(t, err)

	// unset snapshot as if we're bootstrapping from nothing
	repo.Snapshot = nil

	err = client.downloadRoot()
	assert.NoError(t, err)
}

func TestUpdateDownloadRootHappy(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample root, snapshot, and timestamp
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	assert.NoError(t, err)

	// sign snapshot to make root meta in snapshot get updated
	signedOrig, err = repo.SignSnapshot(data.DefaultExpires("snapshot"))

	err = client.downloadRoot()
	assert.NoError(t, err)
}

func TestUpdateDownloadRootBadChecksum(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// sign snapshot to make sure we have a checksum for root
	_, err = repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)

	// create and "upload" sample root, snapshot, and timestamp
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	assert.NoError(t, err)

	// don't sign snapshot again to ensure checksum is out of date (bad)

	err = client.downloadRoot()
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

func TestDownloadTimestampHappy(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample timestamp
	signedOrig, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	assert.NoError(t, err)

	err = client.downloadTimestamp()
	assert.NoError(t, err)
}

func TestDownloadSnapshotHappy(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	signedOrig, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	orig, err = json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	assert.NoError(t, err)

	err = client.downloadSnapshot()
	assert.NoError(t, err)
}

// TestDownloadSnapshotLarge: Check that we can download very large snapshot metadata files,
// which may be caused by adding a large number of delegations.
// This test is slow, so it will not run in short mode.
func TestDownloadSnapshotLarge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// Add a ton of empty delegation roles to targets to make snapshot data huge
	// This can also be done by adding legitimate delegations but it will be much slower
	// 75,000 delegation roles results in > 5MB (~7.3MB on recent runs)
	for i := 0; i < 75000; i++ {
		newRole := &data.SignedTargets{}
		repo.Targets[fmt.Sprintf("targets/%d", i)] = newRole
	}

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	signedOrig, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	orig, err = json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	assert.NoError(t, err)

	// Clear the cache to force an online download
	client.cache.RemoveAll()

	err = client.downloadSnapshot()
	assert.NoError(t, err)
}

// TestDownloadSnapshotNoChecksum: It should never be valid to download a
// snapshot if we don't have a checksum
func TestDownloadSnapshotNoTimestamp(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	repo.Timestamp = nil

	err = client.downloadSnapshot()
	assert.IsType(t, ErrMissingMeta{}, err)
}

func TestDownloadSnapshotNoChecksum(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	delete(repo.Timestamp.Signed.Meta["snapshot"].Hashes, "sha256")

	err = client.downloadSnapshot()
	assert.IsType(t, ErrMissingMeta{}, err)
}

func TestDownloadSnapshotBadChecksum(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// sign timestamp to ensure it has a checksum for snapshot
	_, err = repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	// by not signing timestamp again we ensure it has the wrong checksum

	err = client.downloadSnapshot()
	assert.IsType(t, ErrChecksumMismatch{}, err)
}

// TargetMeta returns the file metadata for a file path in the role subtree,
// if it exists. It also returns the role in that subtree in which the target
// was found. If the path doesn't exist in that role subtree, returns
// nil and an empty string.
func TestTargetMeta(t *testing.T) {
	kdb, repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, nil, kdb, localStorage)

	delegations := []string{
		"targets/level1",
		"targets/level1/a",
		"targets/level1/a/i",
	}

	k, err := cs.Create("", data.ED25519Key)
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	for i, r := range delegations {
		// create role
		role, err := data.NewRole(r, 1, []string{k.ID()}, []string{""}, nil)
		assert.NoError(t, err)

		// add role to repo
		repo.UpdateDelegations(role, []data.PublicKey{k})
		repo.InitTargets(r)

		// add a target to the role
		_, err = repo.AddTargets(r, data.Files{strconv.Itoa(i): f})
		assert.NoError(t, err)
	}

	// returns the right level
	fileMeta, role := client.TargetMeta("targets", "1")
	assert.Equal(t, &f, fileMeta)
	assert.Equal(t, "targets/level1/a", role)

	// looks only in subtree
	fileMeta, role = client.TargetMeta("targets/level1/a", "0")
	assert.Nil(t, fileMeta)
	assert.Equal(t, "", role)

	fileMeta, role = client.TargetMeta("targets/level1/a", "2")
	assert.Equal(t, &f, fileMeta)
	assert.Equal(t, "targets/level1/a/i", role)
}

// If there is no local cache and also no remote timestamp, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoTimestamps(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	err = client.downloadTimestamp()
	assert.Error(t, err)
	notFoundErr, ok := err.(store.ErrMetaNotFound)
	assert.True(t, ok)
	assert.Equal(t, data.CanonicalTimestampRole, notFoundErr.Resource)
}

// If there is no local cache and the remote timestamp is empty, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoLocalTimestampRemoteTimestampEmpty(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: []byte{}}, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	err = client.downloadTimestamp()
	assert.Error(t, err)
	assert.IsType(t, &json.SyntaxError{}, err)
}

// If there is no local cache and the remote timestamp is invalid, downloading the timestamp
// fails with a store.ErrMetaNotFound
func TestDownloadTimestampNoLocalTimestampRemoteTimestampInvalid(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(nil, nil)

	// add a timestamp to the remote cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	tsSigned.Signatures[0].Signature = []byte("12345") // invalidate the signature
	ts, err := json.Marshal(tsSigned)
	assert.NoError(t, err)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts}, nil)

	client := NewClient(repo, remoteStorage, kdb, localStorage)
	err = client.downloadTimestamp()
	assert.Error(t, err)
	assert.IsType(t, signed.ErrRoleThreshold{}, err)
}

// If there is is a local cache and no remote timestamp, we fall back on the cached timestamp
func TestDownloadTimestampLocalTimestampNoRemoteTimestamp(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	// add a timestamp to the local cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	ts, err := json.Marshal(tsSigned)
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts}, nil)

	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	err = client.downloadTimestamp()
	assert.NoError(t, err)
}

// If there is is a local cache and the remote timestamp is invalid, we fall back on the cached timestamp
func TestDownloadTimestampLocalTimestampInvalidRemoteTimestamp(t *testing.T) {
	kdb, repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	// add a timestamp to the local cache
	tsSigned, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	ts, err := json.Marshal(tsSigned)
	assert.NoError(t, err)
	localStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts}, nil)

	// add a timestamp to the remote cache
	tsSigned.Signatures[0].Signature = []byte("12345") // invalidate the signature
	ts, err = json.Marshal(tsSigned)
	assert.NoError(t, err)
	remoteStorage := store.NewMemoryStore(map[string][]byte{data.CanonicalTimestampRole: ts}, nil)

	client := NewClient(repo, remoteStorage, kdb, localStorage)
	err = client.downloadTimestamp()
	assert.NoError(t, err)
}
