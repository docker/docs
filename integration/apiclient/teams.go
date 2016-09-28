package apiclient

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"path"
// )

// // These constants are used by various client methods.
// const (
// 	TeamTypeManaged = "managed"
// 	TeamTypeLDAP    = "ldap"
// )

// // Team represents a team within an organization.
// type Team struct {
// 	ID                       uint   `json:"id"`
// 	OrgID                    uint   `json:"orgID"`
// 	Type                     string `json:"type"`
// 	Name                     string `json:"name"`
// 	Description              string `json:"description"`
// 	LdapDN                   string `json:"ldapDN,omitempty"`
// 	LdapGroupMemberAttribute string `json:"ldapGroupMemberAttribute,omitempty"`
// }

// // createTeamForm containst the fields used in CreateTeam requests.
// type createTeamForm struct {
// 	Name                     string `json:"name"`
// 	Description              string `json:"description"`
// 	Type                     string `json:"type"`
// 	LdapDN                   string `json:"ldapDN"`
// 	LdapGroupMemberAttribute string `json:"ldapGroupMemberAttribute"`
// }

// // createTeam creates a team in the given organization with the data in the
// // given form.
// func (c *apiClient) createTeam(orgname string, form createTeamForm) (*Team, error) {
// 	response, err := c.makeRequest("POST", path.Join(accountsBasePath, orgname, "teams"), form)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusCreated); err != nil {
// 		return nil, err
// 	}

// 	var team Team
// 	if err := json.NewDecoder(response.Body).Decode(&team); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return &team, nil
// }

// // CreateManagedTeam creates a new managed team in an organization with the
// // given teamname and description.
// func (c *apiClient) CreateManagedTeam(orgname, teamname, description string) (*Team, error) {
// 	return c.createTeam(orgname, createTeamForm{
// 		Type:        TeamTypeManaged,
// 		Name:        teamname,
// 		Description: description,
// 	})
// }

// // CreateLDAPGroupSyncedTeam creates a new LDAP group-synced team in an
// // organization with the given teamname, description, groupDN, and memberAttr.
// func (c *apiClient) CreateLDAPGroupSyncedTeam(orgname, teamname, description, groupDN, memberAttr string) (*Team, error) {
// 	return c.createTeam(orgname, createTeamForm{
// 		Type:        TeamTypeLDAP,
// 		Name:        teamname,
// 		Description: description,
// 		LdapDN:      groupDN,
// 		LdapGroupMemberAttribute: memberAttr,
// 	})
// }

// // ListTeams lists all of the teams in an organization.
// func (c *apiClient) ListTeams(orgname string) ([]*Team, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "teams"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var teamsResponse struct {
// 		Teams []*Team `json:"teams"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&teamsResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return teamsResponse.Teams, nil
// }

// // GetTeam retrieves info about a team in an organization.
// func (c *apiClient) GetTeam(orgname, teamname string) (*Team, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "teams", teamname), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var team Team
// 	if err := json.NewDecoder(response.Body).Decode(&team); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return &team, nil
// }

// // TeamUpdateForm is used to update fields of a team.
// // If Name is not blank, the name will be updated.
// // If Description is not nil, the description will be updated.
// // If LdapDN is not blank and the team is synced with an LDAP group, the LDAP
// // DN of the group will be updated.
// // If LdapGroupMemberAttribute is not blank and the team is synced with an LDAP
// // group, the group member attribute of the team will be updated.
// type TeamUpdateForm struct {
// 	Name                     string  `json:"name"`
// 	Description              *string `json:"description"`
// 	LdapDN                   string  `json:"ldapDN"`
// 	LdapGroupMemberAttribute string  `json:"ldapGroupMemberAttribute"`
// }

// // UpdateTeam updates a team in an organization with the field in the given
// // form.
// func (c *apiClient) UpdateTeam(orgname, teamname string, form TeamUpdateForm) (*Team, error) {
// 	response, err := c.makeRequest("PATCH", path.Join(accountsBasePath, orgname, "teams", teamname), form)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var team Team
// 	if err := json.NewDecoder(response.Body).Decode(&team); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return &team, nil
// }

