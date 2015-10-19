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

// HSMKey represents the information for an HSMKey with ObjectHandle for private portion
type HSMKey struct {
	id        string
	algorithm data.KeyAlgorithm
	public    []byte
	private   pkcs11.ObjectHandle
}

// NewHSMKey returns a HSMKey
func NewHSMKey(public []byte, private pkcs11.ObjectHandle) *HSMKey {
	return &HSMKey{
		public:  public,
		private: private,
	}
}

// Algorithm implements a method of the data.Key interface
func (k *HSMKey) Algorithm() data.KeyAlgorithm {
	return k.algorithm
}

// ID implements a method of the data.Key interface
func (k *HSMKey) ID() string {
	if k.id == "" {
		pubK := data.NewPublicKey(k.algorithm, k.public)
		k.id = pubK.ID()
	}
	return k.id
}

// Public implements a method of the data.Key interface
func (k *HSMKey) Public() []byte {
	return k.public
}

// Private implements a method of the data.PrivateKey interface
func (k *HSMKey) Private() []byte {
	// Not possible to return private key bytes from a hardware device
	return nil
}

// PKCS11ObjectHandle returns the PKCS11 object handle stored in the HSMKey
// structure
func (k *HSMKey) PKCS11ObjectHandle() pkcs11.ObjectHandle {
	return k.private
}
