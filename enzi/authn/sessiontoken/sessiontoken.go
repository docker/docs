package sessiontoken

import (
	"fmt"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/schema"
)

const (
	defaultSessionLifetime           time.Duration = time.Hour * 72
	defaultSessionExtensionThreshold time.Duration = time.Hour * 24
)

var errInvalidSession = errors.InvalidAuthentication("invalid session")

// Authenticator extends the default interface with additional methods for
// managing sessions.
type Authenticator interface {
	authn.SessionTokenAuthenticator

	// CreateSession creates a new session for the given user.
	CreateSession(user *authn.Account) (*schema.Session, error)
	// ExtendSession extends the given session according to the default
	// session extension policy.
	ExtendSession(session *schema.Session) error
	// DeleteSession deletes the given session.
	DeleteSession(session *schema.Session) error
	// DeleteSessionsForUser deletes all sessions for the user with the
	// given userID except for the given session (if not nil).
	DeleteSessionsForUser(userID string, session *schema.Session) error
}

type authenticator struct {
	schemaMgr schema.Manager
}

var _ Authenticator = (*authenticator)(nil)

// New creates an authenticator using the given schema manager.
func New(schemaMgr schema.Manager) Authenticator {
	return &authenticator{
		schemaMgr: schemaMgr,
	}
}

// AuthenticateSessionToken follows the same semantics as AuthenticateRequest
// in the authn.Authenticator interface but only attempts session token
// authentication. If the session is invalid the returned error should be
// authn.ErrInvalidSession.
func (a *authenticator) AuthenticateSessionToken(ctx context.Context, token string) (*authn.Account, *errors.APIError) {
	// Token value should correspond to the session secret.
	session, err := a.schemaMgr.GetSession(schema.MakeSessionID(token))
	if err != nil {
		if err == schema.ErrNoSuchSession {
			return nil, errInvalidSession
		}

		// Internal error.
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get session: %s", err))
	}

	if session.Expiration.Before(time.Now()) {
		return nil, errInvalidSession
	}

	user, err := a.schemaMgr.GetUserByID(session.UserID)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			// The user no longer exists and the session hasn't yet
			// been cleaned up. Handle like an invalid session.
			return nil, errInvalidSession
		}

		// Internal error.
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get user for session: %s", err))
	}

	if !user.IsActive {
		return nil, authn.ErrAccountInactive()
	}

	// Attempt to extend the session according to session policy.
	if err := a.ExtendSession(session); err != nil {
		return nil, errors.Internal(ctx, fmt.Errorf("unable to extend session: %s", err))
	}

	return &authn.Account{
		Account: *user,
		Session: session,
	}, nil
}

// CreateSession creates a new session for the given user.
func (a *authenticator) CreateSession(user *authn.Account) (*schema.Session, error) {
	session, err := a.schemaMgr.CreateSession(user.ID, defaultSessionLifetime)
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %s", err)
	}

	return session, nil
}

// ExtendSession extends the expiration of the given session to 72 hours from
// now if the session expires in less than 24 hours from now.
func (a *authenticator) ExtendSession(session *schema.Session) error {
	if session.Expiration.After(time.Now().Add(defaultSessionExtensionThreshold)) {
		// The sesison expires more than the designated threshold of
		// time from now, so there is no need to extend the session.
		return nil
	}

	// Extend the session by the default duration.
	return a.schemaMgr.ExtendSession(session, defaultSessionLifetime)
}

// DeleteSession deletes the given session, making it unusable.
func (a *authenticator) DeleteSession(session *schema.Session) error {
	if err := a.schemaMgr.DeleteSession(session.ID); err != nil {
		return fmt.Errorf("unable to delete session: %s", err)
	}

	return nil
}

// DeleteSessionsForUser deletes all sessions for the user with the given
// userID except for the given session (if not nil).
func (a *authenticator) DeleteSessionsForUser(userID string, session *schema.Session) error {
	var excludeID string
	if session != nil {
		excludeID = session.ID
	}

	if err := a.schemaMgr.DeleteSessionsForUser(userID, excludeID); err != nil {
		return fmt.Errorf("unable to delete sessions for user: %s", err)
	}

	return nil
}
