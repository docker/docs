package authz

import (
	"fmt"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
)

// SharedRepositoriesForUser lists repositories which are "shared" with the
// given user, i.e., those repos not owned by the user but which the user has
// some granted access to (including public access) either via user repo
// access grants, team repo access grants, or team repo namespace access
// grants. Global access is not taken into consideration as global access would
// simply list all repositories in the system.
func (a *authorizer) SharedRepositoriesForUser(user *authn.User, startPublicID string, limit uint) (sharedRepositories []*schema.Repository, nextPublicID string, err error) {
	namespaces := []string{user.Account.ID}
	adminNamespaces := map[string]*responses.Account{user.Account.ID: user.Account}
	// TODO paging
	memberOrgs, _, err := user.EnziSession.ListUserOrganizations("id:"+user.Account.ID, "", 0)
	if err != nil {
		return nil, "", err
	}
	for _, memberOrg := range memberOrgs.MemberOrgs {
		namespaces = append(namespaces, memberOrg.Org.ID)
		if memberOrg.IsAdmin {
			adminNamespaces[memberOrg.Org.ID] = &memberOrg.Org
		}
	}
	nextPublicID = startPublicID
	for uint(len(sharedRepositories)) <= limit || limit == 0 {
		// Get a page of the repos under the user's namespace
		var allPossibleRepos []*schema.Repository
		allPossibleRepos, nextPublicID, err = a.repoMgr.ListRepositoriesInNamespacesOrPublic(namespaces, nextPublicID, limit)
		if err != nil {
			return nil, "", err
		}
		// Add what we can...
		namespaceRepoGrouping := groupReposByNamespace(allPossibleRepos)
		for namespaceID, namespaceRepos := range namespaceRepoGrouping {
			if limit != 0 && limit < uint(len(sharedRepositories)) {
				break
			}

			nsObj, isNamespaceAdmin := adminNamespaces[namespaceID]
			if !isNamespaceAdmin {
				var err error
				if nsObj, err = user.EnziSession.GetAccount("id:" + namespaceID); err != nil {
					apiErrs, _ := err.(*enzierrors.APIErrors)
					if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_ACCOUNT") {
						// TODO maybe remove all the repos under this
						// account?
						continue
					} else {
						return nil, "", err
					}
				}
			}
			newVisibleRepos, err := a.FilterVisibleReposInNamespace(user, nsObj, namespaceRepos, isNamespaceAdmin)
			if err != nil {
				return nil, "", err
			}
			sharedRepositories = append(sharedRepositories, newVisibleRepos...)
		}
		if nextPublicID == "" {
			break
		}
	}

	if limit != 0 && limit < uint(len(sharedRepositories)) {
		nextPublicID = sharedRepositories[limit].PublicID()
		sharedRepositories = sharedRepositories[:limit]
	}
	return sharedRepositories, nextPublicID, nil
}

func (a *authorizer) FilterVisibleReposInNamespace(user *authn.User, namespace *responses.Account, repos []*schema.Repository, isNamespaceAdmin bool) ([]*schema.Repository, error) {
	if isNamespaceAdmin || user.Account.ID == namespace.ID {
		return repos, nil
	}
	// check if member of team with admin, read-only, or read-write access
	// to namespace
	// TODO paginate??
	if namespaceTeamAccesses, _, err := a.accessMgr.ListTeamsWithAccessToNamespace(namespace.ID, "", 0); err != nil {
		return nil, fmt.Errorf("unable to retrieve team accesses to namespace: %s", err)
	} else {
		for _, nsa := range namespaceTeamAccesses {
			if _, err := user.EnziSession.GetTeamMember("id:"+namespace.ID, "id:"+nsa.TeamID, "id:"+user.Account.ID); err != nil {
				apiErrs, _ := err.(*enzierrors.APIErrors)
				if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_MEMBER") {
					continue
				}
				return nil, fmt.Errorf("failed to check team membership from enzi: %s", err)
			}
			if schema.AccessLevelAtLeast(nsa.AccessLevel, schema.AccessLevelReadOnly) {
				return repos, nil
			}
		}
	}

	var memberTeams []responses.MemberTeam
	if namespace.IsOrg {
		if memberTeamsResp, _, err := user.EnziSession.ListOrganizationMemberTeams("id:"+namespace.ID, "id:"+user.Account.ID, "", 0); err != nil {
			apiErrs, _ := err.(*enzierrors.APIErrors)
			if dtrutil.CheckContainsEnziError(apiErrs, "NO_SUCH_MEMBER") {
				// TODO handle bad orgs (remove if NO_SUCH_ACCOUNT)
				return nil, fmt.Errorf("unable to retrieve user's organization membership details: %s", err)
			}
		} else {
			memberTeams = memberTeamsResp.MemberTeams
		}
	}
	var teamIDs []string
	for _, memberTeam := range memberTeams {
		if memberTeam.IsAdmin {
			return repos, nil
		}
		teamIDs = append(teamIDs, memberTeam.Team.ID)
	}

	accesses, err := a.accessMgr.ListRepositoryAccessForTeamsAndRepositories(teamIDs, repos)
	if err != nil {
		return nil, err
	}
	teamVisibleRepos := make(map[string]struct{}, len(accesses))
	for _, access := range accesses {
		teamVisibleRepos[access.RepositoryID] = struct{}{}
	}
	var visibleRepos []*schema.Repository
	for _, repo := range repos {
		_, isTeamVisible := teamVisibleRepos[repo.ID]
		if repo.Visibility == schema.RepositoryVisibilityPublic || isTeamVisible {
			visibleRepos = append(visibleRepos, repo)
		}
	}
	return visibleRepos, nil
}

func groupReposByNamespace(repos []*schema.Repository) map[string][]*schema.Repository {
	if len(repos) == 0 {
		return nil
	}
	grouping := map[string][]*schema.Repository{}
	currentNamespaceID := repos[0].NamespaceAccountID
	i, j := 0, 0
	for ; j < len(repos); j++ {
		if currentNamespaceID != repos[j].NamespaceAccountID {
			grouping[currentNamespaceID] = repos[i:j]
			currentNamespaceID = repos[j].NamespaceAccountID
			i = j
		}
	}
	grouping[currentNamespaceID] = repos[i:j]
	return grouping
}
