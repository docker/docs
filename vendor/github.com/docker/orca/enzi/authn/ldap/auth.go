package ldap

import (
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/schema"
)

type authenticator struct {
	*config.Settings
	schemaMgr schema.Manager
}

var _ authn.UsernamePasswordAuthenticator = (*authenticator)(nil)

// New returns a new authenticator that authenticates all clients using basic
// auth. It searches for a user with the given LDAP settings and attempts to
// BIND as that user.
func New(settings *config.Settings, schemaMgr schema.Manager) authn.UsernamePasswordAuthenticator {
	return &authenticator{
		Settings:  settings,
		schemaMgr: schemaMgr,
	}
}
