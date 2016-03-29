package utils

import (
	"fmt"
	"net/http"
	"time"

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
			if challenge, ok := err.(auth.Challenge); ok {
				// Let the challenge write the response.
				challenge.ServeHTTP(w, r)

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

// CacheControlConfig is an interface for something that knows how to set cache
// control headers
type CacheControlConfig interface {
	// SetHeaders will actually set the cache control headers on a Headers object
	SetHeaders(headers http.Header)
}

// NewCacheControlConfig returns CacheControlConfig interface for either setting
// cache control or disabling cache control entirely
func NewCacheControlConfig(maxAgeInSeconds int, mustRevalidate bool) CacheControlConfig {
	if maxAgeInSeconds > 0 {
		return PublicCacheControl{MustReValidate: mustRevalidate, MaxAgeInSeconds: maxAgeInSeconds}
	}
	return NoCacheControl{}
}

// PublicCacheControl is a set of options that we will set to enable cache control
type PublicCacheControl struct {
	MustReValidate  bool
	MaxAgeInSeconds int
}

// SetHeaders sets the public headers with an optional must-revalidate header
func (p PublicCacheControl) SetHeaders(headers http.Header) {
	cacheControlValue := fmt.Sprintf("public, max-age=%v, s-maxage=%v",
		p.MaxAgeInSeconds, p.MaxAgeInSeconds)

	if p.MustReValidate {
		cacheControlValue = fmt.Sprintf("%s, must-revalidate", cacheControlValue)
	}
	headers.Set("Cache-Control", cacheControlValue)
	// delete the Pragma directive, because the only valid value in HTTP is
	// "no-cache"
	headers.Del("Pragma")
	if headers.Get("Last-Modified") == "" {
		SetLastModifiedHeader(headers, time.Time{})
	}
}

// NoCacheControl is an object which represents a directive to cache nothing
type NoCacheControl struct{}

// SetHeaders sets the public headers cache-control headers and pragma to no-cache
func (n NoCacheControl) SetHeaders(headers http.Header) {
	headers.Set("Cache-Control", "max-age=0, no-cache, no-store")
	headers.Set("Pragma", "no-cache")
}

// cacheControlResponseWriter wraps an existing response writer, and if Write is
// called, will try to set the cache control headers if it can
type cacheControlResponseWriter struct {
	http.ResponseWriter
	config     CacheControlConfig
	statusCode int
}

// WriteHeader stores the header before writing it, so we can tell if it's been set
// to a non-200 status code
func (c *cacheControlResponseWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

// Write will set the cache headers if they haven't already been set and if the status
// code has either not been set or set to 200
func (c *cacheControlResponseWriter) Write(data []byte) (int, error) {
	if c.statusCode == http.StatusOK || c.statusCode == 0 {
		headers := c.ResponseWriter.Header()
		if headers.Get("Cache-Control") == "" {
			c.config.SetHeaders(headers)
		}
	}
	return c.ResponseWriter.Write(data)
}

type cacheControlHandler struct {
	http.Handler
	config CacheControlConfig
}

func (c cacheControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler.ServeHTTP(&cacheControlResponseWriter{ResponseWriter: w, config: c.config}, r)
}

// WrapWithCacheHandler wraps another handler in one that can add cache control headers
// given a 200 response
func WrapWithCacheHandler(ccc CacheControlConfig, handler http.Handler) http.Handler {
	if ccc != nil {
		return cacheControlHandler{Handler: handler, config: ccc}
	}
	return handler
}

// SetLastModifiedHeader takes a time and uses it to set the LastModified header using
// the right date format
func SetLastModifiedHeader(headers http.Header, lmt time.Time) {
	headers.Set("Last-Modified", lmt.Format(time.RFC1123))
}
