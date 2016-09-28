package pipeline

import (
	"net/http"
	"net/http/httptest"

	"github.com/docker/orca/controller/ctx"
)

// MockTestServer creates an httptest server for a given Handler
// All requests are performed as the admin user
func MockTestServer(h Handler) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h(w, ctx.MockAdmin(r))
		}))
}
