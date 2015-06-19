package utils

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/auth"
	"github.com/endophage/gotuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/notary/errors"
)

// contextHandler defines an alterate HTTP handler interface which takes in
// a context for authorization and returns an HTTP application error.
type contextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError

// rootHandler is an implementation of an HTTP request handler which handles
// authorization and calling out to the defined alternate http handler.
type rootHandler struct {
	handler contextHandler
	auth    auth.AccessController
	actions []string
	context context.Context
	trust   signed.CryptoService
}

// RootHandlerFactory creates a new rootHandler factory  using the given
// Context creator and authorizer.  The returned factory allows creating
// new rootHandlers from the alternate http handler contextHandler and
// a scope.
func RootHandlerFactory(auth auth.AccessController, ctx context.Context, trust signed.CryptoService) func(contextHandler, ...string) *rootHandler {
	return func(handler contextHandler, actions ...string) *rootHandler {
		return &rootHandler{
			handler: handler,
			auth:    auth,
			actions: actions,
			context: ctx,
			trust:   trust,
		}
	}
}

// ServeHTTP serves an HTTP request and implements the http.Handler interface.
func (root *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.WithValue(root.context, "repo", vars["imageName"])
	ctx = context.WithValue(ctx, "trust", root.trust)
	ctx = context.WithValue(ctx, "http.request", r)

	//	access := buildAccessRecords(vars["imageName"], root.actions...)
	//	var err error
	//	if ctx, err = root.auth.Authorized(ctx, access...); err != nil {
	//		http.Error(w, err.Error(), http.StatusUnauthorized)
	//		return
	//	}
	if err := root.handler(ctx, w, r); err != nil {
		logrus.Error("[Notary] ", err.Error())
		http.Error(w, err.Error(), err.HTTPStatus)
		return
	}
	return
}

func buildAccessRecords(repo string, actions ...string) []auth.Access {
	requiredAccess := make([]auth.Access, 0, len(actions))
	for _, action := range actions {
		requiredAccess = append(requiredAccess, auth.Access{
			Resource: auth.Resource{
				Type: "repo",
				Name: repo,
			},
			Action: action,
		})
	}
	return requiredAccess
}
