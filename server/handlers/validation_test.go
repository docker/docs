package handlers

import (
	"bytes"
	"encoding/json"
	"path"
	"reflect"
	"testing"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/testutils"
	"github.com/docker/notary/tuf/validation"
	"github.com/stretchr/testify/assert"
)

// this is a fake storage that serves errors
type getFailStore struct {
	errsToReturn map[string]error
	storage.MetaStore
}

// GetCurrent returns the current metadata, or an error depending on whether
// getFailStore is configured to return an error for this role
func (f getFailStore) GetCurrent(gun, tufRole string) ([]byte, error) {
	err := f.errsToReturn[tufRole]
	if err == nil {
		return f.MetaStore.GetCurrent(gun, tufRole)
	}
	return nil, err
}

// GetChecksum returns the metadata with this checksum, or an error depending on
// whether getFailStore is configured to return an error for this role
func (f getFailStore) GetChecksum(gun, tufRole, checksum string) ([]byte, error) {
	err := f.errsToReturn[tufRole]
	if err == nil {
		return f.MetaStore.GetChecksum(gun, tufRole, checksum)
	}
	return nil, err
}

func copyKeys(t *testing.T, from signed.CryptoService, roles ...string) signed.CryptoService {
	memKeyStore := trustmanager.NewKeyMemoryStore(passphrase.ConstantRetriever("pass"))
	for _, role := range roles {
		for _, keyID := range from.ListKeys(role) {
			key, _, err := from.GetPrivateKey(keyID)
			assert.NoError(t, err)
			memKeyStore.AddKey(path.Base(keyID), role, key)
		}
	}
	return cryptoservice.NewCryptoService("", memKeyStore)
}

// Returns a mapping of role name to `MetaUpdate` objects
func getUpdates(r, tg, sn, ts *data.Signed) (
	root, targets, snapshot, timestamp storage.MetaUpdate, err error) {

	rs, tgs, sns, tss, err := testutils.Serialize(r, tg, sn, ts)
	if err != nil {
		return
	}

	root = storage.MetaUpdate{
		Role:    data.CanonicalRootRole,
		Version: 1,
		Data:    rs,
	}
	targets = storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgs,
	}
	snapshot = storage.MetaUpdate{
		Role:    data.CanonicalSnapshotRole,
		Version: 1,
		Data:    sns,
	}
	timestamp = storage.MetaUpdate{
		Role:    data.CanonicalTimestampRole,
		Version: 1,
		Data:    tss,
	}
	return
}

func TestValidateEmptyNew(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)

	// we generated our own timestamp, and did not take the other timestamp,
	// but all other metadata should come from updates
	founds := make(map[string]bool)
	for _, update := range updates {
		switch update.Role {
		case data.CanonicalRootRole:
			assert.True(t, bytes.Equal(update.Data, root.Data))
			founds[data.CanonicalRootRole] = true
		case data.CanonicalSnapshotRole:
			assert.True(t, bytes.Equal(update.Data, snapshot.Data))
			founds[data.CanonicalSnapshotRole] = true
		case data.CanonicalTargetsRole:
			assert.True(t, bytes.Equal(update.Data, targets.Data))
			founds[data.CanonicalTargetsRole] = true
		case data.CanonicalTimestampRole:
			assert.False(t, bytes.Equal(update.Data, timestamp.Data))
			founds[data.CanonicalTimestampRole] = true
		}
	}
	assert.Len(t, founds, 4)
}

func TestValidatePrevTimestamp(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	store := storage.NewMemStorage()
	store.UpdateCurrent("testGUN", timestamp)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)

	// we generated our own timestamp, and did not take the other timestamp,
	// but all other metadata should come from updates
	var foundTimestamp bool
	for _, update := range updates {
		if update.Role == data.CanonicalTimestampRole {
			foundTimestamp = true
			oldTimestamp, newTimestamp := &data.SignedTimestamp{}, &data.SignedTimestamp{}
			assert.NoError(t, json.Unmarshal(timestamp.Data, oldTimestamp))
			assert.NoError(t, json.Unmarshal(update.Data, newTimestamp))
			assert.Equal(t, oldTimestamp.Signed.Version+1, newTimestamp.Signed.Version)
		}
	}
	assert.True(t, foundTimestamp)
}

func TestValidatePreviousTimestampCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	// have corrupt timestamp data in the storage already
	store := storage.NewMemStorage()
	timestamp.Data = timestamp.Data[1:]
	store.UpdateCurrent("testGUN", timestamp)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, &json.SyntaxError{}, err)
}

func TestValidateGetCurrentTimestampBroken(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot}

	store := getFailStore{
		MetaStore:    storage.NewMemStorage(),
		errsToReturn: map[string]error{data.CanonicalTimestampRole: data.ErrNoSuchRole{}},
	}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, data.ErrNoSuchRole{}, err)
}

func TestValidateNoNewRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	updates := []storage.MetaUpdate{targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoNewTargets(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent("testGUN", targets)
	updates := []storage.MetaUpdate{root, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOnlySnapshot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	store.UpdateCurrent("testGUN", targets)

	updates := []storage.MetaUpdate{snapshot}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateOldRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent("testGUN", root)
	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateRootRotation(t *testing.T) {
	repo, crypto, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store.UpdateCurrent("testGUN", root)

	oldRootRole := repo.Root.Signed.Roles["root"]
	oldRootKey := repo.Root.Signed.Keys[oldRootRole.KeyIDs[0]]

	rootKey, err := crypto.Create("root", "", data.ED25519Key)
	assert.NoError(t, err)
	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil)
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
	root, targets, snapshot, timestamp, err = getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, crypto, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateNoRoot(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	_, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrValidation{}, err)
}

func TestValidateSnapshotMissingNoSnapshotKey(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadHierarchy{}, err)
}

func TestValidateSnapshotGenerateNoPrev(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

func TestValidateSnapshotGenerateWithPrev(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	// set the current snapshot in the store manually so we find it when generating
	// the next version
	store.UpdateCurrent("testGUN", snapshot)

	prev, err := data.SnapshotFromSigned(sn)
	assert.NoError(t, err)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)

	for _, u := range updates {
		if u.Role == data.CanonicalSnapshotRole {
			curr := &data.SignedSnapshot{}
			err = json.Unmarshal(u.Data, curr)
			assert.Equal(t, prev.Signed.Version+1, curr.Signed.Version)
			assert.Equal(t, u.Version, curr.Signed.Version)
		}
	}
}

func TestValidateSnapshotGeneratePrevCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	// corrupt the JSON structure of prev snapshot
	snapshot.Data = snapshot.Data[1:]
	// set the current snapshot in the store manually so we find it when generating
	// the next version
	store.UpdateCurrent("testGUN", snapshot)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, &json.SyntaxError{}, err)
}

// Store is broken when getting the current snapshot
func TestValidateSnapshotGenerateStoreGetCurrentSnapshotBroken(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := getFailStore{
		MetaStore:    storage.NewMemStorage(),
		errsToReturn: map[string]error{data.CanonicalSnapshotRole: data.ErrNoSuchRole{}},
	}
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, data.ErrNoSuchRole{}, err)
}

func TestValidateSnapshotGenerateNoTargets(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, _, _, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
}

func TestValidateSnapshotGenerate(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()
	snapRole, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	assert.NoError(t, err)

	for _, k := range snapRole.Keys {
		err := store.SetKey("testGUN", data.CanonicalSnapshotRole, k.Algorithm(), k.Public())
		assert.NoError(t, err)
	}

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, _, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{targets}

	store.UpdateCurrent("testGUN", root)

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole, data.CanonicalSnapshotRole)
	updates, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.NoError(t, err)
}

