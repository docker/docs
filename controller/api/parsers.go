package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"

	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/resources"
)

var ErrAccessDenied = errors.New("access denied")

// NoAccess is a pipeline.Parser that creates a PublicResource.
func (a *Api) nilParser(rc *ctx.OrcaRequestContext) (int, error) {
	tr := resources.NewPublicResource()
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

// adminParser is a pipeline.Parser that creates an AdminResource.
func (a *Api) adminParser(rc *ctx.OrcaRequestContext) (int, error) {
	tr := resources.NewAdminResource()
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

// userParser is a pipeline.Parser that creates an UserResource.
func (a *Api) userParser(rc *ctx.OrcaRequestContext) (int, error) {
	ur := resources.NewUserResource(rc.PathVars["username"])
	if ur == nil {
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, ur)
	rc.MainResource = ur
	return http.StatusOK, nil
}

// teamParser is a pipeline.Parser that creates a teamResource
func (a *Api) teamParser(rc *ctx.OrcaRequestContext) (int, error) {
	tr := resources.NewTeamResource(rc.PathVars["id"])
	if tr == nil {
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

/* Container Resource Parsers */

// containerParser is creates a containerResource from a "id" or "name" Path Variable
func (a *Api) containerParser(rc *ctx.OrcaRequestContext, action resources.ResourceAction) (int, error) {
	// attempt to get the container ID from either the "id" or "name" path variables
	container := rc.PathVars["id"]
	if container == "" {
		container = rc.PathVars["name"]
	}

	tr, err := resources.NewContainerResourceFromID(action, container, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for container %s: %s", rc.Auth.User.Username, container, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

// containerCreateParser parses a container creation requests and creates all appropriate resources
func (a *Api) containerCreateParser(rc *ctx.OrcaRequestContext) (int, error) {
	// Parse a ContainerCreateConfig from the body of the request
	if len(rc.BodyBuffer()) == 0 {
		return http.StatusBadRequest, fmt.Errorf("unable to parse create container request with empty body")
	}

	var containerCreateRequest resources.ContainerCreateRequest
	if err := json.NewDecoder(rc.Body()).Decode(&containerCreateRequest); err != nil {
		return http.StatusBadRequest, err
	}

	// Require that the resulting containerCreateRequest has both a Config and a HostConfig
	if containerCreateRequest.Config == nil || containerCreateRequest.HostConfig == nil {
		return http.StatusBadRequest, fmt.Errorf("unsupported container create payload: Config and HostConfig cannot be nil")
	}

	tr, err := resources.NewContainerResourceFromRequest(&containerCreateRequest, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for a new container: %s", rc.Auth.User.Username, err)
		return http.StatusForbidden, ErrAccessDenied
	}

	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr

	// If the container is attached to a network, add a Network resource to the context and require Restricted Control
	network := string(containerCreateRequest.HostConfig.NetworkMode)
	if network != "" && network != "default" {
		nr, err := resources.NewNetworkResourceFromID(resources.ActionRestricted, network, a.manager)
		if err != nil {
			if rc.Auth.User.Admin {
				// Give the admin full disclosure on whether the network exists or not
				return http.StatusBadRequest, err
			}
			// Other users get an access denied message to avoid leaking information
			log.Infof("Denied access to user %s for network %s: %s", rc.Auth.User.Username, network, err)
			return http.StatusBadRequest, ErrAccessDenied
		}
		rc.Resources = append(rc.Resources, nr)
	}

	return http.StatusOK, nil
}

func (a *Api) containerViewOnly(rc *ctx.OrcaRequestContext) (int, error) {
	return a.containerParser(rc, resources.ActionView)
}

func (a *Api) containerRestrictedControl(rc *ctx.OrcaRequestContext) (int, error) {
	return a.containerParser(rc, resources.ActionRestricted)
}

func (a *Api) containerFullControl(rc *ctx.OrcaRequestContext) (int, error) {
	return a.containerParser(rc, resources.ActionFull)
}

/* Network Resource Parsers */

func (a *Api) networkParser(rc *ctx.OrcaRequestContext, action resources.ResourceAction) (int, error) {
	networkID := rc.PathVars["id"]
	tr, err := resources.NewNetworkResourceFromID(action, networkID, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for network %s: %s", rc.Auth.User.Username, networkID, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

func (a *Api) networkViewOnly(rc *ctx.OrcaRequestContext) (int, error) {
	return a.networkParser(rc, resources.ActionView)
}

func (a *Api) networkRestrictedControl(rc *ctx.OrcaRequestContext) (int, error) {
	return a.networkParser(rc, resources.ActionRestricted)
}

func (a *Api) networkCreateParser(rc *ctx.OrcaRequestContext) (int, error) {
	// Parse a networkCreateRequest from the body of the request
	if len(rc.BodyBuffer()) == 0 {
		return http.StatusBadRequest, fmt.Errorf("unable to parse create network request with empty body")
	}

	var networkCreateRequest types.NetworkCreateRequest
	if err := json.NewDecoder(rc.Body()).Decode(&networkCreateRequest); err != nil {
		return http.StatusBadRequest, err
	}

	tr, err := resources.NewNetworkResourceFromRequest(&networkCreateRequest, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for a new network: %s", rc.Auth.User.Username, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

func (a *Api) networkConnectParser(rc *ctx.OrcaRequestContext) (int, error) {
	// Parse the network resource from the request ID as a main resource
	res, err := a.networkParser(rc, resources.ActionRestricted)
	if err != nil {
		return res, err
	}

	// Extract the body of the network connect request
	if len(rc.BodyBuffer()) == 0 {
		return http.StatusBadRequest, fmt.Errorf("unable to parse network connect request with empty body")
	}

	var networkConnectRequest types.NetworkConnect
	if err := json.NewDecoder(rc.Body()).Decode(&networkConnectRequest); err != nil {
		return http.StatusBadRequest, err
	}

	// Create the container resource for the affected container
	container := networkConnectRequest.Container
	tr, err := resources.NewContainerResourceFromID(resources.ActionRestricted, container, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for container %s: %s", rc.Auth.User.Username, container, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	return http.StatusOK, nil
}

func (a *Api) networkDisconnectParser(rc *ctx.OrcaRequestContext) (int, error) {
	// Parse the network resource from the request ID as a main resource
	res, err := a.networkParser(rc, resources.ActionRestricted)
	if err != nil {
		return res, err
	}

	// Extract the body of the network connect request
	if len(rc.BodyBuffer()) == 0 {
		return http.StatusBadRequest, fmt.Errorf("unable to parse create network disconnect with empty body")
	}

	var networkDisconnectRequest types.NetworkDisconnect
	if err := json.NewDecoder(rc.Body()).Decode(&networkDisconnectRequest); err != nil {
		return http.StatusBadRequest, err
	}

	// Create the container resource for the affected container
	container := networkDisconnectRequest.Container
	tr, err := resources.NewContainerResourceFromID(resources.ActionRestricted, container, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for container %s: %s", rc.Auth.User.Username, container, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	return http.StatusOK, nil
}

/* Service Resource Parsers */

func (a *Api) serviceParser(rc *ctx.OrcaRequestContext, action resources.ResourceAction) (int, error) {
	serviceID := rc.PathVars["id"]
	tr, err := resources.NewServiceResourceFromID(action, serviceID, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for service %s: %s", rc.Auth.User.Username, serviceID, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}

func (a *Api) serviceCreateParser(rc *ctx.OrcaRequestContext) (int, error) {
	// Parse a ServiceCreateRequest from the body of the request
	if len(rc.BodyBuffer()) == 0 {
		return http.StatusBadRequest, fmt.Errorf("unable to parse create service request with empty body")
	}

	var serviceSpec swarm.ServiceSpec
	if err := json.NewDecoder(rc.Body()).Decode(&serviceSpec); err != nil {
		return http.StatusBadRequest, err
	}

	tr, err := resources.NewServiceResourceFromRequest(&serviceSpec, a.manager)
	if err != nil {
		if rc.Auth.User.Admin {
			// Give the admin full disclosure on a malformed request
			return http.StatusBadRequest, err
		}
		// Other users get an access denied message to avoid leaking information
		log.Infof("Denied access to user %s for a new service: %s", rc.Auth.User.Username, err)
		return http.StatusForbidden, ErrAccessDenied
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	rc.RequiresNotary = true
	return http.StatusOK, nil
}

func (a *Api) serviceViewOnly(rc *ctx.OrcaRequestContext) (int, error) {
	return a.serviceParser(rc, resources.ActionView)
}

func (a *Api) serviceRestrictedControl(rc *ctx.OrcaRequestContext) (int, error) {
	return a.serviceParser(rc, resources.ActionRestricted)
}

// userAccessListParser generates a user or admin resource from the username query var
func (a *Api) userAccessListParser(rc *ctx.OrcaRequestContext) (int, error) {
	requestedUsername := rc.QueryVars.Get("username")
	if requestedUsername == "" {
		// If no username is specified as a query parameter, accesslists is admin-only
		tr := resources.NewAdminResource()
		rc.Resources = append(rc.Resources, tr)
		rc.MainResource = tr
		return http.StatusOK, nil
	}

	// Allow a user to view their own accesslists
	tr := resources.NewUserResource(requestedUsername)
	if tr == nil {
		return http.StatusBadRequest, fmt.Errorf("invalid username specified")
	}
	rc.Resources = append(rc.Resources, tr)
	rc.MainResource = tr
	return http.StatusOK, nil
}
