package config

import (
	"strings"

	"github.com/docker/orca/version"
)

var (
	SANNodeLabel = "com.docker.ucp.SANs"

	// TODO: duplicate entry entry until bootstrap/config is merged here
	OrcaReconcileContainerName = "ucp-reconcile"
	RuntimeContainerNames      = []string{
		OrcaControllerContainerName,
		OrcaSwarmManagerContainerName,
		OrcaSwarmJoinContainerName,
		OrcaKvContainerName,
		OrcaProxyContainerName,
		OrcaCAContainerName,
		OrcaSwarmCAContainerName,
		AuthStoreContainerName,
		AuthAPIContainerName,
		AuthWorkerContainerName,
		// The ephemeral containers are added to this list too so that
		// they are cleaned up during uninstall and fresh-install.
		AuthSyncDBContainerName,
		AuthCreateAdminContainerName,
		AuthDrainDBServerContainerName,
		OrcaKvBackupContainerName,
		OrcaKvRestoreContainerName,
	}

	// AgentContainerNames are the names of all containers expected to be
	// running on a non-manager node
	AgentContainerNames = []string{
		OrcaProxyContainerName,
		OrcaSwarmJoinContainerName,
	}

	// ManagerContainerNames are the names of all containers expected to be
	// running on a manager node
	ManagerContainerNames = append(AgentContainerNames,
		OrcaControllerContainerName,
		OrcaSwarmManagerContainerName,
		OrcaKvContainerName,
		OrcaCAContainerName,
		OrcaSwarmCAContainerName,
		AuthStoreContainerName,
		AuthAPIContainerName,
		AuthWorkerContainerName,
	)

	// Images by container name
	Images = map[string]string{
		OrcaControllerContainerName:    "docker/ucp-controller",
		OrcaSwarmManagerContainerName:  "docker/ucp-swarm",
		OrcaSwarmJoinContainerName:     "docker/ucp-swarm",
		OrcaKvContainerName:            "docker/ucp-etcd",
		OrcaKvBackupContainerName:      "docker/ucp-etcd",
		OrcaKvRestoreContainerName:     "docker/ucp-etcd",
		OrcaProxyContainerName:         "docker/ucp-agent",
		OrcaAgentContainerName:         "docker/ucp-agent",
		OrcaReconcileContainerName:     "docker/ucp-agent",
		OrcaCAContainerName:            "docker/ucp-cfssl",
		OrcaSwarmCAContainerName:       "docker/ucp-cfssl",
		AuthStoreContainerName:         "docker/ucp-auth-store",
		AuthAPIContainerName:           "docker/ucp-auth",
		AuthWorkerContainerName:        "docker/ucp-auth",
		AuthSyncDBContainerName:        "docker/ucp-auth",
		AuthCreateAdminContainerName:   "docker/ucp-auth",
		AuthDrainDBServerContainerName: "docker/ucp-auth",
		"addr": "docker/ucp-agent", // Just requires busybox's version of ip
		// Bogus, but to get the image pulled
		"dsinfo":  "docker/ucp-dsinfo",
		"compose": "docker/ucp-compose",
	}

	ImageVersion = strings.Split(version.FullVersion(), " ")[0]

	SwarmGRPCPort     = 2377
	OrcaPort          = 443
	SwarmPort         = 2376
	ProxyPort         = 12376
	KvPort            = 12379
	KvPortPeer        = 12380
	SwarmCAPort       = 12381
	OrcaCAPort        = 12382
	AuthStorePort     = 12383
	AuthStorePeerPort = 12384
	AuthAPIPort       = 12385
	AuthWorkerPort    = 12386

	// TODO rename this controller required ports...
	RequiredPorts = []*int{
		&OrcaPort,
		&SwarmPort,
		&ProxyPort,
		&KvPort,
		&KvPortPeer,
		&SwarmCAPort,
		&OrcaCAPort,
		&AuthStorePort,
		&AuthStorePeerPort,
		&AuthAPIPort,
		&AuthWorkerPort,
	}
)

const (
	// Container Names
	OrcaControllerContainerName   = "ucp-controller"
	OrcaSwarmManagerContainerName = "ucp-swarm-manager"
	OrcaSwarmJoinContainerName    = "ucp-swarm-join"
	OrcaKvContainerName           = "ucp-kv"
	OrcaProxyContainerName        = "ucp-proxy"
	OrcaAgentContainerName        = "ucp-agent"
	OrcaCAContainerName           = "ucp-client-root-ca"
	OrcaSwarmCAContainerName      = "ucp-cluster-root-ca"
	AuthStoreContainerName        = "ucp-auth-store"
	AuthAPIContainerName          = "ucp-auth-api"
	AuthWorkerContainerName       = "ucp-auth-worker"
	// These containers are ephemeral and only used for install, join,
	// upgrade, etc.
	AuthSyncDBContainerName        = "ucp-auth-sync-db"
	AuthCreateAdminContainerName   = "ucp-auth-create-admin"
	AuthDrainDBServerContainerName = "ucp-auth-drain-db-server"
	OrcaKvBackupContainerName      = "ucp-kv-backup"
	OrcaKvRestoreContainerName     = "ucp-kv-restore"
	BootstrapContainerName         = "ucp"
	BootstrapPhase2ContainerName   = "ucp-phase2"

	UCPLabelPrefix        = "com.docker.ucp"
	UCPInstanceIDLabelKey = UCPLabelPrefix + ".InstanceID"

	// Swarm-Mode Node Certs
	SwarmModeNodeCertDir      = "/var/lib/docker/swarm/certificates"
	SwarmModeNodeCAFilename   = "swarm-root-ca.crt"
	SwarmModeNodeKeyFilename  = "swarm-node.key"
	SwarmModeNodeCertFilename = "swarm-node.crt"
	SwarmModeNodeCAPath       = SwarmModeNodeCertDir + "/" + SwarmModeNodeCAFilename
	SwarmModeNodeCertPath     = SwarmModeNodeCertDir + "/" + SwarmModeNodeCertFilename
	SwarmModeNodeKeyPath      = SwarmModeNodeCertDir + "/" + SwarmModeNodeKeyFilename
)
