package client

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// Automatically find the instance ID based on containers.
// XXX: needs work to support multiple UCPs on a node
func (c *EngineClient) FindOrcaInstanceID() (string, error) {
	containers := c.FindContainers(orcaconfig.RuntimeContainerNames)

	ids := GetInstanceIDs(containers)
	if len(ids) == 0 {
		return "", fmt.Errorf("No running UCP instances detected on this engine")
	}
	if len(ids) > 1 {
		log.Warnf("Multiple UCP instances detected: %v", ids)
	}
	instanceID := ids[0]
	return instanceID, nil
}

// Stop and Remove the specified containers, plus any related orca containers by instance ID
func (c *EngineClient) RemoveOrcaContainers(containers []*types.ContainerJSON) error {
	// Find the unique orca instance ID(s) and then query for any related containers
	containerIDSet := make(map[string]struct{})
	orcaIDSet := make(map[string]struct{})
	for _, info := range containers {
		containerIDSet[info.ID] = struct{}{}
		if info.Config == nil {
			continue
		}
		if id, found := info.Config.Labels[fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)]; found {
			orcaIDSet[id] = struct{}{}
		}
	}
	for orcaID := range orcaIDSet {
		containerIDs, err := c.FindContainerIDsByOrcaInstanceID(orcaID)
		if err != nil {
			return err
		}
		for _, id := range containerIDs {
			containerIDSet[id] = struct{}{}
		}
	}

	// Now stop/remove all of the IDs we found
	var ret error
	for id := range containerIDSet {
		log.Debugf("Stopping container %s", id)
		if err := c.client.StopContainer(id, 5); err != nil {
			ret = err
		}
	}
	for id := range containerIDSet {
		log.Debugf("Removing container %s", id)
		var err error
		// Sometimes we'll get a "busy" from the engine, so try a few times before giving up
		for i := 0; i < 5; i++ {
			if err = c.client.RemoveContainer(id, true, true); err != nil {
				log.Debugf("daemon reported: %s", err)
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if err != nil {
			ret = err
		}
	}
	return ret
}

// Stop and Remove the specified containers
func (c *EngineClient) RemoveContainers(containers []*types.ContainerJSON) error {
	var ret error
	for _, info := range containers {
		log.Debugf("Stopping and removing container %s", info.Name)
		if err := c.RemoveContainerByID(info.ID, 5, true); err != nil {
			ret = err
			continue
		}
	}
	return ret
}

// Stop and Remove the specified container
func (c *EngineClient) RemoveContainerByID(containerID string, timeout int, volumes bool) error {
	log.Debugf("Removing container %s", containerID)
	// Ignore failed stop, just proceed with attempt to remove regardless
	if err := c.client.StopContainer(containerID, timeout); err != nil {
		log.Infof("Unable to cleanly stop %s, proceeding with removal", containerID)
	}

	if err := c.client.RemoveContainer(containerID, true, volumes); err != nil {
		return err
	}
	return nil
}

func GetInstanceIDs(containers []*types.ContainerJSON) []string {
	containerIDSet := make(map[string]struct{})
	for _, container := range containers {
		if id, found := container.Config.Labels[fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)]; found {
			log.Debugf(`Container %s has instance ID "%s"`, container.Name, id)
			id = strings.TrimSpace(id)
			if id != "" {
				containerIDSet[id] = struct{}{}
			}
		}
	}
	res := []string{}
	for id := range containerIDSet {
		res = append(res, id)
	}
	return res
}

func (c *EngineClient) StartContainer(id string) error {
	// TODO - might want to do some lookups for logging by name instead of ID
	return c.client.StartContainer(id)
}
func (c *EngineClient) StopContainer(id string) error {
	// TODO - might want to do some lookups for logging by name instead of ID
	return c.client.StopContainer(id, 30)
}
func (c *EngineClient) ContainerRestart(id string) error {
	// TODO - might want to do some lookups for logging by name instead of ID
	return c.client.ContainerRestart(id, 30)
}

func (c *EngineClient) StopContainers(containerIDs []string) error {
	for _, id := range containerIDs {
		log.Debugf("Stopping container %s", id)
		if err := c.StopContainer(id); err != nil {
			// Fail fast when stopping
			return err
		}
	}
	return nil
}
func (c *EngineClient) RestartContainers(containerIDs []string) error {
	// Attempt to resume as many as possible
	var lastError error
	for _, id := range containerIDs {
		log.Debugf("Restarting container %s", id)
		if err := c.ContainerRestart(id); err != nil {
			// TODO - might want to give them the friendly name not id here, but inspect could fail too...
			lastError = fmt.Errorf("Failed to restart container %s: %s", id, err)
		}
	}
	return lastError
}

func (ec *EngineClient) CreateContainer(config *container.Config, hostConfig *container.HostConfig, name string) (types.ContainerCreateResponse, error) {
	return ec.client.CreateContainer(config, hostConfig, name)
}

func (ec *EngineClient) RemoveContainer(id string, force bool, volumes bool) error {
	return ec.client.RemoveContainer(id, force, volumes)
}
