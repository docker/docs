package trustmanager

import (
	"errors"

	"github.com/docker/notary"
)

const (
	visible = notary.PubCertPerms
	private = notary.PrivKeyPerms
)

var (
	// ErrPathOutsideStore indicates that the returned path would be
	// outside the store
	ErrPathOutsideStore = errors.New("path outside file store")
)

// LimitedFileStore implements the bare bones primitives (no hierarchy)
type LimitedFileStore interface {
	Add(fileName string, data []byte) error
	Remove(fileName string) error
	Get(fileName string) ([]byte, error)
	ListFiles() []string
}

// FileStore is the interface for full-featured FileStores
type FileStore interface {
	LimitedFileStore

	RemoveDir(directoryName string) error
	GetPath(fileName string) (string, error)
	ListDir(directoryName string) []string
	BaseDir() string
}
