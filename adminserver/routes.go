package adminserver

import (
	"errors"
	"net/http"
	"net/http/pprof"
	"os"
	"path"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const logsRouteRoot = "logs"

func (a *AdminServer) wireSubroutes(router *mux.Router) {
	router.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir("ui")))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("ui"))))
	router.PathPrefix("/" + deploy.AdminSubroute + "/bin/").Handler(http.StripPrefix("/"+deploy.AdminSubroute+"/bin/", http.FileServer(http.Dir("bin"))))

	enziSubrouter := linkSubrouterWithMiddleware(router, "/"+deploy.EnziSubroute, negroni.New(
		a.authRedirectMiddleware(),
	))
	enziSubrouter.HandleFunc("/{path:.*}", a.enziProxyHandler)

	rethinkSubrouter := linkSubrouterWithMiddleware(router, "/"+deploy.RethinkSubroute, negroni.New(
		a.authRedirectMiddleware(),
	))
	rethinkSubrouter.HandleFunc("/{path:.*}", a.rethinkProxyHandler)

	auditLogsHandler := auditLogsHandler{Writer: a.syslogger}
	adminAPISubrouter := linkSubrouterWithMiddleware(router, "/api/v{version:0(?:\\.[0-9]+)?}/"+deploy.AdminSubroute, negroni.New(
		a.adminAPIAuthMiddleware(),
		auditLogsHandler.logAPICall(),
	))
	adminAPISubrouter.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		writeJSONError(rw, errors.New(http.StatusText(http.StatusNotFound)), http.StatusNotFound)
	})

	adminAPISubrouter.HandleFunc("/info", a.infoHandler).Methods("GET")
	// adminAPISubrouter.HandleFunc("/settings", a.updateSettingsHandler).Methods("PUT")
	// Registry
	// When creating an OPTIONS request return all available configurable fields for a PUT request to /settings/registry
	adminAPISubrouter.HandleFunc("/settings/registry", a.getRegistryConfigurationsHandler).Methods("OPTIONS")
	adminAPISubrouter.HandleFunc("/settings/registry", a.getRegistrySettingsHandler).Methods("GET")
	adminAPISubrouter.HandleFunc("/settings/registry", a.updateRegistrySettingsHandler).Methods("PUT")               // Entire YAML editing
	adminAPISubrouter.HandleFunc("/settings/registry/simple", a.updateRegistrySettingsViaFormHandler).Methods("PUT") // Form editing
	adminAPISubrouter.HandleFunc("/settings/license", a.updateLicenseSettingsHandler).Methods("PUT")
	adminAPISubrouter.HandleFunc("/settings/license", a.getLicenseSettingsHandler).Methods("GET")
	adminAPISubrouter.HandleFunc("/settings/license/toggle", a.toggleLicenseAutoRefreshHandler).Methods("PUT")
	adminAPISubrouter.HandleFunc("/clientAnalytics", a.clientAnalyticsHandler).Methods("GET")

	adminAPISubrouter.HandleFunc("/jobs", a.jobStatus).Methods("GET")
	adminAPISubrouter.HandleFunc("/jobs", a.runJob).Methods("POST")

	adminAPISubrouter.Handle("/upgrade", negroni.New(
		a.upgradeMiddleware(),
		negroni.Wrap(http.HandlerFunc(a.getUpgradeHandler)),
	)).Methods("GET")

	adminAPISubrouter.HandleFunc("/dhlogin", a.dockerHubLoginHandler).Methods("POST")

	// Account, Team, Repository, and Repository Access Control API endpoints.
	// These are separate from adminAPISubrouter because they handle all versions of the api
	// It's important that no other subrouter matches /api/v0/<new api path>
	apiSubrouterRestful := linkSubrouterWithMiddleware(router, "/api", negroni.New(
		auditLogsHandler.logAPICall(),
	))
	_, apiWSContainer := a.apiServer.BuildSubroutes("/api", "api server")
	// we tell the mux subrouter to give us everything because restful appears
	// as a single handler to it and restful takes over routing from here
	apiSubrouterRestful.Handle("/{path:.*}", apiWSContainer)

	// the version endpoint totally shouldn't be secret
	router.HandleFunc(path.Join("/", deploy.AdminSubroute, "version"), a.getVersionHandler).Methods("GET")
	adminSubrouter := linkSubrouterWithMiddleware(router, "/"+deploy.AdminSubroute, negroni.New(
		a.authRedirectMiddleware(),
		a.adminAuthMiddleware(),
	))

	adminSubrouter.HandleFunc("/", a.serveHTML).Methods("GET")

	// Catch alls for loading the UI on any URL
	adminSubrouter.HandleFunc("/{url:.*}", a.serveHTML).Methods("GET", "HEAD")

	// Get the CA
	router.HandleFunc("/ca", a.caHandler).Methods("GET")
	// health check for this specific replica
	router.HandleFunc("/health", a.healthHandler).Methods("GET")
	router.HandleFunc("/ws/events", a.eventWsHandler).Methods("GET")

	pprofEnabled, err := strconv.ParseBool(os.Getenv(deploy.PProfEnvVar))
	if err != nil {
		log.Debugf("Can't read pprof environment variable: %s", err)
	} else {
		if pprofEnabled {
			router.HandleFunc("/debug/pprof/", pprof.Index)
			router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			router.HandleFunc("/debug/pprof/profile", pprof.Profile)
			router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			router.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
			router.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
			router.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
			router.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		}
	}

	router.Handle("/{url:.*}", negroni.New(
		a.authRedirectMiddleware(),
		negroni.Wrap(http.HandlerFunc(a.serveHTML)),
	)).Methods("GET", "HEAD")
}

func linkSubrouterWithMiddleware(router *mux.Router, route string, middleware *negroni.Negroni) *mux.Router {
	subroute := router.PathPrefix(route)
	subrouter := router.PathPrefix(route).Subrouter()
	middleware.UseHandler(subrouter)
	subroute.Handler(middleware)
	return subrouter
}

func (a *AdminServer) buildRouter() *mux.Router {
	auditLogsHandler := auditLogsHandler{Writer: a.syslogger}

	router := mux.NewRouter()
	// This option must be off because if it's on, go-restful's swagger server is not able to serve
	// endpoints that end in / such as /api/docs/
	router.StrictSlash(false)
	a.wireSubroutes(router)

	auditLogsSubrouter := router.PathPrefix(path.Join("/", deploy.EventsEndpointSubroute)).Subrouter()
	auditLogsHandler.WireSubroutes(auditLogsSubrouter)

	return router
}
