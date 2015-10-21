package client

import (
	"errors"
	"fmt"
)

// Simple client errors
var (
	ErrNoRootKeys       = errors.New("tuf: no root keys found in local meta store")
	ErrInsufficientKeys = errors.New("tuf: insufficient keys to meet threshold")
)

// ErrChecksumMismatch - a checksum failed verification
type ErrChecksumMismatch struct {
	role string
}

func (e ErrChecksumMismatch) Error() string {
	return fmt.Sprintf("tuf: checksum for %s did not match", e.role)
}

// ErrMissingMeta - couldn't find the FileMeta object for a role or target
type ErrMissingMeta struct {
	role string
}

func (e ErrMissingMeta) Error() string {
	return fmt.Sprintf("tuf: sha256 checksum required for %s", e.role)
}

// ErrMissingRemoteMetadata - remote didn't have requested metadata
type ErrMissingRemoteMetadata struct {
	Name string
}

func (e ErrMissingRemoteMetadata) Error() string {
	return fmt.Sprintf("tuf: missing remote metadata %s", e.Name)
}

// ErrDownloadFailed - a download failed
type ErrDownloadFailed struct {
	File string
	Err  error
}

func (e ErrDownloadFailed) Error() string {
	return fmt.Sprintf("tuf: failed to download %s: %s", e.File, e.Err)
}

// ErrDecodeFailed - couldn't parse a download
type ErrDecodeFailed struct {
	File string
	Err  error
}

func (e ErrDecodeFailed) Error() string {
	return fmt.Sprintf("tuf: failed to decode %s: %s", e.File, e.Err)
}

func isDecodeFailedWithErr(err, expected error) bool {
	e, ok := err.(ErrDecodeFailed)
	if !ok {
		return false
	}
	return e.Err == expected
}

// ErrNotFound - didn't find a file
type ErrNotFound struct {
	File string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("tuf: file not found: %s", e.File)
}

// IsNotFound - check if an error is an ErrNotFound type
func IsNotFound(err error) bool {
	_, ok := err.(ErrNotFound)
	return ok
}

// ErrWrongSize - the size is wrong
type ErrWrongSize struct {
	File     string
	Actual   int64
	Expected int64
}

func (e ErrWrongSize) Error() string {
	return fmt.Sprintf("tuf: unexpected file size: %s (expected %d bytes, got %d bytes)", e.File, e.Expected, e.Actual)
}

// ErrCorruptedCache - local data is incorrect
type ErrCorruptedCache struct {
	file string
}

func (e ErrCorruptedCache) Error() string {
	return fmt.Sprintf("cache is corrupted: %s", e.file)
}
