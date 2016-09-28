package managed

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
)

type authenticator struct {
	schemaMgr schema.Manager
}

var _ authn.UsernamePasswordAuthenticator = (*authenticator)(nil)

// New returns a new authenticator that authenticates all clients using basic
// auth with hashed passwords stored in the database.
func New(schemaMgr schema.Manager) authn.UsernamePasswordAuthenticator {
	return &authenticator{
		schemaMgr: schemaMgr,
	}
}

// AuthenticateUsernamePassword follows the same semantics as
// AuthenticateRequest in the authn.Authenticator interface but only attempts
// basic authentication with a username and password.
func (a *authenticator) AuthenticateUsernamePassword(ctx context.Context, username, password string) (*authn.Account, *errors.APIError) {
	user, err := a.schemaMgr.GetUserByName(username)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			return nil, authn.ErrInvalidUsernamePassword()
		}

		return nil, errors.Internal(ctx, fmt.Errorf("unable to get user by name: %s", err))
	}

	if !user.IsActive {
		return nil, authn.ErrAccountInactive()
	}

	if !passwords.CheckPassword(ctx, user.PasswordHash, password) {
		return nil, authn.ErrInvalidUsernamePassword()
	}

	return &authn.Account{
		Account: *user,
	}, nil
}
