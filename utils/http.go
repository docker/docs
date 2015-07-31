package utils

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/distribution/registry/auth"
	"github.com/docker/notary/errors"
	"github.com/endophage/gotuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// contextHandler defines an alterate HTTP handler interface which takes in
// a context for authorization and returns an HTTP application error.
type contextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// rootHandler is an implementation of an HTTP request handler which handles
// authorization and calling out to the defined alternate http handler.
type rootHandler struct {
	handler contextHandler
	auth    auth.AccessController
	actions []string
	context context.Context
	trust   signed.CryptoService
	//cachePool redis.Pool
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

	ctx = context.WithValue(ctx, "cryptoService", root.trust)

	ctx = context.WithValue(ctx, "http.request", r)

	if root.auth != nil {
		var err error
		access := buildAccessRecords(vars["imageName"], root.actions...)
		if ctx, err = root.auth.Authorized(ctx, access...); err != nil {
			if err, ok := err.(auth.Challenge); ok {
				err.ServeHTTP(w, r)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			errcode.ServeJSON(w, v2.ErrorCodeUnauthorized)
			return
		}
	}
	if err := root.handler(ctx, w, r); err != nil {
		if err, ok := err.(errcode.Error); ok {
			logrus.Errorf(
				"[Notary Server] %d %s %s",
				err.Code.Descriptor().HTTPStatusCode,
				r.Method,
				r.URL.Path,
			)
		} else {
			logrus.Errorf(
				"[Notary Server] 5XX %s %s %s",
				r.Method,
				r.URL.Path,
				err.Error(),
			)
		}
		e := errcode.ServeJSON(w, err)
		if e != nil {
			logrus.Error(e)
		}
		return
	}
	logrus.Infof("[Notary Server] 200 %s %s", r.Method, r.URL.Path)
	return
}

func buildAccessRecords(repo string, actions ...string) []auth.Access {
	requiredAccess := make([]auth.Access, 0, len(actions))
	for _, action := range actions {
		requiredAccess = append(requiredAccess, auth.Access{
			Resource: auth.Resource{
				Type: "repository",
				Name: repo,
			},
			Action: action,
		})
	}
	return requiredAccess
}

// NotFoundHandler is used as a generic catch all handler to return the ErrMetadataNotFound
// 404 response
func NotFoundHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.ErrMetadataNotFound.WithDetail(nil)
}
