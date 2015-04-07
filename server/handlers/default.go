package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	repo "github.com/docker/go-tuf"
	"github.com/docker/go-tuf/data"
	"github.com/docker/go-tuf/store"
	"github.com/docker/go-tuf/util"
	"github.com/gorilla/mux"
)

var db = util.GetSqliteDB()

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

// AddHandler accepts urls in the form /<imagename>/<tag>
func AddHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("AddHandler")
	vars := mux.Vars(r)
	local := store.DBStore(db, vars["imageName"])
	// parse body for correctness
	meta := data.FileMeta{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&meta)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to Decode JSON"))
		return
	}
	// add to targets
	local.AddBlob(vars["tag"], meta)
	tufRepo, err := repo.NewRepo(local, "sha256", "sha512")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to inistantiate TUF repository"))
		return
	}
	_ = tufRepo.Init(true)
	err = tufRepo.AddTarget(vars["tag"], json.RawMessage{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to add target"))
		log.Print(err)
		return
	}
	err = tufRepo.Sign("targets.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign targets file"))
		log.Print(err)
		return
	}
	tufRepo.Snapshot(repo.CompressionTypeNone)
	err = tufRepo.Sign("snapshot.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign snapshot file"))
		log.Print(err)
		return
	}
	tufRepo.Timestamp()
	err = tufRepo.Sign("timestamp.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign timestamps file"))
		log.Print(err)
		return
	}
	return
}

// RemoveHandler accepts urls in the form /<imagename>/<tag>
func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RemoveHandler")
	// remove tag from tagets list
	vars := mux.Vars(r)
	local := store.DBStore(db, vars["imageName"])
	local.RemoveBlob(vars["tag"])
	tufRepo, err := repo.NewRepo(local, "sha256", "sha512")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to inistantiate TUF repository"))
		return
	}
	_ = tufRepo.Init(true)
	tufRepo.RemoveTarget(vars["tag"])
	err = tufRepo.Sign("targets.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign targets file"))
		log.Print(err)
		return
	}
	tufRepo.Snapshot(repo.CompressionTypeNone)
	err = tufRepo.Sign("snapshot.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign snapshot file"))
		log.Print(err)
		return
	}
	tufRepo.Timestamp()
	err = tufRepo.Sign("timestamp.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to sign timestamps file"))
		log.Print(err)
		return
	}
	return
}

// GetHandler accepts urls in the form /<imagename>/<tuf file>.json
func GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetHandler")
	// generate requested file and serve
	vars := mux.Vars(r)
	local := store.DBStore(db, vars["imageName"])

	meta, err := local.GetMeta()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read TUF metadata"))
		return
	}
	w.Write(meta[vars["tufFile"]])
	return
}

func GenKeysHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GenKeysHandler")
	// remove tag from tagets list
	vars := mux.Vars(r)
	local := store.DBStore(db, vars["imageName"])
	tufRepo, err := repo.NewRepo(local, "sha256", "sha512")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to inistantiate TUF repository"))
		return
	}
	tufRepo.GenKey("root")
	tufRepo.GenKey("targets")
	tufRepo.GenKey("snapshot")
	tufRepo.GenKey("timestamp")
	//tufRepo.Sign("root.json")
	return
}
