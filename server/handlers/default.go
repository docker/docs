package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/docker/vetinari/errors"
)

// MainHandler is the default handler for the server
func MainHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	if r.Method == "GET" {
		err := json.NewEncoder(w).Encode("{}")
		if err != nil {
			w.Write([]byte("{server_error: 'Could not parse error message'}"))
		}
	} else {
		//w.WriteHeader(http.StatusNotFound)
		return &errors.HTTPError{
			HTTPStatus: http.StatusNotFound,
			Code:       9999,
			Err:        nil,
		}
	}
	return nil
}

// AddHandler accepts urls in the form /<imagename>/<tag>
func UpdateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	return nil
}

// GetHandler accepts urls in the form /<imagename>/<tuf file>.json
func GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	return nil
}
