package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/filters"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// RouteCreateAccounts returns a route describing the CreateAccounts endpoint.
func (s *Service) routeCreateAccounts() server.Route {
	return server.Route{
		Method:     "POST",
		Path:       "/",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleCreateAccount),
		Doc:        "Create a user or organization",
		BodySample: forms.CreateAccount{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusCreated,
				Message: "Success, account created.",
				Sample:  responses.Account{},
			},
		},
	}
}

// HandleCreateAccount handles a request for creating a user or organizaiton.
func (s *Service) handleCreateAccount(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	// Gather request data.
	authConfig, errResponse := helpers.AuthConfig(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.CreateAccount)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	newAccount := &schema.Account{
		Name:     form.Name,
		IsOrg:    form.IsOrg,
		FullName: form.FullName,
	}

	// Creating an Organization or User?
	if !form.IsOrg {
		// Can't create users if using LDAP, they must be synced.
		if authConfig.Backend == config.AuthBackendLDAP {
			return responses.APIError(errors.CannotCreateUser("Users are synced with LDAP"))
		}

		passwordHash, err := passwords.HashPassword(form.Password)
		if err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}

		newAccount.IsActive = form.IsActive
		newAccount.IsAdmin = form.IsAdmin
		newAccount.PasswordHash = passwordHash
	}

	if err := s.schemaMgr.CreateAccount(newAccount); err != nil {
		if err == schema.ErrAccountExists {
			return responses.APIError(errors.AccountExists())
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusCreated, responses.MakeAccount(newAccount))
}

// RouteListAccounts returns a route describing the ListAccounts endpoint.
func (s *Service) routeListAccounts() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListAccounts),
		Doc:     "List users and organizations",
		Notes:   "Lists accounts in ascending order by name.",
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("filter", "Filter accounts by type or attribute - either 'users', 'orgs', 'admins', 'non-admins', or 'all' (default).").DefaultValue("all"),
			restful.QueryParameter("start", "Only return accounts with a name greater than or equal to this name."),
			restful.QueryParameter("limit", "Maximum number of accounts per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of accounts listed.",
				Sample:  responses.Accounts{},
			},
		},
	}
}

// HandleListAccounts handles a request for listing users and/or organizations.
func (s *Service) handleListAccounts(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	var (
		accounts         []schema.Account
		nextPageStart    string
		err              error
		startName, limit = helpers.PageParams(r, "start", "limit")
	)

	switch r.QueryParameter("filter") {
	case "users":
		accounts, nextPageStart, err = s.schemaMgr.ListUsers(startName, limit)
	case "orgs":
		accounts, nextPageStart, err = s.schemaMgr.ListOrgs(startName, limit)
	case "admins":
		accounts, nextPageStart, err = s.schemaMgr.ListAdmins(startName, limit)
	case "non-admins":
		accounts, nextPageStart, err = s.schemaMgr.ListNonAdmins(startName, limit)
	default:
		accounts, nextPageStart, err = s.schemaMgr.ListAccounts(startName, limit)
	}

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeAccounts(accounts), r, nextPageStart)
}

// RouteGetAccount returns a route describing the GetAccount endpoint.
func (s *Service) routeGetAccount() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{accountNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetAccount),
		Doc:     "Details for a user or organization",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account to fetch",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, account returned.",
				Sample:  responses.Account{},
			},
		},
	}
}

// HandleGetAccount handles a request for details of a user or organization.
func (s *Service) handleGetAccount(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]

	acct, errResponse := helpers.AccountOrNotFound(ctx, s.schemaMgr, accountNameOrID)
	if errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeAccount(acct))
}

// RouteUpdateAccount returns a route describing the UpdateAccount endpoint.
func (s *Service) routeUpdateAccount() server.Route {
	return server.Route{
		Method:  "PATCH",
		Path:    "/{accountNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleUpdateAccount),
		Doc:     "Update details for a user or organization",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account to update",
		},
		BodySample: forms.UpdateAccount{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, account updated.",
				Sample:  responses.Account{},
			},
		},
	}
}

