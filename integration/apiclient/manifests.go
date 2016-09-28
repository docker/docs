package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

func (c *apiClient) GetRepositoryManifests(namespace, reponame string) ([]responses.Manifest, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "manifests")}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to make API request: %s", err)
	}
	defer response.Body.Close()

	if err := validateStatusCode(response, http.StatusOK); err != nil {
		return nil, err
	}

	var parsed []responses.Manifest
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("unable to decode API response: %s", err)
	}

	return parsed, nil
}

func (c *apiClient) DeleteManifest(namespace, reponame, reference string) error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(repositoriesBasePath, namespace, reponame, "manifests", reference)}, nil)
	if err != nil {
		return fmt.Errorf("unable to make API request: %s", err)
	}
	return validateStatusCode(response, http.StatusNoContent)
}
