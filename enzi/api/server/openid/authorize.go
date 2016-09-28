package openid

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

func (s *Service) routeAuthorize() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/authorize",
		Handler: server.WrapHandlerWithAuthentication(s.authorizer, s.handleAuthorize),
		Doc:     "Oauth 2.0 Authorization",
		Notes: `This endpoint performs authentication of the End-User
                        on behalf of a service. The service should prepare a
                        request and redirect the End-User to this endpoint
                        using request parameters defined by OAuth 2.0 and
                        additional parameters and parameter values defined by
                        OpenID Connect.
                        ` + "\n" +
			`This endpoint will perform validation of the request,
                        optionally prompt the user for consent/authorization,
                        and redirect the End-User to a registered callback URL
                        for the service with an Authorization Code parameter
                        which the service will then be able to exchange it for
                        an ID Token and Access Token directly using the Token
                        endpoint on this service.`,
		Consumes: []string{"application/x-www-form-urlencoded"},
		FormParameterDocs: []*restful.Parameter{
			restful.FormParameter("scope", "MUST contain the value 'openid'. No other scope values are understood at this time."),
			restful.FormParameter("response_type", "MUST contain the value 'code'. No other response_type values are understood at this time."),
			restful.FormParameter("client_id", "MUST contain the ID of the service which is requesting authorization by the End-User."),
			restful.FormParameter("redirect_uri", "Redirection URI to which the End-User will be sent with an authorization code. This value MUST exactly match one of the redirect URI values that the service has pre-registered with this auth provider."),
			restful.FormParameter("state", "RECOMMENDED to be set to a value that is used to maintain session state between this request and the callback. Typically used to mitigate CSRF, this value may be set to a hash of a browser cookie set by the service. This value may also be used by the service to know which page to redirect to next, after the authentication flow is complete. The redirect response will include this value as a request parameter."),
		},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, Redirect Response Returned.",
				Sample:  redirectResponse{},
			},
		},
	}
}

type redirectResponse struct {
	Redirect string `json:"redirect" description:"The URL to redirect to the service to continue the authorization flow."`
}

func (s *Service) handleAuthorize(ctx context.Context, clientAccount *authn.Account, r *restful.Request) responses.APIResponse {
	// Parse the form, ignore errors.
	r.Request.ParseForm()
	params := r.Request.PostForm

	// Validate scope, response type, client_id, and redirect_uri.
	//
	// If the request fails due to a missing, invalid, or mismatching
	// redirection URI, or if the client identifier is missing or invalid,
	// the authorization server SHOULD inform the resource owner of the
	// error and MUST NOT automatically redirect the user-agent to the
	// invalid redirection URI.
	serviceID := params.Get("client_id")
	if serviceID == "" {
		return responses.APIError(errors.InvalidFormField("client_id", "this value MUST be specified"))
	}

	service, err := s.schemaMgr.GetServiceByID(serviceID)
	if err != nil {
		if err == schema.ErrNoSuchService {
			return responses.APIError(errors.InvalidFormField("client_id", err.Error()))
		}

		// We would redirect to the service if we could be we aren't
		// even able to validate the redirect URI yet.
		return responses.APIError(errors.Internal(ctx, fmt.Errorf("unable to get service: %s", err)))
	}

	redirectURI := params.Get("redirect_uri")
	if redirectURI == "" {
		return responses.APIError(errors.InvalidFormField("redirect_uri", "this value MUST be specified"))
	}

	if !validateRedirectURI(service, redirectURI) {
		return responses.APIError(errors.InvalidFormField("redirect_uri", "this value MUST exactly match one of the services registered Redirect URIs"))
	}

	// If the resource owner denies the access request or if the request
	// fails for reasons other than a missing or invalid redirection URI,
	// the authorization server informs the client by adding error
	// parameters to the query component of the redirection URI.
	if scope := params.Get("scope"); scope != "openid" {
		return redirectErrorResponse(redirectURI, params, "invalid_scope", "this value MUST be 'openid'")
	}

	if responseType := params.Get("response_type"); responseType != "code" {
		return redirectErrorResponse(redirectURI, params, "unsupported_response_type", "this value MUST be 'code'")
	}

	// Validate that the service is privileged. Non-privileged services are
	// not yet supported.
	if !service.Privileged {
		return redirectErrorResponse(redirectURI, params, "access_denied", "only privileged services are currently supported")
	}

	// Create service authorization code.
	authCode := &schema.ServiceAuthCode{
		Expiration:  time.Now().Add(time.Minute),
		ServiceID:   service.ID,
		AccountID:   clientAccount.ID,
		RedirectURI: redirectURI,
	}

	if clientAccount.Session != nil {
		// This authorization code can be used to establish a linked
		// service session.
		authCode.SessionID = clientAccount.Session.ID
	}

	if err := s.schemaMgr.CreateServiceAuthCode(authCode); err != nil {
		// We must send the internal server error message to the
		// service. Log it here first though.
		context.GetLogger(ctx).Errorf("unable to create service auth code: %s", err)
		return redirectErrorResponse(redirectURI, params, "server_error", fmt.Sprintf("internal error - requestID: %s", ctx.Value("http.request.id")))
	}

	// Redirect the redirect URI with code and state parameters.
	redirectParams := url.Values{}
	redirectParams.Set("code", authCode.Code)
	if state := params.Get("state"); state != "" {
		redirectParams.Set("state", state)
	}

	redirect := redirectResponse{
		Redirect: fmt.Sprintf("%s?%s", redirectURI, redirectParams.Encode()),
	}

	return responses.JSONResponse(http.StatusOK, redirect)
}

func validateRedirectURI(service *schema.Service, givenURI string) bool {
	for _, redirectURI := range service.RedirectURIs {
		if givenURI == redirectURI {
			return true
		}
	}

	return false
}

func redirectErrorResponse(redirectURI string, params url.Values, errCode, errDescription string) responses.APIResponse {
	redirectParams := url.Values{}
	redirectParams.Set("error", errCode)
	redirectParams.Set("error_description", errDescription)
	if state := params.Get("state"); state != "" {
		redirectParams.Set("state", state)
	}

	redirect := redirectResponse{
		Redirect: fmt.Sprintf("%s?%s", redirectURI, redirectParams.Encode()),
	}

	return responses.JSONResponse(http.StatusOK, redirect)
}
