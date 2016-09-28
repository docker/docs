package action_configs

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
		Doc("ActionConfigs")

	routes := []server.Route{
		s.routeUpdateActionConfig(),
		s.routeListActionConfigs(),
		s.routeGetActionConfig(),
		s.routeDeleteActionConfig(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// routeUpdateActionConfig returns a route describing the UpdateActionConfig endpoint.
func (s *Service) routeUpdateActionConfig() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/",
		Handler: s.wrapper(s.handleUpdateActionConfig),
		Doc:     "Configure actions",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusAccepted,
				Message: "Success.",
				Sample:  tmpresponses.ActionConfig{},
			},
		},
		BodySample: tmpforms.ActionConfigCreate{},
	}
}

func (s *Service) handleUpdateActionConfig(ctx context.Context, r *restful.Request) responses.APIResponse {
	defer r.Request.Body.Close()

	form := new(tmpforms.ActionConfigCreate)
	form.ActionInfo = s.actionInfo
	if formErrs := form.ValidateJSON(r.Request.Body); len(formErrs) > 0 {
		return responses.APIError(formErrs...)
	}

	actionConfig := &schema.ActionConfig{
		Action:           form.Action,
		MaxJobsPerWorker: form.MaxJobsPerWorker,
		HeartbeatTimeout: form.HeartbeatTimeout,
		Parameters:       form.Parameters,
	}

	if _, err := s.schemaMgr.UpdateActionConfig(actionConfig); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusAccepted, tmpresponses.MakeActionConfig(actionConfig))
}

func (s *Service) routeListActionConfigs() server.Route {
	return server.Route{
		Method:             "GET",
		Path:               "/",
		Handler:            s.wrapper(s.handleListActionConfigs),
		Doc:                "List all action configs",
		QueryParameterDocs: []*restful.Parameter{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, list of action configs returned.",
				Sample:  tmpresponses.ActionConfigs{},
			},
		},
	}
}

// handleListActionConfigs handles a request to list actionConfigs
func (s *Service) handleListActionConfigs(ctx context.Context, r *restful.Request) responses.APIResponse {
	actionConfigs, err := s.schemaMgr.ListActionConfigs()

	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeActionConfigs(actionConfigs))
}

func (s *Service) routeGetActionConfig() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{action}",
		Handler: s.wrapper(s.handleGetActionConfig),
		Doc:     "Get info about the actionConfig with the given action",
		PathParameterDocs: map[string]string{
			"action": "name of action to fetch the config for",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, action config info returned.",
				Sample:  tmpresponses.ActionConfig{},
			},
		},
	}
}

// handleGetActionConfig handles a request to get details for a specific actionConfig.
func (s *Service) handleGetActionConfig(ctx context.Context, r *restful.Request) responses.APIResponse {
	action := r.PathParameters()["action"]

	actionConfig, err := s.schemaMgr.GetActionConfig(action)
	if err != nil {
		if err == schema.ErrNoSuchActionConfig {
			return responses.APIError(tmperrors.NoSuchActionConfig(action))
		}

		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, tmpresponses.MakeActionConfig(actionConfig))
}

func (s *Service) routeDeleteActionConfig() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{action}",
		Handler: s.wrapper(s.handleDeleteActionConfig),
		Doc:     "Delete the action config. The defaults will be used.",
		PathParameterDocs: map[string]string{
			"action": "the name of the action to delete the config for",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, action config has been deleted.",
			},
		},
	}
}

// handleDeleteActionConfig handles a request to delete a specific actionConfig.
func (s *Service) handleDeleteActionConfig(ctx context.Context, r *restful.Request) responses.APIResponse {
	action := r.PathParameters()["action"]

	// Finally, delete the actionConfig from the database.
	if err := s.schemaMgr.DeleteActionConfig(action); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
