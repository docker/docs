package config

import (
	"fmt"

	"github.com/docker/orca/enzi/schema"
)

// AuthConfigPropertyKey is the property key which maps to the system's auth
// configuration.
const AuthConfigPropertyKey = "config.auth"

// Auth is a config object which specifies the auth backend to use.
type Auth struct {
	Backend string `json:"backend"`
}

// Supported Auth Backends.
const (
	AuthBackendManaged = "managed"
	AuthBackendLDAP    = "ldap"
)

var (
	defaultAuthConfig = Auth{
		Backend: AuthBackendManaged,
	}

	// SupportedAuthBackends is the set of supported auth backends.
	SupportedAuthBackends = map[string]struct{}{
		AuthBackendManaged: {},
		AuthBackendLDAP:    {},
	}
)

// GetAuthConfig retrieves the current auth configuration using the given
// schema manager.
func GetAuthConfig(mgr schema.Manager) (authConfig *Auth, err error) {
	authConfig = new(Auth)

	if err := mgr.GetProperty(AuthConfigPropertyKey, authConfig); err != nil {
		if err == schema.ErrNoSuchProperty {
			// Use the default backend.
			return &defaultAuthConfig, nil
		}

		return nil, fmt.Errorf("unable to lookup auth config property: %s", err)
	}

	return authConfig, nil
}

// SetAuthConfig sets the current auth configuration to the given auth config
// value. If it is nil, the current auth config will be deleted.
func SetAuthConfig(mgr schema.Manager, authConfig *Auth) error {
	if authConfig == nil {
		return mgr.DeleteProperty(AuthConfigPropertyKey)
	}

	if err := mgr.SetProperty(AuthConfigPropertyKey, authConfig); err != nil {
		return fmt.Errorf("unable to set auth config property: %s", err)
	}

	return nil
}
