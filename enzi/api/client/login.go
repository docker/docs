package client

import (
	"net/http"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// Login submits a form to Login as a user account and create a session. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) Login(form forms.Login) (*responses.LoginSession, error) {
	endpoint := s.buildURL("/v0/id/login", nil)

	var session responses.LoginSession
	if err := s.performRequest("POST", endpoint, form, http.StatusOK, &session, nil); err != nil {
		return nil, err
	}

	return &session, nil
}

// Logout submits a form to Logout, terminating the current session (if the
// current RequestAuthenticator is a SessionTokenAuthenticator). If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) Logout() error {
	endpoint := s.buildURL("/v0/id/logout", nil)

	return s.performRequest("POST", endpoint, nil, http.StatusNoContent, nil, nil)
}
