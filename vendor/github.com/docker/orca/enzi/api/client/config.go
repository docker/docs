package client

import (
	"net/http"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// GetAuthConfig retrieves the current System Auth Configuration. If there is
// an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetAuthConfig() (*responses.AuthConfig, error) {
	endpoint := s.buildURL("/v0/config/auth", nil)

	var config responses.AuthConfig
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetAuthConfig submits a form to set or replace the current System Auth
// Configuration. If there is an API error response then the returned error
// will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) SetAuthConfig(form forms.AuthConfig) (*responses.AuthConfig, error) {
	endpoint := s.buildURL("/v0/config/auth", nil)

	var config responses.AuthConfig
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetLDAPSettings retrieves the current System LDAP Configuration. If there is
// an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetLDAPSettings() (*responses.LDAPSettings, error) {
	endpoint := s.buildURL("/v0/config/auth/ldap", nil)

	var settings responses.LDAPSettings
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &settings, nil); err != nil {
		return nil, err
	}

	return &settings, nil
}

// SetLDAPSettings submits a form to set or replace the current System LDAP
// Configuration. If there is an API error response then the returned error
// will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) SetLDAPSettings(form forms.LDAPSettings) (*responses.LDAPSettings, error) {
	endpoint := s.buildURL("/v0/config/auth/ldap", nil)

	var settings responses.LDAPSettings
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &settings, nil); err != nil {
		return nil, err
	}

	return &settings, nil
}

// TryLDAPLogin submits a form try a Login using the given configuration. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) TryLDAPLogin(form forms.TryLdapLogin) (*responses.Account, error) {
	endpoint := s.buildURL("/v0/config/auth/ldap/tryLogin", nil)

	var account responses.Account
	if err := s.performRequest("POST", endpoint, form, http.StatusOK, &account, nil); err != nil {
		return nil, err
	}

	return &account, nil
}
