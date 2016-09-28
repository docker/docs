package openidtoken

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/enzi/schema"
)

func invalidJWT(reason string) *errors.APIError {
	return errors.InvalidAuthentication(fmt.Sprintf("invalid JWT: %s", reason))
}

// IdentityJWT is a deconstructed JSON Web Token used with the OpenID Connect
// Protocol.
type IdentityJWT struct {
	Header
	Claims

	SigningInput string
	Signature    string
}

func (jwt *IdentityJWT) String() string {
	return fmt.Sprintf("%s.%s", jwt.SigningInput, jwt.Signature)
}

// ExpiresIn returns the lifetime of this JWT in seconds.
func (jwt *IdentityJWT) ExpiresIn() int64 {
	return jwt.Claims.Expiration - jwt.Claims.IssuedAt
}

// Header is the decoded header of a JSON Web Token.
type Header struct {
	Type       string `json:"typ"`
	SigningAlg string `json:"alg"`
	KeyID      string `json:"kid"`
}

// Claims is the decoded body of an OpenID Connect JSON Web Token.
type Claims struct {
	Issuer          string   `json:"iss"`
	Subject         string   `json:"sub"`
	Audience        []string `json:"aud"`
	AuthorizedParty string   `json:"azp"`
	IssuedAt        int64    `json:"iat"`
	Expiration      int64    `json:"exp"`
}

type authenticator struct {
	schemaMgr schema.Manager
}

var _ authn.OpenIDTokenAuthenticator = (*authenticator)(nil)

// New creates an authenticator with the given schema manager.
func New(schemaMgr schema.Manager) authn.OpenIDTokenAuthenticator {
	return &authenticator{
		schemaMgr: schemaMgr,
	}
}

// AuthenticateOpenIDToken follows the same semantics as
// AuthenticateRequest in the authn.Authenticator interface but only attempts
// openid token authentication which authenticates the client as a service
// acting on behalf of an account.
func (a *authenticator) AuthenticateOpenIDToken(ctx context.Context, token string) (*authn.Account, *errors.APIError) {
	// Token value should be an encoded JWT.
	jwt, parseErr := parseJWT(token)
	if parseErr != nil {
		return nil, parseErr
	}

	// Validate token fields and verify signature.
	if validationErr := a.validateJWT(ctx, jwt); validationErr != nil {
		return nil, validationErr
	}

	// Get the account which is the subject of the JWT.
	account, err := a.schemaMgr.GetAccountByID(jwt.Claims.Subject)
	if err != nil {
		if err == schema.ErrNoSuchAccount {
			// The user no longer exists?
			return nil, invalidJWT(fmt.Sprintf("no such account: %s", jwt.Claims.Subject))
		}

		// Internal error.
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get user for session: %s", err))
	}

	// Get the service which is the authorized party of the JWT.
	service, err := a.schemaMgr.GetServiceByID(jwt.Claims.AuthorizedParty)
	if err != nil {
		if err == schema.ErrNoSuchService {
			// The service no longer exists?
			return nil, invalidJWT(fmt.Sprintf("no such service: %s", jwt.Claims.AuthorizedParty))
		}

		// Internal error.
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get user for session: %s", err))
	}

	return &authn.Account{
		Account:           *account,
		AuthorizedService: service,
	}, nil
}

func parseJWT(rawToken string) (*IdentityJWT, *errors.APIError) {
	parts := strings.SplitN(rawToken, ".", 3)
	if len(parts) != 3 {
		return nil, invalidJWT("must have 3 parts separated by '.'")
	}

	rawHeader, err := jose.Base64URLDecode(parts[0])
	if err != nil {
		return nil, invalidJWT(fmt.Sprintf("unable to base64url-decode header: %s", err))
	}

	var header Header
	if err := json.Unmarshal(rawHeader, &header); err != nil {
		return nil, invalidJWT(fmt.Sprintf("unable to JSON-decode header: %s", err))
	}

	rawClaims, err := jose.Base64URLDecode(parts[1])
	if err != nil {
		return nil, invalidJWT(fmt.Sprintf("unable to base64url-decode claims: %s", err))
	}

	var claims Claims
	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return nil, invalidJWT(fmt.Sprintf("unable to JSON-decode claims: %s", err))
	}

	return &IdentityJWT{
		Header:       header,
		Claims:       claims,
		SigningInput: strings.Join(parts[:2], "."),
		Signature:    parts[2],
	}, nil
}

