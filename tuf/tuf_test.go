package tuf

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"
	"github.com/stretchr/testify/assert"
)

func initRepo(t *testing.T, cryptoService signed.CryptoService, keyDB *keys.KeyDB) *Repo {
	rootKey, err := cryptoService.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	targetsKey, err := cryptoService.Create("targets", data.ED25519Key)
	assert.NoError(t, err)
	snapshotKey, err := cryptoService.Create("snapshot", data.ED25519Key)
	assert.NoError(t, err)
	timestampKey, err := cryptoService.Create("timestamp", data.ED25519Key)
	assert.NoError(t, err)

	keyDB.AddKey(rootKey)
	keyDB.AddKey(targetsKey)
	keyDB.AddKey(snapshotKey)
	keyDB.AddKey(timestampKey)

	rootRole := &data.Role{
		Name: "root",
		RootRole: data.RootRole{
			KeyIDs:    []string{rootKey.ID()},
			Threshold: 1,
		},
	}
	targetsRole := &data.Role{
		Name: "targets",
		RootRole: data.RootRole{
			KeyIDs:    []string{targetsKey.ID()},
			Threshold: 1,
		},
	}
	snapshotRole := &data.Role{
		Name: "snapshot",
		RootRole: data.RootRole{
			KeyIDs:    []string{snapshotKey.ID()},
			Threshold: 1,
		},
	}
	timestampRole := &data.Role{
		Name: "timestamp",
		RootRole: data.RootRole{
			KeyIDs:    []string{timestampKey.ID()},
			Threshold: 1,
		},
	}

	keyDB.AddRole(rootRole)
	keyDB.AddRole(targetsRole)
	keyDB.AddRole(snapshotRole)
	keyDB.AddRole(timestampRole)

	repo := NewRepo(keyDB, cryptoService)
	err = repo.InitRepo(false)
	assert.NoError(t, err)
	return repo
}

// we require that at least the base targets role is available when creating
// initializing a snapshot
func TestInitSnapshotNoTargets(t *testing.T) {
	cryptoService := signed.NewEd25519()
	keyDB := keys.NewDB()
	rootKey, err := cryptoService.Create("root", data.ED25519Key)
	assert.NoError(t, err)
	snapshotKey, err := cryptoService.Create("snapshot", data.ED25519Key)
	assert.NoError(t, err)

	keyDB.AddKey(rootKey)
	keyDB.AddKey(snapshotKey)

	rootRole := &data.Role{
		Name: "root",
		RootRole: data.RootRole{
			KeyIDs:    []string{rootKey.ID()},
			Threshold: 1,
		},
	}
	snapshotRole := &data.Role{
		Name: "snapshot",
		RootRole: data.RootRole{
			KeyIDs:    []string{snapshotKey.ID()},
			Threshold: 1,
		},
	}

	keyDB.AddRole(rootRole)
	keyDB.AddRole(snapshotRole)

	repo := NewRepo(keyDB, cryptoService)
	err = repo.InitSnapshot()
	assert.Error(t, err)
	assert.IsType(t, ErrNotLoaded{}, err)
}

func writeRepo(t *testing.T, dir string, repo *Repo) {
	err := os.MkdirAll(dir, 0755)
	assert.NoError(t, err)
	signedRoot, err := repo.SignRoot(data.DefaultExpires("root"))
	assert.NoError(t, err)
	rootJSON, _ := json.Marshal(signedRoot)
	ioutil.WriteFile(dir+"/root.json", rootJSON, 0755)

	for r := range repo.Targets {
		signedTargets, err := repo.SignTargets(r, data.DefaultExpires("targets"))
		assert.NoError(t, err)
		targetsJSON, _ := json.Marshal(signedTargets)
		p := path.Join(dir, r+".json")
		parentDir := filepath.Dir(p)
		os.MkdirAll(parentDir, 0755)
		ioutil.WriteFile(p, targetsJSON, 0755)
	}

	signedSnapshot, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	assert.NoError(t, err)
	snapshotJSON, _ := json.Marshal(signedSnapshot)
	ioutil.WriteFile(dir+"/snapshot.json", snapshotJSON, 0755)

	signedTimestamp, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	assert.NoError(t, err)
	timestampJSON, _ := json.Marshal(signedTimestamp)
	ioutil.WriteFile(dir+"/timestamp.json", timestampJSON, 0755)
}

