package resources

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/swarm"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/manager"
)

var ErrNilServiceRequest = errors.New("cannot parse nil ServiceSpec")
var ErrNoServiceID = errors.New("cannot create service resource from empty service ID")
var ErrServiceOwnerLabelNotFound = errors.New("unable to detect an owner label for the given service")

type ServiceResourceRequest struct {
	ServiceSpec *swarm.ServiceSpec
}

func (r *ServiceResourceRequest) TeamLabel() (string, error) {
	if r.ServiceSpec == nil {
		return "", ErrNilServiceRequest
	}
	if l, ok := r.ServiceSpec.Labels[orca.UCPAccessLabel]; ok {
		log.Debugf("found UCP access label: %s", l)
		return l, nil
	}
	return "", nil
}

func (r *ServiceResourceRequest) OwnerLabel() (string, error) {
	if r.ServiceSpec == nil {
		return "", ErrNilServiceRequest
	}
	if l, ok := r.ServiceSpec.Labels[orca.UCPOwnerLabel]; ok {
		log.Debugf("found UCP owner label: %s", l)
		return l, nil
	}
	return "", ErrServiceOwnerLabelNotFound
}

func NewServiceResourceFromID(action ResourceAction, serviceID string, man manager.Manager) (*CRUDResourceRequest, error) {
	if serviceID == "" {
		// Not enough information for a service
		return nil, ErrNoServiceID
	}
	service, err := man.Service(serviceID)
	if err != nil {
		// TODO: Use a common error with engine-api
		return nil, fmt.Errorf("Error: No such service: %s", serviceID)
	}
	return &CRUDResourceRequest{
		LabelledResource: &ServiceResourceRequest{
			ServiceSpec: &(service.Spec),
		},
		Action:  action,
		manager: man,
	}, nil
}

func NewServiceResourceFromRequest(serviceRequest *swarm.ServiceSpec, man manager.Manager) (*CRUDResourceRequest, error) {
	if serviceRequest == nil {
		return nil, ErrNilServiceRequest
	}
	return &CRUDResourceRequest{
		LabelledResource: &ServiceResourceRequest{
			ServiceSpec: serviceRequest,
		},
		Action:  ActionCreate,
		manager: man,
	}, nil
}
