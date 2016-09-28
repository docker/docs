package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// ping is the unauthenticated version of ping, which has a standard http.Handler signature
func (a *Api) ping(w http.ResponseWriter, r *http.Request) {
	// API spec doesn't define any error states besides 500, so
	// for now just always return OK - in the future consider
	// adding logic for very rudimentary health checking
	w.Header().Set("content-type", "text/plain")

	err := a.manager.GetSelfStatus()

	if err == nil {
		w.Write([]byte("OK"))
	} else {
		log.Debug(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
	}
}
