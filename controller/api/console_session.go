package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
)

func (a *Api) createConsoleSession(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	// generate token
	u4, err := uuid.NewV4()
	if err != nil {
		log.Errorf("error generating console session token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := u4.String()

	cs := &orca.ConsoleSession{
		ContainerID: rc.PathVars["id"],
		Token:       token,
	}

	if err := a.manager.CreateConsoleSession(cs); err != nil {
		log.Errorf("error creating console session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cs); err != nil {
		log.Errorf("error encoding console session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("created console session: container=%s", cs.ContainerID)
}

func (a *Api) consoleSession(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	cs, err := a.manager.ConsoleSession(rc.PathVars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) removeConsoleSession(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	cs, err := a.manager.ConsoleSession(rc.PathVars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.RemoveConsoleSession(cs); err != nil {
		log.Errorf("error removing console session: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
