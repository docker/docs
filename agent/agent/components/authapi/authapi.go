package authapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	engineClient "github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/utils"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

type AuthAPI struct {
}

func (p *AuthAPI) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.AuthAPIContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.AuthAPIContainerName, err)
		return nil
	}

	// Create a container entry for the auth api container
	currentCfg.Containers[config.AuthAPIContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *AuthAPI) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.AuthAPIContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the Auth API is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the Auth API Container is not running
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

func (p *AuthAPI) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	containerName := config.AuthAPIContainerName
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// Clean up existing Auth API container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		containerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// demote flow - don't recreate the container
		return nil
	}

	// Launch Auth API
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}
	if err := utils.DeployAuthAPIContainer(ec); err != nil {
		return err
	}

	log.Info("Auth API component reconciled succesfully")
	return nil

}