// If there is no timestamp key in the store, validation fails.  This could
// happen if pushing an existing repository from one server to another that
// does not have the repo.
func TestValidateRootNoTimestampKey(t *testing.T) {
	oldRepo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	// do not copy the targets key to the storage, and try to update the root
	serverCrypto := signed.NewEd25519()
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)

	// there should still be no timestamp keys - one should not have been
	// created
	assert.Empty(t, serverCrypto.ListAllKeys())
}

// If the timestamp key in the store does not match the timestamp key in
// the root.json, validation fails.  This could happen if pushing an existing
// repository from one server to another that had already initialized the same
// repo.
func TestValidateRootInvalidTimestampKey(t *testing.T) {
	oldRepo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	serverCrypto := signed.NewEd25519()
	_, err = serverCrypto.Create(data.CanonicalTimestampRole, data.ED25519Key)
	assert.NoError(t, err)

	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

// If the timestamp role has a threshold > 1, validation fails.
func TestValidateRootInvalidTimestampThreshold(t *testing.T) {
	oldRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	tsRole, ok := oldRepo.Root.Signed.Roles[data.CanonicalTimestampRole]
	assert.True(t, ok)
	tsRole.Threshold = 2

	r, tg, sn, ts, err := testutils.Sign(oldRepo)
	assert.NoError(t, err)
	root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	store := storage.NewMemStorage()
	updates := []storage.MetaUpdate{root, targets, snapshot}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timestamp role has invalid threshold")
}

// If any role has a threshold < 1, validation fails
func TestValidateRootInvalidZeroThreshold(t *testing.T) {
	for _, role := range data.BaseRoles {
		oldRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
		assert.NoError(t, err)
		tsRole, ok := oldRepo.Root.Signed.Roles[role]
		assert.True(t, ok)
		tsRole.Threshold = 0

		r, tg, sn, ts, err := testutils.Sign(oldRepo)
		assert.NoError(t, err)
		root, targets, snapshot, _, err := getUpdates(r, tg, sn, ts)
		assert.NoError(t, err)

		store := storage.NewMemStorage()
		updates := []storage.MetaUpdate{root, targets, snapshot}

		serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
		_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid threshold")
	}
}

// ### Role missing negative tests ###
// These tests remove a role from the Root file and
// check for a validation.ErrBadRoot
func TestValidateRootRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "root")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "targets")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateSnapshotRoleMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

// ### End role missing negative tests ###

// ### Signature missing negative tests ###
func TestValidateRootSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	delete(repo.Root.Signed.Roles, "snapshot")

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	r.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	tg.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadTargets{}, err)
}

func TestValidateSnapshotSigMissing(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	sn.Signatures = nil

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End signature missing negative tests ###

// ### Corrupted metadata negative tests ###
func TestValidateRootCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	root.Data[0] = root.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	targets.Data[0] = targets.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadTargets{}, err)
}

func TestValidateSnapshotCorrupt(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)
	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	snapshot.Data[0] = snapshot.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End corrupted metadata negative tests ###

// ### Snapshot size mismatch negative tests ###
func TestValidateRootModifiedSize(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	// add another copy of the signature so the hash is different
	r.Signatures = append(r.Signatures, r.Signatures[0])

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	// flip all the bits in the first byte
	root.Data[0] = root.Data[0] ^ 0xff

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadRoot{}, err)
}

