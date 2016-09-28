package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// ListOrganizationTeams lists teams in the organization with the given
// orgNameOrID. If it is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. This API call only returns teams with
// a name greater than or equal to the given start name. The given limit value
// specifies the maximum number of teams to return. If zero, the default is 10.
// The limit can be no greater than 2^16. If there is an API error response
// then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListOrganizationTeams(orgNameOrID, start string, limit uint) (teams *responses.Teams, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams", orgNameOrID), params)

	teams = new(responses.Teams)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, teams, &nextPageStart); err != nil {
		return nil, "", err
	}

	return teams, nextPageStart, nil
}

// CreateTeam submits a form to create a team in the organization with the
// given orgNameOrID. If it is an account ID, it should be prefixed with `id:`
// to disambiguate it from an account name. If there is an API error response
// then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) CreateTeam(orgNameOrID string, form forms.CreateTeam) (*responses.Team, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams", orgNameOrID), nil)

	var team responses.Team
	if err := s.performRequest("POST", endpoint, form, http.StatusCreated, &team, nil); err != nil {
		return nil, err
	}

	return &team, nil
}

// DeleteTeam submits a request to delete the team with the given teamNameOrID
// from the organization with the given orgNameOrID. If etiher is an account
// ID, it should be prefixed with `id:` to disambiguate it from an account
// name. If there is an API error response then the returned error will be of
// the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) DeleteTeam(orgNameOrID, teamNameOrID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s", orgNameOrID, teamNameOrID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetTeam retrieves the team with the given teamNameOrID in the organization
// with the given orgNameOrID. If etiher is an account ID, it should be
// prefixed with `id:` to disambiguate it from an account name. If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetTeam(orgNameOrID, teamNameOrID string) (*responses.Team, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s", orgNameOrID, teamNameOrID), nil)

	var team responses.Team
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &team, nil); err != nil {
		return nil, err
	}

	return &team, nil
}

// UpdateTeam submits a form to update the team with the given teamNameOrID
// in the organization with the given orgNameOrID. If etiher is an account ID,
// it should be prefixed with `id:` to disambiguate it from an account name. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) UpdateTeam(orgNameOrID, teamNameOrID string, form forms.UpdateTeam) (*responses.Team, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s", orgNameOrID, teamNameOrID), nil)

	var team responses.Team
	if err := s.performRequest("PATCH", endpoint, form, http.StatusOK, &team, nil); err != nil {
		return nil, err
	}

	return &team, nil
}

// ListOrganizationMemberTeams lists teams in the organization with the given
// orgNameOrID of which the user with the given memberNameOrID is member. If
// either is an account ID, it should be prefixed with `id:` to disambiguate it
// from an account name. This API call only returns teams with an ID greater
// than or equal to the given start ID. The given limit value specifies the
// maximum number of teams to return. If zero, the default is 10. The limit can
// be no greater than 2^16. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListOrganizationMemberTeams(orgNameOrID, memberNameOrID, start string, limit uint) (teams *responses.MemberTeams, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members/%s/teams", orgNameOrID, memberNameOrID), params)

	teams = new(responses.MemberTeams)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, teams, &nextPageStart); err != nil {
		return nil, "", err
	}

	return teams, nextPageStart, nil
}

// ListTeamMembers lists user accounts which are members of the team with the
// given teamNameOrID in the organization with the given orgNameOrID. If either
// is an account ID, it should be prefixed with `id:` to disambiguate it from
// an account name. The given filter can be either 'admins', 'non-admins', or
// 'all' (default if left empty). This API call only returns members with an
// account ID greater than or equal to the given start ID. The given limit
// value specifies the maximum number of members to return. If zero, the
// default is 10. The limit can be no greater than 2^16. If there is an API
// error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListTeamMembers(orgNameOrID, teamNameOrID, filter, start string, limit uint) (members *responses.Members, nextPageStart string, err error) {
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

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/members", orgNameOrID, teamNameOrID), params)

	members = new(responses.Members)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, members, &nextPageStart); err != nil {
		return nil, "", err
	}

	return members, nextPageStart, nil
}

// ListPublicTeamMembers lists user accounts which are public members of the
// team with the given teamNameOrID in the organization with the given
// orgNameOrID, i.e., those members who have opted to reveal their membership
// in the team to non-members. If either is an account ID, it should be
// prefixed with `id:` to disambiguate it from an account name. This API call
// only returns members with an account ID greater than or equal to the given
// start ID. The given limit value specifies the maximum number of members to
// return. If zero, the default is 10. The limit can be no greater than 2^16.
// If there is an API error response then the returned error will be of the
// type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListPublicTeamMembers(orgNameOrID, teamNameOrID, start string, limit uint) (members *responses.Members, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/publicMembers", orgNameOrID, teamNameOrID), params)

	members = new(responses.Members)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, members, &nextPageStart); err != nil {
		return nil, "", err
	}

	return members, nextPageStart, nil
}

