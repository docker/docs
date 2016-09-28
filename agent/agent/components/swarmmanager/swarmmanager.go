package swarmmanager

import (
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	autils "github.com/docker/orca/agent/agent/utils"
	engineClient "github.com/docker/orca/bootstrap/client"
	bconfig "github.com/docker/orca/bootstrap/config"
	butils "github.com/docker/orca/bootstrap/utils"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
	"github.com/docker/orca/utils"
)

type SwarmManager struct {
}

func (p *SwarmManager) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaSwarmManagerContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaControllerContainerName, err)
		return nil
	}

	// Extract the swarm port from the HostConfig of the `ucp-swarm-manager` container
	ports, err := utils.GetHostPortsFromContainerJSON(containerJSON)
	if err != nil {
		return err
	}
	if len(ports) != 1 {
		return err
	}
	currentCfg.SwarmPort = ports[0]
	currentCfg.HostAddress, err = utils.GetHostAddressFromContainerJSON(containerJSON, "--advertise")
	if err != nil {
		return err
	}

	// Create a container entry for the controller
	currentCfg.Containers[config.OrcaSwarmManagerContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *SwarmManager) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaSwarmManagerContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the Swarm Manager is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the Swarm Manager is not running
	if !ok || !ctr.Running {
		return true, nil
	}
	// Determine if image is correct
	expectedCtr, ok := expectedCfg.Containers[containerName]
	if !ok {
		return true, nil
	}
	if expectedCtr.Image != ctr.Image {
		log.Debugf("Image mismatch for %s: %s != %s", containerName, expectedCtr.Image, ctr.Image)
		return true, nil
	}
	// TODO: return true on mismatch of container DNS settings and UCP labels
	return false, nil

}

func (p *SwarmManager) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	requires, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !requires {
		return nil
	}

	// Clean up existing proxy swarm-manager, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaSwarmManagerContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[config.OrcaSwarmManagerContainerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// demote flow - don't recreate the container
		return nil
	}

	// Create swarm-manager argument for the target KV - should be reachable
	kvEndpoint := fmt.Sprintf("etcd://%s:%d", expectedCfg.HostAddress, config.KvPort)

	// Launch the  container
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}
	_, err = butils.DeploySwarmManagerContainer(ec, kvEndpoint)
	if err != nil {
		return err
	}

	if !autils.IsFreshInstall() {
		kvURL, err := url.Parse(kvEndpoint)
		if err != nil {
			return err
		}
		log.Info(kvURL.String())
		orcaEndpoint := fmt.Sprintf("%s:%s", expectedCfg.HostAddress, expectedCfg.ControllerPort)
		swarmEndpoint := fmt.Sprintf("%s:%s", expectedCfg.HostAddress, expectedCfg.SwarmPort)
		proxyEndpoint := fmt.Sprintf("%s:%d", expectedCfg.HostAddress, config.ProxyPort)

		err = bconfig.AddSwarmManager(kvURL, orcaEndpoint, swarmEndpoint, proxyEndpoint)
		if err != nil {
			return err
		}
	}

	log.Info("Swarm Manager component reconciled successfully")
	return nil
}
