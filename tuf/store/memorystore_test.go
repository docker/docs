package store

import (
	"crypto/sha256"
	"testing"

	"github.com/docker/notary/tuf/utils"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreMetadataOperations(t *testing.T) {
	s := NewMemoryStore(nil)

	// GetMeta of a non-existent metadata fails
	_, err := s.GetMeta("nonexistent", 0)
	require.Error(t, err)
	require.IsType(t, ErrMetaNotFound{}, err)

	// Once SetMeta succeeds, GetMeta with the role name and the consistent name
	// should succeed
	metaContent := []byte("content")
	metaSize := int64(len(metaContent))
	shasum := sha256.Sum256(metaContent)
	invalidShasum := sha256.Sum256([]byte{})

	require.NoError(t, s.SetMeta("exists", metaContent))
	require.NoError(t, s.SetMultiMeta(map[string][]byte{"multi1": metaContent, "multi2": metaContent}))

	for _, metaName := range []string{"exists", "multi1", "multi2"} {
		meta, err := s.GetMeta(metaName, metaSize)
		require.NoError(t, err)
		require.Equal(t, metaContent, meta)

		meta, err = s.GetMeta(utils.ConsistentName(metaName, shasum[:]), metaSize)
		require.NoError(t, err)
		require.Equal(t, metaContent, meta)

		_, err = s.GetMeta(utils.ConsistentName(metaName, invalidShasum[:]), metaSize)
		require.Error(t, err)
		require.IsType(t, ErrMetaNotFound{}, err)
	}

	// Once Metadata is removed, it's no longer accessible
	err = s.RemoveAll()
	require.NoError(t, err)

	_, err = s.GetMeta("exists", 0)
	require.Error(t, err)
	require.IsType(t, ErrMetaNotFound{}, err)
}

func TestMemoryStoreGetMetaSize(t *testing.T) {
	content := []byte("content")
	s := NewMemoryStore(map[string][]byte{"content": content})

	// we can get partial size
	meta, err := s.GetMeta("content", 3)
	require.NoError(t, err)
	require.Equal(t, []byte("con"), meta)

	// we can get zero size
	meta, err = s.GetMeta("content", 0)
	require.NoError(t, err)
	require.Equal(t, []byte{}, meta)

	// we can get the whole thing by passing NoSizeLimit (-1)
	meta, err = s.GetMeta("content", NoSizeLimit)
	require.NoError(t, err)
	require.Equal(t, content, meta)

	// a size much larger than the actual length will return the whole thing
	meta, err = s.GetMeta("content", 8000)
	require.NoError(t, err)
	require.Equal(t, content, meta)
}
