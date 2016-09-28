package idauth

import (
	"fmt"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/jose"
)

// Authenticator is able to authenticate a client which provides an identity
// token issued by an eNZi OpenID Connect Provider.
type Authenticator struct {
	Context          context.Context
	APIClient        *client.Session
	IssuerIdentifier string
	ServiceID        string

	SigningKeyCache KeyCache
}

// AuthenticateIdentityToken attempts to parse, and verify the given raw JSON
// Web Token and returns the Claims from within.
func (a *Authenticator) AuthenticateIdentityToken(rawToken string) (*jose.JWTClaims, error) {
	jwt, err := jose.DecodeJWT(rawToken)
	if err != nil {
		return nil, fmt.Errorf("unable to decode identity token: %v", err)
	}

	key, err := a.getVerifyingKey(jwt.Header.KeyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get verifying key: %v", err)
	}

	if err := jwt.Verify(key, a.ServiceID, a.IssuerIdentifier); err != nil {
		return nil, fmt.Errorf("unable to verify identity token: %v", err)
	}

	return &jwt.Claims, nil
}

func (a *Authenticator) getVerifyingKey(keyID string) (*jose.PublicKey, error) {
	// Check the cache.
	key := a.SigningKeyCache.Get(keyID)
	if key != nil {
		return key, nil
	}

	// Need to get keys from the provider's JWKs endpoint.
	jwkSet, err := a.APIClient.GetSigningKeys()
	if err != nil {
		return nil, fmt.Errorf("unable to get token signing keys from provider: %v", err)
	}

	var targetKey *jose.PublicKey
	for _, jwk := range jwkSet.Keys {
		key, err := jose.NewPublicKeyJWK(jose.JWK(jwk))
		if err != nil {
			// Log the issue and skip.
			context.GetLogger(a.Context).Warnf("unable to decode token signing key from provider: %v", err)
			continue
		}

		// Cache it for 6 hours.
		a.SigningKeyCache.Set(key, 6*time.Hour)

		if key.ID == keyID {
			targetKey = key
		}
	}

	if targetKey == nil {
		return nil, fmt.Errorf("token signing key %q not found at provider", keyID)
	}

	return targetKey, nil
}
