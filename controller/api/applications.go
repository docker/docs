package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) applications(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	apps, err := a.manager.Applications(rc.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(apps); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *Api) application(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	app, err := a.manager.Application(rc.Auth, rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(app); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) removeApplication(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	app, err := a.manager.Application(rc.Auth, rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.RemoveApplication(app); err != nil {
		log.Errorf("error deleting app: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) applicationContainers(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	app, err := a.manager.Application(rc.Auth, rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	containers, err := a.manager.ApplicationContainers(app)
	if err != nil {
		log.Errorf("error getting app containers: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(containers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) restartApplication(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	app, err := a.manager.Application(rc.Auth, rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.RestartApplication(app); err != nil {
		log.Errorf("error restarting app: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) stopApplication(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	app, err := a.manager.Application(rc.Auth, rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.StopApplication(app); err != nil {
		log.Errorf("error stopping app: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