// HandleUpdateAccount handles a request to update details of a user or
// organization.
func (s *Service) handleUpdateAccount(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.UpdateAccount)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	// Check authorization and sanitize.
	if rd.Acct.IsOrg {
		rd.Org = rd.Acct

		// The client needs org admin access to update the org.
		if errResponse := rd.AddFilters(
			rd.GetOrgAccess,
			rd.RequireOrgAdmin,
		).EvaluateFilters(); errResponse != nil {
			return errResponse
		}

		// Ensure these fields aren't set.
		form.IsActive = nil
		form.IsAdmin = nil
	} else {
		// The client needs to be an admin or the account being updated.
		if !(clientAccount.IsAdmin || rd.Acct.ID == clientAccount.ID) {
			return responses.APIError(errors.NotAuthorized("must be a system admin or the user being updated"))
		}

		changingActiveStatus := form.IsActive != nil && rd.Acct.IsActive != *form.IsActive
		changingAdminStatus := form.IsAdmin != nil && rd.Acct.IsAdmin != *form.IsAdmin

		// If we're changing the user's active or admin status, we
		// require a system admin.
		if (changingActiveStatus || changingAdminStatus) && !clientAccount.IsAdmin {
			return responses.APIError(errors.NotAuthorized("must be a system admin to change a user's admin or active status"))
		}

		usingLDAP := rd.AuthConfig.Backend == config.AuthBackendLDAP

		if changingActiveStatus && usingLDAP {
			return responses.APIError(errors.LdapPrecludes("user active status is synced with LDAP"))
		}

		if changingAdminStatus && usingLDAP {
			// With the LDAP backend, admin status can only be
			// changed manually if admins aren't configured to be
			// synced with LDAP.
			ldapSettings, errResponse := helpers.LDAPSettings(ctx, s.authorizer)
			if errResponse != nil {
				return errResponse
			}

			if ldapSettings.AdminSyncOpts.EnableSync {
				return responses.APIError(errors.LdapPrecludes("user admin status is synced with LDAP"))
			}
		}

		if form.IsActive != nil {
			rd.Acct.IsActive = *form.IsActive
		}
		if form.IsAdmin != nil {
			rd.Acct.IsAdmin = *form.IsAdmin
		}
	}

	if form.FullName != nil {
		rd.Acct.FullName = *form.FullName
	}

	updateFields := schema.AccountUpdateFields{
		FullName: form.FullName,
		IsAdmin:  form.IsAdmin,
		IsActive: form.IsActive,
	}

	if err := s.schemaMgr.UpdateAccount(rd.Acct.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeAccount(rd.Acct))
}

// RouteChangePassword returns a route describing the ChangePassword endpoint.
func (s *Service) routeChangePassword() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/{userNameOrID}/changePassword",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleChangePassword),
		Doc:     "Change a user's password",
		PathParameterDocs: map[string]string{
			"userNameOrID": "Name or id of username whose password is to be changed",
		},
		BodySample: forms.ChangePassword{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, password changed.",
				Sample:  responses.Account{},
			},
		},
	}
}

// HandleChangePassword handles a request to change a user's password.
func (s *Service) handleChangePassword(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	userNameOrID := pathParams["userNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetUser(userNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// Passwords are meaningless if we're using LDAP syncing.
	if rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.CannotChangePassword("you cannot use passwords when using LDAP authentication"))
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.ChangePassword)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	// The authenticated client needs to either be an admin or be the user
	// whose password is being changed.
	if !(clientAccount.IsAdmin || clientAccount.ID == rd.User.ID) {
		return responses.APIError(errors.NotAuthorized("must be a system admin or the user whose password is being changed"))
	}

	// If the client is changing their own password, we need to verify the
	// old password. Even if they are an admin, they must provide their
	// current password to change their own password.
	if clientAccount.ID == rd.User.ID && !passwords.CheckPassword(ctx, clientAccount.PasswordHash, form.OldPassword) {
		return responses.APIError(errors.PasswordIncorrect())
	}

	newPasswordHash, err := passwords.HashPassword(form.NewPassword)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	updateFields := schema.AccountUpdateFields{
		PasswordHash: &newPasswordHash,
	}

	if err := s.schemaMgr.UpdateAccount(rd.User.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	// Changing the password should invalidate all of the user's sessions
	// except for the session which was used to authenticate this request
	// (if the client authenticated via a session).
	if err := s.authorizer.SessionTokenAuthenticator().DeleteSessionsForUser(rd.User.ID, clientAccount.Session); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeAccount(rd.User))
}

