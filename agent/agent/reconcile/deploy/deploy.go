package deploy

import (
	"fmt"
	"net"
	"net/url"

	log "github.com/Sirupsen/logrus"

	engineClient "github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
)

// DeployAllManagerContainers sets up all containers for the first manager ofthe cluster
func DeployAllManagerContainers(ec *engineClient.EngineClient, adminUsername, adminPassword string) error {
	// Note: The lock is technically unused in this use-case since we won't lock on the initial creation
	chDone := make(chan bool)
	defer func() {
		chDone <- true
		log.Debug("closing channel for kv lock")
		close(chDone)
	}()

	// Verify filesystem permissions
	if err := utils.VerifyPermissions(true); err != nil {
		return err
	}

	log.Info("Deploying KV Store")
	kvStoreURL, err := utils.DeployKVContainer(ec, "", chDone)
	if err != nil {
		return err
	}
	log.Infof("Internal KV started at %s", kvStoreURL)

	log.Info("Deploying Engine Proxy")
	proxyEndpoint, err := utils.DeployProxyContainer(ec)
	if err != nil {
		log.Errorf(`Failed to start proxy.  Run "docker logs %s" for more details`, orcaconfig.OrcaProxyContainerName)
	}

	log.Info("Deploying Swarm Join")
	if err := ec.StartSwarmJoin(kvStoreURL.String(), proxyEndpoint, []string{}); err != nil {
		log.Errorf(`Failed to start swarm join.  Run "docker logs %s" for more details`, orcaconfig.OrcaSwarmJoinContainerName)
		return err
	}

	// Deploy Swarm Manager
	log.Info("Deploying Swarm Manager")
	managerEndpoint, err := utils.DeploySwarmManagerContainer(ec, kvStoreURL.String())
	if err != nil {
		return err
	}

	// Deploy the two CAs
	log.Info("Deploying CAs")
	if err := utils.DeployCAContainers(ec); err != nil {
		return err
	}

	// Deploy Auth service containers.
	log.Info("Deploying Auth Store")
	if err := utils.DeployAuthStoreContainer(ec); err != nil {
		return err
	}
	log.Info("Deploying Auth API")
	if err := utils.DeployAuthAPIContainer(ec); err != nil {
		return err
	}
	log.Info("Deploying Auth Worker")
	if err := utils.DeployAuthWorkerContainer(ec); err != nil {
		return err
	}

	// Create the initial admin user.
	if err := ec.RunAuthCreateAdmin(adminUsername, adminPassword); err != nil {
		log.Errorf(`Failed to create initial admin user.  Run "docker logs %s" for more details`, orcaconfig.AuthCreateAdminContainerName)
		return err
	}
	orcaEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.OrcaPort)

	// Write out the initial configuration for Orca
	// TODO: disable tracking & usage
	if err := config.BootstrapOrcaConfig(kvStoreURL, adminUsername, adminPassword,
		orcaEndpoint, managerEndpoint, proxyEndpoint, true, true); err != nil {
		log.Error("Failed to set up initial UCP configuration")
		return err
	}

	log.Info("Deploying UCP Controller")
	// Finally, deploy the UCP server replica.
	return utils.DeployControllerContainer(ec, kvStoreURL.String(), managerEndpoint, proxyEndpoint,
		orcaEndpoint)
}

func DeployAgentContainers(ec *engineClient.EngineClient, swarmArgs []string) error {
	log.Info("Deploying Engine Proxy")
	proxyEndpoint, err := utils.DeployProxyContainer(ec)
	if err != nil {
		log.Errorf(`Failed to start proxy.  Run "docker logs %s" for more details`, orcaconfig.OrcaProxyContainerName)
	}

	// The first swarm argument should be a URL for the KV store.
	kvStoreURL, err := url.Parse(swarmArgs[len(swarmArgs)-1])
	if err != nil {
		log.Error("unable to parse first swarm argument as a kv store URL")
		return err
	}

	// Deploy Swarm join
	if err := ec.StartSwarmJoin(kvStoreURL.String(), proxyEndpoint, nil); err != nil {
		log.Error("Failed to start swarm")
		return err
	}
	return nil
}

func releaseLockChannel(c chan bool) {
	c <- true
	log.Debug("closing channel for kv lock")
	close(c)
}

func managerCriticalSection(ec *engineClient.EngineClient, remoteKV string, managers []string) (*url.URL, error) {
	chDone := make(chan bool)

	peerAddrs := make([]string, len(managers))
	for i, peerIP := range managers {
		peerAddrs[i] = net.JoinHostPort(peerIP, fmt.Sprintf("%d", orcaconfig.AuthStorePeerPort))
	}

	defer releaseLockChannel(chDone)
	log.Info("Deploying KV Store")
	kvStoreURL, err := utils.DeployKVContainer(ec, remoteKV, chDone)
	if err != nil {
		return nil, err
	}
	log.Infof("KV started at %s", kvStoreURL)

	log.Info("Deploying Auth Store")
	if err := utils.DeployAuthStoreContainer(ec, peerAddrs...); err != nil {
		return nil, err
	}

	return kvStoreURL, nil
}

// DeployManagerOnlyContainers sets up the manager containers when promoting from an agent node
func DeployManagerOnlyContainers(ec *engineClient.EngineClient, kvURL string, managers []string) error {
	proxyEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.ProxyPort)

	// Verify filesystem permissions
	if err := utils.VerifyPermissions(true); err != nil {
		return err
	}

	// Deploy the KV Store and Auth Store within a cluster-wide lock
	log.Infof("Entering critical section. Remote KV is at %s", kvURL)
	kvStoreURL, err := managerCriticalSection(ec, kvURL, managers)
	if err != nil {
		return err
	}

	// Deploy the two CAs
	log.Info("Deploying CAs")
	if err := utils.DeployCAContainers(ec); err != nil {
		return err
	}

	// Deploy Swarm Manager
	log.Info("Deploying Swarm Manager")
	managerEndpoint, err := utils.DeploySwarmManagerContainer(ec, kvStoreURL.String())
	if err != nil {
		return err
	}

	// Add Swarm Manager and controller entries in the KV store
	orcaEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.OrcaPort)
	if err := config.AddSwarmManager(kvStoreURL, orcaEndpoint, managerEndpoint, proxyEndpoint); err != nil {
		return err
	}

	log.Info("Deploying Auth API")
	if err := utils.DeployAuthAPIContainer(ec); err != nil {
		return err
	}
	log.Info("Deploying Auth Worker")
	if err := utils.DeployAuthWorkerContainer(ec); err != nil {
		return err
	}

	if err := config.UpdateUCPServiceAuthConfig(kvStoreURL); err != nil {
		return fmt.Errorf("unable to update auth configuration: %s", err)
	}

	log.Info("Deploying UCP Controller")
	return utils.DeployControllerContainer(ec, kvStoreURL.String(), managerEndpoint, proxyEndpoint,
		orcaEndpoint)
}
