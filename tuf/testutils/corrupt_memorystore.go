package testutils

import (
	"github.com/docker/notary/tuf/store"
)

// CorruptingMemoryStore corrupts all data returned by GetMeta
type CorruptingMemoryStore struct {
	store.MemoryStore
}

// NewCorruptingMemoryStore returns a new instance of memory store that
// corrupts all data requested from it.
func NewCorruptingMemoryStore(meta map[string][]byte, files map[string][]byte) *CorruptingMemoryStore {
	s := store.NewMemoryStore(meta, files)
	return &CorruptingMemoryStore{MemoryStore: *s}
}

// GetMeta returns up to size bytes of meta identified by string. It will
// always be corrupted by setting the first character to }
func (cm CorruptingMemoryStore) GetMeta(name string, size int64) ([]byte, error) {
	d, err := cm.MemoryStore.GetMeta(name, size)
	if err != nil {
		return nil, err
	}
	d[0] = '}' // all our content is JSON so must start with {
	return d, err
}
