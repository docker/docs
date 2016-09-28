package client

import (
	"net/http"

	"github.com/docker/orca/enzi/api/responses"
)

// GetSigningKeys returns a JWK set of the signing keys which are currently
// used by the OpenID Connect provider. If there is an API error response then
// the returned error will be of the type
// *(github.com/docker/orca/enzi/api/errors).APIErrors
func (s *Session) GetSigningKeys() (*responses.JWKSet, error) {
	endpoint := s.buildURL("/v0/signingKeys", nil)

	var keys responses.JWKSet
	if err := s.performRequest("GET", endpoint, nil, http.StatusOK, &keys, nil); err != nil {
		return nil, err
	}

	return &keys, nil
}