// // DeleteTeam deletes a team from an organization.
// // Note: this method may be called repeatedly with no effect.
// func (c *apiClient) DeleteTeam(orgname, teamname string) error {
// 	response, err := c.makeRequest("DELETE", path.Join(accountsBasePath, orgname, "teams", teamname), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusNoContent)
// }

// // ListTeamMembers lists the users which are a member of a team in an
// // organization.
// func (c *apiClient) ListTeamMembers(orgname, teamname string) ([]*Account, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "teams", teamname, "members"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var membersResponse struct {
// 		Members []*Account `json:"members"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&membersResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return membersResponse.Members, nil
// }

// // ListOrganizationMemberTeams lists the teams which an organization member
// // belongs to.
// func (c *apiClient) ListOrganizationMemberTeams(orgname, username string) ([]*Team, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "members", username, "teams"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var teamsResponse struct {
// 		Teams []*Team `json:"teams"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&teamsResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return teamsResponse.Teams, nil
// }

// // ListOrgMembers lists the users which are a member of a team in an
// // organization.
// func (c *apiClient) ListOrgMembers(orgname string) ([]*Account, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "members"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var membersResponse struct {
// 		Members []*Account `json:"members"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&membersResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return membersResponse.Members, nil
// }

// func (c *apiClient) ListUserOrganizations(username string) ([]*Account, error) {
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

// // CheckTeamMembership checks whether a user is a member of a team in an
// // organization.
// func (c *apiClient) CheckTeamMembership(orgname, teamname, username string) (isMember bool, err error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "teams", teamname, "members", username), nil)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	switch response.StatusCode {
// 	case http.StatusNotFound:
// 		// 404 Not Found is returned when the user is not a member
// 		// (or the org, team, or user dose not exist).
// 		return false, nil
// 	case http.StatusNoContent:
// 		// 204 No Content is returned when the user is a member.
// 		return true, nil
// 	}

// 	return false, validateStatusCode(response, 0)
// }

// // CheckOrganizationMembership checks whether a user is a member of an
// // organization.
// func (c apiClient) CheckOrganizationMembership(orgname, username string) (isMember bool, err error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, orgname, "members", username), nil)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	switch response.StatusCode {
// 	case http.StatusNotFound:
// 		// 404 Not Found is returned when the user is not a member
// 		// (or the org, team, or user dose not exist).
// 		return false, nil
// 	case http.StatusNoContent:
// 		// 204 No Content is returned when the user is a member.
// 		return true, nil
// 	}

// 	return false, validateStatusCode(response, 0)
// }

// // AddTeamMember adds a user to a team within an organization.
// // Note: The team must be a managed (i.e., non-ldap group-synced) team.
// // Note: this method may be called repeatedly with no effect.
// func (c *apiClient) AddTeamMember(orgname, teamname, username string) error {
// 	response, err := c.makeRequest("PUT", path.Join(accountsBasePath, orgname, "teams", teamname, "members", username), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusOK)
// }

// // DeleteTeamMember removes a user from a team within an organization.
// // Note: The team must be a managed (i.e., non-ldap group-synced) team.
// // Note: this method may be called repeatedly with no effect.
// func (c *apiClient) DeleteTeamMember(orgname, teamname, username string) error {
// 	response, err := c.makeRequest("DELETE", path.Join(accountsBasePath, orgname, "teams", teamname, "members", username), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusNoContent)
// }

// // DeleteOrganizationMember removes a user from an organization (i.e., all
// // teams in that organization).
// func (c *apiClient) DeleteOrganizationMember(orgname, username string) error {
// 	response, err := c.makeRequest("DELETE", path.Join(accountsBasePath, orgname, "members", username), nil)
// 	if err != nil {
// 		return fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	return validateStatusCode(response, http.StatusNoContent)
// }
