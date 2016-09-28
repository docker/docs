package client

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	orcatypes "github.com/docker/orca/types"
)

func sanitizeHostname(hostname string) string {
	r := regexp.MustCompile("([^A-Za-z0-9_])")
	return r.ReplaceAllString(hostname, "_")
}

// StartAuthStore starts the Auth Store container.
func (c *EngineClient) StartAuthStore(peerAddrs ...string) error {
	log.Debug("Starting Auth Store")

	portStr := fmt.Sprintf("%d", orcaconfig.AuthStorePort)
	peerPortStr := fmt.Sprintf("%d", orcaconfig.AuthStorePeerPort)

	portSpec := fmt.Sprintf("%s/tcp", portStr)
	peerPortSpec := fmt.Sprintf("%s/tcp", peerPortStr)

	portMap := map[nat.Port]struct{}{
		nat.Port(portSpec):     {},
		nat.Port(peerPortSpec): {},
	}

	bindingMap := nat.PortMap{
		// Port for client connections.
		nat.Port(portSpec): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: portStr,
			},
		},
		// Port for cluster peer connections.
		nat.Port(peerPortSpec): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: peerPortStr,
			},
		},
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.AuthStoreContainerName)
	if err != nil {
		return err
	}

	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.AuthStoreCertsVolumeName, "/tls"),
		fmt.Sprintf("%s:%s:rw", config.AuthStoreDataVolumeName, "/var/data"),
	}

	serverName := sanitizeHostname(fmt.Sprintf("ucp-auth-store-%s", config.OrcaLocalName))
	canonicalAddr := fmt.Sprintf("%s:%s", config.OrcaHostAddress, peerPortStr)

	cfg := &container.Config{
		Image:        imageName,
		ExposedPorts: portMap,
		Cmd: strslice.StrSlice{
			"--bind", "all", // Listen on all network interfaces.
			"--no-http-admin",   // The admin web console is insecure and unused.
			"--no-update-check", // Disables checking for available updates as well as anonymous usage data collection.
			"--server-name", serverName,
			"--cluster-port", peerPortStr,
			"--driver-port", portStr,
			"--canonical-address", canonicalAddr,
			"--directory", "/var/data/rethinkdb",
			"--driver-tls-key", "/tls/key.pem",
			"--driver-tls-cert", "/tls/cert.pem",
			"--driver-tls-ca", "/tls/ca.pem",
			"--cluster-tls-key", "/tls/key.pem",
			"--cluster-tls-cert", "/tls/cert.pem",
			"--cluster-tls-ca", "/tls/ca.pem",
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.AuthStoreContainerName,
		},
	}

	for _, peerAddr := range peerAddrs {
		cfg.Cmd = append(cfg.Cmd, "--join", peerAddr)
	}

	hostConfig := &container.HostConfig{
		Binds:        mounts,
		PortBindings: bindingMap,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.AuthStoreContainerName)
	if err != nil {
		return err
	}

	containerID := resp.ID

	if err := c.client.StartContainer(containerID); err != nil {
		log.Debugf("Failed to launch Auth Store: %s", err)
		return err
	}

	return nil
}

