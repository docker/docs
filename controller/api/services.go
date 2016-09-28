package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/engine-api/types/filters"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/resources"
)

func (a *Api) listServices(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	f, err := filters.FromParam(rc.QueryVars.Get("filters"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	services, err := a.manager.ListUserServices(rc.Auth, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(services); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// createService injects the appropriate owner label to a service create request
func (a *Api) createService(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// Attempt to cast the MainResource as a CRUDResourceRequest
	var serviceReq *resources.ServiceResourceRequest
	var found bool

	crudReq, foundCrud := rc.MainResource.(*resources.CRUDResourceRequest)
	if !foundCrud {
		http.Error(w, fmt.Sprintf("internal error: unable to locate CRUD resource in context"), http.StatusInternalServerError)
	}

	// Check if the underlying Labelled Resource can be cast to a service resource
	serviceReq, found = crudReq.LabelledResource.(*resources.ServiceResourceRequest)
	if !found {
		http.Error(w, fmt.Sprintf("internal error: unable to locate service resource in context"), http.StatusInternalServerError)
		return
	}

	// Set the owner label to the current user
	serviceSpec := serviceReq.ServiceSpec
	// Create the service labels map if not populated already
	if serviceSpec.Labels == nil {
		serviceSpec.Labels = make(map[string]string)
	}
	serviceSpec.Labels[orca.UCPOwnerLabel] = rc.Auth.User.Username

	// Pass the UCP owner label to all created containers
	// Create the containerSpec labels map if not populated already
	if serviceSpec.TaskTemplate.ContainerSpec.Labels == nil {
		serviceSpec.TaskTemplate.ContainerSpec.Labels = make(map[string]string)
	}
	serviceSpec.TaskTemplate.ContainerSpec.Labels[orca.UCPOwnerLabel] = rc.Auth.User.Username

	// Copy the service access label to the container spec, if it exists
	if _, ok := serviceSpec.Labels[orca.UCPAccessLabel]; ok {
		serviceSpec.TaskTemplate.ContainerSpec.Labels[orca.UCPAccessLabel] = serviceSpec.Labels[orca.UCPAccessLabel]
	}

	// Re-forge the request
	data, err := json.Marshal(serviceSpec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	rc.Request.ContentLength = int64(len(data))

	// Throw it to the engine and wrap with Registry header logic
	a.engineRegistryRedirect(w, rc, serviceSpec.TaskTemplate.ContainerSpec.Image, "pull")
}
