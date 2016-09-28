package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/docker/orca/enzi/jose"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	// ErrNoSuchSigningKey conveys that a signing key with the given id
	// does not exist.
	ErrNoSuchSigningKey = errors.New("no such signing key")
	// ErrNoSuchServiceKey conveys that a service key with the given id
	// does not exist.
	ErrNoSuchServiceKey = errors.New("no such service key")
)

// JWK is a JSON Web Key with fields for either RSA or ECDSA public keys.
type JWK struct {
	// Key ID value is arbitrary but should uniquely identify the public
	// key, usually via a hash.
	ID string `gorethink:"id"`
	// Key Type is either "RSA" or "EC"
	KeyType string `gorethink:"keyType"`
	// Fields for RSA public keys.
	Modulus  string `gorethink:"modulus,omitempty"`
	Exponent string `gorethink:"exponent,omitempty"`
	// Fields for ECDSA public keys. Curve is one of "P-256", "P-384", and
	// "P-521".
	Curve       string `gorethink:"curve,omitempty"`
	XCoordinate string `gorethink:"xCoordinate,omitempty"`
	YCoordinate string `gorethink:"yCoordinate,omitempty"`

	Expiration time.Time `gorethink:"expiration"`
}

// NewJWK creates a JWK from the given public key and expiration time that is
// the given diration of time from now.
func NewJWK(key *jose.PublicKey, duration time.Duration) JWK {
	return JWK{
		ID:          key.ID,
		KeyType:     key.KeyType,
		Modulus:     key.Modulus,
		Exponent:    key.Exponent,
		Curve:       key.Curve,
		XCoordinate: key.XCoordinate,
		YCoordinate: key.YCoordinate,
		Expiration:  time.Now().Add(duration),
	}
}

// PublicKey attempts to convert this JWK into a public key which can be used
// to verify JWS signatures.
func (key *JWK) PublicKey() (*jose.PublicKey, error) {
	return jose.NewPublicKeyJWK(jose.JWK{
		ID:          key.ID,
		KeyType:     key.KeyType,
		Modulus:     key.Modulus,
		Exponent:    key.Exponent,
		Curve:       key.Curve,
		XCoordinate: key.XCoordinate,
		YCoordinate: key.YCoordinate,
	})
}

var signingKeysTable = table{
	db:         dbName,
	name:       "signing_keys",
	primaryKey: "id", // Guarantees uniqueness of Key ID. Quick lookups.
	secondaryIndexes: map[string][]string{
		"expiration": nil, // For quickly deleting expired keys.
	},
}

// SaveSigningKey stores the given signing key.
func (m *manager) SaveSigningKey(jwk JWK) error {
	if _, err := signingKeysTable.Term().Insert(jwk).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to save signing key: %s", err)
	}

	return nil
}

// GetSigningKey returns the signing key with the given ID. If no such key
// exists, the error will be ErrNoSuchSigningKey.
func (m *manager) GetSigningKey(id string) (*JWK, error) {
	cursor, err := signingKeysTable.Term().Get(id).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var key JWK
	if err := cursor.One(&key); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchSigningKey
		}

		return nil, fmt.Errorf("unable to scan query result: %s", err)
	}

	return &key, nil
}

// ListSigningKeys returns a list of all currently used keys.
func (m *manager) ListSigningKeys() (keys []JWK, err error) {
	cursor, err := signingKeysTable.Term().Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	if err := cursor.All(&keys); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return keys, nil
}

// ExtendSigningKeyExpiration updates the expiration time of the signing key
// with the given id to be now plus the given additional duration of time.
func (m *manager) ExtendSigningKeyExpiration(id string, duration time.Duration) error {
	if _, err := signingKeysTable.Term().Get(id).Update(
		map[string]interface{}{"expiration": time.Now().Add(duration)},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to update record in db: %s", err)
	}

	return nil
}

// DeleteExpiredSigningKeys removes any signing keys with an expiration time
// which is in the past.
func (m *manager) DeleteExpiredSigningKeys() error {
	if _, err := signingKeysTable.Term().Between(
		rethink.MinVal, time.Now(),
		rethink.BetweenOpts{Index: "expiration"},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete signing keys from db: %s", err)
	}

	return nil
}

// ServiceKey is a public key which is used to verify tokens signed by a
// service.
type ServiceKey struct {
	PK        string `gorethink:"pk"`        // Hash of serviceID and key ID. Primary key.
	ServiceID string `gorethink:"serviceID"` // Foreign key relation to the service which uses this key.
	JWK       JWK    `gorethink:"jwk"`
}

func computeServiceKeyPK(serviceID, keyID string) string {
	hash := sha256.New()

	hash.Write([]byte(serviceID))
	hash.Write([]byte(keyID))

	return hex.EncodeToString(hash.Sum(nil))
}

var serviceKeysTable = table{
	db:         dbName,
	name:       "service_keys",
	primaryKey: "pk", // Guarantees uniqueness of (serviceID, Key ID). Quick lookups.
	secondaryIndexes: map[string][]string{
		"jwk.expiration": nil, // For quickly deleting expired keys.
	},
}

// SaveServiceKeys stores the given keys for the given service. Keys are
// updated if they already exist, extending key expiration in the process.
func (m *manager) SaveServiceKeys(serviceID string, keys ...JWK) error {
	serviceKeys := make([]ServiceKey, len(keys))
	for i, key := range keys {
		serviceKeys[i] = ServiceKey{
			PK:        computeServiceKeyPK(serviceID, key.ID),
			ServiceID: serviceID,
			JWK:       key,
		}
	}

	if _, err := serviceKeysTable.Term().Insert(
		serviceKeys,
		rethink.InsertOpts{Conflict: "replace"},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to save service keys: %s", err)
	}

	return nil
}

// GetServiceKey retrieves the service key with the given serviceID and keyID.
// If no such service key exists the error will be ErrNoSuchServiceKey.
func (m *manager) GetServiceKey(serviceID, keyID string) (*ServiceKey, error) {
	cursor, err := serviceKeysTable.Term().Get(computeServiceKeyPK(serviceID, keyID)).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var serviceKey ServiceKey
	if err := cursor.One(&serviceKey); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchServiceKey
		}

		return nil, fmt.Errorf("unable to scan query result: %s", err)
	}

	return &serviceKey, nil
}

// DeleteExpiredServiceKeys removes any service keys with an expiration time
// which is in the past.
func (m *manager) DeleteExpiredServiceKeys() error {
	if _, err := serviceKeysTable.Term().Between(
		rethink.MinVal, time.Now(),
		rethink.BetweenOpts{Index: "jwk.expiration"},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete service keys from db: %s", err)
	}

	return nil
}
