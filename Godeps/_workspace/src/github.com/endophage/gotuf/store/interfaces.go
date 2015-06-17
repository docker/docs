package store

import (
	"encoding/json"
	"io"

	"github.com/endophage/gotuf/data"
)

type targetsWalkFunc func(path string, meta data.FileMeta) error

type MetadataStore interface {
	GetMeta(name string, size int64) (json.RawMessage, error)
	SetMeta(name string, blob json.RawMessage) error
}

// [endophage] I'm of the opinion this should go away.
type TargetStore interface {
	WalkStagedTargets(paths []string, targetsFn targetsWalkFunc) error
}

type LocalStore interface {
	MetadataStore
	TargetStore
}

type RemoteStore interface {
	MetadataStore
	GetTarget(path string) (io.ReadCloser, error)
}
