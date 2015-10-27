package keys

import (
	"errors"

	"github.com/docker/notary/tuf/data"
	"github.com/miekg/pkcs11"
)

var (
	// ErrExists happens when a Key already exists in a database
	ErrExists = errors.New("notary-signer: key already in db")
	// ErrInvalidKeyID error happens when a key isn't found
	ErrInvalidKeyID = errors.New("notary-signer: invalid key id")
	// ErrFailedKeyGeneration happens when there is a failure in generating a key
	ErrFailedKeyGeneration = errors.New("notary-signer: failed to generate new key")
)

// HSMRSAKey represents the information for an HSMRSAKey with ObjectHandle for private portion
type HSMRSAKey struct {
	id      string
	public  []byte
	private pkcs11.ObjectHandle
}

// NewHSMRSAKey returns a HSMRSAKey
func NewHSMRSAKey(public []byte, private pkcs11.ObjectHandle) *HSMRSAKey {
	return &HSMRSAKey{
		public:  public,
		private: private,
	}
}

// Algorithm implements a method of the data.Key interface
func (rsa *HSMRSAKey) Algorithm() data.KeyAlgorithm {
	return data.RSAKey
}

// ID implements a method of the data.Key interface
func (rsa *HSMRSAKey) ID() string {
	if rsa.id == "" {
		pubK := data.NewPublicKey(rsa.Algorithm(), rsa.Public())
		rsa.id = pubK.ID()
	}
	return rsa.id
}

// Public implements a method of the data.Key interface
func (rsa *HSMRSAKey) Public() []byte {
	return rsa.public
}

// Private implements a method of the data.PrivateKey interface
func (rsa *HSMRSAKey) Private() []byte {
	// Not possible to return private key bytes from a hardware device
	return nil
}

// PKCS11ObjectHandle returns the PKCS11 object handle stored in the HSMRSAKey
// structure
func (rsa *HSMRSAKey) PKCS11ObjectHandle() pkcs11.ObjectHandle {
	return rsa.private
}
