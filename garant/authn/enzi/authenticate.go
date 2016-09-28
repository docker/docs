package enzi

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/libtrust"

	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/garant/auth/common"
	enziclient "github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/client/openid"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/jose"
)

const CSRFTokenHeaderName = "X-Csrf-Token"

type authenticator struct {
	httpClient *http.Client
	enziConfig util.EnziConfig

	openidClient *openid.Client

	settingsStore hubconfig.SettingsStore
}

// NewAuthenticator creates an eNZi authenticator with the given settings.
func NewAuthenticator(settingsStore hubconfig.SettingsStore) (authn.Authenticator, error) {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return nil, err
	}

	enziConfig := util.GetEnziConfig(haConfig)
	client, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	if err != nil {
		return nil, err
	}

	openidClient, err := NewOpenIDClient(client, settingsStore)
	return &authenticator{client, enziConfig, openidClient, settingsStore}, err
}

type openIDClientKey struct {
	signingKey     string
	registrationID string
	redirectURI    string
	enziHost       string
	authorizePath  string
	tokenPath      string
}

var openIDClientCache = map[openIDClientKey]*openid.Client{}

func NewOpenIDClient(httpClient *http.Client, settingsStore hubconfig.SettingsStore) (*openid.Client, error) {
	userHubconfig, err := settingsStore.UserHubConfig()
	if err != nil {
		return nil, err
	}

	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return nil, err
	}

	enziConfig := util.GetEnziConfig(haConfig)

	redirectURI := fmt.Sprintf("https://%s/api/v0/openid/callback", userHubconfig.DTRHost)

	serviceRegistration, err := settingsStore.EnziService()
	if err != nil {
		return nil, err
	}

	enziKey, err := settingsStore.EnziSigningKey()
	if err != nil {
		return nil, err
	}

	signingKey, err := jose.NewPrivateKey(enziKey.CryptoPrivateKey())
	if err != nil {
		return nil, err
	}

	registrationID := serviceRegistration.ID
	authorizePath := fmt.Sprintf("%s/authorize", enziConfig.Prefix)
	tokenPath := fmt.Sprintf("%s/v0/id/token", enziConfig.Prefix)

	// Every time we create a new openid client it creates each HTTP client caches connections
	// and eventually (after 65536 requests) we run out of source IP addresses for making new
	// connections to enzi. We don't know why the HTTP clients are not being cleaned up
	// properly after they are used, so instead of investigating that we are adding a cache
	// so that we don't create more openid clients than necessary. We can't switch to a model
	// where we create this only once at start up because it would require the admin servers
	// to watch for config changes made by other admin servers. I think that's the right long
	// term solution, but this somewhat hacky cache should suffice for the purposes of
	// making a stable 2.0.0 release on time.
	// See #1917 for how we discovered this issue.
	rawEnziKey, err := settingsStore.RawEnziSigningKey()
	if err != nil {
		return nil, err
	}
	key := openIDClientKey{
		signingKey:     rawEnziKey,
		registrationID: registrationID,
		redirectURI:    redirectURI,
		enziHost:       enziConfig.Host,
		authorizePath:  authorizePath,
		tokenPath:      tokenPath,
	}

	client, ok := openIDClientCache[key]
	if !ok {
		client, err = openid.NewClient(httpClient, signingKey, registrationID, redirectURI, enziConfig.Host, authorizePath, tokenPath), nil
		if err != nil {
			return nil, err
		}
		openIDClientCache[key] = client
	}
	return client, nil
}

func (a *authenticator) MakeAnonymousUser() *authn.User {
	username := "anonymous"
	isAdmin := false
	isActive := false

	return &authn.User{
		IsAnonymous: true,
		Account: &enziresponses.Account{
			Name:     username,
			IsAdmin:  &isAdmin,
			IsActive: &isActive,
		},
		EnziSession: enziclient.New(a.httpClient, a.enziConfig.Host, a.enziConfig.Prefix, nil),
	}
}

