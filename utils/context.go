package utils

import (
	"github.com/endophage/go-tuf/signed"
	"net/http"
)

// IContext defines an interface for managing authorizations.
type IContext interface {
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
	Authorization() IAuthorization

	// SetAuthStatus should be called to change the authorization
	// status of the context (and therefore the request)
	SetAuthorization(IAuthorization)

	Signer() *signed.Signer
}

// IContextFactory creates a IContext from an http request.
type IContextFactory func(*http.Request) IContext

// Context represents an authorization context for a resource.
type Context struct {
	resource      string
	authorization IAuthorization
}

// ContextFactory creates a new authorization context with the
// given HTTP request path as the resource.
func ContextFactory(r *http.Request) IContext {
	return &Context{
		resource: r.URL.Path,
	}
}

// Resource returns the resource value for the context.
func (ctx *Context) Resource() string {
	return ctx.resource
}

// Authorization returns an IAuthorization implementation for
// the context.
func (ctx *Context) Authorization() IAuthorization {
	return ctx.authorization
}

// SetAuthorization allows setting an IAuthorization for
// the context.
func (ctx *Context) SetAuthorization(authzn IAuthorization) {
	ctx.authorization = authzn
}

// Signer returns the instantiated signer for the context
func (ctx *Context) Signer() *signed.Signer {
	return nil
}
