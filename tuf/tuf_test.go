package tuf

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/stretchr/testify/require"
)

var testGUN = "gun"

func initRepo(t *testing.T, cryptoService signed.CryptoService) *Repo {
	rootKey, err := cryptoService.Create("root", testGUN, data.ED25519Key)
	require.NoError(t, err)
	return initRepoWithRoot(t, cryptoService, rootKey)
}

func initRepoWithRoot(t *testing.T, cryptoService signed.CryptoService, rootKey data.PublicKey) *Repo {
	targetsKey, err := cryptoService.Create("targets", testGUN, data.ED25519Key)
	require.NoError(t, err)
	snapshotKey, err := cryptoService.Create("snapshot", testGUN, data.ED25519Key)
	require.NoError(t, err)
	timestampKey, err := cryptoService.Create("timestamp", testGUN, data.ED25519Key)
	require.NoError(t, err)

	rootRole := data.NewBaseRole(
		data.CanonicalRootRole,
		1,
		rootKey,
	)
	targetsRole := data.NewBaseRole(
		data.CanonicalTargetsRole,
		1,
		targetsKey,
	)
	snapshotRole := data.NewBaseRole(
		data.CanonicalSnapshotRole,
		1,
		snapshotKey,
	)
	timestampRole := data.NewBaseRole(
		data.CanonicalTimestampRole,
		1,
		timestampKey,
	)

	repo := NewRepo(cryptoService)
	err = repo.InitRoot(rootRole, timestampRole, snapshotRole, targetsRole, false)
	require.NoError(t, err)
	_, err = repo.InitTargets(data.CanonicalTargetsRole)
	require.NoError(t, err)
	err = repo.InitSnapshot()
	require.NoError(t, err)
	err = repo.InitTimestamp()
	require.NoError(t, err)
	return repo
}

func TestInitSnapshotNoTargets(t *testing.T) {
	cs := signed.NewEd25519()
	repo := initRepo(t, cs)

	repo.Targets = make(map[string]*data.SignedTargets)

	err := repo.InitSnapshot()
	require.Error(t, err)
	require.IsType(t, ErrNotLoaded{}, err)
}

func writeRepo(t *testing.T, dir string, repo *Repo) {
	err := os.MkdirAll(dir, 0755)
	require.NoError(t, err)
	signedRoot, err := repo.SignRoot(data.DefaultExpires("root"))
	require.NoError(t, err)
	rootJSON, _ := json.Marshal(signedRoot)
	ioutil.WriteFile(dir+"/root.json", rootJSON, 0755)

	for r := range repo.Targets {
		signedTargets, err := repo.SignTargets(r, data.DefaultExpires("targets"))
		require.NoError(t, err)
		targetsJSON, _ := json.Marshal(signedTargets)
		p := path.Join(dir, r+".json")
		parentDir := filepath.Dir(p)
		os.MkdirAll(parentDir, 0755)
		ioutil.WriteFile(p, targetsJSON, 0755)
	}

	signedSnapshot, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	require.NoError(t, err)
	snapshotJSON, _ := json.Marshal(signedSnapshot)
	ioutil.WriteFile(dir+"/snapshot.json", snapshotJSON, 0755)

	signedTimestamp, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	require.NoError(t, err)
	timestampJSON, _ := json.Marshal(signedTimestamp)
	ioutil.WriteFile(dir+"/timestamp.json", timestampJSON, 0755)
}

func TestInitRepo(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)
	writeRepo(t, "/tmp/tufrepo", repo)
	// after signing a new repo, there are only 4 roles: the 4 base roles
	require.Len(t, repo.Root.Signed.Roles, 4)

	// can't use getBaseRole because it's not a valid real role
	_, err := repo.Root.BuildBaseRole("root.1")
	require.Error(t, err)
}

func TestUpdateDelegations(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	require.False(t, ok, "no empty targets file should be created for deepest delegation")

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, testKey.ID(), keyIDs[0])

	testDeepKey, err := ed25519.Create("targets/test/deep", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test/deep", []data.PublicKey{testDeepKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test/deep", []string{"test/deep"}, []string{}, false)
	require.NoError(t, err)

	// this metadata didn't exist before, but creating targets/test/deep created
	// the targets/test metadata
	r, ok = repo.Targets["targets/test"]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, testDeepKey.ID(), keyIDs[0])
	require.True(t, r.Dirty)

	// no empty delegation metadata is created for targets/test/deep
	_, ok = repo.Targets["targets/test/deep"]
	require.False(t, ok, "no empty targets file should be created for deepest delegation")
}

func TestUpdateDelegationsParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testDeepKey, err := ed25519.Create("targets/test/deep", testGUN, data.ED25519Key)
	err = repo.UpdateDelegationKeys("targets/test/deep", []data.PublicKey{testDeepKey}, []string{}, 1)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 0)

	// no delegation metadata created for non-existent parent
	_, ok = repo.Targets["targets/test"]
	require.False(t, ok, "no targets file should be created for nonexistent parent delegation")
}

// Updating delegations needs to modify the parent of the role being updated.
// If there is no signing key for that parent, the delegation cannot be added.
func TestUpdateDelegationsMissingParentKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// remove the target key (all keys)
	repo.cryptoService = signed.NewEd25519()

	roleKey, err := ed25519.Create("Invalid Role", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/role", []data.PublicKey{roleKey}, []string{}, 1)
	require.Error(t, err)
	require.IsType(t, signed.ErrNoKeys{}, err)

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/role"]
	require.False(t, ok, "no targets file should be created for empty delegation")
}

func TestUpdateDelegationsInvalidRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	roleKey, err := ed25519.Create("Invalid Role", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("root", []data.PublicKey{roleKey}, []string{}, 1)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 0)

	// no delegation metadata created for invalid delegation
	_, ok = repo.Targets["root"]
	require.False(t, ok, "no targets file should be created since delegation failed")
}

// A delegation can be created with a role that is missing a signing key, so
// long as UpdateDelegations is called with the key
func TestUpdateDelegationsRoleThatIsMissingDelegationKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	roleKey, err := ed25519.Create("Invalid Role", testGUN, data.ED25519Key)
	require.NoError(t, err)

	// key should get added to role as part of updating the delegation
	err = repo.UpdateDelegationKeys("targets/role", []data.PublicKey{roleKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/role", []string{""}, []string{}, false)
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, roleKey.ID(), keyIDs[0])
	require.True(t, r.Dirty)

	// no empty delegation metadata created for new delegation
	_, ok = repo.Targets["targets/role"]
	require.False(t, ok, "no targets file should be created for empty delegation")
}

func TestUpdateDelegationsNotEnoughKeys(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	roleKey, err := ed25519.Create("Invalid Role", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/role", []data.PublicKey{roleKey}, []string{}, 2)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	// no delegation metadata created for failed delegation
	_, ok := repo.Targets["targets/role"]
	require.False(t, ok, "no targets file should be created since delegation failed")
}

func TestUpdateDelegationsAddKeyToRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, testKey.ID(), keyIDs[0])

	testKey2, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey2}, []string{}, 1)
	require.NoError(t, err)

	r, ok = repo.Targets["targets"]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 2)
	keyIDs = r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 2)
	// it does an append so the order is deterministic (but not meaningful to TUF)
	require.Equal(t, testKey.ID(), keyIDs[0])
	require.Equal(t, testKey2.ID(), keyIDs[1])
	require.True(t, r.Dirty)
}

func TestDeleteDelegations(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, testKey.ID(), keyIDs[0])

	// ensure that the metadata is there and snapshot is there
	targets, err := repo.InitTargets("targets/test")
	require.NoError(t, err)
	targetsSigned, err := targets.ToSigned()
	require.NoError(t, err)
	require.NoError(t, repo.UpdateSnapshot("targets/test", targetsSigned))
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	require.True(t, ok)

	require.NoError(t, repo.DeleteDelegation("targets/test"))
	require.Len(t, r.Signed.Delegations.Roles, 0)
	require.Len(t, r.Signed.Delegations.Keys, 0)
	require.True(t, r.Dirty)

	// metadata should be deleted
	_, ok = repo.Targets["targets/test"]
	require.False(t, ok)
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	require.False(t, ok)
}

func TestDeleteDelegationsRoleNotExistBecauseNoParentMeta(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/test"]
	require.False(t, ok, "no targets file should be created for empty delegation")

	delRole, err := data.NewRole("targets/test/a", 1, []string{testKey.ID()}, []string{"test"})

	err = repo.DeleteDelegation(delRole.Name)
	require.NoError(t, err)
	// still no metadata
	_, ok = repo.Targets["targets/test"]
	require.False(t, ok)
}

func TestDeleteDelegationsRoleNotExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// initRepo leaves all the roles as Dirty. Set to false
	// to test removing a non-existent role doesn't mark
	// a role as dirty
	repo.Targets[data.CanonicalTargetsRole].Dirty = false

	role, err := data.NewRole("targets/test", 1, []string{}, []string{""})
	require.NoError(t, err)

	err = repo.DeleteDelegation(role.Name)
	require.NoError(t, err)
	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 0)
	require.Len(t, r.Signed.Delegations.Keys, 0)
	require.False(t, r.Dirty)
}

func TestDeleteDelegationsInvalidRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// data.NewRole errors if the role isn't a valid TUF role so use one of the non-delegation
	// valid roles
	invalidRole, err := data.NewRole("root", 1, []string{}, []string{""})
	require.NoError(t, err)

	err = repo.DeleteDelegation(invalidRole.Name)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 0)
}

func TestDeleteDelegationsParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testRole, err := data.NewRole("targets/test/deep", 1, []string{}, []string{""})
	require.NoError(t, err)

	err = repo.DeleteDelegation(testRole.Name)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 0)
}

// Can't delete a delegation if we don't have the parent's signing key
func TestDeleteDelegationsMissingParentSigningKey(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	keyIDs := r.Signed.Delegations.Roles[0].KeyIDs
	require.Len(t, keyIDs, 1)
	require.Equal(t, testKey.ID(), keyIDs[0])

	// ensure that the metadata is there and snapshot is there
	targets, err := repo.InitTargets("targets/test")
	require.NoError(t, err)
	targetsSigned, err := targets.ToSigned()
	require.NoError(t, err)
	require.NoError(t, repo.UpdateSnapshot("targets/test", targetsSigned))
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	require.True(t, ok)

	// delete all signing keys
	repo.cryptoService = signed.NewEd25519()
	err = repo.DeleteDelegation("targets/test")
	require.Error(t, err)
	require.IsType(t, signed.ErrNoKeys{}, err)

	require.Len(t, r.Signed.Delegations.Roles, 1)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	require.True(t, r.Dirty)

	// metadata should be here still
	_, ok = repo.Targets["targets/test"]
	require.True(t, ok)
	_, ok = repo.Snapshot.Signed.Meta["targets/test"]
	require.True(t, ok)
}