func TestInitRepo(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)
	writeRepo(t, "/tmp/tufrepo", repo)
}

func TestUpdateDelegations(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	assert.False(t, ok, "no empty targets file should be created for deepest delegation")

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	testDeepKey, err := ed25519.Create("targets/test/deep", data.ED25519Key)
	assert.NoError(t, err)
	roleDeep, err := data.NewRole("targets/test/deep", 1, []string{testDeepKey.ID()}, []string{"test/deep"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(roleDeep, data.KeyList{testDeepKey})
	assert.NoError(t, err)

	// this metadata didn't exist before, but creating targets/test/deep created
	// the targets/test metadata
	r, ok = repo.Targets["targets/test"]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testDeepKey.ID(), keyIDs[0])
	assert.True(t, r.Dirty)

	// no empty delegation metadata is created for targets/test/deep
	_, ok = repo.Targets["targets/test/deep"]
	assert.False(t, ok, "no empty targets file should be created for deepest delegation")
}

func TestUpdateDelegationsParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testDeepKey, err := ed25519.Create("targets/test/deep", data.ED25519Key)
	assert.NoError(t, err)
	roleDeep, err := data.NewRole("targets/test/deep", 1, []string{testDeepKey.ID()}, []string{"test/deep"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(roleDeep, data.KeyList{testDeepKey})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 0)

	// no delegation metadata created for non-existent parent
	_, ok = repo.Targets["targets/test"]
	assert.False(t, ok, "no targets file should be created for nonexistent parent delegation")
}

// Updating delegations needs to modify the parent of the role being updated.
// If there is no signing key for that parent, the delegation cannot be added.
func TestUpdateDelegationsMissingParentKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// remove the target key (all keys)
	repo.cryptoService = signed.NewEd25519()

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole("targets/role", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{roleKey})
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/role"]
	assert.False(t, ok, "no targets file should be created for empty delegation")
}

func TestUpdateDelegationsInvalidRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	// data.NewRole errors if the role isn't a valid TUF role so use one of the non-delegation
	// valid roles
	invalidRole, err := data.NewRole("root", 1, []string{roleKey.ID()}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(invalidRole, data.KeyList{roleKey})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 0)

	// no delegation metadata created for invalid delgation
	_, ok = repo.Targets["root"]
	assert.False(t, ok, "no targets file should be created since delegation failed")
}

// A delegation can be created with a role that is missing a signing key, so
// long as UpdateDelegations is called with the key
func TestUpdateDelegationsRoleThatIsMissingDelegationKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole("targets/role", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	// key should get added to role as part of updating the delegation
	err = repo.UpdateDelegations(role, data.KeyList{roleKey})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, roleKey.ID(), keyIDs[0])
	assert.True(t, r.Dirty)

	// no empty delegation metadata created for new delegation
	_, ok = repo.Targets["targets/role"]
	assert.False(t, ok, "no targets file should be created for empty delegation")
}

func TestUpdateDelegationsNotEnoughKeys(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole("targets/role", 2, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{roleKey})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	// no delegation metadata created for failed delegation
	_, ok := repo.Targets["targets/role"]
	assert.False(t, ok, "no targets file should be created since delegation failed")
}

func TestUpdateDelegationsReplaceRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	// no empty delegation metadata created for new delegation
	_, ok = repo.Targets["targets/test"]
	assert.False(t, ok, "no targets file should be created for empty delegation")

	// create one now to assert that replacing the delegation doesn't delete the
	// metadata
	repo.InitTargets("targets/test")

	// create another role with the same name and ensure it replaces the
	// previous role
	testKey2, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role2, err := data.NewRole("targets/test", 1, []string{testKey2.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role2, data.KeyList{testKey2})
	assert.NoError(t, err)

	r, ok = repo.Targets["targets"]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey2.ID(), keyIDs[0])
	assert.True(t, r.Dirty)

	// delegation was not deleted
	_, ok = repo.Targets["targets/test"]
	assert.True(t, ok, "targets file should still be here")
}

func TestUpdateDelegationsAddKeyToRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	testKey2, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey2})
	assert.NoError(t, err)

	r, ok = repo.Targets["targets"]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 2)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 2)
	// it does an append so the order is deterministic (but not meaningful to TUF)
	assert.Equal(t, testKey.ID(), keyIDs[0])
	assert.Equal(t, testKey2.ID(), keyIDs[1])
	assert.True(t, r.Dirty)
}

func TestDeleteDelegations(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	// ensure that the metadata is there and snapshot is there
	targets, err := repo.InitTargets("targets/test")
	assert.NoError(t, err)
	targetsSigned, err := targets.ToSigned()
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateSnapshot("targets/test", targetsSigned))
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	assert.True(t, ok)

	assert.NoError(t, repo.DeleteDelegation(*role))
	assert.Len(t, r.Signed.Delegations.Roles, 0)
	assert.Len(t, r.Signed.Delegations.Keys, 0)
	assert.True(t, r.Dirty)

	// metadata should be deleted
	_, ok = repo.Targets["targets/test"]
	assert.False(t, ok)
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	assert.False(t, ok)
}

func TestDeleteDelegationsRoleNotExistBecauseNoParentMeta(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/test"]
	assert.False(t, ok, "no targets file should be created for empty delegation")

	delRole, err := data.NewRole(
		"targets/test/a", 1, []string{testKey.ID()}, []string{"test"}, []string{})

	err = repo.DeleteDelegation(*delRole)
	assert.NoError(t, err)
	// still no metadata
	_, ok = repo.Targets["targets/test"]
	assert.False(t, ok)
}

func TestDeleteDelegationsRoleNotExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// initRepo leaves all the roles as Dirty. Set to false
	// to test removing a non-existant role doesn't mark
	// a role as dirty
	repo.Targets[data.CanonicalTargetsRole].Dirty = false

	role, err := data.NewRole("targets/test", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*role)
	assert.NoError(t, err)
	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 0)
	assert.Len(t, r.Signed.Delegations.Keys, 0)
	assert.False(t, r.Dirty)
}

func TestDeleteDelegationsInvalidRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// data.NewRole errors if the role isn't a valid TUF role so use one of the non-delegation
	// valid roles
	invalidRole, err := data.NewRole("root", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*invalidRole)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 0)
}

func TestDeleteDelegationsParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testRole, err := data.NewRole("targets/test/deep", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*testRole)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 0)
}

// Can't delete a delegation if we don't have the parent's signing key
func TestDeleteDelegationsMissingParentSigningKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	// ensure that the metadata is there and snapshot is there
	targets, err := repo.InitTargets("targets/test")
	assert.NoError(t, err)
	targetsSigned, err := targets.ToSigned()
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateSnapshot("targets/test", targetsSigned))
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	assert.True(t, ok)

	// delete all signing keys
	repo.cryptoService = signed.NewEd25519()
	err = repo.DeleteDelegation(*role)
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)

	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	assert.True(t, r.Dirty)

	// metadata should be here still
	_, ok = repo.Targets["targets/test"]
	assert.True(t, ok)
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	assert.True(t, ok)
}

func TestDeleteDelegationsMidSliceRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	role2, err := data.NewRole("targets/test2", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role2, data.KeyList{testKey})
	assert.NoError(t, err)

	role3, err := data.NewRole("targets/test3", 1, []string{}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role3, data.KeyList{testKey})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*role2)
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	assert.Len(t, r.Signed.Delegations.Roles, 2)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	assert.True(t, r.Dirty)
}

// If the parent exists, the metadata exists, and the delegation is in it,
// returns the role that was found
func TestGetDelegationRoleAndMetadataExistDelegationExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("meh", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole(
		"targets/level1", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateDelegations(role, data.KeyList{testKey}))

	role, err = data.NewRole(
		"targets/level1/level2", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateDelegations(role, data.KeyList{testKey}))

	gottenRole, gottenKeys, err := repo.GetDelegation("targets/level1/level2")
	assert.NoError(t, err)
	assert.Equal(t, role, gottenRole)
	_, ok := gottenKeys[testKey.ID()]
	assert.True(t, ok)
}

