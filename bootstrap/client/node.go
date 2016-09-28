package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/orca/config"
)

type NodeType int

const (
	NonUCP NodeType = iota
	Controller
	Worker
)

func (c *EngineClient) DetectNodeType() NodeType {
	// Note: get all containers so we can detect stopped controllers and correctly identify them
	containers, err := c.client.ListContainers(types.ContainerListOptions{All: true, Size: false})
	if err != nil {
		log.Warnf("Failed to retrieve container list: %s", err)
		return NonUCP
	}

	anyUCPFound := false
	// Search for controller container
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+config.OrcaControllerContainerName {
				return Controller
			} else if name == "/"+config.OrcaSwarmJoinContainerName || name == "/"+config.OrcaProxyContainerName {
				// NOTE: this will get trickier to detect in swarm v2
				anyUCPFound = true
			}
		}
	}
	if anyUCPFound {
		return Worker
	}
	return NonUCP
}
