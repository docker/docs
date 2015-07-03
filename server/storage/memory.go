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

// memStorage is really just designed for dev and testing. It is very
// inefficient in many scenarios
type memStorage struct {
	lock    sync.Mutex
	tufMeta map[string][]*ver
	tsKeys  map[string]*key
}

// NewMemStorage instantiates a memStorage instance
func NewMemStorage() *memStorage {
	return &memStorage{
		tufMeta: make(map[string][]*ver),
		tsKeys:  make(map[string]*key),
	}
}

func (st *memStorage) UpdateCurrent(gun, role string, version int, data []byte) error {
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

func (st *memStorage) GetCurrent(gun, role string) (data []byte, err error) {
	id := entryKey(gun, role)
	st.lock.Lock()
	defer st.lock.Unlock()
	space, ok := st.tufMeta[id]
	if !ok {
		return nil, &ErrNotFound{}
	}
	return space[len(st.tufMeta[id])-1].data, nil
}

func (st *memStorage) Delete(gun string) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	for k, _ := range st.tufMeta {
		if strings.HasPrefix(k, gun) {
			delete(st.tufMeta, k)
		}
	}
	return nil
}

func (st *memStorage) GetTimestampKey(gun string) (cipher string, public []byte, err error) {
	// no need for lock. It's ok to return nil if an update
	// wasn't observed
	if k, ok := st.tsKeys[gun]; !ok {
		return "", nil, &ErrNoKey{gun: gun}
	} else {
		return k.cipher, k.public, nil
	}
}

func (st *memStorage) SetTimestampKey(gun, cipher string, public []byte) error {
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
