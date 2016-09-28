// Package config holds the global bootstrap configuration
package config

// XXX Should we read this in from some config file?  Kinda nice having it "compiled in" though...

import (
	"crypto/tls"
	"path/filepath"
	"syscall"
	"time"
)

type InteractivePromptData struct {
	Echo   bool
	Prompt string
}

var (
	MinVersion       = "1.12.0"
	MinKernelVersion = "3.10" // We only warn if it's too old
	MinMemoryGB      = float64(2.0)
	MinStorageGB     = float64(3.0) // TODO - Get some realistic numbers before GA

	DockerSock        = "unix:///var/run/docker.sock"
	AliveCheckTimeout = 60 * time.Second
	// AuthStoreAliveCheckTimeout is 5 times longer than normal for the
	// due to a DNS timeout issue on some hosts/distros.
	AuthStoreAliveCheckTimeout               = 5 * 60 * time.Second
	DockerTimeout                            = 30 * time.Second
	DockerTls                  (*tls.Config) = nil

	// Unique ID for this orca instance (generated in phase 1, fed to phase 2)
	OrcaInstanceID  = ""
	OrcaInstanceKey = ""
	OrcaLabelPrefix = "com.docker.ucp"

	// TODO - container names were here.

	OrcaSANs  = []string{}
	DNS       = []string{}
	DNSOpt    = []string{}
	DNSSearch = []string{}

	// KvPath is the location in the KV store where both engine and swarm discovery take place
	KvPath = "/docker/nodes"

	Phase2VolumeMounts = []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		"/var/lib/docker/swarm:/var/lib/docker/swarm",
	}

	PullBehavior = ""
	KvType       = "etcd"

	// Volumes
	OrcaRootCAVolumeName          = "ucp-client-root-ca"
	SwarmRootCAVolumeName         = "ucp-cluster-root-ca"
	OrcaServerCertVolumeName      = "ucp-controller-server-certs"
	SwarmNodeCertVolumeName       = "ucp-node-certs"
	SwarmKvCertVolumeName         = "ucp-kv-certs"
	SwarmControllerCertVolumeName = "ucp-controller-client-certs"
	OrcaKVVolumeName              = "ucp-kv"
	OrcaKVBackupVolumeName        = "ucp-kv-backup"
	AuthStoreCertsVolumeName      = "ucp-auth-store-certs"
	AuthAPICertsVolumeName        = "ucp-auth-api-certs"
	AuthWorkerCertsVolumeName     = "ucp-auth-worker-certs"
	AuthStoreDataVolumeName       = "ucp-auth-store-data"
	AuthWorkerDataVolumeName      = "ucp-auth-worker-data"

	AllVolumesMap = map[string]interface{}{
		OrcaRootCAVolumeName:          struct{}{},
		SwarmRootCAVolumeName:         struct{}{},
		OrcaServerCertVolumeName:      struct{}{},
		SwarmNodeCertVolumeName:       struct{}{},
		SwarmKvCertVolumeName:         struct{}{},
		SwarmControllerCertVolumeName: struct{}{},
		OrcaKVVolumeName:              struct{}{},
		AuthStoreCertsVolumeName:      struct{}{},
		AuthAPICertsVolumeName:        struct{}{},
		AuthWorkerCertsVolumeName:     struct{}{},
		AuthStoreDataVolumeName:       struct{}{},
		AuthWorkerDataVolumeName:      struct{}{},
		OrcaKVBackupVolumeName:        struct{}{},
	}

	// Mounts within the bootstrapper to keep them isolated
	OrcaRootCAVolumeMount          = filepath.Join(Phase2VolMountDir, OrcaRootCAVolumeName)
	SwarmRootCAVolumeMount         = filepath.Join(Phase2VolMountDir, SwarmRootCAVolumeName)
	OrcaServerCertVolumeMount      = filepath.Join(Phase2VolMountDir, OrcaServerCertVolumeName)
	SwarmNodeCertVolumeMount       = filepath.Join(Phase2VolMountDir, SwarmNodeCertVolumeName)
	SwarmKvCertVolumeMount         = filepath.Join(Phase2VolMountDir, SwarmKvCertVolumeName)
	SwarmControllerCertVolumeMount = filepath.Join(Phase2VolMountDir, SwarmControllerCertVolumeName)
	OrcaKVVolumeMount              = filepath.Join(Phase2VolMountDir, OrcaKVVolumeName)
	AuthStoreCertsVolumeMount      = filepath.Join(Phase2VolMountDir, AuthStoreCertsVolumeName)
	AuthAPICertsVolumeMount        = filepath.Join(Phase2VolMountDir, AuthAPICertsVolumeName)
	AuthWorkerCertsVolumeMount     = filepath.Join(Phase2VolMountDir, AuthWorkerCertsVolumeName)
	AuthStoreDataVolumeMount       = filepath.Join(Phase2VolMountDir, AuthStoreDataVolumeName)
	AuthWorkerDataVolumeMount      = filepath.Join(Phase2VolMountDir, AuthWorkerDataVolumeName)
	EngineLibDir                   = "/var/lib/docker"

	EngineConfigDir          = "/etc/docker/"
	EngineConfigFile         = EngineConfigDir + "daemon.json"
	EngineConfigReloadSignal = syscall.SIGHUP
	EnginePidDir             = "/var/run/"

	InteractiveArgs = map[string]InteractivePromptData{
		"REGISTRY_USERNAME":  {true, "Please enter your Docker Hub username"},
		"REGISTRY_PASSWORD":  {false, "Please enter your Docker Hub password"},
		"REGISTRY_EMAIL":     {true, "Please enter your Docker Hub e-mail address"},
		"UCP_URL":            {true, "Please enter the URL to your UCP server"},
		"UCP_ADMIN_USER":     {true, "Please enter your UCP Admin username"},
		"UCP_ADMIN_PASSWORD": {false, "Please enter your UCP Admin password"},
	}

	InPhase2 = false

	// Where volumes get mounted in phase 2
	Phase2VolMountDir = "/var/lib/docker/ucp"

	// Hostname related settings, discovered, or provided by user
	OrcaHostAddress = ""         // The primary hostname/IP for this node
	OrcaLocalName   = ""         // The Name of the node (extracted from dockerclient.Info by deafult)
	OrcaHostnames   = []string{} // Used for certificates

	// Certificate settings
	CertDir             = "/etc/docker/ssl"
	CAFilename          = "ca.pem"
	CertFilename        = "cert.pem"
	KeyFilename         = "key.pem"
	ControllerSwarmCACN = "UCP Cluster Root CA"
	ControllerOrcaCACN  = "UCP Client Root CA"
	ExternalOrcaCA      = false
	KeyAlgo             = "rsa"
	KeySize             = 4096

	// Swarm experimental features
	SwarmExperimental = false

	// Swarm scheduler choice (default if unspec'd is spread)
	SwarmBinpack = false
	SwarmRandom  = false

	// User injection of initial license
	LicenseFile = "/docker_subscription.lic"
	BackupFile  = "/backup.tar"

	// KV install-time configuration
	KVTimeout = 1000 // ms

	// Reconcile Mount Paths
	ReconcileOrcaKeyFileMount  = SwarmRootCAVolumeMount + "/ucp-instance-key.pem"
	ReconcileOrcaCredFileMount = SwarmRootCAVolumeMount + "/ucp-credentials"
)
