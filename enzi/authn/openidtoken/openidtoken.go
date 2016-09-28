package openidtoken

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/enzi/schema"
)

func invalidJWT(reason string) *errors.APIError {
	return errors.InvalidAuthentication(fmt.Sprintf("invalid JWT: %s", reason))
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
	jwt, err := jose.DecodeJWT(token)
	if err != nil {
		return nil, invalidJWT(fmt.Sprintf("unable to decode JWT: %v", err))
	}

	// Validate token fields and verify signature.
	if validationErr := a.verifyJWT(ctx, jwt); validationErr != nil {
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

	if !account.IsActive {
		return nil, authn.ErrAccountInactive()
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

// verifyJWT verifies the signature of the given JWT.
func (a *authenticator) verifyJWT(ctx context.Context, jwt *jose.JWT) *errors.APIError {
	openIDConfig, err := config.GetOpenIDConfig(a.schemaMgr)
	if err != nil {
		return errors.Internal(ctx, fmt.Errorf("unable to get issuer identifier: %s", err))
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
	if err := jwt.Verify(pubKey, openIDConfig.IssuerIdentifier, openIDConfig.IssuerIdentifier); err != nil {
		return invalidJWT(fmt.Sprintf("invalid JWT signature: %s", err))
	}

	return nil
}

func checkStringSliceContains(target string, candidates ...string) bool {
	for _, candidate := range candidates {
		if target == candidate {
			return true
		}
	}

	return false
}
