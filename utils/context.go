package utils

import (
	"github.com/endophage/go-tuf/signed"
	"net/http"
)

// Context defines an interface for managing authorizations.
type Context interface {
	// TODO: define a set of standard getters. Using getters
	//       will allow us to easily and transparently cache
	//       fields or load them on demand. Using this interface
	//       will allow people to define their own context struct
	//       that may handle things like caching and lazy loading
	//       differently.

	// Resource return the QDN of the resource being accessed
	Resource() string

	// Authorized returns a boolean indicating whether the user
	// has been successfully authorized for this request.
	Authorization() Authorization

	// SetAuthStatus should be called to change the authorization
	// status of the context (and therefore the request)
	SetAuthorization(Authorization)

	// Trust returns the trust service to be used
	Trust() signed.TrustService
}

// ContextFactory creates a IContext from an http request.
type ContextFactory func(*http.Request, signed.TrustService) Context

// Context represents an authorization context for a resource.
type context struct {
	resource      string
	authorization Authorization
	trust         signed.TrustService
}

// NewContext creates a new authorization context with the
// given HTTP request path as the resource.
func NewContext(r *http.Request, trust signed.TrustService) Context {
	return &context{
		resource: r.URL.Path,
		trust:    trust,
	}
}

// Resource returns the resource value for the context.
func (ctx *context) Resource() string {
	return ctx.resource
}

// Authorization returns an IAuthorization implementation for
// the context.
func (ctx *context) Authorization() Authorization {
	return ctx.authorization
}

// SetAuthorization allows setting an IAuthorization for
// the context.
func (ctx *context) SetAuthorization(authzn Authorization) {
	ctx.authorization = authzn
}

// Trust returns the instantiated TrustService for the context
func (ctx *context) Trust() signed.TrustService {
	return ctx.trust
}
