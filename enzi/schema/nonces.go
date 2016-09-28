package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	// ErrNoSuchServiceAuthCode conveys that an authorization code does not
	// exist.
	ErrNoSuchServiceAuthCode = errors.New("no such service auth code")
)

// ServiceAuthCode represents a nonce used during service authorization by a
// user as part of an authorization or login flow. When a user authorizes a
// service, a random authorizaiton code is generated as a nonce and stored here
// with a short expiration time (1 minute) with reference to the service that
// was authorized, the ID of the account granting authorization, and the ID of
// the session of the user to facilitate service session linking.
type ServiceAuthCode struct {
	ID          string    `gorethink:"id"` // Hash of the code. Primary key.
	Code        string    `gorethink:"-"`  // Not stored in the DB.
	Expiration  time.Time `gorethink:"expiration"`
	ServiceID   string    `gorethink:"serviceID"`
	AccountID   string    `gorethink:"accountID"`
	SessionID   string    `gorethink:"sessionID"`
	RedirectURI string    `gorethink:"redirectURI"`
}

// makeServiceAuthCodeID makes a service auth code ID the given secret code.
func makeServiceAuthCodeID(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))

	return hex.EncodeToString(hasher.Sum(nil))
}

var serviceAuthCodesTable = table{
	db:         dbName,
	name:       "service_auth_codes",
	primaryKey: "id", // Guarantees uniqueness of auth code. Quick lookups.
	secondaryIndexes: map[string][]string{
		"expiration": nil, // For quickly deleting expired auth codes.
	},
}

// CreateServiceAuthCode creates a new service authcode using the values from
// the given serviceAuthCode. To prevent secrets from being stored in the
// backend, the service auth code ID is set to a hash of a random secret. This
// secret is not stored in the backend and can only be produced once by this
// function which sets it as the Code field on the given serviceAuthCode value.
func (m *manager) CreateServiceAuthCode(serviceAuthCode *ServiceAuthCode) error {
	serviceAuthCode.Code = uuid.NewV4().String()
	serviceAuthCode.ID = makeServiceAuthCodeID(serviceAuthCode.Code)

	if _, err := serviceAuthCodesTable.Term().Insert(serviceAuthCode).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to insert service auth code in db: %s", err)
	}

	return nil
}

// GetServiceAuthCode returns the service auth code with the given code. If no
// such service auth code exists the error will be ErrNoSuchServiceAuthCode.
func (m *manager) GetServiceAuthCode(code string) (*ServiceAuthCode, error) {
	cursor, err := serviceAuthCodesTable.Term().Get(makeServiceAuthCodeID(code)).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var serviceAuthCode ServiceAuthCode
	if err := cursor.One(&serviceAuthCode); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchServiceAuthCode
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &serviceAuthCode, nil
}

// DeleteServiceAuthCode deletes the service auth code with the given ID. This
// should be done to ensure that auth codes are not used more than once.
func (m *manager) DeleteServiceAuthCode(id string) error {
	if _, err := serviceAuthCodesTable.Term().Get(id).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete service auth code from db: %s", err)
	}

	return nil
}

// DeleteExpiredServiceAuthCodes removes any service auth codes with an
// expiration time which is in the past.
func (m *manager) DeleteExpiredServiceAuthCodes() error {
	if _, err := serviceAuthCodesTable.Term().Between(
		rethink.MinVal, time.Now(),
		rethink.BetweenOpts{Index: "expiration"},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete service auth codes from db: %s", err)
	}

	return nil
}
