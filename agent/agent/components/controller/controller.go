package controller

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

type Controller struct {
}

func (p *Controller) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaControllerContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaControllerContainerName, err)
		return nil
	}

	// Extract the controller port from the HostConfig of the `ucp-controller` container
	ports, err := utils.GetHostPortsFromContainerJSON(containerJSON)
	if err != nil {
		return err
	}
	if len(ports) != 1 {
		return err
	}
	currentCfg.ControllerPort = ports[0]

	// Create a container entry for the controller
	currentCfg.Containers[config.OrcaControllerContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *Controller) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaControllerContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the controller is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the controller is not running
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

func (p *Controller) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// Clean up existing Controller container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaControllerContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[config.OrcaControllerContainerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// demote flow - don't recreate the container
		return nil
	}

	// Launch Controller container
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}

	kvEndpoint := fmt.Sprintf("etcd://%s:%d", expectedCfg.HostAddress, config.KvPort)
	orcaEndpoint := fmt.Sprintf("%s:%s", expectedCfg.HostAddress, expectedCfg.ControllerPort)
	swarmEndpoint := fmt.Sprintf("%s:%s", expectedCfg.HostAddress, expectedCfg.SwarmPort)
	proxyEndpoint := fmt.Sprintf("%s:%d", expectedCfg.HostAddress, config.ProxyPort)
	kvURL, err := url.Parse(kvEndpoint)
	if err != nil {
		return err
	}

	if autils.IsFreshInstall() {
		log.Info("Putting initial Controller configuration in the KV store")
		// Write out the initial configuration for Orca
		adminUsername, adminPassword, err := autils.GetUCPCredentials()
		if err != nil {
			return err
		}
		bconfig.OrcaInstanceKey, err = autils.GetUCPInstanceKey()
		if err != nil {
			return err
		}
		// TODO: usage/tracking settings
		if err := bconfig.BootstrapOrcaConfig(kvURL, adminUsername, adminPassword,
			orcaEndpoint, swarmEndpoint, proxyEndpoint, true, true); err != nil {
			log.Error("Failed to set up initial UCP configuration")
			return err
		}
	}

	log.Info("Deploying UCP Controller Container")
	err = butils.DeployControllerContainer(ec, kvEndpoint, swarmEndpoint, proxyEndpoint, orcaEndpoint)
	if err != nil {
		return err
	}

	log.Info("Controller component reconciled succesfully")
	return nil
}
