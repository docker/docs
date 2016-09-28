package authn

import (
	"net/http"

	"github.com/docker/distribution/context"
)

// Authenticator is a subset of the Garant Authorizer interface which only
// implements the Authenticate method. This will allow us to authenticate
// using multiple different methods (hashed password in DB, ldap bind, etc.)
// and still share the same Authorize() implementation.
type Authenticator interface {
	// AuthenticateRequestUser should attempt to authenticate an account for
	// the given HTTP Request. If the request does not attempt authentication
	// (i.e. an anonymous request), the Authenticator should return a nil User
	// and a nil error. If the request does attempt authentication but fails
	// (e.g., invalid username/password) then the Authenticator should return a
	// non-nil error. The returned User is non-nil and error is nil iff
	// authentication succeeds. If the returned error is meant to indicate an
	// error by the client (e.g., status code 401) the error type must
	// implement the Challenge interface. Any other errors will be interpreted
	// as server errors.
	AuthenticateRequestUser(ctx context.Context, r *http.Request) (*User, error)
}
