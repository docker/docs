package apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpforms"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpresponses"
)

func (c *apiClient) GetJobStatus(id string) (string, error) {
	response, err := c.makeRequest("GET", url.URL{Path: "/api/v0/jobs/" + id}, nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return "", fmt.Errorf("job not found")
	}

	resp := &tmpresponses.Job{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return "", err
	}

	return resp.Status, nil
}

// RunJob starts a new job within the jobrunner framework
func (c *apiClient) RunJob(job tmpforms.JobSubmission) (*tmpresponses.Job, error) {

	response, err := c.makeRequest("POST", url.URL{Path: "/api/v0/jobs"}, job)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resp := &tmpresponses.Job{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// RunJob starts a new job within the jobrunner framework
func (c *apiClient) RunJobByAction(action string) (*tmpresponses.Job, error) {
	job := map[string]string{
		"action": action,
	}

	response, err := c.makeRequest("POST", url.URL{Path: "/api/v0/jobs"}, job)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := tmpresponses.Job{}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
