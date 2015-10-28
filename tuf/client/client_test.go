package client

import (
	"crypto/sha256"
	"encoding/json"
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("targets", orig)
	assert.NoError(t, err)

	// call repo.SignSnapshot to update the targets role in the snapshot
	repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)

	err = client.downloadTargets("targets")
	assert.NoError(t, err)
}

func TestDownloadTargetChecksumMismatch(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample targets
	signedOrig, err := repo.SignTargets("targets", data.DefaultExpires("targets"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample root
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample root, snapshot, and timestamp
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"), nil)
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("root", orig)
	assert.NoError(t, err)

	// sign snapshot to make root meta in snapshot get updated
	signedOrig, err = repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)

	err = client.downloadRoot()
	assert.NoError(t, err)
}

func TestUpdateDownloadRootBadChecksum(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// sign snapshot to make sure we have a checksum for root
	_, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)

	// create and "upload" sample root, snapshot, and timestamp
	signedOrig, err := repo.SignRoot(data.DefaultExpires("root"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample timestamp
	signedOrig, err := repo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	assert.NoError(t, err)

	err = client.downloadTimestamp()
	assert.NoError(t, err)
}

func TestDownloadSnapshotHappy(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	signedOrig, err = repo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)
	orig, err = json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("timestamp", orig)
	assert.NoError(t, err)

	err = client.downloadSnapshot()
	assert.NoError(t, err)
}

// TestDownloadSnapshotNoChecksum: It should never be valid to download a
// snapshot if we don't have a checksum
func TestDownloadSnapshotNoTimestamp(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
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
	kdb, repo, _ := testutils.EmptyRepo()
	localStorage := store.NewMemoryStore(nil, nil)
	remoteStorage := store.NewMemoryStore(nil, nil)
	client := NewClient(repo, remoteStorage, kdb, localStorage)

	// sign timestamp to ensure it has a checksum for snapshot
	_, err := repo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)

	// create and "upload" sample snapshot and timestamp
	signedOrig, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)
	orig, err := json.Marshal(signedOrig)
	assert.NoError(t, err)
	err = remoteStorage.SetMeta("snapshot", orig)
	assert.NoError(t, err)

	// by not signing timestamp again we ensure it has the wrong checksum

	err = client.downloadSnapshot()
	assert.IsType(t, ErrChecksumMismatch{}, err)
}