// If the parent exists, the metadata exists, and the delegation isn't in it,
// returns an ErrNoSuchRole
func TestGetDelegationRoleAndMetadataExistDelegationDoesntExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("meh", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole(
		"targets/level1", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateDelegations(role, data.KeyList{testKey}))

	// ensure metadata exists
	repo.InitTargets("targets/level1")

	_, _, err = repo.GetDelegation("targets/level1/level2")
	assert.Error(t, err)
	assert.IsType(t, data.ErrNoSuchRole{}, err)
}

// If the parent exists but the metadata doesn't exist, returns an ErrNoSuchRole
func TestGetDelegationRoleAndMetadataDoesntExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("meh", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole(
		"targets/level1", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)
	assert.NoError(t, repo.UpdateDelegations(role, data.KeyList{testKey}))

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/test"]
	assert.False(t, ok, "no targets file should be created for empty delegation")

	_, _, err = repo.GetDelegation("targets/level1/level2")
	assert.Error(t, err)
	assert.IsType(t, data.ErrNoSuchRole{}, err)
}

// If the parent role doesn't exist, GetDelegation fails with an ErrInvalidRole
func TestGetDelegationParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	_, _, err := repo.GetDelegation("targets/level1/level2")
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// Adding targets to a role that exists and has metadata (like targets)
// correctly adds the target
func TestAddTargetsRoleAndMetadataExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	_, err := repo.AddTargets(data.CanonicalTargetsRole, data.Files{"f": f})
	assert.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	assert.True(t, ok)
	targetsF, ok := r.Signed.Targets["f"]
	assert.True(t, ok)
	assert.Equal(t, f, targetsF)
}

// Adding targets to a role that exists and has not metadata first creates the
// metadata and then correctly adds the target
func TestAddTargetsRoleExistsAndMetadataDoesntExist(t *testing.T) {
	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	assert.False(t, ok, "no empty targets file should be created")

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	assert.NoError(t, err)

	r, ok := repo.Targets["targets/test"]
	assert.True(t, ok)
	targetsF, ok := r.Signed.Targets["f"]
	assert.True(t, ok)
	assert.Equal(t, f, targetsF)
}

// Adding targets to a role that doesn't exist fails
func TestAddTargetsRoleDoesntExist(t *testing.T) {
	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	_, err := repo.AddTargets("targets/test", data.Files{"f": f})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// Adding targets to a role that we don't have signing keys for fails
func TestAddTargetsNoSigningKeys(t *testing.T) {
	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// now delete the signing key (all keys)
	repo.cryptoService = signed.NewEd25519()

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)
}

// Removing targets from a role that exists, has targets, and is signable
// should succeed, even if we also want to remove targets that don't exist.
func TestRemoveExistingAndNonexistingTargets(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	assert.False(t, ok, "no empty targets file should be created")

	// now remove a target
	assert.NoError(t, repo.RemoveTargets("targets/test", "f"))

	// still no metadata
	_, ok = repo.Targets["targets/test"]
	assert.False(t, ok)
}

// Removing targets from a role that exists but without metadata succeeds.
func TestRemoveTargetsNonexistentMetadata(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	err := repo.RemoveTargets("targets/test", "f")
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// Removing targets from a role that doesn't exist fails
func TestRemoveTargetsRoleDoesntExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	err := repo.RemoveTargets("targets/test", "f")
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// Removing targets from a role that we don't have signing keys for fails
func TestRemoveTargetsNoSigningKeys(t *testing.T) {
	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey.ID()}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	assert.NoError(t, err)

	r, ok := repo.Targets["targets/test"]
	assert.True(t, ok)
	_, ok = r.Signed.Targets["f"]
	assert.True(t, ok)

	// now delete the signing key (all keys)
	repo.cryptoService = signed.NewEd25519()

	// now remove the target - it should fail
	err = repo.RemoveTargets("targets/test", "f")
	assert.Error(t, err)
	assert.IsType(t, signed.ErrNoKeys{}, err)
}

