package manager

import (
	"crypto"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/docker/libtrust"
	"golang.org/x/net/context"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/dockerhub"
	"github.com/docker/orca/enzi/jose"
	orcatypes "github.com/docker/orca/types"
)

const (
	datastoreKey   = "orca"
	MixpanelToken  = "a5e078749467d17c4cb534a78d5f171e"
	MixpanelUrl    = "https://api.mixpanel.com"
	NodeHealthUp   = "up"
	NodeHealthDown = "down"
	RetryLimit     = 5 * time.Minute // How long we'll keep retrying on startup before giving up
)

var (
	ErrConnectionError                = errors.New("unable to connect to host")
	ErrAccessListDoesNotExist         = errors.New("access list does not exist")
	ErrInvalidAuthToken               = errors.New("invalid auth token")
	ErrExtensionDoesNotExist          = errors.New("extension does not exist")
	ErrWebhookKeyDoesNotExist         = errors.New("webhook key does not exist")
	ErrTokenDoesNotExist              = errors.New("token does not exist")
	ErrConsoleSessionDoesNotExist     = errors.New("console session does not exist")
	ErrContainerLogsTokenDoesNotExist = errors.New("containerlogs token does not exist")
	ErrRegistryDoesNotExist           = errors.New("registry does not exist")
	ErrRegistryExists                 = errors.New("registry already exists")
	ErrUnlicensed                     = errors.New("the system is unlicensed")
	authenticatorSecret               = "defaultorca"
	datastoreVersion                  = datastoreKey + "/v1"
	PeriodicInterval                  = time.Minute * 60 // TODO - make this configurable in the KV store
)

