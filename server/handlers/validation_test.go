package handlers

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/server/storage"
)

func copyTimestampKey(t *testing.T, fromKeyDB *keys.KeyDB,
	toStore storage.MetaStore, gun string) {

	role := fromKeyDB.GetRole(data.CanonicalTimestampRole)
	assert.NotNil(t, role, "No timestamp role in the KeyDB")
	assert.Len(t, role.KeyIDs, 1, fmt.Sprintf(
		"Expected 1 timestamp key in timestamp role, got %d", len(role.KeyIDs)))

	pubTimestampKey := fromKeyDB.GetKey(role.KeyIDs[0])
	assert.NotNil(t, pubTimestampKey,
		"Timestamp key specified by KeyDB role not in KeysDB")

	err := toStore.SetTimestampKey(gun, pubTimestampKey.Algorithm(),
		pubTimestampKey.Public())
	assert.NoError(t, err)
}

func TestValidateEmptyNew(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}
	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoNewRoot(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
	)

	updates := []storage.MetaUpdate{
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}
	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoNewTargets(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
	)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}
	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOnlySnapshot(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
	)
	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
	)

	updates := []storage.MetaUpdate{
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
	}
	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOldRoot(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
	)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}
	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateRootRotation(t *testing.T) {
	kdb, repo, crypto := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent(
		"testGUN",
		storage.MetaUpdate{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
	)

	oldRootRole := repo.Root.Signed.Roles["root"]
	oldRootKey := repo.Root.Signed.Keys[oldRootRole.KeyIDs[0]]

	rootKey, err := crypto.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	assert.NoError(t, err)

	delete(repo.Root.Signed.Keys, oldRootRole.KeyIDs[0])

	repo.Root.Signed.Roles["root"] = &rootRole.RootRole
	repo.Root.Signed.Keys[rootKey.ID()] = rootKey

	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	assert.NoError(t, err)
	err = signed.Sign(crypto, r, rootKey, oldRootKey)
	assert.NoError(t, err)

	rt, err := data.RootFromSigned(r)
	assert.NoError(t, err)
	repo.SetRoot(rt)

	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err = testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoRoot(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	_, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrValidation{}, err)
}

func TestValidateSnapshotMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, _, _, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadHierarchy{}, err)
}

// If there is no timestamp key in the store, validation fails.  This could
// happen if pushing an existing repository from one server to another that
// does not have the repo.
func TestValidateRootNoTimestampKey(t *testing.T) {
	_, oldRepo, _ := testutils.EmptyRepo()

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
	}

	// sanity check - no timestamp keys for the GUN
	_, _, err = store.GetTimestampKey("testGUN")
	assert.Error(t, err)
	assert.IsType(t, &storage.ErrNoKey{}, err)

	// do not copy the targets key to the storage, and try to update the root
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)

	// there should still be no timestamp keys - one should not have been
	// created
	_, _, err = store.GetTimestampKey("testGUN")
	assert.Error(t, err)
}

// If the timestamp key in the store does not match the timestamp key in
// the root.json, validation fails.  This could happen if pushing an existing
// repository from one server to another that had already initialized the same
// repo.
func TestValidateRootInvalidTimestampKey(t *testing.T) {
	_, oldRepo, _ := testutils.EmptyRepo()

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
	}

	key, err := trustmanager.GenerateECDSAKey(rand.Reader)
	assert.NoError(t, err)
	err = store.SetTimestampKey("testGUN", key.Algorithm(), key.Public())
	assert.NoError(t, err)

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

// ### Role missing negative tests ###
// These tests remove a role from the Root file and
// check for a ErrBadRoot
func TestValidateRootRoleMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "root")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsRoleMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "targets")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateSnapshotRoleMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

// ### End role missing negative tests ###

// ### Signature missing negative tests ###
func TestValidateRootSigMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	r.Signatures = nil

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsSigMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	tg.Signatures = nil

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadTargets{}, err)
}

func TestValidateSnapshotSigMissing(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	sn.Signatures = nil

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End signature missing negative tests ###

// ### Corrupted metadata negative tests ###
func TestValidateRootCorrupt(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	root[0] = root[0] ^ 0xff

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsCorrupt(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	targets[0] = targets[0] ^ 0xff

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadTargets{}, err)
}

func TestValidateSnapshotCorrupt(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	snapshot[0] = snapshot[0] ^ 0xff

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End corrupted metadata negative tests ###

// ### Snapshot size mismatch negative tests ###
func TestValidateRootModifiedSize(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	// add another copy of the signature so the hash is different
	r.Signatures = append(r.Signatures, r.Signatures[0])

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	root[0] = root[0] ^ 0xff

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsModifiedSize(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	// add another copy of the signature so the hash is different
	tg.Signatures = append(tg.Signatures, tg.Signatures[0])

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End snapshot size mismatch negative tests ###

// ### Snapshot hash mismatch negative tests ###
func TestValidateRootModifiedHash(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	assert.NoError(t, err)
	snap.Signed.Meta["root"].Hashes["sha256"][0] = snap.Signed.Meta["root"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	assert.NoError(t, err)

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

func TestValidateTargetsModifiedHash(t *testing.T) {
	kdb, repo, _ := testutils.EmptyRepo()
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	assert.NoError(t, err)
	snap.Signed.Meta["targets"].Hashes["sha256"][0] = snap.Signed.Meta["targets"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	assert.NoError(t, err)

	root, targets, snapshot, timestamp, err := testutils.Serialize(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{
		{
			Role:    "root",
			Version: 1,
			Data:    root,
		},
		{
			Role:    "targets",
			Version: 1,
			Data:    targets,
		},
		{
			Role:    "snapshot",
			Version: 1,
			Data:    snapshot,
		},
		{
			Role:    "timestamp",
			Version: 1,
			Data:    timestamp,
		},
	}

	copyTimestampKey(t, kdb, store, "testGUN")
	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End snapshot hash mismatch negative tests ###
