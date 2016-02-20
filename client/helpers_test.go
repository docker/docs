package client

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/testutils"
	"github.com/stretchr/testify/assert"
)

func TestApplyTargetsChange(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	_, err = repo.InitTargets(data.CanonicalTargetsRole)
	assert.NoError(t, err)
	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	addChange := &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}
	err = applyTargetsChange(repo, addChange)
	assert.NoError(t, err)
	assert.NotNil(t, repo.Targets["targets"].Signed.Targets["latest"])

	removeChange := &changelist.TufChange{
		Actn:       changelist.ActionDelete,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       nil,
	}
	err = applyTargetsChange(repo, removeChange)
	assert.NoError(t, err)
	_, ok := repo.Targets["targets"].Signed.Targets["latest"]
	assert.False(t, ok)
}

// Adding the same target twice doesn't actually add it.
func TestApplyAddTargetTwice(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	_, err = repo.InitTargets(data.CanonicalTargetsRole)
	assert.NoError(t, err)
	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))

	assert.NoError(t, applyChangelist(repo, cl))
	assert.Len(t, repo.Targets["targets"].Signed.Targets, 1)
	assert.NotEmpty(t, repo.Targets["targets"].Signed.Targets["latest"])

	assert.NoError(t, applyTargetsChange(repo, &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))
	assert.Len(t, repo.Targets["targets"].Signed.Targets, 1)
	assert.NotEmpty(t, repo.Targets["targets"].Signed.Targets["latest"])
}

func TestApplyChangelist(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	_, err = repo.InitTargets(data.CanonicalTargetsRole)
	assert.NoError(t, err)
	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	addChange := &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}
	cl.Add(addChange)
	err = applyChangelist(repo, cl)
	assert.NoError(t, err)
	assert.NotNil(t, repo.Targets["targets"].Signed.Targets["latest"])

	cl.Clear("")

	removeChange := &changelist.TufChange{
		Actn:       changelist.ActionDelete,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       nil,
	}
	cl.Add(removeChange)
	err = applyChangelist(repo, cl)
	assert.NoError(t, err)
	_, ok := repo.Targets["targets"].Signed.Targets["latest"]
	assert.False(t, ok)
}

func TestApplyChangelistMulti(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)
	_, err = repo.InitTargets(data.CanonicalTargetsRole)
	assert.NoError(t, err)
	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	addChange := &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}

	removeChange := &changelist.TufChange{
		Actn:       changelist.ActionDelete,
		Role:       changelist.ScopeTargets,
		ChangeType: "target",
		ChangePath: "latest",
		Data:       nil,
	}

	cl.Add(addChange)
	cl.Add(removeChange)

	err = applyChangelist(repo, cl)
	assert.NoError(t, err)
	_, ok := repo.Targets["targets"].Signed.Targets["latest"]
	assert.False(t, ok)
}

func TestApplyTargetsDelegationCreateDelete(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	tgts := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, tgts.Signed.Delegations.Roles, 1)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	_, ok := tgts.Signed.Delegations.Keys[newKey.ID()]
	assert.True(t, ok)

	role := tgts.Signed.Delegations.Roles[0]
	assert.Len(t, role.KeyIDs, 1)
	assert.Equal(t, newKey.ID(), role.KeyIDs[0])
	assert.Equal(t, "targets/level1", role.Name)
	assert.Equal(t, "level1", role.Paths[0])

	// delete delegation
	td = &changelist.TufDelegation{
		RemoveKeys: []string{newKey.ID()},
	}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)
	ch = changelist.NewTufChange(
		changelist.ActionDelete,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	assert.Len(t, tgts.Signed.Delegations.Roles, 0)
	assert.Len(t, tgts.Signed.Delegations.Keys, 0)
}