type (
	DefaultManagerConfig struct {
		SwarmClassicURL *url.URL
		EngineProxyURL  *url.URL
		DiscoveryAddr   string
		HostAddr        string
		Client          *client.Client
		ClientTransport *http.Transport
		ProxyClient     *client.Client
		ProxyTransport  *http.Transport
		Authenticator   auth.Authenticator
		SupportTimeout  int
		SupportImage    string
		DtrUrl          string
		DtrInsecure     bool
		DtrAdmin        string
		ClusterCAAddr   string
		ClientCAAddr    string

		// TLS Config for KV Store. It's used for pretty much
		// everything.
		DiscoveryTLSCaCertPath string
		DiscoveryTLSCertPath   string
		DiscoveryTLSKeyPath    string

		// TLS CA bundle for controller https server. Used for:
		// - generating docker client bundles
		// - authorizing joins of new nodes
		// - updating controller registration in KV store
		// The cert and key are used for signing/verifying tokens with
		// the builtin auth backend.
		ControllerCAPEM   []byte
		ControllerCertPEM []byte
		ControllerKeyPEM  []byte

		// TLS Config for swarm client and signing authentication JWSs
		// for the eNZi auth provider.
		SwarmCAPEM   []byte
		SwarmCertPEM []byte
		SwarmKeyPEM  []byte
	}

	DefaultManager struct {
		swarmClassicURL      *url.URL
		engineProxyURL       *url.URL
		datastore            kvstore.Store
		datastoreConfig      *kvstore.Config
		datastoreAddr        string
		hostAddr             string
		authenticator        *auth.Authenticator
		client               *client.Client
		clientTransport      *http.Transport
		proxyClient          *client.Client
		proxyTransport       *http.Transport
		disableUsageInfoCh   chan<- bool
		disableUsageInfo     *bool
		disableTracking      *bool
		anonymizeTracking    *bool
		trustKey             libtrust.PrivateKey
		swarmTLSConfig       *tls.Config
		httpClient           http.Client
		clusterCAChain       string
		clientCAChain        string
		supportTimeout       int
		supportImage         string
		configSubsystems     map[string]ConfigSubsystem
		configSubsystemCtors map[string]NewConfigSubsystem
		configStopCh         chan struct{}
		mutex                *sync.Mutex
		registry             orca.Registry
		selfManagerNode      *ManagerNode

		// The key and certificate chain used to produce authentication
		// tokens for the eNZi auth provider.
		enziTokenSigningKey *jose.PrivateKey
		enziTokenCertChain  []string

		// The controller key and cert are used for signing and
		// verifying session tokens produced by the legacy builtin
		// authenticator backend.
		controllerCertPEM []byte
		controllerKeyPEM  []byte
	}

	Manager interface {
		DockerClient() *client.Client
		DockerClientTransport() *http.Transport
		ProxyClient() *client.Client
		ProxyClientTransport() *http.Transport
		SwarmClassicURL() *url.URL
		EngineProxyURL() *url.URL
		GetAuthenticator() auth.Authenticator
		ID() string
		// conf
		GetTrackingDisabled() bool
		GetUsageInfoDisabled() bool
		AnonymizeTracking() bool
		TrackClientInfo(req *http.Request)
		GenerateClientBundle(ctx *auth.Context, host, publicKeyLabel string) ([]byte, error)
		AddTrustedCert(caCert string)
		EnableAdminUCPScheduling() bool
		EnableUserUCPScheduling() bool
		// accounts
		Accounts(ctx *auth.Context) ([]*auth.Account, error)
		Account(ctx *auth.Context, username string) (*auth.Account, error)
		AccountsForTeam(ctx *auth.Context, teamID string) ([]*auth.Account, error)
		AuthenticateUsernamePassword(username, password, remoteAddr string) (*auth.Context, error)
		AuthenticatePublicKey(key crypto.PublicKey, remoteAddr string) (*auth.Context, error)
		SaveAccount(ctx *auth.Context, updateAccount *auth.Account) (string, error)
		DeleteAccount(ctx *auth.Context, account *auth.Account) error
		ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error
		AuthSyncMessages(ctx *auth.Context) string
		AuthSync(ctx *auth.Context) string

		// access lists
		AccessLists() ([]*auth.AccessList, error)
		AccessListsForTeam(teamId string) ([]*auth.AccessList, error)
		AccessListsForAccount(ctx *auth.Context, username string) ([]*auth.AccessList, error)
		SaveAccessList(list *auth.AccessList) (string, error)
		AccessList(teamId, id string) (*auth.AccessList, error)
		RemoveAccessList(teamId, id string) error
		GetAccess(ctx *auth.Context) (map[string]auth.Role, error)

		// teams
		Team(ctx *auth.Context, id string) (*auth.Team, error)
		Teams(ctx *auth.Context) ([]*auth.Team, error)
		SaveTeam(ctx *auth.Context, team *auth.Team) (string, error)
		DeleteTeam(ctx *auth.Context, team *auth.Team) error
		AddMemberToTeam(ctx *auth.Context, teamID, username string) error
		RemoveMemberFromTeam(ctx *auth.Context, teamID, username string) error

		// token
		AuthenticateSessionToken(tokenStr string) (*auth.Context, error)
		Logout(tokenStr string) error
		// services
		Service(id string) (swarm.Service, error)
		ListServices(filters filters.Args) ([]swarm.Service, error)
		ListUserServices(ctx *auth.Context, filters filters.Args) ([]swarm.Service, error)
		// containers
		Container(id string) (types.ContainerJSON, error)
		ScaleContainer(id string, numInstances int) ScaleResult
		ListContainers(all bool, size bool, filters string) ([]types.Container, error)
		ListUserContainers(ctx *auth.Context, all bool, size bool, filters string) ([]types.Container, error)
		ContainerLogs(id string, opts types.ContainerLogsOptions) (io.ReadCloser, error)
		// networks
		Network(id string) (types.NetworkResource, error)
		ListUserNetworks(ctx *auth.Context, filters filters.Args) ([]types.NetworkResource, error)
		// images
		ListImages(all bool, f filters.Args) ([]types.Image, error)
		ListUserImages(ctx *auth.Context, all bool, f filters.Args) ([]types.Image, error)
		// events
		SaveEvent(event *orca.Event) error
		// webhook keys
		WebhookKey(key string) (*dockerhub.WebhookKey, error)
		WebhookKeys() ([]*dockerhub.WebhookKey, error)
		NewWebhookKey(image string) (*dockerhub.WebhookKey, error)
		SaveWebhookKey(key *dockerhub.WebhookKey) error
		DeleteWebhookKey(id string) error
		// nodes
		Nodes() ([]*orca.Node, error)
		Node(name string) (*orca.Node, error)
		ListNodes() ([]swarm.Node, error)
		InspectNode(nodeID string) (swarm.Node, error)
		SetControllerServerCerts(name string, req orcatypes.ServerCertRequest) error
		PromoteNode(name string, replicateCerts bool) error
		AuthorizeNodeRequest(r *orca.NodeRequest) (*orca.NodeConfiguration, error)
		// registry
		Registry(name string) (orca.Registry, error)
		// console sessions
		CreateConsoleSession(c *orca.ConsoleSession) error
		RemoveConsoleSession(c *orca.ConsoleSession) error
		ConsoleSession(token string) (*orca.ConsoleSession, error)
		ValidateConsoleSessionToken(containerId, token string) bool
		// catalog
		GetDefaultCatalog() ([]orca.CatalogItem, error)
		SearchCatalog(query string) ([]orca.CatalogItem, error)
		// applications
		Applications(ctx *auth.Context) ([]*orca.Application, error)
		Application(ctx *auth.Context, name string) (*orca.Application, error)
		RemoveApplication(app *orca.Application) error
		RestartApplication(app *orca.Application) error
		StopApplication(app *orca.Application) error
		ApplicationContainers(app *orca.Application) ([]types.Container, error)
		// support
		SupportDump(w io.Writer) error
		// datastore
		Datastore() kvstore.Store
		// High Availability
		GetStatus(node *ManagerNode) string
		GetSelfStatus() error
		GetManagers() []*ManagerNode
		// Configuration
		ListConfigSubsystems() []string
		GetSubsystemConfig(subsystemName string) (string, error)
		UserConfigUpdate(subsystemName, jsonConfig string) error
		// Licensing
		GetLicense() LicenseSubsystemConfig
		GetLicenseKeyID() string
		GetLicenseTier() string
		// websocket tokens
		CreateContainerLogsToken(c *orca.ContainerLogsToken) error
		RemoveContainerLogsToken(c *orca.ContainerLogsToken) error
		ContainerLogsToken(token string) (*orca.ContainerLogsToken, error)
		ValidateContainerLogsToken(containerId, token string) bool
		// notary
		RequireContentTrustForDTR() *bool
		RequireContentTrustForHub() *bool

		GetBanner(ctx *auth.Context) []Banner
	}

	ScaleResult struct {
		Scaled []string
		Errors []string
	}
)

