package signed

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

// ListKeys only returns the keys for that role
func TestListKeys(t *testing.T) {
	c := NewEd25519()
	tskey, err := c.Create(data.CanonicalTimestampRole, "", data.ED25519Key)
	require.NoError(t, err)

	_, err = c.Create(data.CanonicalRootRole, "", data.ED25519Key)
	require.NoError(t, err)

	tsKeys := c.ListKeys(data.CanonicalTimestampRole)
	require.Len(t, tsKeys, 1)
	require.Equal(t, tskey.ID(), tsKeys[0])

	require.Len(t, c.ListKeys(data.CanonicalTargetsRole), 0)
}
