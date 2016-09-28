package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
)

func (a *Api) addAccessList(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var l *auth.AccessList

	if err := json.NewDecoder(rc.Body()).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if eventType, err := a.manager.SaveAccessList(l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Infof("created access list: role=%d label=%s", l.Role, l.Label)
		w.Header().Set("Location", fmt.Sprintf("/api/accesslists/%s/%s", l.TeamId, l.Id))
		if eventType == "add-access-list" {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (a *Api) accessLists(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	username := rc.QueryVars.Get("username")

	lists := []*auth.AccessList{}

	if username != "" {
		l, err := a.manager.AccessListsForAccount(rc.Auth, username)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lists = l
	} else {
		l, err := a.manager.AccessLists()
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lists = l
	}
	if err := json.NewEncoder(w).Encode(lists); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *Api) accessList(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	teamId := rc.PathVars["teamId"]
	id := rc.PathVars["id"]

	list, err := a.manager.AccessList(teamId, id)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) removeAccessList(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	teamId := rc.PathVars["teamId"]
	id := rc.PathVars["id"]

	if err := a.manager.RemoveAccessList(teamId, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("removed access list: id=%s", id)
	w.WriteHeader(http.StatusNoContent)
}
