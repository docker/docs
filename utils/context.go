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

	// Scopes returns the scopes required for the request to be
	// successfully completed.
	Scopes() []Scope

	// Authorized returns a boolean indicating whether the user
	// has been successfully authorized for this request.
	Authorized() bool

	// SetAuthStatus should be called to change the authorization
	// status of the context (and therefore the request)
	SetAuthStatus(bool)
}

type Context struct {
	resource   string
	scopes     []Scope
	authorized bool
}

func generateContext(r *http.Request) Context {
	return Context{authorized: false}
}

func (ctx *Context) Resource() string {
	return ctx.resource
}

func (ctx *Context) Scopes() string {
	return ctx.scopes
}

func (ctx *Context) Authorized() string {
	return ctx.authorized
}

func (ctx *Context) SetAuthStatus(newStatus bool) string {
	ctx.authorized = newStatus
}
