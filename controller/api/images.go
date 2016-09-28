package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/engine-api/types/filters"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) listImages(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	allParam := rc.QueryVars.Get("all")
	all := false
	if allParam != "" {
		all = true
	}

	f, err := filters.FromParam(rc.QueryVars.Get("filters"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	images, err := a.manager.ListUserImages(rc.Auth, all, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(images); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
