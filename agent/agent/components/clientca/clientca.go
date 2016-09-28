package clientca

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	engineClient "github.com/docker/orca/bootstrap/client"
	bconfig "github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

type ClientCA struct {
}

func (p *ClientCA) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaCAContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaCAContainerName, err)
		return nil
	}

	// Create a container entry for the client CA
	currentCfg.Containers[config.OrcaCAContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *ClientCA) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaCAContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the Client CA is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the Client CA is not running
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

func (p *ClientCA) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// Clean up existing Client CA container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaCAContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[config.OrcaCAContainerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// demote flow - don't recreate the container
		return nil
	}

	// Launch client CA container
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}

	err = ec.StartCA(config.OrcaCAContainerName, bconfig.OrcaRootCAVolumeName, config.OrcaCAPort)
	if err != nil {
		return err
	}

	log.Info("Client CA component reconciled succesfully")
	return nil
}
