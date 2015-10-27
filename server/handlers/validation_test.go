package handlers

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/docker/notary/server/storage"
)

func TestValidateEmptyNew(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoNewRoot(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoNewTargets(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOnlySnapshot(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOldRoot(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateRootRotation(t *testing.T) {
	_, repo, crypto := testutils.EmptyRepo()
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
	repo.Root.Signed.Keys[rootKey.ID()] = data.NewPrivateKey(rootKey.Algorithm(), rootKey.Public(), nil)

	r, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole), nil)
	assert.NoError(t, err)
	err = signed.Sign(crypto, r, rootKey, oldRootKey)
	assert.NoError(t, err)

	rt, err := data.RootFromSigned(r)
	assert.NoError(t, err)
	repo.SetRoot(rt)

	sn, err = repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole), nil)
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

	err = validateUpdate("testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoRoot(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrValidation{}, err)
}

func TestValidateSnapshotMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadHierarchy{}, err)
}

// ### Role missing negative tests ###
// These tests remove a role from the Root file and
// check for a ErrBadRoot
func TestValidateRootRoleMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsRoleMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateSnapshotRoleMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

// ### End role missing negative tests ###

// ### Signature missing negative tests ###
func TestValidateRootSigMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsSigMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadTargets{}, err)
}

func TestValidateSnapshotSigMissing(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End signature missing negative tests ###

// ### Corrupted metadata negative tests ###
func TestValidateRootCorrupt(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsCorrupt(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadTargets{}, err)
}

func TestValidateSnapshotCorrupt(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End corrupted metadata negative tests ###

// ### Snapshot size mismatch negative tests ###
func TestValidateRootModifiedSize(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadRoot{}, err)
}

func TestValidateTargetsModifiedSize(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End snapshot size mismatch negative tests ###

// ### Snapshot hash mismatch negative tests ###
func TestValidateRootModifiedHash(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

func TestValidateTargetsModifiedHash(t *testing.T) {
	_, repo, _ := testutils.EmptyRepo()
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

	err = validateUpdate("testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, ErrBadSnapshot{}, err)
}

// ### End snapshot hash mismatch negative tests ###
