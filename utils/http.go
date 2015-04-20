package utils

import (
	"net/http"

	"github.com/docker/vetinari/errors"
	"github.com/endophage/go-tuf/signed"
)

// BetterHandler defines an alterate HTTP handler interface which takes in
// a context for authorization and returns an HTTP application error.
type BetterHandler func(ctx IContext, w http.ResponseWriter, r *http.Request) *errors.HTTPError

// RootHandler is an implementation of an HTTP request handler which handles
// authorization and calling out to the defined alternate http handler.
type RootHandler struct {
	handler BetterHandler
	auth    IAuthorizer
	scopes  []IScope
	context IContextFactory
	trust   signed.TrustService
}

// RootHandlerFactory creates a new RootHandler factory  using the given
// Context creator and authorizer.  The returned factory allows creating
// new RootHandlers from the alternate http handler BetterHandler and
// a scope.
func RootHandlerFactory(auth IAuthorizer, ctxFac IContextFactory, trust signed.TrustService) func(BetterHandler, ...IScope) *RootHandler {
	return func(handler BetterHandler, scopes ...IScope) *RootHandler {
		return &RootHandler{handler, auth, scopes, ctxFac, trust}
	}
}

// ServeHTTP serves an HTTP request and implements the http.Handler interface.
func (root *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
