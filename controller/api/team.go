package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
)

func (a *Api) teams(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	teams, err := a.manager.Teams(rc.Auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(teams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) saveTeam(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var team *auth.Team
	if err := json.NewDecoder(rc.Body()).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if eventType, err := a.manager.SaveTeam(rc.Auth, team); err != nil {
		log.Errorf("error saving team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Debugf("saved team: name=%s", team.Name)
		w.Header().Set("Location", "/api/teams/"+team.Id)
		if eventType == "add-team" {
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(&team); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (a *Api) updateTeam(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var team *auth.Team
	if err := json.NewDecoder(rc.Body()).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// override the team id in case an incorrect one was sent
	team.Id = rc.PathVars["id"]

	if _, err := a.manager.SaveTeam(rc.Auth, team); err != nil {
		log.Errorf("error updating team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("updated team: name=%s", team.Name)
	w.Header().Set("Location", "/api/teams/"+team.Id)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) team(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	team, err := a.manager.Team(rc.Auth, rc.PathVars["id"])
	if err != nil {
		log.Errorf("error retrieving team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) deleteTeam(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	team, err := a.manager.Team(rc.Auth, rc.PathVars["id"])
	if err != nil {
		log.Errorf("error retrieving team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.manager.DeleteTeam(rc.Auth, team); err != nil {
		log.Errorf("error deleting team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("deleted team: name=%s", team.Name)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) addMemberToTeam(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	id := rc.PathVars["id"]
	username := rc.PathVars["username"]

	if err := a.manager.AddMemberToTeam(rc.Auth, id, username); err != nil {
		log.Errorf("error retrieving team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("updated team: id=%s", id)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) removeMemberFromTeam(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	id := rc.PathVars["id"]
	username := rc.PathVars["username"]

	if err := a.manager.RemoveMemberFromTeam(rc.Auth, id, username); err != nil {
		log.Errorf("error retrieving team: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("updated team: id=%s", id)
	w.WriteHeader(http.StatusNoContent)
}