func TestApplyTargetsDelegationCreate2SharedKey(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create first delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	// create second delegation
	kl = data.KeyList{newKey}
	td = &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level2"},
	}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)

	ch = changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level2",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	tgts := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, tgts.Signed.Delegations.Roles, 2)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	role1 := tgts.Signed.Delegations.Roles[0]
	assert.Len(t, role1.KeyIDs, 1)
	assert.Equal(t, newKey.ID(), role1.KeyIDs[0])
	assert.Equal(t, "targets/level1", role1.Name)
	assert.Equal(t, "level1", role1.Paths[0])

	role2 := tgts.Signed.Delegations.Roles[1]
	assert.Len(t, role2.KeyIDs, 1)
	assert.Equal(t, newKey.ID(), role2.KeyIDs[0])
	assert.Equal(t, "targets/level2", role2.Name)
	assert.Equal(t, "level2", role2.Paths[0])

	// delete one delegation, ensure shared key remains
	td = &changelist.TufDelegation{
		RemoveKeys: []string{newKey.ID()},
	}
	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)
	ch = changelist.NewTufChange(
		changelist.ActionDelete,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	assert.Len(t, tgts.Signed.Delegations.Roles, 1)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	// delete other delegation, ensure key cleaned up
	ch = changelist.NewTufChange(
		changelist.ActionDelete,
		"targets/level2",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	assert.Len(t, tgts.Signed.Delegations.Roles, 0)
	assert.Len(t, tgts.Signed.Delegations.Keys, 0)
}

func TestApplyTargetsDelegationCreateEdit(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	// edit delegation
	newKey2, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	kl = data.KeyList{newKey2}
	td = &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		RemoveKeys:   []string{newKey.ID()},
	}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)

	ch = changelist.NewTufChange(
		changelist.ActionUpdate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	tgts := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, tgts.Signed.Delegations.Roles, 1)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	_, ok := tgts.Signed.Delegations.Keys[newKey2.ID()]
	assert.True(t, ok)

	role := tgts.Signed.Delegations.Roles[0]
	assert.Len(t, role.KeyIDs, 1)
	assert.Equal(t, newKey2.ID(), role.KeyIDs[0])
	assert.Equal(t, "targets/level1", role.Name)
	assert.Equal(t, "level1", role.Paths[0])
}

func TestApplyTargetsDelegationEditNonExisting(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionUpdate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

func TestApplyTargetsDelegationCreateAlreadyExisting(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)
	// we have sufficient checks elsewhere we don't need to confirm that
	// creating fresh works here via more asserts.

	extraKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl = data.KeyList{extraKey}
	td = &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)

	ch = changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	// when attempting to create the same role again, check that we added a key
	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)
	delegation, err := repo.GetDelegationRole("targets/level1")
	assert.NoError(t, err)
	assert.Contains(t, delegation.Paths, "level1")
	assert.Equal(t, len(delegation.ListKeyIDs()), 2)
}

func TestApplyTargetsDelegationAlreadyExistingMergePaths(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)
	// we have sufficient checks elsewhere we don't need to confirm that
	// creating fresh works here via more asserts.

	// Use different path for this changelist
	td.AddPaths = []string{"level2"}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)

	ch = changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	// when attempting to create the same role again, check that we
	// merged with previous details
	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)
	delegation, err := repo.GetDelegationRole("targets/level1")
	assert.NoError(t, err)
	// Assert we have both paths
	assert.Contains(t, delegation.Paths, "level2")
	assert.Contains(t, delegation.Paths, "level1")
}

func TestApplyTargetsDelegationInvalidRole(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"bad role",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
}

func TestApplyTargetsDelegationInvalidJSONContent(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON[1:],
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
}

func TestApplyTargetsDelegationInvalidAction(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		"bad action",
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		nil,
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
}

func TestApplyTargetsChangeInvalidType(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		"badType",
		"",
		nil,
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
}

func TestApplyTargetsDelegationCreate2Deep(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1"},
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	tgts := repo.Targets[data.CanonicalTargetsRole]
	assert.Len(t, tgts.Signed.Delegations.Roles, 1)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	_, ok := tgts.Signed.Delegations.Keys[newKey.ID()]
	assert.True(t, ok)

	role := tgts.Signed.Delegations.Roles[0]
	assert.Len(t, role.KeyIDs, 1)
	assert.Equal(t, newKey.ID(), role.KeyIDs[0])
	assert.Equal(t, "targets/level1", role.Name)
	assert.Equal(t, "level1", role.Paths[0])

	// init delegations targets file. This would be done as part of a publish
	// operation
	repo.InitTargets("targets/level1")

	td = &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{"level1/level2"},
	}

	tdJSON, err = json.Marshal(td)
	assert.NoError(t, err)

	ch = changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1/level2",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.NoError(t, err)

	tgts = repo.Targets["targets/level1"]
	assert.Len(t, tgts.Signed.Delegations.Roles, 1)
	assert.Len(t, tgts.Signed.Delegations.Keys, 1)

	_, ok = tgts.Signed.Delegations.Keys[newKey.ID()]
	assert.True(t, ok)

	role = tgts.Signed.Delegations.Roles[0]
	assert.Len(t, role.KeyIDs, 1)
	assert.Equal(t, newKey.ID(), role.KeyIDs[0])
	assert.Equal(t, "targets/level1/level2", role.Name)
	assert.Equal(t, "level1/level2", role.Paths[0])
}

