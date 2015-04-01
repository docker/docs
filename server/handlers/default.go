package handlers

import (
	"encoding/json"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode("{}")
		if err != nil {
			w.Write([]byte("{server_error: 'Could not parse error message'}"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
