package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

// Different repository access levels.
const (
	AccessLevelReadOnly  = "read-only"
	AccessLevelReadWrite = "read-write"
	AccessLevelAdmin     = "admin"
)

// UserAccess describes a user account's access level to some resource.
type RepositoryUserAccess struct {
	AccessLevel string  `json:"accessLevel"`
	User        Account `json:"user"`
}

// GetUserRepositoryAccess gets the access level that a user has to a specific
// repository. Account ownership, user access grants, and team access grants
// are all considered in determining the user's access level to the repository.
// Note: A user can only lookup repository access for themself.
// Note: Superuser status is not considered.
func (c *apiClient) GetUserRepositoryAccess(username, namespace, reponame string) (*responses.RepoUserAccess, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join("/api"+accountsBasePath, username, "repositoryAccess", namespace, reponame)}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var rua responses.RepoUserAccess
	if err := json.NewDecoder(response.Body).Decode(&rua); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &rua, nil
}

// // GetUserRepositoryNamespaceAccess gets the access level that a user has to a namespace.
// func (c *apiClient) GetUserRepositoryNamespaceAccess(username, namespace string) (*responses.RepoNamespaceUserAccess, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, username, "repositoryNamespaceAccess", namespace), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var rua responses.RepoNamespaceUserAccess
// 	if err := json.NewDecoder(response.Body).Decode(&rua); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return &rua, nil
// }

// // ListRepositoryUserAccess lists the users which have been granted access to
// // a user-owned repository along with their granted access levels.
// func (c *apiClient) ListRepositoryUserAccess(namespace, reponame string) (*Repository, []*UserAccess, error) {
// 	response, err := c.makeRequest("GET", path.Join(repositoriesBasePath, namespace, reponame, "userAccess"), nil)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, nil, err
// 	}

// 	var repositoryUserAccessList struct {
// 		Repository     *Repository   `json:"repository"`
// 		UserAccessList []*UserAccess `json:"userAccessList"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&repositoryUserAccessList); err != nil {
// 		return nil, nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return repositoryUserAccessList.Repository, repositoryUserAccessList.UserAccessList, nil
// }

type accessLevelForm struct {
	AccessLevel string `json:"accessLevel"`
}

// SetRepositoryUserAccess sets a user's access level on a user-owned
// repository.
func (c *apiClient) SetRepositoryUserAccess(namespace, reponame, username, accessLevel string) (*RepositoryUserAccess, error) {
	response, err := c.makeRequest(
		"PUT", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "userAccess", username)},
		accessLevelForm{AccessLevel: accessLevel},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var rua RepositoryUserAccess
	if err := json.NewDecoder(response.Body).Decode(&rua); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &rua, nil
}

// RevokeRepositoryUserAccess removes access to a user-owned repository for a
// specific user.
func (c *apiClient) RevokeRepositoryUserAccess(namespace, reponame, username string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "userAccess", username)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	return validateStatusCode(response, http.StatusNoContent)
}
