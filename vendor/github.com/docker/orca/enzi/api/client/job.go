package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

// ListJobs lists all jobs ordered by most-recently scheduled. Jobs may be
// filtered to only a specific action by specifying an action. Jobs may also
// be filtered to those which ran on a specific worker by specifying a
// workerID. The action and worker filters may be combined. This API call only
// returns matching jobs after skipping the given offset number of jobs. This
// allows simple pagination of jobs but there may be gaps or duplicates if jobs
// are created or deleted between requests with paginated offsets. The given
// limit value specifies the maximum number of jobs to return. If zero, the
// default is 10. The limit can be no greater than 2^16. If there is an API
// error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListJobs(action, workerID string, offset, limit uint) (*responses.Jobs, error) {
	params := url.Values{}
	if action != "" {
		params.Set("action", action)
	}
	if workerID != "" {
		params.Set("worker", workerID)
	}
	if offset > 0 {
		params.Set("start", fmt.Sprintf("%d", offset))
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	endpoint := s.buildURL("/v0/jobs", params)

	var jobs responses.Jobs
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &jobs, nil); err != nil {
		return nil, err
	}

	return &jobs, nil
}

// CreateJob submits a form to run a new job now. If there is an API error
// response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) CreateJob(form forms.JobSubmission) (*responses.Job, error) {
	endpoint := s.buildURL("/v0/jobs", nil)

	var job responses.Job
	if err := s.performRequest("POST", endpoint, form, http.StatusAccepted, &job, nil); err != nil {
		return nil, err
	}

	return &job, nil
}

// DeleteJob submits a request to delete the job with the given ID. If
// there is an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) DeleteJob(jobID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/jobs/%s", jobID), nil)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetJob retrieves the job with the given jobID. If there is an API error
// response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetJob(jobID string) (*responses.Job, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/jobs/%s", jobID), nil)

	var job responses.Job
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &job, nil); err != nil {
		return nil, err
	}

	return &job, nil
}

// CancelJob submits a form to cancel the job with the given jobID. If there is
// an API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) CancelJob(jobID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/jobs/%s/cancel", jobID), nil)

	return s.performRequest("POST", endpoint, nil, http.StatusNoContent, nil, nil)
}

// GetJobLogs retrieves the logs for the job with the given jobID. It is the
// caller's responsibility to close the returned io.ReadCloser. If there is an
// API error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetJobLogs(jobID string) (io.ReadCloser, error) {
	endpoint := s.buildURL(fmt.Sprintf("/v0/jobs/%s/logs", jobID), nil)

	resp, err := s.performRequestRawResponse("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, s.handleUnexpectedResponse(http.StatusOK, resp)
	}

	return resp.Body, nil
}
