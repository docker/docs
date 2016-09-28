package resources

import (
	"errors"
	"fmt"

	"github.com/cloudflare/cfssl/log"
	"github.com/docker/engine-api/types"
	"github.com/docker/orca"
	"github.com/docker/orca/controller/manager"
)

var ErrNilNetworkCreateRequest = errors.New("cannot parse nil network create request")
var ErrNoNetworkID = errors.New("cannot create network resource from empty network ID")
var ErrNetworkOwnerLabelNotFound = errors.New("unable to detect an owner label for the given network")

type NetworkResourceRequest struct {
	NetworkCreateRequest *types.NetworkCreateRequest
}

func (r *NetworkResourceRequest) TeamLabel() (string, error) {
	if r.NetworkCreateRequest == nil {
		return "", ErrNilNetworkCreateRequest
	}
	if l, ok := r.NetworkCreateRequest.NetworkCreate.Labels[orca.UCPAccessLabel]; ok {
		log.Debugf("found UCP access label: %s", l)
		return l, nil
	}
	return "", nil
}

func (r *NetworkResourceRequest) OwnerLabel() (string, error) {
	if r.NetworkCreateRequest == nil {
		return "", ErrNilNetworkCreateRequest
	}
	if l, ok := r.NetworkCreateRequest.NetworkCreate.Labels[orca.UCPOwnerLabel]; ok {
		log.Debugf("found UCP owner label: %s", l)
		return l, nil
	}
	return "", ErrNetworkOwnerLabelNotFound
}

func NewNetworkResourceFromID(action ResourceAction, networkID string, man manager.Manager) (*CRUDResourceRequest, error) {
	if networkID == "" {
		// Not enough information for a network
		return nil, ErrNoNetworkID
	}
	network, err := man.Network(networkID)
	if err != nil {
		// TODO: Use a common error with engine-api
		return nil, fmt.Errorf("Error: No such network: %s", networkID)
	}

	// Create a NetworkCreate from the inspected NetworkResource
	netCreate := types.NetworkCreate{
		Labels:     network.Labels,
		Options:    network.Options,
		Internal:   network.Internal,
		IPAM:       network.IPAM,
		EnableIPv6: network.EnableIPv6,
		Driver:     network.Driver,
	}
	netCreateRequest := types.NetworkCreateRequest{
		NetworkCreate: netCreate,
		Name:          network.Name,
	}

	return &CRUDResourceRequest{
		LabelledResource: &NetworkResourceRequest{
			NetworkCreateRequest: &netCreateRequest,
		},
		Action:  action,
		manager: man,
	}, nil
}

func NewNetworkResourceFromRequest(networkCreateRequest *types.NetworkCreateRequest, man manager.Manager) (*CRUDResourceRequest, error) {
	if networkCreateRequest == nil {
		return nil, ErrNilNetworkCreateRequest
	}
	return &CRUDResourceRequest{
		LabelledResource: &NetworkResourceRequest{
			NetworkCreateRequest: networkCreateRequest,
		},
		Action:  ActionCreate,
		manager: man,
	}, nil
}
