package openid

import (
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/config"
	"github.com/emicklei/go-restful"
)

// routeLogin returns the route describing the Login endpoint.
func (s *Service) routeLogin() server.Route {
	return server.Route{
		Method:     "POST",
		Path:       "/login",
		Handler:    s.handleLogin,
		Doc:        "Submit a Login Form in exchange for a Session Token",
		Consumes:   []string{restful.MIME_JSON},
		BodySample: forms.Login{},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, login response returned.",
				Sample:  responses.LoginSession{},
			},
		},
	}
}

// handleLogin handles a user login form submission.
func (s *Service) handleLogin(ctx context.Context, request *restful.Request) responses.APIResponse {
	// If the client is already authenticated with a session, we need to
	// invalidate it.
	clientAccount, apiErr := s.authorizer.AuthenticateRequest(ctx, request.Request)
	if apiErr != nil {
		return responses.APIError(apiErr)
	}

	if clientAccount.Session != nil {
		// The client already had a valid session. This shouldn't
		// happen in a controlled environment, but it might if the
		// client is doing something unexpected so we need to delete
		// the old session.
		if err := s.authorizer.SessionTokenAuthenticator().DeleteSession(clientAccount.Session); err != nil {
			return responses.APIError(errors.Internal(ctx, err))
		}
	}

	defer request.Request.Body.Close()

	form := new(forms.Login)
	if apiErrs := form.ValidateJSON(request.Request.Body); apiErrs != nil {
		return responses.APIError(apiErrs...)
	}

	authConfig, errResponse := helpers.AuthConfig(ctx, s.authorizer)
	if errResponse != nil {
		return errResponse
	}

	var ldapSettings *ldapconfig.Settings
	if authConfig.Backend == config.AuthBackendLDAP {
		ldapSettings, errResponse = helpers.LDAPSettings(ctx, s.authorizer)
		if errResponse != nil {
			return errResponse
		}
	}

	authenticator := s.authorizer.UsernamePasswordAuthenticator(ldapSettings)

	user, apiErr := authenticator.AuthenticateUsernamePassword(ctx, form.Username, form.Password)
	if apiErr != nil {
		if apiErr.IsInternal() {
			responses.APIError(apiErr)
		}

		return responses.APIError(errors.InvalidFormField("authentication", apiErr.Message))
	}

	session, err := s.authorizer.SessionTokenAuthenticator().CreateSession(user)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	return responses.JSONResponse(http.StatusOK, responses.MakeLoginSession(&user.Account, session.Secret))
}

// routeLogout returns the route describing the Logout endpoint.
func (s *Service) routeLogout() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/logout",
		Handler: s.handleLogout,
		Doc:     "Delete the current session is use.",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusNoContent,
				Message: "Success, current session deleted.",
			},
		},
	}
}

// handleLogout handles a user logout form submission.
func (s *Service) handleLogout(ctx context.Context, request *restful.Request) responses.APIResponse {
	// It's okay if the client is not authenticated to thes endpoint. It
	// becomes a nop effectively.
	clientAccount, apiErr := s.authorizer.AuthenticateRequest(ctx, request.Request)
	if apiErr != nil {
		return responses.APIError(apiErr)
	}

	if clientAccount.Session == nil {
		// The client already does not have a session. Just redirect
		// them to the next URL.
		return responses.JSONResponse(http.StatusNoContent, nil)
	}

	if err := s.authorizer.SessionTokenAuthenticator().DeleteSession(clientAccount.Session); err != nil {
		return responses.APIError(errors.Internal(ctx, err))
	}

	// IMPORTANT: set the client's session to nil so that any
	// wrapping handler does not use it.
	clientAccount.Session = nil

	return responses.JSONResponse(http.StatusNoContent, nil)
}
