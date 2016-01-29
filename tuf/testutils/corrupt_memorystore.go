package testutils

import (
	"github.com/docker/notary/tuf/store"
)

// CorruptingMemoryStore corrupts all data returned by GetMeta
type CorruptingMemoryStore struct {
	store.MemoryStore
}

func NewCorruptingMemoryStore(meta map[string][]byte, files map[string][]byte) *CorruptingMemoryStore {
	s := store.NewMemoryStore(meta, files)
	return &CorruptingMemoryStore{MemoryStore: *s}
}

func (cm CorruptingMemoryStore) GetMeta(name string, size int64) ([]byte, error) {
	d, err := cm.MemoryStore.GetMeta(name, size)
	if err != nil {
		return nil, err
	}
	d[0] = '}' // all our content is JSON so must start with {
	return d, err
}
