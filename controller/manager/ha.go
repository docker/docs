package manager

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"

	"github.com/docker/orca/config"
	"github.com/docker/orca/config/discover"
)

// Represents one of the replicated cluster manager nodes
type ManagerNode struct {
	// TODO - consider moving this up to types and reconciling it with the other node/cluster types
	Name                   string
	ControllerURL          string
	SwarmClassicManagerURL string
	EngineProxyURL         string
	KVURL                  string
	AuthProviderURL        string
}

// Stored in the KV store to represent the controllers in the cluster
const (
	Healthy  = "Healthy"
	Degraded = "Degraded"
	Failed   = "Failed"
	// Use two levels of timeouts so we can timeout at the correct level
	// (caller of _ping versus swarm and KV store)
	PingTestTimeout         = 5 * time.Second
	InnerServiceTestTimeout = PingTestTimeout - 500*time.Millisecond
)

// Load up ourself for healthchecking locally
func (m DefaultManager) verifySelfManagerPresent() {
	dclient := m.ProxyClient()
	// Figure out the cluster config from ourself
	selfManager, err := discover.DiscoverNodeConfig(dclient)
	if err != nil {
		log.Errorf("Failed to get current manager config: %s", err)
		return
	}

	m.selfManagerNode.Name = selfManager.HostAddress // TODO Inspect the node so we can get a friendly name for it
	m.selfManagerNode.ControllerURL = fmt.Sprintf("https://%s:%s", selfManager.HostAddress, selfManager.ControllerPort)
	m.selfManagerNode.SwarmClassicManagerURL = fmt.Sprintf("tcp://%s:%s", selfManager.HostAddress, selfManager.SwarmPort)
	m.selfManagerNode.EngineProxyURL = fmt.Sprintf("tcp://%s:%d", selfManager.HostAddress, config.ProxyPort) // TODO - this is likely to change
	m.selfManagerNode.KVURL = fmt.Sprintf("etcd://%s:%d", selfManager.HostAddress, config.KvPort)
	m.selfManagerNode.AuthProviderURL = fmt.Sprintf("https://%s:%d", selfManager.HostAddress, config.AuthAPIPort)
}

// Return the list of managers
func (m DefaultManager) GetManagers() []*ManagerNode {
	nodes := []*ManagerNode{}
	dclient := m.ProxyClient()

	info, err := dclient.Info(context.TODO())
	if err != nil {
		log.Errorf("Failed to retrieve cluster info: %s", err)
		return nodes
	}
	managers, err := config.GetManagers(info)
	if err != nil {
		log.Errorf("Failed to get managers: %s", err)
		return nodes
	}

	// Figure out the cluster config from ourself
	selfManager, err := discover.DiscoverNodeConfig(dclient)
	if err != nil {
		log.Errorf("Failed to get current manager config: %s", err)
		return nodes
	}

	// Now we can build up the details for each manager:
	for _, manager := range managers {
		node := &ManagerNode{
			Name:                   manager, // TODO Inspect the node so we can get a friendly name for it
			ControllerURL:          fmt.Sprintf("https://%s:%s", manager, selfManager.ControllerPort),
			SwarmClassicManagerURL: fmt.Sprintf("tcp://%s:%s", manager, selfManager.SwarmPort),
			EngineProxyURL:         fmt.Sprintf("tcp://%s:%d", manager, config.ProxyPort), // TODO - this is likely to change
			KVURL:                  fmt.Sprintf("etcd://%s:%d", manager, config.KvPort),
			AuthProviderURL:        fmt.Sprintf("https://%s:%d", manager, config.AuthAPIPort),
		}

		nodes = append(nodes, node)
	}
	return nodes

}

// Report an overall status for this node
// Note an edge case: since we delegate checking of swarm and KV store to the
// node in question, we cannot report a situation where the controller is down
// but swarm and/or KV store are up
func (m *DefaultManager) GetStatus(node *ManagerNode) string {
	log.Debugf("Checking status for Orca node %s", node.Name)

	client := &http.Client{
		Timeout: PingTestTimeout, // Very short timeout for health check
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TODO - remove this and replace with trusting the CA
			},
		},
	}
	// Ping the controller
	resp, err := client.Get(fmt.Sprintf("%s/_ping", node.ControllerURL))
	if err != nil {
		// Controller is down
		return Failed
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return Healthy
	} else {
		// This means KV store and/or swarm is down but controller is up
		return Degraded
	}
}

