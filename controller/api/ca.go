package api

import (
	"net/http"
)

func (a *Api) ca(w http.ResponseWriter, r *http.Request) {
	if a.tlsCACert != "" {
		w.Header().Set("content-type", "application/x-pem-file")
		w.Write([]byte(a.tlsCACert))
	} else {
		// This shouldn't happen in production, but dev mode might not be TLS enabled.
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("content-type", "text/plain")
		w.Write([]byte("No CA available"))
	}
}
