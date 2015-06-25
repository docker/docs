package storage

import (
	"fmt"
)

type ErrOldVersion struct{}

func (err ErrOldVersion) Error() string {
	return fmt.Sprintf("Error updating metadata. A newer version is already available")
}

type ErrNotFound struct{}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("No record found")
}

type ErrTimestampKeyExists struct {
	gun string
}

func (err ErrTimestampKeyExists) Error() string {
	return fmt.Sprintf("Error, timestamp key already exists for %s", err.gun)
}

type ErrNoKey struct {
	gun string
}

func (err ErrNoKey) Error() string {
	return fmt.Sprintf("Error, no timestamp key found for %s", err.gun)
}
