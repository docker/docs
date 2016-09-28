package client

import (
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// Start the Cert Service
func (c *EngineClient) StartCA(name, volume string, caPort int) error {
	log.Debugf("Starting %s CA service", name)

	// We'll use linking so we don't have to lock down the cert service
	imageName, err := orcaconfig.GetContainerImage(name)
	if err != nil {
		return err
	}
	portMap := make(map[nat.Port]struct{})
	portMap[nat.Port(fmt.Sprintf("%d/tcp", caPort))] = struct{}{}
	bindingMap := nat.PortMap{
		nat.Port(fmt.Sprintf("%d/tcp", caPort)): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", caPort),
			},
		},
	}
	cfg := &container.Config{
		Image:        imageName,
		ExposedPorts: portMap,
		Cmd: strslice.StrSlice{
			"serve",
			"-address", "0.0.0.0",
			"-port", fmt.Sprintf("%d", caPort),
			"-ca", "cert.pem",
			"-ca-key", "key.pem",
			"-tls-cert", path.Join(config.CertDir, config.CertFilename),
			"-tls-key", path.Join(config.CertDir, config.KeyFilename),
			"-mutual-tls-ca", path.Join(config.CertDir, config.CAFilename),
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         name,
		},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/etc/cfssl", volume),
			fmt.Sprintf("%s:%s:ro", config.SwarmNodeCertVolumeName, config.CertDir),
		},
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
	resp, err := c.client.CreateContainer(cfg, hostConfig, name)
	if err != nil {
		return err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch Cert Service: %s", err)
		return err
	}

	if err := WaitForCfssl(fmt.Sprintf("https://%s:%d", config.OrcaHostAddress, caPort), 30*time.Second); err != nil {
		log.Error("CA didn't come up in time")
		return err
	}
	return nil
}
