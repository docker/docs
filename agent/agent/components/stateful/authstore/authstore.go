package authstore

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	engineTypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"

	"github.com/docker/orca/agent/agent/utils"
	engineClient "github.com/docker/orca/bootstrap/client"
	butils "github.com/docker/orca/bootstrap/utils"
	"github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

type AuthStore struct {
}

func (p *AuthStore) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	containerJSON, err := dclient.ContainerInspect(context.TODO(), config.AuthStoreContainerName)
	if err != nil {
		log.Debug("Error inspecting %s container: %s", config.AuthStoreContainerName, err)
		return nil
	}

	// Create a container entry for the kv container
	currentCfg.Containers[config.AuthStoreContainerName] = &types.OrcaContainer{
		Image:   containerJSON.Config.Image,
		Running: containerJSON.State.Running,
	}
	return nil
}

func (p *AuthStore) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	containerName := config.AuthStoreContainerName
	ctr, ok := currentCfg.Containers[containerName]
	if !expectedCfg.IsManager {
		if ok && ctr.Running {
			// The current node is not meant to be a manager but the Auth Store is running - demote
			return true, nil
		}
		return false, nil
	}

	// The current node is meant to be a manager, but the Auth Store is not running
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

func (p *AuthStore) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig, chDone chan bool) error {
	reconcile, err := p.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	if !reconcile {
		return nil
	}

	// Clean up existing Auth Store container, if present
	_ = dclient.ContainerRemove(context.TODO(),
		config.AuthStoreContainerName, engineTypes.ContainerRemoveOptions{Force: true})

	ctr, ok := currentCfg.Containers[config.AuthStoreContainerName]
	if !expectedCfg.IsManager && ok && ctr.Running {
		// TODO: demote flow - gracefully spin down auth store
		return nil
	}

	ec, err := engineClient.NewBareClient()
	if err != nil {
		return err
	}

	// Get a list of managers where this node is excluded
	managers := []string{}
	for _, manager := range expectedCfg.Managers {
		if manager != expectedCfg.HostAddress {
			managers = append(managers, manager)
		}
	}

	// TODO: ugly - refactor with KV peers
	peerAddrs := make([]string, len(managers))
	for i, peerIP := range managers {
		peerAddrs[i] = net.JoinHostPort(peerIP, fmt.Sprintf("%d", config.AuthStorePeerPort))
	}

	log.Info("Deploying Auth Store")
	err = butils.DeployAuthStoreContainer(ec, peerAddrs...)
	if err != nil {
		return err
	}

	// In a fresh install, also create the admin user within this lock
	if utils.IsFreshInstall() {
		adminUsername, adminPassword, err := utils.GetUCPCredentials()
		if err != nil {
			return err
		}
		err = ec.RunAuthCreateAdmin(adminUsername, adminPassword)
		if err != nil {
			log.Errorf(`Failed to create initial admin user.  Run "docker logs %s" for more details`, config.AuthCreateAdminContainerName)
			return err
		}
		log.Info("Created Admin user")
	}
	return nil
}
