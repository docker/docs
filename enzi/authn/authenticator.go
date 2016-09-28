package authn

import (
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
)

// ErrInvalidUsernamePassword returns an error indicating that an
// authentication attempt has an invalid username or password.
func ErrInvalidUsernamePassword() *errors.APIError {
	return errors.InvalidAuthentication("incorrect username or password")
}

// ErrAccountInactive returns an error indicating that an authentication
// attempt is invalid because the account in not active.
func ErrAccountInactive() *errors.APIError {
	return errors.InvalidAuthentication("account is inactive")
}

// RequestAuthenticator is able to authenticate the request using multiple
// different methods.
type RequestAuthenticator interface {
	// AuthenticateRequest should attempt to authenticate an account for
	// the given HTTP Request. If the request does not attempt
	// authentication (i.e. an anonymous request), the Authenticator should
	// return Account with IsAnonymous set to true, and a nil error. The
	// Authenticator should attempt to authenticate the request in any as
	// many ways as it is able in whatever order it sees fit until either
	// authentication succeeds or an error is encountered in which case it
	// should return without trying successive authentication methods. If
	// the returned error is meant to indicate an error by the client
	// (e.g., status code 401) the error should have an HTTP status code of
	// 401 (Unauthorized).
	AuthenticateRequest(ctx context.Context, r *http.Request) (*Account, *errors.APIError)
}

// UsernamePasswordAuthenticator is able to authenticate a client with a
// username and password.
type UsernamePasswordAuthenticator interface {
	// AuthenticateUsernamePassword follows the same semantics as
	// AuthenticateRequest in the Authenticator interface but only attempts
	// basic authentication with a username and password. The
	// implementation should not check if the account is active, that is
	// left to the root authenticator.
	AuthenticateUsernamePassword(ctx context.Context, username, password string) (*Account, *errors.APIError)
}

// OpenIDTokenAuthenticator is able to authenticate a client with an OpenID
// Token.
type OpenIDTokenAuthenticator interface {
	// AuthenticateOpenIDToken follows the same semantics as
	// AuthenticateRequest in the Authenticator interface but only attempts
	// openid token authentication which authenticates the client as a
	// service acting on behalf of an account. The implementation should
	// not check if the account is active, that is left to the root
	// authenticator.
	AuthenticateOpenIDToken(ctx context.Context, token string) (*Account, *errors.APIError)
}

// SessionTokenAuthenticator is able to authenticate a client with a session
// token.
type SessionTokenAuthenticator interface {
	// AuthenticateSessionToken follows the same semantics as
	// AuthenticateRequest in the Authenticator interface but only attempts
	// session token authentication. If the session is invalid the returned
	// error should be ErrInvalidSession. The implementation should not
	// check if the account is active, that is left to the root
	// authenticator.
	AuthenticateSessionToken(ctx context.Context, token string) (*Account, *errors.APIError)
}
