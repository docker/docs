package proxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleHTTPer is an interface that permits the handling of HTTP
type HandleHTTPer interface {
	HandleHTTP(writer http.ResponseWriter, r *http.Request) error
}

func withError(handler HandleHTTPer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		if err := handler.HandleHTTP(writer, req); err != nil {
			log.Println(err)
			http.Error(writer, err.Error(), 500)
		}
	}
}

func addRoutes(router *mux.Router, path string, handler HandleHTTPer) {
	addRoute(router, "/v{version:[0-9.]+}"+path, handler)
	addRoute(router, path, handler)
}

func addRoute(router *mux.Router, path string, handler HandleHTTPer) {
	router.HandleFunc(path, withError(handler))
}