func TestValidateTargetsModifiedSize(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	// add another copy of the signature so the hash is different
	tg.Signatures = append(tg.Signatures, tg.Signatures[0])

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End snapshot size mismatch negative tests ###

// ### Snapshot hash mismatch negative tests ###
func TestValidateRootModifiedHash(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	assert.NoError(t, err)
	snap.Signed.Meta["root"].Hashes["sha256"][0] = snap.Signed.Meta["root"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	assert.NoError(t, err)

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadSnapshot{}, err)
}

func TestValidateTargetsModifiedHash(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	r, tg, sn, ts, err := testutils.Sign(repo)
	assert.NoError(t, err)

	snap, err := data.SnapshotFromSigned(sn)
	assert.NoError(t, err)
	snap.Signed.Meta["targets"].Hashes["sha256"][0] = snap.Signed.Meta["targets"].Hashes["sha256"][0] ^ 0xff

	sn, err = snap.ToSigned()
	assert.NoError(t, err)

	root, targets, snapshot, timestamp, err := getUpdates(r, tg, sn, ts)
	assert.NoError(t, err)

	updates := []storage.MetaUpdate{root, targets, snapshot, timestamp}

	serverCrypto := copyKeys(t, cs, data.CanonicalTimestampRole)
	_, err = validateUpdate(serverCrypto, "testGUN", updates, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadSnapshot{}, err)
}

// ### End snapshot hash mismatch negative tests ###

// ### generateSnapshot tests ###
func TestGenerateSnapshotRootNotLoaded(t *testing.T) {
	repo := tuf.NewRepo(nil)
	_, err := generateSnapshot("gun", repo, storage.NewMemStorage())
	assert.Error(t, err)
	assert.IsType(t, validation.ErrValidation{}, err)
}

func TestGenerateSnapshotNoKey(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	// delete snapshot key in the cryptoservice
	for _, keyID := range cs.ListKeys(data.CanonicalSnapshotRole) {
		assert.NoError(t, cs.RemoveKey(keyID))
	}

	_, err = generateSnapshot("gun", repo, store)
	assert.Error(t, err)
	assert.IsType(t, validation.ErrBadHierarchy{}, err)
}

// ### End generateSnapshot tests ###

// ### Target validation with delegations tests
func TestLoadTargetsFromStore(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	st, err := repo.SignTargets(
		data.CanonicalTargetsRole,
		data.DefaultExpires(data.CanonicalTargetsRole),
	)
	assert.NoError(t, err)

	tgs, err := json.Marshal(st)
	assert.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgs,
	}
	store.UpdateCurrent("gun", update)

	generated := repo.Targets[data.CanonicalTargetsRole]
	delete(repo.Targets, data.CanonicalTargetsRole)
	_, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.False(t, ok)

	err = loadTargetsFromStore("gun", data.CanonicalTargetsRole, repo, store)
	assert.NoError(t, err)
	loaded, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.True(t, reflect.DeepEqual(generated.Signatures, loaded.Signatures))
	assert.Len(t, loaded.Signed.Targets, 0)
	assert.Equal(t, len(generated.Signed.Targets), len(loaded.Signed.Targets))
	assert.Len(t, loaded.Signed.Delegations.Roles, 0)
	assert.Equal(t, len(generated.Signed.Delegations.Roles), len(loaded.Signed.Delegations.Roles))
	assert.Len(t, loaded.Signed.Delegations.Keys, 0)
	assert.Equal(t, len(generated.Signed.Delegations.Keys), len(loaded.Signed.Delegations.Keys))
	assert.True(t, generated.Signed.Expires.Equal(loaded.Signed.Expires))
	assert.Equal(t, generated.Signed.Type, loaded.Signed.Type)
	assert.Equal(t, generated.Signed.Version, loaded.Signed.Version)
}

func TestValidateTargetsLoadParent(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	assert.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	assert.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	assert.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	// we're not going to validate things loaded from storage, so no need
	// to sign the base targets, just Marshal it and set it into storage
	tgtsJSON, err := json.Marshal(baseRepo.Targets["targets"])
	assert.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	// generate the update object we're doing to use to call loadAndValidateTargets
	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	assert.NoError(t, err)
	delJSON, err := json.Marshal(del)
	assert.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{"targets/level1": delUpdate}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	assert.NoError(t, err)
	assert.Len(t, updates, 1)
	assert.Equal(t, "targets/level1", updates[0].Role)
	assert.Equal(t, delJSON, updates[0].Data)
}

