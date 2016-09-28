package recertify_root

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

func (s *TestSuite) TestBasicInstall() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
}

// This does all the heavy lifting
// TODO - in the future, do the recertify in the suite, then we can run the individual scenarios
//        in different tests, instead of one big long test
func (s *TestSuite) TestRecertifyWithCA() {
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	externalIP, err := s.ControllerMachines[0].GetIP()
	require.Nil(s.T(), err)

	log.Info("Gathering existing certs")
	filesToCompare := []string{
		"/var/lib/docker/volumes/ucp-client-root-ca/_data/cert.pem",  // Not regenerated in this test case
		"/var/lib/docker/volumes/ucp-cluster-root-ca/_data/cert.pem", // Not regenerated in this test case
		"/var/lib/docker/discovery_certs/cert.pem",
		"/var/lib/docker/volumes/ucp-controller-client-certs/_data/cert.pem",
		"/var/lib/docker/volumes/ucp-controller-server-certs/_data/cert.pem",
		"/var/lib/docker/volumes/ucp-kv-certs/_data/cert.pem",
		"/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem",
	}
	initialValues := map[string]string{}
	for _, filename := range filesToCompare {
		content, err := s.ControllerMachines[0].CatHostFile(filename)
		require.Nil(s.T(), err)
		initialValues[filename] = string(content)
	}

	// Get a client we can verify fails after regen
	oldAdminClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)

	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Info("Regenerating CA certs on this controller")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"regen-certs", "-D", "--root-ca-only", "--id", id}, []string{}, []string{}))

	log.Info("Now regenerating server certs on this controller")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"regen-certs", "-D", "--id", id, "--san", externalIP}, []string{}, []string{}))

	// NOTE: we could skip bouncing the daemon to speed this up, but it's not a thorough test then...
	log.Info("Restarting the docker daemon for discovery certs to take effect - this'll take a while...")
	require.Nil(s.T(), utils.RestartDockerDaemon(s.ControllerMachines[0]))

	// And orca should be running
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 240))
	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))
	log.Info("orca server is running")

	// Now make sure the certs actually changed
	log.Info("Comparing the new certs to make sure they changed")
	for _, filename := range filesToCompare {
		content, err := s.ControllerMachines[0].CatHostFile(filename)
		require.Nil(s.T(), err)
		require.NotEqual(s.T(), initialValues[filename], string(content), filename)
	}
	log.Info("Certs look like they got regenerated.  Poking at the system to make sure it still looks OK")

	// Make sure the old client (aka, old bundle) fails
	log.Info("Verifying old bundle does not work")
	version, err := oldAdminClient.Version()
	require.NotNil(s.T(), err)

	// Now get a new one and make sure it works
	log.Info("Verifying new bundle works")
	newAdminClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	version, err = newAdminClient.Version()
	require.Nil(s.T(), err)
	require.True(s.T(), strings.Contains(version.Version, "ucp"))

	// Run a few more scenarios to make sure things didn't go wonky
	utils.TestNonAdminUserNoSwarmForYou(s.T(), []string{serverURL})
	utils.TestNonAdminUserNoProxyForYou(s.T(), []string{serverURL})
	utils.TestAdminHasSwarm(s.T(), []string{serverURL})
	utils.TestEscalationsBlocked(s.T(), serverURL)
	engineVersion, err := utils.GetEngineVersion(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestSimpleCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), engineVersion)
	log.Info("Everything looks good on the recertified system")
}

func TestRecertifyWithCA(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
