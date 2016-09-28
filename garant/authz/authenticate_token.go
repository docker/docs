package authz

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/manager/schema"

	"github.com/docker/distribution/context"
	"github.com/docker/garant/auth"
	"github.com/docker/garant/auth/common"

	log "github.com/Sirupsen/logrus"
)

// This file extends authorizer with methods implementing garant.TokenAuthorizer
// which is needed for logging in via tokens from the docker client.

// AuthenticateWithToken attempts to authenticate a HTTP request via an auth
// token
// This is called by Garant's handlers when authenticating with a token.
func (a *authorizer) AuthenticateWithToken(ctx context.Context, token string) (auth.Account, error) {
	ct, err := a.clientTokenMgr.GetClientToken(token)
	if err != nil || ct == nil {
		return nil, &common.ClientError{
			Err:  fmt.Errorf("token not found"),
			Code: http.StatusUnauthorized,
		}
	}

	// Construct a new request which will be used within the enzi authenticator's
	// AuthenticateRequestUser method to determine logging in via a token.
	req := new(http.Request)
	req.Form = url.Values{}
	req.Form.Add("refresh_token", token)

	// Set the client token in the current context so that DTR's enzi authenticator
	// can load the token in its AuthenticateRequestUser method. This will allow
	// the method to get the accountID/refresh_token from the token.
	//
	// By setting this in context and not in the request above we can guarantee
	// that requests can't be forged unless the token is known.
	ctx = context.WithValue(ctx, "clientToken", ct)

	return a.Authenticate(ctx, req)
}

// AuthenticateWithPassword attempts to authenticate a HTTP request via an auth
// token.
// This is called by Garant's handlers when authenticating by password.
func (a *authorizer) AuthenticateWithPassword(ctx context.Context, username, password string) (auth.Account, error) {
	// We're going to use the enzi authenticator to determine whether the given
	// username and password is correct from a docker CLI login.
	// In order to do this we must create a *new* request with basic auth
	// headers set so that the enzi auth can proceed
	r := new(http.Request)
	r.SetBasicAuth(username, password)
	// call AuthenticateRequestUser, which is the same as docker's standard
	// login procedure with basic auth.
	account, err := a.AuthenticateRequestUser(ctx, r)

	if account == nil {
		log.Errorf("AuthenticateRequestUser returned a nil account instead of an anonymous one. This should be impossible!")
	}

	return account, err
}

// GetToken creates a new token for the given client ID
// This is called by Garant's handlers to generate a **new** refresh token after
// a successful authentication request, allowing docker to store the token in
// its credential store.
func (a *authorizer) GetToken(ctx context.Context, acct auth.Account, clientID string) (auth.RefreshToken, error) {
	// At this point acct is the value returned from AuthenticateWithPassword, which
	// has been generated from authenticating via enzi with the user/password combo
	// to generate a refresh token.
	user, ok := acct.(*authn.User)
	if !ok {
		return auth.RefreshToken{}, fmt.Errorf("unknown user type %T", acct)
	}

	t := schema.NewClientToken()
	t.ClientID = clientID
	t.AccountID = user.Account.ID
	if req, err := context.GetRequest(ctx); err == nil {
		t.CreatorIP = req.RemoteAddr
		t.CreatorUA = req.Header.Get("User-Agent")
	}

	if err := a.clientTokenMgr.CreateClientToken(t); err != nil {
		return auth.RefreshToken{}, err
	}

	return t.RefreshToken(), nil
}

// AccountTokens returns all refresh tokens for a given account
func (a *authorizer) AccountTokens(ctx context.Context, acct auth.Account) (tokens []auth.RefreshToken, err error) {
	return
}

// RevokeToken revokes a client token given the unhashed token value
func (a *authorizer) RevokeToken(ctx context.Context, token string) error {
	return nil
}
