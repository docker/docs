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
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// RouteCreateTeam returns a route describing the CreateTeam endpoint.
func (s *Service) routeCreateTeam() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/{orgNameOrID}/teams",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleCreateTeam),
		Doc:     "Create a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization in which the team will be created",
		},
		BodySample: forms.CreateTeam{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusCreated,
				Message: "Success, team created.",
				Sample:  responses.Team{},
			},
		},
	}
}

// HandleCreateTeam handles a request for creating a team in an organization.
func (s *Service) handleCreateTeam(ctx context.Context, clientUser *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientUser)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.CreateTeam)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	newTeam := &schema.Team{
		OrgID:       rd.Org.ID,
		Name:        form.Name,
		Description: form.Description,
	}

	if err := s.schemaMgr.CreateTeam(newTeam); err != nil {
		if err == schema.ErrTeamExists {
			return responses.APIError(errors.TeamExists())
		}
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusCreated, responses.MakeTeam(newTeam))
}

// RouteListTeams returns a route describing the ListTeams endpoint.
func (s *Service) routeListTeams() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListTeams),
		Doc:     "List teams in an organization",
		Notes:   "Lists teams in ascending order by name.",
		PathParameterDocs: map[string]string{
			"orgNameOrID": "Name or id of organization whose teams will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return teams with a name greater than or equal to this name."),
			restful.QueryParameter("limit", "Maximum number of teams per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of teams listed.",
				Sample:  responses.Teams{},
			},
		},
	}
}

// HandleListTeams handles a request for listing teams in an organization.
func (s *Service) handleListTeams(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
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

	startName, limit := helpers.PageParams(r, "start", "limit")

	teams, nextPageStart, err := s.schemaMgr.ListTeamsInOrg(rd.Org.ID, startName, limit)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeTeams(teams), r, nextPageStart)
}

// RouteGetTeam returns a route describing the GetTeam endpoint.
func (s *Service) routeGetTeam() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetTeam),
		Doc:     "Details for a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization in which the team will be retrieved",
			"teamNameOrID": "Name or id of team which will be retrieved",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, team returned.",
				Sample:  responses.Team{},
			},
		},
	}
}

// HandleGetTeam handles a request for getting details of a team in an
// organization.
func (s *Service) handleGetTeam(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeTeam(rd.Team))
}

// RouteUpdateTeam returns a route describing the UpdateTeam endpoint.
func (s *Service) routeUpdateTeam() server.Route {
	return server.Route{
		Method:  "PATCH",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleUpdateTeam),
		Doc:     "Update details for a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization in which the team will be updated",
			"teamNameOrID": "Name or id of team which will be updated",
		},
		BodySample: forms.UpdateTeam{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, team updated.",
				Sample:  responses.Team{},
			},
		},
	}
}

// HandleUpdateTeam handles a request for updating details of a team in an
// organization.
func (s *Service) handleUpdateTeam(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.UpdateTeam)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	var updateFields schema.TeamUpdateFields

	if form.Description != nil {
		rd.Team.Description = *form.Description
		updateFields.Description = form.Description
	}

	if err := s.schemaMgr.UpdateTeam(rd.Team.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	if form.Name != nil && *form.Name != rd.Team.Name {
		// The client wants to change the name of the team. This is a
		// separate task that requires changing the PK of the team.
		if err := s.schemaMgr.RenameTeam(rd.Team, *form.Name); err != nil {
			if err == schema.ErrTeamExists {
				return responses.APIError(errors.TeamExists())
			}
			return responses.APIError(errors.Internal(ctx, err))
		}
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeTeam(rd.Team))
}

// RouteDeleteTeam returns a route describing the DeleteTeam endpoint.
func (s *Service) routeDeleteTeam() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleDeleteTeam),
		Doc:     "Delete a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization in which the team will be deleted",
			"teamNameOrID": "Name or id of team which will be deleted",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, team deleted.",
			},
		},
	}
}

// HandleDeleteTeam handles a request for deleting a team in an organization.
func (s *Service) handleDeleteTeam(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgAdmin,
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	team, err := helpers.TeamByNameOrID(s.schemaMgr, rd.Org.ID, teamNameOrID)
	if err != nil {
		if err == schema.ErrNoSuchTeam {
			// Already deleted.
			return responses.JSONResponse(http.StatusNoContent, nil)
		}
		return responses.APIError(errors.Internal(ctx, err))
	}

	if err := s.schemaMgr.DeleteTeam(team.OrgID, team.Name); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}

