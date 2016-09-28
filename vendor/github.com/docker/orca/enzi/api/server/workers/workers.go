package workers

import (
	"net/http"
	"net/url"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to worker management.
type Service struct {
	server.Service

	schemaMgr     schema.Manager
	authenticator authn.RequestAuthenticator
	workerClient  *http.Client
}

// NewService returns a new Worker Service.
func NewService(baseContext context.Context, schemaMgr schema.Manager, workerClient *http.Client, rootPath string) *Service {
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

	return service
}

// connectRoutes registers all API endpoints on this service with paths
// relative to the given rootPath.
func (s *Service) connectRoutes(rootPath string) {
	s.WebService.Path(rootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Workers")

	routes := []server.Route{
		s.routeListWorkers(),
		s.routePingWorker(),
		s.routeDeleteWorker(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

func (s *Service) routeListWorkers() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleListWorkers),
		Doc:     "List all workers ordered by ID",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, list of workers returned.",
				Sample:  responses.Workers{},
			},
		},
	}
}

// handleListWorkers handles a request to list workers.
func (s *Service) handleListWorkers(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	workers, err := s.schemaMgr.ListWorkers()
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeWorkers(workers))
}

func (s *Service) routePingWorker() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{workerID}/ping",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handlePingWorker),
		Doc:     "Check to make sure that a job's worker is responsive",
		PathParameterDocs: map[string]string{
			"workerID": "ID of worker to ping",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, job's worker is responsive.",
			},
		},
	}
}

// handlePingWorker handles a request to ping a worker.
func (s *Service) handlePingWorker(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	workerID := r.PathParameters()["workerID"]

	worker, err := s.schemaMgr.GetWorker(workerID)
	if err != nil {
		if err == schema.ErrNoSuchWorker {
			return responses.APIError(errors.NoSuchWorker(workerID))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	pingURL := url.URL{
		Scheme: "https",
		Host:   worker.Address,
		Path:   "/v0/ping",
	}

	req, err := http.NewRequest("GET", pingURL.String(), nil)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	resp, err := s.workerClient.Do(req)
	if err != nil {
		return responses.APIError(errors.WorkerUnavailable(ctx, workerID, err))
	}

	// Don't bother checking the response. What matters is that we got one.
	resp.Body.Close()

	return responses.JSONResponse(http.StatusNoContent, nil)
}

func (s *Service) routeDeleteWorker() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{workerID}",
		Handler: server.WrapHandlerWithAdminAccount(s.authenticator, s.handleDeleteWorker),
		Doc:     "Delete a worker",
		PathParameterDocs: map[string]string{
			"workerID": "ID of worker to delete",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("force", "Force removal of the worker record. For use when the worker has been manually shutdown.").DataType("boolean").DefaultValue("false"),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, worker has been deleted.",
			},
		},
	}
}

// handleDeleteWorker handles a request to delete a specific worker.
func (s *Service) handleDeleteWorker(ctx context.Context, adminAccount *authn.Account, r *restful.Request) responses.APIResponse {
	workerID := r.PathParameters()["workerID"]
	force := helpers.ParseBoolQueryParam(r, "force", false)

	worker, err := s.schemaMgr.GetWorker(workerID)
	if err != nil {
		if err == schema.ErrNoSuchWorker {
			// The worker was already deleted.
			return responses.JSONResponse(http.StatusNoContent, nil)
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	shutdownURL := url.URL{
		Scheme: "https",
		Host:   worker.Address,
		Path:   "/v0/shutdown",
	}

	req, err := http.NewRequest("POST", shutdownURL.String(), nil)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	resp, err := s.workerClient.Do(req)
	if err != nil && !force {
		return responses.APIError(errors.WorkerUnavailable(ctx, workerID, err))
	}

	if resp != nil {
		// Don't bother checking the response. What matters is that we got one.
		resp.Body.Close()
	}

	if err := s.schemaMgr.DeleteWorker(workerID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
