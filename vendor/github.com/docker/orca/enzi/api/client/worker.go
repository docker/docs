package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/orca/enzi/api/responses"
)

// ListWorkers lists all workers ordered by ID. If there is an API error
// response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) ListWorkers() (*responses.Workers, error) {
	endpoint := s.buildURL("/v0/workers", nil)

	var workers responses.Workers
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &workers, nil); err != nil {
		return nil, err
	}

	return &workers, nil
}

// PingWorker pings the worker with the given workerID. If there is an API
// error response then the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
// If the worker is unavailable or not responding, the API Error will have a
// code of "WORKER_UNAVAILABLE".
func (s *Session) PingWorker(workerID string) error {
	endpoint := s.buildURL(fmt.Sprintf("/v0/workers/%s/ping", workerID), nil)

	return s.performRequest("GET", endpoint, nil, http.StatusNoContent, nil, nil)
}

// DeleteWorker submits a request to delete the worker with the given ID. If
// the worker has already been shutdown or is otherwise not responsive, set
// force to true. If there is an API error response then the returned error
// will be of the type *(github.com/docker/orca/enzi/api/errors).APIErrors
// If force is false and the worker is unavailable or not responding, the API
// Error will have a code of "WORKER_UNAVAILABLE".
func (s *Session) DeleteWorker(workerID string, force bool) error {
	params := url.Values{}
	if force {
		params.Set("force", "true")
	}

	endpoint := s.buildURL(fmt.Sprintf("/v0/worker/%s", workerID), params)

	return s.performRequest("DELETE", endpoint, nil, http.StatusNoContent, nil, nil)
}