// RouteListOrganizationMemberTeams returns a route describing the
// ListOrganizationMemberTeams endpoint.
func (s *Service) routeListOrganizationMemberTeams() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/members/{memberNameOrID}/teams",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListOrganizationMemberTeams),
		Doc:     "List a user's team membership in an organization",
		Notes:   "Lists team memberships in ascending order by team ID.",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the member's team memberships will be listed",
			"memberNameOrID": "Name or id of user whose memberships will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return team memberships with a team ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of team memberships per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of member's teams listed.",
				Sample:  responses.MemberTeams{},
			},
		},
	}
}

// HandleListOrganizationMemberTeams handles a request for listing a user's
// team memberships in an organization.
func (s *Service) handleListOrganizationMemberTeams(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
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

	// Only an admin or the user in question may list the user's team
	// memberships in an organization.
	if !(clientAccount.IsAdmin || clientAccount.ID == rd.User.ID) {
		return responses.APIError(errors.NotAuthorized("must be a system admin or the user whose teams are being listed"))
	}

	startID, limit := helpers.PageParams(r, "start", "limit")

	memberTeams, nextPageStart, err := s.schemaMgr.ListTeamsInOrgForUser(rd.Org.ID, rd.User.ID, startID, limit)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMemberTeams(memberTeams), r, nextPageStart)
}

// RouteListTeamMembers returns a route describing the ListTeamMembers
// endpoint.
func (s *Service) routeListTeamMembers() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/members",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListTeamMembers),
		Doc:     "List members of a team",
		Notes:   "Lists memberships in ascending order by user ID.",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization in which the team's members will be listed'",
			"teamNameOrID": "Name or id of team whose members will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("filter", "Filter members by type - either 'admins', 'non-admins', or 'all' (default).").DefaultValue("all"),
			restful.QueryParameter("start", "Only return members with a user ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of members per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of team members listed.",
				Sample:  responses.Members{},
			},
		},
	}
}

// HandleListTeamMembers handles a request for listing memberships in
// a team.
func (s *Service) handleListTeamMembers(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamMember,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	var (
		teamMembers    []schema.MemberInfo
		nextPageStart  string
		err            error
		startID, limit = helpers.PageParams(r, "start", "limit")
	)

	switch r.QueryParameter("filter") {
	case "admins":
		teamMembers, nextPageStart, err = s.schemaMgr.ListTeamAdmins(rd.Team.ID, startID, limit)
	case "non-admins":
		teamMembers, nextPageStart, err = s.schemaMgr.ListTeamNonAdmins(rd.Team.ID, startID, limit)
	default:
		teamMembers, nextPageStart, err = s.schemaMgr.ListTeamMembers(rd.Team.ID, startID, limit)
	}

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMembers(teamMembers), r, nextPageStart)
}

// RouteListTeamPublicMembers returns a route describing the
// ListTeamPublicMembers endpoint.
func (s *Service) routeListTeamPublicMembers() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/publicMembers",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListTeamPublicMembers),
		Doc:     "List public members of a team",
		Notes:   "Lists public members in ascending order by user ID.",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization in which the team's public members will be listed'",
			"teamNameOrID": "Name or id of team whose public members will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return members with a user ID greater than or equal to this ID."),
			restful.QueryParameter("limit", "Maximum number of members per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of public team members listed.",
				Sample:  responses.Members{},
			},
		},
	}
}

// HandleListTeamPublicMembers handles a request for listing public members of
// a team.
func (s *Service) handleListTeamPublicMembers(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	startID, limit := helpers.PageParams(r, "start", "limit")

	publicMembers, nextPageStart, err := s.schemaMgr.ListPublicTeamMembers(rd.Team.ID, startID, limit)

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeMembers(publicMembers), r, nextPageStart)
}

// routeGetTeamMemberSyncConfig returns a route describing the
// GetTeamMemberSyncConfig endpoint.
func (s *Service) routeGetTeamMemberSyncConfig() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/memberSyncConfig",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetTeamMemberSyncConfig),
		Doc:     "Get options for syncing members of a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization to which the team belongs",
			"teamNameOrID": "Name or id of team whose LDAP sync config will be retrieved",
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

