package server

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/adminserver/util"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/garant/authz"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/distribution/context"
	enzierrors "github.com/docker/orca/enzi/api/errors"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/emicklei/go-restful"
)

var (
	validRepositoryNamePattern   = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*$`)
	maxRepoFullNameLength        = 255 // This is restricted by the registry.
	maxRepoShortDescriptonLength = 140 // Should be succinct enough to be a Tweet.
)

// getAuthenticatedUser retrieves the currently authenticated user for the
// given request. If required is true, then there must be an authenticatetd
// user (not anonymous). Returns a non-nil responses.APIResponse describing any
// error.
func (a *APIServer) getAuthenticatedUser(r *restful.Request, required bool) (user *authn.User, errResponse responses.APIResponse) {
	if user = util.GetAuthenticatedUser(r.Request); user == nil {
		errResponse = responses.JSONResponse(http.StatusInternalServerError, nil, nil, errors.ErrorCodeNotAuthenticated)
	} else if required && user.IsAnonymous {
		errResponse = responses.JSONResponse(http.StatusUnauthorized, nil, nil, errors.ErrorCodeNotAuthenticated)
	}

	return user, errResponse
}

// getUserOrNotFound attempts to retrieve the user account with the given name.
// If there is no such account or some other error occurs then a non-nil
// responses.APIResponse is returned indicating the error with appropriate status code.
// Returns a non-nil account on success.
func (a *APIServer) getUserOrNotFound(ctx context.Context, user *authn.User, name string) (acct *enziresponses.Account, errResponse responses.APIResponse) {
	acct, err := user.EnziSession.GetAccount(name)
	if err == nil {
		return acct, nil
	}

	apiErrs, _ := err.(*enzierrors.APIErrors)
	if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
		return nil, responses.APIError(errors.NoSuchAccountError(name))
	}

	return nil, responses.APIError(errors.InternalError(ctx, err))
}

func (a *APIServer) getOrgOrNotFound(ctx context.Context, user *authn.User, name string) (acct *enziresponses.Account, errResponse responses.APIResponse) {
	acct, err := user.EnziSession.GetAccount(name)
	if err == nil {
		return acct, nil
	}

	apiErrs, _ := err.(*enzierrors.APIErrors)
	if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
		return nil, responses.APIError(errors.NoSuchOrganizationError(name))
	}

	return nil, responses.APIError(errors.InternalError(ctx, err))
}

// validateRepositoryName validates that the given repository name matches
// our pattern for a valid repository.
func validateRepositoryName(namespace, name string) (errResponse responses.APIResponse) {
	fullName := fmt.Sprintf("%s/%s", namespace, name)
	if len(fullName) > maxRepoFullNameLength || !validRepositoryNamePattern.MatchString(name) {
		return responses.APIError(errors.ErrorCodeInvalidRepositoryName)
	}

	return nil
}

// validateRepoShortDescription validates that the given repository short
// description meets our maximum description length requirement.
func validateRepoShortDescription(description string) (errResponse responses.APIResponse) {
	if len(description) > maxRepoShortDescriptonLength {
		return responses.APIError(errors.InvalidRepositoryShortDescriptionError(maxRepoShortDescriptonLength))
	}

	return nil
}

// validateRepoVisibility validates that the given repository visibility value
// is either "public", or "private".
func validateRepoVisibility(visibility string) (errResponse responses.APIResponse) {
	switch visibility {
	case schema.RepositoryVisibilityPublic, schema.RepositoryVisibilityPrivate:
		return nil
	default:
		return responses.APIError(errors.ErrorCodeInvalidRepositoryVisibility)
	}
}

// validateAdminOrOrgOwner validates that the given user is either a system
// admin or a member of the owners group in the given organization.
func (a *APIServer) validateAdminOrOrgOwner(ctx context.Context, user *authn.User, org *enziresponses.Account) (errResponse responses.APIResponse) {
	err := a.authorizer.CheckAdminOrOrgOwner(user, org.ID)

	if err == nil {
		return nil
	}

	if err == authz.ErrNotAdminOrOrgOwner {
		return responses.APIError(errors.NotAuthorizedError("The client must be authenticated as either a system admin or organization owner."))
	}

	return responses.APIError(errors.InternalError(ctx, err))
}

// validateAdminOrOrgMember validates that the given user is either a system
// admin or a member of any team in the given organization.
func (a *APIServer) validateAdminOrOrgMember(ctx context.Context, user *authn.User, orgID string) (errResponse responses.APIResponse) {
	err := a.authorizer.CheckAdminOrOrgMember(user, orgID)

	if err == nil {
		return nil
	}

	if err == authz.ErrNotAdminOrOrgMember {
		return responses.APIError(
			errors.NotAuthorizedError("The client must be authenticated as either a system admin or an organization member."),
		)
	}

	return responses.APIError(errors.InternalError(ctx, err))
}

func (a *APIServer) getUserIsTeamMember(ctx context.Context, user *authn.User, namespaceID, teamID string) (isTeamMember bool, errResponse responses.APIResponse) {
	_, err := user.EnziSession.GetTeamMember("id:"+namespaceID, "id:"+teamID, "id:"+user.Account.ID)
	if err != nil {
		apiErrs, _ := err.(*enzierrors.APIErrors)
		if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_MEMBER") {
			return false, nil
		} else {
			return false, responses.APIError(errors.InternalError(ctx, err))
		}
	}
	return true, nil
}

func (a *APIServer) getTeamOrNotFound(ctx context.Context, user *authn.User, orgID, name string) (team *enziresponses.Team, errResponse responses.APIResponse) {
	team, err := user.EnziSession.GetTeam("id:"+orgID, name)
	if err == nil {
		return team, nil
	}

	apiErrs, _ := err.(*enzierrors.APIErrors)
	if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_TEAM") {
		return nil, responses.APIError(errors.NoSuchTeamError(name))
	}

	return nil, responses.APIError(errors.InternalError(ctx, err))
}

func (a *APIServer) getRepoNamespaceOrNotFound(ctx context.Context, user *authn.User, namespace string) (ns *enziresponses.Account, errResponse responses.APIResponse) {
	ns, err := user.EnziSession.GetAccount(namespace)
	if err == nil {
		return ns, nil
	}

	apiErrs, _ := err.(*enzierrors.APIErrors)
	if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
		return nil, responses.APIError(errors.NoSuchAccountError(namespace))
	}

	return nil, responses.APIError(errors.InternalError(ctx, err))
}

func (a *APIServer) getRepositoryOrNotFound(ctx context.Context, namespace *enziresponses.Account, name string) (repo *schema.Repository, errResponse responses.APIResponse) {
	repo, err := a.repoMgr.GetRepositoryByName(namespace.ID, name)

	if err == nil {
		return repo, nil
	}

	if err == schema.ErrNoSuchRepository {
		return nil, responses.APIError(errors.NoSuchRepositoryError(namespace.Name, name))
	}

	return nil, responses.APIError(errors.InternalError(ctx, err))
}

// validateAccessLevel returns an error APIRespnse if the given access level is
// not a valid choice of read-only, read-write, or admin.
func validateAccessLevel(accessLevel string) (errResponse responses.APIResponse) {
	switch accessLevel {
	case schema.AccessLevelReadOnly, schema.AccessLevelReadWrite, schema.AccessLevelAdmin:
		return nil // These are valid choices.
	default:
		return responses.APIError(errors.InvalidAccessLevelError(schema.AccessLevelReadOnly, schema.AccessLevelReadWrite, schema.AccessLevelAdmin))
	}
}

func validateRepoIsOrgOwned(ns *enziresponses.Account) (errResponse responses.APIResponse) {
	if ns.IsOrg {
		return nil
	}

	return responses.APIError(errors.MakeError(errors.ErrorCodeInvalidRepoContext, "The repository namespace must belong to an organization account."))
}
