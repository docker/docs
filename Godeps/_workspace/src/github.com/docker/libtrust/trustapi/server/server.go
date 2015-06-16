package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/docker/libtrust/trustapi"
	"github.com/gorilla/mux"
)

var baseGraphs map[string][]byte
var baseDirectory string

func init() {
	baseGraphs = map[string][]byte{}
	flag.StringVar(&baseDirectory, "base", "", "Base directory for base graph files")
}

func register(r *mux.Router) {
	r.Get("graphbase").HandlerFunc(handleBaseGraph)
}

func main() {
	flag.Parse()
	files, err := filepath.Glob(filepath.Join(baseDirectory, "*.json"))
	if err != nil {
		log.Fatalf("Error loading json files: %s", err)
	}
	for _, f := range files {
		name := filepath.Base(f)
		name = name[:len(name)-5]
		b, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}
		baseGraphs[name] = b
	}

	r := trustapi.NewRouter("localhost")
	register(r)
	http.ListenAndServe("localhost:8092", r)
}

func handleBaseGraph(rw http.ResponseWriter, r *http.Request) {
	graphname := mux.Vars(r)["graphname"]
	log.Printf("Getting graph: %s", graphname)
	b, ok := baseGraphs[graphname]
	if !ok {
		rw.WriteHeader(404)
		fmt.Fprintf(rw, "base graph not found")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
