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

// ErrNoSuchSession conveys that a session with the given ID does not exist.
var ErrNoSuchSession = errors.New("no such session")

// Session represents a user's login session on the server. Note that the
// Secret field is only available when the session is first created.
type Session struct {
	ID         string    `gorethink:"id"`
	Secret     string    `gorethink:"-"` // Not stored in the database.
	UserID     string    `gorethink:"userID"`
	CSRFToken  string    `gorethink:"csrfToken"`
	Expiration time.Time `gorethink:"expiration"`
}

var sessionsTable = table{
	db:         dbName,
	name:       "sessions",
	primaryKey: "id", // Guarantees uniqueness of session IDs. Quick lookups.
	secondaryIndexes: map[string][]string{
		"userID_id":  {"userID", "id"}, // For quickly listing sessions for a user.
		"expiration": nil,              // For quickly deleting expired sessions.
	},
}

// MakeSessionID makes a session ID or service session ID for the given secret.
func MakeSessionID(secret string) string {
	hasher := sha256.New()
	hasher.Write([]byte(secret))

	return hex.EncodeToString(hasher.Sum(nil))
}

// CreateSession inserts a new session into the *sessions* table of the
// database using the given user ID and duration. To prevent secrets from being
// stored in the backend, the session ID is set to a hash of a random secret
// UUID. This secret is not stored in the backend and can only be produced once
// by this function. The CSRF token is also set to a random UUID but is not
// kept secret.
func (m *manager) CreateSession(userID string, duration time.Duration) (*Session, error) {
	secret := uuid.NewV4().String()

	session := &Session{
		ID:         MakeSessionID(secret),
		Secret:     secret,
		UserID:     userID,
		CSRFToken:  uuid.NewV4().String(),
		Expiration: time.Now().Add(duration),
	}

	if _, err := sessionsTable.Term().Insert(session).RunWrite(m.session); err != nil {
		return nil, fmt.Errorf("unable to insert session in db: %s", err)
	}

	return session, nil
}

// GetSession retrieves the session corresponding to the given ID. If no
// such session exists, the error will be ErrNoSuchSession.
func (m *manager) GetSession(id string) (*Session, error) {
	cursor, err := sessionsTable.Term().Get(id).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var session Session
	if err := cursor.One(&session); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchSession
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &session, nil
}

// ExtendSession updates the expiration time of the given session to be now
// plus the given additional duration of time.
func (m *manager) ExtendSession(session *Session, duration time.Duration) error {
	session.Expiration = time.Now().Add(duration)

	if _, err := sessionsTable.Term().Get(session.ID).Update(
		map[string]interface{}{"expiration": session.Expiration},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to update session in db: %s", err)
	}

	return nil
}

// DeleteSession removes the session with the given ID.
func (m *manager) DeleteSession(id string) error {
	if err := m.DeleteServiceSessions(id); err != nil {
		return err
	}

	if _, err := sessionsTable.Term().Get(id).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete session from db: %s", err)
	}

	return nil
}

// DeleteSessionsForUser deletes all sessions for the given userID except for
// the session which matches the given excludeID.
func (m *manager) DeleteSessionsForUser(userID string, excludeID string) error {
	if _, err := sessionsTable.Term().Between(
		[]interface{}{userID, rethink.MinVal},
		[]interface{}{userID, rethink.MaxVal},
		rethink.BetweenOpts{Index: "userID_id"},
	).Filter(func(row rethink.Term) rethink.Term {
		return row.Field("id").Ne(excludeID)
	}).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete sessions for user: %s", err)
	}

	return nil
}

// DeleteExpiredSessions removes any sessions with an expiration time which is
// in the past.
func (m *manager) DeleteExpiredSessions() error {
	if _, err := sessionsTable.Term().Between(
		rethink.MinVal, time.Now(),
		rethink.BetweenOpts{Index: "expiration"},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete sessions from db: %s", err)
	}

	return nil
}

// ServiceSession represents a link to a user's login session on another
// service. Note that the Secret field is only available when the session link
// is first created.
type ServiceSession struct {
	PK        string `gorethink:"pk"`        // Hash of sessionID and serviceID. Primary key.
	ID        string `gorethink:"id"`        // Hash of randomly generated secret.
	Secret    string `gorethink:"-"`         // Not stored in the database. The service session has its own secret.
	SessionID string `gorethink:"sessionID"` // Foreign key reference to root session.
	ServiceID string `gorethink:"serviceID"` // Foreign key reference to service.
	CSRFToken string `gorethink:"csrfToken"` // The service session has its own CSRF token too.
}

var serviceSessionsTable = table{
	db:         dbName,
	name:       "service_sessions",
	primaryKey: "pk", // Guarantees uniqueness of (sessionID, serviceID). Quick lookups.
	secondaryIndexes: map[string][]string{
		"id": nil, // For quick lookups by ID.
		"sessionID_serviceID": {"sessionID", "serviceID"}, // For quickly listing service sessions associated with a root session.
	},
}

func computeServiceSessionPK(sessionID, serviceID string) string {
	hash := sha256.New()

	hash.Write([]byte(sessionID))
	hash.Write([]byte(serviceID))

	return hex.EncodeToString(hash.Sum(nil))
}

// CreateServiceSession creates a new service session linked to the given root
// sessionID for the service with the given serviceID. To prevent secrets from
// being stored in the backend, the service session ID is set to a hash of a
// random secret UUID. This secret is not stored in the backend and can only be
// produced once by this function. The CSRF token is also set to a random UUID
// but is not kept secret.
func (m *manager) CreateServiceSession(sessionID, serviceID string) (*ServiceSession, error) {
	secret := uuid.NewV4().String()

	session := &ServiceSession{
		PK:        computeServiceSessionPK(sessionID, serviceID),
		ID:        MakeSessionID(secret),
		Secret:    secret,
		SessionID: sessionID,
		ServiceID: serviceID,
		CSRFToken: uuid.NewV4().String(),
	}

	if _, err := serviceSessionsTable.Term().Insert(
		session, rethink.InsertOpts{Conflict: "replace"},
	).RunWrite(m.session); err != nil {
		return nil, fmt.Errorf("unable to insert session in db: %s", err)
	}

	return session, nil
}

// GetServiceSession retrieves the service session corresponding to the given
// ID. If no such session exists, the error will be ErrNoSuchSession.
func (m *manager) GetServiceSession(id string) (*ServiceSession, error) {
	cursor, err := serviceSessionsTable.Term().GetAllByIndex("id", id).Limit(1).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var session ServiceSession
	if err := cursor.One(&session); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchSession
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &session, nil
}

// DeleteServiceSessions removes the service sessions associated with the given
// root session ID.
func (m *manager) DeleteServiceSessions(sessionID string) error {
	if _, err := serviceSessionsTable.Term().Between(
		[]interface{}{sessionID, rethink.MinVal},
		[]interface{}{sessionID, rethink.MaxVal},
		rethink.BetweenOpts{Index: "sessionID_serviceID"},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete service sessions from db: %s", err)
	}

	return nil
}
