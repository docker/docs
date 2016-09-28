package client

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// StartAuthWorker starts the Auth Worker container.
func (c *EngineClient) StartAuthWorker() error {
	log.Debug("Starting Auth Worker")

	// The worker server listens on port 4443 internally.
	exposedPorts := map[nat.Port]struct{}{"4443/tcp": {}}
	bindingMap := nat.PortMap{
		nat.Port("4443/tcp"): []nat.PortBinding{
			{
				HostIP: "0.0.0.0",
				// The external Port for the Worker server.
				HostPort: fmt.Sprintf("%d", orcaconfig.AuthWorkerPort),
			},
		},
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.AuthWorkerContainerName)
	if err != nil {
		return err
	}

	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.AuthWorkerCertsVolumeName, "/tls"),
		fmt.Sprintf("%s:%s:rw", config.AuthWorkerDataVolumeName, "/work"),
	}

	authStoreAddr := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.AuthStorePort)
	workerAddr := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.AuthWorkerPort)

	cfg := &container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts,
		Cmd: strslice.StrSlice{
			// Connect to the local auth store.
			fmt.Sprintf("--db-addr=%s", authStoreAddr),
			"--jsonlog", // Always format logs in JSON.
			"worker",
			fmt.Sprintf("--addr=%s", workerAddr),
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.AuthWorkerContainerName,
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

	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.AuthWorkerContainerName)
	if err != nil {
		return err
	}

	containerID := resp.ID

	if err := c.client.StartContainer(containerID); err != nil {
		log.Debugf("Failed to launch Auth Worker: %s", err)
		return err
	}

	return nil
}
