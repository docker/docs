package keys

import (
	"errors"

	"github.com/endophage/go-tuf/data"
)

var (
	ErrWrongType        = errors.New("tuf: invalid key type")
	ErrExists           = errors.New("tuf: key already in db")
	ErrWrongID          = errors.New("tuf: key id mismatch")
	ErrInvalidKey       = errors.New("tuf: invalid key")
	ErrInvalidRole      = errors.New("tuf: invalid role")
	ErrInvalidKeyID     = errors.New("tuf: invalid key id")
	ErrInvalidThreshold = errors.New("tuf: invalid role threshold")
)

type PublicKey struct {
	data.Key
	ID string
}

func NewPublicKey(keyType string, public []byte) *PublicKey {
	// create a copy so the private key is not included
	key := data.Key{
		Type:  keyType,
		Value: data.KeyValue{Public: public},
	}
	return &PublicKey{key, key.ID()}
}

type PrivateKey struct {
	PublicKey
	Private []byte
}

type DB struct {
	types map[string]int
	roles map[string]*data.Role
	keys  map[string]*PublicKey
}

func NewDB() *DB {
	return &DB{
		roles: make(map[string]*data.Role),
		keys:  make(map[string]*PublicKey),
	}
}

func (db *DB) AddKey(k *PublicKey) error {
	//if _, ok := db.types[k.Type]; !ok {
	//	return ErrWrongType
	//}
	//if len(k.Value.Public) != ed25519.PublicKeySize {
	//	return ErrInvalidKey
	//}

	key := PublicKey{
		Key: data.Key{
			Type: k.Type,
			Value: data.KeyValue{
				Public: make([]byte, len(k.Value.Public)),
			},
		},
		ID: k.ID,
	}

	copy(key.Value.Public, k.Value.Public)

	db.keys[k.ID] = &key
	return nil
}

var validRoles = map[string]struct{}{
	"root":      {},
	"targets":   {},
	"snapshot":  {},
	"timestamp": {},
}

func ValidRole(name string) bool {
	_, ok := validRoles[name]
	return ok
}

func (db *DB) AddRole(name string, r *data.Role) error {
	if !ValidRole(name) {
		return ErrInvalidRole
	}
	if r.Threshold < 1 {
		return ErrInvalidThreshold
	}

	// validate all key ids have the correct length
	for _, id := range r.KeyIDs {
		if len(id) != data.KeyIDLength {
			return ErrInvalidKeyID
		}
	}

	db.roles[name] = r
	return nil
}

func (db *DB) GetKey(id string) *PublicKey {
	return db.keys[id]
}

func (db *DB) GetRole(name string) *data.Role {
	return db.roles[name]
}
