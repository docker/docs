package jobs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to job management.
type Service struct {
	server.Service

	schemaMgr     schema.Manager
	authenticator authn.RequestAuthenticator
	workerClient  *http.Client
}

// NewService returns a new Jobs Service.
func NewService(baseContext context.Context, schemaMgr schema.Manager, workerClient *http.Client, rootPath string) (*Service, error) {
	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		schemaMgr:     schemaMgr,
		authenticator: authz.NewAuthorizer(schemaMgr),
		workerClient:  workerClient,
	}

	service.connectRoutes(rootPath)

	return service, nil
}

// connectRoutes registers all API endpoints on this service with paths
// relative to the given rootPath.
func (s *Service) connectRoutes(rootPath string) {
	s.WebService.Path(rootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Jobs")

	routes := []server.Route{
		s.routeCreateJob(),
		s.routeListJobs(),
		s.routeGetJob(),
		s.routeCancelJob(),
		s.routeDeleteJob(),
		s.routeGetJobLogs(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// routeCreateJob returns a route describing the CreateJob endpoint.
func (s *Service) routeCreateJob() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleCreateJob),
		Doc:     "Schedule a job to be run immediately",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusAccepted,
				Message: "Success, job waiting to be claimed.",
				Sample:  responses.Job{},
			},
		},
		BodySample: forms.JobSubmission{},
	}
}

// handleCreateJob handles a request to submit a new task to run as a job on a
// worker.
func (s *Service) handleCreateJob(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	form := new(forms.JobSubmission)
	if formErrs := form.ValidateJSON(r.Request.Body); len(formErrs) > 0 {
		return responses.APIError(formErrs...)
	}

	now := time.Now().UTC()

	job := &schema.Job{
		Status:      schema.JobStatusWaiting,
		ScheduledAt: now,
		LastUpdated: now,
		Action:      form.Action,
	}

	// CreateJob will generate and set an ID for the job.
	if err := s.schemaMgr.CreateJob(job); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusAccepted, responses.MakeJob(job))
}

func (s *Service) routeListJobs() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleListJobs),
		Doc:     "List all jobs ordered by most recently scheduled",
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("action", "Filter jobs by action.").DefaultValue("any"),
			restful.QueryParameter("worker", "Filter jobs by worker ID.").DefaultValue("any"),
			restful.QueryParameter("start", "Return most recently scheduled jobs starting from this offset index.").DataType("int").DefaultValue("0"),
			restful.QueryParameter("limit", "Maximum number of jobs per page of results.").DataType("int").DefaultValue(string(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, list of jobs returned.",
				Sample:  responses.Jobs{},
			},
		},
	}
}

// handleListJobs handles a request to list jobs
func (s *Service) handleListJobs(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	var (
		jobs             []schema.Job
		err              error
		offsetStr, limit = helpers.PageParams(r, "start", "limit")
		offset           = helpers.ParseOffsetAsUint(offsetStr)
	)

	action, workerID := r.QueryParameter("action"), r.QueryParameter("worker")

	any := "any"
	if action == "" {
		action = any
	}
	if workerID == "" {
		workerID = any
	}

	switch {
	case action == any && workerID == any:
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobs(offset, limit)
	case action == any:
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobsForWorker(workerID, offset, limit)
	case workerID == any:
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobsWithAction(action, offset, limit)
	default:
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobsForWorkerWithAction(workerID, action, offset, limit)
	}

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeJobs(jobs))
}

func (s *Service) routeGetJob() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{jobID}",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleGetJob),
		Doc:     "Get info about the job with the given ID",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to fetch",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, job info returned.",
				Sample:  responses.Job{},
			},
		},
	}
}

// handleGetJob handles a request to get details for a specific job.
func (s *Service) handleGetJob(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	job, err := s.schemaMgr.GetJob(jobID)
	if err != nil {
		if err == schema.ErrNoSuchJob {
			return responses.APIError(errors.NoSuchJob(jobID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeJob(job))
}

func (s *Service) routeCancelJob() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/{jobID}/cancel",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleCancelJob),
		Doc:     "Signal this job's worker to cancel the job",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to cancel",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, job has been canceled.",
			},
		},
	}
}