// RouteDeleteAccount returns a route describing the DeleteAccount endpoint.
func (s *Service) routeDeleteAccount() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{accountNameOrID}",
		Handler: server.WrapHandlerWithAdminAccount(s.authorizer, s.handleDeleteAccount),
		Doc:     "Delete a user or organization",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account to delete",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, account deleted.",
			},
		},
	}
}

// HandleDeleteAccount handles a request for deleting a user or organization.
func (s *Service) handleDeleteAccount(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]

	account, err := helpers.AccountByNameOrID(s.schemaMgr, accountNameOrID)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			// This is OK. The account is already deleted.
			return responses.JSONResponse(http.StatusNoContent, nil)
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	if err := s.schemaMgr.DeleteAccount(account.ID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}

// RouteListUserOrganizations returns a route describing the
// ListUserOrganizations endpoint.
func (s *Service) routeListUserOrganizations() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{userNameOrID}/organizations",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListUserOrganizations),
		Doc:     "List a user's organization memberships",
		Notes:   "Lists organization memberships in ascending order by organization ID.",
		PathParameterDocs: map[string]string{
			"userNameOrID": "Name or id of user to whose organizations will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return memberships with an org ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of organizations per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of user's organizations listed.",
				Sample:  responses.MemberOrgs{},
			},
		},
	}
}

// HandleListUserOrganizations handles a request for listing a user's organization memberships.
func (s *Service) handleListUserOrganizations(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	userNameOrID := pathParams["userNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetUser(userNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// Only an admin or the user in question may list organization
	// memberships.
	if !(clientAccount.IsAdmin || clientAccount.ID == rd.User.ID) {
		return responses.APIError(errors.NotAuthorized("must be a system admin or the user whose organizations are being listed"))
	}

	startID, limit := helpers.PageParams(r, "start", "limit")

	memberOrgs, nextPageStart, err := s.schemaMgr.ListOrgsForUser(rd.User.ID, startID, limit)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMemberOrgs(memberOrgs), r, nextPageStart)
}

// RouteListOrganizationMembers returns a route describing the
// ListOrganizationMembers endpoint.
func (s *Service) routeListOrganizationMembers() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/members",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListOrganizationMembers),
		Doc:     "List members of an organization",
		Notes:   "Lists memberships in ascending order by user ID.",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization whose members will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("filter", "Filter members by type - either 'admins', 'non-admins', or 'all' (default).").DefaultValue("all"),
			restful.QueryParameter("start", "Only return members with a user ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of members per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of organization members listed.",
				Sample:  responses.Members{},
			},
		},
	}
}

// HandleListOrganizationMembers handles a request for listing memberships in
// an organization.
func (s *Service) handleListOrganizationMembers(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	var (
		orgMembers     []schema.MemberInfo
		nextPageStart  string
		err            error
		startID, limit = helpers.PageParams(r, "start", "limit")
	)

	switch r.QueryParameter("filter") {
	case "admins":
		orgMembers, nextPageStart, err = s.schemaMgr.ListOrgAdmins(rd.Org.ID, startID, limit)
	case "non-admins":
		orgMembers, nextPageStart, err = s.schemaMgr.ListOrgNonAdmins(rd.Org.ID, startID, limit)
	default:
		orgMembers, nextPageStart, err = s.schemaMgr.ListOrgMembers(rd.Org.ID, startID, limit)
	}

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMembers(orgMembers), r, nextPageStart)
}

// RouteListOrganizationPublicMembers returns a route describing the
// ListOrganizationPublicMembers endpoint.
func (s *Service) routeListOrganizationPublicMembers() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/publicMembers",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListOrganizationPublicMembers),
		Doc:     "List public members of an organization",
		Notes:   "Lists public members in ascending order by user ID.",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization whose public members will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return members with a user ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of members per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of organization public members listed.",
				Sample:  responses.Members{},
			},
		},
	}
}

// HandleListOrganizationPublicMembers handles a request for listing public
// members in an organization.
func (s *Service) handleListOrganizationPublicMembers(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	startID, limit := helpers.PageParams(r, "start", "limit")

	publicMembers, nextPageStart, err := s.schemaMgr.ListPublicOrgMembers(rd.Org.ID, startID, limit)

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMembers(publicMembers), r, nextPageStart)
}

