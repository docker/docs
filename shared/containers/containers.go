package containers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/shared/certstore"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
)

type Volume struct {
	Name     string
	Location string
}

func (v Volume) ReplicaName(replica string) string {
	return fmt.Sprintf("%s-%s", v.Name, replica)
}

func (v Volume) FormatForReplica(replica string) string {
	return fmt.Sprintf("%s:%s", v.ReplicaName(replica), v.Location)
}

var (
	EtcdVolume = Volume{
		"dtr-etcd",
		"/data",
	}
	RethinkVolume = Volume{
		"dtr-rethink",
		"/data",
	}
	RegistryVolume = Volume{
		"dtr-registry",
		"/storage",
	}
	CAVolume = Volume{
		"dtr-ca",
		"/ca",
	}
	NotaryVolume = Volume{
		"dtr-notary",
		"/notary",
	}
	Volumes = []Volume{
		EtcdVolume,
		RethinkVolume,
		RegistryVolume,
		CAVolume,
		NotaryVolume,
	}
)

var EtcdCACertStore = certstore.New(filepath.Join(CAVolume.Location, "etcd"))
var EtcdCertStore = certstore.New(filepath.Join(CAVolume.Location, "etcd-client"))
var RethinkCACertStore = certstore.New(filepath.Join(CAVolume.Location, "rethink"))
var RethinkCertStore = certstore.New(filepath.Join(CAVolume.Location, "rethink-client"))
var NotaryCACertStore = certstore.New(filepath.Join(CAVolume.Location, "notary"))
var NotaryCertStore = certstore.New(filepath.Join(CAVolume.Location, "notary-client"))
var NotarySignerStore = certstore.New(filepath.Join(NotaryVolume.Location, "notary-signer"))
var NotaryServerStore = certstore.New(filepath.Join(NotaryVolume.Location, "notary-server"))

type NetworkConfig struct {
	Name    string
	Aliases []string
}

type ContainerConfig struct {
	OpenStdin    bool
	StdinOnce    bool
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Environment  map[string]string
	Constraints  []string
	Entrypoint   strslice.StrSlice
	ExcludeNodes []string
	Image        string
	IPAddress    []string
	Name         string
	Networks     []NetworkConfig
	Node         string
	Ports        []string
	Restart      string
	Tty          bool
	Volumes      []Volume
	DumbVolumes  []string
	LogConfig    container.LogConfig
}

type DTRContainer struct {
	Name       string
	Image      string
	BaseConfig ContainerConfig
}

var (
	Etcd = DTRContainer{
		Name:  "etcd",
		Image: deploy.EtcdImageName,
		BaseConfig: ContainerConfig{
			Environment: map[string]string{
				"ETCD_DATA_DIR": "/data/dtr.etcd",

				"ETCD_LISTEN_CLIENT_URLS": fmt.Sprintf("https://0.0.0.0:%d,https://0.0.0.0:%d", EtcdClientPort1, EtcdClientPort2),

				"ETCD_LISTEN_PEER_URLS":      fmt.Sprintf("https://0.0.0.0:%d", EtcdPeerPort1),
				"ETCD_INITIAL_CLUSTER_TOKEN": "etcd-cluster-1",

				"ETCD_CLIENT_CERT_AUTH": "true",
				"ETCD_TRUSTED_CA_FILE":  EtcdCACertStore.CertPath(),
				"ETCD_CERT_FILE":        EtcdCertStore.CertPath(),
				"ETCD_KEY_FILE":         EtcdCertStore.KeyPath(),

				"ETCD_PEER_CLIENT_CERT_AUTH": "true",
				"ETCD_PEER_TRUSTED_CA_FILE":  EtcdCACertStore.CertPath(),
				"ETCD_PEER_CERT_FILE":        EtcdCertStore.CertPath(),
				"ETCD_PEER_KEY_FILE":         EtcdCertStore.KeyPath(),
			},
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}, {Name: deploy.OverlayNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				EtcdVolume,
				CAVolume,
			},
		},
	}
	Registry = DTRContainer{
		Name:  "registry",
		Image: deploy.RegistryRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				RegistryVolume,
				CAVolume,
			},
		},
	}
	Jobrunner = DTRContainer{
		Name:  "jobrunner",
		Image: deploy.JobrunnerRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				CAVolume,
				RegistryVolume,
			},
		},
	}
	LoadBalancer = DTRContainer{
		Name:  "nginx",
		Image: deploy.NginxRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				CAVolume,
			},
		},
	}
	APIServer = DTRContainer{
		Name:  "api",
		Image: deploy.APIServerRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				// this is needed for filesystem storage driver GC for now
				RegistryVolume,
				CAVolume,
			},
		},
	}
	Rethinkdb = DTRContainer{
		Name:  "rethinkdb",
		Image: deploy.RethinkRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Networks:    []NetworkConfig{{Name: deploy.BridgeNetworkName}, {Name: deploy.OverlayNetworkName}},
			Restart:     "unless-stopped",
			Environment: map[string]string{},
			Volumes: []Volume{
				RethinkVolume,
				CAVolume,
			},
			Entrypoint: []string{
				"/start.sh",
			},
		},
	}
	NotaryServer = DTRContainer{
		Name:  "notary-server",
		Image: deploy.NotaryServerRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Environment: map[string]string{},
			Networks:    []NetworkConfig{{Name: deploy.BridgeNetworkName}, {Name: deploy.OverlayNetworkName}},
			Restart:     "unless-stopped",
			Volumes: []Volume{
				CAVolume,
				NotaryVolume,
			},
		},
	}
	NotarySigner = DTRContainer{
		Name:  "notary-signer",
		Image: deploy.NotarySignerRepo.TaggedName(),
		BaseConfig: ContainerConfig{
			Environment: map[string]string{
				// TODO: these should be configurable
				"NOTARY_SIGNER_DEFAULT_ALIAS": "timestamp_1",
				"NOTARY_SIGNER_TIMESTAMP_1":   "testpassword",
			},
			Networks: []NetworkConfig{{Name: deploy.BridgeNetworkName}, {Name: deploy.OverlayNetworkName}},
			Restart:  "unless-stopped",
			Volumes: []Volume{
				CAVolume,
				NotaryVolume,
			},
		},
	}
)

