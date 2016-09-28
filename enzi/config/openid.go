package config

import (
	"fmt"

	"github.com/docker/orca/enzi/schema"
)

// OpenIDConfigPropertyKey is the property key which maps to the system's
// OpenID Connect Provider configuration.
const OpenIDConfigPropertyKey = "config.openid"

// OpenID is a config object which specifies parameters for the system's
// OpenID Connect Provider configuration.
type OpenID struct {
	IssuerIdentifier string `json:"issuerIdentifier"`
}

var defaultOpenIDConfig = OpenID{
	IssuerIdentifier: "https://enzi.example.com",
}

// GetOpenIDConfig retrieves the current openID configuration using the
// given schema manager.
func GetOpenIDConfig(mgr schema.Manager) (openIDConfig *OpenID, err error) {
	openIDConfig = new(OpenID)

	if err := mgr.GetProperty(OpenIDConfigPropertyKey, openIDConfig); err != nil {
		if err == schema.ErrNoSuchProperty {
			// Use the default value.
			return &defaultOpenIDConfig, nil
		}

		return nil, fmt.Errorf("unable to lookup openID config property: %s", err)
	}

	return openIDConfig, nil
}

// SetOpenIDConfig sets the current openID configuration to the given openID
// config value. If it is nil, the current openID config will be deleted.
func SetOpenIDConfig(mgr schema.Manager, openIDConfig *OpenID) error {
	if openIDConfig == nil {
		return mgr.DeleteProperty(OpenIDConfigPropertyKey)
	}

	if err := mgr.SetProperty(OpenIDConfigPropertyKey, openIDConfig); err != nil {
		return fmt.Errorf("unable to set openID config property: %s", err)
	}

	return nil
}
