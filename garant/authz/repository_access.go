package authz

import (
	"fmt"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"

	enziresponses "github.com/docker/orca/enzi/api/responses"
)

// RepositoryAccess handles Repository Authorization:
//
// To determine the level of access that a user has to a given repository, the
// following flow is used:
//
// 	Note: If the repository has "public" visibility, return at least "read"
// 	level access if no access would be returned otherwise.
//
// 	Get the user's access to the repository's namespace using the method
// 	implemented in (*Authorizer).NamespaceAccess.
//
// 	Is the namespace access level "admin"?
// 		yes? -> return "admin" level access
// 		no? -> continue
// 	Note: this handles the case of the user being a global admin, the
// 	repository being owned by the user, the user being an owner of the
// 	organization (if an organization's repository), and other global roles.
//
// 	If the user is in the global "read-write" or "read-only" role, remember
// 	that (highest) corresponding level of access.
//
// 	Is the repository owned by some other user account?
// 		yes? -> return the higher of the namespace access (if any)
// 			and the owner's granted access to this repository for
// 			the user (if any).
// 		no? -> continue
//
// 	Then the repository must be owned by an organization account. The
// 	teams of this organization that the user is a member of are then
// 	queried:
//
// 	Get the highest level of repository access that the user's teams are
// 	granted to this repository. If that team level of access is higher
// 	than the user's namespace access level (if any), return that, otherwise
// 	return the user's namespace access (which may be none).
func (a *authorizer) RepositoryAccess(user *authn.User, repo *schema.Repository, ns *enziresponses.Account) (accessLevel string, err error) {
	if repo.Visibility == schema.RepositoryVisibilityPublic {
		// Grant "read" access at the very least.
		accessLevel = schema.AccessLevelReadOnly
	}

	if user.IsAnonymous {
		// Anonymous users can have access only if the repo is public.
		return accessLevel, nil
	}

	namespaceAccessLevel, err := a.NamespaceAccess(user, ns)
	if err != nil {
		return "", fmt.Errorf("unable to get repository namespace access level: %s", err)
	}

	// Compare with possible public "read" access.
	accessLevel = schema.HighestRankingAccessLevel(accessLevel, namespaceAccessLevel)

	if accessLevel == schema.AccessLevelAdmin {
		// The user is either a global admin, namespace owner, or has "admin"
		// access to the repository namespace via teams.
		return schema.AccessLevelAdmin, nil
	}

	if !ns.IsOrg {
		// The namespace belongs to some other user account.
		return accessLevel, nil
	}

	member, err := user.EnziSession.GetOrganizationMember("id:"+ns.ID, "id:"+user.Account.ID)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user's organization membership details: %s", err)
	}

	if member == nil {
		return accessLevel, nil
	}

	if member.IsAdmin {
		return schema.AccessLevelAdmin, nil
	}

	memberTeams, _, err := user.EnziSession.ListOrganizationMemberTeams("id:"+ns.ID, "id:"+user.Account.ID, "", 0)
	if err != nil {
		return "", fmt.Errorf("unable to determine if user is a member of the owners team: %s", err)
	}

	rtas, _, err := a.accessMgr.ListTeamsWithAccessToRepository(repo.ID, "", 0)
	if err != nil {
		return "", fmt.Errorf("unable to list team accesses for repository: %s", err)
	}

	rtaMap := make(map[string]schema.RepositoryTeamAccess, len(rtas)) // map from teamID to RepositoryTeamAccess
	for _, rta := range rtas {
		rtaMap[rta.TeamID] = rta
	}

	for _, memberTeam := range memberTeams.MemberTeams {
		if access, ok := rtaMap[memberTeam.Team.ID]; ok {
			accessLevel = schema.HighestRankingAccessLevel(accessLevel, access.AccessLevel)
		}
	}

	return accessLevel, nil
}
