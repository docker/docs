package openid

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/api/server"
	"github.com/emicklei/go-restful"
)

func (s *Service) routeListSigningKeys() server.Route {
	return server.Route{
		Method:  "GET",
		Path:    "/signingKeys",
		Handler: s.handleListSigningKeys,
		Doc:     "Get a cacheable JSON Web Key Set of all signing keys currently in use",
		ResponseDocs: []server.ResponseDoc{
			{
				Code:    http.StatusOK,
				Message: "Success, JWK set returned.",
				Sample:  responses.JWKSet{},
			},
		},
	}
}

func (s *Service) handleListSigningKeys(ctx context.Context, r *restful.Request) responses.APIResponse {
	keys, err := s.schemaMgr.ListSigningKeys()
	if err != nil {
		return responses.APIError(errors.Internal(ctx, fmt.Errorf("unable to list signing keys: %s", err)))
	}

	cacheControlHeader := http.Header{
		"Cache-Control": []string{fmt.Sprintf("max-age=%d", uint(s.signingKeyCacheMaxAge.Seconds()))},
	}

	return responses.JSONResponseWithHeaders(http.StatusOK, responses.MakeJWKSet(keys), cacheControlHeader)
}