// handleGetTeamMemberSyncConfig handles a request to get options for syncing
// members of a team.
func (s *Service) handleGetTeamMemberSyncConfig(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMemberSyncOpts(rd.Team.MemberSyncConfig))
}

// routeSetTeamMemberSyncConfig returns a route describing the
// routeSetTeamMemberSyncConfig endpoint.
func (s *Service) routeSetTeamMemberSyncConfig() server.Route {
	return server.Route{
		Method:  "PUT",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/memberSyncConfig",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleSetTeamMemberSyncConfig),
		Doc:     "Set options for syncing members of a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":  "Name or id of organization to which the team belongs",
			"teamNameOrID": "Name or id of team whose LDAP sync config will be set",
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

// handleSetTeamMemberSyncConfig handles a request to set options for syncing
// members of a team.
func (s *Service) handleSetTeamMemberSyncConfig(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.MemberSyncOpts)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	rd.Team.MemberSyncConfig = schema.MemberSyncOpts{
		EnableSync:         form.EnableSync,
		SelectGroupMembers: form.SelectGroupMembers,
		GroupDN:            form.GroupDN,
		GroupMemberAttr:    form.GroupMemberAttr,
		SearchBaseDN:       form.SearchBaseDN,
		SearchScopeSubtree: form.SearchScopeSubtree,
		SearchFilter:       form.SearchFilter,
	}

	updateFields := schema.TeamUpdateFields{
		MemberSyncConfig: &rd.Team.MemberSyncConfig,
	}

	if err := s.schemaMgr.UpdateTeam(rd.Team.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMemberSyncOpts(rd.Team.MemberSyncConfig))
}

// RouteAddTeamMember returns a route describing the AddTeamMember endpoint.
func (s *Service) routeAddTeamMember() server.Route {
	return server.Route{
		Method:  "PUT",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleAddTeamMember),
		Doc:     "Add a member to a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the team membership will be added",
			"teamNameOrID":   "Name or id of the team in which the membership will be added",
			"memberNameOrID": "Name or id of user which will be added as a member",
		},
		BodySample: forms.SetMembership{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, team membership set.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleAddTeamMember handles a request for adding a member to a team in an
// organization.
func (s *Service) handleAddTeamMember(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamAdmin,
		rd.MakeFilterGetUser(memberNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// If the team should be synced with LDAP, clients can't manually add
	// members to the team.
	if rd.Team.MemberSyncConfig.EnableSync && rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.LdapPrecludes("this team's membership can only be changed via LDAP syncing"))
	}

	// At this point, the client is authorized to add users to the team,
	// and the user in question is already a member of the organization.

	defer r.Request.Body.Close()

	var form forms.SetMembership
	if err := json.NewDecoder(r.Request.Body).Decode(&form); err != nil {
		return responses.APIError(errors.InvalidJSON(err))
	}

	if err := s.schemaMgr.AddTeamMembership(rd.Org.ID, rd.Team.ID, rd.User.ID, form.IsAdmin, form.IsPublic); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	teamMembership, err := s.schemaMgr.GetTeamMembership(rd.Team.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	memberInfo := schema.MemberInfo{
		Member:   *rd.User,
		IsAdmin:  teamMembership.IsAdmin,
		IsPublic: teamMembership.IsPublic,
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMember(&memberInfo))
}

// RouteUpdateTeamMembership returns a route describing the
// UpdateTeamMembership endpoint.
func (s *Service) routeUpdateTeamMembership() server.Route {
	return server.Route{
		Method:  "PATCH",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleUpdateTeamMembership),
		Doc:     "Update details of a user's membership in a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the team membership will be updated",
			"teamNameOrID":   "Name or id of the team in which the membership will be updated",
			"memberNameOrID": "Name or id of user whose team membership will be updated",
		},
		BodySample: forms.SetMembership{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, team membership updated.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleUpdateTeamMembership handles a request for updating attributes of a
// team membership.
func (s *Service) handleUpdateTeamMembership(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.MakeFilterGetUser(memberNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	isTeamAdmin := rd.ClientTeamAccess != nil && rd.ClientTeamAccess.IsAdmin
	if !(isTeamAdmin || clientAccount.ID == rd.User.ID) {
		return responses.APIError(errors.NotAuthorized("must have team admin access or authenticate as the user whose membership is being updated"))
	}

	if member, err := s.schemaMgr.GetTeamMembership(rd.Team.ID, rd.User.ID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	} else if member == nil {
		return responses.APIError(errors.NoSuchMember(memberNameOrID))
	}

	defer r.Request.Body.Close()

	var form forms.SetMembership
	if err := json.NewDecoder(r.Request.Body).Decode(&form); err != nil {
		return responses.APIError(errors.InvalidJSON(err))
	}

	// The client needs team admin access to edit the isAdmin status.
	if !isTeamAdmin && form.IsAdmin != nil {
		return responses.APIError(errors.NotAuthorized("must have team admin access to edit membership admin status"))
	}

	// Can't alter admin status if the team is synced with LDAP.
	if form.IsAdmin != nil && rd.Team.MemberSyncConfig.EnableSync && rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.LdapPrecludes("cannot set membership admin status if membership is synced with LDAP"))
	}

	if err := s.schemaMgr.AddTeamMembership(rd.Org.ID, rd.Team.ID, rd.User.ID, form.IsAdmin, form.IsPublic); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	teamMembership, err := s.schemaMgr.GetTeamMembership(rd.Team.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	memberInfo := schema.MemberInfo{
		Member:   *rd.User,
		IsAdmin:  teamMembership.IsAdmin,
		IsPublic: teamMembership.IsPublic,
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeMember(&memberInfo))
}

// RouteGetTeamMembership returns a route describing the GetTeamMembership
// endpoint.
func (s *Service) routeGetTeamMembership() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGeTeamMembership),
		Doc:     "Details of a user's membership in a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the team membership will be retrieved",
			"teamNameOrID":   "Name or id of the team in which the membership will be retrieved",
			"memberNameOrID": "Name or id of user whose team membership will be retrieved",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, team membership retuned.",
				Sample:  responses.Member{},
			},
		},
	}
}

