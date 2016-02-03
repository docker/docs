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
func NewCorruptingMemoryStore(meta map[string][]byte) *CorruptingMemoryStore {
	s := store.NewMemoryStore(meta)
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

// LongMemoryStore corrupts all data returned by GetMeta
type LongMemoryStore struct {
	store.MemoryStore
}

// NewLongMemoryStore returns a new instance of memory store that
// returns one byte too much data on any request to GetMeta
func NewLongMemoryStore(meta map[string][]byte) *LongMemoryStore {
	s := store.NewMemoryStore(meta)
	return &LongMemoryStore{MemoryStore: *s}
}

// GetMeta returns one byte too much
func (lm LongMemoryStore) GetMeta(name string, size int64) ([]byte, error) {
	d, err := lm.MemoryStore.GetMeta(name, size)
	if err != nil {
		return nil, err
	}
	d = append(d, ' ')
	return d, err
}

// ShortMemoryStore corrupts all data returned by GetMeta
type ShortMemoryStore struct {
	store.MemoryStore
}

// NewShortMemoryStore returns a new instance of memory store that
// returns one byte too little data on any request to GetMeta
func NewShortMemoryStore(meta map[string][]byte) *ShortMemoryStore {
	s := store.NewMemoryStore(meta)
	return &ShortMemoryStore{MemoryStore: *s}
}

// GetMeta returns one byte too few
func (sm ShortMemoryStore) GetMeta(name string, size int64) ([]byte, error) {
	d, err := sm.MemoryStore.GetMeta(name, size)
	if err != nil {
		return nil, err
	}
	return d[1:], err
}
