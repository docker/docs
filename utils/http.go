package utils

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/distribution/registry/auth"
	"github.com/docker/notary/tuf/signed"
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
	ctx := ctxu.WithRequest(root.context, r)
	ctx, w = ctxu.WithResponseWriter(ctx, w)
	ctx = ctxu.WithLogger(ctx, ctxu.GetRequestLogger(ctx))
	ctx = context.WithValue(ctx, "repo", vars["imageName"])
	ctx = context.WithValue(ctx, "cryptoService", root.trust)

	defer func() {
		ctxu.GetResponseLogger(ctx).Info("response completed")
	}()

	if root.auth != nil {
		access := buildAccessRecords(vars["imageName"], root.actions...)
		var authCtx context.Context
		var err error
		if authCtx, err = root.auth.Authorized(ctx, access...); err != nil {
			if err, ok := err.(auth.Challenge); ok {
				err.ServeHTTP(w, r)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			errcode.ServeJSON(w, v2.ErrorCodeUnauthorized)
			return
		}
		ctx = authCtx
	}
	if err := root.handler(ctx, w, r); err != nil {
		e := errcode.ServeJSON(w, err)
		if e != nil {
			logrus.Error(e)
		}
		return
	}
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
