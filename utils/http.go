// http.go contains useful http utilities.
package utils

import (
	"net/http"

	"github.com/docker/vetinari/errors"
)

type BetterHandler func(ctx IContext, w http.ResponseWriter, r *http.Request) *errors.HTTPError

type RootHandler struct {
	handler BetterHandler
	auth    IAuthorizer
	scopes  []IScope
	context IContextFactory
}

func RootHandlerFactory(auth IAuthorizer, ctxFac IContextFactory) func(BetterHandler, ...IScope) *RootHandler {
	return func(handler BetterHandler, scopes ...IScope) *RootHandler {
		return &RootHandler{handler, auth, scopes, ctxFac}
	}
}

func (root *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := root.context(r)
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
