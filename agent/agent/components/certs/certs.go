package certs

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"

	"github.com/docker/orca/agent/agent/utils"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/config"
	caconfig "github.com/docker/orca/ca/config"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

// The Certs package attempts to verify whether the root key material is present and to sign all certificates
type Certs struct {
}

func (p *Certs) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	return nil
}

func (p *Certs) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	// Always reconcile certs on promotion
	if !currentCfg.IsManager && expectedCfg.IsManager {
		return true, nil
	}
	return false, nil
}

func (p *Certs) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// TODO: delete root key material on demotion

	// Seed rand so that we can contact remote managers in a random order
	rand.Seed(time.Now().UTC().UnixNano())

	// Create an HTTP Client to communicate with an existing manager through the UCP API.
	// The swarm-mode node cert is loaded as a client cert
	httpClient, err := utils.GetHTTPClient()
	if err != nil {
		return err
	}

	replicateServerCerts := true
	if utils.AreServerCertsPresent() {
		log.Info("Server certs detected in ucp-controller-server-certs, will not request replication")
		replicateServerCerts = false
	}

	// Ask a manager for the root key material until it is available in the target volume
	// TODO: verify the root material corresponds to the certs signed through CSR
	targetMgr := expectedCfg.Managers[rand.Intn(len(expectedCfg.Managers))]
	for !utils.IsRootMaterialPresent() {
		time.Sleep(2 * time.Second) // TODO: not needed

		// The first hostname is the actual hostname
		err = sendNodePromoteRequest(httpClient, targetMgr, config.OrcaHostAddress, replicateServerCerts)
		if err != nil {
			log.Error(err)
			targetMgr = expectedCfg.Managers[rand.Intn(len(expectedCfg.Managers))]
		}
	}

	err = extractSANsFromLabels(dclient)
	if err != nil {
		return err
	}

	info, err := dclient.Info(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to retrieve engine info: %s", err)
	}

	// Root material obtained, self-sign the manager certs
	err = setupManagerOnlyCerts(info.Swarm.NodeID)
	if err != nil {
		return err
	}

	log.Info("Certs component reconciled successfully")
	return nil
}

// TODO: move the result out of the config package
func extractSANsFromLabels(dclient *client.Client) error {
	info, err := dclient.Info(context.TODO())
	if err != nil {
		return err
	}

	node, _, err := dclient.NodeInspectWithRaw(context.TODO(), info.Swarm.NodeID)
	if err != nil {
		return err
	}

	if node.Spec.Annotations.Labels == nil {
		return nil
	}

	sanString, ok := node.Spec.Annotations.Labels["com.docker.ucp.SANs"]
	if !ok {
		return nil
	}

	config.OrcaHostnames = append(config.OrcaHostnames, strings.Split(sanString, ",")...)
	return nil
}

// setupManagerOnlyCerts self-signs all of the manager node's certs
// assumes the cluster and client root key material is in place
func setupManagerOnlyCerts(nodeID string) error {
	for _, certConfig := range []struct {
		ou    string
		mount string
		uid   int
		gid   int
	}{
		// KV store.
		{"kv", config.SwarmKvCertVolumeMount, 65534, 65534},
		// For UCP controller to talk to CA.
		{"ucp", config.SwarmControllerCertVolumeMount, 65534, 65534},
		// Auth datastore.
		{"auth-store", config.AuthStoreCertsVolumeMount, 65534, 65534},
		// Auth API server.
		{"auth-api", config.AuthAPICertsVolumeMount, 65534, 65534},
		// Auth worker server.
		{"auth-worker", config.AuthWorkerCertsVolumeMount, 65534, 65534},
	} {
		if err := certs.InitLocalNode(
			config.SwarmRootCAVolumeMount,
			certConfig.mount,
			nodeID,
			certConfig.ou,
			config.OrcaLocalName,
			config.OrcaHostnames,
			certConfig.uid,
			certConfig.gid,
		); err != nil {
			log.Errorf("Failed to setup %s cert", certConfig.ou)
			return err
		}
	}

	if !utils.AreServerCertsPresent() {
		if err := certs.InitLocalNode(config.OrcaRootCAVolumeMount, config.OrcaServerCertVolumeMount, "ucp", "", config.OrcaLocalName, config.OrcaHostnames, 65534, 65534); err != nil {
			return err
		}
	}

	// Make sure orca trusts swarm too
	if err := certs.SetupOrcaTrust(); err != nil {
		return err
	}

	return nil
}

// TODO: wire this up
func addExternalCA(dclient *client.Client) error {
	clusterCAAddr := net.JoinHostPort(config.OrcaHostAddress, strconv.Itoa(orcaconfig.SwarmCAPort))
	clusterCASignURL := fmt.Sprintf("https://%s%s", clusterCAAddr, caconfig.APISignPath)

	swarmCluster, err := dclient.SwarmInspect(context.TODO())
	if err != nil {
		return err
	}

	externalCAs := swarmCluster.Spec.CAConfig.ExternalCAs
	externalCAs = append(externalCAs, &swarm.ExternalCA{
		Protocol: swarm.ExternalCAProtocolCFSSL,
		URL:      clusterCASignURL,
	})
	swarmCluster.Spec.CAConfig.ExternalCAs = externalCAs

	err = dclient.SwarmUpdate(context.TODO(), swarmCluster.Meta.Version, swarmCluster.Spec, swarm.UpdateFlags{})
	if err != nil {
		return err
	}
	return nil
}

func sendNodePromoteRequest(httpClient *http.Client, manager, hostname string, replicate bool) error {
	orcaURL := &url.URL{}
	orcaURL.Host = manager
	orcaURL.Path = "/api/nodes/promote"
	orcaURL.Scheme = "https"

	if replicate {
		q := orcaURL.Query()
		q.Set("replicate", "true")
		orcaURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("POST", orcaURL.String(), nil)
	if err != nil {
		log.Debug("Failed to build request")
		return err
	}
	req.Header.Set("Node", hostname)
	log.Infof("Attemping Node Promotion as node %s towards controller at %s", hostname, manager)
	_, err = httpClient.Do(req)
	if err != nil {
		log.Debug("Failed to send request")
		return err
	}
	return nil
}
