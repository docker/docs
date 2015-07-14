package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/notary/errors"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/server/timestamp"
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

// AtomicUpdateHandler will accept multiple TUF files and ensure that the storage
// backend is atomically updated with all the new records.
func AtomicUpdateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	defer r.Body.Close()
	s := ctx.Value("metaStore")
	if s == nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store is nil"),
		}
	}
	store, ok := s.(storage.MetaStore)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	vars := mux.Vars(r)
	gun := vars["imageName"]
	reader, err := r.MultipartReader()
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusBadRequest,
			Code:       9999,
			Err:        err,
		}
	}
	var updates []storage.MetaUpdate
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		role := strings.TrimSuffix(part.FileName(), ".json")
		if role == "" {
			return &errors.HTTPError{
				HTTPStatus: http.StatusBadRequest,
				Code:       9999,
				Err:        fmt.Errorf("Empty filename provided. No updates performed"),
			}
		} else if !data.ValidRole(role) {
			return &errors.HTTPError{
				HTTPStatus: http.StatusBadRequest,
				Code:       9999,
				Err:        fmt.Errorf("Invalid role: %s. No updates performed", role),
			}
		}
		meta := &data.SignedTargets{}
		var input []byte
		inBuf := bytes.NewBuffer(input)
		dec := json.NewDecoder(io.TeeReader(part, inBuf))
		err = dec.Decode(meta)
		if err != nil {
			return &errors.HTTPError{
				HTTPStatus: http.StatusBadRequest,
				Code:       9999,
				Err:        err,
			}
		}
		version := meta.Signed.Version
		updates = append(updates, storage.MetaUpdate{
			Role:    role,
			Version: version,
			Data:    inBuf.Bytes(),
		})
	}
	err = store.UpdateMany(gun, updates)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	return nil
}

// UpdateHandler adds the provided json data for the role and GUN specified in the URL
func UpdateHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	defer r.Body.Close()
	s := ctx.Value("metaStore")
	if s == nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store is nil"),
		}
	}
	store, ok := s.(storage.MetaStore)
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
	update := storage.MetaUpdate{
		Role:    tufRole,
		Version: meta.Signed.Version,
		Data:    input,
	}
	err = store.UpdateCurrent(gun, update)
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
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
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
	out, err := store.GetCurrent(gun, tufRole)
	if err != nil {
		if _, ok := err.(*storage.ErrNotFound); ok {
			return &errors.HTTPError{
				HTTPStatus: http.StatusNotFound,
				Code:       9999,
				Err:        err,
			}
		}
		logrus.Errorf("[Notary Server] 500 GET repository: %s, role: %s", gun, tufRole)
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	if out == nil {
		logrus.Errorf("[Notary Server] 404 GET repository: %s, role: %s", gun, tufRole)
		return &errors.HTTPError{
			HTTPStatus: http.StatusNotFound,
			Code:       9999,
			Err:        err,
		}
	}
	logrus.Debug("Writing data")
	w.Write(out)
	return nil
}

// DeleteHandler deletes all data for a GUN. A 200 responses indicates success.
func DeleteHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
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

// GetTimestampHandler returns a timestamp.json given a GUN
func GetTimestampHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	cryptoServiceVal := ctx.Value("cryptoService")
	cryptoService, ok := cryptoServiceVal.(signed.CryptoService)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("CryptoService not configured"),
		}
	}

	vars := mux.Vars(r)
	gun := vars["imageName"]

	out, err := timestamp.GetOrCreateTimestamp(gun, store, cryptoService)
	if err != nil {
		if _, ok := err.(*storage.ErrNoKey); ok {
			return &errors.HTTPError{
				HTTPStatus: http.StatusNotFound,
				Code:       9999,
				Err:        err,
			}
		}
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}

	logrus.Debug("Writing data")
	w.Write(out)
	return nil
}

// GetTimestampKeyHandler returns a timestamp public key, creating a new key-pair
// it if it doesn't yet exist
func GetTimestampKeyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	s := ctx.Value("metaStore")
	store, ok := s.(storage.MetaStore)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Version store not configured"),
		}
	}
	c := ctx.Value("cryptoService")
	crypto, ok := c.(signed.CryptoService)
	if !ok {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("CryptoService not configured"),
		}
	}

	vars := mux.Vars(r)
	gun := vars["imageName"]

	key, err := timestamp.GetOrCreateTimestampKey(gun, store, crypto, data.ECDSAKey)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}

	out, err := json.Marshal(key)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        fmt.Errorf("Error serializing key."),
		}
	}
	w.Write(out)
	return nil
}
