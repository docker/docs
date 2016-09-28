package worker

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/worker"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to job management.
type Service struct {
	server.Service
	worker.Worker
}

// NewService returns a new Worker Service.
func NewService(baseContext context.Context, w worker.Worker, rootPath string) *Service {
	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		Worker: w,
	}

	service.connectRoutes(rootPath)

	return service
}

// connectRoutes registers all API endpoints on this service with paths
// relative to the given rootPath.
func (s *Service) connectRoutes(rootPath string) {
	s.WebService.Path(rootPath).Doc("Worker")

	routes := []server.Route{
		s.routePing(),
		s.routeShutdown(),
		s.routeCancelJob(),
		s.routeDeleteJob(),
		s.routeGetJobLogs(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// routePing returns a route describing the ping endpoint.
func (s *Service) routePing() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/ping",
		Handler: s.handlePing,
		Doc:     "Ping the worker",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, worker is running.",
			},
		},
	}
}

// handlePing is a simple handler which just writes a No Content response.
func (s *Service) handlePing(ctx context.Context, r *restful.Request) responses.APIResponse {
	return responses.JSONResponse(http.StatusNoContent, nil)
}

// routeShutdown returns a route describing the Shutdown endpoint.
func (s *Service) routeShutdown() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/shutdown",
		Handler: s.handleShutdown,
		Doc:     "Shutdown the worker",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, worker is shutting down.",
			},
		},
	}
}

// handleShutdown shuts down the worker.
func (s *Service) handleShutdown(ctx context.Context, r *restful.Request) responses.APIResponse {
	context.GetLogger(ctx).Info("handling shutdown request ...")

	s.Worker.Shutdown()

	// Force the process to exit after 5 seconds, giving enough time to
	// write the response.
	time.AfterFunc(5*time.Second, func() {
		context.GetLogger(ctx).Info("goodbye.")

		os.Exit(0)
	})

	return responses.JSONResponse(http.StatusNoContent, nil)
}

// routeCancelJob returns a route describing the CancelJob endpoint.
func (s *Service) routeCancelJob() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/jobs/{jobID}/cancel",
		Handler: s.handleCancelJob,
		Doc:     "Cancel a job",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to cancel",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, worker canceled job.",
			},
		},
	}
}

// handleCancelJob cancels the job with the given ID.
func (s *Service) handleCancelJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]
	s.Worker.CancelJob(jobID)

	return responses.JSONResponse(http.StatusNoContent, nil)
}

// routeDeleteJob returns a route describing the DeleteJob endpoint.
func (s *Service) routeDeleteJob() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/jobs/{jobID}",
		Handler: s.handleDeleteJob,
		Doc:     "Delete a job",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to delete",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, worker deleted job.",
			},
		},
	}
}

// handleDeleteJob deletes the job with the given ID.
func (s *Service) handleDeleteJob(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]
	s.Worker.DeleteJob(jobID)

	return responses.JSONResponse(http.StatusNoContent, nil)
}

// routeGetJobLogs returns a route describing the GetJobLogs endpoint.
func (s *Service) routeGetJobLogs() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/jobs/{jobID}/logs",
		Handler: s.handleGetJobLogs,
		Doc:     "Get logs for a job",
		PathParameterDocs: map[string]string{
			"jobID": "ID of job to whose logs to retrieve",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, job logs returned.",
			},
		},
		Produces: []string{"text/plain", restful.MIME_JSON},
	}
}

func (s *Service) handleGetJobLogs(ctx context.Context, r *restful.Request) responses.APIResponse {
	jobID := r.PathParameters()["jobID"]

	logs, err := s.Worker.GetJobLogs(jobID)
	if err != nil {
		// Unable to get log file?
		return responses.APIError(errors.Internal(ctx, err))
	}

	if logs == nil {
		return responses.APIError(errors.NoSuchJob(jobID))
	}

	return &logResponse{
		ReadCloser: logs,
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
