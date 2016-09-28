package crons

import (
	"net/http"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmperrors"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpforms"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/tmpresponses"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	"github.com/docker/distribution/context"
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

// NewService returns a new Crons Service.
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
		Doc("Crons")

	routes := []server.Route{
		s.routeUpdateCron(),
		s.routeListCrons(),
		s.routeGetCron(),
		s.routeDeleteCron(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// routeUpdateCron returns a route describing the UpdateCron endpoint.
func (s *Service) routeUpdateCron() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/",
		Handler: s.wrapper(s.handleUpdateCron),
		Doc:     "Create / update a periodic task",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusAccepted,
				Message: "Success.",
				Sample:  tmpresponses.Cron{},
			},
		},
		BodySample: tmpforms.CronCreate{},
	}
}

func (s *Service) handleUpdateCron(ctx context.Context, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	form := new(tmpforms.CronCreate)
	form.ActionInfo = s.actionInfo
	if formErrs := form.ValidateJSON(r.Request.Body); len(formErrs) > 0 {
		return responses.APIError(formErrs...)
	}

	cron := &schema.Cron{
		Action:      form.Action,
		Schedule:    form.Schedule,
		Retries:     form.Retries,
		Parameters:  form.Parameters,
		Deadline:    form.Deadline,
		StopTimeout: form.StopTimeout,
	}

	if _, err := s.schemaMgr.UpdateCron(cron); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusAccepted, tmpresponses.MakeCron(cron))
}

func (s *Service) routeListCrons() server.Route {
	return server.Route{
		Method:             "GET",
		Path:               "/",
		Handler:            s.wrapper(s.handleListCrons),
		Doc:                "List all crons",
		QueryParameterDocs: []*restful.Parameter{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, list of crons returned.",
				Sample:  tmpresponses.Crons{},
			},
		},
	}
}

// handleListCrons handles a request to list crons
func (s *Service) handleListCrons(ctx context.Context, r *restful.Request) responses.APIResponse {
	crons, err := s.schemaMgr.ListCrons()

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeCrons(crons))
}

func (s *Service) routeGetCron() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{action}",
		Handler: s.wrapper(s.handleGetCron),
		Doc:     "Get info about the cron with the given action",
		PathParameterDocs: map[string]string{
			"action": "action of the cron to fetch",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, cron info returned.",
				Sample:  tmpresponses.Cron{},
			},
		},
	}
}

// handleGetCron handles a request to get details for a specific cron.
func (s *Service) handleGetCron(ctx context.Context, r *restful.Request) responses.APIResponse {
	action := r.PathParameters()["action"]

	cron, err := s.schemaMgr.GetCron(action)
	if err != nil {
		if err == schema.ErrNoSuchCron {
			return responses.APIError(tmperrors.NoSuchCron(action))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeCron(cron))
}

func (s *Service) routeDeleteCron() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{cronID}",
		Handler: s.wrapper(s.handleDeleteCron),
		Doc:     "Delete the cron. Jobs created from it will not be cancelled.",
		PathParameterDocs: map[string]string{
			"cronID": "ID of cron to delete",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, cron has been deleted.",
			},
		},
	}
}

// handleDeleteCron handles a request to delete a specific cron.
func (s *Service) handleDeleteCron(ctx context.Context, r *restful.Request) responses.APIResponse {
	cronID := r.PathParameters()["cronID"]

	// Finally, delete the cron from the database.
	if err := s.schemaMgr.DeleteCron(cronID); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
