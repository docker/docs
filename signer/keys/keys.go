package keys

import "errors"

var (
	// ErrExists happens when a Key already exists in a database
	ErrExists = errors.New("notary-signer: key already in db")
	// ErrInvalidKeyID error happens when a key isn't found
	ErrInvalidKeyID = errors.New("notary-signer: invalid key id")
	// ErrFailedKeyGeneration happens when there is a failure in generating a key
	ErrFailedKeyGeneration = errors.New("notary-signer: failed to generate new key")
)
