package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteMeta(t *testing.T) {
	sp := &SignedSnapshot{
		Signatures: make([]Signature, 0),
		Signed: Snapshot{
			Type:    TUFTypes["snapshot"],
			Version: 0,
			Expires: DefaultExpires("snapshot"),
			Meta: Files{
				CanonicalRootRole:    FileMeta{},
				CanonicalTargetsRole: FileMeta{},
			},
		},
	}
	_, ok := sp.Signed.Meta["root"]
	assert.True(t, ok)
	sp.DeleteMeta("root")
	_, ok = sp.Signed.Meta["root"]
	assert.False(t, ok)
	assert.True(t, sp.Dirty)
}
