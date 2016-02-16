package data

import "fmt"

// ErrMissingMeta - couldn't find the FileMeta object for a role or target
type ErrMissingMeta struct {
	Role string
}

func (e ErrMissingMeta) Error() string {
	return fmt.Sprintf("tuf: sha256 checksum required for %s", e.Role)
}
