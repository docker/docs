package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

// ListTeamRepositoryAccess gets the access levels that a team has to any
// repositories in its organization.
func (c *apiClient) ListTeamRepositoryAccess(orgname, teamname string) (*responses.Team, []responses.RepoAccess, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join("/api"+accountsBasePath, orgname, "teams", teamname, "repositoryAccess")}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, nil, err
	}

	var teamRepositoryAccessList responses.ListTeamRepoAccess
	if err := json.NewDecoder(response.Body).Decode(&teamRepositoryAccessList); err != nil {
		return nil, nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &teamRepositoryAccessList.Team, teamRepositoryAccessList.RepositoryAccessList, nil
}

// ListRepositoryTeamAccess lists the teams which have been granted access to
// an organization-owned repository along with their granted access levels.
func (c *apiClient) ListRepositoryTeamAccess(namespace, reponame string) (*responses.Repository, []responses.TeamAccess, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "teamAccess")}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, nil, err
	}

	var repositoryTeamAccessList responses.ListRepoTeamAccess
	if err := json.NewDecoder(response.Body).Decode(&repositoryTeamAccessList); err != nil {
		return nil, nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &repositoryTeamAccessList.Repository, repositoryTeamAccessList.TeamAccessList, nil
}

// SetRepositoryTeamAccess sets a team's access level on an organization-owned
// repository.
func (c *apiClient) SetRepositoryTeamAccess(namespace, reponame, teamname, accessLevel string) (*responses.RepoTeamAccess, error) {
	response, err := c.makeRequest(
		"PUT", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "teamAccess", teamname)},
		accessLevelForm{AccessLevel: accessLevel},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var rta responses.RepoTeamAccess
	if err := json.NewDecoder(response.Body).Decode(&rta); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &rta, nil
}

// RevokeRepositoryTeamAccess removes access to an organization-owned
// repository for a specific team.
func (c *apiClient) RevokeRepositoryTeamAccess(namespace, reponame, teamname string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "teamAccess", teamname)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	return validateStatusCode(response, http.StatusNoContent)
}
