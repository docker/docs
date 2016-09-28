package adminserver

import (
	"encoding/json"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/docker/dhe-deploy"
	apiserver "github.com/docker/dhe-deploy/adminserver/api/server"
	"github.com/docker/dhe-deploy/adminserver/healthcheck"
	"github.com/docker/dhe-deploy/garant/authz"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/jobs"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/manager/schema/interfaces"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/mailgun/oxy/forward"

	log "github.com/Sirupsen/logrus"
	"github.com/benbjohnson/clock"
	"github.com/codegangsta/negroni"
	"github.com/docker/distribution/context"
	garantauth "github.com/docker/garant/auth"
	gorillacontext "github.com/gorilla/context"
	mixpanel "github.com/pdevine/go-mixpanel"
	rethink "gopkg.in/dancannon/gorethink.v2"
	"gopkg.in/tylerb/graceful.v1"
)

var TemplatesDir = "/ui"

type AdminServer struct {
	settingsStore           hubconfig.SettingsStore
	kvStore                 hubconfig.KeyValueStore
	licenseChecker          hubconfig.LicenseChecker
	versionChecker          versions.Checker
	managerLock             sync.Mutex
	etcdClient              *http.Client
	notaryClient            *http.Client
	syslogger               SyslogWriter
	rethinkSession          *rethink.Session
	eventManager            schema.EventManager
	eventWebSocketManager   *webSocketManager
	repositoryManager       *schema.RepositoryManager
	repositoryAccessManager *schema.RepositoryAccessManager
	propertyManager         interfaces.PropertyManager
	authorizer              authz.Authorizer
	enziFwd                 *forward.Forwarder
	apiServer               *apiserver.APIServer
	jobRunner               jobs.Runner
	mixpanelClient          *mixpanel.Mixpanel
	*alerts
}

func New(settingsStore hubconfig.SettingsStore, kvStore hubconfig.KeyValueStore, etcdClient, notaryClient *http.Client, licenseChecker hubconfig.LicenseChecker, versionChecker versions.Checker, syslogger SyslogWriter, authorizer authz.Authorizer, session *rethink.Session) (*AdminServer, error) {
	alerts := &alerts{settingsStore: settingsStore, licenseChecker: licenseChecker, kvStore: kvStore}

	eventManager := schema.NewEventManager(session)
	eventWebSocketManager := NewWebSocketManager()
	repositoryManager := schema.NewRepositoryManager(session)
	repositoryAccessManager := schema.NewRepositoryAccessManager(session)
	propertyManager := schema.NewRethinkPropertyManager(session)

	healthCheckJob, err := healthcheck.NewJob(settingsStore, kvStore, etcdClient, notaryClient, propertyManager, session, log.StandardLogger())
	if err != nil {
		return nil, err
	}

	jobRunner := jobs.NewRunner(clock.New())
	jobRunner.AddJob("healthCheck", healthCheckJob)
	_ = jobRunner.RunNow("healthCheck")
	backgroundContext := context.Background()
	contextLogger := log.New()
	contextLogger.Formatter = new(log.JSONFormatter)
	context.WithLogger(backgroundContext, contextLogger)

	as := &AdminServer{
		settingsStore:           settingsStore,
		kvStore:                 kvStore,
		licenseChecker:          licenseChecker,
		versionChecker:          versionChecker,
		managerLock:             sync.Mutex{},
		alerts:                  alerts,
		etcdClient:              etcdClient,
		notaryClient:            notaryClient,
		syslogger:               syslogger,
		rethinkSession:          session,
		eventManager:            eventManager,
		eventWebSocketManager:   eventWebSocketManager,
		repositoryManager:       repositoryManager,
		repositoryAccessManager: repositoryAccessManager,
		propertyManager:         propertyManager,
		authorizer:              authorizer,
		apiServer:               apiserver.NewAPIServer(backgroundContext, authorizer, settingsStore, kvStore, propertyManager, session),
		jobRunner:               jobRunner,
	}

	as.setupTracking()

	return as, nil
}

