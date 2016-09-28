package client

import (
	"bytes"
	"io"
	"time"

	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"
)

// Shim for the DockerClient to facilitate testing
type clientShim interface {
	InspectImage(id string) (types.ImageInspect, error)
	PullImage(name string, options types.ImagePullOptions) (io.ReadCloser, error)
	RemoveImage(name string, force bool) ([]types.ImageDelete, error)
	Version() (types.Version, error)
	InspectContainer(id string) (types.ContainerJSON, error)
	CreateContainer(config *container.Config, hostConfig *container.HostConfig, name string) (types.ContainerCreateResponse, error)
	StartContainer(id string) error
	RemoveContainer(id string, force, volumes bool) error
	StopContainer(id string, timeout int) error
	ContainerRestart(id string, timeout int) error
	ContainerLogs(containerID string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	ContainerAttach(containerID string, options types.ContainerAttachOptions) (types.HijackedResponse, error)
	Info() (types.Info, error)
	ListContainers(options types.ContainerListOptions) ([]types.Container, error)
	CreateVolume(request types.VolumeCreateRequest) (types.Volume, error)
	VolumeInspect(volumeID string) (types.Volume, error)
	InspectNode(nodeID string) (swarm.Node, error)
	VolumeList() ([]*types.Volume, error)
	RemoveVolume(name string) error
	ContainerWait(containerID string) (int, error)
	SwarmInit(swarm.InitRequest) (string, error)
	SwarmLeave(force bool) error
	ServiceCreate(swarm.ServiceSpec) (types.ServiceCreateResponse, error)
	ServiceList(options types.ServiceListOptions) ([]swarm.Service, error)
	ServiceRemove(serviceID string) error
	GetClient() *dockerclient.Client
}

type realShim struct {
	client *dockerclient.Client
}

func (c realShim) GetClient() *dockerclient.Client {
	return c.client
}

func (c realShim) InspectImage(id string) (types.ImageInspect, error) {
	res, _, err := c.client.ImageInspectWithRaw(context.TODO(), id, false)
	return res, err
}
func (c realShim) PullImage(name string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return c.client.ImagePull(context.TODO(), name, options)
}
func (c realShim) RemoveImage(name string, force bool) ([]types.ImageDelete, error) {
	return c.client.ImageRemove(context.TODO(), name, types.ImageRemoveOptions{Force: force})
}
func (c realShim) Version() (types.Version, error) {
	return c.client.ServerVersion(context.TODO())
}
func (c realShim) InspectContainer(id string) (types.ContainerJSON, error) {
	return c.client.ContainerInspect(context.TODO(), id)
}
func (c realShim) CreateContainer(config *container.Config, hostConfig *container.HostConfig, name string) (types.ContainerCreateResponse, error) {
	return c.client.ContainerCreate(context.TODO(), config, hostConfig, nil, name)
}
func (c realShim) StartContainer(id string) error {
	return c.client.ContainerStart(context.TODO(), id, types.ContainerStartOptions{
		CheckpointID: "",
	})
}
func (c realShim) RemoveContainer(id string, force, volumes bool) error {
	return c.client.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{Force: force, RemoveVolumes: volumes})
}
func (c realShim) StopContainer(id string, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	return c.client.ContainerStop(context.TODO(), id, &timeout)
}
func (c realShim) ContainerRestart(id string, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	return c.client.ContainerRestart(context.TODO(), id, &timeout)
}
func (c realShim) ContainerLogs(containerID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return c.client.ContainerLogs(context.TODO(), containerID, options)
}
func (c realShim) ContainerAttach(containerID string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	return c.client.ContainerAttach(context.TODO(), containerID, options)
}

func (c realShim) Info() (types.Info, error) {
	return c.client.Info(context.TODO())
}
func (c realShim) InspectNode(nodeID string) (swarm.Node, error) {
	node, _, err := c.client.NodeInspectWithRaw(context.TODO(), nodeID)
	return node, err
}
func (c realShim) ListContainers(options types.ContainerListOptions) ([]types.Container, error) {
	return c.client.ContainerList(context.TODO(), options)
}
func (c realShim) CreateVolume(request types.VolumeCreateRequest) (types.Volume, error) {
	return c.client.VolumeCreate(context.TODO(), request)
}
func (c realShim) RemoveVolume(name string) error {
	return c.client.VolumeRemove(context.TODO(), name)
}
func (c realShim) VolumeInspect(volumeID string) (types.Volume, error) {
	return c.client.VolumeInspect(context.TODO(), volumeID)
}
func (c realShim) VolumeList() ([]*types.Volume, error) {
	ret, err := c.client.VolumeList(context.TODO(), filters.NewArgs())
	if err != nil {
		return nil, err
	}
	return ret.Volumes, nil
}

func (c realShim) ContainerWait(containerID string) (int, error) {
	return c.client.ContainerWait(context.TODO(), containerID)
}

