package storage

import (
	"fmt"
)

// ErrOldVersion is returned when a newer version of TUF metadada is already available
type ErrOldVersion struct{}

// ErrOldVersion is returned when a newer version of TUF metadada is already available
func (err ErrOldVersion) Error() string {
	return fmt.Sprintf("Error updating metadata. A newer version is already available")
}

// ErrNotFound is returned when TUF metadata isn't found for a specific record
type ErrNotFound struct{}

// Error implements error
func (err ErrNotFound) Error() string {
	return fmt.Sprintf("No record found")
}

// ErrTimestampKeyExists is returned when a timestamp key already exists
type ErrTimestampKeyExists struct {
	gun string
}

// ErrTimestampKeyExists is returned when a timestamp key already exists
func (err ErrTimestampKeyExists) Error() string {
	return fmt.Sprintf("Error, timestamp key already exists for %s", err.gun)
}

// ErrNoKey is returned when no timestamp key is found
type ErrNoKey struct {
	gun string
}

// ErrNoKey is returned when no timestamp key is found
func (err ErrNoKey) Error() string {
	return fmt.Sprintf("Error, no timestamp key found for %s", err.gun)
}