func init() {
	consul.Register()
	etcd.Register()
	zookeeper.Register()
	boltdb.Register()
}

func getKVStore(addr string, options *kvstore.Config) (kvstore.Store, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	kvType := strings.ToLower(u.Scheme)
	kvHost := u.Host
	var backend kvstore.Backend

	switch kvType {
	case "consul":
		backend = kvstore.CONSUL
	case "etcd":
		backend = kvstore.ETCD
	case "zk":
		backend = kvstore.ZK
	case "boltdb":
		backend = kvstore.BOLTDB
	}

	kv, err := libkv.NewStore(
		backend,
		[]string{kvHost},
		options,
	)

	if err != nil {
		return nil, err
	}

	return kv, nil
}

func NewDefaultManager(config *DefaultManagerConfig) (Manager, error) {
	startTime := time.Now()
	var swarmTLSConfig *tls.Config

	clusterRootCA := string(config.SwarmCAPEM)
	clientRootCA := string(config.ControllerCAPEM)

	// init kv
	kvOpts := &kvstore.Config{
		ConnectionTimeout: time.Second * 10,
	}
	if config.DiscoveryTLSCaCertPath != "" && config.DiscoveryTLSCertPath != "" && config.DiscoveryTLSKeyPath != "" {
		tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
			CAFile:   config.DiscoveryTLSCaCertPath,
			CertFile: config.DiscoveryTLSCertPath,
			KeyFile:  config.DiscoveryTLSKeyPath,
		})
		if err != nil {
			return nil, err
		}
		swarmTLSConfig = tlsConfig

		log.Debug("Setting up KV with TLS")
		kvOpts.TLS = tlsConfig
	} else {
		log.Debug("Setting up KV without TLS")
	}

	// Note: this call doesn't immediately contact the KV
	kv, err := getKVStore(config.DiscoveryAddr, kvOpts)
	if err != nil {
		return nil, err
	}

	// Helper to retry things during startup
	keepTrying := func(retryMsg string, fp func() error) error {
		err := fp()
		for !time.Now().After(startTime.Add(RetryLimit)) && err != nil {
			err = fp()
			if err != nil {
				log.Infof("%s - %s", retryMsg, err)
				time.Sleep(2 * time.Second)
			}
		}
		if err != nil {
			log.Error("Unable to contact dependent services.  Giving up.")
		}
		return err
	}

	// This may hit errors if the KV store isn't responding (yet)
	var trustKey libtrust.PrivateKey
	if err := keepTrying("KV not responding", func() error {
		trustKey, err = LoadOrCreateTrustKey(kv)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// For the eNZi token signing key, we use our swarm client key/cert.
	keyPair, err := tls.X509KeyPair(config.SwarmCertPEM, config.SwarmKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("unable to parse swarm client key pair: %s", err)
	}

	enziTokenSigningKey, err := jose.NewPrivateKey(keyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load JOSE private key: %s", err)
	}

	// There is guaranteed to be at least one in the chain or else the
	// call to tls.X509KeyPair() above would have failed.
	enziTokenCertChain := make([]string, len(keyPair.Certificate))
	for i, certDerBytes := range keyPair.Certificate {
		enziTokenCertChain[i] = base64.StdEncoding.EncodeToString(certDerBytes)
	}

	// If we got this far, we know the KV store is up and available

	// Make sure swarm has finished coming up (might be electing a leader...)
	if err := keepTrying("Swarm manager not responding", func() error {
		_, err := config.Client.ServerVersion(context.TODO())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	swarmClient := http.Client{
		Timeout: InnerServiceTestTimeout, // Very short timeout for health check
		Transport: &http.Transport{
			// The Auth API server uses a server cert signed by the
			// cluster CA.
			TLSClientConfig: swarmTLSConfig,
		},
	}

	m := &DefaultManager{
		swarmClassicURL:      config.SwarmClassicURL,
		engineProxyURL:       config.EngineProxyURL,
		client:               config.Client,
		clientTransport:      config.ClientTransport,
		proxyClient:          config.ProxyClient,
		proxyTransport:       config.ProxyTransport,
		httpClient:           swarmClient,
		swarmTLSConfig:       swarmTLSConfig,
		clusterCAChain:       clusterRootCA,
		clientCAChain:        clientRootCA,
		supportTimeout:       config.SupportTimeout,
		supportImage:         config.SupportImage,
		trustKey:             trustKey,
		datastore:            kv,
		datastoreConfig:      kvOpts,
		datastoreAddr:        config.DiscoveryAddr,
		hostAddr:             config.HostAddr,
		configSubsystems:     make(map[string]ConfigSubsystem),
		configSubsystemCtors: make(map[string]NewConfigSubsystem),
		mutex:                &sync.Mutex{},
		selfManagerNode:      &ManagerNode{},
		enziTokenSigningKey:  enziTokenSigningKey,
		enziTokenCertChain:   enziTokenCertChain,
		controllerCertPEM:    config.ControllerCertPEM,
		controllerKeyPEM:     config.ControllerKeyPEM,
	}
	// These will panic if they can't get setup properly
	// TODO - refactor for better error handling
	setupAuthenticator(m)
	setupLogging(m)
	setupLicensing(m)
	setupRegistry(m)
	setupTracking(m)
	setupCA(m)
	setupScheduling(m)
	setupTrust(m)

	m.setupConfigWatcher()

	m.migrateLegacyCAConfig()

	m.verifySelfManagerPresent()
	m.startPeriodicTasks()

	return m, nil
}

func (m DefaultManager) Datastore() kvstore.Store {
	return m.datastore
}

func (m DefaultManager) ID() string {
	return m.trustKey.PublicKey().KeyID()
}

func (m DefaultManager) DockerClient() *client.Client {
	return m.client
}

func (m DefaultManager) DockerClientTransport() *http.Transport {
	return m.clientTransport
}

func (m DefaultManager) ProxyClient() *client.Client {
	return m.proxyClient
}

func (m DefaultManager) ProxyClientTransport() *http.Transport {
	return m.proxyTransport
}

func (m DefaultManager) SwarmClassicURL() *url.URL {
	return m.swarmClassicURL
}

func (m DefaultManager) EngineProxyURL() *url.URL {
	return m.engineProxyURL
}

func (m DefaultManager) createAdmin(username string, setPassword bool) error {
	if _, err := m.Account(nil, username); err == auth.ErrAccountDoesNotExist {
		// create roles
		var acct *auth.Account
		if username == "admin" && setPassword == true {
			acct = &auth.Account{
				Username:  "admin",
				Password:  "orca",
				FirstName: "Orca",
				LastName:  "Admin",
				Admin:     true,
			}
			log.Info("created admin user: username: admin password: orca")
		} else {
			acct = &auth.Account{
				Username: username,
				Admin:    true,
			}
			log.Infof("default admin user: %s", username)
		}
		ctx := &auth.Context{User: acct}
		if _, err := m.SaveAccount(ctx, acct); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (m DefaultManager) logEvent(eventType, message string, tags []string) {
	evt := &orca.Event{
		Type:    eventType,
		Time:    time.Now(),
		Message: message,
		Tags:    tags,
	}

	if err := m.SaveEvent(evt); err != nil {
		log.Errorf("error logging event: %s", err)
	}
}

func (m DefaultManager) RequireContentTrustForDTR() *bool {
	if configSubsystem, ok := m.configSubsystems[SingletonTrustKvKey].(TrustConfigSubsystem); ok {
		return &configSubsystem.cfg.RequireContentTrustForDTR
	}
	ret := false
	return &ret
}

func (m DefaultManager) RequireContentTrustForHub() *bool {
	if configSubsystem, ok := m.configSubsystems[SingletonTrustKvKey].(TrustConfigSubsystem); ok {
		return &configSubsystem.cfg.RequireContentTrustForHub
	}
	ret := false
	return &ret
}
