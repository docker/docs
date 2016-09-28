package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// // These constants are used by various client methods.
const (
	// 	AccountTypeUser         = "user"
	// 	AccountTypeOrganization = "organization"
	// 	GlobalOrgName           = "_global"

	accountsBasePath = "/v0/accounts"
)

// Account contains fields from account API responses.
type Account struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	LdapLogin string `json:"ldapLogin,omitempty"`
	IsActive  bool   `json:"isActive"`
}

// // createAccountForm contains fields used when creating an account.
// // Note: some fields may be unused depending on the type of account being
// // created.
// type createAccountForm struct {
// 	Type      string `json:"type"`
// 	Name      string `json:"name"`
// 	Password  string `json:"password"`
// 	LdapLogin string `json:"ldapLogin"`
// }

// // createAccount creates an account using the given form.
// func (c *apiClient) createAccount(form createAccountForm) (*Account, error) {
// 	response, err := c.makeRequest("POST", accountsBasePath, form)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusCreated); err != nil {
// 		return nil, err
// 	}

// 	var acct Account
// 	if err := json.NewDecoder(response.Body).Decode(&acct); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return &acct, nil
// }

// // CreateUser creates a managed user account with the given username and
// // password.
// func (c *apiClient) CreateUser(username, password string) (*Account, error) {
// 	return c.createAccount(createAccountForm{
// 		Type:     AccountTypeUser,
// 		Name:     username,
// 		Password: password,
// 	})
// }

// // CreateUserWithLDAPLogin creates a user account with the given ldapLogin
// // (which may be different) and password.
// func (c *apiClient) CreateUserWithLDAPLogin(username, ldapLogin, password string) (*Account, error) {
// 	return c.createAccount(createAccountForm{
// 		Type:      AccountTypeUser,
// 		Name:      username,
// 		Password:  password,
// 		LdapLogin: ldapLogin,
// 	})
// }

// // CreateOrganization creates an organization account with the given name.
// // Note: the client must authenticate as a superuser to create an organization.
// func (c *apiClient) CreateOrganization(name string) (*Account, error) {
// 	return c.createAccount(createAccountForm{
// 		Type: AccountTypeOrganization,
// 		Name: name,
// 	})
// }

// func (c *apiClient) listAccounts(query string) ([]*Account, error) {
// 	response, err := c.makeRequest("GET", accountsBasePath+"?"+query, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var accountsResponse struct {
// 		Accounts []*Account `json:"accounts"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&accountsResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return accountsResponse.Accounts, nil
// }

// // ListAccounts lists all accounts on the server.
// func (c *apiClient) ListAccounts() ([]*Account, error) {
// 	return c.listAccounts("")
// }

// // ListUserAccounts lists all user accounts on the server.
// func (c *apiClient) ListUserAccounts() ([]*Account, error) {
// 	return c.listAccounts(url.Values{"type": {"user"}}.Encode())
// }

// // ListOrgAccounts lists all user accounts on the server.
// func (c *apiClient) ListOrgAccounts() ([]*Account, error) {
// 	return c.listAccounts(url.Values{"type": {"organization"}}.Encode())
// }

// GetAccount retrieves the account with the given name.
func (c *apiClient) GetAccount(name string) (*Account, error) {
	u := path.Join(accountsBasePath, name)
	if c.ucpAsEnzi {
		u = "/enzi" + u
	}

	resp, err := c.makeRequest("GET", url.URL{Path: u}, nil)
	if err != nil {
		return nil, fmt.Errorf("error executing request to account endpoint %s", err)
	}
	defer resp.Body.Close()

	if err := validateStatusCode(resp, http.StatusOK); err != nil {
		return nil, err
	}

	acct := Account{}
	if err := json.NewDecoder(resp.Body).Decode(&acct); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &acct, nil
}

// // DeleteAccount deletes the account with the given name.
// // Note: this method may be called repeatedly with no effect.
// func (c *apiClient) DeleteAccount(name string) error {
// 	response, err := c.makeRequest("DELETE", path.Join(accountsBasePath, name), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusNoContent)
// }

// // ChangePassword changes the password of a user account.
// // Note: If authenticated as a superuser, oldPassword may be left blank.
// func (c *apiClient) ChangePassword(username, oldPassword, newPassword string) error {
// 	changePasswordForm := map[string]string{
// 		"oldPassword": oldPassword,
// 		"newPassword": newPassword,
// 	}

// 	response, err := c.makeRequest("POST", path.Join(accountsBasePath, username, "changePassword"), changePasswordForm)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusOK)
// }

// // ActivateUser activates the user account with the given name.
// // Note: the client must be authenticated as a superuser.
// func (c *apiClient) ActivateUser(username string) error {
// 	response, err := c.makeRequest("PUT", path.Join(accountsBasePath, username, "activate"), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusOK)
// }

// // DeactivateUser deactivates the user account with the given name.
// // Note: the client must be authenticated as a superuser.
// func (c *apiClient) DeactivateUser(username string) error {
// 	response, err := c.makeRequest("PUT", path.Join(accountsBasePath, username, "deactivate"), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusOK)
// }

// // ListOrganizations lists all of the organizations that the given user is a
// // member of.
// func (c *apiClient) ListOrganizations(username string) ([]*Account, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, username, "organizations"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var organizationsResponse struct {
// 		Organizations []*Account `json:"organizations"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&organizationsResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return organizationsResponse.Organizations, nil
// }
