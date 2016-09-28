package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// ListAccounts lists accounts from the eNZi API server using the given filter.
// Filter can be either 'users', 'orgs', 'admins', 'non-admins', or 'all'
// (default if left empty). This API call only returns accounts with a name
// greater than or equal to the given start name. The given limit value
// specifies the maximum number of accounts to return. If zero, the default is
// 10. The limit can be no greater than 2^16. If there is an API error response
// then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListAccounts(filter, start string, limit uint) (accounts *responses.Accounts, nextPageStart string, err error) {
	params := url.Values{}
	if filter != "" {
		params.Set("filter", filter)
	}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL("/v0/accounts", params)

	accounts = new(responses.Accounts)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, accounts, &nextPageStart); err != nil {
		return nil, "", err
	}

	return accounts, nextPageStart, nil
}

// CreateAccount submits the given form to create an account. If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) CreateAccount(form forms.CreateAccount) (*responses.Account, error) {
	endpoint := s.buildURL("/v0/accounts", nil)

	var account responses.Account
	if err := s.performRequest("POST", endpoint, form, http.StatusCreated, &account, nil); err != nil {
		return nil, err
	}

	return &account, nil
}

// DeleteAccount submits a request to delete the account with the given name or
// ID. If it is an account ID, it should be prefixed with `id:` to disambiguate
// it from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) DeleteAccount(accountNameOrID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s", accountNameOrID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetAccount submits a request to retrieve the account with the given name or
// ID. If it is an account ID, it should be prefixed with `id:` to disambiguate
// it from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetAccount(accountNameOrID string) (*responses.Account, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s", accountNameOrID), nil)

	var account responses.Account
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &account, nil); err != nil {
		return nil, err
	}

	return &account, nil
}

// UpdateAccount submits a form to update the account with the given name or
// ID. If it is an account ID, it should be prefixed with `id:` to disambiguate
// it from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) UpdateAccount(accountNameOrID string, form forms.UpdateAccount) (*responses.Account, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s", accountNameOrID), nil)

	var account responses.Account
	if err := s.performRequest("PATCH", endpoint, form, http.StatusOK, &account, nil); err != nil {
		return nil, err
	}

	return &account, nil
}

// ChangePassword submits a form to change the password of the user with the
// given name or ID. If it is an account ID, it should be prefixed with `id:`
// to disambiguate it from an account name. If there is an API error response
// then the returned error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ChangePassword(usernameOrID string, form forms.ChangePassword) (*responses.Account, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/changePassword", usernameOrID), nil)

	var account responses.Account
	if err := s.performRequest("POST", endpoint, form, http.StatusOK, &account, nil); err != nil {
		return nil, err
	}

	return &account, nil
}
