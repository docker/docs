package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/orca/auth"
)

func (a *Api) login(w http.ResponseWriter, r *http.Request) {
	a.manager.TrackClientInfo(r)

	var creds *Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, err := a.manager.AuthenticateUsernamePassword(creds.Username, creds.Password, r.RemoteAddr)
	if err == auth.ErrInvalidPassword || err == auth.ErrAccountDoesNotExist {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := auth.AuthToken{
		Token: ctx.SessionToken,
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) logout(w http.ResponseWriter, r *http.Request) {
	authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if err := a.manager.Logout(authToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
