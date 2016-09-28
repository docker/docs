// Package client implements a wrapper on top of engine-api tailored for the bootstrapper
package client

import (
	"errors"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

type EngineClient struct {
	client       clientShim
	bootstrapper *types.ContainerJSON
}

var (
	ErrNoBootstrapContainerFound = errors.New(
		"Unable to find the bootstrap container.  Did you forget to run it with \"--name " + orcaconfig.BootstrapContainerName + "\" perhaps?")
	ErrMultipleBootstrapsRunning = errors.New(
		"Multiple bootstrap containers detected.  Is another orca install already running?")
)

// NewBareClient creates a generic EngineClient that can exist outside the bootstrapper package
func NewBareClient() (*EngineClient, error) {
	ret := &EngineClient{}

	// Check the socket before we try to connect
	if err := ret.CheckSocket(); err != nil {
		return nil, err
	}

	log.Debugf("Connecting to docker %s", config.DockerSock)
	docker, err := dockerclient.NewClient(config.DockerSock, "", nil, nil)
	if err != nil {
		log.Error("Unable to connect to the docker engine")
		return nil, err
	}
	shim := realShim{client: docker}
	ret.client = clientShim(shim)
	ret.bootstrapper = &types.ContainerJSON{}
	ret.bootstrapper.Config = &container.Config{
		Tty: false,
	}
	return ret, nil
}

// Generate a new EngineClient, and perform the common pre-flight checks common to all bootstrap scenarios
func New() (*EngineClient, error) {
	ret, err := NewBareClient()
	if err != nil {
		return nil, err
	}

	// Perform pre-checks that always apply
	if err := ret.CheckKernelVersion(); err != nil {
		return nil, err
	}
	if err := ret.CheckDockerVersion(); err != nil {
		return nil, err
	}

	infos := ret.FindContainers([]string{orcaconfig.BootstrapPhase2ContainerName})
	if infos == nil {
		infos = ret.FindContainers([]string{orcaconfig.BootstrapContainerName})
		if infos == nil {
			return nil, ErrNoBootstrapContainerFound
		}
	}
	if len(infos) > 1 {
		return nil, ErrMultipleBootstrapsRunning
	}

	ret.bootstrapper = infos[0]

	// Detect if we appear to be running via UCP and fail fast
	labels := ret.bootstrapper.Config.Labels
	if labels["com.docker.ucp.access.owner"] != "" {
		return nil, fmt.Errorf("You appear to be running the UCP tool connected to UCP.  This is not supported.   You must change your DOCKER_HOST to point directly to the underlying engine, or run this command locally on the system you want to target.")
	} else if labels["com.docker.swarm.id"] != "" {
		return nil, fmt.Errorf("You appear to be running the UCP tool connected to Swarm.  This is not supported.   You must change your DOCKER_HOST to point directly to the underlying engine, or run this command locally on the system you want to target.")
	}

	return ret, nil
}

// TODO: make the EngineClient implement a clientShim
func (ec *EngineClient) InspectContainer(name string) (types.ContainerJSON, error) {
	return ec.client.InspectContainer(name)
}

// Check the TTY state of phase 1
func (ec *EngineClient) IsTty() bool {
	return ec.bootstrapper.Config.Tty
}

func (ec *EngineClient) PullImage(imageName string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return ec.client.PullImage(imageName, options)
}

func (ec *EngineClient) ContainerLogs(containerID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return ec.client.ContainerLogs(containerID, options)
}

// This is a first step to avoiding bloating the client shim and rather rely on engine-api
func (ec *EngineClient) GetClient() *dockerclient.Client {
	return ec.client.GetClient()
}
