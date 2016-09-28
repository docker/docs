package client

import (
	"errors"
	"fmt"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

var ErrNoKvStoreFound = errors.New("Could not determine KV store's location.")

// Start the Orca Server
func (c *EngineClient) StartOrcaServer(kvEndpoint, swarmEndpoint, engineEndpoint string) error {

	log.Debug("Starting UCP controller")

	bindingMap := nat.PortMap{
		nat.Port("8080/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", orcaconfig.OrcaPort),
			},
		},
	}

	portMap := make(map[nat.Port]struct{})
	portMap["8080/tcp"] = struct{}{}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaControllerContainerName)
	if err != nil {
		return err
	}
	mounts := []string{
		fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, filepath.Join(config.CertDir, "orca")),
		fmt.Sprintf("%s:%s:ro", config.SwarmControllerCertVolumeName, filepath.Join(config.CertDir, "swarm")),
	}
	cfg := &container.Config{
		Hostname:     "ucp-controller-" + config.OrcaLocalName,
		Image:        imageName,
		ExposedPorts: portMap,
		Cmd: strslice.StrSlice{
			//"--debug",
			"server",
			"--discovery", kvEndpoint,
			"--swarm-url", fmt.Sprintf("tcp://%s", swarmEndpoint), // Pin to the local swarm V1 manager
			"--docker-url", fmt.Sprintf("tcp://%s", engineEndpoint), // Pin to the local engine manager
			"--orca-tls-ca-cert", filepath.Join(config.CertDir, "orca", config.CAFilename),
			"--orca-tls-cert", filepath.Join(config.CertDir, "orca", config.CertFilename),
			"--orca-tls-key", filepath.Join(config.CertDir, "orca", config.KeyFilename),
			"--tls-ca-cert", filepath.Join(config.CertDir, "swarm", config.CAFilename),
			"--tls-cert", filepath.Join(config.CertDir, "swarm", config.CertFilename),
			"--tls-key", filepath.Join(config.CertDir, "swarm", config.KeyFilename),
			"--discovery-tls-ca-cert", filepath.Join(config.CertDir, "swarm", config.CAFilename),
			"--discovery-tls-cert", filepath.Join(config.CertDir, "swarm", config.CertFilename),
			"--discovery-tls-key", filepath.Join(config.CertDir, "swarm", config.KeyFilename),
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
		},
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
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaControllerContainerName)
	if err != nil {
		return err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch UCP: %s", err)
		return err
	}

	return nil
}

func (c *EngineClient) AutodetectControllers() ([]string, error) {
	log.Debugf("Trying to find KV store from controller...")
	kvStoreURL, err := c.FindKv()
	if err != nil {
		// Check the swarm node if there is no controller
		log.Debugf("Trying to find KV store from swarm...")
		kvStoreURL, err = c.FindKvNonController()
		if err != nil {
			// Shouldn't happen, doubly so on non-replica join
			// (because we've found it earlier in the join process)
			return []string{}, ErrNoKvStoreFound
		}
	}

	controllers, err := config.GetControllers(kvStoreURL)
	if err != nil {
		return []string{}, err
	}
	controllerAddresses := types.GetIPsFromControllers(controllers)
	log.Debugf("Auto-detected controllers: %s", controllerAddresses)

	return controllerAddresses, nil
}
