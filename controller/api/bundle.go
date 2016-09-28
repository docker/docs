package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) generateClientBundle(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	r := rc.Request
	host := r.Header.Get("X-Docker-Host")
	if host == "" {
		host = r.Host
	}

	// check for port; assume 443 if not given
	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:443", host)
	}

	// Set up an initial label which the user can change later
	initialLabel := "Generated on " + time.Now().Format(time.RFC1123)

	data, err := a.manager.GenerateClientBundle(rc.Auth, host, initialLabel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Encoding", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename='ucp-bundle-%s.zip'", rc.Auth.User.Username))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write(data)
}
