package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore(nil)
	_, err := s.GetMeta("nonexistent", 0)
	require.Error(t, err)
	require.IsType(t, ErrMetaNotFound{}, err)

	metaContent := []byte("content")
	metaSize := int64(7)
	err = s.SetMeta("exists", metaContent)
	require.NoError(t, err)

	meta, err := s.GetMeta("exists", metaSize)
	require.NoError(t, err)
	require.Equal(t, metaContent, meta)

	meta, err = s.GetMeta("exists", -1)
	require.NoError(t, err)
	require.Equal(t, metaContent, meta)

	err = s.RemoveAll()
	require.NoError(t, err)

	_, err = s.GetMeta("exists", 0)
	require.Error(t, err)
	require.IsType(t, ErrMetaNotFound{}, err)
}
