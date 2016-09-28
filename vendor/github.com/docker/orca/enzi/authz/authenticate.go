package authz

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/authn"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/config"

	"github.com/docker/distribution/context"
)

// AuthenticateRequestUser attempts to authenticate the request and return a
// user object. This authorizer will never return (nil, nil) - the user will
// be an anonymous user object instead.
func (a *authorizer) AuthenticateRequest(ctx context.Context, r *http.Request) (*authn.Account, *errors.APIError) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return authn.MakeAnonymousAccount(), nil
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 {
		return nil, errors.InvalidAuthentication(fmt.Sprintf("invalid Authorization header: %q", authorizationHeader))
	}

	authScheme := strings.ToLower(parts[0])

	// Determine which auth backend we're using.
	authConfig, err := a.AuthConfig()
	if err != nil {
		return nil, errors.Internal(ctx, fmt.Errorf("unable to get current auth config: %s", err))
	}

	// If we're using the LDAP, load the LDAP settings.
	var ldapSettings *ldapconfig.Settings
	if authConfig.Backend == config.AuthBackendLDAP {
		ldapSettings, err = ldapconfig.GetLDAPConfig(a.schemaMgr)
		if err != nil {
			return nil, errors.Internal(ctx, fmt.Errorf("unable to get current LDAP config: %s", err))
		}
	}

	var (
		account *authn.Account
		apiErr  *errors.APIError
	)

	switch authScheme {
	case "basic":
		account, apiErr = a.handleBasicAuthentication(ctx, ldapSettings, parts[1])
	case "openidtoken":
		account, apiErr = a.OpenIDTokenAuthenticator().AuthenticateOpenIDToken(ctx, parts[1])
	case "sessiontoken", "bearer":
		// Handle both auth schemes as Session Token.
		account, apiErr = a.SessionTokenAuthenticator().AuthenticateSessionToken(ctx, parts[1])
	default:
		return nil, errors.InvalidAuthentication("unknown authentication scheme")
	}

	if apiErr != nil {
		return nil, apiErr
	}

	// If we're using the LDAP backend, check if this is the recovery admin
	// user. They should always be active and always be an admin.
	if ldapSettings != nil && ldapSettings.RecoveryAdminUsername == account.Name {
		account.IsActive = true
		account.IsAdmin = true
	}

	// Check if the account is active.
	if !account.IsActive {
		return nil, errors.InactiveAccount()
	}

	return account, nil
}

func (a *authorizer) handleBasicAuthentication(ctx context.Context, ldapSettings *ldapconfig.Settings, encoded string) (*authn.Account, *errors.APIError) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.InvalidAuthentication("unable to decode basic authentication credentials")
	}

	decodedParts := strings.SplitN(string(decoded), ":", 2)
	if len(decodedParts) < 2 {
		return nil, errors.InvalidAuthentication("unable to decode basic authentication credentials")
	}

	username, password := decodedParts[0], decodedParts[1]

	basicAuthenticator := a.UsernamePasswordAuthenticator(ldapSettings)

	return basicAuthenticator.AuthenticateUsernamePassword(ctx, username, password)
}
