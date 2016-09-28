package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/ctx"
	orcatypes "github.com/docker/orca/types"
)

func (a *Api) nodes(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	nodes, err := a.manager.Nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) node(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	node, err := a.manager.Node(rc.PathVars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) authorizeNodeRequest(w http.ResponseWriter, r *http.Request) {
	// Ensure a client cert is used and it corresponds a swarm-mode node
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 || r.TLS.PeerCertificates[0] == nil {
		http.Error(w, "Client authorization required", http.StatusInternalServerError)
		return
	}
	ous := r.TLS.PeerCertificates[0].Subject.OrganizationalUnit
	if len(ous) != 1 || (ous[0] != "swarm-manager" && ous[0] != "swarm-worker") {
		log.Warnf("CSR received on /api/nodes/authorize with unrecognized OUs: %s", ous)
		http.Error(w, "Invalid client certificate", http.StatusInternalServerError)
		return
	}

	log.Info("Authorizing Node CSR Request with Client Certificate OU: \"%s\"", ous[0])
	var req orca.NodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Check for v1 client, and give a specific error
		var oldReq orca.V1NodeRequest
		if err := json.NewDecoder(r.Body).Decode(&oldReq); err == nil {
			http.Error(w, "Your UCP tool is too old.  Please pull an updated copy of docker/ucp before joining", http.StatusBadRequest)
			return
		}
		// Malformed, so let the original error percolate out
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	log.Infof("Join node request from %s - %s", r.RemoteAddr, r.UserAgent())

	nodeConfig, err := a.manager.AuthorizeNodeRequest(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(nodeConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// setControllerServerCerts changes the server certificates of the controller located at the
// hostname specified under the "name" query parameter. If "name" is not specified, then the
// certificates will be copied to all controllers
func (a *Api) setControllerServerCerts(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var serverCertRequest orcatypes.ServerCertRequest
	err := json.NewDecoder(rc.Body()).Decode(&serverCertRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = a.manager.SetControllerServerCerts(rc.QueryVars.Get("name"), serverCertRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Api) promoteNode(w http.ResponseWriter, r *http.Request) {
	// Extract the classic swarm node name for the remote host
	node := ""
	nodeSlice, ok := r.Header["Node"]
	if ok {
		node = nodeSlice[0]
	}

	// Determine whether the node requesting promotion also requested for the
	// server certs to be replicated.
	// This will be false if externally signed certs were already placed in the
	// remote node's ucp-controller-server-certs volume
	replicateCerts := false
	if r.URL.Query().Get("replicate") != "" {
		replicateCerts = true
	}
	log.Infof("Promoting node %s, with cert replication as %b", node, replicateCerts)

	err := a.manager.PromoteNode(node, replicateCerts)
	if err != nil {
		log.Warnf("Could not promote node: %s", err)
	}

	// Always respond with 400
	http.Error(w, "", http.StatusBadRequest)
}

func (a *Api) listNodes(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	nodelist, err := a.manager.ListNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(nodelist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Api) inspectNode(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	node, err := a.manager.InspectNode(rc.PathVars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
