package deploy

import (
	"fmt"
	"strings"
)

const (
	/* Timeout in seconds */
	ContainerRestartTimeout = 10
)

const (
	Phase2EnvVar                 = "DTR_PHASE2"
	ReplicaIDEnvVar              = "DTR_REPLICA_ID"
	PProfEnvVar                  = "DTR_PPROF_ENABLED"
	ComposeProjectName           = "dockertrustedregistry"
	DTRImagePrefix               = "dtr"
	DTRPrefixNoDash              = "dtr"
	DTRPrefix                    = DTRPrefixNoDash + "-"
	BridgeNetworkName            = "dtr-br"
	OverlayNetworkName           = "dtr-ol"
	DockerIndexURL               = "https://index.docker.io/v1/"
	DummyHubUserEmail            = "don.hcd@gmail.com" // just need an email that already has a docker hub account
	UnconfiguredCertSentinelCN   = "domain_name_required"
	BootstrapPhase2ContainerName = "dtr-phase2"
	BootstrapHelperContainerName = "dtr-helper"
	MixpanelToken                = "e273cb1bc7309af318b4c89f62c3ec7d"
)

const AllowUpgradesStartingFrom = "2.0.0"

const DockerVersionRequired = "1.10.0"
const UCPVersionRequired = "1.1.0"

const OverlayTestContainer1Name = "dtr-overlay-test-1"
const OverlayTestContainer2Name = "dtr-overlay-test-2"

const (
	TLSPEMFilename             = "server.pem"
	UCPCAPEMFilename           = "ucp.pem"
	EnziCAPEMFilename          = "enzi.pem"
	UserHubConfigFilename      = "hub.yml"
	NginxConfigFilename        = "nginx.conf"
	NotaryCertFilename         = "notary.pem"
	JobrunnerConfigFilename    = "jobrunner.yml"
	AuthBypassCAFilename       = "auth_bypass_ca.pem"
	RegistryConfigFilename     = "storage.yml"
	NotaryServerConfigFilename = "server_config.json"
	NotarySignerConfigFilename = "signer_config.json"
	AuthConfigFilename         = "garant.yml"
	LicenseConfigFilename      = "license.json"
	HAConfigFilename           = "ha.json"
	LoggingConfigFilename      = "logging.yml"
	HubCredentialsFilename     = ".dockercfg"
	EtcdDirectory              = "etcd"
	RethinkdbDirectory         = "rethink"
)

const ReplicaHealthPropertyPrefix = "health_"
const ReplicaHealthTimestampPropertyPrefix = "health_timestamp_"
const RethinkdbPort = uint16(28015)

const (
	RegistryRestartLockPath = "lock/registry-restart"
	ConfigDirPath           = "/config"
	EtcdPath                = "/dtr/configs"
	RegistryROStatePath     = "read-only-registry"
	RegistryGCStopPath      = "gc-stop-semaphore"
	GCRunLockPath           = "lock/gc-run"
	GCRunCountPath          = "gc-run-count" // see https://github.com/docker/dhe-deploy/issues/1537
	GeneratedConfigsDir     = "generatedConfigs"
	LogsCertPathInContainer = "/dtr-logging-certs"
	LogsCertPathInHost      = "/etc/docker/dtr-logging-certs"
	DataDirPath             = "/var/local/dtr"
	IndexDirName            = "/index"
	IndexPath               = DataDirPath + IndexDirName
	DatabasesDirName        = "/databases"
	DatabasesPath           = DataDirPath + DatabasesDirName
	EtcdDataPath            = DataDirPath + "/etcd"
	ImageStorageRootPath    = DataDirPath + "/image-storage"
	DockerSocketPath        = "/var/run/docker.sock"
)

const (
	DTRDBName          = "dtr2"
	JobrunnerDBName    = "jobrunner"
	NotaryServerDBName = "notaryserver"
	NotarySignerDBName = "notarysigner"
)

const (
	HasMigratedToTagstore = "hasMigratedToTagstore" // Whether the migration to tagstore has completed
	MigrationRepo         = "migrationRepo"         // The current repository that we're migrating
)

const (
	DefaultAuthBypassOU = "ucp"
)

const (
	DockerSecurityOpt = "label:type:docker_t"
)

const ReplicaIDLen = 12

const (
	StorageContainerPort uint16 = 5000
	RethinkAdminPort     uint16 = 8080
	NotaryServerHTTPPort uint16 = 4443
	NotarySignerGRPCPort uint16 = 7899
	NotarySignerHTTPPort uint16 = 4444
)

const (
	GarantRootCertFilename   = "token_roots.pem"
	GarantSigningKeyFilename = "signing_key.json"
	GarantSubroute           = "auth"
	GarantPort               = 80
	AdminSubroute            = "admin"
	EnziSubroute             = "enzi"
	RethinkSubroute          = "db"
	EventsEndpointSubroute   = "events"
	AdminPort                = 80  // XXX - this should be renamed
	AdminTlsPort             = 443 // XXX - this should be renamed
)

const (
	EnziSigningKeyFilename = "enzi_signing_key.json"
	EnziServiceFilename    = "enzi_service.json"
)

const (
	RegistryEventsHeaderName = "X-Registry-Events"
)

func IsProduction() bool {
	return ParseReleaseChannel(DefaultReleaseChannel).Namespace == "docker"
}

type repoName string

func (name repoName) Name() string {
	return fmt.Sprintf("%s/%s", DockerHubOrg, name)
}

func (name repoName) TaggedName() string {
	return fmt.Sprintf("%s:%s", name.Name(), Version)
}

var RepoNames = []repoName{
	BootstrapRepo,
	APIServerRepo,
	NginxRepo,
	RethinkRepo,
	RegistryRepo,
	NotaryServerRepo,
	NotarySignerRepo,
}

type ReleaseChannel struct {
	Namespace string
	Suffix    string
}

var defaultReleaseChannel = ReleaseChannel{
	Namespace: "docker",
}

// this is used for checking for updates, it needs to be updated
func (rc ReleaseChannel) ManagerRepoName() string {
	if rc.Suffix == "" {
		return rc.Namespace + "/" + DTRImagePrefix
	}
	return rc.Namespace + "/" + DTRImagePrefix + "-" + rc.Suffix
}

func ParseReleaseChannel(rc string) ReleaseChannel {
	namespace, suffix := "docker", "stable"
	channelParts := strings.Split(rc, "/")
	switch len(channelParts) {
	case 1:
		suffix = channelParts[0]
	case 2:
		namespace, suffix = channelParts[0], channelParts[1]
	}

	if (namespace == "" || suffix == "") ||
		(namespace == "docker" && suffix == "stable") {
		return defaultReleaseChannel
	}

	return ReleaseChannel{
		Namespace: namespace,
		Suffix:    suffix,
	}
}
