package authz

import (
	"net/http"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/garant/authn/enzi"

	"github.com/docker/distribution/context"
	garantauth "github.com/docker/garant/auth"
)

// Authenticate implements the Authenticate method for the Garant
// authn.Authorizer interface. Since it must return the garantauth.Account
// interface type, this function simply wraps the AuthenticateUserRequest
// method.
func (a *authorizer) Authenticate(ctx context.Context, r *http.Request) (garantauth.Account, error) {
	return a.AuthenticateRequestUser(ctx, r)
}

// TODO get to the bottom of this, before I edited this AuthenticateRequestUser could definitely return nil, nil
// AuthenticateRequestUser attempts to authenticate the request and return a
// user object. The DTR authorizer will never return (nil, nil) - the user will
// be an anonymous user object instead.
func (a *authorizer) AuthenticateRequestUser(ctx context.Context, r *http.Request) (user *authn.User, err error) {
	authenticator, err := enzi.NewAuthenticator(a.settingsStore)
	if err != nil {
		return nil, err
	}

	return authenticator.AuthenticateRequestUser(ctx, r)
}
