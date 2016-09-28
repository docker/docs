package common

import (
	"github.com/docker/garant/auth"
)

// RepositoryScopeSets maps a repository name or prefix to a set of scopes.
type RepositoryScopeSets map[string]ScopeSet

// GetValidRepositoryScopeSets makes a set of valid repository scope sets from
// the given resource access list.
func GetValidRepositoryScopeSets(accessList []auth.Access) RepositoryScopeSets {
	repoScopeSets := make(RepositoryScopeSets, len(accessList))

	for _, access := range accessList {
		if access.Resource.Type != "repository" {
			continue // Ignore non-repository resources.
		}

		repo := access.Resource.Name

		scopes, ok := repoScopeSets[repo]
		if !ok {
			scopes = NewScopeSet()
			repoScopeSets[repo] = scopes
		}

		scopes.Add(access.Action)
	}

	return repoScopeSets
}

// MakeRepoAccessList makes a resource access list from the given repository
// scope sets.
func MakeRepoAccessList(repoScopeSets RepositoryScopeSets) []auth.Access {
	accessList := make([]auth.Access, 0, len(repoScopeSets))

	for repo, scopes := range repoScopeSets {
		for level := range scopes {
			accessList = append(accessList, auth.Access{
				Resource: auth.Resource{
					Type: "repository",
					Name: repo,
				},
				Action: level,
			})
		}
	}

	return accessList
}
