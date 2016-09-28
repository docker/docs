package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

// These constants are used by various client methods.
const (
	RepositoryVisibilityPublic  = "public"
	RepositoryVisibilityPrivate = "private"

	repositoriesBasePath = "/api/v0/repositories"
)

// Repository contains fields from repository API responses.
// createRepositoryForm contains fields used when creating a repository.
type createRepositoryForm struct {
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	LongDescription  string `json:"longDescription"`
	Visibility       string `json:"visibility"`
}

// CreateRepository creates a new repository within an account's namespace
// with the given name, short and long description, and visibility.
func (c *apiClient) CreateRepository(namespace, reponame, shortDescription, longDescription, visibility string) (*responses.Repository, error) {
	response, err := c.makeRequest("POST", url.URL{Path: path.Join(repositoriesBasePath, namespace)}, createRepositoryForm{
		Name:             reponame,
		ShortDescription: shortDescription,
		LongDescription:  longDescription,
		Visibility:       visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusCreated); err != nil {
		return nil, err
	}

	var repo responses.Repository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &repo, nil
}

// NOTE: the shared repos endpoint is dead for now.
// // ListSharedRepositories lists all repositories you can see
// func (c *apiClient) ListSharedRepositories(username string) ([]*responses.Repository, error) {
// 	response, err := c.makeRequest("GET", path.Join(accountsBasePath, username, "sharedRepositories"), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to make API request: %s", err)
// 	}
// 	defer response.Body.Close()

// 	if err := validateStatusCode(response, http.StatusOK); err != nil {
// 		return nil, err
// 	}

// 	var repositoriesResponse struct {
// 		Repositories []*responses.Repository `json:"repositories"`
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&repositoriesResponse); err != nil {
// 		return nil, fmt.Errorf("unable to decode API response: %s", err)
// 	}

// 	return repositoriesResponse.Repositories, nil
// }

// ListAllRepositories lists all repositories you can see
func (c *apiClient) ListAllRepositories() ([]*responses.Repository, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath)}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var repositoriesResponse struct {
		Repositories []*responses.Repository `json:"repositories"`
	}
	if err := json.NewDecoder(response.Body).Decode(&repositoriesResponse); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return repositoriesResponse.Repositories, nil
}

// ListRepositories lists repositories in the given account namespace.
// Note: Only visible repositories (those which the client has at least read
// access to) will be returned.
func (c *apiClient) ListRepositories(namespace string) ([]*responses.Repository, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace)}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var repositoriesResponse struct {
		Repositories []*responses.Repository `json:"repositories"`
	}
	if err := json.NewDecoder(response.Body).Decode(&repositoriesResponse); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return repositoriesResponse.Repositories, nil
}

// GetRepository retrieves info about the specified repository.
func (c *apiClient) GetRepository(namespace, reponame string) (*responses.Repository, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame)}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var repo responses.Repository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &repo, nil
}

// RepositoryUpdateForm is used to update fields of a repository.
// If ShortDescription is not nil, the name will be updated.
// If LongDescription is not nil, the description will be updated.
// If Visibility is not nil, the visibility will be updated.
type RepositoryUpdateForm struct {
	ShortDescription *string `json:"shortDescription"`
	LongDescription  *string `json:"longDescription"`
	Visibility       *string `json:"visibility"`
}

// UpdateRepository updates a repository with the fields in the given form.
func (c *apiClient) UpdateRepository(namespace, reponame string, form RepositoryUpdateForm) (*responses.Repository, error) {
	response, err := c.makeRequest("PATCH", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame)}, form)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var repo responses.Repository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &repo, nil
}

// DeleteRepository deletes a repository.
// Note: this method may be called repeatedly with no effect.
func (c *apiClient) DeleteRepository(namespace, reponame string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}

	return validateStatusCode(response, http.StatusNoContent)
}
