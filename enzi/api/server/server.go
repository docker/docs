package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/authn"
	"github.com/emicklei/go-restful"
)

// Service is a base type for an API server.
type Service struct {
	*restful.WebService

	BaseContext context.Context
}

type (
	// Handler is the basic method type for all API request handlers.
	Handler func(context.Context, *restful.Request) responses.APIResponse
	// HandlerWithClientAccount includes context about the client user,
	// which may be anonymous.
	HandlerWithClientAccount func(context.Context, *authn.Account, *restful.Request) responses.APIResponse
)

// Route links method and path to its handler function.
type Route struct {
	Method  string
	Path    string // Relative to service path.
	Handler Handler

	Consumes []string
	Produces []string

	// Documentation
	Doc                string
	Notes              string
	PathParameterDocs  map[string]string
	QueryParameterDocs []*restful.Parameter
	FormParameterDocs  []*restful.Parameter
	BodySample         interface{}
	ResponseDocs       []ResponseDoc
	JSONErrs           []errors.APIError
}

// Register sets this route on the given service.
func (r *Route) Register(service *Service) {
	builder := service.WebService.
		Method(r.Method).
		Path(r.Path).
		To(service.WrapHandler(r.Handler)).
		Consumes(r.Consumes...).
		Produces(r.Produces...).
		Doc(r.Doc).
		Notes(r.Notes)

	for name, description := range r.PathParameterDocs {
		builder.Param(restful.PathParameter(name, description))
	}

	for _, queryParam := range r.QueryParameterDocs {
		builder.Param(queryParam)
	}

	for _, formParam := range r.FormParameterDocs {
		builder.Param(formParam)
	}

	if r.BodySample != nil {
		builder.Reads(r.BodySample)
	}

	for _, doc := range r.ResponseDocs {
		builder.Returns(doc.Code, doc.Message, doc.Sample)
	}

	service.WebService.Route(builder)
}

// ResponseDoc documents response and error types for an API endpoint.
type ResponseDoc struct {
	Code    int
	Message string
	Sample  interface{}
}

// WrapHandler wraps the given APIHandler by setting up the request context
// from this APIServer's baseContext, calling the given handler, and writing
// the response. It also intercepts a panic, logging it with a stack traces,
// writing a JSON error response, and re-panics so that any wrapping panic
// handlers may also handle the panic. Returns a restful.RouteFunction.
func (s *Service) WrapHandler(handler Handler) restful.RouteFunction {
	return restful.RouteFunction(func(request *restful.Request, response *restful.Response) {
		ctx := context.WithRequest(s.BaseContext, request.Request)
		logger := context.GetRequestLogger(ctx)
		ctx = context.WithLogger(ctx, logger)

		// insert the garant response writer in the middle between restful and its own response writer
		ctx, response.ResponseWriter = context.WithResponseWriter(ctx, response.ResponseWriter)

		defer func() {
			err := recover()
			if err == nil {
				// Not currently panicking.
				return
			}

			// Write a simple error response to the client.
			jsonResponse := responses.APIError(errors.Internal(ctx, fmt.Errorf("runtime panic: %v", err)))

			jsonResponse.WriteResponse(ctx, response)

			var stack []byte
			if stacker, ok := err.(errors.Stacker); ok {
				stack = stacker.Stack()
			} else {
				stack = debug.Stack()
			}

			// Push the stacktrace onto the context to be logged.
			ctx = context.WithValue(ctx, "stackTrace", string(stack))
			ctx = context.WithLogger(ctx, context.GetLogger(ctx, "stackTrace"))
			context.GetResponseLogger(ctx).Errorf("runtime panic: %v", err)
		}()

		handler(ctx, request).WriteResponse(ctx, response)
	})
}

// WrapHandlerWithAuthentication wraps the given handler with a regular Handler
// which first attempts to authenticate the client account using the given
// authenticator. The client MUST authenticate.
func WrapHandlerWithAuthentication(authenticator authn.RequestAuthenticator, handlerWithClientAccount HandlerWithClientAccount) Handler {
	return func(ctx context.Context, request *restful.Request) responses.APIResponse {
		account, apiErr := authenticator.AuthenticateRequest(ctx, request.Request)
		if apiErr != nil {
			return responses.APIError(apiErr)
		}

		if account.IsAnonymous {
			return responses.APIError(errors.AuthenticationRequired())
		}

		return handlerWithClientAccount(ctx, account, request)
	}
}

// WrapHandlerWithAdminAccount wraps the given handler with a regular Handler
// which first attempts to authenticate the client account using the given
// authenticator. The client account MUST be an admin to continue handling the
// request.
func WrapHandlerWithAdminAccount(authenticator authn.RequestAuthenticator, handlerWithAdminAccount HandlerWithClientAccount) Handler {
	return WrapHandlerWithAuthentication(authenticator, requireAdminWrapper(handlerWithAdminAccount))
}

func requireAdminWrapper(handlerWithAdminAccount HandlerWithClientAccount) HandlerWithClientAccount {
	return func(ctx context.Context, clientAccount *authn.Account, request *restful.Request) responses.APIResponse {
		if !clientAccount.IsAdmin {
			return responses.APIError(errors.NotAuthorized("system admin access is required"))
		}

		return handlerWithAdminAccount(ctx, clientAccount, request)
	}
}

// NewHTTPClient returns an HTTP client for accessing external API endpoints. If
// the given TLS config is nil, the system default config will be used.
func NewHTTPClient(tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSClientConfig:       tlsConfig,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			MaxIdleConnsPerHost:   5,
		},
	}
}
