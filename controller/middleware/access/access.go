package access

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/controller/ctx"
)

var ErrAccessDenied = errors.New("access denied")
var ErrBadRequest = errors.New("Malformed request")

func LayerHandler(rc *ctx.OrcaRequestContext) (int, error) {
	// Admins have access to everything
	if rc.Auth.User.Admin {
		return http.StatusOK, nil
	}

	for _, resource := range rc.Resources {
		if resource == nil {
			return http.StatusBadRequest, ErrBadRequest
		}
		if !resource.HasAccess(rc.Auth) {
			log.Warnf("access denied for %s to %s from %s", rc.Auth.User.Username, rc.Request.URL.Path, rc.Request.RemoteAddr)
			return http.StatusForbidden, ErrAccessDenied
		}
	}

	return http.StatusOK, nil
}
