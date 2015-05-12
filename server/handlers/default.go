package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/vetinari/errors"
	"github.com/docker/vetinari/utils"
	repo "github.com/endophage/go-tuf"
	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/store"
	"github.com/endophage/go-tuf/util"
	"github.com/gorilla/mux"
)

// TODO: This is just for PoC. The real DB should be injected as part of
// the context for a final version.
var db = util.GetSqliteDB()

// MainHandler is the default handler for the server
func MainHandler(ctx utils.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
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
func AddHandler(ctx utils.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	vars := mux.Vars(r)
	log.Printf("AddHandler request for images name: %s and tag: %s", vars["imageName"], vars["tag"])

	local := store.DBStore(db, vars["imageName"])
	// parse body for correctness
	meta := data.FileMeta{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&meta)
	defer r.Body.Close()
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	// add to targets
	local.AddBlob(vars["tag"], meta)
	tufRepo, err := repo.NewRepo(ctx.Trust(), local, "sha256", "sha512")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	//	err = tufRepo.Init(false)
	//if err != nil {
	//	return &errors.HTTPError{
	//		HTTPStatus: http.StatusInternalServerError,
	//		Code:       9996,
	//		Err:        err,
	//	}
	//}
	err = tufRepo.AddTarget(vars["tag"], json.RawMessage{})
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9997,
			Err:        err,
		}
	}
	err = tufRepo.Sign("targets.json")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9998,
			Err:        err,
		}
	}
	tufRepo.Snapshot(repo.CompressionTypeNone)
	err = tufRepo.Sign("snapshot.json")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	tufRepo.Timestamp()
	err = tufRepo.Sign("timestamp.json")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	return nil
}

// RemoveHandler accepts urls in the form /<imagename>/<tag>
func RemoveHandler(ctx utils.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	// remove tag from tagets list
	vars := mux.Vars(r)
	log.Printf("RemoveHandler request for images name: %s and tag: %s", vars["imageName"], vars["tag"])

	local := store.DBStore(db, vars["imageName"])
	local.RemoveBlob(vars["tag"])
	tufRepo, err := repo.NewRepo(ctx.Trust(), local, "sha256", "sha512")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	_ = tufRepo.Init(false)
	tufRepo.RemoveTarget(vars["tag"])
	err = tufRepo.Sign("targets.json")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	tufRepo.Snapshot(repo.CompressionTypeNone)
	err = tufRepo.Sign("snapshot.json")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	tufRepo.Timestamp()
	err = tufRepo.Sign("timestamp.json")
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
func GetHandler(ctx utils.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	// generate requested file and serve
	vars := mux.Vars(r)
	log.Printf("GetHandler request for image name: %s and tuf-file: %s", vars["imageName"], vars["tufFile"])

	local := store.DBStore(db, vars["imageName"])

	meta, err := local.GetMeta()
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	w.Write(meta[vars["tufFile"]])
	return nil
}

// GenKeysHandler is the handler for generate keys endpoint
func GenKeysHandler(ctx utils.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	// remove tag from tagets list
	vars := mux.Vars(r)
	log.Printf("GenKeysHandler request for: %s", vars["imageName"])

	local := store.DBStore(db, vars["imageName"])
	tufRepo, err := repo.NewRepo(ctx.Trust(), local, "sha256", "sha512")
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	// init repo
	err = tufRepo.Init(false)
	if err != nil {
		return &errors.HTTPError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       9999,
			Err:        err,
		}
	}
	// gen keys

	// generate empty targets file

	// snapshot

	// timestamp
	return nil
}
