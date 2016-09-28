package ctx

import (
	"github.com/docker/orca/auth"
	"net/http"
)

// MockUser(account, req) returns a mock request context representing a user
func MockUser(account *auth.Account, req *http.Request) *OrcaRequestContext {
	rc := &OrcaRequestContext{
		Auth:    &auth.Context{User: account},
		Request: req,
	}
	rc.ParseVars()
	return rc
}

// MockAdmin(req) returns a mock request context representing admin
func MockAdmin(req *http.Request) *OrcaRequestContext {
	admin := &auth.Account{
		Username: "admin",
		Role:     auth.Admin,
		Admin:    true,
	}
	return MockUser(admin, req)
}