// AuthenticateRequestUser should attempt to authenticate an account for
// the given HTTP Request. If the request does not attempt authentication
// (i.e. an anonymous request), the Authenticator should return an anonymous User
// and a nil error. If the request does attempt authentication but fails
// (e.g., invalid username/password) then the Authenticator should return a
// non-nil error and an anonymous user. If the returned error is meant to indicate an
// error by the client (e.g., status code 401) the error type must
// implement the Challenge interface. Any other errors will be interpreted
// as server errors.
func (a *authenticator) AuthenticateRequestUser(ctx context.Context, r *http.Request) (*authn.User, error) {
	var tokenResponse *openid.TokenResponse
	var err error

	if user, pass, ok := r.BasicAuth(); ok {
		bypass, err := authn.ShouldBypassAuth(r)
		if err != nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("encountered an error during ucp auth bypass: %s", err),
				Code: http.StatusUnauthorized,
			}
		}

		if bypass {
			tokenResponse, err = a.openidClient.GetTokenWithRefreshToken(user)
			if err != nil {
				return a.MakeAnonymousUser(), &common.ClientError{
					Err:  fmt.Errorf("unable to get identity token from refresh token in ucp auth bypass: %s", err),
					Code: http.StatusUnauthorized,
				}
			}
		} else {
			tokenResponse, err = a.openidClient.GetTokenWithUsernamePasword(user, pass, false)
			if err != nil {
				return a.MakeAnonymousUser(), &common.ClientError{
					Err:  fmt.Errorf("unable to get identity token from basic auth credentials: %s", err),
					Code: http.StatusUnauthorized,
				}
			}
		}

		if tokenResponse == nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("received a nil tokenResponse from basic auth credentials"),
				Code: http.StatusInternalServerError,
			}
		}
	} else if parts := strings.Split(r.Header.Get("Authorization"), " "); len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		tokenResponse, err = a.GetTokenResponseFromRawToken(ctx, parts[1])
		if err != nil {
			return a.MakeAnonymousUser(), err
		}

		if tokenResponse == nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("received a nil tokenResponse from service session"),
				Code: http.StatusInternalServerError,
			}
		}
	} else if sessionCookie, _ := r.Cookie("session"); sessionCookie != nil && sessionCookie.Value != "" {
		tokenResponse, err = a.openidClient.GetTokenWithServiceSession(sessionCookie.Value)
		if err != nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("unable to get identity token from service session: %s", err),
				Code: http.StatusUnauthorized,
			}
		}

		if tokenResponse == nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("received a nil tokenResponse from service session"),
				Code: http.StatusInternalServerError,
			}
		}
	} else if token := r.FormValue("refresh_token"); token != "" {
		// If the request has no basic auth check for the refresh_token parameter
		// passed by logging in via the docker CLI using a token instead of a password.
		//
		// When this path is hit this authenticate function is being called from
		// AuthenticateWithToken which has already loaded the ClientToken from DTR and
		// has set the token in the current context.
		token, ok := ctx.Value("clientToken").(*schema.ClientToken)
		if !ok {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("unable to get client token from context"),
				Code: http.StatusInternalServerError,
			}
		}
		// Note that when using an ID as a refresh token it needs an "id:" prefix
		tokenResponse, err = a.openidClient.GetTokenWithRefreshToken("id:" + token.AccountID)
		if err != nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("unable to get identity token from client token: %s", err),
				Code: http.StatusUnauthorized,
			}
		}
		if tokenResponse == nil {
			return a.MakeAnonymousUser(), &common.ClientError{
				Err:  fmt.Errorf("received a nil tokenResponse from service session"),
				Code: http.StatusInternalServerError,
			}
		}
	} else {
		return a.MakeAnonymousUser(), nil
	}

	if tokenResponse.Account == nil {
		// Handle like an unauthenticated request.
		return a.MakeAnonymousUser(), nil
	} else if tokenResponse.Account.IsAdmin == nil {
		return a.MakeAnonymousUser(), &common.ClientError{
			Err:  fmt.Errorf("received a tokenResponse for an org while authenticating a user"),
			Code: http.StatusInternalServerError,
		}
	}

	if tokenResponse.ErrorResponse != nil && tokenResponse.ErrorResponse.Code != "" {
		return a.MakeAnonymousUser(), &common.ClientError{
			Err:  fmt.Errorf("token request error: %s - %s", tokenResponse.ErrorResponse.Code, tokenResponse.ErrorResponse.Description),
			Code: http.StatusUnauthorized,
		}
	}

	csrfToken := tokenResponse.SessionCSRFToken
	// We don't require csrf for:
	// 1. the database proxy
	// 2. GET requests
	// 3. HEAD requests
	// We accept the CSRF token as:
	// Header: CSRFTokenHeaderName
	// Form value "csrftoken"
	if !((r.URL != nil && strings.HasPrefix(r.URL.Path, "/db")) || r.Method == "GET" || r.Method == "HEAD" || csrfToken == r.Header.Get(CSRFTokenHeaderName) || csrfToken == r.FormValue("csrftoken")) {
		// Handle like an unauthenticated request.
		return a.MakeAnonymousUser(), nil
	}

	return &authn.User{
		EnziSession: enziclient.New(a.httpClient, a.enziConfig.Host, a.enziConfig.Prefix, tokenResponse),
		Token:       tokenResponse,
		Account:     tokenResponse.Account,
	}, nil
}

