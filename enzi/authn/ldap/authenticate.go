package ldap

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
	"github.com/go-ldap/ldap"
)

// AuthenticateUsernamePassword follows the same semantics as
// AuthenticateRequest in the authn.Authenticator interface but only attempts
// basic authentication with a username and password. The LDAP authenticator
// validates that the given username or password are not blank, gets the
// corresponding user account, and performs a bind against the configured LDAP
// server.
func (a *authenticator) AuthenticateUsernamePassword(ctx context.Context, username, password string) (*authn.Account, *errors.APIError) {
	if password == "" {
		// A BIND with an empty password may succeed as an Anonymous
		// BIND. We DO NOT want that to happen.
		return nil, authn.ErrInvalidUsernamePassword()
	}

	// Lookup the user in the database.
	user, err := a.schemaMgr.GetUserByName(username)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			return nil, authn.ErrInvalidUsernamePassword()
		}

		return nil, errors.Internal(ctx, fmt.Errorf("unable to get user by name: %s", err))
	}

	// We can try password hash authentication for the special recovery
	// admin user account.
	if user.Name == a.RecoveryAdminUsername && passwords.CheckPassword(ctx, user.PasswordHash, password) {
		// Password matches. Explicitly set as admin, active.
		user.IsAdmin = true
		user.IsActive = true

		return &authn.Account{
			Account: *user,
		}, nil
	}

	if user.LdapDN == "" {
		context.GetLogger(ctx).Warnf("user %s has an empty LDAP DN", user.Name)
		return nil, authn.ErrInvalidUsernamePassword()
	}

	if !user.IsActive {
		return nil, authn.ErrAccountInactive()
	}

	// Attempt to connect to the LDAP server and test the DN/password
	// combination.
	ldapConn, err := GetConn(a.ServerURL, a.Settings)
	if err != nil {
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get LDAP connection: %s", err))
	}
	defer ldapConn.Close()

	if err := ldapConn.Bind(user.LdapDN, password); err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials {
			return nil, authn.ErrInvalidUsernamePassword()
		}

		return nil, errors.Internal(ctx, fmt.Errorf("unable to perform LDAP bind request: %s", err))
	}

	return &authn.Account{
		Account: *user,
	}, nil
}
