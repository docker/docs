package auth

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/garant/config"
)

// Account represents an account object. Its underlying type depends on the
// Authorizer backend.
type Account interface {
	// Subject should supply a value which will be used in the subject field of
	// a JSON Web Token, e.g., username, email, keyID.
	Subject() string
	// Info should return basic info (at the discretion of the backend) about
	// the given account. This will typically include fields like the account
	// name, admin status, etc.
	Info() map[string]interface{}
}

// Challenge is a special error type which is used for HTTP 401 Unauthorized
// responses and is able to write the response with WWW-Authenticate challenge
// header values based on the error.
type Challenge interface {
	error
	// ServeHTTP prepares the request to conduct the appropriate challenge
	// response. For most implementations, simply calling ServeHTTP should be
	// sufficient. Because no body is written, users may write a custom body
	// after calling ServeHTTP, but any headers must be written before the call
	// and may be overwritten.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Resource describes a resource by type and name.
type Resource struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Access describes a specific action that is
// requested or allowed for a given resource.
type Access struct {
	Resource
	Action string `json:"action"`
}

// Authorizer is the interface for account authentication and authorization.
type Authorizer interface {
	// Authenticate should attempt to authenticate an account for the given
	// HTTP Request. If the request does not attempt authentication (i.e. an
	// anonymous request), the Authorizer should return a nil account and a nil
	// error. If the request does attempt authentication but fails (e.g.,
	// invalid username/password) then the Authorizer should return a non-nil
	// error. The returned account is non-nil and error is nil iff
	// authentication succeeds. If the returned error is meant to indicate an
	// error by the client (e.g., status code 401) the error type must
	// implement the Challenge interface. Any other errors will be interpreted
	// as server errors.
	Authenticate(ctx context.Context, r *http.Request) (Account, error)
	// Authorize should attempt to authorize the given account for the
	// requested access to resources hosted by the specified service. The
	// authorizer should return a subset of the requested access set which
	// represents the access which the account has been granted. The authorizer
	// should ignore access requests for resources that do not exist or for
	// actions that do not exist in the context of the resource type. Any
	// non-nil error returned will be interpreted as a server error.
	Authorize(ctx context.Context, acct Account, service string, requestedAccess ...Access) (grantedAccess []Access, err error)
}

// TokenAuthorizer is the interface for OAuth-compatible token authentication
// and token administration.
type TokenAuthorizer interface {
	Authorizer
	// AuthenticateWithToken attempts to authenticate an account from a provided
	// refresh token.
	AuthenticateWithToken(ctx context.Context, token string) (Account, error)
	// AuthenticateWithPassword attempts to authenticate an account from the
	// provided username and password, as defined in the OAuth-compatible token
	// flow specification.
	AuthenticateWithPassword(ctx context.Context, username, password string) (Account, error)
	// GetToken returns the token associated with the given account and clientID
	// combination or generates a new one if none exists.
	GetToken(ctx context.Context, acct Account, clientID string) (RefreshToken, error)
	// AccountTokens returns a list of all refresh tokens associated with the
	// provided user account.
	AccountTokens(ctx context.Context, acct Account) ([]RefreshToken, error)
	// RevokeToken invalidates and removes the provided token from the backend
	// storage, preventing future use with AuthenticateWithToken and
	// AccountTokens.
	RevokeToken(ctx context.Context, token string) error
}

// RefreshToken represents an auditable long-lived authentication token.
type RefreshToken struct {
	Token    string
	ClientID string
}

// InitFunc is the type of an Authorizer factory function and is used to
// register the constructor for different Authorizer backends.
type InitFunc func(params config.Parameters) (Authorizer, error)

var authorizers map[string]InitFunc

func init() {
	authorizers = make(map[string]InitFunc)
}

// Register is used to register an InitFunc for an Authorizer backend with the
// given name.
func Register(name string, initFunc InitFunc) error {
	if _, exists := authorizers[name]; exists {
		return fmt.Errorf("name already registered: %s", name)
	}

	authorizers[name] = initFunc

	return nil
}

// NewAuthorizer constructs an Authorizer with the given parameters using the
// backend registered with the given name.
func NewAuthorizer(name string, parameters config.Parameters) (Authorizer, error) {
	if initFunc, exists := authorizers[name]; exists {
		return initFunc(parameters)
	}

	return nil, fmt.Errorf("no authorizers registered with name: %s", name)
}
