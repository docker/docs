package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) listConfig(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	subsystems := a.manager.ListConfigSubsystems()
	if err := json.NewEncoder(w).Encode(subsystems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Api) getConfig(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	if config, err := a.manager.GetSubsystemConfig(rc.PathVars["subsystem"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(config))
	}
}

func (a *Api) updateConfig(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	if err := a.manager.UserConfigUpdate(rc.PathVars["subsystem"], string(rc.BodyBuffer())); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
