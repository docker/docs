package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/orca/controller/ctx"
)

type (
	Passwords struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}
)

func (a *Api) changePassword(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var passwds *Passwords
	if err := json.NewDecoder(rc.Body()).Decode(&passwds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := rc.Auth.User.Username
	if username == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := a.manager.ChangePassword(rc.Auth, username, passwds.OldPassword, passwds.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
