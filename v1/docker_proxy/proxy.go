package proxy

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	verbose bool
)

// SetVerboseLogging sets the logging to verbose
func SetVerboseLogging(value bool) {
	verbose = value
}

// Serve starts the proxy
func Serve(listener net.Listener, underlyingPath string, approver Approver, portRewriter PortRewriter, envRewriter EnvRewriter) error {
	backendDialer, err := NewBackendDialer("unix", underlyingPath)
	if err != nil {
		log.Fatal(err)
	}

	m := newRouter(backendDialer, approver, portRewriter, envRewriter, nil)

	return http.Serve(listener, m)
}

// ServeBackend starts the backend proxy
func ServeBackend(listener net.Listener, backendDialer BackendDialer, portRewriter PortRewriter, mountRewriter MountRewriter, envRewriter EnvRewriter) error {
	m := newRouter(backendDialer, nil, portRewriter, envRewriter, mountRewriter)

	return http.Serve(listener, m)
}

// ServeWindows starts the windows proxy
func ServeWindows(addr string, backendDialer BackendDialer, portRewriter PortRewriter, mountRewriter MountRewriter, envRewriter EnvRewriter) error {
	m := newRouter(backendDialer, nil, portRewriter, envRewriter, mountRewriter)

	return http.ListenAndServe(addr, m)
}

func newRouter(backendDialer BackendDialer, approver Approver, portRewriter PortRewriter, envRewriter EnvRewriter, mountRewriter MountRewriter) http.Handler {
	router := mux.NewRouter()

	addRoutes(router, "/containers/create", &proxyCreate{
		passthru:      newPassthru(backendDialer, &nopRewriter{}),
		mountRewriter: mountRewriter,
		envRewriter:   envRewriter,
	})
	addRoutes(router, "/containers/{name:.*}/start", &proxyStart{
		passthru:      newPassthru(backendDialer, &nopRewriter{}),
		approver:      approver,
		backendDialer: backendDialer,
	})
	addRoutes(router, "/containers/{name:.*}/json", newPassthru(backendDialer, NewInspectRewriter(portRewriter)))
	addRoutes(router, "/containers/json", newPassthru(backendDialer, NewPsRewriter(portRewriter)))

	// From https://github.com/docker/docker/search?utf8=%E2%9C%93&q=%22router.Cancellable%28%22&type=Code
	addRoutes(router, "/build", newCancellablePassthru(backendDialer, &nopRewriter{}))
	addRoutes(router, "/images/create", newCancellablePassthru(backendDialer, &nopRewriter{}))
	addRoutes(router, "/images/{name:.*}/push", newCancellablePassthru(backendDialer, &nopRewriter{}))
	addRoutes(router, "/events", newCancellablePassthru(backendDialer, &nopRewriter{}))
	addRoutes(router, "/containers/{name:.*}/logs", newCancellablePassthru(backendDialer, &nopRewriter{}))
	addRoutes(router, "/containers/{name:.*}/stats", newCancellablePassthru(backendDialer, &nopRewriter{}))

	addRoute(router, "/{any:.*}", newPassthru(backendDialer, &nopRewriter{}))

	return router
}
