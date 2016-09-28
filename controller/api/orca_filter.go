package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/orca/controller/ctx"
)

func (a *Api) orcaInstanceFilter(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	req := rc.Request

	err := detectOrcaInstance(a.manager.DockerClient(), a.manager.ID(), w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	req.URL, err = url.ParseRequestURI(a.swarmClassicURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.fwd.ServeHTTP(w, req)
}

// Detect if the container in question is a UCP infrastructure container and require the force flag
func detectOrcaInstance(client *client.Client, ucpID string, w http.ResponseWriter, req *http.Request) error {
	req.ParseForm()
	force := strings.ToLower(req.FormValue("force"))

	// If force is set, no need to lookup if this is one of our infrastructure containers
	if force == "true" || force == "1" {
		return nil
	}

	vars := mux.Vars(req)
	name := vars["name"]

	// TODO Might make sense to cache the known container IDs so this isn't slowing down every !force request,
	//      but that should be part of a broader caching strategy...
	info, err := client.ContainerInspect(context.TODO(), name)
	if err == nil && info.Config != nil {
		if info.Config.Labels["com.docker.ucp.InstanceID"] == ucpID {
			return fmt.Errorf("WARNING: The container '%s' is part of the Universal Control Plane infrastructure.  You must use the force flag (-f) to proceed", name)
		}
	} else {
		// Could be simple user error, so no reason to be too alarmed, but we'll let the back-end reject it...
		log.Debug("Failed to lookup container %s - %s", name, err)
	}
	return nil
}
