package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
)

func (a *Api) createContainerLogsToken(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	// generate token
	u4, err := uuid.NewV4()
	if err != nil {
		log.Errorf("error generating containerlogs token token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := u4.String()

	cs := &orca.ContainerLogsToken{
		ContainerID: rc.PathVars["id"],
		Token:       token,
	}

	if err := a.manager.CreateContainerLogsToken(cs); err != nil {
		log.Errorf("error creating containerlogs token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cs); err != nil {
		log.Errorf("error encoding containerlogs token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf("created containerlogs token: container=%s", cs.ContainerID)
}

func (a *Api) containerLogsToken(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	cs, err := a.manager.ContainerLogsToken(rc.PathVars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) removeContainerLogsToken(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	cs, err := a.manager.ContainerLogsToken(rc.PathVars["token"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.manager.RemoveContainerLogsToken(cs); err != nil {
		log.Errorf("error removing containerlogs token: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