// routeGetOrganizationAdminSyncConfig returns a route describing the
// GetOrganizationAdminSyncConfig endpoint.
func (s *Service) routeGetOrganizationAdminSyncConfig() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/adminMemberSyncConfig",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetOrganizationAdminSyncConfig),
		Doc:     "Get options for syncing admin members of an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization whose LDAP sync options to be retrieved",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, LDAP sync options retrieved.",
				Sample:  responses.MemberSyncOpts{},
			},
		},
	}
}

// handleGetOrganizationAdminSyncConfig handles a request to get options for
// syncing admin members of an organization.
func (s *Service) handleGetOrganizationAdminSyncConfig(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMemberSyncOpts(rd.Org.AdminSyncConfig))
}

// routeSetOrganizationAdminSyncConfig returns a route describing the
// SetOrganizationAdminSyncConfig endpoint.
func (s *Service) routeSetOrganizationAdminSyncConfig() server.Route {
	return server.Route{
		Method:  "PUT",
		Path:    "/{orgNameOrID}/adminMemberSyncConfig",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleSetOrganizationAdminSyncConfig),
		Doc:     "Set options for syncing admin members of an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization whose LDAP sync options to set",
		},
		BodySample: forms.MemberSyncOpts{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, LDAP sync options set.",
				Sample:  responses.MemberSyncOpts{},
			},
		},
	}
}

// handleSetOrganizationAdminSyncConfig handles a request to set options for
// syncing admin members of an organization.
func (s *Service) handleSetOrganizationAdminSyncConfig(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.MemberSyncOpts)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	rd.Org.AdminSyncConfig = schema.MemberSyncOpts{
		EnableSync:         form.EnableSync,
		GroupDN:            form.GroupDN,
		GroupMemberAttr:    form.GroupMemberAttr,
		SearchBaseDN:       form.SearchBaseDN,
		SearchScopeSubtree: form.SearchScopeSubtree,
		SearchFilter:       form.SearchFilter,
	}

	updateFields := schema.AccountUpdateFields{
		AdminSyncConfig: &rd.Org.AdminSyncConfig,
	}

	if err := s.schemaMgr.UpdateAccount(rd.Org.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMemberSyncOpts(rd.Org.AdminSyncConfig))
}

