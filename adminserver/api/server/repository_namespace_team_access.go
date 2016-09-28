package server

import (
	"encoding/json"
	"net/http"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/manager/schema"

	"github.com/docker/distribution/context"
	"github.com/emicklei/go-restful"
)

// handleListRepoNamespaceTeamAccess lists teams and their granted level of
// access to the specified org-owned repository namespace. The client must be
// authenticated as a user with "admin" access to the namespace.
func (a *APIServer) handleListRepoNamespaceTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetOrganization(namespace),
		makeFilterGetRepoNamespace(namespace),
		ensureUserIsSuperuserOrOrgMember,
		getPagerParams,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// Go ahead and list the repository namespace team access.
	ntas, next, err := a.repoAccessMgr.ListTeamsWithAccessToNamespace(rd.namespace.ID, rd.start, rd.limit)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	// Get whether the client user is a member of the teams.
	memberTeams, _, err := rd.user.EnziSession.ListOrganizationMemberTeams(namespace, "id:"+rd.user.Account.ID, "", 0)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	membershipSet := map[string]struct{}{} // set of TeamIDs of teams the user is a member of
	for _, memberTeam := range memberTeams.MemberTeams {
		membershipSet[memberTeam.Team.ID] = struct{}{}
	}

	teamAccessList := make([]responses.TeamAccess, 0, len(ntas))
	for _, nta := range ntas {
		_, isMember := membershipSet[nta.TeamID]
		teamAccessList = append(teamAccessList, responses.MakeTeamAccessForNamespace(&nta, isMember))
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.ListRepoNamespaceTeamAccess{responses.MakeNamespace(rd.namespace), teamAccessList}, r, next, 0)
}

// handleGetRepoNamespaceTeamAccess looks up the granted level of access for a
// team on the specified org-owned repository namespace. The client must be
// authenticated as a user with "admin" access to the namespace or be a member
// of the team in question.
func (a *APIServer) handleGetRepoNamespaceTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetOrganization(namespace),
		ensureUserIsSuperuserOrOrgMember,
		makeFilterGetTeam(teamName),
		makeFilterGetRepoNamespace(namespace),
		getUserIsTeamMember,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// Go ahead and look up the namespace access that the grantee team has.
	nta, err := a.repoAccessMgr.GetNamespaceAccessForTeam(rd.namespace.ID, rd.team.ID)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeNamespaceTeamAccess(nta, rd.isTeamMember, rd.namespace))
}

// handleGrantRepoNamespaceTeamAccess grants a team some access to the
// specified org-owned repository namespace. The client must be authenticated
// as a user with "admin" access to the namespace.
func (a *APIServer) handleGrantRepoNamespaceTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetOrganization(namespace),
		makeFilterGetRepoNamespace(namespace),
		getRepoNamespaceAccess,
		ensureRepoNamespaceAccessLevelAtLeastAdmin,
		makeFilterGetTeam(teamName),
		getUserIsTeamMember,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// Parse the request body for an access level.
	form := forms.Access{}
	if err := json.NewDecoder(r.Request.Body).Decode(&form); err != nil {
		return responses.APIError(errors.MakeError(errors.ErrorCodeInvalidJSON, err))
	}

	if errResponse := validateAccessLevel(form.AccessLevel); errResponse != nil {
		return errResponse
	}

	nta := &schema.NamespaceTeamAccess{
		NamespaceID: rd.namespace.ID,
		TeamID:      rd.team.ID,
		AccessLevel: form.AccessLevel,
	}

	if err := a.repoAccessMgr.SetNamespaceTeamAccess(nta); err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeNamespaceTeamAccess(nta, rd.isTeamMember, rd.namespace))
}

// handleRevokeRepoNamespaceTeamAccess revokes a team's access to the
// specified org-owned repository namespace. The client must be authenticated
// as a user with "admin" access to the namespace.
func (a *APIServer) handleRevokeRepoNamespaceTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetOrganization(namespace),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetTeam(teamName),
		getRepoNamespaceAccess,
		ensureRepoNamespaceAccessLevelAtLeastAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// This operation is idempotent so it's okay if the grantee team does
	// not exist or does not have access to begin with.
	if err := a.repoAccessMgr.DeleteNamespaceTeamAccess(rd.namespace.ID, rd.team.ID); err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil, nil, nil)
}
