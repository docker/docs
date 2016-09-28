package openid

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/client/openid"
	"github.com/docker/orca/enzi/api/client/openid/oautherrors"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/internal/helpers"
	"github.com/docker/orca/enzi/authn"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
)

func (s *Service) routeToken() server.Route {
	return server.Route{
		Method:  "POST",
		Path:    "/token",
		Handler: s.handleToken,
		Doc:     "Get an Identity Token",
		Notes: `An identity token allows a service to authenticate to
			this auth service and to other services on behalf of an
			account. These tokens are JWTs which can be verified by
			one of the keys in use by this auth service which are
			advertised by the Signing Keys endpoint. Usage and
			acceptance of the token depends on the values of the
			claimset of the token.
			` + "\n" +
			`The "iss" value is the issuer of the token which is
			always the Issuer Identifier of this auth service.
			` + "\n" +
			`The "azp" value is identity of the authorized party and
			corresponds to the ID of the service to which the
			token was issued.
			` + "\n" +
			`The "sub" value is the subject of the token which will
			be the ID of the account for which the token grants
			access.
			` + "\n" +
			`The "aud" value is an array of identifiers for those
			services which are the intended audience of the token.
			Audience values may include the domain name of this
			auth service in order to access the accounts API on
			behalf of the subject account or may be an ID of any
			other service to which the authorized party needs to
			authenticate on behalf of the subject account.
			` + "\n" +
			`The "iat" value indicates the time at which the JWT
			was issued. Its value is a JSON number representing the
			number of secconds since the UNIX epoch.
			` + "\n" +
			`The "exp" value indicates the time after which the
			token must no longer be accepted. Its value is a JSON
			number representing the number of seconds since the
			UNIX epoch. The tokens issued by this endpoint are
			relatively short-lived (10 minutes) but once an account
			has authorized a service, that service may request
			additional tokens as needed as long as that
			authorization has not been revoked.
			` + "\n" +
			`There are several grant types which may be used to
			request a token:
			` + "\n" +
			`The "authorization_code" grant type is used for an
			account to signup or login to a service. See the notes
			for the Oauth Authorize endpoint for how to get an
			authorization code. This grant type requires setting
			the "code" and "redirect_uri" form parameters. If the
			"link_session" form parameter is set to a "true" value,
			a service session will be created and associated with
			the session by the user who performed the initial
			authorization request.
			` + "\n" +
			`The "service_session" grant type can then be used by
			including a "session_secret" form parameter which
			corresponds to a valid session by a user.
			` + "\n" +
			`The "refresh_token" grant type can be used by any
			service which has already been granted authorization by
			an account. The "refresh_token" form parameter value is
			set to the ID of the account you wish to get an ID
			token for.
			` + "\n" +
			`All grant types require that a service authenticate to
			the endpoint by providing a signed JWT with the service
			ID as both the "iss" and "sub" value. It must be singed
			by a key which is advertized at the service's
			registered JWKs URI`,
		Consumes: []string{"application/x-www-form-urlencoded"},
		FormParameterDocs: []*restful.Parameter{
			restful.FormParameter("client_assertion_type", "Services should set this value to 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer' to use a JWT Bearer Token for client authentication."),
			restful.FormParameter("client_assertion", "Services should set this value to a single JWT to use for client authentication; See endpoint notes for more details."),
			restful.FormParameter("grant_type", "The Oauth 2.0 grant type being used; Required."),
			restful.FormParameter("code", "The authorization code; Required if 'grant_type' is 'authorization_code'."),
			restful.FormParameter("redirect_uri", "The 'redirect_uri' parameter used in the authorization request; Required if 'grant_type' is 'authorization_code'."),
			restful.FormParameter("link_session", "Whether to create a shared session linked to the session used during the initial authorization request. Applies only to 'authorization_code' and 'password' grant types. Accepted values are: '1', 't', 'T', 'TRUE', 'true', 'True'. All other values are treated as 'false'."),
			restful.FormParameter("username", "The username; Required if 'grant_type' is 'password'."),
			restful.FormParameter("password", "The password; Required if 'grant_type' is 'password'."),
			restful.FormParameter("refresh_token", "The refresh token; Required if 'grant_type' is 'refresh_token'. This value is usually set to an account ID."),
			restful.FormParameter("audience", "When using the refresh_token grant type, this value may be set to the ID of another service to which the authorizing party wishes to authenticate on behalf of the subject account."),
			restful.FormParameter("session_secret", "The service session secret; Required if 'grant_type' is 'service_session' or 'root_session'."),
		},
		Produces: []string{restful.MIME_JSON},
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, Token response returned.",
				Sample:  openid.TokenResponse{},
			},
		},
	}
}

