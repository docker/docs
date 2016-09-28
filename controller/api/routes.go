package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"net/url"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/docker/orca/auth/enzi"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/middleware/pipeline"
)

type swarmPipeline struct {
	pipeline pipeline.MiddlewarePipeline
}

func newSwarmPipeline(p pipeline.MiddlewarePipeline) *swarmPipeline {
	return &swarmPipeline{
		pipeline: p,
	}
}

func (p *swarmPipeline) Route(path string, method string, parser pipeline.Parser, handler pipeline.Handler) {
	versionedParser := func(rc *ctx.OrcaRequestContext) (int, error) {
		version, ok := rc.PathVars["version"]
		if !ok {
			return http.StatusInternalServerError, fmt.Errorf("Unable to parse the Version URL argument")
		}

		if !isSupported(version) {
			return http.StatusBadRequest, fmt.Errorf("API version %s is not supported by UCP", version)
		}

		return parser(rc)
	}
	versionedPath := string("/v{version:[0-9.]+}" + path)

	p.pipeline.Route(path, method, parser, handler)
	p.pipeline.Route(versionedPath, method, versionedParser, handler)
}

// registerOrcaRoutes defines the UCP API routes and handlers
func (a *Api) RegisterOrcaRoutes(p pipeline.MiddlewarePipeline) {
	p.Route("/api/account", "GET", a.nilParser, a.currentAccount)

	p.Route("/api/accounts", "GET", a.adminParser, a.accounts)
	// POST /api/accounts is used to set a user's own Public Key.
	// Authorization is deferred to the authenticator in this case
	p.Route("/api/accounts", "POST", a.nilParser, a.saveAccount)
	p.Route("/api/accounts/{username}", "GET", a.userParser, a.account)
	p.Route("/api/accounts/{username}", "DELETE", a.adminParser, a.deleteAccount)

	p.Route("/api/accesslists", "GET", a.userAccessListParser, a.accessLists)
	p.Route("/api/accesslists", "POST", a.adminParser, a.addAccessList)
	p.Route("/api/accesslists/{teamId}/{id}", "GET", a.adminParser, a.accessList)
	p.Route("/api/accesslists/{teamId}/{id}", "DELETE", a.adminParser, a.removeAccessList)

	p.Route("/api/teams", "GET", a.adminParser, a.teams)
	p.Route("/api/teams", "POST", a.adminParser, a.saveTeam)
	p.Route("/api/teams/{id}", "GET", a.teamParser, a.team)
	p.Route("/api/teams/{id}", "DELETE", a.adminParser, a.deleteTeam)
	p.Route("/api/teams/{id}", "PUT", a.adminParser, a.updateTeam)
	p.Route("/api/teams/{id}/members/add/{username}", "PUT", a.adminParser, a.addMemberToTeam)
	p.Route("/api/teams/{id}/members/remove/{username}", "DELETE", a.adminParser, a.removeMemberFromTeam)

	p.Route("/api/config", "GET", a.adminParser, a.listConfig)
	p.Route("/api/config/{subsystem}", "GET", a.adminParser, a.getConfig)
	p.Route("/api/config/{subsystem}", "POST", a.adminParser, a.updateConfig)

	p.Route("/api/banner", "GET", a.nilParser, a.banner)
	p.Route("/api/catalog", "GET", a.nilParser, a.catalog)
	p.Route("/api/clientbundle", "GET", a.nilParser, a.generateClientBundle)
	p.Route("/api/clientbundle", "POST", a.nilParser, a.generateClientBundle)

	p.Route("/api/applications", "GET", a.nilParser, a.applications)
	p.Route("/api/applications/{name}", "GET", a.nilParser, a.application)
	p.Route("/api/applications/{name}/containers", "GET", a.nilParser, a.applicationContainers)
	p.Route("/api/applications/{name}/restart", "POST", a.nilParser, a.restartApplication)
	p.Route("/api/applications/{name}/stop", "POST", a.nilParser, a.stopApplication)
	p.Route("/api/applications/{name}", "DELETE", a.nilParser, a.removeApplication)

	p.Route("/api/nodes", "GET", a.nilParser, a.nodes)
	p.Route("/api/nodes/{name}", "GET", a.nilParser, a.node)
	p.Route("/api/nodes/certs", "POST", a.adminParser, a.setControllerServerCerts)

	// TODO(alexmavr): determine expected behavior for webhookkeys
	p.Route("/api/webhookkeys", "GET", a.adminParser, a.webhookKeys)
	p.Route("/api/webhookkeys", "POST", a.adminParser, a.addWebhookKey)
	p.Route("/api/webhookkeys/{key}", "GET", a.adminParser, a.webhookKey)
	p.Route("/api/webhookkeys/{key}", "DELETE", a.adminParser, a.deleteWebhookKey)

	p.Route("/api/support", "POST", a.adminParser, a.supportDump)
	p.Route("/api/authsync", "GET", a.adminParser, a.authSyncMessages)
	p.Route("/api/authsync", "POST", a.adminParser, a.doAuthSync)

	p.Route("/api/consolesession/{id}", "POST", a.containerFullControl, a.createConsoleSession)
	p.Route("/api/consolesession/{token}", "GET", a.nilParser, a.consoleSession)
	p.Route("/api/consolesession/{token}", "DELETE", a.nilParser, a.removeConsoleSession)
	p.Route("/api/containers/{id}/scale", "POST", a.containerRestrictedControl, a.scaleContainer)
	p.Route("/api/containerlogs/{id}", "POST", a.containerViewOnly, a.createContainerLogsToken)
	p.Route("/api/containerlogs/{token}", "GET", a.nilParser, a.containerLogsToken)
	p.Route("/api/containerlogs/{token}", "DELETE", a.nilParser, a.removeContainerLogsToken)

	p.Route("/account/changepassword", "POST", a.nilParser, a.changePassword)
}