func (a *AdminServer) SubscribeToEvents() {
	// register eventWebSocketManager as a listener for db events
	_ = a.eventManager.Subscribe(func(e schema.Event) {
		permissionFilter := func(c *webSocketClient, data interface{}) bool {
			evt, ok := data.(schema.Event)
			if !ok || c.user == nil || c.user.Account == nil {
				return false
			}

			// if the user is an admin
			if (c.user.Account.IsAdmin != nil && *c.user.Account.IsAdmin) ||
				// if the user committed the action
				c.user.Account.ID == evt.Actor {
				return true
			}

			return false
		}
		a.eventWebSocketManager.FilteredPublish(e, permissionFilter)
	})
	go a.eventWebSocketManager.Listen()
}

func (a *AdminServer) jobStatus(writer http.ResponseWriter, request *http.Request) {
	jobStatus, err := a.jobRunner.Status()
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}
	writeJSON(writer, jobStatus)
}

func (a *AdminServer) runJob(writer http.ResponseWriter, request *http.Request) {
	options := struct {
		Job string `json:"job"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&options); err != nil {
		log.WithField("error", err).Warn("Failed to decode job options")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}
	if err := a.jobRunner.RunNow(options.Job); err != nil {
		status := http.StatusInternalServerError
		if err == jobs.ErrJobAlreadyRunning || err == interfaces.ErrPropertyNotSet {
			status = http.StatusConflict
		} else if err == jobs.ErrJobNotFound {
			status = http.StatusNotFound
		}
		writeJSONError(writer, err, status)
		return
	}
	writeJSONStatus(writer, nil, http.StatusAccepted)
}

// authenticateRequest attempts to authenticate the client using either a
// session cookie or using basic authentication (or any other authenticaiton
// schema supported by pluggable backends). If authentication is successful,
// a "user" value is set on the request using gorilla/context. If
// authenticaiton fails due to invalid credentials, the error is ignored. A
// non-nil error is only returned if there is an internal server error.
func (a *AdminServer) authenticateRequest(writer http.ResponseWriter, request *http.Request) error {
	// The authorizer is configured to authenticate using cookies or basic
	// authentication. Any error returned is either an internal error OR
	// an authentication challenge error.
	user, err := a.authorizer.AuthenticateRequestUser(context.Background(), request)
	if err != nil {
		if _, ok := err.(garantauth.Challenge); ok && user.IsAnonymous {
			// Failed authentication attempt or invalid creds.
			// Ignore the authentication failure and set the user
			// to an unauthenticated user
			gorillacontext.Set(request, "user", user)
			return nil
		}

		// A legitimate internal error.
		log.WithField("error", err).Info("Unable to authenticate")
		return err
	}

	gorillacontext.Set(request, "user", user)
	return nil
}

func (a *AdminServer) wrapHandler(handler http.Handler) http.Handler {
	n := negroni.New(
		a.clearContext(),
		a.recoveryMiddleware(),
		a.authMiddleware(),
		a.noCacheMiddleware(),
	)

	n.UseHandler(handler)
	return n
}

func (a *AdminServer) ListenAndServe(addr string) {
	router := a.buildRouter()
	trackedDirectories := []string{"/"}
	registryConfig, err := a.settingsStore.RegistryConfig()
	if registryConfig.Storage.Type() == "filesystem" {
		userStoragePath, _ := registryConfig.Storage.Parameters()["rootdirectory"].(string)
		storageDir := path.Join(deploy.ImageStorageRootPath, userStoragePath)
		trackedDirectories = append(trackedDirectories, storageDir)
	}

	if err = a.jobRunner.Start(); err != nil {
		log.WithField("error", err).Error("Failed to start job system")
	}

	handler := a.wrapHandler(router)

	gracefulServer := graceful.Server{
		Timeout: 2 * time.Second,
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	log.WithField("address", addr).Info("Admin server about to listen for connections")
	err = gracefulServer.ListenAndServe()

	select {
	case <-gracefulServer.StopChan():
		a.jobRunner.Stop()
		time.Sleep(30 * time.Second)
		log.Warn("Expected process to exit!")
	default:
		if err != nil {
			log.WithField("error", err).Error("Admin server exited with error")
		} else {
			log.Warn("Admin server exited with no error")
		}
	}
}
