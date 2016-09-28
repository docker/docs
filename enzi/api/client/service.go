package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// ListAccountOwnedServices lists services owned by the account with the given
// accountNameOrID. If it is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. This API call only returns services
// with a name greater than or equal to the given start name. The given limit
// value specifies the maximum number of services to return. If zero, the
// default is 10. The limit can be no greater than 2^16. If there is an API
// error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListAccountOwnedServices(accountNameOrID, start string, limit uint) (services *responses.Services, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/services", accountNameOrID), params)

	services = new(responses.Services)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, services, &nextPageStart); err != nil {
		return nil, "", err
	}

	return services, nextPageStart, nil
}

// CreateService submits a form to create a service to be owned by the account
// with the given accountNameOrID. If it is an account ID, it should be
// prefixed with `id:` to disambiguate it from an account name. If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) CreateService(accountNameOrID string, form forms.CreateService) (*responses.Service, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/services", accountNameOrID), nil)

	var service responses.Service
	if err := s.performRequest("POST", endpoint, form, http.StatusCreated, &service, nil); err != nil {
		return nil, err
	}

	return &service, nil
}

// DeleteService submits a request to delete the service with the given
// serviceNameOrID owned by the account with the given accountNameOrID. If
// etiher is an account ID, it should be prefixed with `id:` to disambiguate it
// from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) DeleteService(accountNameOrID, serviceNameOrID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/services/%s", accountNameOrID, serviceNameOrID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetService retrieves the service with the given serviceNameOrID owned by the
// account with the given accountNameOrID. If etiher is an account ID, it
// should be prefixed with `id:` to disambiguate it from an account name. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetService(accountNameOrID, serviceNameOrID string) (*responses.Service, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/services/%s", accountNameOrID, serviceNameOrID), nil)

	var service responses.Service
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &service, nil); err != nil {
		return nil, err
	}

	return &service, nil
}

// UpdateService submits a form to update the service with the given
// serviceNameOrID owned by the account with the given accountNameOrID. If
// etiher is an account ID, it should be prefixed with `id:` to disambiguate it
// from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) UpdateService(accountNameOrID, serviceNameOrID string, form forms.UpdateService) (*responses.Service, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/services/%s", accountNameOrID, serviceNameOrID), nil)

	var service responses.Service
	if err := s.performRequest("PATCH", endpoint, form, http.StatusOK, &service, nil); err != nil {
		return nil, err
	}

	return &service, nil
}
