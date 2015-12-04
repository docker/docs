package signed

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

// ListKeys only returns the keys for that role
func TestListKeys(t *testing.T) {
	c := NewEd25519()
	tskey, err := c.Create(data.CanonicalTimestampRole, data.ED25519Key)
	assert.NoError(t, err)

	_, err = c.Create(data.CanonicalRootRole, data.ED25519Key)
	assert.NoError(t, err)

	tsKeys := c.ListKeys(data.CanonicalTimestampRole)
	assert.Len(t, tsKeys, 1)
	assert.Equal(t, tskey.ID(), tsKeys[0])

	assert.Len(t, c.ListKeys(data.CanonicalTargetsRole), 0)
}