// RegisterEngineRoutes handles the Remote API routes that are redirected
// directly to the local engine proxy.
func (a *Api) RegisterEngineRoutes(p pipeline.MiddlewarePipeline) {
	// a swarmPipeline is used to handle API versioning of routes
	sp := newSwarmPipeline(p)

	sp.Route("/swarm", "GET", a.nilParser, a.engineRedirect)
	sp.Route("/swarm/init", "POST", a.nilParser, a.swarmModeInit)
	sp.Route("/swarm/join", "POST", a.nilParser, a.swarmModeJoin)
	sp.Route("/swarm/leave", "POST", a.nilParser, a.swarmModeLeave)
	sp.Route("/swarm/update", "POST", a.adminParser, a.engineRedirect)

	sp.Route("/services", "GET", a.nilParser, a.listServices)
	sp.Route("/services/{id:.*}", "GET", a.serviceViewOnly, a.engineRedirect)
	sp.Route("/services/create", "POST", a.serviceCreateParser, a.createService)
	sp.Route("/services/{id:.*}/update", "POST", a.serviceRestrictedControl, a.engineRedirect)
	sp.Route("/services/{id:.*}", "DELETE", a.serviceRestrictedControl, a.engineRedirect)

	sp.Route("/nodes", "GET", a.nilParser, a.listNodes)
	sp.Route("/nodes/{id}", "GET", a.nilParser, a.inspectNode)
	sp.Route("/nodes/{id}", "DELETE", a.adminParser, a.engineRedirect)
	sp.Route("/nodes/{id}/update", "POST", a.adminParser, a.engineRedirect)

	// TODO: Role-Based filtering for task endpoints
	sp.Route("/tasks", "GET", a.nilParser, a.engineRedirect)
	sp.Route("/tasks/{id:.*}", "GET", a.nilParser, a.engineRedirect)
}

