package jose

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Audience is a special string slice which can be JSON decoded as either a
// single string or an array of strings.
type Audience []string

// Contains returns whether or not this audience set contains the given
// expected audience value.
func (aud Audience) Contains(expectedAudience string) bool {
	for _, audience := range aud {
		if audience == expectedAudience {
			return true
		}
	}

	return false
}

// MarshalJSON encodes this audience value as either a single string or an
// array of strings depending on the number of audiences.
func (aud Audience) MarshalJSON() ([]byte, error) {
	// If there is only a single audience, encode as a JSON string.
	if len(aud) == 1 {
		return json.Marshal(aud[0])
	}

	// Otherwise encode as a JSON array.
	var audiences []string = aud
	return json.Marshal(audiences)
}

// UnmarshalJSON decodes the given JSON data as either a string or a slice
// of strings.
func (aud *Audience) UnmarshalJSON(data []byte) error {
	// Try the common special case first: a single string.
	var audStr string
	if err := json.Unmarshal(data, &audStr); err == nil {
		*aud = []string{audStr}
		return nil
	}

	// Fallback to the general case of a slice of strings.
	return json.Unmarshal(data, (*[]string)(aud))
}

// JWTHeader is the decoded header of a JSON Web Token.
type JWTHeader struct {
	Type       string `json:"typ"`
	SigningAlg string `json:"alg"`
	KeyID      string `json:"kid"`
}

// JWTClaims is the decoded claim set of an OpenID Connect JSON Web Token.
type JWTClaims struct {
	Issuer          string   `json:"iss"`
	Subject         string   `json:"sub"`
	Audience        Audience `json:"aud"`
	AuthorizedParty string   `json:"azp"`
	IssuedAt        int64    `json:"iat"`
	Expiration      int64    `json:"exp"`

	// Extra Standard Claims.
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
}

// JWT is a deconstructed JSON Web Token used with the OpenID Connect Protocol.
type JWT struct {
	Header JWTHeader
	Claims JWTClaims

	SigningInput string
	Signature    string
}

// DecodeJWT parses and decodes this JWT and performs basic pre-validation. The
// signature IS NOT verified. Token expiration is checked with 1 minute of
// leeway.
func DecodeJWT(rawToken string) (*JWT, error) {
	parts := strings.SplitN(rawToken, ".", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("must have 3 parts separated by '.'")
	}

	rawHeader, err := Base64URLDecode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to base64url-decode header: %s", err)
	}

	var header JWTHeader
	if err := json.Unmarshal(rawHeader, &header); err != nil {
		return nil, fmt.Errorf("unable to JSON-decode header: %s", err)
	}

	if header.Type != "JWT" {
		return nil, fmt.Errorf("JWS header type value must be %q not %q", "JWT", header.Type)
	}

	rawClaims, err := Base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("unable to base64url-decode claims: %s", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return nil, fmt.Errorf("unable to JSON-decode claims: %s", err)
	}

	return &JWT{
		Header:       header,
		Claims:       claims,
		SigningInput: strings.Join(parts[:2], "."),
		Signature:    parts[2],
	}, nil
}

// NewJWT uses the given Private Key and claims to create and sign a new JWT.
func NewJWT(key *PrivateKey, claims JWTClaims) (*JWT, error) {
	signer, err := key.Signer()
	if err != nil {
		return nil, fmt.Errorf("unable to get signer for key: %s", err)
	}

	jwt := &JWT{
		Header: JWTHeader{
			Type:       "JWT",
			SigningAlg: signer.Name(),
			KeyID:      key.ID,
		},
		Claims: claims,
	}

	rawHeader, err := json.Marshal(jwt.Header)
	if err != nil {
		return nil, fmt.Errorf("unable to encode JWT header: %s", err)
	}
	rawClaimset, err := json.Marshal(jwt.Claims)
	if err != nil {
		return nil, fmt.Errorf("unable to encode JWT claims: %s", err)
	}

	jwt.SigningInput = fmt.Sprintf("%s.%s", Base64URLEncode(rawHeader), Base64URLEncode(rawClaimset))
	jwt.Signature, err = signer.Sign(strings.NewReader(jwt.SigningInput))
	if err != nil {
		return nil, fmt.Errorf("unable to sign JWT: %s", err)
	}

	return jwt, nil
}

func (jwt *JWT) String() string {
	return fmt.Sprintf("%s.%s", jwt.SigningInput, jwt.Signature)
}

// ExpiresIn returns the lifetime of this JWT in seconds.
func (jwt *JWT) ExpiresIn() int64 {
	return jwt.Claims.Expiration - jwt.Claims.IssuedAt
}

// Verify uses the given public key to verify the signature of this JWT.
func (jwt *JWT) Verify(key *PublicKey, expectedAudience, trustedIssuer string) error {
	if !jwt.Claims.Audience.Contains(expectedAudience) {
		return fmt.Errorf("JWT intended for a different audience")
	}

	if jwt.Claims.Issuer != trustedIssuer {
		return fmt.Errorf("JWT issuer untrusted")
	}

	// Validate that the token is not yet expired.
	expiration := time.Unix(jwt.Claims.Expiration, 0)
	now := time.Now().Round(time.Second)

	// Check if expiration is more than one minute in the past.
	if now.After(expiration.Add(time.Minute)) {
		return fmt.Errorf("JWT expired %s ago", now.Sub(expiration))
	}

	verifier, err := key.Verifier(jwt.Header.SigningAlg)
	if err != nil {
		return fmt.Errorf("unable to get verifier for key: %v", err)
	}

	if err := verifier.Verify(strings.NewReader(jwt.SigningInput), jwt.Signature); err != nil {
		return fmt.Errorf("invalid JWT signature: %v", err)
	}

	return nil
}
