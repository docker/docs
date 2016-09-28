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

func (a *Api) listNetworks(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	f, err := filters.FromParam(rc.QueryVars.Get("filters"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	networks, err := a.manager.ListUserNetworks(rc.Auth, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(networks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// createNetwork injects the appropriate owner label to a network create request
func (a *Api) createNetwork(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// Attempt to cast the MainResource as a CRUDResourceRequest
	var networkReq *resources.NetworkResourceRequest
	var found bool

	crudReq, foundCrud := rc.MainResource.(*resources.CRUDResourceRequest)
	if !foundCrud {
		http.Error(w, fmt.Sprintf("internal error: unable to locate CRUD resource in context"), http.StatusInternalServerError)
	}

	// Check if the underlying Labelled Resource can be cast to a network resource
	networkReq, found = crudReq.LabelledResource.(*resources.NetworkResourceRequest)
	if !found {
		http.Error(w, fmt.Sprintf("internal error: unable to locate network resource in context"), http.StatusInternalServerError)
		return
	}

	networkCreateRequest := networkReq.NetworkCreateRequest

	// Set the owner label to the current user
	// Create the network labels map if not populated already
	if networkCreateRequest.NetworkCreate.Labels == nil {
		networkCreateRequest.NetworkCreate.Labels = make(map[string]string)
	}
	networkCreateRequest.NetworkCreate.Labels[orca.UCPOwnerLabel] = rc.Auth.User.Username

	// Re-forge the request
	data, err := json.Marshal(networkCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	rc.Request.ContentLength = int64(len(data))

	// Throw it to the engine
	a.engineRedirect(w, rc)
}
