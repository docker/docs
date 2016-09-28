package proxy

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

type Proxy struct {
}

func (p *Proxy) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.OrcaProxyContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.OrcaProxyContainerName, err)
		return nil
	}

	// Extract the UCP Instance ID
	instanceID, ok := containerJSON.Config.Labels[config.UCPInstanceIDLabelKey]
	if ok {
		currentCfg.ClusterConfig.UCPInstanceID = instanceID
	}

	// Populate DNS configs
	currentCfg.ClusterConfig.DNS = containerJSON.HostConfig.DNS
	currentCfg.ClusterConfig.DNSOpt = containerJSON.HostConfig.DNSOptions
	currentCfg.ClusterConfig.DNSSearch = containerJSON.HostConfig.DNSSearch

	// Create a container entry for the proxy
	currentCfg.Containers[config.OrcaProxyContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *Proxy) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.OrcaProxyContainerName
	// TODO: detect when a new CSR is triggered and return yes
	ctr, ok := currentCfg.Containers[containerName]
	if !ok || !ctr.Running {
		// The proxy is either not present or not running
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

func (p *Proxy) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	requires, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !requires {
		return nil
	}

	// Clean up existing proxy container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.OrcaProxyContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	// Launch proxy container
	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}
	_, err = utils.DeployProxyContainer(ec)
	if err != nil {
		return err
	}
	log.Info("Proxy component reconciled succesfully")
	return nil
}
