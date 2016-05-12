package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/health"
	"github.com/docker/distribution/registry/api/errcode"
	"github.com/docker/distribution/registry/auth"
	"github.com/docker/notary"
	"github.com/docker/notary/server/errors"
	"github.com/docker/notary/server/handlers"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

func init() {
	data.SetDefaultExpiryTimes(notary.NotaryDefaultExpiries)
}

func prometheusOpts(operation string) prometheus.SummaryOpts {
	return prometheus.SummaryOpts{
		Namespace:   "notary_server",
		Subsystem:   "http",
		ConstLabels: prometheus.Labels{"operation": operation},
	}
}

// Config tells Run how to configure a server
type Config struct {
	Addr                         string
	TLSConfig                    *tls.Config
	Trust                        signed.CryptoService
	AuthMethod                   string
	AuthOpts                     interface{}
	RepoPrefixes                 []string
	ConsistentCacheControlConfig utils.CacheControlConfig
	CurrentCacheControlConfig    utils.CacheControlConfig
}

// Run sets up and starts a TLS server that can be cancelled using the
// given configuration. The context it is passed is the context it should
// use directly for the TLS server, and generate children off for requests
func Run(ctx context.Context, conf Config) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.Addr)
	if err != nil {
		return err
	}
	var lsnr net.Listener
	lsnr, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	if conf.TLSConfig != nil {
		logrus.Info("Enabling TLS")
		lsnr = tls.NewListener(lsnr, conf.TLSConfig)
	}

	var ac auth.AccessController
	if conf.AuthMethod == "token" {
		authOptions, ok := conf.AuthOpts.(map[string]interface{})
		if !ok {
			return fmt.Errorf("auth.options must be a map[string]interface{}")
		}
		ac, err = auth.GetAccessController(conf.AuthMethod, authOptions)
		if err != nil {
			return err
		}
	}

	svr := http.Server{
		Addr: conf.Addr,
		Handler: RootHandler(
			ac, ctx, conf.Trust,
			conf.ConsistentCacheControlConfig, conf.CurrentCacheControlConfig,
			conf.RepoPrefixes),
	}

	logrus.Info("Starting on ", conf.Addr)

	err = svr.Serve(lsnr)

	return err
}

// assumes that required prefixes is not empty
func filterImagePrefixes(requiredPrefixes []string, err error, handler http.Handler) http.Handler {
	if len(requiredPrefixes) == 0 {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageName := mux.Vars(r)["imageName"]

		for _, prefix := range requiredPrefixes {
			if strings.HasPrefix(imageName, prefix) {
				handler.ServeHTTP(w, r)
				return
			}
		}

		errcode.ServeJSON(w, err)
	})
}

type _serverEndpoint struct {
	OperationName       string
	ServerHandler       utils.ContextHandler
	ErrorIfGUNInvalid   error
	IncludeCacheHeaders bool
	CacheControlConfig  utils.CacheControlConfig
	PermissionsRequired []string
}

// RootHandler returns the handler that routes all the paths from / for the
// server.
func RootHandler(ac auth.AccessController, ctx context.Context, trust signed.CryptoService,
	consistent, current utils.CacheControlConfig, repoPrefixes []string) http.Handler {

	authWrapper := utils.RootHandlerFactory(ac, ctx, trust)

	createHandler := func(opts _serverEndpoint) http.Handler {
		var wrapped http.Handler
		wrapped = authWrapper(opts.ServerHandler, opts.PermissionsRequired...)
		if opts.IncludeCacheHeaders {
			wrapped = utils.WrapWithCacheHandler(opts.CacheControlConfig, wrapped)
		}
		wrapped = filterImagePrefixes(repoPrefixes, opts.ErrorIfGUNInvalid, wrapped)
		return prometheus.InstrumentHandlerWithOpts(prometheusOpts(opts.OperationName), wrapped)
	}

	invalidGUNErr := errors.ErrInvalidGUN.WithDetail(fmt.Sprintf("Require GUNs with prefix: %v", repoPrefixes))
	notFoundError := errors.ErrMetadataNotFound.WithDetail(nil)

	r := mux.NewRouter()
	r.Methods("GET").Path("/v2/").Handler(authWrapper(handlers.MainHandler))

	r.Methods("POST").Path("/v2/{imageName:.*}/_trust/tuf/").Handler(createHandler(_serverEndpoint{
		OperationName:       "UpdateTuf",
		ErrorIfGUNInvalid:   invalidGUNErr,
		ServerHandler:       handlers.AtomicUpdateHandler,
		PermissionsRequired: []string{"push", "pull"},
	}))
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/{tufRole:root|targets(?:/[^/\\s]+)*|snapshot|timestamp}.{checksum:[a-fA-F0-9]{64}|[a-fA-F0-9]{96}|[a-fA-F0-9]{128}}.json").Handler(createHandler(_serverEndpoint{
		OperationName:       "GetRoleByHash",
		ErrorIfGUNInvalid:   notFoundError,
		IncludeCacheHeaders: true,
		CacheControlConfig:  consistent,
		ServerHandler:       handlers.GetHandler,
		PermissionsRequired: []string{"pull"},
	}))
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/{tufRole:root|targets(?:/[^/\\s]+)*|snapshot|timestamp}.json").Handler(createHandler(_serverEndpoint{
		OperationName:       "GetRole",
		ErrorIfGUNInvalid:   notFoundError,
		IncludeCacheHeaders: true,
		CacheControlConfig:  current,
		ServerHandler:       handlers.GetHandler,
		PermissionsRequired: []string{"pull"},
	}))
	r.Methods("GET").Path(
		"/v2/{imageName:.*}/_trust/tuf/{tufRole:snapshot|timestamp}.key").Handler(createHandler(_serverEndpoint{
		OperationName:       "GetKey",
		ErrorIfGUNInvalid:   notFoundError,
		ServerHandler:       handlers.GetKeyHandler,
		PermissionsRequired: []string{"push", "pull"},
	}))
	r.Methods("DELETE").Path("/v2/{imageName:.*}/_trust/tuf/").Handler(createHandler(_serverEndpoint{
		OperationName:       "DeleteTuf",
		ErrorIfGUNInvalid:   notFoundError,
		ServerHandler:       handlers.DeleteHandler,
		PermissionsRequired: []string{"push", "pull"},
	}))

	r.Methods("GET").Path("/_notary_server/health").HandlerFunc(health.StatusHandler)
	r.Methods("GET").Path("/metrics").Handler(prometheus.Handler())
	r.Methods("GET", "POST", "PUT", "HEAD", "DELETE").Path("/{other:.*}").Handler(
		authWrapper(handlers.NotFoundHandler))

	return r
}
