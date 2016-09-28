package client

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"

	"github.com/docker/orca/bootstrap/config"
)

// Find multiple containers by name
func (c *EngineClient) FindContainers(names []string) []*types.ContainerJSON {
	res := make([]*types.ContainerJSON, 0, len(names))
	for _, name := range names {
		log.Debugf("Looking for container %s", name)
		info, err := c.client.InspectContainer(name)
		if err != nil {
			log.Debugf("Container %s not found: %s", name, err)
			continue
		}
		res = append(res, &info)
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

func (c *EngineClient) FindContainerIDsByOrcaInstanceID(orcaInstanceID string) ([]string, error) {
	res := []string{}
	containers, err := c.client.ListContainers(types.ContainerListOptions{All: true, Size: false})
	if err != nil {
		log.Warnf("Failed to list containers %s", err)
		return nil, err
	}
	for _, container := range containers {
		//Skip the bootstrap container so we don't delete ourself while running
		if container.Image == c.bootstrapper.Image {
			log.Debug("Detected bootstrapper, omitting")
			continue
		}
		if id, found := container.Labels[fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)]; found {
			if id == orcaInstanceID {
				log.Debugf("Found UCP:%s container %v (%s)", orcaInstanceID, container.Names, container.ID)
				res = append(res, container.ID)
			}
		}
	}
	return res, nil
}

func (c *EngineClient) FindServiceIDsByOrcaInstanceID(orcaInstanceID string) ([]string, error) {
	res := []string{}
	services, err := c.client.ServiceList(types.ServiceListOptions{})
	if err != nil {
		log.Warnf("Failed to list services %s", err)
		return nil, err
	}
	for _, service := range services {
		if id, found := service.Spec.TaskTemplate.ContainerSpec.Labels[fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)]; found {
			if id == orcaInstanceID {
				log.Debugf("Found UCP:%s service %v (%s)", orcaInstanceID, service.Spec.Name, service.ID)
				res = append(res, service.ID)
			}
		}
	}
	return res, nil
}

func (c *EngineClient) FindSwarmV2Cluster() (string, error) {
	info, err := c.client.Info()
	if err != nil {
		return "", err
	}

	if info.Swarm.LocalNodeState == swarm.LocalNodeStateInactive {
		// The node is not currently participating in a swarm
		return "", nil
	} else if info.Swarm.LocalNodeState == swarm.LocalNodeStatePending {
		return "", fmt.Errorf("the current node is in a pending swarm state")
	} else if info.Swarm.LocalNodeState == swarm.LocalNodeStateError {
		return "", fmt.Errorf("the current node is in an error swarm state")
	}

	node, err := c.client.InspectNode(info.Swarm.NodeID)
	if err != nil {
		return "", err
	}

	switch node.Spec.Role {
	case swarm.NodeRoleManager:
		return "manager", nil
	case swarm.NodeRoleWorker:
		return "worker", nil
	default:
		return "unknown role", nil
	}
}

func (c *EngineClient) FindSwarmV2NodeID() (string, error) {
	info, err := c.client.Info()
	if err != nil {
		return "", err
	}

	return info.Swarm.NodeID, nil
}

func (c *EngineClient) FindIPFromNode(node string) (string, error) {
	info, err := c.client.InspectNode(node)
	if err != nil {
		return "", fmt.Errorf("Could not inspect cluster node %s: %s", node, err)
	}

	hostport := info.ManagerStatus.Addr
	ip, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", fmt.Errorf("Could not parse address from swarm: %s", err)
	}

	return ip, nil
}
