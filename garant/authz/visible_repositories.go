package authz

import (
	"fmt"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"

	enziresponses "github.com/docker/orca/enzi/api/responses"
)

func (a *authorizer) AllVisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account) ([]*schema.Repository, error) {
	visibleRepositories, _, err := a.VisibleRepositoriesInNamespace(user, ns, "", 0)
	return visibleRepositories, err
}

// VisibleRepositoriesInNamespace lists repositories in the given namespace
// which are visible to the given user.
//
// To determine this list of repositories, the following flow is used:
//
// 	Does the user have at least global "read" access?
// 	  yes?  return all repositories in the namespace
// 	  no?   continue
//
// 	Get the user's namespace access level.
//
// 	Is the namespace access level at least "read" access?
// 	  yes?  return all repositories in the namespace
// 	  no?   continue
//
// 	Get a list of all of the public repositories in the namespace. The user
// 	will have visibility to at least these.
//
// 	Determine the repositories that the user has been granted explicit
// 	access to in this namespace. A list of repository user access should
// 	be retrieved using the following flow:
//
//	  Is the namespace owned by a user or an organization?
//	    user?  get the repositories that the namespace owner has explicitly
// 	           granted access to this user.
//	    org?   get the repositories that the user has been granted
// 	           access to via team membership in this organization.
//
// 	Return the merged set of public repositories in this namespace with
// 	those that the user is granted explicit access to.
//
//  For pagination we take the first `limit` results after sorting by id
//  and discard any extras that were returned. We also make sure we correctly
//  set the `next` ID depending on whether or not there are more results on
//  the next page.
func (a *authorizer) VisibleRepositoriesInNamespace(user *authn.User, ns *enziresponses.Account, startPublicID string, limit uint) (visibleRepositories []*schema.Repository, nextPublicID string, err error) {
	if *user.Account.IsAdmin {
		return a.repoMgr.ListRepositoriesInNamespace(ns.ID, startPublicID, limit)
	}

	isNamespaceAdmin := false
	if ns.IsOrg {
		if member, err := user.EnziSession.GetOrganizationMember("id:"+ns.ID, "id:"+user.Account.ID); err != nil {
			return nil, "", fmt.Errorf("Could not check organization membership: %s:", err)
		} else if member != nil {
			isNamespaceAdmin = member.IsAdmin
		}
	}

	nextPublicID = startPublicID
	for uint(len(visibleRepositories)) <= limit || limit == 0 {
		repos, nextPublicID, err := a.repoMgr.ListRepositoriesInNamespace(ns.ID, nextPublicID, limit)
		if err != nil {
			return nil, "", err
		}
		newVisibleRepos, err := a.FilterVisibleReposInNamespace(user, ns, repos, isNamespaceAdmin)
		if err != nil {
			return nil, "", err
		}
		visibleRepositories = append(visibleRepositories, newVisibleRepos...)
		if nextPublicID == "" {
			break
		}
	}

	if limit != 0 && limit < uint(len(visibleRepositories)) {
		nextPublicID = visibleRepositories[limit].PublicID()
		visibleRepositories = visibleRepositories[:limit]
	}

	return visibleRepositories, nextPublicID, nil
}
