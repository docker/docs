package changelist

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/stretchr/testify/assert"
)

func TestTufDelegation(t *testing.T) {
	cs := signed.NewEd25519()
	key, err := cs.Create("targets/new_name", "gun", data.ED25519Key)
	assert.NoError(t, err)
	kl := data.KeyList{key}
	td := TufDelegation{
		NewName:      "targets/new_name",
		NewThreshold: 1,
		AddKeys:      kl,
		AddPaths:     []string{""},
	}

	r, err := td.ToNewRole("targets/old_name")
	assert.NoError(t, err)
	assert.Equal(t, td.NewName, r.Name)
	assert.Len(t, r.KeyIDs, 1)
	assert.Equal(t, kl[0].ID(), r.KeyIDs[0])
	assert.Len(t, r.Paths, 1)
}
