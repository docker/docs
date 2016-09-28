package authz

import (
	"errors"
	"fmt"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"

	enziresponses "github.com/docker/orca/enzi/api/responses"
)

var (
	// ErrNoSuchNamespace indicates that the specified namespace does not exist
	ErrNoSuchNamespace = errors.New("the specified namespace does not exist")
)

// NamespaceAccess handles Repository Namespace Authorization:
//
// To determine the level of access that a user has to a given repository
// namespace, the following flow is used:
//
// 	Is the user a global admin?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	If the user is in the global "read-write" or "read-only" role, remember
// 	that (highest) corresponding level of access.
//
// 	Is the namespace the user's own (same name)?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	Is the namespace owned by some other user account?
// 		yes? -> return the highest global access level (if any) or
// 			none.
// 		no? -> continue
//
// 	Then the namespace must be owned by an organization account. The teams
// 	of this organization that the user is a member of are then queried:
//
// 	Is the user a member of the "owners" team in this organization?
// 		yes? -> return "admin" level access
// 		no? -> continue
//
// 	Get the highest level of namespace access that the user's teams are
// 	granted to this namespace. If that team level of access is higher than
// 	the user's global access level (if any), return that, otherwise return
// 	the user's global access level (which may be none).
func (a *authorizer) NamespaceAccess(user *authn.User, ns *enziresponses.Account) (accessLevel string, err error) {
	if *user.Account.IsAdmin || ns.Name == user.Account.Name {
		// User is a global admin or direct owner of this namespace.
		return schema.AccessLevelAdmin, nil
	}

	if !ns.IsOrg {
		// The namespace belongs to some other user account.
		return "", nil
	}

	member, err := user.EnziSession.GetOrganizationMember("id:"+ns.ID, "id:"+user.Account.ID)
	if err != nil {
		return "", fmt.Errorf("unable to determine if user is a member of the owners team: %s", err)
	} else if member == nil {
		return "", nil
	}

	if member.IsAdmin {
		return schema.AccessLevelAdmin, nil
	}

	//TODO paging
	memberTeams, _, err := user.EnziSession.ListOrganizationMemberTeams("id:"+ns.ID, "id:"+user.Account.ID, "", 0)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user's organization membership details: %s", err)
	}

	ntas, _, err := a.accessMgr.ListTeamsWithAccessToNamespace(ns.ID, "", 0)
	if err != nil {
		return "", fmt.Errorf("unable to list team accesses for namespace: %s", err)
	}

	ntaMap := make(map[string]schema.NamespaceTeamAccess, len(ntas)) // map from teamID to NamespaceTeamAccess
	for _, nta := range ntas {
		ntaMap[nta.TeamID] = nta
	}

	for _, memberTeam := range memberTeams.MemberTeams {
		if access, ok := ntaMap[memberTeam.Team.ID]; ok {
			accessLevel = schema.HighestRankingAccessLevel(accessLevel, access.AccessLevel)
		}
	}

	return accessLevel, nil
}
