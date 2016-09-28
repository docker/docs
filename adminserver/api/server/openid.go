package server

import (
	"net/http"
	"time"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/garant/authn/enzi"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/client/openid"
	"github.com/docker/orca/enzi/api/client/openid/oautherrors"
	"github.com/docker/orca/enzi/jose"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

func (a *APIServer) getOpenIDClient() (*openid.Client, error) {
	haConfig, err := a.settingsStore.HAConfig()
	if err != nil {
		return nil, err
	}

	enziConfig := util.GetEnziConfig(haConfig)
	client, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	if err != nil {
		return nil, err
	}

	return enzi.NewOpenIDClient(client, a.settingsStore)
}

func (a *APIServer) handleOpenIDBegin(ctx context.Context, r *restful.Request) responses.APIResponse {
	openidClient, err := a.getOpenIDClient()
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	authorizeURL, redirectStateKey := openidClient.PrepareAuthorizationRequest(r.Request.URL.Query().Get("next"))

	cookies := []*http.Cookie{{
		Name:     "csrf_redirect_state_key",
		Value:    redirectStateKey,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 5),
	}}
	return responses.JSONResponse(http.StatusFound, map[string][]string{"Location": {authorizeURL}}, cookies, nil)
}

func validateRedirectState(r *http.Request) *openid.RedirectState {
	encodedState := r.URL.Query().Get("state")
	if encodedState == "" {
		log.Errorf("no redirect state in OpenID Connect authorization callback")
		return nil
	}

	cookie, _ := r.Cookie("csrf_redirect_state_key")
	if cookie == nil || cookie.Value == "" {
		log.Errorf("OpenID Connect authorization callback csrf redirect state key cookie not set")
		return nil
	}

	redirectState, err := openid.DecodeRedirectState(encodedState, cookie.Value)
	if err != nil {
		log.Errorf("unable to decode OpenID Connect authorization callback redirect state: %s", err)
		return nil
	}

	return redirectState
}

func (a *APIServer) handleOpenIDCallback(ctx context.Context, r *restful.Request) responses.APIResponse {
	openidClient, err := a.getOpenIDClient()
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	params := r.Request.URL.Query()

	state := validateRedirectState(r.Request)
	if state == nil {
		// TODO where to redirect to?
		return responses.JSONResponse(http.StatusFound, map[string][]string{"Location": {"/"}}, nil, nil)
	}

	if errCode := params.Get("error"); errCode != "" {
		return responses.APIError(errors.OpenIDError(errCode, params.Get("error_description")))
	}

	code := params.Get("code")
	if code == "" {
		log.Errorf("OpenID Connect Provider authorize endpoint returned no code")
		return responses.JSONResponse(http.StatusFound, map[string][]string{"Location": {"/"}}, nil, nil)
	}

	tokenResponse, err := openidClient.GetTokenWithAuthorizationCode(code, state, true)
	if err != nil {
		if oauthErr, ok := err.(*oautherrors.ErrorResponse); ok {
			return responses.APIError(errors.OpenIDError(oauthErr.Code, oauthErr.Description))
		} else {
			log.Errorf("unable to get identity token: %s", err)
			return responses.APIError(errors.OpenIDError("unable to get token", "see server logs for details"))
		}
	}

	// Set browser cookies to expire at some point that is very, very far
	// in the future. Sessions should be managed by users with the auth
	// provider, not their browser, unless they wish to delete cookies.
	expiration := time.Now().Add(time.Hour * 24 * 365 * 5) // 5 years.

	cookies := []*http.Cookie{{
		Name:     "session",
		Value:    tokenResponse.SessionSecret,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  expiration,
	}, {
		Name:     "csrftoken",
		Value:    tokenResponse.SessionCSRFToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: false, // Available to javascript.
		Expires:  expiration,
	}}

	if state.RedirectNext == "" {
		//TODO homepage?
		state.RedirectNext = "/"
	}

	return responses.JSONResponse(http.StatusFound, map[string][]string{"Location": {state.RedirectNext}}, cookies, nil)
}

func (a *APIServer) handleOpenIDKeys(ctx context.Context, r *restful.Request) responses.APIResponse {
	enziKey, err := a.settingsStore.EnziSigningKey()
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	signingKey, err := jose.NewPrivateKey(enziKey.CryptoPrivateKey())
	if err != nil {
		return responses.APIError(errors.InternalError(ctx, err))
	}

	header := map[string][]string{
		"Cache-Control": {"max-age=3600"},
	}
	return responses.JSONResponse(http.StatusOK, header, nil, responses.OpenIDKeys{Keys: []jose.PublicKey{signingKey.PublicKey}})
}
