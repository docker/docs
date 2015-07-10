package storage

import (
	"fmt"
	"strings"
	"sync"
)

type key struct {
	cipher string
	public []byte
}

type ver struct {
	version int
	data    []byte
}

// MemStorage is really just designed for dev and testing. It is very
// inefficient in many scenarios
type MemStorage struct {
	lock    sync.Mutex
	tufMeta map[string][]*ver
	tsKeys  map[string]*key
}

// NewMemStorage instantiates a memStorage instance
func NewMemStorage() *MemStorage {
	return &MemStorage{
		tufMeta: make(map[string][]*ver),
		tsKeys:  make(map[string]*key),
	}
}

// UpdateCurrent updates the meta data for a specific role
func (st *MemStorage) UpdateCurrent(gun, role string, version int, data []byte) error {
	id := entryKey(gun, role)
	st.lock.Lock()
	defer st.lock.Unlock()
	if space, ok := st.tufMeta[id]; ok {
		for _, v := range space {
			if v.version >= version {
				return &ErrOldVersion{}
			}
		}
	}
	st.tufMeta[id] = append(st.tufMeta[id], &ver{version: version, data: data})
	return nil
}

// GetCurrent returns the metadada for a given role, under a GUN
func (st *MemStorage) GetCurrent(gun, role string) (data []byte, err error) {
	id := entryKey(gun, role)
	st.lock.Lock()
	defer st.lock.Unlock()
	space, ok := st.tufMeta[id]
	if !ok {
		return nil, &ErrNotFound{}
	}
	return space[len(st.tufMeta[id])-1].data, nil
}

// Delete delets all the metadata for a given GUN
func (st *MemStorage) Delete(gun string) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	for k := range st.tufMeta {
		if strings.HasPrefix(k, gun) {
			delete(st.tufMeta, k)
		}
	}
	return nil
}

// GetTimestampKey returns the public key material of the timestamp key of a given gun
func (st *MemStorage) GetTimestampKey(gun string) (cipher string, public []byte, err error) {
	// no need for lock. It's ok to return nil if an update
	// wasn't observed
	k, ok := st.tsKeys[gun]
	if !ok {
		return "", nil, &ErrNoKey{gun: gun}
	}

	return k.cipher, k.public, nil
}

// SetTimestampKey sets a Timestamp key under a gun
func (st *MemStorage) SetTimestampKey(gun, cipher string, public []byte) error {
	k := &key{cipher: cipher, public: public}
	st.lock.Lock()
	defer st.lock.Unlock()
	if _, ok := st.tsKeys[gun]; ok {
		return &ErrTimestampKeyExists{gun: gun}
	}
	st.tsKeys[gun] = k
	return nil
}

func entryKey(gun, role string) string {
	return fmt.Sprintf("%s.%s", gun, role)
}
