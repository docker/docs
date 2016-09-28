package api

import (
	"net/http"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) authSyncMessages(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(a.manager.AuthSyncMessages(rc.Auth)))
}

func (a *Api) doAuthSync(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(a.manager.AuthSync(rc.Auth)))
}