func TestDeleteDelegationsMidSliceRole(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{""}, []string{}, false)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test2", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test2", []string{""}, []string{}, false)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test3", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test3", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	err = repo.DeleteDelegation("targets/test2")
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	require.Len(t, r.Signed.Delegations.Roles, 2)
	require.Len(t, r.Signed.Delegations.Keys, 1)
	require.True(t, r.Dirty)
}

// If the parent exists, the metadata exists, and the delegation is in it,
// returns the role that was found
func TestGetDelegationRoleAndMetadataExistDelegationExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("meh", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/level1", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/level1/level2", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/level1/level2", []string{""}, []string{}, false)
	require.NoError(t, err)

	gottenRole, err := repo.GetDelegationRole("targets/level1/level2")
	require.NoError(t, err)
	require.Equal(t, "targets/level1/level2", gottenRole.Name)
	require.Equal(t, 1, gottenRole.Threshold)
	require.Equal(t, []string{""}, gottenRole.Paths)
	_, ok := gottenRole.Keys[testKey.ID()]
	require.True(t, ok)
}

// If the parent exists, the metadata exists, and the delegation isn't in it,
// returns an ErrNoSuchRole
func TestGetDelegationRoleAndMetadataExistDelegationDoesntExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("meh", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/level1", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	// ensure metadata exists
	repo.InitTargets("targets/level1")

	_, err = repo.GetDelegationRole("targets/level1/level2")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

// If the parent exists but the metadata doesn't exist, returns an ErrNoSuchRole
func TestGetDelegationRoleAndMetadataDoesntExists(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("meh", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/level1", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/level1", []string{""}, []string{}, false)
	require.NoError(t, err)

	// no empty delegation metadata created for new delegation
	_, ok := repo.Targets["targets/test"]
	require.False(t, ok, "no targets file should be created for empty delegation")

	_, err = repo.GetDelegationRole("targets/level1/level2")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

// If the parent role doesn't exist, GetDelegation fails with an ErrInvalidRole
func TestGetDelegationParentMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	_, err := repo.GetDelegationRole("targets/level1/level2")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

// Adding targets to a role that exists and has metadata (like targets)
// correctly adds the target
func TestAddTargetsRoleAndMetadataExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	hash := sha256.Sum256([]byte{})
	f := data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}

	_, err := repo.AddTargets(data.CanonicalTargetsRole, data.Files{"f": f})
	require.NoError(t, err)

	r, ok := repo.Targets[data.CanonicalTargetsRole]
	require.True(t, ok)
	targetsF, ok := r.Signed.Targets["f"]
	require.True(t, ok)
	require.Equal(t, f, targetsF)
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
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{""}, []string{}, false)
	require.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	require.False(t, ok, "no empty targets file should be created")

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	require.NoError(t, err)

	r, ok := repo.Targets["targets/test"]
	require.True(t, ok)
	targetsF, ok := r.Signed.Targets["f"]
	require.True(t, ok)
	require.Equal(t, f, targetsF)
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
	repo := initRepo(t, ed25519)

	_, err := repo.AddTargets("targets/test", data.Files{"f": f})
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
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
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	// now delete the signing key (all keys)
	repo.cryptoService = signed.NewEd25519()

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	require.Error(t, err)
	require.IsType(t, signed.ErrNoKeys{}, err)
}

// Removing targets from a role that exists, has targets, and is signable
// should succeed, even if we also want to remove targets that don't exist.
func TestRemoveExistingAndNonexistingTargets(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"test"}, []string{}, false)
	require.NoError(t, err)

	// no empty metadata is created for this role
	_, ok := repo.Targets["targets/test"]
	require.False(t, ok, "no empty targets file should be created")

	// now remove a target
	require.NoError(t, repo.RemoveTargets("targets/test", "f"))

	// still no metadata
	_, ok = repo.Targets["targets/test"]
	require.False(t, ok)
}

// Removing targets from a role that exists but without metadata succeeds.
func TestRemoveTargetsNonexistentMetadata(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	err := repo.RemoveTargets("targets/test", "f")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

// Removing targets from a role that doesn't exist fails
func TestRemoveTargetsRoleDoesntExist(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	err := repo.RemoveTargets("targets/test", "f")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
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
	repo := initRepo(t, ed25519)

	testKey, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{""}, []string{}, false)
	require.NoError(t, err)

	// adding the targets to the role should create the metadata though
	_, err = repo.AddTargets("targets/test", data.Files{"f": f})
	require.NoError(t, err)

	r, ok := repo.Targets["targets/test"]
	require.True(t, ok)
	_, ok = r.Signed.Targets["f"]
	require.True(t, ok)

	// now delete the signing key (all keys)
	repo.cryptoService = signed.NewEd25519()

	// now remove the target - it should fail
	err = repo.RemoveTargets("targets/test", "f")
	require.Error(t, err)
	require.IsType(t, signed.ErrNoKeys{}, err)
}

