package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
)

// Deploy a container based on an existing ContainerInfo dump, which is already "fixed up"
func (c *EngineClient) DeployFromInfo(info *types.ContainerJSON) error {
	log.Debugf("Starting %s", info.Name)

	resp, err := c.client.CreateContainer(info.Config, info.HostConfig, info.Name)
	if err != nil {
		return err
	}
	containerId := resp.ID
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch %s: %s", info.Name, err)
		return err
	}
	return nil
}
