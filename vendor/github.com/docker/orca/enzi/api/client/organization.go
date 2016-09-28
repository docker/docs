package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// ListUserOrganizations lists organization accounts which the user with the
// given usernameOrID is a member of. If it is an account ID, it should be
// prefixed with `id:` to disambiguate it from an account name. This API call
// only returns organizations with a name greater than or equal to the given
// start name. The given limit value specifies the maximum number of
// organizations to return. If zero, the default is 10. The limit can be no
// greater than 2^16. If there is an API error response then the returned error
// will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListUserOrganizations(usernameOrID, start string, limit uint) (organizations *responses.MemberOrgs, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/organizations", usernameOrID), params)

	organizations = new(responses.MemberOrgs)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, organizations, &nextPageStart); err != nil {
		return nil, "", err
	}

	return organizations, nextPageStart, nil
}

// ListOrganizationMembers lists user accounts which are members of the
// organization with the given orgNameOrID. If it is an account ID, it should
// be prefixed with `id:` to disambiguate it from an account name. The given
// filter can be either 'admins', 'non-admins', or 'all' (default if left
// empty). This API call only returns members with an account ID greater than
// or equal to the given start ID. The given limit value specifies the maximum
// number of members to return. If zero, the default is 10. The limit can be no
// greater than 2^16. If there is an API error response then the returned error
// will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListOrganizationMembers(orgNameOrID, filter, start string, limit uint) (members *responses.Members, nextPageStart string, err error) {
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

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members", orgNameOrID), params)

	members = new(responses.Members)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, members, &nextPageStart); err != nil {
		return nil, "", err
	}

	return members, nextPageStart, nil
}

// ListPublicOrganizationMembers lists user accounts which are public members
// of the organization with the given orgNameOrID, i.e., those members who have
// opted to reveal their membership in the organization to non-members. If it
// is an account ID, it should be prefixed with `id:` to disambiguate it from
// an account name. This API call only returns members with an account ID
// greater than or equal to the given start ID. The given limit value specifies
// the maximum number of members to return. If zero, the default is 10. The
// limit can be no greater than 2^16. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListPublicOrganizationMembers(orgNameOrID, start string, limit uint) (members *responses.Members, nextPageStart string, err error) {
	params := url.Values{}
	if start != "" {
		params.Set("start", start)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/publicMembers", orgNameOrID), params)

	members = new(responses.Members)
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, members, &nextPageStart); err != nil {
		return nil, "", err
	}

	return members, nextPageStart, nil
}

// GetOrganizationAdminMemberSyncConfig retrieves the configuration options for
// syncing the admin members of an organization with an LDAP group or search
// filter. If orgNameOrID is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is not configured to use an LDAP backend for syncing
// accounts and membership then there will be an APIError with the code
// "LDAP_REQUIRED".
func (s *Session) GetOrganizationAdminMemberSyncConfig(orgNameOrID string) (*responses.MemberSyncOpts, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/adminMemberSyncConfig", orgNameOrID), nil)

	var config responses.MemberSyncOpts
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetOrganizationAdminMemberSyncConfig submits a form to set or replace the
// configuration options for syncing the admin members of an organization with
// an LDAP group or search filter. If orgNameOrID is an account ID, it should
// be prefixed with `id:` to disambiguate it from an account name. If there is
// an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is not configured to use an LDAP backend for syncing
// accounts and membership then there will be an APIError with the code
// "LDAP_REQUIRED".
func (s *Session) SetOrganizationAdminMemberSyncConfig(orgNameOrID string, form forms.MemberSyncOpts) (*responses.MemberSyncOpts, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/adminMemberSyncConfig", orgNameOrID), nil)

	var config responses.MemberSyncOpts
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &config, nil); err != nil {
		return nil, err
	}

	return &config, nil
}

// DeleteOrganizationMember submits a request to delete the member with the
// given memberNameOrID from the organization with the given orgNameOrID. If
// either is an account ID, it should be prefixed with `id:` to disambiguate
// it from an account name. If there is an API error response then the returned
// error will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) DeleteOrganizationMember(orgNameOrID, memberNameOrID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members/%s", orgNameOrID, memberNameOrID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetOrganizationMember submits a request to retrieve membership info for
// the user with the given memberNameOrID in the organization with the given
// orgNameOrID. If either is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the authenticated client does not have access to view the member's info
// then the returned APIError will have the code "NOT_AUTHORIZED". If the user
// is not a member of the organization then (nil, nil) is returned.
func (s *Session) GetOrganizationMember(orgNameOrID, memberNameOrID string) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members/%s", orgNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &member, nil); err != nil {
		if apiErrs, ok := err.(*errors.APIErrors); ok && apiErrs.HTTPStatusCode == http.StatusNotFound {
			// Indicates that the user is not a member (or that the
			// org doesn't exist or that the user doesn't exist).
			return nil, nil
		}

		return nil, err
	}

	return &member, nil
}

// UpdateOrganizationMember submits a form to update membership info for
// the user with the given memberNameOrID in the organization with the given
// orgNameOrID. If either is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) UpdateOrganizationMember(orgNameOrID, memberNameOrID string, form forms.SetMembership) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members/%s", orgNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("PATCH", endpoint, form, http.StatusOK, &member, nil); err != nil {
		return nil, err
	}

	return &member, nil
}

// AddOrganizationMember submits a form to set or replace membership info for
// the user with the given memberNameOrID in the organization with the given
// orgNameOrID. If either is an account ID, it should be prefixed with `id:` to
// disambiguate it from an account name. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the auth service is configured to use an LDAP backend for syncing
// accounts and membership and the organization is configured to sync admin
// members with an LDAP group or search filter then then there will be an
// APIError with the code "LDAP_PRECLUDES".
func (s *Session) AddOrganizationMember(orgNameOrID, memberNameOrID string, form forms.SetMembership) (*responses.Member, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/accounts/%s/members/%s", orgNameOrID, memberNameOrID), nil)

	var member responses.Member
	if err := s.performRequest("PUT", endpoint, form, http.StatusOK, &member, nil); err != nil {
		return nil, err
	}

	return &member, nil
}