func (c realShim) SwarmInit(swarmInitRequest swarm.InitRequest) (string, error) {
	return c.client.SwarmInit(context.TODO(), swarmInitRequest)
}

func (c realShim) SwarmLeave(force bool) error {
	return c.client.SwarmLeave(context.TODO(), force)
}

func (c realShim) ServiceCreate(spec swarm.ServiceSpec) (types.ServiceCreateResponse, error) {
	return c.client.ServiceCreate(context.TODO(), spec, types.ServiceCreateOptions{})
}

func (c realShim) ServiceList(options types.ServiceListOptions) ([]swarm.Service, error) {
	return c.client.ServiceList(context.TODO(), options)
}

func (c realShim) ServiceRemove(service string) error {
	return c.client.ServiceRemove(context.TODO(), service)
}

// Mock shim for testing
type testShim struct {
	inspectImageOut       types.ImageInspect
	inspectImageErr       error
	pullOut               io.ReadCloser
	pullErr               error
	removeImageOut        []types.ImageDelete
	removeImageErr        error
	versionOut            types.Version
	versionErr            error
	inspectContainerOut   types.ContainerJSON
	inspectContainerErr   error
	createOut             types.ContainerCreateResponse
	createErr             error
	startErr              error
	removeErr             error
	stopErr               error
	restartErr            error
	logsOut               io.ReadCloser
	logsErr               error
	infoOut               types.Info
	infoErr               error
	listOut               []types.Container
	listErr               error
	createVolOut          types.Volume
	createVolErr          error
	listVolOut            []*types.Volume
	listVolErr            error
	inspVolOut            types.Volume
	inspVolErr            error
	removeVolErr          error
	attachOut             types.HijackedResponse
	attachErr             error
	waitOut               int
	waitErr               error
	serviceCreateResponse types.ServiceCreateResponse
	serviceListResponse   []swarm.Service
}

func (c testShim) GetClient() *dockerclient.Client {
	return nil
}
func (c testShim) InspectImage(id string) (types.ImageInspect, error) {
	return c.inspectImageOut, c.inspectImageErr
}
func (c testShim) PullImage(name string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return c.pullOut, c.pullErr
}
func (c testShim) RemoveImage(name string, force bool) ([]types.ImageDelete, error) {
	return c.removeImageOut, c.removeImageErr
}
func (c testShim) Version() (types.Version, error) {
	return c.versionOut, c.versionErr
}
func (c testShim) InspectContainer(id string) (types.ContainerJSON, error) {
	return c.inspectContainerOut, c.inspectContainerErr
}
func (c testShim) CreateContainer(config *container.Config, hostConfig *container.HostConfig, name string) (types.ContainerCreateResponse, error) {
	return c.createOut, c.createErr
}
func (c testShim) StartContainer(id string) error {
	return c.startErr
}
func (c testShim) RemoveContainer(id string, force, volumes bool) error {
	return c.removeErr
}
func (c testShim) StopContainer(id string, timeout int) error {
	return c.stopErr
}
func (c testShim) ContainerRestart(id string, timeout int) error {
	return c.restartErr
}
func (c testShim) ContainerLogs(containerID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return c.logsOut, c.logsErr
}
func (c testShim) Info() (types.Info, error) {
	return c.infoOut, c.infoErr
}

// TODO: adapt to the rest of the codebase or deprecate the shim entirely
func (c testShim) InspectNode(nodeID string) (swarm.Node, error) {
	return swarm.Node{}, nil
}
func (c testShim) ListContainers(options types.ContainerListOptions) ([]types.Container, error) {
	return c.listOut, c.listErr
}
func (c testShim) CreateVolume(request types.VolumeCreateRequest) (types.Volume, error) {
	return c.createVolOut, c.createVolErr
}
func (c testShim) VolumeInspect(volumeID string) (types.Volume, error) {
	return c.inspVolOut, c.inspVolErr
}
func (c testShim) VolumeList() ([]*types.Volume, error) {
	return c.listVolOut, c.listVolErr
}
func (c testShim) RemoveVolume(name string) error {
	return c.removeVolErr
}
func (c testShim) ContainerAttach(containerID string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	return c.attachOut, c.attachErr
}

func (c testShim) ContainerWait(containerID string) (int, error) {
	return c.waitOut, c.waitErr
}

func (c testShim) SwarmInit(swarmInitRequest swarm.InitRequest) (string, error) {
	return "default", nil
}

func (c testShim) SwarmLeave(force bool) error {
	return nil
}

func (c testShim) ServiceCreate(spec swarm.ServiceSpec) (types.ServiceCreateResponse, error) {
	return c.serviceCreateResponse, nil
}

func (c testShim) ServiceList(options types.ServiceListOptions) ([]swarm.Service, error) {
	return c.serviceListResponse, nil
}

func (c testShim) ServiceRemove(service string) error {
	return nil
}

type closingBuffer struct {
	*bytes.Buffer
}

func (cb *closingBuffer) Close() error {
	// no-op
	return nil
}
