package authz

import (
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authn/ldap"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/authn/managed"
	"github.com/docker/orca/enzi/authn/openidtoken"
	"github.com/docker/orca/enzi/authn/sessiontoken"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
)

// Authorizer defines an interface with all of the methods we need for
// being a garant authorizer backend, authenticationg request users in DTR,
// and authorizing access to repositories in DTR.
type Authorizer interface {
	authn.RequestAuthenticator

	SessionTokenAuthenticator() sessiontoken.Authenticator
	OpenIDTokenAuthenticator() authn.OpenIDTokenAuthenticator
	UsernamePasswordAuthenticator(ldapSettings *ldapconfig.Settings) authn.UsernamePasswordAuthenticator

	AuthConfig() (*config.Auth, error)
	LDAPSettings() (*ldapconfig.Settings, error)
	AccountAccess(acct *schema.Account, account *authn.Account) (accountAccess *MembershipAccess, err error)
	OrgMembershipAccess(orgID string, account *authn.Account) (*MembershipAccess, error)
	TeamMembershipAccess(teamID string, orgAccess MembershipAccess, account *authn.Account) (teamAccess *MembershipAccess, err error)
}

// Authorizer is capable of authenticating a user and determining access
// levels for that user.
type authorizer struct {
	schemaMgr schema.Manager
}

// Assert that *authorizer implements the Authorizer interface.
var _ Authorizer = (*authorizer)(nil)

// NewAuthorizer creates a new authorizer using the given database connection.
func NewAuthorizer(schemaMgr schema.Manager) Authorizer {
	return &authorizer{
		schemaMgr: schemaMgr,
	}
}

func (a *authorizer) AuthConfig() (*config.Auth, error) {
	return config.GetAuthConfig(a.schemaMgr)
}

func (a *authorizer) LDAPSettings() (*ldapconfig.Settings, error) {
	return ldapconfig.GetLDAPConfig(a.schemaMgr)
}

// UsernamePasswordAuthenticator returns a UsernamePasswordAuthenticator. If
// the given ldapSettings value is not nil, an LDAP authenticator is returned,
// otherwise a managed authenticator is returned.
func (a *authorizer) UsernamePasswordAuthenticator(ldapSettings *ldapconfig.Settings) authn.UsernamePasswordAuthenticator {
	if ldapSettings != nil {
		return ldap.New(ldapSettings, a.schemaMgr)
	}

	return managed.New(a.schemaMgr)
}

func (a *authorizer) SessionTokenAuthenticator() sessiontoken.Authenticator {
	return sessiontoken.New(a.schemaMgr)
}

func (a *authorizer) OpenIDTokenAuthenticator() authn.OpenIDTokenAuthenticator {
	return openidtoken.New(a.schemaMgr)
}
