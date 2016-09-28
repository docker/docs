package api

import (
	"expvar"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mailgun/oxy/forward"

	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/controller/middleware/access"
	"github.com/docker/orca/controller/middleware/audit"
	mAuth "github.com/docker/orca/controller/middleware/auth"
	"github.com/docker/orca/controller/middleware/notary"
	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/docker/orca/controller/middleware/tracking"
	"github.com/docker/orca/tlsutils"
)

type (
	Api struct {
		listenAddr         string
		manager            manager.Manager
		authWhitelistCIDRs []string
		serverVersion      string
		allowInsecure      bool
		swarmCAPEM         []byte
		controllerCAPEM    []byte
		controllerCertPEM  []byte
		controllerKeyPEM   []byte
		swarmClassicURL    string
		engineProxyURL     string
		fwd                *forward.Forwarder
		enableProfiling    bool
		Router             *mux.Router
		tlsCACert          string
	}

	ApiConfig struct {
		ListenAddr         string
		Manager            manager.Manager
		AuthWhiteListCIDRs []string
		AllowInsecure      bool

		// CA cert/key for creating and signing client certs
		SwarmCAPEM      []byte
		ControllerCAPEM []byte
		// server cert/key for TLS
		ControllerCertPEM []byte
		ControllerKeyPEM  []byte

		// Turn on remote pprof support
		EnableProfiling bool
	}

	Credentials struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
)

func NewApi(config ApiConfig) (*Api, error) {
	return &Api{
		listenAddr:         config.ListenAddr,
		manager:            config.Manager,
		authWhitelistCIDRs: config.AuthWhiteListCIDRs,
		allowInsecure:      config.AllowInsecure,
		swarmCAPEM:         config.SwarmCAPEM,
		controllerCAPEM:    config.ControllerCAPEM,
		controllerCertPEM:  config.ControllerCertPEM,
		controllerKeyPEM:   config.ControllerKeyPEM,
		enableProfiling:    config.EnableProfiling,
		Router:             mux.NewRouter(),
	}, nil
}

// Initialize wires up all API routes and middleware layers together,
// using the provided MiddlewarePipeline
// Initialize is called by Run() and all routes are matched by a.Router
func (a *Api) Initialize(m pipeline.MiddlewarePipeline) {
	// Set up middleware layers
	auditExcludes := []string{
		"^/.*/containers/json",
		"^/.*/images/json",
	}
	trackingExcludes := []string{
		"^/containers/.*/json",
		"^/.*/containers/json",
		"^/.*/images/json",
	}
	authRequired := mAuth.NewAuthRequired(a.manager, a.authWhitelistCIDRs)
	apiAuditor := audit.NewAuditor(a.manager, auditExcludes)
	apiTracker := tracking.NewTracker(a.manager, trackingExcludes)
	tmpDir, err := ioutil.TempDir("", "trust")
	if err != nil {
		log.Errorf("Could not create tmp dir for notary middleware in middleware initializer: %s", err.Error())
		tmpDir = "/tmp/trust"
	}

	notaryMW := notary.NewNotaryMiddleware(
		a.manager,
		tmpDir,
		a.getAuthHeaderForGUN,
	)

	// The pipeline object itself is discarded after route registration
	m.AddInitializer(authRequired.Initializer)
	m.AddLayer(access.LayerHandler)
	m.AddLayer(notaryMW.LayerHandler)
	m.AddLayer(apiAuditor.LayerHandler)
	m.AddLayer(apiTracker.LayerHandler)

	// Register all routes in this specific order.
	// Do not reorder these definitions.
	log.Info("Registering UCP Routes")
	a.RegisterOrcaRoutes(m)
	log.Info("Registering Docker Remote API Routes")
	a.RegisterSwarmRoutes(m)
	a.RegisterEngineRoutes(m)
	log.Info("Registering Public Routes")
	a.RegisterPublicRoutes()

	// Do not add any new routes to a.Router after this point
	// Doing so will result in a 404, as all remaining names are served
	// by the static file handler
}

func (a *Api) Run() error {
	// Set up logging
	proxyLogger := log.New()
	proxyLogger.Level = log.WarnLevel
	proxyLogger.Formatter = &log.JSONFormatter{}
	fl := forward.Logger(proxyLogger)

	if a.manager.SwarmClassicURL() == nil {
		return fmt.Errorf("Unable to initialize controller with nil swarm V1 manager URL")
	}
	a.swarmClassicURL = "https://" + a.manager.SwarmClassicURL().Host

	if a.manager.EngineProxyURL() == nil {
		return fmt.Errorf("Unable to initialize controller with nil docker engine URL")
	}
	a.engineProxyURL = "https://" + a.manager.EngineProxyURL().Host

	t := a.manager.DockerClientTransport()

	// check if TLS is enabled and configure if so
	if t.TLSClientConfig == nil || len(t.TLSClientConfig.Certificates) == 0 {
		return fmt.Errorf("TLS configuration invalid or missing certificates. Aborting API initialization.")
	}

	// setup custom roundtripper with TLS transport
	r := forward.RoundTripper(
		&http.Transport{
			TLSClientConfig: t.TLSClientConfig,
		})
	f, err := forward.New(r, fl)
	if err != nil {
		log.Fatal(err)
	}
	a.fwd = f

	// Create the middleware pipeline
	m := pipeline.New(a.Router)

	// Set up all middleware and routes on the pipeline.
	a.Initialize(m)

	// Start the HTTP Server using a.Router
	log.Infof("controller listening on %s", a.listenAddr)
	s := &http.Server{
		Addr:    a.listenAddr,
		Handler: a.Router,
	}

	// Handle the Non-TLS case
	if len(a.controllerCertPEM) == 0 || len(a.controllerKeyPEM) == 0 {
		return s.ListenAndServe()
	}

	// setup TLS config
	var caCert []byte
	if len(a.controllerCAPEM) > 0 {
		caCert = a.controllerCAPEM
		a.tlsCACert = string(caCert)

		// Make sure if the user added additional external trust, we include it into cert bundles
		a.manager.AddTrustedCert(string(caCert))
	}
	if len(a.swarmCAPEM) > 0 {

		caCert = append(caCert, byte('\n'))
		caCert = append(caCert, a.swarmCAPEM...)

		// No need to add this as a new trusted cert, since it will be added automatically
	}

	tlsConfig, err := tlsutils.GetServerTLSConfig(caCert, a.controllerCertPEM, a.controllerKeyPEM, a.allowInsecure)
	if err != nil {
		return err
	}

	s.TLSConfig = tlsConfig

	// We do not need to specify the cert/key filenames here because they
	// have already been set in s.TLSConfig.
	return s.ListenAndServeTLS("", "")
}

// Replicated from expvar.go as not public.
func expVars(w http.ResponseWriter, r *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