func (m DefaultManager) GetSelfStatus() error {
	node := m.selfManagerNode
	if node.ControllerURL == "" {
		return fmt.Errorf("Internal Error")
	}

	// Parallelize the tests so we get it done as quickly as possible
	errChan := make(chan error, 4)
	var wg sync.WaitGroup
	wg.Add(1)

	// Test the Classic Swarm Manager
	go func() {
		defer wg.Done()
		err := m.getEngineStatus(node.SwarmClassicManagerURL, "Classic Swarm Manager")
		if err != nil {
			// bubble up errors to log
			log.Warnf("Swarm status check error: %s", err)
			errChan <- err
		}
	}()

	wg.Add(1)
	// Test the Swarm Mode Manager
	go func() {
		defer wg.Done()
		// Attempt to _ping the local engine
		err := m.getEngineStatus(node.EngineProxyURL, "Engine Proxy")
		if err != nil {
			// bubble up errors to log
			log.Warnf("Engine proxy status check error: %s", err)
			errChan <- err
		}

	}()

	wg.Add(1)
	go func() {
		// Get the local node's status from the Swarm section of /info
		defer wg.Done()
		err := m.getSwarmModeStatus()
		if err != nil {
			log.Warnf("Swarm mode status check error: %s", err)
			errChan <- err
		}
	}()

	// Test the KV store
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := m.getKVStatus(node)
		if err != nil {
			log.Warnf("KV store status check error: %s", err)
			errChan <- err
		}
	}()

	// Test the auth API provider
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := m.getAuthProviderStatus(node)
		if err != nil {
			log.Warnf("Auth API provider status check error: %s", err)
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return fmt.Errorf("ERR")
	}
	return nil
}

func (m DefaultManager) getAuthProviderStatus(node *ManagerNode) error {
	log.Debugf("Checking Auth API Provider status for Orca node %s", node.Name)

	if node.AuthProviderURL == "" {
		return fmt.Errorf("Auth API Provider URL is unknown")
	}

	providerURL, err := url.Parse(node.AuthProviderURL)
	if err != nil {
		return fmt.Errorf("unable to parse Auth Provider URL: %s", err)
	}
	providerURL.Path = "/enzi/_ping"

	resp, err := m.httpClient.Get(providerURL.String())
	if err != nil {
		return fmt.Errorf("Auth API Provider error: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Auth API Provider error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Auth API Provider error: %s", body)
	}

	return nil
}

// Etcd struct response for the /health endpoint
type KVHealthResp struct {
	Health string `json:"health,omitempty"`
}

func (m DefaultManager) getKVStatus(node *ManagerNode) error {
	// Replace the etcd:// prefix with https://
	etcdURL, err := url.Parse(node.KVURL)
	if err != nil {
		return fmt.Errorf("Unable to determine KV Store url %s", err)
	}
	etcdURL.Scheme = "https"
	etcdURL.Path = "/health"

	resp, err := m.httpClient.Get(etcdURL.String())
	if err != nil {
		return fmt.Errorf("Unable to reach KV Store at %s: %s", etcdURL.String(), err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Unable to read KV Store response at %s: %s", etcdURL.String(), err)
	}
	var healthResp KVHealthResp
	err = json.Unmarshal(body, &healthResp)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal KV Store response at %s: %s", etcdURL.String(), err)
	}

	if healthResp.Health != "true" {
		return fmt.Errorf("KV Store reported unhealthy status")
	}
	return nil
}

// Report the status of a Docker Remote API engine for this node
// This will be invoked for both the Swarm V1 manager and the
// underlying Docker Engine
func (m DefaultManager) getEngineStatus(engineURL, errName string) error {
	log.Debugf("Checking status for %s", errName)

	swarmURL, err := url.Parse(engineURL)
	if err != nil {
		return fmt.Errorf("%s health check URL parsing error: %s", errName, err)
	}
	swarmURL.Scheme = "https"
	swarmURL.Path = "/_ping"

	resp, err := m.httpClient.Get(swarmURL.String())
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("%s health check error: %s", errName, err)
		}
		if resp.StatusCode == 200 {
			return nil
		} else {
			return fmt.Errorf("%s health check error: %s", errName, body)
		}
	} else {
		return fmt.Errorf("%s error: %s", errName, err)
	}
}

func (m DefaultManager) getSwarmModeStatus() error {
	info, err := m.proxyClient.Info(context.TODO())
	if err != nil {
		return fmt.Errorf("Swarm Mode Manager health check error: info: %s", err)
	}

	if info.Swarm.LocalNodeState != swarm.LocalNodeStateActive {
		return fmt.Errorf("Swarm Mode Manager reported a local node state of %s", info.Swarm.LocalNodeState)
	}
	return nil
}

func (m *DefaultManager) getHABanners() []Banner {
	res := []Banner{}

	// Note, this will delay at most 2 seconds per call in the case of an unhealthy system
	managers := m.GetManagers()

	if len(managers) == 1 {
		// HA health status is moot for a single controller setup
		return res
	} else if len(managers) == 2 {
		res = append(res, Banner{
			Level:   BannerWARN,
			Message: "HA Degraded: You appear to only have two controller nodes, which does not provide high availability.  Please add a 3rd controller for a miminum HA configuration.  For more information visit https://success.docker.com/?cid=UCP1",
		})
	}

	anyBad := false
	var wg sync.WaitGroup
	for _, mn := range managers {
		wg.Add(1)
		go func(mn *ManagerNode) {
			defer wg.Done()
			status := m.GetStatus(mn)
			if status != Healthy {
				anyBad = true
			}
		}(mn)
	}
	wg.Wait()
	if anyBad {
		res = append(res,
			Banner{
				Level:   BannerCRIT,
				Message: "HA Degraded: One or more controllers are unhealthy.  For more information visit https://success.docker.com/?cid=UCP2",
			})
	}

	return res
}
