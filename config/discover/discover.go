package discover

import (
	"github.com/docker/engine-api/client"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/agent/agent/components"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

func DiscoverNodeConfig(dclient *client.Client) (*types.NodeConfig, error) {
	cfg := &types.NodeConfig{
		Containers: make(map[string]*types.OrcaContainer),
	}
	for _, component := range components.ComponentList {
		err := component.BuildCurrentConfig(dclient, cfg)
		if err != nil {
			return cfg, err
		}
	}

	// The current node config represents a manager if the controller is present
	if _, found := cfg.Containers[config.OrcaControllerContainerName]; found {
		log.Debugf("Detected %s, marking node as a controller", config.OrcaControllerContainerName)
		cfg.IsManager = true
	} else {
		log.Debugf("Did not detect %s, marking node as a worker", config.OrcaControllerContainerName)
	}
	return cfg, nil
}
