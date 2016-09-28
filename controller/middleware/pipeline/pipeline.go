package pipeline

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/version"
)

// A MiddlewareInitializer generates an OrcaRequestContext from an HTTP request.
type MiddlewareInitializer func(http.ResponseWriter, *http.Request) (*ctx.OrcaRequestContext, error)

// A MiddlewareLayer processes the OrcaRequestContext and returns an error and an HTTP status code
type MiddlewareLayer func(*ctx.OrcaRequestContext) (int, error)

// Parser is the first middleware layer to be invoked and is explictly defined together with
// the final handler
type Parser MiddlewareLayer
type Handler func(http.ResponseWriter, *ctx.OrcaRequestContext)

type MiddlewarePipeline interface {
	// The Auth layer act as an initializer which creates the requestContext
	AddInitializer(MiddlewareInitializer) error

	// All other layers act upon the request context itself
	AddLayer(MiddlewareLayer)

	// Route is used to define a new Route in this pipeline
	Route(string, string, Parser, Handler)
}

// middlewarePipeline is a concrete implementation of MiddlewarePipeline
// An existing mux.Router is used to register all Routes
type middlewarePipeline struct {
	router      *mux.Router
	initializer func(http.ResponseWriter, *http.Request) (*ctx.OrcaRequestContext, error)
	layers      []MiddlewareLayer
}

var (
	ErrInitializerPresent = errors.New("an initializer is already present in the pipeline")
	ErrNilInitializer     = errors.New("the provided initializer is nil")
)

func (p *middlewarePipeline) AddInitializer(init MiddlewareInitializer) error {
	if init == nil {
		return ErrNilInitializer
	}
	if p.initializer != nil {
		return ErrInitializerPresent
	}
	p.initializer = init
	return nil
}

func (p *middlewarePipeline) AddLayer(ml MiddlewareLayer) {
	p.layers = append(p.layers, ml)
}

func (p *middlewarePipeline) Route(path string, method string, parser Parser, handler Handler) {
	if p.initializer == nil {
		log.Errorf("attempted route registration before a pipeline initializer was configured")
	}

	matcher := func(w http.ResponseWriter, r *http.Request) {
		rc, err := p.initializer(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Invoke the Parser
		status, err := parser(rc)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		for _, layer := range p.layers {
			status, err := layer(rc)
			if err != nil {
				// TODO: introduce error handlers
				http.Error(w, err.Error(), status)
				return
			}
		}

		// Inject a UCP Version header to every response
		w.Header().Set("UCP-Version", version.FullVersion())

		// Invoke the Handler and let it deal with the http ResponseWriter
		handler(w, rc)

	}

	p.router.HandleFunc(path, matcher).Methods(method)
}

func New(r *mux.Router) MiddlewarePipeline {
	return &middlewarePipeline{
		router:      r,
		initializer: nil,
	}
}
