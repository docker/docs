package signed

import (
	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/keys"
)

type SigningService interface {
	Sign(keyIDs []string, data []byte) ([]data.Signature, error)
}

type KeyService interface {
	Create() (*keys.PublicKey, error)
	PublicKeys(keyIDs ...string) (map[string]*keys.PublicKey, error)
}

type TrustService interface {
	SigningService
	KeyService
}