// adding a key to a role marks root as dirty as well as the role
func TestAddBaseKeysToRoot(t *testing.T) {
	for _, role := range data.BaseRoles {
		ed25519 := signed.NewEd25519()
		repo := initRepo(t, ed25519)

		origKeyIDs := ed25519.ListKeys(role)
		require.Len(t, origKeyIDs, 1)

		key, err := ed25519.Create(role, testGUN, data.ED25519Key)
		require.NoError(t, err)

		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 1)

		require.NoError(t, repo.AddBaseKeys(role, key))

		_, ok := repo.Root.Signed.Keys[key.ID()]
		require.True(t, ok)
		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 2)
		require.True(t, repo.Root.Dirty)

		switch role {
		case data.CanonicalSnapshotRole:
			require.True(t, repo.Snapshot.Dirty)
		case data.CanonicalTargetsRole:
			require.True(t, repo.Targets[data.CanonicalTargetsRole].Dirty)
		case data.CanonicalTimestampRole:
			require.True(t, repo.Timestamp.Dirty)
		case data.CanonicalRootRole:
			require.NoError(t, err)
			require.Len(t, repo.originalRootRole.Keys, 1)
			require.Contains(t, repo.originalRootRole.ListKeyIDs(), origKeyIDs[0])
		}
	}
}

// removing one or more keys from a role marks root as dirty as well as the role
func TestRemoveBaseKeysFromRoot(t *testing.T) {
	for _, role := range data.BaseRoles {
		ed25519 := signed.NewEd25519()
		repo := initRepo(t, ed25519)

		origKeyIDs := ed25519.ListKeys(role)
		require.Len(t, origKeyIDs, 1)

		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 1)

		require.NoError(t, repo.RemoveBaseKeys(role, origKeyIDs...))

		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 0)
		require.True(t, repo.Root.Dirty)

		switch role {
		case data.CanonicalSnapshotRole:
			require.True(t, repo.Snapshot.Dirty)
		case data.CanonicalTargetsRole:
			require.True(t, repo.Targets[data.CanonicalTargetsRole].Dirty)
		case data.CanonicalTimestampRole:
			require.True(t, repo.Timestamp.Dirty)
		case data.CanonicalRootRole:
			require.Len(t, repo.originalRootRole.Keys, 1)
			require.Contains(t, repo.originalRootRole.ListKeyIDs(), origKeyIDs[0])
		}
	}
}

// replacing keys in a role marks root as dirty as well as the role
func TestReplaceBaseKeysInRoot(t *testing.T) {
	for _, role := range data.BaseRoles {
		ed25519 := signed.NewEd25519()
		repo := initRepo(t, ed25519)

		origKeyIDs := ed25519.ListKeys(role)
		require.Len(t, origKeyIDs, 1)

		key, err := ed25519.Create(role, testGUN, data.ED25519Key)
		require.NoError(t, err)

		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 1)

		require.NoError(t, repo.ReplaceBaseKeys(role, key))

		_, ok := repo.Root.Signed.Keys[key.ID()]
		require.True(t, ok)
		require.Len(t, repo.Root.Signed.Roles[role].KeyIDs, 1)
		require.True(t, repo.Root.Dirty)

		switch role {
		case data.CanonicalSnapshotRole:
			require.True(t, repo.Snapshot.Dirty)
		case data.CanonicalTargetsRole:
			require.True(t, repo.Targets[data.CanonicalTargetsRole].Dirty)
		case data.CanonicalTimestampRole:
			require.True(t, repo.Timestamp.Dirty)
		case data.CanonicalRootRole:
			require.Len(t, repo.originalRootRole.Keys, 1)
			require.Contains(t, repo.originalRootRole.ListKeyIDs(), origKeyIDs[0])
		}

		origNumRoles := len(repo.Root.Signed.Roles)
		// sign the root and assert the number of roles after
		_, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
		require.NoError(t, err)

		switch role {
		case data.CanonicalRootRole:
			// root role changed, so the old role should have been saved
			require.Len(t, repo.Root.Signed.Roles, origNumRoles+1)
		default:
			// number of roles should not have changed
			require.Len(t, repo.Root.Signed.Roles, origNumRoles)
		}
	}
}

func TestGetAllRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// After we init, we get the base roles
	roles := repo.GetAllLoadedRoles()
	require.Len(t, roles, len(data.BaseRoles))
}

func TestGetBaseRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// After we init, we get the base roles
	for _, role := range data.BaseRoles {
		baseRole, err := repo.GetBaseRole(role)
		require.NoError(t, err)

		require.Equal(t, role, baseRole.Name)
		keyIDs := repo.cryptoService.ListKeys(role)
		for _, keyID := range keyIDs {
			_, ok := baseRole.Keys[keyID]
			require.True(t, ok)
			require.Contains(t, baseRole.ListKeyIDs(), keyID)
		}
		// initRepo should set all key thresholds to 1
		require.Equal(t, 1, baseRole.Threshold)
	}
}

func TestGetBaseRolesInvalidName(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	_, err := repo.GetBaseRole("invalid")
	require.Error(t, err)

	_, err = repo.GetBaseRole("targets/delegation")
	require.Error(t, err)
}