// handleCancelJob handles a request to cancel a specific job.
func (s *Service) handleCancelJob(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	job, err := s.schemaMgr.GetJob(jobID)
	if err != nil {
		if err == schema.ErrNoSuchJob {
			return responses.APIError(errors.NoSuchJob(jobID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	worker, err := s.schemaMgr.GetWorker(job.WorkerID)
	if err != nil {
		if err == schema.ErrNoSuchWorker {
			return responses.APIError(errors.NoSuchWorker(job.WorkerID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	cancelURL := url.URL{
		Scheme: "https",
		Host:   worker.Address,
		Path:   fmt.Sprintf("/v0/jobs/%s/cancel", job.ID),
	}

	req, err := http.NewRequest("POST", cancelURL.String(), nil)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	resp, err := s.workerClient.Do(req)
	if err != nil {
		return responses.APIError(errors.WorkerUnavailable(ctx, job.WorkerID, err))
	}

	// Don't bother checking the response. What matters is that we got one.
	resp.Body.Close()

	return responses.JSONResponse(http.StatusNoContent, nil)
}

func (s *Service) routeDeleteJob() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{jobID}",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleDeleteJob),
		Doc:     "Signal this job's worker to cancel and delete the job",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to delete",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, job has been deleted.",
			},
		},
	}
}

// handleDeleteJob handles a request to delete a specific job.
func (s *Service) handleDeleteJob(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	job, err := s.schemaMgr.GetJob(jobID)
	if err != nil {
		if err == schema.ErrNoSuchJob {
			// The job must have already been deleted.
			return responses.JSONResponse(http.StatusNoContent, nil)
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	worker, err := s.schemaMgr.GetWorker(job.WorkerID)
	if err != nil {
		if err == schema.ErrNoSuchWorker {
			return responses.APIError(errors.NoSuchWorker(job.WorkerID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	deleteURL := url.URL{
		Scheme: "https",
		Host:   worker.Address,
		Path:   fmt.Sprintf("/v0/jobs/%s", job.ID),
	}

	req, err := http.NewRequest("DELETE", deleteURL.String(), nil)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	resp, err := s.workerClient.Do(req)
	if err != nil {
		return responses.APIError(errors.WorkerUnavailable(ctx, job.WorkerID, err))
	}

	// Don't bother checking the response. What matters is that we got one.
	resp.Body.Close()

	// Finally, delete the job from the database.
	if err := s.schemaMgr.DeleteJob(jobID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}

func (s *Service) routeGetJobLogs() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{jobID}/logs",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleGetJobLogs),
		Doc:     "Retrieve logs for this job from its worker",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job whose logs to retrieve",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, job's logs returned.",
			},
		},
		Produces: []string{"text/plain", restful.MIME_JSON},
	}
}

// handleGetJobLogs handles a request to get logs of a specific job.
func (s *Service) handleGetJobLogs(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	job, err := s.schemaMgr.GetJob(jobID)
	if err != nil {
		if err == schema.ErrNoSuchJob {
			return responses.APIError(errors.NoSuchJob(jobID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	worker, err := s.schemaMgr.GetWorker(job.WorkerID)
	if err != nil {
		if err == schema.ErrNoSuchWorker {
			return responses.APIError(errors.NoSuchWorker(job.WorkerID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	logsURL := url.URL{
		Scheme: "https",
		Host:   worker.Address,
		Path:   fmt.Sprintf("/v0/jobs/%s/logs", job.ID),
	}

	req, err := http.NewRequest("GET", logsURL.String(), nil)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	resp, err := s.workerClient.Do(req)
	if err != nil {
		return responses.APIError(errors.WorkerUnavailable(ctx, job.WorkerID, err))
	}

	// The status code will be 200 on success.
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)

		return responses.APIError(errors.Internal(ctx, fmt.Errorf("unexpected logs response code from worker: %s output: %s", resp.StatusCode, buf.String())))
	}

	return &logResponse{
		ReadCloser: resp.Body,
	}
}

// logResponse is used as an HTTP response for job logs.
type logResponse struct {
	io.ReadCloser
}

func (r *logResponse) StatusCode() int {
	return http.StatusOK
}

func (r *logResponse) WriteResponse(ctx context.Context, response *restful.Response) {
	defer r.Close()

	if _, err := io.Copy(response, r); err != nil {
		context.GetLogger(ctx).Errorf("unable to copy job logs to response: %s", err)
	}
}

func (r *logResponse) AddCookies(cookies ...*http.Cookie) {
	// Not supported.
}