func TestValidateTargetsParentInUpdate(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	assert.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	assert.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	assert.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	targets, err := baseRepo.SignTargets("targets", data.DefaultExpires(data.CanonicalTargetsRole))

	tgtsJSON, err := json.Marshal(targets)
	assert.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	assert.NoError(t, err)
	delJSON, err := json.Marshal(del)
	assert.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{
		"targets/level1": delUpdate,
		"targets":        update,
	}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	// because we sort the roles, the list of returned updates
	// will contain shallower roles first, in this case "targets",
	// and then "targets/level1"
	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	assert.NoError(t, err)
	assert.Len(t, updates, 2)
	assert.Equal(t, "targets", updates[0].Role)
	assert.Equal(t, tgtsJSON, updates[0].Data)
	assert.Equal(t, "targets/level1", updates[1].Role)
	assert.Equal(t, delJSON, updates[1].Data)
}

func TestValidateTargetsParentNotFound(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	k, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	assert.NoError(t, err)

	err = baseRepo.UpdateDelegationKeys("targets/level1", []data.PublicKey{k}, []string{}, 1)
	assert.NoError(t, err)
	err = baseRepo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	assert.NoError(t, err)

	// no targets file is created for the new delegations, so force one
	baseRepo.InitTargets("targets/level1")

	// generate the update object we're doing to use to call loadAndValidateTargets
	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	assert.NoError(t, err)
	delJSON, err := json.Marshal(del)
	assert.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	roles := map[string]storage.MetaUpdate{"targets/level1": delUpdate}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	_, err = loadAndValidateTargets("gun", valRepo, roles, store)
	assert.Error(t, err)
	assert.IsType(t, storage.ErrNotFound{}, err)
}

func TestValidateTargetsRoleNotInParent(t *testing.T) {
	baseRepo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	store := storage.NewMemStorage()

	level1Key, err := cs.Create("targets/level1", "docker.com/notary", data.ED25519Key)
	assert.NoError(t, err)
	r, err := data.NewRole("targets/level1", 1, []string{level1Key.ID()}, []string{""})

	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Roles = []*data.Role{r}
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Keys = data.Keys{
		level1Key.ID(): level1Key,
	}

	baseRepo.InitTargets("targets/level1")

	del, err := baseRepo.SignTargets("targets/level1", data.DefaultExpires(data.CanonicalTargetsRole))
	assert.NoError(t, err)
	delJSON, err := json.Marshal(del)
	assert.NoError(t, err)

	delUpdate := storage.MetaUpdate{
		Role:    "targets/level1",
		Version: 1,
		Data:    delJSON,
	}

	// set back to empty so stored targets doesn't have reference to level1
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Roles = nil
	baseRepo.Targets[data.CanonicalTargetsRole].Signed.Delegations.Keys = nil
	targets, err := baseRepo.SignTargets(data.CanonicalTargetsRole, data.DefaultExpires(data.CanonicalTargetsRole))

	tgtsJSON, err := json.Marshal(targets)
	assert.NoError(t, err)
	update := storage.MetaUpdate{
		Role:    data.CanonicalTargetsRole,
		Version: 1,
		Data:    tgtsJSON,
	}
	store.UpdateCurrent("gun", update)

	roles := map[string]storage.MetaUpdate{
		"targets/level1":          delUpdate,
		data.CanonicalTargetsRole: update,
	}

	valRepo := tuf.NewRepo(nil)
	valRepo.SetRoot(baseRepo.Root)

	// because we sort the roles, the list of returned updates
	// will contain shallower roles first, in this case "targets",
	// and then "targets/level1"
	updates, err := loadAndValidateTargets("gun", valRepo, roles, store)
	assert.NoError(t, err)
	assert.Len(t, updates, 1)
	assert.Equal(t, data.CanonicalTargetsRole, updates[0].Role)
	assert.Equal(t, tgtsJSON, updates[0].Data)
}

// ### End target validation with delegations tests
