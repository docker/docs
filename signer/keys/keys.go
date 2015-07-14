package keys

import (
	"errors"

	"github.com/agl/ed25519"
	"github.com/miekg/pkcs11"
)

var (
	// ErrExists happens when a Key already exists in a database
	ErrExists = errors.New("rufus: key already in db")
	// ErrInvalidKeyID error happens when a key isn't found
	ErrInvalidKeyID = errors.New("rufus: invalid key id")
	// ErrFailedKeyGeneration happens when there is a failure in generating a key
	ErrFailedKeyGeneration = errors.New("rufus: failed to generate new key")
)

// Key represents all the information of a key, including the private and public bits
type Key struct {
	ID        string
	Algorithm string
	Public    [ed25519.PublicKeySize]byte
	Private   *[ed25519.PrivateKeySize]byte
}

// HSMRSAKey represents the information for an HSMRSAKey with ObjectHandle for private portion
type HSMRSAKey struct {
	ID        string
	Algorithm string
	Public    []byte
	Private   pkcs11.ObjectHandle
}
