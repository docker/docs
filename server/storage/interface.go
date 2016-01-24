package storage

// KeyStore provides a minimal interface for managing key persistence
type KeyStore interface {
	// GetKey returns the algorithm and public key for the given GUN and role.
	// If the GUN+role don't exist, returns an error.
	GetKey(gun, role string) (algorithm string, public []byte, err error)

	// SetKey sets the algorithm and public key for the given GUN and role if
	// it doesn't already exist.  Otherwise an error is returned.
	SetKey(gun, role, algorithm string, public []byte) error
}

// MetaStore holds the methods that are used for a Metadata Store
type MetaStore interface {
	// UpdateCurrent adds new metadata version for the given GUN if and only
	// if it's a new role, or the version is greater than the current version
	// for the role. Otherwise an error is returned.
	UpdateCurrent(gun string, update MetaUpdate) error

	// UpdateMany adds multiple new metadata for the given GUN.  It can even
	// add multiple versions for the same role, so long as those versions are
	// all unique and greater than any current versions.  Otherwise,
	// none of the metadata is added, and an error is be returned.
	UpdateMany(gun string, updates []MetaUpdate) error

	// GetCurrent returns the data part of the metadata for the latest version
	// of the given GUN and role.  If there is no data for the given GUN and
	// role, an error is returned.
	GetCurrent(gun, tufRole string) (data []byte, err error)

	// GetChecksum return the given tuf role file for the GUN with the
	// provided checksum. If the given (gun, role, checksum) are not
	// found, it returns storage.ErrNotFound
	GetChecksum(gun, tufRole, checksum string) (data []byte, err error)

	// Delete removes all metadata for a given GUN.  It does not return an
	// error if no metadata exists for the given GUN.
	Delete(gun string) error

	KeyStore
}
