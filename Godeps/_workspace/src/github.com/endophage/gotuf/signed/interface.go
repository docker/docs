package signed

import (
	"github.com/endophage/gotuf/data"
)

// SigningService defines the necessary functions to determine
// if a user is able to sign with a key, and to perform signing.
type SigningService interface {
	// Sign takes a slice of keyIDs and a piece of data to sign
	// and returns a slice of signatures and an error
	Sign(keyIDs []string, data []byte) ([]data.Signature, error)

	// CanSign takes a single keyID and returns a boolean indicating
	// whether the caller is able to sign with the keyID (i.e. does
	// this signing service hold the private key associated with
	// they keyID)
	CanSign(keyID string) bool
}

// KeyService provides management of keys locally. It will never
// accept or provide private keys. Communication between the KeyService
// and a SigningService happen behind the Create function.
type KeyService interface {
	// Create issues a new key pair and is responsible for loading
	// the private key into the appropriate signing service.
	Create() (*data.PublicKey, error)

	// PublicKeys return the PublicKey instances for the given KeyIDs
	PublicKeys(keyIDs ...string) (map[string]*data.PublicKey, error)
}

// CryptoService defines a unified Signing and Key Service as this
// will be most useful for most applications.
type CryptoService interface {
	SigningService
	KeyService
}

// Verifier defines an interface for verfying signatures. An implementer
// of this interface should verify signatures for one and only one
// signing scheme.
type Verifier interface {
	Verify(key data.Key, sig []byte, msg []byte) error
}
