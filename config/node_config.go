package config

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	orcatypes "github.com/docker/orca/types"
)

func GetExpectedNodeConfig(dclient *client.Client, clcfg orcatypes.ClusterConfig) (*orcatypes.NodeConfig, error) {
	log.Debug("Obtaining docker info")
	info, err := dclient.Info(context.TODO())
	if err != nil {
		return nil, err
	}

	log.Debug("Extracting the set of managers from info")
	remoteManagerIPs, err := GetManagers(info)
	if err != nil {
		return nil, err
	}

	hostAddress := info.Swarm.NodeAddr
	log.Debugf("Using swarm defined host address %s", hostAddress)

	log.Debug("Initializing list of expected containers")
	containers := make(map[string]*orcatypes.OrcaContainer)
	for _, ctrName := range ManagerContainerNames {
		expectedImage, err := GetContainerImage(ctrName)
		if err != nil {
			return nil, err
		}
		containers[ctrName] = &orcatypes.OrcaContainer{
			Image:   expectedImage,
			Running: true,
		}
	}

	if !info.Swarm.ControlAvailable {
		// If this node is not a swarm-mode manager, mark the expected state of the
		// manager containers as not running
		for _, ctrName := range AgentContainerNames {
			ctr, _ := containers[ctrName]
			ctr.Running = false
		}
	}

	return &orcatypes.NodeConfig{
		ClusterConfig: clcfg,
		IsManager:     info.Swarm.ControlAvailable,
		CertsExpiring: false, // TODO
		Managers:      remoteManagerIPs,
		HostAddress:   hostAddress,
		Containers:    containers,
	}, nil
}

func GetManagers(info types.Info) ([]string, error) {
	managers := []string{}
	for _, peer := range info.Swarm.RemoteManagers {
		mgr, _, err := net.SplitHostPort(peer.Addr)
		if err != nil {
			return managers, fmt.Errorf("unable to parse Manager address from cluster config: %s", err)
		}
		managers = append(managers, mgr)
	}
	return managers, nil
}
