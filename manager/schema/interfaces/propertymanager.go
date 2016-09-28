package interfaces

import (
	"errors"
)

// ErrPropertyNotSet conveys that the given property is not set in the database.
var ErrPropertyNotSet = errors.New("property not set")

// TODO(andrewnguyen): Merge this with hubconfig.KeyValueStore? This is silly
type PropertyManager interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type MockPropertyManager struct {
	backingMap map[string]string
}

func NewMockPropertyManager() *MockPropertyManager {
	return &MockPropertyManager{
		backingMap: make(map[string]string),
	}
}

func (m *MockPropertyManager) Set(key, value string) error {
	m.backingMap[key] = value
	return nil
}

func (m *MockPropertyManager) Get(key string) (string, error) {
	value, ok := m.backingMap[key]
	if !ok {
		return "", ErrPropertyNotSet
	}
	return value, nil
}