// RegisterSwarmRoutes defines the Docker Swarm API routes that are handled by UCP
func (a *Api) RegisterSwarmRoutes(p pipeline.MiddlewarePipeline) {
	sp := newSwarmPipeline(p)

	sp.Route("/events", "GET", a.nilParser, a.swarmHijack)
	sp.Route("/info", "GET", a.nilParser, a.swarmInfo)
	sp.Route("/version", "GET", a.nilParser, a.swarmVersion)

	sp.Route("/images/json", "GET", a.nilParser, a.listImages)
	sp.Route("/images/viz", "GET", a.nilParser, a.swarmImages)
	sp.Route("/images/search", "GET", a.nilParser, a.swarmImages)
	sp.Route("/images/get", "GET", a.nilParser, a.swarmImages)
	sp.Route("/images/{name:.*}/get", "GET", a.nilParser, a.swarmImages)
	sp.Route("/images/{name:.*}/history", "GET", a.nilParser, a.swarmImages)
	sp.Route("/images/{name:.*}/json", "GET", a.nilParser, a.swarmImagesInfo)
	sp.Route("/images/create", "POST", a.nilParser, a.swarmRegistryImagesCreate)
	sp.Route("/images/load", "POST", a.nilParser, a.swarmImages)
	sp.Route("/images/{name:.*}/push", "POST", a.nilParser, a.swarmRegistryImagesPush)
	sp.Route("/images/{name:.*}/tag", "POST", a.nilParser, a.swarmImages)
	sp.Route("/images/{name:.*}", "DELETE", a.nilParser, a.swarmImages)

	// TODO: RBAC for volumes
	sp.Route("/volumes", "GET", a.nilParser, a.swarmRedirect)
	sp.Route("/volumes/{name:.*}", "GET", a.nilParser, a.swarmRedirect)
	sp.Route("/volumes/create", "POST", a.nilParser, a.swarmRedirect)
	sp.Route("/volumes/{name:.*}", "DELETE", a.nilParser, a.swarmRedirect)

	// Network Resource API endpoints
	// The connect and disconnect endpoints have special parsers as they involve specific containers
	sp.Route("/networks", "GET", a.nilParser, a.listNetworks)
	sp.Route("/networks/{id:.*}", "GET", a.networkViewOnly, a.engineRedirect)
	sp.Route("/networks/create", "POST", a.networkCreateParser, a.createNetwork)
	sp.Route("/networks/{id:.*}/connect", "POST", a.networkConnectParser, a.engineRedirect)
	sp.Route("/networks/{id:.*}/disconnect", "POST", a.networkDisconnectParser, a.engineRedirect)
	sp.Route("/networks/{id:.*}", "DELETE", a.networkRestrictedControl, a.engineRedirect)

	// All targetted container resource endpoints wrap around containerParser to create a containerResource
	sp.Route("/containers/ps", "GET", a.nilParser, a.swarmContainers)
	sp.Route("/containers/json", "GET", a.nilParser, a.listContainers)
	sp.Route("/containers/create", "POST", a.containerCreateParser, a.createContainer)
	sp.Route("/containers/{name:.*}/export", "GET", a.containerViewOnly, a.swarmContainers)
	sp.Route("/containers/{name:.*}/changes", "GET", a.containerViewOnly, a.swarmContainers)
	sp.Route("/containers/{name:.*}/json", "GET", a.containerViewOnly, a.swarmContainersInfo)
	sp.Route("/containers/{name:.*}/top", "GET", a.containerViewOnly, a.swarmContainers)
	sp.Route("/containers/{name:.*}/logs", "GET", a.containerViewOnly, a.swarmHijack)
	sp.Route("/containers/{name:.*}/stats", "GET", a.containerViewOnly, a.swarmContainers)
	sp.Route("/containers/{name:.*}/archive", "GET", a.containerViewOnly, a.swarmContainers)
	sp.Route("/containers/{name:.*}/attach/ws", "GET", a.containerViewOnly, a.swarmHijack)
	sp.Route("/containers/{name:.*}/kill", "POST", a.containerRestrictedControl, a.orcaInstanceFilter)
	sp.Route("/containers/{name:.*}/pause", "POST", a.containerRestrictedControl, a.orcaInstanceFilter)
	sp.Route("/containers/{name:.*}/unpause", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/rename", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/restart", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/start", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/stop", "POST", a.containerRestrictedControl, a.orcaInstanceFilter)
	sp.Route("/containers/{name:.*}/wait", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/resize", "POST", a.containerRestrictedControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/attach", "POST", a.containerRestrictedControl, a.swarmHijack)
	sp.Route("/containers/{name:.*}", "DELETE", a.containerRestrictedControl, a.orcaInstanceFilter)
	sp.Route("/containers/{name:.*}/copy", "POST", a.containerFullControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/exec", "POST", a.containerFullControl, a.swarmContainers)
	sp.Route("/containers/{name:.*}/archive", "PUT", a.containerFullControl, a.swarmContainers)

	sp.Route("/exec/{name:.*}/json", "GET", a.nilParser, a.swarmRedirect)
	sp.Route("/exec/{name:.*}/start", "POST", a.nilParser, a.swarmHijack)
	sp.Route("/exec/{name:.*}/resize", "POST", a.nilParser, a.swarmRedirect)

	sp.Route("/auth", "POST", a.nilParser, a.swarmImages)
	sp.Route("/commit", "POST", a.nilParser, a.swarmImages)
	sp.Route("/build", "POST", a.nilParser, a.swarmImages)

	sp.Route("", "OPTIONS", a.nilParser, a.swarmRedirect)
}

// registerPublicRoutes defines all routes and handlers that are unauthorized
func (a *Api) RegisterPublicRoutes() {
	a.Router.HandleFunc("/ca", a.ca).Methods("GET")
	a.Router.HandleFunc("/_ping", a.ping).Methods("GET")
	a.Router.HandleFunc("/v{version:[0-9.]+}/_ping", a.ping).Methods("GET")
	a.Router.HandleFunc("/auth/login", a.login).Methods("POST")
	a.Router.HandleFunc("/auth/logout", a.logout).Methods("POST")
	a.Router.HandleFunc("/hub/webhook/{key}", a.hubWebhook).Methods("POST")

	// Node authorization is unauthorized as it requires a short-lived secret
	// The secret is maintained at a cluster-level in the KV store
	a.Router.HandleFunc("/api/nodes/authorize", a.authorizeNodeRequest).Methods("POST")

	// Node promotion is unauthorized as it uses an existing Classic Swarm channel
	// To transfer the root key material
	a.Router.HandleFunc("/api/nodes/promote", a.promoteNode).Methods("POST")

	a.Router.Handle("/exec", websocket.Handler(a.execContainer))
	a.Router.Handle("/containerlogs", websocket.Handler(a.containerLogs))

	// Warning, remote profiling is unauthenticated!
	if a.enableProfiling {
		log.Info("Registering pprof entrypoint at /debug/pprof/...")
		a.Router.HandleFunc("/debug/vars", expVars) // Borrowed from docker engine... do we need this?
		a.Router.HandleFunc("/debug/pprof/", pprof.Index)
		a.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		a.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		a.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		a.Router.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
		a.Router.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
		a.Router.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
		a.Router.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

		// TODO This will probably need tuning...
		runtime.SetBlockProfileRate(100) // 100% will slow us down so we will probably want to reduce this
	}

	// Register a handler to proxy requests to eNZi if the manager is using
	// an eNZi authenticator. The manager's authenticator needs to be
	// already setup using the eNZi backend.
	authenticator := a.manager.GetAuthenticator()
	if enziAuthenticator, ok := authenticator.(*enzi.Authenticator); ok {
		// Add proxy to the local eNZi API server:
		//     /enzi/ -> https://provider/enzi/
		a.Router.PathPrefix("/enzi/").Handler(makeEnziProxy(enziAuthenticator))
	}

	// Also register a handler for eNZi OpenID authentication keys.
	a.Router.HandleFunc("/openid_keys", a.openidKeys).Methods("GET")

	// Static Files. Don't add anything to a.Router after this
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("static/dist")))
}

func makeEnziProxy(enziAuthenticator *enzi.Authenticator) http.Handler {
	providerURL := &url.URL{
		Scheme: "https",
		Host:   enziAuthenticator.ProviderAddr(),
		Path:   "/",
	}
	tlsConfig := &tls.Config{
		RootCAs: enziAuthenticator.ProviderCAs(),
	}

	proxy := httputil.NewSingleHostReverseProxy(providerURL)
	proxy.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		// The default is 2 which is too small. We may need to adjust
		// this value as we get results from load/stress tests.
		MaxIdleConnsPerHost: 5,
	}

	return proxy
}