func (s *Service) handleToken(ctx context.Context, r *restful.Request) responses.APIResponse {
	if err := r.Request.ParseForm(); err != nil {
		return responses.JSONResponse(http.StatusBadRequest, oautherrors.InvalidRequest(fmt.Errorf("unable to parse request body: %s", err)))
	}

	grantType := r.Request.PostFormValue("grant_type")
	if grantType == "" {
		return responses.JSONResponse(http.StatusBadRequest, oautherrors.InvalidRequest(fmt.Errorf("form parameter 'grant_type' MUST be specified")))
	}

	// Authenticate the Client (service) using the 'private_key_jwt'
	// client authentication method.
	service, err := s.authenticateService(ctx, r)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return responses.APIError(apiErr)
		}

		// Otherwise, err is of type oautherrors.Response
		return responses.JSONResponse(http.StatusBadRequest, err)
	}

	var (
		account            *schema.Account
		serviceSession     *schema.ServiceSession
		rootSession        *schema.Session
		additionalAudience string
	)

	// Handle the grant type.
	switch grantType {
	case "authorization_code":
		account, serviceSession, err = s.handleAuthorizationCodeGrant(ctx, r, service)
	case "password":
		account, rootSession, err = s.handlePasswordGrant(ctx, r, service)
	case "root_session":
		account, rootSession, err = s.handleRootSessionGrant(ctx, r, service)
	case "service_session":
		account, serviceSession, err = s.handleServiceSessionGrant(ctx, r, service)
	case "refresh_token":
		account, additionalAudience, err = s.handleRefreshTokenCodeGrant(ctx, r, service)
	default:
		return responses.JSONResponse(http.StatusBadRequest, oautherrors.UnsupportedGrantType(fmt.Errorf("unsupported grant type: %q")))
	}

	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return responses.APIError(apiErr)
		}

		return responses.JSONResponse(http.StatusBadRequest, oautherrors.InvalidGrant(err))
	}

	openIDConfig, err := config.GetOpenIDConfig(s.schemaMgr)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, fmt.Errorf("unable to get issuer identifier: %s", err)))
	}

	now := time.Now()
	claimSet := jose.JWTClaims{
		Issuer:            openIDConfig.IssuerIdentifier,
		Subject:           account.ID,
		Audience:          []string{openIDConfig.IssuerIdentifier}, // ID tokens can always be presented back to the issuer.
		AuthorizedParty:   service.ID,
		IssuedAt:          now.Unix(),
		Expiration:        now.Add(10 * time.Minute).Unix(),
		Name:              account.FullName,
		PreferredUsername: account.Name,
	}

	if additionalAudience != "" {
		// The additional audience will the ID of a service which the
		// authorized party wishes to authenticate to on behalf of the
		// subject account.
		claimSet.Audience = append(claimSet.Audience, additionalAudience)
	}

	idToken, err := jose.NewJWT(s.signingKey, claimSet)
	if err != nil {
		return responses.APIError(errors.Internal(ctx, fmt.Errorf("unable to make ID token: %s", err)))
	}

	responseHeader := http.Header{
		"Cache-Control": []string{"no-store"},
		"Pragma":        []string{"no-cache"},
	}

	accountResponseObject := responses.MakeAccount(account)

	responseBody := openid.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  idToken.String(),
		IDToken:      idToken.String(),
		ExpiresIn:    idToken.ExpiresIn(),
		RefreshToken: account.ID,
		Account:      &accountResponseObject,
	}

	if serviceSession != nil {
		responseBody.SessionSecret = serviceSession.Secret
		responseBody.SessionCSRFToken = serviceSession.CSRFToken
	} else if rootSession != nil {
		responseBody.SessionSecret = rootSession.Secret
		responseBody.SessionCSRFToken = rootSession.CSRFToken
	}

	return responses.JSONResponseWithHeaders(http.StatusOK, responseBody, responseHeader)
}