func TestGetDelegationValidRoles(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey1, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey1}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"path", "anotherpath"}, []string{}, false)
	require.NoError(t, err)

	delgRole, err := repo.GetDelegationRole("targets/test")
	require.NoError(t, err)
	require.Equal(t, "targets/test", delgRole.Name)
	require.Equal(t, 1, delgRole.Threshold)
	require.Equal(t, []string{testKey1.ID()}, delgRole.ListKeyIDs())
	require.Equal(t, []string{"path", "anotherpath"}, delgRole.Paths)
	require.Equal(t, testKey1, delgRole.Keys[testKey1.ID()])

	testKey2, err := ed25519.Create("targets/a", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/a", []data.PublicKey{testKey2}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/a", []string{""}, []string{}, false)
	require.NoError(t, err)

	delgRole, err = repo.GetDelegationRole("targets/a")
	require.NoError(t, err)
	require.Equal(t, "targets/a", delgRole.Name)
	require.Equal(t, 1, delgRole.Threshold)
	require.Equal(t, []string{testKey2.ID()}, delgRole.ListKeyIDs())
	require.Equal(t, []string{""}, delgRole.Paths)
	require.Equal(t, testKey2, delgRole.Keys[testKey2.ID()])

	testKey3, err := ed25519.Create("targets/test/b", testGUN, data.ED25519Key)
	require.NoError(t, err)
	err = repo.UpdateDelegationKeys("targets/test/b", []data.PublicKey{testKey3}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test/b", []string{"path/subpath", "anotherpath"}, []string{}, false)
	require.NoError(t, err)

	delgRole, err = repo.GetDelegationRole("targets/test/b")
	require.NoError(t, err)
	require.Equal(t, "targets/test/b", delgRole.Name)
	require.Equal(t, 1, delgRole.Threshold)
	require.Equal(t, []string{testKey3.ID()}, delgRole.ListKeyIDs())
	require.Equal(t, []string{"path/subpath", "anotherpath"}, delgRole.Paths)
	require.Equal(t, testKey3, delgRole.Keys[testKey3.ID()])

	testKey4, err := ed25519.Create("targets/test/c", testGUN, data.ED25519Key)
	require.NoError(t, err)
	// Try adding empty paths, ensure this is valid
	err = repo.UpdateDelegationKeys("targets/test/c", []data.PublicKey{testKey4}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test/c", []string{}, []string{}, false)
	require.NoError(t, err)
}

func TestGetDelegationRolesInvalidName(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	_, err := repo.GetDelegationRole("invalid")
	require.Error(t, err)

	for _, role := range data.BaseRoles {
		_, err = repo.GetDelegationRole(role)
		require.Error(t, err)
		require.IsType(t, data.ErrInvalidRole{}, err)
	}
	_, err = repo.GetDelegationRole("targets/doesnt_exist")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

func TestGetDelegationRolesInvalidPaths(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	testKey1, err := ed25519.Create("targets/test", testGUN, data.ED25519Key)
	require.NoError(t, err)

	err = repo.UpdateDelegationKeys("targets/test", []data.PublicKey{testKey1}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test", []string{"path", "anotherpath"}, []string{}, false)
	require.NoError(t, err)

	testKey2, err := ed25519.Create("targets/test/b", testGUN, data.ED25519Key)
	require.NoError(t, err)
	// Now we add a delegation with a path that is not prefixed by its parent delegation, the invalid path can't be added so there is an error
	err = repo.UpdateDelegationKeys("targets/test/b", []data.PublicKey{testKey2}, []string{}, 1)
	require.NoError(t, err)
	err = repo.UpdateDelegationPaths("targets/test/b", []string{"invalidpath"}, []string{}, false)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)

	delgRole, err := repo.GetDelegationRole("targets/test")
	require.NoError(t, err)
	require.Contains(t, delgRole.Paths, "path")
	require.Contains(t, delgRole.Paths, "anotherpath")
}

func TestDelegationRolesParent(t *testing.T) {
	delgA := data.DelegationRole{
		BaseRole: data.BaseRole{
			Keys:      nil,
			Name:      "targets/a",
			Threshold: 1,
		},
		Paths: []string{"path", "anotherpath"},
	}

	delgB := data.DelegationRole{
		BaseRole: data.BaseRole{
			Keys:      nil,
			Name:      "targets/a/b",
			Threshold: 1,
		},
		Paths: []string{"path/b", "anotherpath/b", "b/invalidpath"},
	}

	// Assert direct parent relationship
	require.True(t, delgA.IsParentOf(delgB))
	require.False(t, delgB.IsParentOf(delgA))
	require.False(t, delgA.IsParentOf(delgA))

	delgC := data.DelegationRole{
		BaseRole: data.BaseRole{
			Keys:      nil,
			Name:      "targets/a/b/c",
			Threshold: 1,
		},
		Paths: []string{"path/b", "anotherpath/b/c", "c/invalidpath"},
	}

	// Assert direct parent relationship
	require.True(t, delgB.IsParentOf(delgC))
	require.False(t, delgB.IsParentOf(delgB))
	require.False(t, delgA.IsParentOf(delgC))
	require.False(t, delgC.IsParentOf(delgB))
	require.False(t, delgC.IsParentOf(delgA))
	require.False(t, delgC.IsParentOf(delgC))

	// Check that parents correctly restrict paths
	restrictedDelgB, err := delgA.Restrict(delgB)
	require.NoError(t, err)
	require.Contains(t, restrictedDelgB.Paths, "path/b")
	require.Contains(t, restrictedDelgB.Paths, "anotherpath/b")
	require.NotContains(t, restrictedDelgB.Paths, "b/invalidpath")

	_, err = delgB.Restrict(delgA)
	require.Error(t, err)
	_, err = delgA.Restrict(delgC)
	require.Error(t, err)
	_, err = delgC.Restrict(delgB)
	require.Error(t, err)
	_, err = delgC.Restrict(delgA)
	require.Error(t, err)

	// Make delgA have no paths and check that it changes delgB and delgC accordingly when chained
	delgA.Paths = []string{}
	restrictedDelgB, err = delgA.Restrict(delgB)
	require.NoError(t, err)
	require.Empty(t, restrictedDelgB.Paths)
	restrictedDelgC, err := restrictedDelgB.Restrict(delgC)
	require.NoError(t, err)
	require.Empty(t, restrictedDelgC.Paths)
}

func TestGetBaseRoleEmptyRepo(t *testing.T) {
	repo := NewRepo(nil)
	_, err := repo.GetBaseRole(data.CanonicalRootRole)
	require.Error(t, err)
	require.IsType(t, ErrNotLoaded{}, err)
}

func TestGetBaseRoleKeyMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// change root role to have a KeyID that doesn't exist
	repo.Root.Signed.Roles[data.CanonicalRootRole].KeyIDs = []string{"abc"}

	_, err := repo.GetBaseRole(data.CanonicalRootRole)
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

func TestGetDelegationRoleKeyMissing(t *testing.T) {
	ed25519 := signed.NewEd25519()
	repo := initRepo(t, ed25519)

	// add a delegation that has a KeyID that doesn't exist
	// in the relevant key map
	tar := repo.Targets[data.CanonicalTargetsRole]
	tar.Signed.Delegations.Roles = []*data.Role{
		{
			RootRole: data.RootRole{
				KeyIDs:    []string{"abc"},
				Threshold: 1,
			},
			Name:  "targets/missing_key",
			Paths: []string{""},
		},
	}

	_, err := repo.GetDelegationRole("targets/missing_key")
	require.Error(t, err)
	require.IsType(t, data.ErrInvalidRole{}, err)
}

func verifySignatureList(t *testing.T, signed *data.Signed, expectedKeys ...data.PublicKey) {
	require.Equal(t, len(expectedKeys), len(signed.Signatures))
	usedKeys := make(map[string]struct{}, len(signed.Signatures))
	for _, sig := range signed.Signatures {
		usedKeys[sig.KeyID] = struct{}{}
	}
	for _, key := range expectedKeys {
		_, ok := usedKeys[key.ID()]
		require.True(t, ok)
		verifyRootSignatureAgainstKey(t, signed, key)
	}
}

func verifyRootSignatureAgainstKey(t *testing.T, signedRoot *data.Signed, key data.PublicKey) error {
	roleWithKeys := data.BaseRole{Name: data.CanonicalRootRole, Keys: data.Keys{key.ID(): key}, Threshold: 1}
	return signed.VerifySignatures(signedRoot, roleWithKeys)
}

func TestSignRootOldKeyCertExists(t *testing.T) {
	gun := "docker/test-sign-root"
	referenceTime := time.Now()

	cs := cryptoservice.NewCryptoService(trustmanager.NewKeyMemoryStore(
		passphrase.ConstantRetriever("password")))

	rootPublicKey, err := cs.Create(data.CanonicalRootRole, gun, data.ECDSAKey)
	require.NoError(t, err)
	rootPrivateKey, _, err := cs.GetPrivateKey(rootPublicKey.ID())
	require.NoError(t, err)
	oldRootCert, err := cryptoservice.GenerateCertificate(rootPrivateKey, gun, referenceTime.AddDate(-9, 0, 0),
		referenceTime.AddDate(1, 0, 0))
	require.NoError(t, err)
	oldRootCertKey := trustmanager.CertToKey(oldRootCert)

	repo := initRepoWithRoot(t, cs, oldRootCertKey)

	// Create a first signature, using the old key.
	signedRoot, err := repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedRoot, oldRootCertKey)
	err = verifyRootSignatureAgainstKey(t, signedRoot, oldRootCertKey)
	require.NoError(t, err)

	// Create a new certificate
	newRootCert, err := cryptoservice.GenerateCertificate(rootPrivateKey, gun, referenceTime, referenceTime.AddDate(10, 0, 0))
	require.NoError(t, err)
	newRootCertKey := trustmanager.CertToKey(newRootCert)
	require.NotEqual(t, oldRootCertKey.ID(), newRootCertKey.ID())

	// Only trust the new certificate
	err = repo.ReplaceBaseKeys(data.CanonicalRootRole, newRootCertKey)
	require.NoError(t, err)
	updatedRootRole, err := repo.GetBaseRole(data.CanonicalRootRole)
	require.NoError(t, err)
	updatedRootKeyIDs := updatedRootRole.ListKeyIDs()
	require.Equal(t, 1, len(updatedRootKeyIDs))
	require.Equal(t, newRootCertKey.ID(), updatedRootKeyIDs[0])

	// Create a second signature
	signedRoot, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedRoot, oldRootCertKey, newRootCertKey)

	// Verify that the signature can be verified when trusting the old certificate
	err = verifyRootSignatureAgainstKey(t, signedRoot, oldRootCertKey)
	require.NoError(t, err)
	// Verify that the signature can be verified when trusting the new certificate
	err = verifyRootSignatureAgainstKey(t, signedRoot, newRootCertKey)
	require.NoError(t, err)
}

