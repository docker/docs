package kv

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	engineClient "github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/utils"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

type KV struct {
}

func (p *KV) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaKvContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaKvContainerName, err)
		return nil
	}
	// Create a container entry for the kv container
	currentCfg.Containers[config.OrcaKvContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *KV) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaKvContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the KV is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the KV is not running
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

func (p *KV) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig, chDone chan bool) error {
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// Clean up existing KV container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaKvContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[config.OrcaKvContainerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// TODO: demote flow - gracefully spin down KV
		return nil
	}

	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}

	// Get a list of managers where this node is excluded
	// TODO: DRY
	filteredMgrs := []string{}
	for _, manager := range expectedCfg.Managers {
		if manager != expectedCfg.HostAddress {
			filteredMgrs = append(filteredMgrs, manager)
		}
	}
	log.Info(filteredMgrs)

	targetKV := ""
	if len(filteredMgrs) > 0 {
		targetKV = filteredMgrs[0]
	}
	log.Infof("Deploying KV Container with target KV of \"%s\"", targetKV)

	_, err = utils.DeployKVContainer(ec, targetKV, chDone)
	if err != nil {
		return err
	}
	log.Info("KV Deployed")

	// HACK: Restart the swarm-join container, if it exists
	timeout := time.Minute
	_ = dclient.ContainerRestart(context.TODO(), config.OrcaSwarmJoinContainerName, &timeout)

	return nil
}