// adding a key to a role marks root as dirty as well as the role
func TestAddBaseKeysToRoot(t *testing.T) {
	for _, role := range data.BaseRoles {
		ed25519 := signed.NewEd25519()
		keyDB := keys.NewDB()
		repo := initRepo(t, ed25519, keyDB)

		key, err := ed25519.Create(role, data.ED25519Key)
		assert.NoError(t, err)

		assert.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 1)

		assert.NoError(t, repo.AddBaseKeys(role, key))

		_, ok := repo.Root.Signed.Keys[key.ID()]
		assert.True(t, ok)
		assert.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 2)
		assert.True(t, repo.Root.Dirty)

		switch role {
		case data.CanonicalSnapshotRole:
			assert.True(t, repo.Snapshot.Dirty)
		case data.CanonicalTargetsRole:
			assert.True(t, repo.Targets[data.CanonicalTargetsRole].Dirty)
		case data.CanonicalTimestampRole:
			assert.True(t, repo.Timestamp.Dirty)
		}
	}
}

func TestGetAllRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// After we init, we get the base roles
	roles := repo.GetAllLoadedRoles()
	assert.Len(t, roles, len(data.BaseRoles))

	// Clear the keysDB, check that we get an empty list
	repo.keysDB = keys.NewDB()
	roles = repo.GetAllLoadedRoles()
	assert.Len(t, roles, 0)
}

func TestGetBaseRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// After we init, we get the base roles
	for _, role := range data.BaseRoles {
		baseRole, err := repo.GetBaseRole(role)
		assert.NoError(t, err)

		assert.Equal(t, role, baseRole.GetName())
		keyIDs := repo.cryptoService.ListKeys(role)
		for _, keyID := range keyIDs {
			_, ok := baseRole.Keys[keyID]
			assert.True(t, ok)
			assert.Contains(t, baseRole.ListKeyIDs(), keyID)
		}
		// initRepo should set all key thresholds to 1
		assert.Equal(t, 1, baseRole.GetThreshold())
		assert.False(t, baseRole.IsDelegationRole())
		assert.True(t, baseRole.IsBaseRole())
		_, err = baseRole.ListPathHashPrefixes()
		assert.Error(t, err)
		_, err = baseRole.ListPaths()
		assert.Error(t, err)
	}
}

func TestGetBaseRolesInvalidName(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	_, err := repo.GetBaseRole("invalid")
	assert.Error(t, err)

	_, err = repo.GetBaseRole("targets/delegation")
	assert.Error(t, err)
}

func TestGetDelegationValidRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey1, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey1.ID()}, []string{"path", "anotherpath"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey1})
	assert.NoError(t, err)

	delgRole, err := repo.GetDelegationRole("targets/test")
	assert.NoError(t, err)
	assert.True(t, delgRole.IsDelegationRole())
	assert.False(t, delgRole.IsBaseRole())
	assert.Equal(t, "targets/test", delgRole.GetName())
	assert.Equal(t, 1, delgRole.GetThreshold())
	assert.Equal(t, []string{testKey1.ID()}, delgRole.ListKeyIDs())
	delgPaths, err := delgRole.ListPaths()
	assert.NoError(t, err)
	delgPathPrefixes, err := delgRole.ListPathHashPrefixes()
	assert.NoError(t, err)
	assert.Empty(t, delgPathPrefixes)
	assert.Equal(t, []string{"path", "anotherpath"}, delgPaths)
	assert.Equal(t, testKey1, delgRole.Keys[testKey1.ID()])

	testKey2, err := ed25519.Create("targets/a", data.ED25519Key)
	assert.NoError(t, err)
	role, err = data.NewRole(
		"targets/a", 1, []string{testKey2.ID()}, []string{""}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey2})
	assert.NoError(t, err)

	delgRole, err = repo.GetDelegationRole("targets/a")
	assert.NoError(t, err)
	assert.True(t, delgRole.IsDelegationRole())
	assert.False(t, delgRole.IsBaseRole())
	assert.Equal(t, "targets/a", delgRole.GetName())
	assert.Equal(t, 1, delgRole.GetThreshold())
	assert.Equal(t, []string{testKey2.ID()}, delgRole.ListKeyIDs())
	delgPaths, err = delgRole.ListPaths()
	assert.NoError(t, err)
	assert.Equal(t, []string{""}, delgPaths)
	delgPathPrefixes, err = delgRole.ListPathHashPrefixes()
	assert.NoError(t, err)
	assert.Empty(t, delgPathPrefixes)
	assert.Equal(t, testKey2, delgRole.Keys[testKey2.ID()])

	testKey3, err := ed25519.Create("targets/test/b", data.ED25519Key)
	assert.NoError(t, err)
	role, err = data.NewRole(
		"targets/test/b", 1, []string{testKey3.ID()}, []string{"path/subpath", "anotherpath"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey3})
	assert.NoError(t, err)

	delgRole, err = repo.GetDelegationRole("targets/test/b")
	assert.NoError(t, err)
	assert.True(t, delgRole.IsDelegationRole())
	assert.False(t, delgRole.IsBaseRole())
	assert.Equal(t, "targets/test/b", delgRole.GetName())
	assert.Equal(t, 1, delgRole.GetThreshold())
	assert.Equal(t, []string{testKey3.ID()}, delgRole.ListKeyIDs())
	delgPaths, err = delgRole.ListPaths()
	assert.NoError(t, err)
	assert.Equal(t, []string{"path/subpath", "anotherpath"}, delgPaths)
	delgPathPrefixes, err = delgRole.ListPathHashPrefixes()
	assert.NoError(t, err)
	assert.Empty(t, delgPathPrefixes)
	assert.Equal(t, testKey3, delgRole.Keys[testKey3.ID()])

	testKey4, err := ed25519.Create("targets/test/c", data.ED25519Key)
	assert.NoError(t, err)
	// Try adding empty paths, ensure this is valid
	role, err = data.NewRole(
		"targets/test/c", 1, []string{testKey4.ID()}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey3})
	assert.NoError(t, err)
}

func TestGetDelegationRolesInvalidName(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	_, err := repo.GetDelegationRole("invalid")
	assert.Error(t, err)

	for _, role := range data.BaseRoles {
		_, err = repo.GetDelegationRole(role)
		assert.Error(t, err)
	}
}

func TestGetDelegationRolesInvalidPaths(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey1, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey1.ID()}, []string{"path", "anotherpath"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey1})
	assert.NoError(t, err)

	testKey2, err := ed25519.Create("targets/test/b", data.ED25519Key)
	assert.NoError(t, err)
	// Now we add a delegation with a path that is not prefixed by its parent delegation
	role, err = data.NewRole(
		"targets/test/b", 1, []string{testKey2.ID()}, []string{"invalidpath"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey2})
	assert.NoError(t, err)

	// Getting this delegation should fail path verification, so it'll have 0 paths
	delgRole, err := repo.GetDelegationRole("targets/test/b")
	assert.NoError(t, err)
	delgPaths, err := delgRole.ListPaths()
	assert.NoError(t, err)
	assert.Empty(t, delgPaths)
}

func TestGetDelegationRolesInvalidPathHashPrefix(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey1, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole(
		"targets/test", 1, []string{testKey1.ID()}, []string{}, []string{"pathhash", "anotherpathhash"})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey1})
	assert.NoError(t, err)

	testKey2, err := ed25519.Create("targets/test/b", data.ED25519Key)
	assert.NoError(t, err)
	// Now we add a delegation with a path that is not prefixed by its parent delegation
	role, err = data.NewRole(
		"targets/test/b", 1, []string{testKey2.ID()}, []string{}, []string{"invalidpathhash"})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey2})
	assert.NoError(t, err)

	// Getting this delegation should fail path verification
	delgRole, err := repo.GetDelegationRole("targets/test/b")
	assert.NoError(t, err)
	delgPathHashPrefixes, err := delgRole.ListPathHashPrefixes()
	assert.NoError(t, err)
	assert.Empty(t, delgPathHashPrefixes)
}
