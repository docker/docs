package enzi

import (
	"crypto"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	apierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
)

func (a *Authenticator) AuthenticateUsernamePassword(username, password string) (*auth.Context, error) {
	// Make a session with nil creds (unauthenticated).
	session := a.getSession(nil)

	loginForm := forms.Login{
		Username: username,
		Password: password,
	}

	loginSessionResp, err := session.Login(loginForm)
	if err != nil {
		// The error could be an eNZi API error.
		apiErrs, ok := err.(*apierrors.APIErrors)
		if ok && apiErrs.HTTPStatusCode == http.StatusBadRequest {
			return nil, auth.ErrInvalidPassword
		}

		log.Errorf("unable to login user on auth provider: %s", err)

		return nil, fmt.Errorf("unable to login user on auth provider: %s", err)
	}

	// We can use a LoginSession eNZi API response instance as a Request
	// Authenticator for the eNZi API client.
	clientCreds := loginSessionResp
	user, err := a.populateUserFields(clientCreds, &loginSessionResp.Account)
	if err != nil {
		return nil, err
	}

	return &auth.Context{
		User:         user,
		SessionToken: loginSessionResp.SessionToken,
		ClientCreds:  clientCreds,
	}, nil
}

func (a *Authenticator) AuthenticateSessionToken(sessionToken string) (*auth.Context, error) {
	tokenResp, err := a.openidClient.GetTokenWithRootSession(sessionToken)
	if err != nil {
		log.Errorf("unable to authenticate user with session token on auth provider: %s", err)

		return nil, fmt.Errorf("unable to authenticate user with session token on auth provider: %s", err)
	}

	// We could use the token response instance as a Request Authenticator
	// that uses the OpenID Token, but we want to keep using the session
	// token so we can just create a login session authenticator instead.
	clientCreds := &responses.LoginSession{
		Account:      *tokenResp.Account,
		SessionToken: sessionToken,
	}

	user, err := a.populateUserFields(clientCreds, tokenResp.Account)
	if err != nil {
		return nil, err
	}

	return &auth.Context{
		User: user,
		// Keep the session token around in case we need it.
		SessionToken: sessionToken,
		ClientCreds:  clientCreds,
	}, nil
}

func (a *Authenticator) AuthenticatePublicKey(publicKey crypto.PublicKey) (*auth.Context, error) {
	accountKey, err := a.getPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	// Use the user ID from the public key as a refresh token to get an
	// OpenID token.
	tokenResp, err := a.openidClient.GetTokenWithAccountID(accountKey.UserID)
	if err != nil {
		log.Errorf("unable to authenticate user with refresh token on auth provider: %s", err)

		return nil, fmt.Errorf("unable to authenticate user with refresh token on auth provider: %s", err)
	}

	// We can use the token response instance as a Request Authenticator
	// that uses the OpenID Token.
	clientCreds := tokenResp
	user, err := a.populateUserFields(clientCreds, tokenResp.Account)
	if err != nil {
		return nil, err
	}

	return &auth.Context{
		User: user,
		// No Session Token used when authenticating via PublicKey.
		SessionToken: "",
		ClientCreds:  clientCreds,
	}, nil
}

func (a *Authenticator) Logout(ctx *auth.Context) error {
	session := a.getSession(ctx.ClientCreds)

	if err := session.Logout(); err != nil {
		return fmt.Errorf("unable to logout on auth provider: %s", err)
	}

	return nil
}