// HandleGeTeamMembership handles a request for getting a member of a team.
func (s *Service) handleGeTeamMembership(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.MakeFilterGetUser(memberNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	member, err := s.schemaMgr.GetTeamMembership(rd.Team.ID, rd.User.ID)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	// To have access, the membership must be public, the client must be
	// the user in question, or the client must be a member of the team.
	if !((member != nil && member.IsPublic) || clientAccount.ID == rd.User.ID) {
		teamAccess, err := s.authorizer.TeamMembershipAccess(rd.Team.ID, *rd.ClientOrgAccess, clientAccount)
		if err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}

		// Value is nil if the client user is not a member of the team.
		if teamAccess == nil {
			return responses.APIError(errors.NotAuthorized("must have team member access, authenticate as the user in question, or the membership must be public"))
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

// RouteDeleteTeamMember returns a route describing the DeleteTeamMember
// endpoint.
func (s *Service) routeDeleteTeamMember() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{orgNameOrID}/teams/{teamNameOrID}/members/{memberNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleDeleteTeamMember),
		Doc:     "Remove a member from a team",
		PathParameterDocs: map[string]string{
			"orgNameOrID":    "Name or id of organization in which the team membership will be deleted",
			"teamNameOrID":   "Name or id of the team in which the membership will be deleted",
			"memberNameOrID": "Name or id of user whose team membership will be deleted",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, team membership deleted.",
			},
		},
	}
}

// HandleDeleteTeamMember handles a request for removing a member of a team.
func (s *Service) handleDeleteTeamMember(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	orgNameOrID := pathParams["orgNameOrID"]
	teamNameOrID := pathParams["teamNameOrID"]
	memberNameOrID := pathParams["memberNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetOrganization(orgNameOrID),
		rd.GetOrgAccess,
		rd.RequireOrgMember,
		rd.MakeFilterGetTeam(teamNameOrID),
		rd.GetTeamAccess,
		rd.RequireTeamAdmin,
		rd.MakeFilterGetUser(memberNameOrID),
		rd.GetAuthConfig,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	// If the team should be synced with LDAP, clients can't manually
	// delete members from the team.
	if rd.Team.MemberSyncConfig.EnableSync && rd.AuthConfig.Backend == config.AuthBackendLDAP {
		return responses.APIError(errors.LdapPrecludes("this team's membership can only be changed via LDAP syncing"))
	}

	if err := s.schemaMgr.DeleteTeamMembership(rd.Team.ID, rd.User.ID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
