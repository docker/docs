package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/orca/bootstrap/config"
)

// Prepare all the volumes we need
func (c *EngineClient) PrepareAllVolumes() error {
	for _, name := range []string{
		config.OrcaRootCAVolumeName,
		config.SwarmRootCAVolumeName,
		config.OrcaServerCertVolumeName,
		config.SwarmNodeCertVolumeName,
	} {
		if err := c.PrepareVolume(name); err != nil {
			return err
		}
	}
	return nil
}

// Remove all the data volumes
func (c *EngineClient) RemoveVolume(name string) error {
	if err := c.client.RemoveVolume(name); err != nil {
		log.Warnf("Failed to remove volume %s: %s", name, err)
		return err
	}
	return nil
}

// Create a volume if it doesn't already exist
func (c *EngineClient) PrepareVolume(name string) error {
	if c.VolumeExists(name) {
		log.Debugf("Re-using existing volume %s", name)
		return nil
	}
	log.Debugf("Creating volume %s", name)
	_, err := c.client.CreateVolume(types.VolumeCreateRequest{
		Name: name,
		// Tunables omitted - users can pre-create their own volumes if they have special requirements
		// Driver:     "",
		// DriverOpts: "",
	})
	return err
}

func (c *EngineClient) VolumeExists(name string) bool {
	_, err := c.client.VolumeInspect(name)
	if err != nil {
		return false
	}
	return true
}

// Return a list of the Orca volumes that already exist on this engine
func (c *EngineClient) ListExistingOrcaVolumes() ([]string, error) {
	volumes, err := c.client.VolumeList()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, volume := range volumes {
		if config.AllVolumesMap[volume.Name] != nil {
			ret = append(ret, volume.Name)
		}
	}
	return ret, nil
}