func (s *Service) handleAuthorizationCodeGrant(ctx context.Context, r *restful.Request, service *schema.Service) (*schema.Account, *schema.ServiceSession, error) {
	// Get authorization code.
	rawCode := r.Request.PostFormValue("code")
	if rawCode == "" {
		return nil, nil, fmt.Errorf("form parameter 'code' MUST be specified")
	}

	redirectURI := r.Request.PostFormValue("redirect_uri")
	if redirectURI == "" {
		return nil, nil, fmt.Errorf("form parameter 'redirect_uri' MUST be specified")
	}

	authCode, err := s.schemaMgr.GetServiceAuthCode(rawCode)
	if err != nil && err != schema.ErrNoSuchServiceAuthCode {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get service auth code: %s", err))
	}

	if err == schema.ErrNoSuchServiceAuthCode || authCode.ServiceID != service.ID {
		return nil, nil, fmt.Errorf("no such auth code issued to service: %q", rawCode)
	}

	defer s.schemaMgr.DeleteServiceAuthCode(authCode.ID)

	// Check that the auth code hasn't expired.
	now := time.Now()
	if now.After(authCode.Expiration) {
		return nil, nil, fmt.Errorf("auth code expired at %d - current time is %d", authCode.Expiration.Unix(), now.Unix())
	}

	// Check that the authorization used the same redirect_uri value.
	if authCode.RedirectURI != redirectURI {
		return nil, nil, fmt.Errorf("form parameter 'redirect_uri' does not match that which was used in the initial authorization request")
	}

	account, err := s.schemaMgr.GetAccountByID(authCode.AccountID)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			return nil, nil, fmt.Errorf("account associated with authorization code no longer exists")
		}

		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get account: %s", err))
	}

	if !account.IsActive {
		apiErr := authn.ErrAccountInactive()
		return nil, nil, fmt.Errorf(apiErr.Message)
	}

	// Determine if we need to create a shared service session which is
	// linked to the session used during the initial authorization request.
	linkSession, _ := strconv.ParseBool(r.Request.PostFormValue("link_session"))
	if !linkSession {
		// Do not attempt to create a service session.
		return account, nil, nil
	}

	if authCode.SessionID == "" {
		return nil, nil, fmt.Errorf("unable to link session: authorization code grants no session access")
	}

	rootSession, err := s.schemaMgr.GetSession(authCode.SessionID)
	if err != nil {
		if err == schema.ErrNoSuchSession {
			// No session to link to.
			return nil, nil, fmt.Errorf("unable to link session: session for initial authorization request no longer exists")
		}

		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get session linked to authorization code: %s", err))
	}

	if rootSession.UserID != account.ID {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("authorization code account does not match its session"))
	}

	serviceSession, err := s.schemaMgr.CreateServiceSession(rootSession.ID, service.ID)
	if err != nil {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to create linked service session: %s", err))
	}

	return account, serviceSession, nil
}

func (s *Service) handlePasswordGrant(ctx context.Context, r *restful.Request, service *schema.Service) (*schema.Account, *schema.Session, error) {
	// The service needs to be a priviledged service.
	if !service.Privileged {
		return nil, nil, fmt.Errorf("must be a privileged 1st party service to use password grant")
	}

	username := r.Request.PostFormValue("username")
	if username == "" {
		return nil, nil, fmt.Errorf("form parameter 'username' MUST be specified")
	}

	password := r.Request.PostFormValue("password")
	if password == "" {
		return nil, nil, fmt.Errorf("form parameter 'password' MUST be specified")
	}

	authConfig, err := s.authorizer.AuthConfig()
	if err != nil {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get current auth config: %s", err))
	}

	var ldapSettings *ldapconfig.Settings
	if authConfig.Backend == config.AuthBackendLDAP {
		ldapSettings, err = s.authorizer.LDAPSettings()
		if err != nil {
			return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get current LDAP config: %s", err))
		}
	}

	authenticator := s.authorizer.UsernamePasswordAuthenticator(ldapSettings)

	user, apiErr := authenticator.AuthenticateUsernamePassword(ctx, username, password)
	if apiErr != nil {
		if apiErr.IsInternal() {
			return nil, nil, apiErr
		}

		return nil, nil, fmt.Errorf(apiErr.Message)
	}

	// Determine if we need to create a root session for this user.
	linkSession, _ := strconv.ParseBool(r.Request.PostFormValue("link_session"))
	if !linkSession {
		// Do not attempt to create a service session.
		return &user.Account, nil, nil
	}

	session, err := s.authorizer.SessionTokenAuthenticator().CreateSession(user)
	if err != nil {
		return nil, nil, errors.Internal(ctx, err)
	}

	return &user.Account, session, nil
}

