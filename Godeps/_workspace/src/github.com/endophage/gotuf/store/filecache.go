package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FileCacheStore implements a super simple wrapper around RemoteStore
// to handle local file caching of metadata
type FileCacheStore struct {
	RemoteStore
	cachePath string
}

func NewFileCacheStore(remote RemoteStore, cachePath string) *FileCacheStore {
	return &FileCacheStore{
		RemoteStore: remote,
		cachePath:   cachePath,
	}
}

func (s FileCacheStore) cacheFile(name string, data json.RawMessage) error {
	path := path.Join(s.cachePath, name)
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0600)
	return ioutil.WriteFile(path+".json", data, 0600)
}

func (s FileCacheStore) useCachedFile(name string) (json.RawMessage, error) {
	path := path.Join(s.cachePath, name+".json")
	return ioutil.ReadFile(path)
}

func (s FileCacheStore) GetMeta(name string, size int64) (json.RawMessage, error) {
	data, err := s.useCachedFile(name)
	if err == nil || data != nil {
		return data, nil
	}
	data, err = s.RemoteStore.GetMeta(name, size)
	if err != nil {
		return nil, err
	}
	s.cacheFile(name, data)
	return data, nil
}
