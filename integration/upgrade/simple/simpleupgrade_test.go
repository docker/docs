package doubleinstall

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/samalba/dockerclient"
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

func (s *TestSuite) TestUpgradeWithoutIDFails() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	log.Debug("Attempting upgrade without ID...")
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	require.NotNil(s.T(), utils.RunBootstrapper(client, []string{"upgrade", "-D"}, []string{}, []string{}))
}

func (s *TestSuite) TestUpgrade() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Gather some information about the current install so we can compare after upgrade
	oldFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Infof("Upgrading: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))

	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	log.Info("Verifying upgraded system...")
	// Now verify things have changed.
	newFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)

	require.Equal(s.T(), oldFingerprint, newFingerprint)

	// Now poke at the system a little to make sure it seems healthy
	newUsername := "newadminuser1"
	newPassword := "supersecret"
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	log.Debug("Creating user")
	require.Nil(s.T(), utils.CreateNewUser(nil, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), newUsername, newPassword, true, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	log.Debug("getting docker client as new user")
	client, err = utils.GetUserDockerClient(serverURL, newUsername, newPassword)
	require.Nil(s.T(), err)
	log.Debug("getting version")

	version, err := client.Version()
	require.Nil(s.T(), err)
	log.Info("Get version succeeded")
	require.True(s.T(), strings.Contains(version.Version, "ucp"))

	// Now try to point to Swarm with the same stuff
	swarmURL := *client.URL
	swarmURL.Host = strings.Split(swarmURL.Host, ":")[0] + ":3376" // WARNING - fragile assumption
	swarmClient, err := dockerclient.NewDockerClient(swarmURL.String(), client.TLSConfig)
	require.Nil(s.T(), err)
	version, err = swarmClient.Version()
	require.Nil(s.T(), err)
	require.True(s.T(), strings.Contains(version.Version, "swarm"))
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
