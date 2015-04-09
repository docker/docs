package utils

import (
	"net/http"
)

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
}

type IContextFactory func(*http.Request) IContext

type Context struct {
	resource      string
	authorization IAuthorization
}

func ContextFactory(r *http.Request) IContext {
	return &Context{
		resource: r.URL.Path,
	}
}

func (ctx *Context) Resource() string {
	return ctx.resource
}

func (ctx *Context) Authorization() IAuthorization {
	return ctx.authorization
}

func (ctx *Context) SetAuthorization(authzn IAuthorization) {
	ctx.authorization = authzn
}
