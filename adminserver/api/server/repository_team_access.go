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

// handleListRepoTeamAccess lists teams and their granted level of access for
// the specified org-owned repository. The client must be authenticated as a
// user with "admin" access to the repository. If the repository is not visible
// to the client, then a "404 Not Found" should be returned.
func (a *APIServer) handleListRepoTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		ensureRepoIsOrgOwned,
		getRepositoryAccess,
		ensureAccessLevelAtLeastAdmin,
		getPagerParams,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// Go ahead and list the repository team access.
	rtas, next, err := a.repoAccessMgr.ListTeamsWithAccessToRepository(rd.repo.ID, rd.start, rd.limit)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	// Get whether the client user is a member of the teams.
	memberTeams, _, err := rd.user.EnziSession.ListOrganizationMemberTeams("id:"+rd.namespace.ID, "id:"+rd.user.Account.ID, "", 0)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}
	membershipSet := map[string]struct{}{}
	for _, memberTeam := range memberTeams.MemberTeams {
		membershipSet[memberTeam.Team.ID] = struct{}{}
	}

	teamAccessList := make([]responses.TeamAccess, len(rtas))
	for i, rta := range rtas {
		_, isMember := membershipSet[rta.TeamID]
		teamAccessList[i] = responses.MakeTeamAccess(&rta, isMember)
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.ListRepoTeamAccess{responses.MakeRepository(rd.namespace.Name, rd.namespace.IsOrg, rd.repo, false), teamAccessList}, r, next, 0)
}

// handleListTeamRepoAccess lists the repositories and the respective access
// levels that have been granted to the specified team. The client must be
// authenticated as either a system admin, a member of the organization's
// "owners" team or be a member of the team in question.
func (a *APIServer) handleListTeamRepoAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	orgname := vars["orgname"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetOrganization(orgname),
		ensureUserIsSuperuserOrOrgMember,
		makeFilterGetRepoNamespace(orgname),
		makeFilterGetTeam(teamName),
		ensureUserIsTeamMeberOrNamespaceAdmin,
		getPagerParams,
		getUserIsTeamMember,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// The user is authorized. Go ahead and list the repository team access.
	rtas, next, err := a.repoAccessMgr.ListRepositoryAccessForTeam(rd.team.ID, rd.start, rd.limit)
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	repoAccessList := make([]responses.RepoAccess, 0, len(rtas))
	for _, rta := range rtas {
		repoAccessList = append(repoAccessList, responses.MakeRepoAccess(rd.namespace.Name, rd.namespace.IsOrg, &rta))
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.ListTeamRepoAccess{responses.MakeTeam(rd.team.ID, rd.isTeamMember), repoAccessList}, r, next, 0)
}

// handleGrantRepoTeamAccess grants a team some access to the specified org-
// owned repository. The client must be authenticated as a user with "admin"
// access to the repository. If the repository is not visible to the client
// then a "404 Not Found" should be returned.
func (a *APIServer) handleGrantRepoTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		ensureRepoIsOrgOwned,
		getRepositoryAccess,
		ensureAccessLevelAtLeastAdmin,
		makeFilterGetOrganization(namespace),
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

	if err := a.repoAccessMgr.AddRepositoryTeamAccess(rd.repo.ID, rd.team.ID, form.AccessLevel); err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	rta := &schema.RepositoryTeamAccess{
		AccessLevel: form.AccessLevel,
		Repository:  *rd.repo,
		TeamID:      rd.team.ID,
	}
	return responses.JSONResponse(http.StatusOK, nil, nil, responses.MakeRepoTeamAccess(rd.namespace.Name, rd.namespace.IsOrg, rta, rd.isTeamMember))
}

// handleRevokeRepoTeamAccess revokes a team's access to the specified org-
// owned repository. The client must be authenticated as a user with "admin"
// access to the repository. If the repository is not visible to the client
// then a "404 Not Found" should be returned.
func (a *APIServer) handleRevokeRepoTeamAccess(ctx context.Context, r *restful.Request) responses.APIResponse {
	vars := r.PathParameters()
	namespace := vars["namespace"]
	repoName := vars["reponame"]
	teamName := vars["teamname"]

	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		makeFilterGetRepoNamespace(namespace),
		makeFilterGetRepository(repoName),
		ensureRepoIsOrgOwned,
		getRepositoryAccess,
		ensureAccessLevelAtLeastAdmin,
		makeFilterGetOrganization(namespace),
		makeFilterGetTeam(teamName),
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	// This operation is idempotent so it's okay if the grantee did not have
	// access to begin with.
	if err := a.repoAccessMgr.DeleteRepositoryTeamAccess(rd.repo.ID, rd.team.ID); err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil, nil, nil)
}
