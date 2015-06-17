package store

import (
	"bytes"
	"encoding/json"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/errors"
)

func MemoryStore(meta map[string]json.RawMessage, files map[string][]byte) LocalStore {
	if meta == nil {
		meta = make(map[string]json.RawMessage)
	}
	return &memoryStore{
		meta:  meta,
		files: files,
		keys:  make(map[string][]*data.Key),
	}
}

type memoryStore struct {
	meta  map[string]json.RawMessage
	files map[string][]byte
	keys  map[string][]*data.Key
}

func (m *memoryStore) GetMeta(name string, size int64) (json.RawMessage, error) {
	return m.meta[name], nil
}

func (m *memoryStore) SetMeta(name string, meta json.RawMessage) error {
	m.meta[name] = meta
	return nil
}

func (m *memoryStore) AddBlob(path string, meta data.FileMeta) {

}

func (m *memoryStore) WalkStagedTargets(paths []string, targetsFn targetsWalkFunc) error {
	if len(paths) == 0 {
		for path, dat := range m.files {
			meta, err := data.NewFileMeta(bytes.NewReader(dat), "sha256")
			if err != nil {
				return err
			}
			if err = targetsFn(path, meta); err != nil {
				return err
			}
		}
		return nil
	}

	for _, path := range paths {
		dat, ok := m.files[path]
		if !ok {
			return errors.ErrFileNotFound{path}
		}
		meta, err := data.NewFileMeta(bytes.NewReader(dat), "sha256")
		if err != nil {
			return err
		}
		if err = targetsFn(path, meta); err != nil {
			return err
		}
	}
	return nil
}

func (m *memoryStore) Commit(map[string]json.RawMessage, bool, map[string]data.Hashes) error {
	return nil
}

func (m *memoryStore) GetKeys(role string) ([]*data.Key, error) {
	return m.keys[role], nil
}

func (m *memoryStore) SaveKey(role string, key *data.Key) error {
	if _, ok := m.keys[role]; !ok {
		m.keys[role] = make([]*data.Key, 0)
	}
	m.keys[role] = append(m.keys[role], key)
	return nil
}

func (m *memoryStore) Clean() error {
	return nil
}
