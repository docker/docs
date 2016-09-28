package apiclient

import (
// "encoding/json"
// "net/http"

// "github.com/docker/dhe-deploy/adminserver"
)

// const ldapCheckSubroute = "/api/v0/admin/settings/ldapcheck"
// const authSubroute = "/api/v0/admin/settings/auth"
// const sqlusersSubroute = "/api/v0/admin/settings/sqlusers"

func (c *apiClient) Username() string {
	return c.loginSession.Account.Name
}

// func (c *apiClient) Password() string {
// 	return c.password
// }

// func (c *apiClient) LDAPCheck(params *adminserver.LDAPCheckSettings) (*map[string]interface{}, error) {
// 	if resp, err := c.makeRequest("POST", ldapCheckSubroute, params); err != nil {
// 		return nil, err
// 	} else {
// 		defer resp.Body.Close()

// 		if err := validateStatusCode(resp, http.StatusOK); err != nil {
// 			return nil, err
// 		}

// 		ldapCheck := new(map[string]interface{})
// 		if err := json.NewDecoder(resp.Body).Decode(ldapCheck); err != nil {
// 			return nil, err
// 		} else {
// 			return ldapCheck, nil
// 		}
// 	}
// }

// func (c *apiClient) GetAuthSettings() (*adminserver.AuthSettings, error) {
// 	if resp, err := c.GetAuthSettingsResponse(); err != nil {
// 		return nil, err
// 	} else {
// 		defer resp.Body.Close()

// 		if err := validateStatusCode(resp, http.StatusOK); err != nil {
// 			return nil, err
// 		}

// 		var authSettings adminserver.AuthSettings
// 		if err := json.NewDecoder(resp.Body).Decode(&authSettings); err != nil {
// 			return nil, err
// 		} else {
// 			return &authSettings, nil
// 		}
// 	}
// }

// func (c *apiClient) GetAuthSettingsResponse() (*http.Response, error) {
// 	return c.makeRequest("GET", authSubroute, nil)
// }

// func (c *apiClient) SetAuthSettings(settings *adminserver.AuthSettings) (*adminserver.AuthSettings, error) {
// 	resp, err := c.SetAuthSettingsResponse(settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	if err := validateStatusCode(resp, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var authSettings adminserver.AuthSettings
// 	if err := json.NewDecoder(resp.Body).Decode(&authSettings); err != nil {
// 		return nil, err
// 	}

// 	return &authSettings, nil
// }

// func (c *apiClient) SetAuthSettingsResponse(settings *adminserver.AuthSettings) (*http.Response, error) {
// 	return c.makeRequest("PUT", authSubroute, settings)
// }

// func (c *apiClient) ListSqlUsers() ([]*adminserver.FormUser, error) {
// 	if resp, err := c.ListSqlUsersResponse(); err != nil {
// 		return nil, err
// 	} else {
// 		defer resp.Body.Close()
// 		if err := validateStatusCode(resp, http.StatusOK); err != nil {
// 			return nil, err
// 		}

// 		var managedAuthSettings adminserver.ManagedAuthSettings
// 		if err := json.NewDecoder(resp.Body).Decode(&managedAuthSettings); err != nil {
// 			return nil, err
// 		}

// 		return managedAuthSettings.Users, nil
// 	}
// }

// func (c *apiClient) ListSqlUsersResponse() (*http.Response, error) {
// 	return c.makeRequest("GET", sqlusersSubroute, nil)
// }
