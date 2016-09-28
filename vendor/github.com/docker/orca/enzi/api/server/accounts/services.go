package accounts

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/filters"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// routeCreateService returns a route describing the CreateService endpoint.
func (s *Service) routeCreateService() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/{accountNameOrID}/services",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleCreateService),
		Doc:     "Create a Service",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or ID of account which will own the service",
		},
		BodySample: forms.CreateService{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusCreated,
				Message: "Success, service created.",
				Sample:  responses.Service{},
			},
		},
	}
}

// handleCreateService handles a request for creating a service.
func (s *Service) handleCreateService(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAccountAccess,
		rd.RequireAccountAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.CreateService)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	if form.Privileged && !clientAccount.IsAdmin {
		return responses.APIError(errors.NotAuthorized("must be a system admin to make a privileged service"))
	}

	newService := &schema.Service{
		OwnerID:            rd.Acct.ID,
		Name:               form.Name,
		Description:        form.Description,
		URL:                form.URL,
		Privileged:         form.Privileged,
		RedirectURIs:       form.RedirectURIs,
		GrantTypes:         []string{"authorization_code", "refresh_token", "service_session"},
		ResponseTypes:      []string{"code"},
		JWKsURIs:           form.JWKsURIs,
		ProviderIdentities: form.ProviderIdentities,
		CABundle:           form.CABundle,
	}

	if newService.Privileged {
		newService.GrantTypes = append(newService.GrantTypes, "password", "root_session")
	}

	if err := s.schemaMgr.CreateService(newService); err != nil {
		if err == schema.ErrServiceExists {
			return responses.APIError(errors.ServiceExists())
		}
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusCreated, responses.MakeService(newService))
}

// routeListServices returns a route describing the ListServices endpoint.
func (s *Service) routeListServices() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{accountNameOrID}/services",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleListServices),
		Doc:     "List services owned by an account",
		Notes:   "Lists services in ascending order by name.",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account whose owned services will be listed",
		},
		QueryParameterDocs: []*restful.Parameter{
			restful.QueryParameter("start", "Only return services with a name greater than or equal to this name."),
			restful.QueryParameter("limit", "Maximum number of services per page of results.").DataType("int").DefaultValue(fmt.Sprint(api.DefaultPerPageLimit)),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, page of services listed.",
				Sample:  responses.Services{},
			},
		},
	}
}

// handleListServices handles a request for listing services owned by an
// account.
func (s *Service) handleListServices(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAccountAccess,
		rd.RequireAccountAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	startName, limit := helpers.PageParams(r, "start", "limit")

	services, nextPageStart, err := s.schemaMgr.ListServicesForAccount(rd.Acct.ID, startName, limit)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, responses.MakeServices(services), r, nextPageStart)
}

// routeGetService returns a route describing the GetService endpoint.
func (s *Service) routeGetService() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/{accountNameOrID}/services/{serviceNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleGetService),
		Doc:     "Details for a service",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account which owns the service",
			"serviceNameOrID": "Name or id of service which will be retrieved",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, service returned.",
				Sample:  responses.Service{},
			},
		},
	}
}

// handleGetService handles a request for getting details of a service owned by
// an account.
func (s *Service) handleGetService(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]
	serviceNameOrID := pathParams["serviceNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAccountAccess,
		rd.RequireAccountAdmin,
		rd.MakeFilterGetService(serviceNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeService(rd.Service))
}

// routeUpdateService returns a route describing the UpdateServiec endpoint.
func (s *Service) routeUpdateService() server.Route {
	return server.Route{
		Method:  "PATCH",
		Path:    "/{accountNameOrID}/services/{serviceNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleUpdateService),
		Doc:     "Update details for a service",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of account which owns the service",
			"serviceNameOrID": "Name or id of service which will be updated",
		},
		BodySample: forms.UpdateService{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, service updated.",
				Sample:  responses.Service{},
			},
		},
	}
}

// handleUpdateService handles a request for updating details of a service
// owned by an account.
func (s *Service) handleUpdateService(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]
	serviceNameOrID := pathParams["serviceNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAccountAccess,
		rd.RequireAccountAdmin,
		rd.MakeFilterGetService(serviceNameOrID),
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.UpdateService)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	if form.Privileged != nil && !clientAccount.IsAdmin {
		return responses.APIError(errors.NotAuthorized("must be a system admin to modify a service's privileged status"))
	}

	var updateFields schema.ServiceUpdateFields

	if form.Description != nil {
		rd.Service.Description = *form.Description
		updateFields.Description = form.Description
	}
	if form.URL != nil {
		rd.Service.URL = *form.URL
		updateFields.URL = form.URL
	}
	if form.Privileged != nil {
		rd.Service.Privileged = *form.Privileged
		updateFields.Privileged = form.Privileged
	}
	if form.RedirectURIs != nil {
		rd.Service.RedirectURIs = *form.RedirectURIs
		updateFields.RedirectURIs = form.RedirectURIs
	}
	if form.JWKsURIs != nil {
		rd.Service.JWKsURIs = *form.JWKsURIs
		updateFields.JWKsURIs = form.JWKsURIs
	}
	if form.ProviderIdentities != nil {
		rd.Service.ProviderIdentities = *form.ProviderIdentities
		updateFields.ProviderIdentities = form.ProviderIdentities
	}
	if form.CABundle != nil {
		rd.Service.CABundle = *form.CABundle
		updateFields.CABundle = form.CABundle
	}

	if err := s.schemaMgr.UpdateService(rd.Service.ID, updateFields); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeService(rd.Service))
}

// routeDeleteService returns a route describing the DeleteService endpoint.
func (s *Service) routeDeleteService() server.Route {
	return server.Route{
		Method:  "DELETE",
		Path:    "/{accountNameOrID}/services/{serviceNameOrID}",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleDeleteService),
		Doc:     "Delete a service",
		PathParameterDocs: map[string]string{
			"accountNameOrID": "Name or id of which owns the service",
			"serviceNameOrID": "Name or id of service which will be deleted",
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, service deleted.",
			},
		},
	}
}

// handleDeleteService handles a request for deleting a service in an organization.
func (s *Service) handleDeleteService(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	pathParams := r.PathParameters()
	accountNameOrID := pathParams["accountNameOrID"]
	serviceNameOrID := pathParams["serviceNameOrID"]

	// Gather request data.
	rd := filters.NewRequestData(ctx, s.schemaMgr, s.authorizer, r, clientAccount)

	if errResponse := rd.AddFilters(
		rd.MakeFilterGetAccount(accountNameOrID),
		rd.GetAccountAccess,
		rd.RequireAccountAdmin,
	).EvaluateFilters(); errResponse != nil {
		return errResponse
	}

	service, err := helpers.ServiceByNameOrID(s.schemaMgr, rd.Acct.ID, serviceNameOrID)
	if err != nil {
		if err == schema.ErrNoSuchService {
			// Already deleted.
			return responses.JSONResponse(http.StatusNoContent, nil)
		}
		return responses.APIError(errors.Internal(ctx, err))
	}

	if err := s.schemaMgr.DeleteService(service.OwnerID, service.Name); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusNoContent, nil)
}
