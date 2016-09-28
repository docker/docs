package server

import (
	"fmt"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"

	"github.com/docker/distribution/context"
	enzierrors "github.com/docker/orca/enzi/api/errors"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/emicklei/go-restful"
)

type requestData struct {
	a   *APIServer
	ctx context.Context
	r   *restful.Request

	user         *authn.User
	org          *enziresponses.Account
	team         *enziresponses.Team
	repo         *schema.Repository
	namespace    *enziresponses.Account
	accessLevel  string
	authMethod   string
	start        string
	limit        uint
	isTeamMember bool

	filters     []func(*requestData)
	errResponse responses.APIResponse
}

func newRequestData(a *APIServer, ctx context.Context, r *restful.Request) *requestData {
	return &requestData{
		a:   a,
		ctx: ctx,
		r:   r,
	}
}

// Append the given function to the filters list on this requestData object.
func (rd *requestData) addFilters(filters ...func(*requestData)) *requestData {
	rd.filters = append(rd.filters, filters...)
	return rd
}

// Evaluates all filters in the requestData filter list, stopping if any error
// is encountered. The errResponse field of this requestData object will be set
// if any error is encountered.
func (rd *requestData) evaluateFilters() {
	for _, filter := range rd.filters {
		if filter(rd); rd.errResponse != nil {
			break
		}
	}
}

// Set the authenticated user on the requestData object. If authentication is
// required, it is an error if the client is anonymous.
func makeFilterGetAuthenticatedUser(authenticationRequired bool) func(*requestData) {
	return func(rd *requestData) {
		rd.user, rd.errResponse = rd.a.getAuthenticatedUser(rd.r, authenticationRequired)
	}
}

// Get the organization in question or set a NotFound error response.
func makeFilterGetOrganization(orgName string) func(*requestData) {
	return func(rd *requestData) {
		rd.org, rd.errResponse = rd.a.getOrgOrNotFound(rd.ctx, rd.user, orgName)
	}
}

// Get the team in question in the current organization or NotFound error.
func makeFilterGetTeam(teamName string) func(*requestData) {
	return func(rd *requestData) {
		rd.team, rd.errResponse = rd.a.getTeamOrNotFound(rd.ctx, rd.user, rd.org.ID, teamName)
	}
}

// Set the repo field on the requestData object or produce a NotFound error.
// This must be called after makeFilterGetRepoNamespace.
func makeFilterGetRepository(repoName string) func(*requestData) {
	return func(rd *requestData) {
		rd.repo, rd.errResponse = rd.a.getRepositoryOrNotFound(rd.ctx, rd.namespace, repoName)
	}
}

func makeFilterGetRepoNamespace(namespace string) func(*requestData) {
	return func(rd *requestData) {
		rd.namespace, rd.errResponse = rd.a.getRepoNamespaceOrNotFound(rd.ctx, rd.user, namespace)
	}
}

// Attempt to get the access level that the authenticated user has on the
// repository from the requestData.
func getRepositoryAccess(rd *requestData) {
	var err error
	if rd.accessLevel, err = rd.a.authorizer.RepositoryAccess(rd.user, rd.repo, rd.namespace); err != nil {
		rd.errResponse = responses.APIError(errors.InternalError(rd.ctx, err))
	}
}

func getPagerParams(rd *requestData) {
	rd.start, rd.limit = pagerParams(rd.r)
}

// Get the user's access to the repository namespace.
func getRepoNamespaceAccess(rd *requestData) {
	var err error
	rd.accessLevel, err = rd.a.authorizer.NamespaceAccess(rd.user, rd.namespace)
	if err != nil {
		rd.errResponse = responses.APIError(errors.InternalError(rd.ctx, err))
	}
}

// If the user has no access, return Not Found - No such repository.
func ensureAccessLevelAtLeastReadOnly(rd *requestData) {
	if !schema.AccessLevelAtLeast(rd.accessLevel, schema.AccessLevelReadOnly) {
		rd.errResponse = responses.APIError(errors.NoSuchRepositoryError(rd.namespace.Name, rd.repo.Name))
	}
}

