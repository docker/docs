package bootstrap

import (
	"net/http"

	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/shared/containers"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
)

type (
	Bootstrap interface {
		ImageExists(imageID string) (bool, error)
		PullImage(imageID string) error
		PullImageWithChecks(imageID string) error
		PullImages() error
		ValidateOrGetNewReplicaID(string) (string, error)
		GetNewReplicaID() (string, error)
		ExistingReplicaFlagPicker(string, bool) (*Replica, error)
		ExistingReplicaPicker(string, string, bool) (*Replica, error)
		ConfirmNodeIsOkay() error
		CreateNetwork(name, driver string) (types.NetworkCreateResponse, error)
		CreateNodeConstraintsFromPort(publicPorts []int) (*[]string, error)
		NetworkConnect(networkID, containerID string, endpointCfg *network.EndpointSettings) error
		NetworkInspect(netID string) (types.NetworkResource, error)
		NetworkExists(networkID string) bool
		NetworkID(name string) (string, error)
		ContainerCreate(cfg *container.Config, hostCfg *container.HostConfig, netCfg *network.NetworkingConfig, name string) (types.ContainerCreateResponse, error)
		// TODO: make ContainerConfig a pointer
		ContainerCreateFromContainerConfig(containers.ContainerConfig) (*types.ContainerCreateResponse, error)
		ContainerInspect(containerID string) (types.ContainerJSON, error)
		ContainerKill(containerID, signal string) error
		ContainerList(options types.ContainerListOptions) ([]types.Container, error)
		ContainerNode(containerID string) (string, error)
		ContainerRemove(id string, options types.ContainerRemoveOptions) error
		// TODO: make ContainerConfig a pointer
		ContainerRunFromContainerConfig(containers.ContainerConfig) (string, error)
		ContainerStart(containerID string) error
		ContainerWait(containerID string) (int, error)
		ContainerStop(containerID string) error
		StartDTRContainers(ignoreContainers map[string]bool, replicaID string) (int, error)
		StopDTRContainers(ignoreContainers map[string]bool, replicaID string) error
		RemoveDTRContainers(ignoreContainers map[string]bool, replicaID string) error
		RemoveDTRVolumes(string, bool) error
		RunContainer(cfg *container.Config, hostCfg *container.HostConfig, netCfg *network.NetworkingConfig, name string) (string, error)
		ContainerDiff(container containers.ContainerConfig) (bool, error)
		// returns container name->id map
		GetRunningDTRContainers(replicaID string) (map[string]string, error)
		VolumeCreate(name string) error
		VolumeExists(volumeID string) bool
		VolumeRemove(volumeID string) error
		ListReplicaIDs(bool) ([]string, error)
		ListReplicas(bool) (Replicas, error)
		GetHost() string
		GetUsername() string
		GetPassword() string
		GetTlsVerify() bool
		SetReplicaID(string)
		GetReplicaID() string
		GetNodeName() string
		SetNodeName(string)
		SetUCPClient(*ucpclient.APIClient)
		GetUCPClient() *ucpclient.APIClient
		GetDockerClient() *dc.Client
		VersionCheck() error
	}

	BootstrapConfig struct {
		AuthConfig    string
		Context       context.Context
		Client        *dc.Client
		CertPath      string
		Host          string
		TlsVerify     bool
		Username      string
		Password      string
		HubUsername   string
		HubPassword   string
		SettingsStore hubconfig.SettingsStore
		ReplicaID     string
		Node          string
		UCPClient     *ucpclient.APIClient
	}
)

func NewFromSocket(hubUsername, hubPassword string) (Bootstrap, error) {
	b := &BootstrapConfig{
		Context:    context.Background(),
		TlsVerify:  false,
		AuthConfig: MakeRegistryAuth(hubUsername, hubPassword),
	}

	client, err := DockerClientFromSocket(hubUsername, hubPassword)
	if err != nil {
		return nil, err
	}
	b.Client = client

	return b, nil
}

func NewFromJWT(host, username, password, jwt string, httpClient *http.Client, hubUsername, hubPassword string) (Bootstrap, error) {
	b := &BootstrapConfig{
		Context:    context.Background(),
		CertPath:   "",
		Host:       host,
		Username:   username,
		Password:   password,
		TlsVerify:  true,
		AuthConfig: MakeRegistryAuth(hubUsername, hubPassword),
	}

	client, err := DockerClientFromJWT(host, jwt, httpClient, hubUsername, hubPassword)
	if err != nil {
		return nil, err
	}
	b.Client = client

	return b, nil
}

func NewFromBundle(host, username, password, certPath, hubUsername, hubPassword string) (Bootstrap, error) {
	b := &BootstrapConfig{
		Context:    context.Background(),
		CertPath:   certPath,
		Host:       host,
		Username:   username,
		Password:   password,
		TlsVerify:  true,
		AuthConfig: MakeRegistryAuth(hubUsername, hubPassword),
	}

	client, err := DockerClientFromBundle(host, certPath, hubUsername, hubPassword)
	if err != nil {
		return nil, err
	}
	b.Client = client

	return b, nil
}
