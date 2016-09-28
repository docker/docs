package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth/enzi"
)

func (a *Api) openidKeys(w http.ResponseWriter, r *http.Request) {
	authenticator := a.manager.GetAuthenticator()
	enziAuthenticator, ok := authenticator.(*enzi.Authenticator)
	if !ok {
		// We're not currently using the eNZi authenticator backend.
		// Just return a 404.
		http.NotFound(w, r)
		return
	}

	signingKeys, err := enziAuthenticator.ListSigningKeys()
	if err != nil {
		log.Errorf("unable to list signing keys for eNZi authenticator: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", uint(enzi.SigningKeyCacheMaxAge.Seconds())))

	if err := json.NewEncoder(w).Encode(signingKeys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
