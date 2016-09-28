package hubconfig

import (
	"fmt"
	"time"
)

type KeyValueStore interface {
	Get(key string) ([]byte, error)
	Put(key string, content []byte) error
	Delete(key string) error
	Lock(string, []byte, time.Duration) error
	List(string) ([]string, error)
	Ping() error
	SemSignal(string) error
	SemWait(string, <-chan struct{}) (bool, error)
}

type MockKeyValueStore struct {
	backingMap map[string][]byte
}

func NewMockKeyValueStore() *MockKeyValueStore {
	return &MockKeyValueStore{
		backingMap: make(map[string][]byte),
	}
}

func (mkvs *MockKeyValueStore) Get(key string) ([]byte, error) {
	value, ok := mkvs.backingMap[key]
	if !ok {
		return nil, fmt.Errorf("Could not retrieve key %v", key)
	}
	return value, nil
}

func (mkvs *MockKeyValueStore) Put(key string, content []byte) error {
	mkvs.backingMap[key] = content
	return nil
}

// XXX: doesn't actually lock anything for now. Implement with mutex if needed.
func (mkvs *MockKeyValueStore) Lock(key string, content []byte, ttl time.Duration) error {
	mkvs.backingMap[key] = content
	return nil
}

// TODO: implement this
func (mkvs *MockKeyValueStore) List(key string) ([]string, error) {
	return []string{}, nil
}

func (mkvs *MockKeyValueStore) Delete(key string) error {
	delete(mkvs.backingMap, key)
	return nil
}

func (mkvs *MockKeyValueStore) Ping() error {
	return nil
}

func (mkvs *MockKeyValueStore) SemSignal(key string) error {
	return nil
}

func (mkvs *MockKeyValueStore) SemWait(key string, ch <-chan struct{}) (bool, error) {
	return false, nil
}
