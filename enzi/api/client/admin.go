package client

import (
	"net/http"

	"github.com/docker/orca/enzi/api/forms"
)

// ImportAccounts submits the given form to import an account. If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/enzi/api/errors).APIErrors
func (s *Session) ImportAccounts(form forms.ImportAccounts) error {
	endpoint := s.buildURL("/v0/admin/importAccounts", nil)

	return s.performRequest("POST", endpoint, form, http.StatusCreated, nil, nil)
}
