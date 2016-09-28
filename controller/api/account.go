package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
)

func (a *Api) currentAccount(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(rc.Auth.User); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) accounts(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")
	teamID := rc.QueryVars.Get("teamId")

	// Only admins should be able to list accounts
	if !rc.Auth.User.Admin {
		http.Error(w, auth.ErrUnauthorized.Error(), http.StatusInternalServerError)
		return
	}

	accounts := []*auth.Account{}
	if teamID != "" {
		accts, err := a.manager.AccountsForTeam(rc.Auth, teamID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accounts = accts

	} else {
		accts, err := a.manager.Accounts(rc.Auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accounts = accts
	}

	// Never expose the password externally
	for _, acct := range accounts {
		acct.Password = ""
	}

	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) saveAccount(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// When saving an account, there are two accounts involved
	// the person making the API call, and the account being saved.

	var updateAccount *auth.Account
	if err := json.NewDecoder(rc.Body()).Decode(&updateAccount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if eventType, err := a.manager.SaveAccount(rc.Auth, updateAccount); err != nil {
		log.Errorf("error saving account: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {

		log.Debugf("updated account: name=%s", updateAccount.Username)
		w.Header().Set("Location", "/api/accounts/"+updateAccount.Username)
		if eventType == "add-account" {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (a *Api) account(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	username := rc.PathVars["username"]

	account, err := a.manager.Account(rc.Auth, username)
	if err != nil {
		log.Errorf("error retrieving account: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Never expose the password externally
	account.Password = ""

	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) deleteAccount(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// Only admins should be able to delete accounts
	if !rc.Auth.User.Admin {
		http.Error(w, auth.ErrUnauthorized.Error(), http.StatusInternalServerError)
		return
	}

	username := rc.PathVars["username"]
	account, err := a.manager.Account(rc.Auth, username)
	if err != nil {
		log.Errorf("error retrieving account: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.manager.DeleteAccount(rc.Auth, account); err != nil {
		log.Errorf("error deleting account: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("deleted account: username=%s", account.Username)
	w.WriteHeader(http.StatusNoContent)
}
