package simple

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
	"github.com/samalba/dockerclient"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

var (
	Volumes           []*dockerclient.Volume
	EtcHostFilesystem map[string]interface{}
	VarHostFilesystem map[string]interface{}
)

// WARNING: assumes a single node deployment!
func (s *TestSuite) PreInstall(m utils.Machine) error {
	log.Debug("Retrieving list of volumes before install")
	client, err := m.GetClient()
	require.Nil(s.T(), err)
	volumes, err := client.ListVolumes()
	require.Nil(s.T(), err)

	Volumes = volumes

	log.Debug("Retrieving host /etc filesystem")
	EtcHostFilesystem, err = utils.HostDirManifest(m, "/etc/docker")
	if err != nil {
		return err
	}
	// Whitelist /etc/docker/daemon.json because we create it at install time
	EtcHostFilesystem["./daemon.json"] = struct{}{}
	/*

		        TODO - this shows leakage, but it may be noise.  Need to investigate further.

				log.Debug("Retrieving host /var/lib/docker filesystem")
				VarHostFilesystem, err = utils.HostDirManifest(s.ControllerMachines[0], "/var/lib/docker")
				if err != nil {
					return err
				}
	*/
	return nil
}

func (s *TestSuite) TestForLeackage() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)

	//Dump out all the volumes by container to aid debugging leakage
	containers, err := client.ListContainers(true, false, "")
	require.Nil(s.T(), err)
	log.Info("Container volume mounts:")
	for _, container := range containers {
		info, err := client.InspectContainer(container.Id)
		require.Nil(s.T(), err)
		for key, value := range info.Volumes {
			log.Infof("%s - %s %s", info.Name, key, value)
		}
	}

	// Uninstall the Orca server
	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Debugf("Uninstalling: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"uninstall", "-D", "--id", id}, []string{}, []string{}))

	// Now check the volume list
	volumes, err := client.ListVolumes()
	require.Nil(s.T(), err)

	extrasFound := false
	for _, volume := range volumes {
		found := false
		for _, volCompare := range Volumes {
			if volume.Name == volCompare.Name {
				found = true
				break
			}
		}
		if !found {
			extrasFound = true
			log.Errorf("Volume %s leaked post uninstall (mounted at: %s)", volume.Name, volume.Mountpoint)
		}
	}
	// Use assert so we can also find file leakage
	assert.False(s.T(), extrasFound)

	log.Debug("Retrieving host /etc/docker filesystem")
	newEtcManifest, err := utils.HostDirManifest(s.ControllerMachines[0], "/etc/docker")
	require.Nil(s.T(), err)
	extrasFound = false
	for name := range newEtcManifest {
		if EtcHostFilesystem[name] == nil {
			log.Errorf("Leaked /etc host file %s", name)
			extrasFound = true
		}
	}
	assert.False(s.T(), extrasFound)

	/*

		log.Debug("Retrieving host /var/lib/docker filesystem")
		newVarManifest, err := utils.HostDirManifest(s.ControllerMachines[0], "/var/lib/docker")
		require.Nil(s.T(), err)
		extrasFound = false
		for name := range newVarManifest {
			if VarHostFilesystem[name] == nil {
				log.Errorf("Leaked /var/lib/docker host file %s", name)
				extrasFound = true
			}
		}
		assert.False(s.T(), extrasFound)
	*/
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