func (a *authenticator) GetTokenResponseFromRawToken(ctx context.Context, rawToken string) (*openid.TokenResponse, error) {
	hubConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		return nil, err
	}

	authConfig := util.GetRegistryAuthConfig(hubConfig.DTRHost)

	// check issuer - TODO replace this with jwt.go from enzi once implemented (with auth bypass bypass)
	jwtToken, err := token.NewToken(rawToken)
	if err != nil {
		return nil, fmt.Errorf("could not decode jwt token: %s", err)
	}

	if jwtToken.Claims.Issuer != authConfig.Issuer {
		return nil, fmt.Errorf("the token issuer was not DTR")
	}

	certBundle, err := a.settingsStore.GarantRootCert()
	if err != nil {
		return nil, fmt.Errorf("could not get garant cert bundle: %s", err)
	}

	rawCertBundle := []byte(certBundle)
	var rootCerts []*x509.Certificate
	pemBlock, rawCertBundle := pem.Decode(rawCertBundle)
	for pemBlock != nil {
		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse token auth root certificate: %s", err)
		}

		rootCerts = append(rootCerts, cert)

		pemBlock, rawCertBundle = pem.Decode(rawCertBundle)
	}

	if len(rootCerts) == 0 {
		return nil, errors.New("token auth requires at least one token signing root certificate")
	}

	rootPool := x509.NewCertPool()
	trustedKeys := make(map[string]libtrust.PublicKey, len(rootCerts))
	for _, rootCert := range rootCerts {
		rootPool.AddCert(rootCert)
		pubKey, err := libtrust.FromCryptoPublicKey(crypto.PublicKey(rootCert.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("unable to get public key from token auth root certificate: %s", err)
		}
		trustedKeys[pubKey.KeyID()] = pubKey
	}

	verifyOpts := token.VerifyOptions{
		TrustedIssuers:    []string{authConfig.Issuer},
		AcceptedAudiences: []string{authConfig.Service},
		Roots:             rootPool,
		TrustedKeys:       trustedKeys,
	}

	if err = jwtToken.Verify(verifyOpts); err != nil {
		return nil, fmt.Errorf("Failed to verify jwt token: %s", err)
	}

	token, err := a.openidClient.GetTokenWithRefreshToken(jwtToken.Claims.Subject)
	if err != nil {
		return nil, &common.ClientError{
			Err:  fmt.Errorf("unable to get identity token from refresh token: %s", err),
			Code: http.StatusUnauthorized,
		}
	}

	return token, nil
}
