package resources

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/manager"
)

var ErrNilContainerCreateRequest = errors.New("cannot parse nil ContainerCreateRequest")
var ErrNoContainerID = errors.New("cannot create container resource from empty container ID")
var ErrContainerOwnerLabelNotFound = errors.New("unable to detect an owner label for the given container")

// POST body of a container create request
// TODO: submit a PR to engine-api for this
type ContainerCreateRequest struct {
	*container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

type ContainerResourceRequest struct {
	ContainerCreateRequest *ContainerCreateRequest
}

func (r *ContainerResourceRequest) TeamLabel() (string, error) {
	if r.ContainerCreateRequest == nil {
		return "", ErrNilContainerCreateRequest
	}
	if l, ok := r.ContainerCreateRequest.Config.Labels[orca.UCPAccessLabel]; ok {
		log.Debugf("found UCP access label: %s", l)
		return l, nil
	}
	return "", nil
}

func (r *ContainerResourceRequest) OwnerLabel() (string, error) {
	if r.ContainerCreateRequest == nil {
		return "", ErrNilContainerCreateRequest
	}
	if l, ok := r.ContainerCreateRequest.Config.Labels[orca.UCPOwnerLabel]; ok {
		log.Debugf("found UCP owner label: %s", l)
		return l, nil
	}
	return "", ErrContainerOwnerLabelNotFound
}

func NewContainerResourceFromID(action ResourceAction, containerID string, man manager.Manager) (*CRUDResourceRequest, error) {
	if containerID == "" {
		// Not enough information for a ContainerResource
		return nil, ErrNoContainerID
	}
	containerJSON, err := man.Container(containerID)
	if err != nil {
		// Use a common error with engine-api
		// See https://github.com/docker/engine-api/blob/master/client/errors.go#L35
		return nil, fmt.Errorf("Error: No such container: %s", containerID)
	}
	containerCreateRequest := &ContainerCreateRequest{
		Config:     containerJSON.Config,
		HostConfig: containerJSON.HostConfig,
	}
	if containerCreateRequest.Config.Labels == nil {
		containerCreateRequest.Config.Labels = make(map[string]string)
	}
	return &CRUDResourceRequest{
		LabelledResource: &ContainerResourceRequest{
			ContainerCreateRequest: containerCreateRequest,
		},
		Action:  action,
		manager: man,
	}, nil
}

func NewContainerResourceFromRequest(containerCreateRequest *ContainerCreateRequest, man manager.Manager) (*CRUDResourceRequest, error) {
	if containerCreateRequest == nil {
		return nil, ErrNilContainerCreateRequest
	}
	if containerCreateRequest.Config.Labels == nil {
		containerCreateRequest.Config.Labels = make(map[string]string)
	}
	return &CRUDResourceRequest{
		LabelledResource: &ContainerResourceRequest{
			ContainerCreateRequest: containerCreateRequest,
		},
		Action:  ActionCreate,
		manager: man,
	}, nil
}
