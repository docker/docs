package manager

import (
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
)

// The manager.resources file contains wrappers around the engine-api clients related to
// obtaining engine-api types of cluster resources

// Classic Swarm resources
func (m DefaultManager) Container(id string) (types.ContainerJSON, error) {
	return m.client.ContainerInspect(context.TODO(), id)
}

// Engine-Provided resources

func (m DefaultManager) Network(id string) (types.NetworkResource, error) {
	return m.proxyClient.NetworkInspect(context.TODO(), id)
}

func (m DefaultManager) ListNetworks(f filters.Args) ([]types.NetworkResource, error) {
	return m.proxyClient.NetworkList(context.TODO(), types.NetworkListOptions{
		Filters: f,
	})
}
func (m DefaultManager) Service(id string) (swarm.Service, error) {
	service, _, err := m.proxyClient.ServiceInspectWithRaw(context.TODO(), id)
	return service, err
}

func (m DefaultManager) ListServices(f filters.Args) ([]swarm.Service, error) {
	return m.proxyClient.ServiceList(context.TODO(), types.ServiceListOptions{
		Filter: f,
	})
}

func (m DefaultManager) ListUserServices(ctx *auth.Context, f filters.Args) ([]swarm.Service, error) {
	allServices, err := m.ListServices(f)
	if err != nil {
		return nil, err
	}

	acct := ctx.User

	// return all services for admin users
	if acct.Admin {
		// Work around the issue where an empty array is read by engine-api as a nil array
		// This prevents further json encoding as a "null" value
		// TODO: resolve in engine-api
		if allServices == nil {
			return []swarm.Service{}, nil
		}
		return allServices, nil
	}

	access, err := m.GetAccess(ctx)
	if err != nil {
		return nil, err
	}

	userServices := []swarm.Service{}

	for _, s := range allServices {
		// check access labels
		matched := false
		for k, v := range s.Spec.Labels {
			if k != orca.UCPAccessLabel {
				continue
			}

			// check for label access
			if k == orca.UCPAccessLabel {
				if lvl, ok := access[v]; ok && lvl > auth.None {
					userServices = append(userServices, s)
					matched = true
					break
				}
			}
		}

		if matched {
			continue
		}

		for k, v := range s.Spec.Labels {
			if k != orca.UCPOwnerLabel {
				continue
			}

			// check for ownership
			if k == orca.UCPOwnerLabel && v == acct.Username {
				userServices = append(userServices, s)
				break
			}
		}
	}

	return userServices, nil
}

func (m DefaultManager) ListUserNetworks(ctx *auth.Context, f filters.Args) ([]types.NetworkResource, error) {
	allNetworks, err := m.ListNetworks(f)
	if err != nil {
		return nil, err
	}

	acct := ctx.User

	// return all services for admin users
	if acct.Admin {
		// Work around the issue where an empty array is read by engine-api as a nil array
		// This prevents further json encoding as a "null" value
		if allNetworks == nil {
			return []types.NetworkResource{}, nil
		}
		return allNetworks, nil
	}

	access, err := m.GetAccess(ctx)
	if err != nil {
		return nil, err
	}

	userNetworks := []types.NetworkResource{}

	for _, s := range allNetworks {
		// check access labels
		matched := false
		for k, v := range s.Labels {
			if k != orca.UCPAccessLabel {
				continue
			}

			// check for label access
			if k == orca.UCPAccessLabel {
				if lvl, ok := access[v]; ok && lvl > auth.None {
					userNetworks = append(userNetworks, s)
					matched = true
					break
				}
			}
		}

		if matched {
			continue
		}

		for k, v := range s.Labels {
			if k != orca.UCPOwnerLabel {
				continue
			}

			// check for ownership
			if k == orca.UCPOwnerLabel && v == acct.Username {
				userNetworks = append(userNetworks, s)
				break
			}
		}
	}

	return userNetworks, nil
}
