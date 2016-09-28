package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

func (c *apiClient) GetRepositoryTags(namespace, reponame string) ([]responses.Tag, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "tags")}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var parsed []responses.Tag
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return parsed, nil
}

func (c *apiClient) GetTagTrust(namespace, reponame, tag string) (*responses.Tag, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "tags", tag, "trust")}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var parsed responses.Tag
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return &parsed, nil
}

func (c *apiClient) DeleteTag(namespace, reponame, reference string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "tags", reference)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}
	return validateStatusCode(response, http.StatusNoContent)
}
