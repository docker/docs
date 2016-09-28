package garant

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/garant/auth"
)

// resolveScopeSpecifiers converts a list of scope specifiers from a token
// request's `scope` query parameters into a list of standard access objects.
func resolveScopeSpecifiers(scopeSpecs []string) []auth.Access {
	requestedAccessSet := make(map[auth.Access]struct{}, 2*len(scopeSpecs))

	for _, scopeSpecifier := range scopeSpecs {
		// There should be 3 parts, separated by a `:` character.  The resource
		// name might have extra colons in it though.  So take the first common
		// delimited token, the last comma-delimited token, and assume every
		// other token should be bundled together.
		parts := strings.Split(scopeSpecifier, ":")

		if len(parts) < 3 {
			// Ignore malformed scope specifiers.
			continue
		}

		resourceType := parts[0]
		resourceName := strings.Join(parts[1:len(parts)-1], ":")
		actions := parts[len(parts)-1]

		// Actions should be a comma-separated list of actions.
		for _, action := range strings.Split(actions, ",") {
			requestedAccess := auth.Access{
				Resource: auth.Resource{
					Type: resourceType,
					Name: resourceName,
				},
				Action: action,
			}

			// Add this access to the requested access set.
			requestedAccessSet[requestedAccess] = struct{}{}
		}
	}

	requestedAccessList := make([]auth.Access, 0, len(requestedAccessSet))
	for requestedAccess := range requestedAccessSet {
		requestedAccessList = append(requestedAccessList, requestedAccess)
	}

	return requestedAccessList
}

// accessEntry represents an access entry in a JWT.
type accessEntry struct {
	auth.Resource
	Actions []string `json:"actions"`
}

// CreateJWT creates and signs a JSON Web Token for the given account and
// audience with the granted access.
func (app *App) CreateJWT(account auth.Account, audience string, grantedAccessList []auth.Access, expiresIn time.Duration) (token string) {
	// Make a set of access entries to put in the token's claimset.
	resourceActionSets := make(map[auth.Resource]map[string]struct{}, len(grantedAccessList))
	for _, access := range grantedAccessList {
		actionSet, exists := resourceActionSets[access.Resource]
		if !exists {
			actionSet = map[string]struct{}{}
			resourceActionSets[access.Resource] = actionSet
		}
		actionSet[access.Action] = struct{}{}
	}

	accessEntries := make([]accessEntry, 0, len(resourceActionSets))
	for resource, actionSet := range resourceActionSets {
		actions := make([]string, 0, len(actionSet))
		for action := range actionSet {
			actions = append(actions, action)
		}

		accessEntries = append(accessEntries, accessEntry{
			Resource: resource,
			Actions:  actions,
		})
	}

	// Now build up the JWT!

	randomBytes := make([]byte, 15)
	io.ReadFull(rand.Reader, randomBytes) // Assume it does not fail.
	randomID := base64.URLEncoding.EncodeToString(randomBytes)

	now := time.Now()

	joseHeader := map[string]interface{}{
		"typ": "JWT",
		"alg": "ES256",
	}

	if x5c := app.signingKey.GetExtendedField("x5c"); x5c != nil {
		joseHeader["x5c"] = x5c
	} else {
		joseHeader["jwk"] = app.signingKey.PublicKey()
	}

	var subject string
	if account != nil {
		subject = account.Subject()
	}

	claimSet := map[string]interface{}{
		"iss": app.config.Issuer,
		"sub": subject,
		"aud": audience,
		"exp": now.Add(expiresIn).Unix(),
		"nbf": now.Add(-expiresIn).Unix(),
		"iat": now.Unix(),
		"jti": randomID,

		"access": accessEntries,
	}

	var (
		joseHeaderBytes, claimSetBytes []byte
		err                            error
	)

	if joseHeaderBytes, err = json.Marshal(joseHeader); err != nil {
		panic(fmt.Errorf("unable to encode jose header: %s", err))
	}
	if claimSetBytes, err = json.Marshal(claimSet); err != nil {
		panic(fmt.Errorf("unable to encode claim set: %s", err))
	}

	encodedJoseHeader := joseBase64Encode(joseHeaderBytes)
	encodedClaimSet := joseBase64Encode(claimSetBytes)
	encodingToSign := fmt.Sprintf("%s.%s", encodedJoseHeader, encodedClaimSet)

	var signatureBytes []byte
	if signatureBytes, _, err = app.signingKey.Sign(strings.NewReader(encodingToSign), crypto.SHA256); err != nil {
		panic(fmt.Errorf("unable to sign jwt payload: %s", err))
	}

	signature := joseBase64Encode(signatureBytes)

	return fmt.Sprintf("%s.%s", encodingToSign, signature)
}

func joseBase64Encode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}
