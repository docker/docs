package data

import "fmt"

// ErrInvalidMeta is the error to be returned when metadata is invalid
type ErrInvalidMeta struct {
	Role string
	Msg  string
}

func (e ErrInvalidMeta) Error() string {
	return fmt.Sprintf("%s type metadata invalid: %s", e.Role, e.Msg)
}

// ErrMissingMeta - couldn't find the FileMeta object for a role or target
type ErrMissingMeta struct {
	Role string
}

func (e ErrMissingMeta) Error() string {
	return fmt.Sprintf("tuf: sha256 checksum required for %s", e.Role)
}