// ScaleDownAuthStore stops the local Auth Storage container and connects to
// another node in the cluster (if this is not the last one) to sync the
// database with the new replication count.
// NOTE: this should be called *before* the local kv store config has been
// unwound. The given kvStoreURL should be the URL of the local KV Store.
func (c *EngineClient) ScaleDownAuthStore(kvStoreURL *url.URL) error {
	// Infer what our local IP is from the local kvStore URL.
	localIP, _, err := net.SplitHostPort(kvStoreURL.Host)
	if err != nil {
		log.Errorf("unable to split host/port from kv store address: %s", kvStoreURL.Host)
		return err
	}

	// Need to get IP addresses of remote controller nodes.
	controllers, err := config.GetControllers(kvStoreURL)
	if err != nil {
		log.Errorf("unable to list controllers: %s", err)
		return err
	}

	if len(controllers) == 0 {
		// We would expect there to be at least one controller
		// registered in the KV store representing this node. This
		// would indicate a problem with configuration of the UCP
		// cluster.
		return fmt.Errorf("unable to find any controllers registered in the kv store")
	}

	controllerIPs := orcatypes.GetIPsFromControllers(controllers)
	log.Debugf("controller IP addresses: %s", controllerIPs)

	if len(controllerIPs) < 2 {
		log.Debug("not draining local auth store because it appears to be the last node remaining")
		return nil
	}

	authStoreAddrs := make([]string, 0, len(controllers))
	for _, controllerIP := range controllerIPs {
		authStoreAddr := net.JoinHostPort(controllerIP, fmt.Sprintf("%d", orcaconfig.AuthStorePort))
		authStoreAddrs = append(authStoreAddrs, authStoreAddr)
	}

	return c.RunAuthDrainDBServer(localIP, uint(orcaconfig.AuthStorePeerPort), authStoreAddrs...)
}

func (c *EngineClient) RunAuthDrainDBServer(serverHostname string, serverPeerPort uint, authStoreAddrs ...string) error {
	log.Debug("Running Auth DrainDBServer")

	// There must be at least one db driver address given. If there are
	// more in the cluster they will automatically be discovered by the
	// client.
	if len(authStoreAddrs) == 0 {
		return errNoAuthStoreAddrsGiven
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.AuthDrainDBServerContainerName)
	if err != nil {
		return err
	}

	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.AuthAPICertsVolumeName, "/tls"),
	}

	dbAddrArgs := make([]string, len(authStoreAddrs))
	for i, authStoreAddr := range authStoreAddrs {
		dbAddrArgs[i] = fmt.Sprintf("--db-addr=%s", authStoreAddr)
	}

	cmd := append(dbAddrArgs,
		fmt.Sprintf("--debug=%t", log.GetLevel() == log.DebugLevel),
		"--jsonlog",
		"drain-db-server",
		fmt.Sprintf("--hostname=%s", serverHostname),
		fmt.Sprintf("--cluster-port=%d", serverPeerPort),
	)

	cfg := &container.Config{
		Image: imageName,
		Tty:   c.IsTty(),
		Cmd:   cmd,
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.AuthSyncDBContainerName,
		},
	}

	hostConfig := &container.HostConfig{
		Binds:      mounts,
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	createResp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.AuthSyncDBContainerName)
	if err != nil {
		log.Debugf("Failed to create auth sync-db container: %s", err)
		return err
	}

	containerID := createResp.ID

	// Attach to the container and pass through all the container
	// output.
	if err := c.streamContainerOutput(orcaconfig.AuthSyncDBContainerName, c.IsTty()); err != nil {
		log.Errorf("Failed to stream container logs: %s", err)
		return err
	}

	if err := c.client.StartContainer(containerID); err != nil {
		log.Debugf("Failed to start auth sync-db container: %s", err)
		return err
	}

	// Wait for the drain-db-server task to complete. How long it takes
	// depends on whether or not data needs to be moved from the target
	// server to another. Usually, data shouldn't need to be moved as we
	// are only moving data off of this server that is already replicated
	// to others. This should be pretty quick but we don't want to wait
	// any longer than 2 minutes.
	timeout := 2 * time.Minute
	done := make(chan error, 1)

	timer := time.AfterFunc(timeout, func() {
		done <- fmt.Errorf("auth sync-db timed out after %s", timeout)
	})

	go func() {
		exitStatus, err := c.client.ContainerWait(containerID)
		timer.Stop()

		if err != nil {
			err = fmt.Errorf("unable to wait for auth drain-db-server container: %s", err)
		} else if exitStatus != 0 {
			err = fmt.Errorf("auth drain-db-server returned non-zero exit status: %d", exitStatus)
		}

		done <- err
	}()

	if err := <-done; err != nil {
		return err
	}

	c.client.RemoveContainer(containerID, true, false)

	return nil
}
