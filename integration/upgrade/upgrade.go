package upgrade

// Common upgrade test logic for upgrades, based on acceptance tests, with a few additional upgrade specific tests

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/integration/acceptance"
	"github.com/docker/orca/integration/utils"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

type UpgradeSuite struct {
	acceptance.AcceptanceSuite
	BaseVersion string // TODO - make this a list so we can do chained upgrades
}

func (s *UpgradeSuite) TestUpgradeUsersAndTeam() {
	serverAddr, err := utils.GetOrcaAddr(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestUpgradeUsersAndTeam(s.T(), serverAddr)
}

func (s *UpgradeSuite) PreInstall(m utils.Machine) error {

	s.IsDiscoveryMissing = true // TODO - detect the old version and set this conditionally

	// Temporarily replace the bootstrapper image with the official build
	// of the base version...
	utils.BootstrapImage = "docker/ucp:" + s.BaseVersion
	log.Debugf("Setting bootstrap image to %s", utils.BootstrapImage)

	// Pull the desired images on all the controllers/workers
	log.Debugf("Pulling %s to %s", utils.BootstrapImage, m.GetName())
	client, err := m.GetClient()
	if err != nil {
		return err
	}
	err = client.PullImage(utils.BootstrapImage, &dockerclient.AuthConfig{
		Username: os.Getenv("REGISTRY_USERNAME"),
		Password: os.Getenv("REGISTRY_PASSWORD"),
		Email:    os.Getenv("REGISTRY_EMAIL"),
	})
	if err != nil {
		log.Errorf("Failed to pull %s: %s", utils.BootstrapImage, err)
		return err
	}
	return nil
}

func (s *UpgradeSuite) InstallArgs(m utils.Machine) []string {
	log.Debugf("Wiring install up for version %s", s.BaseVersion)
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "-D", "--swarm-port", "3376", "--san", externalIP, "--image-version", s.BaseVersion}
}

func (s *UpgradeSuite) JoinReplicaArgs(m utils.Machine, serverURL, fingerprint string) [3][]string {
	args := s.JoinWorkerArgs(m, serverURL, fingerprint)
	args[0] = append(args[0], "--replica")
	return args
}
func (s *UpgradeSuite) JoinWorkerArgs(m utils.Machine, serverURL, fingerprint string) [3][]string {
	res := s.AcceptanceSuite.JoinWorkerArgs(m, serverURL, fingerprint)
	res[0] = append(res[0], "--image-version", s.BaseVersion)
	return res
}

func (s *UpgradeSuite) InitialConfig() error {
	utils.BootstrapImage = utils.BuildImageString("ucp")
	log.Infof("Upgrading from %s to %s", s.BaseVersion, utils.BootstrapImage)

	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	// Validate we're a healthy cluster before proceeding
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))
	utils.SetupUpgradeUsersAndTeam(s.T(), serverURL) // Set up upgrade teams
	adminClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Nil(s.T(), utils.VerifyControllerCount(adminClient, len(s.ControllerMachines), 30))

	for _, m := range append(s.ControllerMachines, s.WorkerMachines...) {

		// Now upgrade all the nodes
		log.Infof("Upgrading node: %s", m.GetName())
		client, err := m.GetClient()
		require.Nil(s.T(), err)
		require.Nil(s.T(), utils.RunBootstrapper(client, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))
		require.Nil(s.T(), utils.VerifyControllerCount(adminClient, len(s.ControllerMachines), 30))
	}
	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))

	log.Info("System looks good, proceeding with tests")
	return nil
}