func TestSignRootOldKeyCertMissing(t *testing.T) {
	gun := "docker/test-sign-root"
	referenceTime := time.Now()

	cs := cryptoservice.NewCryptoService(trustmanager.NewKeyMemoryStore(
		passphrase.ConstantRetriever("password")))

	rootPublicKey, err := cs.Create(data.CanonicalRootRole, gun, data.ECDSAKey)
	require.NoError(t, err)
	rootPrivateKey, _, err := cs.GetPrivateKey(rootPublicKey.ID())
	require.NoError(t, err)
	oldRootCert, err := cryptoservice.GenerateCertificate(rootPrivateKey, gun, referenceTime.AddDate(-9, 0, 0),
		referenceTime.AddDate(1, 0, 0))
	require.NoError(t, err)
	oldRootCertKey := trustmanager.CertToKey(oldRootCert)

	repo := initRepoWithRoot(t, cs, oldRootCertKey)

	// Create a first signature, using the old key.
	signedRoot, err := repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedRoot, oldRootCertKey)
	err = verifyRootSignatureAgainstKey(t, signedRoot, oldRootCertKey)
	require.NoError(t, err)

	// Create a new certificate
	newRootCert, err := cryptoservice.GenerateCertificate(rootPrivateKey, gun, referenceTime, referenceTime.AddDate(10, 0, 0))
	require.NoError(t, err)
	newRootCertKey := trustmanager.CertToKey(newRootCert)
	require.NotEqual(t, oldRootCertKey.ID(), newRootCertKey.ID())

	// Only trust the new certificate
	err = repo.ReplaceBaseKeys(data.CanonicalRootRole, newRootCertKey)
	require.NoError(t, err)
	updatedRootRole, err := repo.GetBaseRole(data.CanonicalRootRole)
	require.NoError(t, err)
	updatedRootKeyIDs := updatedRootRole.ListKeyIDs()
	require.Equal(t, 1, len(updatedRootKeyIDs))
	require.Equal(t, newRootCertKey.ID(), updatedRootKeyIDs[0])

	// Now forget all about the old certificate: drop it from the Root carried keys
	delete(repo.Root.Signed.Keys, oldRootCertKey.ID())
	repo2 := NewRepo(cs)
	repo2.Root = repo.Root
	repo2.originalRootRole = updatedRootRole

	// Create a second signature
	signedRoot, err = repo2.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedRoot, newRootCertKey) // Without oldRootCertKey

	// Verify that the signature can be verified when trusting the new certificate
	err = verifyRootSignatureAgainstKey(t, signedRoot, newRootCertKey)
	require.NoError(t, err)
	err = verifyRootSignatureAgainstKey(t, signedRoot, oldRootCertKey)
	require.Error(t, err)
}

