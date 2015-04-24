package utils

import (
	"net/http"

	"github.com/endophage/go-tuf/signed"
	"github.com/gorilla/mux"

	"github.com/docker/vetinari/errors"
)

// contextHandler defines an alterate HTTP handler interface which takes in
// a context for authorization and returns an HTTP application error.
type contextHandler func(ctx Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError

// rootHandler is an implementation of an HTTP request handler which handles
// authorization and calling out to the defined alternate http handler.
type rootHandler struct {
	handler contextHandler
	auth    Authorizer
	scopes  []Scope
	context ContextFactory
	trust   signed.TrustService
}

// RootHandlerFactory creates a new rootHandler factory  using the given
// Context creator and authorizer.  The returned factory allows creating
// new rootHandlers from the alternate http handler contextHandler and
// a scope.
func RootHandlerFactory(auth Authorizer, ctxFac ContextFactory, trust signed.TrustService) func(contextHandler, ...Scope) *rootHandler {
	return func(handler contextHandler, scopes ...Scope) *rootHandler {
		return &rootHandler{handler, auth, scopes, ctxFac, trust}
	}
}

// ServeHTTP serves an HTTP request and implements the http.Handler interface.
func (root *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := root.context(r, root.trust)
	if err := root.auth.Authorize(ctx, root.scopes...); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err := root.handler(ctx, w, r); err != nil {
		// TODO: Log error
		http.Error(w, err.Error(), err.HTTPStatus)
		return
	}
	return
}

func ResourceName(r *http.Request) string {
	params := mux.Vars(r)
	if resource, ok := params["imageName"]; ok {
		return resource
	}
	return ""
}