// validateJWT validates token fields and verifies its signature.
func (a *authenticator) validateJWT(ctx context.Context, jwt *IdentityJWT) *errors.APIError {
	// Check that the required header values have been specified.
	if jwt.Header.Type != "JWT" {
		return invalidJWT("header value 'typ' must be 'JWT'")
	}

	if jwt.Header.SigningAlg == "" {
		return invalidJWT("header value 'alg' must be specified")
	}

	if jwt.Header.KeyID == "" {
		return invalidJWT("header value 'kid' must be specified")
	}

	// Validate that the token is not yet expired.
	expiration := time.Unix(jwt.Claims.Expiration, 0)
	now := time.Now()

	// Check if expiration is in the past.
	if expiration.Before(now) {
		return invalidJWT(fmt.Sprintf("expired at %d - current time is %d", expiration.Unix(), now.Unix()))
	}

	// Get the service which was issued the token.
	service, err := a.schemaMgr.GetServiceByID(jwt.AuthorizedParty)
	if err != nil {
		if err == schema.ErrNoSuchService {
			return invalidJWT(fmt.Sprintf("no such authorized service: %q", jwt.AuthorizedParty))
		}

		return errors.Internal(ctx, fmt.Errorf("unable to get authorized service: %s", err))
	}

	// Ensure that we are the issuer and in the audience.
	if !checkAcceptedIssuer(service.ProviderIdentities, jwt.Issuer) {
		return invalidJWT(fmt.Sprintf("unacceptable issuer in JWT: %s", jwt.Issuer))
	}

	if !checkAcceptedAudiences(service.ProviderIdentities, jwt.Audience) {
		return invalidJWT(fmt.Sprintf("unacceptable audience in JWT: %s", jwt.Audience))
	}

	// Get the key which was used to sign the token.
	jwk, err := a.schemaMgr.GetSigningKey(jwt.Header.KeyID)
	if err != nil {
		return errors.Internal(ctx, fmt.Errorf("unable to get JWT signing key: %s", err))
	}

	pubKey, err := jwk.PublicKey()
	if err != nil {
		return errors.Internal(ctx, fmt.Errorf("unable to decode JWK into public key: %s", err))
	}

	// Verify the signature.
	verifier, err := pubKey.Verifier(jwt.Header.SigningAlg)
	if err != nil {
		return invalidJWT(fmt.Sprintf("unable to verify token signature: %s", err))
	}

	if err := verifier.Verify(strings.NewReader(jwt.SigningInput), jwt.Signature); err != nil {
		return invalidJWT(fmt.Sprintf("invalid JWT signature: %s", err))
	}

	return nil
}

func checkAcceptedIssuer(acceptedIssuers []string, tokenIssuer string) bool {
	for _, acceptedIssuer := range acceptedIssuers {
		if tokenIssuer == acceptedIssuer {
			return true
		}
	}

	return false
}

func checkAcceptedAudiences(acceptedAudiences, jwtAudiences []string) bool {
	accpetedAudienceSet := make(map[string]struct{}, len(acceptedAudiences))
	for _, acceptedAudience := range acceptedAudiences {
		accpetedAudienceSet[acceptedAudience] = struct{}{}
	}

	for _, audience := range jwtAudiences {
		if _, ok := accpetedAudienceSet[audience]; ok {
			return true
		}
	}

	return false
}
