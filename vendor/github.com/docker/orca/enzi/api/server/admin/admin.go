package admin

import (
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authz"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

// Service handles various API requests relating to admin activities. Currently,
// the only "admin activity" is the batch import accounts endpoint.
type Service struct {
	server.Service

	schemaMgr  schema.Manager
	authorizer authz.Authorizer
}

// NewService returns a new Worker Service.
func NewService(baseContext context.Context, schemaMgr schema.Manager, rootPath string) *Service {
	service := &Service{
		Service: server.Service{
			WebService:  new(restful.WebService),
			BaseContext: baseContext,
		},
		schemaMgr:  schemaMgr,
		authorizer: authz.NewAuthorizer(schemaMgr),
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
		Doc("Admin")

	routes := []server.Route{
		s.routeImportAccounts(),
	}

	for _, route := range routes {
		route.Register(&s.Service)
	}
}

// RouteImportAccounts returns a route describing the ImportAccounts endpoint.
func (s *Service) routeImportAccounts() server.Route {
	return server.Route{
		Method:     "POST",
		Path:       "/importAccounts",
		Handler:    server.WrapHandlerWithAdminAccount(s.authorizer, s.handleImportAccounts),
		Doc:        "Admin-only endpoint to import many users and organizations. This endpoint allows adding managed users with password hashes and is used for our migrations processes.",
		BodySample: forms.ImportAccounts{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusCreated,
				Message: "Success, accounts created.",
				Sample:  nil,
			},
		},
	}
}

// handleImportAccounts handles a request by an admin for importing many
// users or organizations. This endpoint allows creation of users with
// PasswordHashes instead of raw passwords, which enables migrations.
func (s *Service) handleImportAccounts(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	// Gather request data.
	authConfig, errResponse := helpers.AuthConfig(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	defer r.Request.Body.Close()

	// Decode and validate the form.
	form := new(forms.ImportAccounts)
	if formErrors := form.ValidateJSON(r.Request.Body); len(formErrors) > 0 {
		return responses.APIError(formErrors...)
	}

	var orgAccounts []forms.ImportAccount
	var userAccounts []forms.ImportAccount
	for _, acc := range form.Accounts {
		if acc.IsOrg {
			orgAccounts = append(orgAccounts, acc)
		} else {
			userAccounts = append(userAccounts, acc)
		}
	}
	if errs := s.importAccounts(ctx, orgAccounts); len(errs) > 0 {
		return responses.APIError(errs...)
	}

	if len(userAccounts) > 0 {
		if authConfig.Backend == config.AuthBackendLDAP {
			// Can't create users if using LDAP, they must be synced.
			return responses.APIError(errors.CannotCreateUser("Users are synced with LDAP"))
		} else if errs := s.importAccounts(ctx, userAccounts); len(errs) > 0 {
			return responses.APIError(errs...)
		}
	}

	return responses.JSONResponse(http.StatusCreated, nil)
}

func (s *Service) importAccounts(ctx context.Context, forms []forms.ImportAccount) (errs []*errors.APIError) {
	var dupAccNames []string
	var internalErr error

	for _, newAccForm := range forms {
		newAccount := &schema.Account{
			Name:         newAccForm.Name,
			IsOrg:        newAccForm.IsOrg,
			FullName:     newAccForm.FullName,
			IsActive:     newAccForm.IsActive,
			IsAdmin:      newAccForm.IsAdmin,
			PasswordHash: newAccForm.PasswordHash,
		}

		if err := s.schemaMgr.CreateAccount(newAccount); err != nil {
			if err == schema.ErrAccountExists {
				dupAccNames = append(dupAccNames, newAccForm.Name)
			} else {
				internalErr = err
				break
			}
		}
	}

	if internalErr != nil {
		errs = append(errs, errors.Internal(ctx, internalErr))
	}
	if len(dupAccNames) > 0 {
		errs = append(errs, errors.AccountsExist(dupAccNames))
	}
	return errs
}