// If the user doesn't have at least "read-only" access, results in NotFound.
// If the user doesn't have at least "read-write" access, results in Forbidden.
func ensureAccessLevelAtLeastReadWrite(rd *requestData) {
	ensureAccessLevelAtLeastReadOnly(rd)
	if rd.errResponse == nil && !schema.AccessLevelAtLeast(rd.accessLevel, schema.AccessLevelReadWrite) {
		rd.errResponse = responses.APIError(errors.NotAuthorizedError("The client must be authenticated as a user with read-write level access to the repository."))
	}
}

// If the user doesn't have at least "read-only" access, results in NotFound.
// If the user doesn't have at least "admin" access, results in Forbidden.
func ensureAccessLevelAtLeastAdmin(rd *requestData) {
	ensureAccessLevelAtLeastReadOnly(rd)
	if rd.errResponse == nil && !schema.AccessLevelAtLeast(rd.accessLevel, schema.AccessLevelAdmin) {
		rd.errResponse = responses.APIError(errors.NotAuthorizedError("The client must be authenticated as a user with admin level access to the repository."))
	}
}

func ensureIsAdmin(rd *requestData) {
	if rd.errResponse == nil && !*rd.user.Account.IsAdmin {
		rd.errResponse = responses.APIError(errors.NotAuthorizedError("The client must be authenticated as a user with admin level access."))
	}
}

// If, in the scope of a repository namespace, the user does not have admin
// access, results in a Forbidden error.
func ensureRepoNamespaceAccessLevelAtLeastAdmin(rd *requestData) {
	if !schema.AccessLevelAtLeast(rd.accessLevel, schema.AccessLevelAdmin) {
		rd.errResponse = responses.APIError(errors.NotAuthorizedError("The client must be authenticated as a user with admin level access to the repository namespace."))
	}
}

// Ensure that the repository is owned by an organization account. If not, then
// return a "400 Bad Request" because the repository must be owned by an org.
func ensureRepoIsOrgOwned(rd *requestData) {
	rd.errResponse = validateRepoIsOrgOwned(rd.namespace)
}

// Ensure that the client user has global admin access or is in the "owners"
// team of the organization.
func ensureUserIsSuperuserOrOrgOwner(rd *requestData) {
	rd.errResponse = rd.a.validateAdminOrOrgOwner(rd.ctx, rd.user, rd.org)
}

// Ensure that the client is a global admin or org member.
func ensureUserIsSuperuserOrOrgMember(rd *requestData) {
	rd.errResponse = rd.a.validateAdminOrOrgMember(rd.ctx, rd.user, rd.org.ID)
}

// Get whether the user is a member of the team in the namespace
func getUserIsTeamMember(rd *requestData) {
	rd.isTeamMember, rd.errResponse = rd.a.getUserIsTeamMember(rd.ctx, rd.user, rd.namespace.ID, rd.team.ID)
}

// Ensure that the user is a member of the team in question OR that the
// user's access level to the namespace is "admin".
func ensureUserIsTeamMeberOrNamespaceAdmin(rd *requestData) {
	// Check if the user is a member of the team
	// query so we do this first and may not have to check if they are a
	// namespace admin.
	if _, err := rd.user.EnziSession.GetTeamMember("id:"+rd.org.ID, "id:"+rd.team.ID, "id:"+rd.user.Account.ID); err != nil {
		apiErrs, ok := err.(*enzierrors.APIErrors)
		if ok {
			for _, apiErr := range apiErrs.Errors {
				switch apiErr.Code {
				case "NO_SUCH_ACCOUNT", "NO_SUCH_TEAM", "NO_SUCH_MEMBER":
				// these are okay. just continue??
				// TODO maybe throw an error or handle removed
				// namespace/teams?
				default:
					rd.errResponse = responses.APIError(errors.InternalError(rd.ctx, apiErr))
					return
				}
			}
		} else {
			rd.errResponse = responses.APIError(errors.InternalError(rd.ctx, fmt.Errorf("received a non-enzi error from GetTeamMember: %s", err)))
			return
		}
	} else {
		// user is a team member
		return
	}

	// Get the user's access to the repository namespace.
	accessLevel, err := rd.a.authorizer.NamespaceAccess(rd.user, rd.namespace)
	if err != nil {
		rd.errResponse = responses.APIError(errors.InternalError(rd.ctx, err))
	} else if !schema.AccessLevelAtLeast(accessLevel, schema.AccessLevelAdmin) {
		rd.errResponse = responses.APIError(errors.NotAuthorizedError("The client must be authenticated as a user with admin level access to the repository namespace or be a member of the team in question."))
	}
}
