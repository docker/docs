package errors

import (
	"errors"
	"fmt"
	"time"
)

// ErrInitNotAllowed - repo has already been initialized
var ErrInitNotAllowed = errors.New("tuf: repository already initialized")

// ErrMissingMetadata - cannot find the file meta being requested.
// Specifically, could not find the FileMeta object in the expected
// location.
type ErrMissingMetadata struct {
	Name string
}

func (e ErrMissingMetadata) Error() string {
	return fmt.Sprintf("tuf: missing metadata %s", e.Name)
}

// ErrFileNotFound - could not find a file
type ErrFileNotFound struct {
	Path string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("tuf: file not found %s", e.Path)
}

// ErrInsufficientKeys - did not have enough keys to sign when requested
type ErrInsufficientKeys struct {
	Name string
}

func (e ErrInsufficientKeys) Error() string {
	return fmt.Sprintf("tuf: insufficient keys to sign %s", e.Name)
}

// ErrInsufficientSignatures - do not have enough signatures on a piece of
// metadata
type ErrInsufficientSignatures struct {
	Name string
	Err  error
}

func (e ErrInsufficientSignatures) Error() string {
	return fmt.Sprintf("tuf: insufficient signatures for %s: %s", e.Name, e.Err)
}

// ErrInvalidRole - role is wrong. Typically we're missing the public keys for it
type ErrInvalidRole struct {
	Role string
}

func (e ErrInvalidRole) Error() string {
	return fmt.Sprintf("tuf: invalid role %s", e.Role)
}

// ErrInvalidExpires - the expiry time for a metadata file is invalid
type ErrInvalidExpires struct {
	Expires time.Time
}

func (e ErrInvalidExpires) Error() string {
	return fmt.Sprintf("tuf: invalid expires: %s", e.Expires)
}

// ErrKeyNotFound - could not find a given key on a role
type ErrKeyNotFound struct {
	Role  string
	KeyID string
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf(`tuf: no key with id "%s" exists for the %s role`, e.KeyID, e.Role)
}

// ErrNotEnoughKeys - there are not enough keys to ever meet the signature threshold
type ErrNotEnoughKeys struct {
	Role      string
	Keys      int
	Threshold int
}

func (e ErrNotEnoughKeys) Error() string {
	return fmt.Sprintf("tuf: %s role has insufficient keys for threshold (has %d keys, threshold is %d)", e.Role, e.Keys, e.Threshold)
}

// ErrPassphraseRequired - a passphrase is needed and wasn't provided
type ErrPassphraseRequired struct {
	Role string
}

func (e ErrPassphraseRequired) Error() string {
	return fmt.Sprintf("tuf: a passphrase is required to access the encrypted %s keys file", e.Role)
}
