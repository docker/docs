package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOfflineStore(t *testing.T) {
	s := OfflineStore{}
	_, err := s.GetMeta("", 0)
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)

	err = s.SetMeta("", nil)
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)

	err = s.SetMultiMeta(nil)
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)

	_, err = s.GetKey("")
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)

	_, err = s.GetTarget("")
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)

	err = s.RemoveAll()
	require.Error(t, err)
	require.IsType(t, ErrOffline{}, err)
}

func TestErrOffline(t *testing.T) {
	var _ error = ErrOffline{}
}
