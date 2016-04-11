package storage

import (
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

func ConsistentEmptyGetCurrentTest(t *testing.T, s MetaStore) {
	_, byt, err := s.GetCurrent("testGUN", data.CanonicalRootRole)
	require.Nil(t, byt)
	require.Error(t, err, "There should be an error Getting an empty table")
	require.IsType(t, ErrNotFound{}, err, "Should get a not found error")
}

func ConsistentMissingTSAndSnapGetCurrentTest(t *testing.T, s MetaStore) {
	_, byt, err := s.GetCurrent("testGUN", data.CanonicalRootRole)
	require.Nil(t, byt)
	require.Error(t, err, "There should be an error because there is no timestamp or snapshot to use on GetCurrent")
}

func GetChecksumFoundTest(t *testing.T, s MetaStore, rec TUFFile) {
	_, _, err := s.GetChecksum("testGUN", data.CanonicalRootRole, rec.Sha256)
	require.NoError(t, err, "There should no error for GetChecksum")
}

func ConsistentGetCurrentFoundTest(t *testing.T, s MetaStore, rec TUFFile) {
	_, byt, err := s.GetCurrent("testGUN", rec.Role)
	require.NoError(t, err)
	require.Equal(t, rec.Data, byt)
}