func (s *Service) handleRootSessionGrant(ctx context.Context, r *restful.Request, service *schema.Service) (*schema.Account, *schema.Session, error) {
	// The service needs to be a priviledged service.
	if !service.Privileged {
		return nil, nil, fmt.Errorf("must be a privileged 1st party service to use root session grant")
	}

	sessionSecret := r.Request.PostFormValue("session_secret")
	if sessionSecret == "" {
		return nil, nil, fmt.Errorf("form parameter 'session_secret' MUST be specified")
	}

	sessionAuthenticator := s.authorizer.SessionTokenAuthenticator()

	user, apiErr := sessionAuthenticator.AuthenticateSessionToken(ctx, sessionSecret)
	if apiErr != nil {
		if apiErr.IsInternal() {
			return nil, nil, apiErr
		}

		return nil, nil, fmt.Errorf(apiErr.Message)
	}

	return &user.Account, user.Session, nil
}

func (s *Service) handleServiceSessionGrant(ctx context.Context, r *restful.Request, service *schema.Service) (*schema.Account, *schema.ServiceSession, error) {
	sessionSecret := r.Request.PostFormValue("session_secret")
	if sessionSecret == "" {
		return nil, nil, fmt.Errorf("form parameter 'session_secret' MUST be specified")
	}

	serviceSession, err := s.schemaMgr.GetServiceSession(schema.MakeSessionID(sessionSecret))
	if err != nil && err != schema.ErrNoSuchSession {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get service session: %s", err))
	}
	if err == schema.ErrNoSuchSession || serviceSession.ServiceID != service.ID {
		return nil, nil, fmt.Errorf("no such service session")
	}

	// IMPORTANT! The service session object retrieved from storage does
	// not have the secret field set, so set it here while we have it.
	serviceSession.Secret = sessionSecret

	session, err := s.schemaMgr.GetSession(serviceSession.SessionID)
	if err != nil {
		if err == schema.ErrNoSuchSession {
			return nil, nil, err
		}

		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get session: %s", err))
	}

	now := time.Now()

	if now.After(session.Expiration) {
		return nil, nil, fmt.Errorf("session expired at %d - current time is %d", session.Expiration.Unix(), now.Unix())
	}

	user, err := s.schemaMgr.GetUserByID(session.UserID)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			return nil, nil, err
		}

		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to get account: %s", err))
	}

	if !user.IsActive {
		apiErr := authn.ErrAccountInactive()
		return nil, nil, fmt.Errorf(apiErr.Message)
	}

	sessionAuthenticator := s.authorizer.SessionTokenAuthenticator()
	if err := sessionAuthenticator.ExtendSession(session); err != nil {
		return nil, nil, errors.Internal(ctx, fmt.Errorf("unable to extend session: %s", err))
	}

	return user, serviceSession, nil
}

func (s *Service) handleRefreshTokenCodeGrant(ctx context.Context, r *restful.Request, service *schema.Service) (account *schema.Account, additionalAudience string, err error) {
	// The refresh token is simply the account ID.
	accountNameOrID := r.Request.PostFormValue("refresh_token")
	if accountNameOrID == "" {
		return nil, "", fmt.Errorf("form parameter 'refresh_token' MUST be specified")
	}

	// The service needs to be a priviledged service (for now, until we
	// have records for 3rd party service authorizations).
	if !service.Privileged {
		return nil, "", fmt.Errorf("must be a privileged 1st party service to use refresh tokens")
	}

	// Privileged services are always authorized, so just get the account.
	account, err = helpers.GetAccountByNameOrID(s.schemaMgr, accountNameOrID)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			return nil, "", fmt.Errorf("no such account for refresh token")
		}

		return nil, "", errors.Internal(ctx, fmt.Errorf("unable to get account: %s", err))
	}

	if !account.IsActive {
		apiErr := authn.ErrAccountInactive()
		return nil, "", fmt.Errorf(apiErr.Message)
	}

	additionalAudience = r.Request.PostForm.Get("audience")

	return account, additionalAudience, nil
}