// GetTeamMemberSyncConfig retrieves the configuration options for syncing the
// members of the team with the given teamNameOrID in the organization with the
// given orgNameOrID with an LDAP group or search filter. If either is an
// account ID, it should be prefixed with `id:` to disambiguate it from an
// account name. If there is an API error response then the returned error will
// be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is not configured to use an LDAP backend for syncing
// accounts and membership then there will be an APIError with the code
// "LDAP_REQUIRED".
func (s *Session) GetTeamMemberSyncConfig(orgNameOrID, teamNameOrID string) (*responses.MemberSyncOpts, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/memberSyncConfig", orgNameOrID, teamNameOrID), nil)

	var config responses.MemberSyncOpts
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetTeamMemberSyncConfig submits a form to set or replace the configuration
// options for syncing the members of the team with the given teamNameOrID in
// the organization with the given orgNameOrID with an LDAP group or search
// filter. If either is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is not configured to use an LDAP backend for syncing
// accounts and membership then there will be an APIError with the code
// "LDAP_REQUIRED".
func (s *Session) SetTeamMemberSyncConfig(orgNameOrID, teamNameOrID string, form forms.MemberSyncOpts) (*responses.MemberSyncOpts, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/memberSyncConfig", orgNameOrID, teamNameOrID), nil)

	var config responses.MemberSyncOpts
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// DeleteTeamMember submits a request to delete the member with the
// given memberNameOrID from the team with the given teamNameOrID in the
// organization with the given orgNameOrID. If either is an account ID, it
// should be prefixed with `id:` to disambiguate it from an account name. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is configured to use an LDAP backend for syncing
// accounts and membership and the team is configured to sync members with an
// LDAP group or search filter then then there will be an APIError with the
// code "LDAP_PRECLUDES".
func (s *Session) DeleteTeamMember(orgNameOrID, teamNameOrID, memberNameOrID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/members/%s", orgNameOrID, teamNameOrID, memberNameOrID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetTeamMember submits a request to retrieve membership info for the user
// with the given memberNameOrID in the team with the given teamNameOrID in the
// organization with the given orgNameOrID. If either is an account ID, it
// should be prefixed with `id:` to disambiguate it from an account name. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the authenticated client does not have access to view the member's info
// then the returned APIError will have the code "NOT_AUTHORIZED". If the user
// is not a member of the team then (nil, nil) is returned.
func (s *Session) GetTeamMember(orgNameOrID, teamNameOrID, memberNameOrID string) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/members/%s", orgNameOrID, teamNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &member, nil); err != nil {
		if apiErrs, ok := err.(*errors.APIErrors); ok && apiErrs.HTTPStatusCode == http.StatusNotFound {
			// Indicates that the user is not a member (or that the
			// org, team, or user doesn't exist).
			return nil, nil
		}

		return nil, err
	}

	return &member, nil
}

// UpdateTeamMember submits a form to update membership info for the user with
// the given memberNameOrID in the team with the given teamNameOrID in the
// organization with the given orgNameOrID. If either is an account ID, it
// should be prefixed with `id:` to disambiguate it from an account name. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) UpdateTeamMember(orgNameOrID, teamNameOrID, memberNameOrID string, form forms.SetMembership) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/members/%s", orgNameOrID, teamNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("PATCH", endpoint, form, http.StatusOK, &member, nil); err != nil {
		return nil, err
	}

	return &member, nil
}

// AddTeamMember submits a form to set or replace membership info for
// the user with the given memberNameOrID in the team with the given
// teamNameOrID in the organization with the given orgNameOrID. If either is an
// account ID, it should be prefixed with `id:` to disambiguate it from an
// account name. If there is an API error response then the returned error will
// be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is configured to use an LDAP backend for syncing
// accounts and membership and the team is configured to sync members with an
// LDAP group or search filter then then there will be an APIError with the
// code "LDAP_PRECLUDES".
func (s *Session) AddTeamMember(orgNameOrID, teamNameOrID, memberNameOrID string, form forms.SetMembership) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/teams/%s/members/%s", orgNameOrID, teamNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &member, nil); err != nil {
		return nil, err
	}

	return &member, nil
}
