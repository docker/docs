package tuf

import (
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

	r := repo.Targets[data.CanonicalTargetsRole]
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

	r = repo.Targets["targets/test"]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testDeepKey.ID(), keyIDs[0])
	assert.True(t, r.Dirty)
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

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 0)
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

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 0)
}

func TestUpdateDelegationsRoleMissingKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole("targets/role", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	// key should get added to role as part of updating the delegation
	err = repo.UpdateDelegations(role, data.KeyList{roleKey})
	assert.NoError(t, err)

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, roleKey.ID(), keyIDs[0])
	assert.True(t, r.Dirty)
}

func TestUpdateDelegationsNotEnoughKeys(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	roleKey, err := ed25519.Create("Invalid Role", data.ED25519Key)
	assert.NoError(t, err)

	role, err := data.NewRole("targets/role", 2, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	// key should get added to role as part of updating the delegation
	err = repo.UpdateDelegations(role, data.KeyList{roleKey})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
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

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	// create another role with the same name and ensure it replaces the
	// previous role
	testKey2, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role2, err := data.NewRole("targets/test", 1, []string{testKey2.ID()}, []string{"test"}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role2, data.KeyList{testKey2})
	assert.NoError(t, err)

	r = repo.Targets["targets"]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey2.ID(), keyIDs[0])
	assert.True(t, r.Dirty)
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

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	testKey2, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey2})
	assert.NoError(t, err)

	r = repo.Targets["targets"]
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

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 1)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	assert.Len(t, keyIDs, 1)
	assert.Equal(t, testKey.ID(), keyIDs[0])

	err = repo.DeleteDelegation(*role)
	assert.Len(t, r.Signed.Delegations.Roles, 0)
	assert.Len(t, r.Signed.Delegations.Keys, 0)
	assert.True(t, r.Dirty)
}

func TestDeleteDelegationsRoleNotExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	// initRepo leaves all the roles as Dirty. Set to false
	// to test removing a non-existant role doesn't mark
	// a role as dirty
	repo.Targets[data.CanonicalTargetsRole].Dirty = false

	role, err := data.NewRole("targets/test", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*role)
	assert.NoError(t, err)
	r := repo.Targets[data.CanonicalTargetsRole]
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
	invalidRole, err := data.NewRole("root", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*invalidRole)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 0)
}

func TestDeleteDelegationsParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testRole, err := data.NewRole("targets/test/deep", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*testRole)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 0)
}

func TestDeleteDelegationsMidSliceRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	assert.NoError(t, err)
	role, err := data.NewRole("targets/test", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role, data.KeyList{testKey})
	assert.NoError(t, err)

	role2, err := data.NewRole("targets/test2", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role2, data.KeyList{testKey})
	assert.NoError(t, err)

	role3, err := data.NewRole("targets/test3", 1, []string{}, []string{}, []string{})
	assert.NoError(t, err)

	err = repo.UpdateDelegations(role3, data.KeyList{testKey})
	assert.NoError(t, err)

	err = repo.DeleteDelegation(*role2)
	assert.NoError(t, err)

	r := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, r.Signed.Delegations.Roles, 2)
	assert.Len(t, r.Signed.Delegations.Keys, 1)
	assert.True(t, r.Dirty)
}

func TestGetDelegationParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	_, err := repo.GetDelegation("targets/level1/level2")
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}
