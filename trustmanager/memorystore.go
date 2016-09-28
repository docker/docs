package trustmanager

import "os"

// MemoryFileStore is an implementation of Storage that keeps the contents in
// memory. It is not thread-safe and should be used by a higher-level interface
// that provides locking.
type MemoryFileStore struct {
	files map[string][]byte
}

// NewMemoryFileStore creates a MemoryFileStore
func NewMemoryFileStore() *MemoryFileStore {
	return &MemoryFileStore{
		files: make(map[string][]byte),
	}
}

// Add writes data to a file with a given name
func (f *MemoryFileStore) Add(name string, data []byte) error {
	f.files[name] = data
	return nil
}

// Remove removes a file identified by name
func (f *MemoryFileStore) Remove(name string) error {
	if _, present := f.files[name]; !present {
		return os.ErrNotExist
	}
	delete(f.files, name)

	return nil
}

// Get returns the data given a file name
func (f *MemoryFileStore) Get(name string) ([]byte, error) {
	fileData, present := f.files[name]
	if !present {
		return nil, os.ErrNotExist
	}

	return fileData, nil
}

// ListFiles lists all the files inside of a store
func (f *MemoryFileStore) ListFiles() []string {
	var list []string

	for name := range f.files {
		list = append(list, name)
	}

	return list
}
