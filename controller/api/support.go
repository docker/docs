package api

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) supportDump(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// TODO(alexmavr): use a Resource to control access to this
	if !rc.Auth.User.Admin {
		http.Error(w, "cannot create support dump for non-admin user", http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/zip")
	n := time.Now()
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("inline; filename='docker-support-%04d%02d%02d-%02d:%02d:%02d.zip'",
			n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second()))
	if err := a.manager.SupportDump(w); err != nil {
		log.Errorf("Could not create support dump: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