func (c DTRContainer) FullName() string {
	return fmt.Sprintf("%s%s", deploy.DTRPrefix, c.Name)
}

func (c DTRContainer) ReplicaName(replicaID string) string {
	return fmt.Sprintf("%s%s-%s", deploy.DTRPrefix, c.Name, replicaID)
}

// rethinkdb server tags don't allow dashes
func (c DTRContainer) RethinkServerTagName(replicaID string) string {
	return fmt.Sprintf("%s_%s_%s", deploy.DTRPrefixNoDash, c.Name, replicaID)
}

func (c DTRContainer) BridgeName(replicaID string) string {
	return fmt.Sprintf("%s%s-%s.%s", deploy.DTRPrefix, c.Name, replicaID, deploy.BridgeNetworkName)
}

func (c DTRContainer) BridgeNameLocalReplica() string {
	return fmt.Sprintf("%s%s-%s.%s", deploy.DTRPrefix, c.Name, os.Getenv(deploy.ReplicaIDEnvVar), deploy.BridgeNetworkName)
}

func (c DTRContainer) OverlayName(replicaID string) string {
	return fmt.Sprintf("%s%s-%s.%s", deploy.DTRPrefix, c.Name, replicaID, deploy.OverlayNetworkName)
}

// TODO: add cluster name to params?
func (c DTRContainer) ContainerConfig(cmdName, replicaID string, allReplicaIDs []string, nodeName string, haConfig *hubconfig.HAConfig) ContainerConfig {
	newConfig := c.BaseConfig
	if newConfig.DumbVolumes == nil {
		newConfig.DumbVolumes = []string{}
	}
	if newConfig.Environment == nil {
		newConfig.Environment = map[string]string{}
	}

	// dynamic configs per container
	switch c.Name {
	case APIServer.Name:
		if !deploy.IsProduction() && deploy.BindAssets {
			newConfig.DumbVolumes = append(newConfig.DumbVolumes, deploy.CurrentDirectory+"/adminserver/ui/src:/ui")
			logrus.Info("deploying in bind-assets mode")
		}
	case LoadBalancer.Name:
		if newConfig.Ports == nil {
			newConfig.Ports = []string{}
		}
		newConfig.Ports = append(newConfig.Ports, fmt.Sprintf("80/tcp:%d", haConfig.ReplicaConfig[replicaID].HTTPPort))
		newConfig.Ports = append(newConfig.Ports, fmt.Sprintf("443/tcp:%d", haConfig.ReplicaConfig[replicaID].HTTPSPort))
	case Etcd.Name:
		if cmdName == "install" {
			newConfig.Environment["ETCD_INITIAL_CLUSTER_STATE"] = "new"
		} else {
			newConfig.Environment["ETCD_INITIAL_CLUSTER_STATE"] = "existing"
		}

		// This is here because it's a bit of a chicken and egg problem
		newConfig.Environment["ETCD_ADVERTISE_CLIENT_URLS"] = fmt.Sprintf("https://%s:%d,https://%s:%d", Etcd.BridgeName(replicaID), EtcdClientPort1, Etcd.BridgeName(replicaID), EtcdClientPort2)
		newConfig.Environment["ETCD_NAME"] = c.ReplicaName(replicaID)
		newConfig.Environment["ETCD_INITIAL_ADVERTISE_PEER_URLS"] = fmt.Sprintf("https://%s:%d", c.OverlayName(replicaID), EtcdPeerPort1)

		if haConfig.EtcdHeartbeatInterval != 0 {
			newConfig.Environment["ETCD_HEARTBEAT_INTERVAL"] = fmt.Sprintf("%d", haConfig.EtcdHeartbeatInterval)
		}
		if haConfig.EtcdElectionTimeout != 0 {
			newConfig.Environment["ETCD_ELECTION_TIMEOUT"] = fmt.Sprintf("%d", haConfig.EtcdElectionTimeout)
		}
		if haConfig.EtcdSnapshotCount != 0 {
			newConfig.Environment["ETCD_SNAPSHOT_COUNT"] = fmt.Sprintf("%d", haConfig.EtcdSnapshotCount)
		}
	}

	newConfig.Environment[deploy.ReplicaIDEnvVar] = replicaID

	newConfig.Environment[deploy.PProfEnvVar] = strconv.FormatBool(haConfig.EnablePProf)

	if haConfig.LogHost != "" && haConfig.LogProtocol != "internal" {
		var logConfig = map[string]string{
			"syslog-address": fmt.Sprintf("%s://%s", haConfig.LogProtocol, haConfig.LogHost),
			// https://github.com/docker/docker/blob/4b98193beab00bc6cf48762858570a1bd418c9ef/docs/reference/logging/log_tags.md#log-tags
			"tag": "{{.Name}}",
		}

		if haConfig.LogProtocol == "tcp+tls" {
			if haConfig.LogTLSCACert != "" {
				logConfig["syslog-tls-ca-cert"] = filepath.Join(deploy.LogsCertPathInHost, "ca.pem")
			}
			if haConfig.LogTLSCert != "" {
				logConfig["syslog-tls-cert"] = filepath.Join(deploy.LogsCertPathInHost, "cert.pem")
			}
			if haConfig.LogTLSKey != "" {
				logConfig["syslog-tls-key"] = filepath.Join(deploy.LogsCertPathInHost, "key.pem")
			}

			logConfig["syslog-tls-skip-verify"] = strconv.FormatBool(haConfig.LogTLSSkipVerify)
		}

		newConfig.LogConfig = container.LogConfig{
			Type:   "syslog",
			Config: logConfig,
		}
	} else {
		newConfig.LogConfig = container.LogConfig{Type: "json-file"}
	}

	// Add all container names to the common noproxy environment variable
	noproxy := []string{}
	if haConfig.NoProxy != "" {
		noproxy = append(noproxy, haConfig.NoProxy)
	}
	for _, container := range AllContainers {
		// Golang uses "." prefixes for domain names in the no_proxy environment variable to match subdomains and domains
		//   .foo.com matches bar.foo.com and foo.com.
		//    foo.com matches bar.foo.com
		// so add em both
		noproxy = append(noproxy, container.BridgeName(replicaID), "."+container.BridgeName(replicaID))
		noproxy = append(noproxy, container.ReplicaName(replicaID), "."+container.ReplicaName(replicaID))
		noproxy = append(noproxy, container.OverlayName(replicaID), "."+container.OverlayName(replicaID))
	}
	ucpHost := strings.Split(haConfig.UCPHost, ":")[0]
	if ucpHost != "" {
		noproxy = append(noproxy, ucpHost, "."+ucpHost)
	}
	enziHost := strings.Split(haConfig.EnziHost, ":")[0]
	if enziHost != "" {
		noproxy = append(noproxy, enziHost, "."+enziHost)
	}
	// NOTE: these must always be set in order for us to be able to diff correctly
	// whey their values are empty
	newConfig.Environment["NO_PROXY"] = strings.Join(noproxy, ", ")
	newConfig.Environment["HTTP_PROXY"] = haConfig.HTTPProxy
	newConfig.Environment["HTTPS_PROXY"] = haConfig.HTTPSProxy

	newConfig.Name = c.ReplicaName(replicaID)
	newConfig.Image = c.Image
	newConfig.Node = nodeName
	return newConfig
}

func ContainerConfigs(cmdName, replicaID string, allReplicaIDs []string, nodeName string, haConfig *hubconfig.HAConfig, dtrContainers []DTRContainer) []ContainerConfig {
	var configs []ContainerConfig

	if len(dtrContainers) == 0 {
		dtrContainers = AllContainers
	}

	for _, container := range dtrContainers {
		configs = append(configs, container.ContainerConfig(cmdName, replicaID, allReplicaIDs, nodeName, haConfig))
	}
	return configs
}

var AllContainers = []DTRContainer{Etcd, Rethinkdb, Registry, APIServer, NotaryServer, LoadBalancer, Jobrunner, NotarySigner}

const (
	EtcdClientPort1 = 2379
	EtcdClientPort2 = 4001
	EtcdPeerPort1   = 2380
)

var EtcdUrls = func() []string {
	return []string{
		fmt.Sprintf("%s:%d", Etcd.BridgeName(os.Getenv(deploy.ReplicaIDEnvVar)), EtcdClientPort1),
		fmt.Sprintf("%s:%d", Etcd.BridgeName(os.Getenv(deploy.ReplicaIDEnvVar)), EtcdClientPort2),
	}
}
