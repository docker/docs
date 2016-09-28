package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/orca/controller/ctx"
)

// If the list is empty, there's nothing to report (i.e., no banner)
func (a *Api) banner(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	banner := a.manager.GetBanner(rc.Auth)
	if err := json.NewEncoder(w).Encode(banner); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