// SignRoot signs with all old roles with valid keys, and also optionally any old
// signatures we have keys for even if they aren't in an old root.  It ignores any
// root role whose version is higher than the current version.  If signing fails,
// it reverts back.
func TestSignRootOldRootRolesAndOldSigs(t *testing.T) {
	gun := "docker/test-sign-root"
	referenceTime := time.Now()

	cs := cryptoservice.NewCryptoService(trustmanager.NewKeyMemoryStore(
		passphrase.ConstantRetriever("password")))

	rootCertKeys := make([]data.PublicKey, 9)
	rootPrivKeys := make([]data.PrivateKey, cap(rootCertKeys))
	for i := 0; i < cap(rootCertKeys); i++ {
		rootPublicKey, err := cs.Create(data.CanonicalRootRole, gun, data.ECDSAKey)
		require.NoError(t, err)
		rootPrivateKey, _, err := cs.GetPrivateKey(rootPublicKey.ID())
		require.NoError(t, err)
		rootCert, err := cryptoservice.GenerateCertificate(rootPrivateKey, gun, referenceTime.AddDate(-9, 0, 0),
			referenceTime.AddDate(1, 0, 0))
		require.NoError(t, err)
		rootCertKeys[i] = trustmanager.CertToKey(rootCert)
		rootPrivKeys[i] = rootPrivateKey
	}

	repo := initRepoWithRoot(t, cs, rootCertKeys[6])
	// sign with key 0, which represents the key for the a version of the root we
	// no longer have a record of
	signedObj, err := repo.Root.ToSigned()
	require.NoError(t, err)
	signedObj, err = repo.sign(signedObj, nil, []data.PublicKey{rootCertKeys[0]})
	require.NoError(t, err)
	// should be signed with key 0
	verifySignatureList(t, signedObj, rootCertKeys[0])
	repo.Root.Signatures = signedObj.Signatures

	// bump root version and also add the above keys and extra roles to root
	repo.Root.Signed.Version = 6
	oldExpiry := repo.Root.Signed.Expires
	// add every key to the root's key list except 1
	for i, key := range rootCertKeys {
		if i != 1 {
			repo.Root.Signed.Keys[key.ID()] = key
		}
	}
	// invalid root role because key not included in the key map - valid root version name
	repo.Root.Signed.Roles["root.1"] = &data.RootRole{KeyIDs: []string{rootCertKeys[1].ID()}, Threshold: 1}
	// invalid root versions names, but valid roles
	repo.Root.Signed.Roles["2.root"] = &data.RootRole{KeyIDs: []string{rootCertKeys[2].ID()}, Threshold: 1}
	repo.Root.Signed.Roles["root3"] = &data.RootRole{KeyIDs: []string{rootCertKeys[3].ID()}, Threshold: 1}
	repo.Root.Signed.Roles["root.4a"] = &data.RootRole{KeyIDs: []string{rootCertKeys[4].ID()}, Threshold: 1}
	// valid old root role and version
	repo.Root.Signed.Roles["root.5"] = &data.RootRole{KeyIDs: []string{rootCertKeys[5].ID()}, Threshold: 1}
	// greater or equal to the current root version, so invalid name, but valid root role
	repo.Root.Signed.Roles["root.6"] = &data.RootRole{KeyIDs: []string{rootCertKeys[7].ID()}, Threshold: 1}

	lenRootRoles := len(repo.Root.Signed.Roles)

	// rotate the current key to key 8
	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalRootRole, rootCertKeys[8]))

	requiredKeys := []data.PrivateKey{
		rootPrivKeys[5], // we need an old valid root role - this was specified in root5
		rootPrivKeys[6], // we need the previous valid key prior to root rotation
		rootPrivKeys[8], // we need the new root key we've rotated to
	}

	for _, privKey := range requiredKeys {
		// if we can't sign with a previous root, we fail
		require.NoError(t, cs.RemoveKey(privKey.ID()))
		_, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
		require.Error(t, err)
		require.IsType(t, signed.ErrInsufficientSignatures{}, err)
		require.Contains(t, err.Error(), "signing keys not available")

		// add back for next test
		require.NoError(t, cs.AddKey(data.CanonicalRootRole, gun, privKey))
	}
	// we haven't saved any unsaved roles because there was an error signing,
	// nor have we bumped the version or altered the expiry
	require.Equal(t, 6, repo.Root.Signed.Version)
	require.Equal(t, oldExpiry, repo.Root.Signed.Expires)
	require.Len(t, repo.Root.Signed.Roles, lenRootRoles)

	// remove all the keys we don't need and demonstrate we can still sign
	for _, index := range []int{1, 2, 3, 4, 7} {
		require.NoError(t, cs.RemoveKey(rootPrivKeys[index].ID()))
	}

	// SignRoot will sign with all the old keys based on old root roles as well
	// as any old signatures
	signedObj, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	expectedSigningKeys := []data.PublicKey{
		rootCertKeys[0], // old signature key, not in any role
		rootCertKeys[5], // root.5 key which is valid
		rootCertKeys[6], // previous key before rotation,
		rootCertKeys[8], //  newly rotated key
	}
	verifySignatureList(t, signedObj, expectedSigningKeys...)
	// verify that we saved the previous root (which overwrote an invalid saved root),
	// since it wasn't in the list of old valid roots, and we didn't save the newest
	// role
	require.NotNil(t, repo.Root.Signed.Roles["root.6"])
	require.Equal(t, data.RootRole{KeyIDs: []string{rootCertKeys[6].ID()}, Threshold: 1},
		*repo.Root.Signed.Roles["root.6"])
	require.Nil(t, repo.Root.Signed.Roles["root.7"])

	// bumped version, 1 new roles, but one overwrote the previous root.6, so actually no
	// additional roles
	require.Equal(t, 7, repo.Root.Signed.Version)
	require.Len(t, repo.Root.Signed.Roles, lenRootRoles)
	require.True(t, oldExpiry.Before(repo.Root.Signed.Expires))
	lenRootRoles = len(repo.Root.Signed.Roles)

	// remove the optional key
	require.NoError(t, cs.RemoveKey(rootPrivKeys[0].ID()))

	// SignRoot will still succeed even if the key that wasn't in a root isn't
	// available
	oldExpiry = repo.Root.Signed.Expires
	signedObj, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedObj, expectedSigningKeys[1:]...)

	// no additional roles were added
	require.Len(t, repo.Root.Signed.Roles, lenRootRoles)
	require.Equal(t, 8, repo.Root.Signed.Version)               // bumped version
	require.True(t, oldExpiry.Before(repo.Root.Signed.Expires)) // expiry updated

	// now rotate a non-root key
	newTargetsKey, err := cs.Create(data.CanonicalTargetsRole, gun, data.ECDSAKey)
	require.NoError(t, err)
	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalTargetsRole, newTargetsKey))

	// we still sign with all old roles no additional roles were added
	oldExpiry = repo.Root.Signed.Expires
	signedObj, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.NoError(t, err)
	verifySignatureList(t, signedObj, expectedSigningKeys[1:]...)
	require.Len(t, repo.Root.Signed.Roles, lenRootRoles)
	require.Equal(t, 9, repo.Root.Signed.Version)               // bumped version
	require.True(t, oldExpiry.Before(repo.Root.Signed.Expires)) // expiry updated

	// rotating a targets key again, if we are missing the previous root's keys, signing will fail
	newTargetsKey, err = cs.Create(data.CanonicalTargetsRole, gun, data.ECDSAKey)
	require.NoError(t, err)
	require.NoError(t, repo.ReplaceBaseKeys(data.CanonicalTargetsRole, newTargetsKey))

	require.NoError(t, cs.RemoveKey(rootPrivKeys[6].ID()))

	oldExpiry = repo.Root.Signed.Expires
	_, err = repo.SignRoot(data.DefaultExpires(data.CanonicalRootRole))
	require.Error(t, err)
	require.IsType(t, signed.ErrInsufficientSignatures{}, err)
	require.Contains(t, err.Error(), "signing keys not available")

	// no additional roles were saved, version has not changed
	require.Len(t, repo.Root.Signed.Roles, lenRootRoles)
	require.Equal(t, 9, repo.Root.Signed.Version) // version has not changed
	require.Equal(t, oldExpiry, repo.Root.Signed.Expires)
}
