package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

type key struct {
	algorithm string
	public    []byte
}

type ver struct {
	version      int
	data         []byte
	createupdate time.Time
}

// MemStorage is really just designed for dev and testing. It is very
// inefficient in many scenarios
type MemStorage struct {
	lock      sync.Mutex
	tufMeta   map[string][]*ver
	keys      map[string]map[string]*key
	checksums map[string]map[string]ver
}

// NewMemStorage instantiates a memStorage instance
func NewMemStorage() *MemStorage {
	return &MemStorage{
		tufMeta:   make(map[string][]*ver),
		keys:      make(map[string]map[string]*key),
		checksums: make(map[string]map[string]ver),
	}
}

// UpdateCurrent updates the meta data for a specific role
func (st *MemStorage) UpdateCurrent(gun string, update MetaUpdate) error {
	id := entryKey(gun, update.Role)
	st.lock.Lock()
	defer st.lock.Unlock()
	if space, ok := st.tufMeta[id]; ok {
		for _, v := range space {
			if v.version >= update.Version {
				return &ErrOldVersion{}
			}
		}
	}
	version := ver{version: update.Version, data: update.Data, createupdate: time.Now()}
	st.tufMeta[id] = append(st.tufMeta[id], &version)
	checksumBytes := sha256.Sum256(update.Data)
	checksum := hex.EncodeToString(checksumBytes[:])

	_, ok := st.checksums[gun]
	if !ok {
		st.checksums[gun] = make(map[string]ver)
	}
	st.checksums[gun][checksum] = version
	return nil
}

// UpdateMany updates multiple TUF records
func (st *MemStorage) UpdateMany(gun string, updates []MetaUpdate) error {
	for _, u := range updates {
		if err := st.UpdateCurrent(gun, u); err != nil {
			return err
		}
	}
	return nil
}

// GetCurrent returns the createupdate date metadata for a given role, under a GUN.
func (st *MemStorage) GetCurrent(gun, role string) (*time.Time, []byte, error) {
	id := entryKey(gun, role)
	st.lock.Lock()
	defer st.lock.Unlock()
	space, ok := st.tufMeta[id]
	if !ok || len(space) == 0 {
		return nil, nil, ErrNotFound{}
	}
	return &(space[len(space)-1].createupdate), space[len(space)-1].data, nil
}

// GetChecksum returns the createupdate date and metadata for a given role, under a GUN.
func (st *MemStorage) GetChecksum(gun, role, checksum string) (*time.Time, []byte, error) {
	st.lock.Lock()
	defer st.lock.Unlock()
	space, ok := st.checksums[gun][checksum]
	if !ok || len(space.data) == 0 {
		return nil, nil, ErrNotFound{}
	}
	return &(space.createupdate), space.data, nil
}

// Delete deletes all the metadata for a given GUN
func (st *MemStorage) Delete(gun string) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	for k := range st.tufMeta {
		if strings.HasPrefix(k, gun) {
			delete(st.tufMeta, k)
		}
	}
	delete(st.checksums, gun)
	return nil
}

// GetKey returns the public key material of the timestamp key of a given gun
func (st *MemStorage) GetKey(gun, role string) (algorithm string, public []byte, err error) {
	// no need for lock. It's ok to return nil if an update
	// wasn't observed
	g, ok := st.keys[gun]
	if !ok {
		return "", nil, &ErrNoKey{gun: gun}
	}
	k, ok := g[role]
	if !ok {
		return "", nil, &ErrNoKey{gun: gun}
	}

	return k.algorithm, k.public, nil
}

// SetKey sets a key under a gun and role
func (st *MemStorage) SetKey(gun, role, algorithm string, public []byte) error {
	k := &key{algorithm: algorithm, public: public}
	st.lock.Lock()
	defer st.lock.Unlock()

	// we hold the lock so nothing will be able to race to write a key
	// between checking and setting
	_, _, err := st.GetKey(gun, role)
	if _, ok := err.(*ErrNoKey); !ok {
		return &ErrKeyExists{gun: gun, role: role}
	}
	_, ok := st.keys[gun]
	if !ok {
		st.keys[gun] = make(map[string]*key)
	}
	st.keys[gun][role] = k
	return nil
}

func entryKey(gun, role string) string {
	return fmt.Sprintf("%s.%s", gun, role)
}
