package swarmjoin

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	engineClient "github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

type SwarmJoin struct {
}

func (p *SwarmJoin) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaSwarmJoinContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaSwarmJoinContainerName, err)
		return nil
	}

	// Get the last argument and loose the "etcd://" prefix
	lastArgument := containerJSON.Config.Cmd[len(containerJSON.Config.Cmd)-1][7:]

	// Split by comma to get the list of advertise targets, which represent the current managers
	advertiseTargets := strings.Split(lastArgument, ",")
	for _, target := range advertiseTargets {
		currentCfg.Managers = append(currentCfg.Managers, strings.Split(target, ":")[0])
	}

	currentCfg.Containers[config.OrcaSwarmJoinContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *SwarmJoin) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaSwarmJoinContainerName
	// TODO: detect when a new CSR is triggered and return yes
	ctr, ok := currentCfg.Containers[containerName]
	if !ok || !ctr.Running {
		// The swarm-join container is either not present or not running
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
	// TODO: return true on mismatch of managers and args
	return false, nil
}

func (p *SwarmJoin) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	requires, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !requires {
		return nil
	}

	// Clean up existing proxy swarm-join, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaSwarmJoinContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	// Create swarm-join argument of target KV stores - some may be unavailable at launch time
	// TODO: make sure this won't break anything
	kvStoreList := []string{}
	for _, manager := range expectedCfg.Managers {
		kvStoreList = append(kvStoreList, fmt.Sprintf("%s:%d", manager, config.KvPort))
	}
	kvStores := fmt.Sprintf("etcd://%s", strings.Join(kvStoreList, ","))

	// Create swarm-join argument of the target proxy - should be available
	proxyEndpoint := fmt.Sprintf("%s:%d", expectedCfg.HostAddress, config.ProxyPort)

	// Launch swarmJoin container
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}
	err = ec.StartSwarmJoin(kvStores, proxyEndpoint, []string{})
	if err != nil {
		return err
	}
	log.Info("Swarm Join component reconciled successfully")
	return nil
}
