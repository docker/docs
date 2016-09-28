package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) catalog(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	repos, err := a.manager.GetDefaultCatalog()
	if err != nil {
		http.Error(w, "unable to connect to official registry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(repos); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
