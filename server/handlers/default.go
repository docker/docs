package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf/data"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/vetinari/errors"
	"github.com/docker/vetinari/server/version"
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
	defer r.Body.Close()
	s := ctx.Value("versionStore")
	if s == nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store is nil"),
		}
	}
	store, ok := s.(*version.VersionDB)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	qdn := vars["imageName"]
	tufRole := vars["tufRole"]
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusBadRequest,
			Code:       9999,
			Err:        err,
		}
	}
	meta := &data.SignedTargets{}
	err = json.Unmarshal(input, meta)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusBadRequest,
			Code:       9999,
			Err:        err,
		}
	}
	version := meta.Signed.Version
	err = store.UpdateCurrent(qdn, tufRole, version, input)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	return nil
}

// GetHandler accepts urls in the form /<imagename>/<tuf file>.json
func GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("versionStore")
	store, ok := s.(*version.VersionDB)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	qdn := vars["imageName"]
	tufRole := vars["tufRole"]
	data, err := store.GetCurrent(qdn, tufRole)
	logrus.Debug("JSON: ", string(data))
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	if data == nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusNotFound,
			Code:       9999,
			Err:        err,
		}
	}
	logrus.Debug("Writing data")
	w.Write(data)
	return nil
}
