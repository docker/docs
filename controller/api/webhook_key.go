package api

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/dockerhub"
)

func (a *Api) hubWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["key"]

	key, err := a.manager.WebhookKey(k)
	if err != nil {
		log.Errorf("invalid webook key: key=%s from %s", k, r.RemoteAddr)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var webhook *dockerhub.Webhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		log.Errorf("error parsing webhook: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Index(webhook.Repository.RepoName, key.Image) == -1 {
		log.Errorf("webhook key image does not match: repo=%s image=%s", webhook.Repository.RepoName, key.Image)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	log.Infof("received webhook notification for %s", webhook.Repository.RepoName)
	// TODO @ehazlett - redeploy containers
}

func (a *Api) webhookKeys(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	keys, err := a.manager.WebhookKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) webhookKey(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	key, err := a.manager.WebhookKey(rc.PathVars["key"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) addWebhookKey(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var k *dockerhub.WebhookKey
	if err := json.NewDecoder(rc.Body()).Decode(&k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := a.manager.NewWebhookKey(k.Image)
	if err != nil {
		log.Errorf("error generating webhook key: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("saved webhook key image=%s", key.Image)
	if err := json.NewEncoder(w).Encode(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) deleteWebhookKey(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	k := rc.PathVars["key"]
	if err := a.manager.DeleteWebhookKey(k); err != nil {
		log.Errorf("error deleting webhook key: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("removed webhook key key=%s", k)
	w.WriteHeader(http.StatusNoContent)
}
