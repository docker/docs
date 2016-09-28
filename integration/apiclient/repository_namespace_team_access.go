package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

const repositoryNamespacesBasePath = "/api/v0/repositoryNamespaces"

// ListRepositoryNamespaceTeamAccess lists the teams which have been granted
// access to an organization's repository namespace along with their granted
// access levels.
func (c *apiClient) ListRepositoryNamespaceTeamAccess(namespace string) (*responses.Namespace, []responses.TeamAccess, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoryNamespacesBasePath, namespace, "teamAccess")}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, nil, err
	}

	var namespaceTeamAccessList responses.ListRepoNamespaceTeamAccess
	if err := json.NewDecoder(response.Body).Decode(&namespaceTeamAccessList); err != nil {
		return nil, nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &namespaceTeamAccessList.Namespace, namespaceTeamAccessList.TeamAccessList, nil
}

// GetRepositoryNamespaceTeamAccess gets a team's granted access to its
// organization's namespace of repositories.
func (c *apiClient) GetRepositoryNamespaceTeamAccess(namespace, teamname string) (*responses.NamespaceTeamAccess, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoryNamespacesBasePath, namespace, "teamAccess", teamname)}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var nta responses.NamespaceTeamAccess
	if err := json.NewDecoder(response.Body).Decode(&nta); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &nta, nil
}

// SetRepositoryNamespaceTeamAccess sets a team's access level to its
// organization's namespace of repositories.
func (c *apiClient) SetRepositoryNamespaceTeamAccess(namespace, teamname, accessLevel string) (*responses.NamespaceTeamAccess, error) {
	response, err := c.makeRequest(
		"PUT", url.URL{Path: path.Join(repositoryNamespacesBasePath, namespace, "teamAccess", teamname)},
		accessLevelForm{AccessLevel: accessLevel},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var nta responses.NamespaceTeamAccess
	if err := json.NewDecoder(response.Body).Decode(&nta); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &nta, nil
}

// RevokeRepositoryNamespaceTeamAccess removes a team's access to its
// organization's namespace of repositories.
func (c *apiClient) RevokeRepositoryNamespaceTeamAccess(namespace, teamname string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoryNamespacesBasePath, namespace, "teamAccess", teamname)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	return validateStatusCode(response, http.StatusNoContent)
}
