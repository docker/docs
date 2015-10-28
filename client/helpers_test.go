package client

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/docker/notary/client/changelist"
	tuf "github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/stretchr/testify/assert"
)

func TestApplyTargetsChange(t *testing.T) {
	kdb := keys.NewDB()
	role, err := data.NewRole("targets", 1, nil, nil, nil)
	assert.NoError(t, err)
	kdb.AddRole(role)

	repo := tuf.NewRepo(kdb, nil)
	err = repo.InitTargets()
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

func TestApplyChangelist(t *testing.T) {
	kdb := keys.NewDB()
	role, err := data.NewRole("targets", 1, nil, nil, nil)
	assert.NoError(t, err)
	kdb.AddRole(role)

	repo := tuf.NewRepo(kdb, nil)
	err = repo.InitTargets()
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
	kdb := keys.NewDB()
	role, err := data.NewRole("targets", 1, nil, nil, nil)
	assert.NoError(t, err)
	kdb.AddRole(role)

	repo := tuf.NewRepo(kdb, nil)
	err = repo.InitTargets()
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
