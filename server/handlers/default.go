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

	"github.com/docker/notary/errors"
	"github.com/docker/notary/server/storage"
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

// AddHandler adds the provided json data for the role and GUN specified in the URL
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
	store, ok := s.(*storage.MySQLStorage)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	gun := vars["imageName"]
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
	err = store.UpdateCurrent(gun, tufRole, version, input)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	return nil
}

// GetHandler returns the json for a specified role and GUN.
func GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("versionStore")
	store, ok := s.(*storage.MySQLStorage)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	gun := vars["imageName"]
	tufRole := vars["tufRole"]
	data, err := store.GetCurrent(gun, tufRole)
	logrus.Debug("JSON: ", string(data))
	if err != nil {
		logrus.Errorf("[Notary Server] 500 GET repository: %s, role: %s", gun, tufRole)
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	if data == nil {
		logrus.Errorf("[Notary Server] 404 GET repository: %s, role: %s", gun, tufRole)
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

// DeleteHandler deletes all data for a GUN. A 200 responses indicates success.
func DeleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("versionStore")
	store, ok := s.(*storage.MySQLStorage)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	gun := vars["imageName"]
	err := store.Delete(gun)
	if err != nil {
		logrus.Errorf("[Notary Server] 500 DELETE repository: %s", gun)
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	return nil
}
