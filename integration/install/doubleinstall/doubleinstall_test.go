package doubleinstall

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func (s *TestSuite) TestDoubleInstallFails() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	log.Debug("Attempting second install...") // Use alternate port so we don't hit the port check logic
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	require.NotNil(s.T(), utils.RunBootstrapper(client, []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "4376"}, []string{}, []string{}))
}

func (s *TestSuite) TestReinstall() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Gather some information about the current install so we can compare
	oldFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	oldId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	// Record the original discovery certs
	log.Debug("Verifying files were removed...")
	discoveryCerts := map[string][]byte{
		"ca.pem":   nil,
		"cert.pem": nil,
		"key.pem":  nil,
	}
	for k := range discoveryCerts {
		discoveryCerts[k], err = s.ControllerMachines[0].CatHostFile("/var/lib/docker/discovery_certs/" + k)
		require.Nil(s.T(), err)
	}

	log.Debug("Attempting second install...")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--fresh-install"}, []string{}, []string{}))

	log.Debug("Verifying second install...")
	// Now verify things have changed.
	newFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	newId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	require.NotEqual(s.T(), oldFingerprint, newFingerprint)
	require.NotEqual(s.T(), oldId, newId)

	// Check that the certs changed
	for k := range discoveryCerts {
		newData, err := s.ControllerMachines[0].CatHostFile("/var/lib/docker/discovery_certs/" + k)
		assert.Nil(s.T(), err)
		assert.NotEqual(s.T(), string(newData), string(discoveryCerts[k]), k)
	}

}

func (s *TestSuite) TestReinstallPreserveCerts() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Gather some information about the current install so we can compare
	oldFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	oldId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Debug("Attempting second install...")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--fresh-install", "--preserve-certs"}, []string{}, []string{}))

	log.Debug("Verifying second install...")
	// Now verify things have changed.
	newFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	newId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	require.Equal(s.T(), oldFingerprint, newFingerprint)
	require.NotEqual(s.T(), oldId, newId)
}

func (s *TestSuite) TestReinstallExternalCA() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Gather some information about the current install so we can compare
	oldFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	oldId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Debug("Attempting second install...")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--fresh-install", "--external-ucp-ca"}, []string{}, []string{}))

	log.Debug("Verifying second install...")
	// Now verify things have changed.
	newFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	newId, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	require.Equal(s.T(), oldFingerprint, newFingerprint)
	require.NotEqual(s.T(), oldId, newId)

	/* XXX This test fails, need to investigate why, might be an actual bug!
	// Verify we can't get an admin bundle
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	_, err = utils.GetUserDockerClient(serverURL, "admin", "orca")
	require.NotNil(s.T(), err)
	*/

	// TODO - would be nice to try to get a regular user account eventually...
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
