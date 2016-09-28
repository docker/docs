package jobs

import (
	"net/http"
	"time"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/helpers"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpforms"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpresponses"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to job management.
type Service struct {
	server.Service

	actionInfo map[string]worker.ActionInfo
	schemaMgr  schema.JobrunnerManager
	wrapper    func(server.Handler) server.Handler
}

// NewService returns a new Jobs Service.
// TODO: the action definition should be passed into here so we can make sure no one tries to create jobs with invalid actions
func NewService(baseContext context.Context, schemaMgr schema.JobrunnerManager, actionInfo map[string]worker.ActionInfo, rootPath string, wrapper func(server.Handler) server.Handler) (*Service, error) {
	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		actionInfo: actionInfo,
		schemaMgr:  schemaMgr,
		wrapper:    wrapper,
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
		Handler: s.wrapper(s.handleCreateJob),
		Doc:     "Schedule a job to be run immediately",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusAccepted,
				Message: "Success, job waiting to be claimed.",
				Sample:  responses.Job{},
			},
		},
		BodySample: tmpforms.JobSubmission{},
	}
}

// handleCreateJob handles a request to submit a new task to run as a job on a
// worker.
func (s *Service) handleCreateJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	form := new(tmpforms.JobSubmission)
	form.ActionInfo = s.actionInfo
	if formErrs := form.ValidateJSON(r.Request.Body); len(formErrs) > 0 {
		return responses.APIError(formErrs...)
	}

	now := time.Now().UTC()

	job := &schema.Job{
		Status:      schema.JobStatusWaiting,
		ScheduledAt: now,
		LastUpdated: now,
		Action:      form.Action,
		Parameters:  form.Parameters,
		// RetriesLeft: form.RetriesLeft,
		Deadline:    form.Deadline,
		StopTimeout: form.StopTimeout,
	}

	// CreateJob will generate and set an ID for the job.
	if err := s.schemaMgr.CreateJob(job); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusAccepted, tmpresponses.MakeJob(job))
}

func (s *Service) routeListJobs() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/",
		Handler: s.wrapper(s.handleListJobs),
		Doc:     "List all jobs ordered by most recently scheduled",
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("action", "Filter jobs by action.").DefaultValue("any"),
			restful.QueryParameter("worker", "Filter jobs by worker ID.").DefaultValue("any"),
			// TODO: implement list of running jobs
			restful.QueryParameter("running", "Show only jobs that are running.").DefaultValue("any"),
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
func (s *Service) handleListJobs(ctx context.Context, r *restful.Request) responses.APIResponse {
	var (
		jobs             []schema.Job
		err              error
		offsetStr, limit = helpers.PageParams(r, "start", "limit")
		offset           = helpers.ParseOffsetAsUint(offsetStr)
	)

	action, workerID := r.QueryParameter("action"), r.QueryParameter("worker")

	switch {
	case action == "any" && workerID == "any":
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobs(offset, limit)
	case action == "any":
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobsForWorker(workerID, offset, limit)
	case workerID == "any":
		jobs, err = s.schemaMgr.GetMostRecentlyScheduledJobsWithAction(action, offset, limit)
	default:
		jobs, err = s.schemaMgr.GetLeastRecentlyScheduledJobsForWorkerWithAction(workerID, action, offset, limit)
	}

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeJobs(jobs))
}

func (s *Service) routeGetJob() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{jobID}",
		Handler: s.wrapper(s.handleGetJob),
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
func (s *Service) handleGetJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	job, err := s.schemaMgr.GetJob(jobID)
	if err != nil {
		if err == schema.ErrNoSuchJob {
			return responses.APIError(errors.NoSuchJob(jobID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeJob(job))
}

func (s *Service) routeCancelJob() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/{jobID}/cancel",
		Handler: s.wrapper(s.handleCancelJob),
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
func (s *Service) handleCancelJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	if err := s.schemaMgr.CancelJob(jobID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}

func (s *Service) routeDeleteJob() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{jobID}",
		Handler: s.wrapper(s.handleDeleteJob),
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
// TODO: consider canceling the job if we can
func (s *Service) handleDeleteJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	// TODO: do an atomic deletion only if the state is "running" or unclaimed
	// if not, delete the logs and the job all from here right now

	// delete the job from the database.
	if err := s.schemaMgr.DeleteJob(jobID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}

func (s *Service) routeGetJobLogs() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{jobID}/logs",
		Handler: s.wrapper(s.handleGetJobLogs),
		Doc:     "Retrieve logs for this job from its worker",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job whose logs to retrieve",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, job's logs returned.",
				Sample:  tmpresponses.JobLogs{},
			},
		},
	}
}

// handleGetJobLogs handles a request to get logs of a specific job.
func (s *Service) handleGetJobLogs(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]
	// TODO: wire up pagination
	jobLogs, err := s.schemaMgr.GetJobLogs(jobID, 0, 0)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	// convert data into the right response type
	resp := tmpresponses.JobLogs{}
	for _, log := range jobLogs {
		resp.JobLogs = append(resp.JobLogs, tmpresponses.JobLog{
			Data: log.Data,
		})
	}

	return responses.JSONResponse(http.StatusOK, jobLogs)
}