// Applying a delegation whose parent doesn't exist fails.
func TestApplyTargetsDelegationParentDoesntExist(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	// make sure a key exists for the previous level, so it's not a missing
	// key error, but we don't care about this key
	_, err = cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1/level2", data.ED25519Key)
	assert.NoError(t, err)

	// create delegation
	kl := data.KeyList{newKey}
	td := &changelist.TufDelegation{
		NewThreshold: 1,
		AddKeys:      kl,
	}

	tdJSON, err := json.Marshal(td)
	assert.NoError(t, err)

	ch := changelist.NewTufChange(
		changelist.ActionCreate,
		"targets/level1/level2",
		changelist.TypeTargetsDelegation,
		"",
		tdJSON,
	)

	err = applyTargetsChange(repo, ch)
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// If there is no delegation target, ApplyTargets creates it
func TestApplyChangelistCreatesDelegation(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	err = repo.UpdateDelegations("targets/level1", changelist.TufDelegation{AddKeys: []data.PublicKey{newKey}, AddPaths: []string{""}, NewThreshold: 1})
	assert.NoError(t, err)
	delete(repo.Targets, "targets/level1")

	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       "targets/level1",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))

	assert.NoError(t, applyChangelist(repo, cl))
	_, ok := repo.Targets["targets/level1"]
	assert.True(t, ok, "Failed to create the delegation target")
	_, ok = repo.Targets["targets/level1"].Signed.Targets["latest"]
	assert.True(t, ok, "Failed to write change to delegation target")
}

// Each change applies only to the role specified
func TestApplyChangelistTargetsToMultipleRoles(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	err = repo.UpdateDelegations("targets/level1", changelist.TufDelegation{AddKeys: []data.PublicKey{newKey}, AddPaths: []string{""}, NewThreshold: 1})
	assert.NoError(t, err)

	err = repo.UpdateDelegations("targets/level2", changelist.TufDelegation{AddKeys: []data.PublicKey{newKey}, AddPaths: []string{""}, NewThreshold: 1})
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       "targets/level1",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionDelete,
		Role:       "targets/level2",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       nil,
	}))

	assert.NoError(t, applyChangelist(repo, cl))
	_, ok := repo.Targets["targets/level1"].Signed.Targets["latest"]
	assert.True(t, ok)
	_, ok = repo.Targets["targets/level2"]
	assert.False(t, ok, "no change to targets/level2, so metadata not created")
}

// ApplyTargets falls back to role that exists when adding or deleting a change
func TestApplyChangelistTargetsFallbackRoles(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	cl := changelist.NewMemChangelist()
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       "targets/level1/level2/level3/level4",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	}))

	assert.NoError(t, applyChangelist(repo, cl))
	_, ok := repo.Targets[data.CanonicalTargetsRole].Signed.Targets["latest"]
	assert.True(t, ok)

	// now delete and assert it applies to
	cl = changelist.NewMemChangelist()
	assert.NoError(t, cl.Add(&changelist.TufChange{
		Actn:       changelist.ActionDelete,
		Role:       "targets/level1/level2/level3/level4",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       nil,
	}))

	assert.NoError(t, applyChangelist(repo, cl))
	assert.Empty(t, repo.Targets[data.CanonicalTargetsRole].Signed.Targets)
}

// changeTargetMeta fallback fails with ErrInvalidRole if role is invalid
func TestChangeTargetMetaFallbackFailsInvalidRole(t *testing.T) {
	repo, _, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	err = changeTargetMeta(repo, &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       "ruhroh",
		ChangeType: "target",
		ChangePath: "latest",
		Data:       fjson,
	})
	assert.Error(t, err)
	assert.IsType(t, data.ErrInvalidRole{}, err)
}

// If applying a change fails due to a prefix error, it does not fall back
// on the parent.
func TestChangeTargetMetaDoesntFallbackIfPrefixError(t *testing.T) {
	repo, cs, err := testutils.EmptyRepo("docker.com/notary")
	assert.NoError(t, err)

	newKey, err := cs.Create("targets/level1", data.ED25519Key)
	assert.NoError(t, err)

	err = repo.UpdateDelegations("targets/level1", changelist.TufDelegation{AddKeys: []data.PublicKey{newKey}, AddPaths: []string{"pathprefix"}, NewThreshold: 1})
	assert.NoError(t, err)

	hash := sha256.Sum256([]byte{})
	f := &data.FileMeta{
		Length: 1,
		Hashes: map[string][]byte{
			"sha256": hash[:],
		},
	}
	fjson, err := json.Marshal(f)
	assert.NoError(t, err)

	err = changeTargetMeta(repo, &changelist.TufChange{
		Actn:       changelist.ActionCreate,
		Role:       "targets/level1",
		ChangeType: "target",
		ChangePath: "notPathPrefix",
		Data:       fjson,
	})
	assert.Error(t, err)

	// no target in targets or targets/latest
	assert.Empty(t, repo.Targets[data.CanonicalTargetsRole].Signed.Targets)
	assert.Empty(t, repo.Targets["targets/level1"].Signed.Targets)
}
