package csr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/orca"
	"github.com/docker/orca/agent/agent/utils"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/types"
)

type CSR struct {
}

func (p *CSR) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	return nil
}

func (p *CSR) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	return false, nil
}

func (p *CSR) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	// Check if node certs are present in the current node
	certExists := true
	keyExists := true
	caCertExists := true
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.SwarmNodeCertVolumeMount, config.CertFilename)); os.IsNotExist(err) {
		certExists = false
	}
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.SwarmNodeCertVolumeMount, config.CAFilename)); os.IsNotExist(err) {
		caCertExists = false
	}
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.SwarmNodeCertVolumeMount, config.KeyFilename)); os.IsNotExist(err) {
		keyExists = false
	}

	// Don't reconcile CSR if the certs are all there
	// TODO: perform better verification of the cert chain
	if certExists && keyExists && caCertExists {
		return nil
	}

	// Check if root key material is there and self-sign everything if so
	if utils.IsRootMaterialPresent() {
		log.Info("TODO")
		return nil
	}

	// Root material is not present, attempt CSR with other managers
	info, err := dclient.Info(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to retrieve engine info: %s", err)
	}

	// Obtain an http.Client with the swarm-mode node cert as a client cert
	httpClient, err := utils.GetHTTPClient()
	if err != nil {
		return err
	}

	// Seed rand so that we can contact remote managers in a random order
	rand.Seed(time.Now().UTC().UnixNano())
	// Keep trying to perform a node join CSR
	for {
		// Select a random manager to contact
		targetMgr := expectedCfg.Managers[rand.Intn(len(expectedCfg.Managers))]

		err = doCSR(httpClient, targetMgr, info.Swarm.NodeID)
		if err == nil {
			break
		}
		log.Infof("Unable to perform CSR request against manager at %s: %s", targetMgr, err)
		time.Sleep(3 * time.Second) //TODO: fine-tune backoff
	}
	log.Info("CSR component reconciled successfully")
	return nil
}

func doCSR(httpClient *http.Client, manager, nodeID string) error {
	var nodeConfig *orca.NodeConfiguration
	var err error

	uid := 65534
	gid := 65534

	// generate the request
	req := orca.NodeRequest{
		ClusterCertificateRequests: map[string]string{},
		UserCertificateRequests:    map[string]string{},
	}
	csr, key, err := certs.GenerateCSR(nodeID, "swarm", config.OrcaHostnames)
	if err != nil {
		log.Debug("Failed to generate CSR")
		return err
	}
	// Write out the keys to the proper cert locations
	path := filepath.Join(config.SwarmNodeCertVolumeMount, config.KeyFilename)
	if err := ioutil.WriteFile(path, key, 0600); err != nil {
		return err
	}
	if err := os.Chown(path, uid, gid); err != nil {
		return err
	}
	// Use the mount point so we can write out to the location specified
	req.ClusterCertificateRequests[config.SwarmNodeCertVolumeMount] = string(csr)

	nodeConfig, err = sendNodeCSR(&req, httpClient, manager)
	if err != nil {
		return err
	}

	log.Infof("Joining UCP ID: %s", nodeConfig.OrcaID)
	if nodeConfig == nil || nodeConfig.OrcaID == "" {
		return fmt.Errorf("Unable to retrieve a valid node configuration from an existing manager")
	}

	if config.OrcaInstanceID != nodeConfig.OrcaID {
		return fmt.Errorf("UCP instance ID mismatch: %s against %s", config.OrcaInstanceID, nodeConfig.OrcaID)
	}

	// Write out the results
	for dir, cert := range nodeConfig.UserCertificates {
		if err := ioutil.WriteFile(filepath.Join(dir, config.CertFilename), []byte(cert), 0644); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(dir, config.CAFilename), []byte(nodeConfig.UserCertificateChain), 0644); err != nil {
			return err
		}
		if err := os.Chown(filepath.Join(dir, config.CertFilename), uid, gid); err != nil {
			return err
		}
		if err := os.Chown(filepath.Join(dir, config.CAFilename), uid, gid); err != nil {
			return err
		}
	}

	for dir, cert := range nodeConfig.ClusterCertificates {
		if err := ioutil.WriteFile(filepath.Join(dir, config.CertFilename), []byte(cert), 0644); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(dir, config.CAFilename), []byte(nodeConfig.ClusterCertificateChain), 0644); err != nil {
			return err
		}
		if err := os.Chown(filepath.Join(dir, config.CertFilename), uid, gid); err != nil {
			return err
		}
		if err := os.Chown(filepath.Join(dir, config.CAFilename), uid, gid); err != nil {
			return err
		}
	}

	return nil
}

func sendNodeCSR(nodeReq *orca.NodeRequest, httpClient *http.Client, manager string) (*orca.NodeConfiguration, error) {
	log.Infof("Sending CSR Request to manager at %s", manager)

	reqJson, err := json.Marshal(*nodeReq)
	if err != nil {
		log.Debug("Failed to generate json csr")
		return nil, err
	}

	orcaURL := &url.URL{}
	orcaURL.Host = manager
	orcaURL.Path = "/api/nodes/authorize"
	orcaURL.Scheme = "https"

	// TODO - figure out what the common failure modes are and set up better messages
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(reqJson))
	if err != nil {
		log.Debug("Failed to build request")
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Debug("Failed to send request")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Debugf("Response code: %d", resp.StatusCode)
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Errorf("Server response: %s", string(body))
			return nil, fmt.Errorf("Failed to add host to UCP: %s", string(body))
		}
		return nil, errors.New("Failed to add host to Orca")
	}

	var nodeConfig orca.NodeConfiguration
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		if err := json.Unmarshal(body, &nodeConfig); err != nil {
			return nil, fmt.Errorf("Failed to add host to UCP: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to add host to UCP: %s", err)
	}
	return &nodeConfig, nil
}
