package openid

import (
	"crypto"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"strings"

	// Register crypto.SHA256 on init()
	_ "crypto/sha256"

	"github.com/docker/orca/enzi/jose"
)

// RedirectState is used to maintain state between the initial authorization
// request and the callback.
type RedirectState struct {
	RedirectNext string `json:"redirectNext"`
	RedirectURI  string `json:"redirectURI"`
}

// Encode encodes this state and appends an HMAC-SHA256 signature using the
// given key.
func (state *RedirectState) Encode(key string) string {
	// This SHOULD NOT error.
	jsonBytes, _ := json.Marshal(state)

	payload := jose.Base64URLEncode(jsonBytes)

	hasher := hmac.New(crypto.SHA256.New, []byte(key))
	hasher.Write([]byte(payload))

	signature := jose.Base64URLEncode(hasher.Sum(nil))

	return fmt.Sprintf("%s.%s", payload, signature)
}

// DecodeRedirectState decodes the given encoded state by first verifying its
// HMAC-SHA256 signature using the given key and then decodes the payload into
// the returned RedirectState value.
func DecodeRedirectState(encoded, key string) (*RedirectState, error) {
	parts := strings.SplitN(encoded, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid encoding")
	}

	payload := parts[0]

	sigBytes, err := jose.Base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %s", err)
	}

	hasher := hmac.New(crypto.SHA256.New, []byte(key))
	hasher.Write([]byte(payload))

	if !hmac.Equal(sigBytes, hasher.Sum(nil)) {
		return nil, fmt.Errorf("invalid HMAC signature")
	}

	jsonBytes, err := jose.Base64URLDecode(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding: %s", err)
	}

	var state RedirectState
	if err := json.Unmarshal(jsonBytes, &state); err != nil {
		return nil, fmt.Errorf("unable to decode JSON: %s", err)
	}

	return &state, nil
}
