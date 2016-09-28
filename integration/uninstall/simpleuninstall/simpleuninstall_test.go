package simpleuninstall

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

func (s *TestSuite) TestBasicUninstall() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)

	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Debugf("Uninstalling: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"uninstall", "-D", "--id", id}, []string{}, []string{}))

	log.Debug("Verifying files were removed...")
	for _, filename := range []string{
		// We'll spot check a few specific files that should have been removed
		"/etc/docker/ssl/orca/orca_ca.pem",
		"/var/lib/docker/orca/orca/trust.key",
		"/var/lib/docker/orca/consul/raft/peers.json",
	} {
		_, err := s.ControllerMachines[0].CatHostFile(filename)
		require.Equal(s.T(), err, utils.PathDoesNotExist, filename)
	}

	log.Debug("Verifying images were removed...")
	c, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	for _, imageName := range []string{
		"docker/ucp-controller",
	} {
		_, err := c.InspectImage(imageName)
		// XXX This seems to have changed... there might be an underlying dockerclient bug...
		// We used to see "Image not found" and now we see "Not Found" - so handle both
		require.True(s.T(), strings.Contains(err.Error(), "ot found"))
	}

}

func TestUninstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