// RouteAddOrganizationMember returns a route describing the
// AddOrganizationMember endpoint.
func (s *Service) routeAddOrganizationMember() server.Route {
	return server.Route{
		Method:  "PUT",
		Path:    "/{orgNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleAddOrganizationMember),
		Doc:     "Add a member to an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the membership will be added",
			"memberNameOrID": "Name or id of user which will be added as a member",
		},
		BodySample: forms.SetMembership{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, membership set.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleAddOrganizationMember handles a request for adding a member to an
// organization.
func (s *Service) handleAddOrganizationMember(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
		rd.MakeFilterGetUser(memberNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// If the org admins should be synced with LDAP, clients can't manually
	// add members to the org.
	if rd.Org.AdminSyncConfig.EnableSync && rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.LdapPrecludes("this organizations's membership can only be changed via LDAP syncing and team membership"))
	}

	defer r.Request.Body.Close()

	var form forms.SetMembership
	if err := json.NewDecoder(r.Request.Body).Decode(&form); err != nil {
		return responses.APIError(errors.InvalidJSON(err))
	}

	if err := s.schemaMgr.AddOrgMembership(rd.Org.ID, rd.User.ID, form.IsAdmin, form.IsPublic); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	orgMembership, err := s.schemaMgr.GetOrgMembership(rd.Org.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	memberInfo := schema.MemberInfo{
		Member:   *rd.User,
		IsAdmin:  orgMembership.IsAdmin,
		IsPublic: orgMembership.IsPublic,
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMember(&memberInfo))
}

// RouteUpdateOrganizationMembership returns a route describing the
// UpdateOrganizationMember endpoint.
func (s *Service) routeUpdateOrganizationMembership() server.Route {
	return server.Route{
		Method:  "PATCH",
		Path:    "/{orgNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleUpdateOrganizationMembership),
		Doc:     "Update details of a user's membership in an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the membership will be updated",
			"memberNameOrID": "Name or id of user whose membership will be updated",
		},
		BodySample: forms.SetMembership{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, membership updated.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleUpdateOrganizationMembership handles a request for updating attributes
// of an org membership.
func (s *Service) handleUpdateOrganizationMembership(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.MakeFilterGetUser(memberNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// The client needs to either be an org admin or be the user whose
	// membership is being updated.
	isOrgAdmin := rd.ClientOrgAccess != nil && rd.ClientOrgAccess.IsAdmin
	if !(isOrgAdmin || clientAccount.ID == rd.User.ID) {
		return responses.APIError(errors.NotAuthorized("must have org admin access or authenticate as the user whose membership is being updated"))
	}

	if member, err := s.schemaMgr.GetOrgMembership(rd.Org.ID, rd.User.ID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	} else if member == nil {
		return responses.APIError(errors.NoSuchMember(memberNameOrID))
	}

	defer r.Request.Body.Close()

	var form forms.SetMembership
	if err := json.NewDecoder(r.Request.Body).Decode(&form); err != nil {
		return responses.APIError(errors.InvalidJSON(err))
	}

	// The client needs org admin access to edit the isAdmin status.
	if !isOrgAdmin && form.IsAdmin != nil {
		return responses.APIError(errors.NotAuthorized("must have organization admin access to edit membership admin status"))
	}

	// Can't alter admin status if the admin members are synced with LDAP.
	if form.IsAdmin != nil && rd.Org.AdminSyncConfig.EnableSync && rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.LdapPrecludes("cannot set membership admin status if membership is synced with LDAP"))
	}

	if err := s.schemaMgr.AddOrgMembership(rd.Org.ID, rd.User.ID, form.IsAdmin, form.IsPublic); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	orgMembership, err := s.schemaMgr.GetOrgMembership(rd.Org.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	memberInfo := schema.MemberInfo{
		Member:   *rd.User,
		IsAdmin:  orgMembership.IsAdmin,
		IsPublic: orgMembership.IsPublic,
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMember(&memberInfo))
}

// RouteGetOrganizationMembership returns a route describing the
// GetOrganizationMembership endpoint.
func (s *Service) routeGetOrganizationMembership() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetOrganizationMembership),
		Doc:     "Details of a user's membership in an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the membership will be retrieved",
			"memberNameOrID": "Name or id of user whose membership will be retrieved",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, membership returned.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleGetOrganizationMembership handles a request for getting a member of an
// organization.
func (s *Service) handleGetOrganizationMembership(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.MakeFilterGetUser(memberNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	member, err := s.schemaMgr.GetOrgMembership(rd.Org.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	// To have access, the membership must be public, the client must be
	// the user in question, or the client must be a member of the org.
	if !((member != nil && member.IsPublic) || clientAccount.ID == rd.User.ID) {
		orgAccess, err := s.authorizer.OrgMembershipAccess(rd.Org.ID, clientAccount)
		if err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}

		// Value is nil if the client user is not a member of the org.
		if orgAccess == nil {
			return responses.APIError(errors.NotAuthorized("must have org member access, authenticate as the user in question, or the membership must be public"))
		}
	}

	if member == nil {
		return responses.APIError(errors.NoSuchMember(memberNameOrID))
	}

	memberInfo := schema.MemberInfo{
		Member:   *rd.User,
		IsAdmin:  member.IsAdmin,
		IsPublic: member.IsPublic,
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMember(&memberInfo))
}

// RouteDeleteOrganizationMember returns a route describing the
// DeleteOrganizationMember endpoint.
func (s *Service) routeDeleteOrganizationMember() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{orgNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleDeleteOrganizationMember),
		Doc:     "Remove a member from an organization",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the membership will be deleted",
			"memberNameOrID": "Name or id of user whose membership will be deleted",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, membership removed.",
			},
		},
	}
}

// HandleDeleteOrganizationMember handles a request for removing a member of an
// organization.
func (s *Service) handleDeleteOrganizationMember(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
		rd.MakeFilterGetUser(memberNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	if err := s.schemaMgr.DeleteOrgMembership(rd.Org.ID, rd.User.ID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
