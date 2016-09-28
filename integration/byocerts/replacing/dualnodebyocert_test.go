package simple

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/integration/utils"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 1
	return
}
func (s *TestSuite) PreInstall(m utils.Machine) error {
	// Only wire up certs for initial controller
	if m == s.ControllerMachines[0] {
		return utils.BYOServerCertInit(m)
	}
	return nil
}

func (s *TestSuite) InstallArgs(m utils.Machine) []string {
	log.Debug("Wiring install up for external-ucp-ca")
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--external-server-cert", "--san", externalIP}
}

func (s *TestSuite) TestReplaceServerCertWithSameCA() {
	// Get the existing CA
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	caPEM, err := utils.GetCAPEM(serverURL)
	require.Nil(s.T(), err)

	// Get a bundle for admin and user
	adminClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	username := "dualnodebyouser2"
	password := "secret"
	require.Nil(s.T(), utils.CreateNewUser(nil, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), username, password, false, auth.View))
	userClient, err := utils.GetUserDockerClient(serverURL, username, password)
	require.Nil(s.T(), err)

	// TODO - consider calling out to the new bootstrapper "stop" routine once available
	log.Info("Stopping ucp-controller")
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	require.Nil(s.T(), client.StopContainer("ucp-controller", 20))

	log.Info("Re-running the cert generation logic")
	// Re-run the cert generation logic
	require.Nil(s.T(), utils.BYOServerCertInitWithSameCA(s.ControllerMachines[0]))

	// Bounce the controller
	log.Info("Starting ucp-controller")
	require.Nil(s.T(), client.StartContainer("ucp-controller", nil))

	log.Info("Waiting for ucp-controller to recover")
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))

	log.Info("Inspecting system to verify new ca/cert picked up")
	// Make sure the CA.pem changed
	newCaPEM, err := utils.GetCAPEM(serverURL)
	require.Nil(s.T(), err)
	require.Equal(s.T(), caPEM, newCaPEM, "ca.pem changed!")

	// Verify both admin and user bundles still work
	log.Info("verifying old bundles still work")
	_, err = adminClient.Version()
	require.Nil(s.T(), err, "Failed to connect with old admin bundle: %s", err)
	_, err = userClient.Version()
	require.Nil(s.T(), err, "Failed to connect with old client bundle: %s", err)

	log.Info("verifying new bundles work")
	// Get new bundles, make sure they work
	adminClient, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	userClient, err = utils.GetUserDockerClient(serverURL, username, password)
	require.Nil(s.T(), err)
	_, err = adminClient.Version()
	assert.Nil(s.T(), err, "Failed to connect with admin bundle: %s", err)
	_, err = userClient.Version()
	assert.Nil(s.T(), err, "Failed to connect with client bundle: %s", err)

	log.Info("Verifying node and controller counts")
	err = utils.VerifyControllerCount(adminClient, 1, 1)
	assert.Nil(s.T(), err)
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(nodes))

	log.Info("System looks good with new certs")
}

func (s *TestSuite) TestReplaceServerCertWithNewCA() {
	// Get the existing CA
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	caPEM, err := utils.GetCAPEM(serverURL)
	require.Nil(s.T(), err)

	// Get a bundle for admin and user
	adminClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	username := "dualnodebyouser"
	password := "secret"
	require.Nil(s.T(), utils.CreateNewUser(nil, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), username, password, false, auth.View))
	userClient, err := utils.GetUserDockerClient(serverURL, username, password)
	require.Nil(s.T(), err)

	// TODO - consider calling out to the new bootstrapper "stop" routine once available
	log.Info("Stopping ucp-controller")
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	require.Nil(s.T(), client.StopContainer("ucp-controller", 20))

	log.Info("Re-running the cert generation logic")
	// Re-run the cert generation logic
	require.Nil(s.T(), utils.BYOServerCertInit(s.ControllerMachines[0]))

	// Bounce the controller
	log.Info("Starting ucp-controller")
	require.Nil(s.T(), client.StartContainer("ucp-controller", nil))

	log.Info("Waiting for ucp-controller to recover")
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))

	log.Info("Inspecting system to verify new ca/cert picked up")
	// Make sure the CA.pem changed
	newCaPEM, err := utils.GetCAPEM(serverURL)
	require.Nil(s.T(), err)
	require.NotEqual(s.T(), caPEM, newCaPEM, "ca.pem didn't change!")

	// Verify both admin and user bundles can't be used again
	log.Info("verifying old bundles no longer work")
	_, err = adminClient.Version()
	require.NotNil(s.T(), err)
	_, err = userClient.Version()
	require.NotNil(s.T(), err)

	log.Info("verifying new bundles work")
	// Get new bundles, make sure they work
	adminClient, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	userClient, err = utils.GetUserDockerClient(serverURL, username, password)
	require.Nil(s.T(), err)
	_, err = adminClient.Version()
	assert.Nil(s.T(), err, "Failed to connect with admin bundle: %s", err)

	_, err = userClient.Version()
	assert.Nil(s.T(), err, "Failed to connect with client bundle: %s", err)

	log.Info("Verifying node and controller counts")
	err = utils.VerifyControllerCount(adminClient, 1, 1)
	assert.Nil(s.T(), err)
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(nodes))

	log.Info("System looks good with new certs")
}

// Regenerate the certs without touching the root CAs
func (s *TestSuite) TestRegenerateCertsWithoutRoot() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
		return
	}
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	externalIP, err := s.ControllerMachines[0].GetIP()
	require.Nil(s.T(), err)

	log.Info("Gathering existing certs")
	filesToCompare := []string{
		//"/var/lib/docker/volumes/ucp-client-root-ca/_data/cert.pem", // Not regenerated in this test case
		//"/var/lib/docker/volumes/ucp-cluster-root-ca/_data/cert.pem", // Not regenerated in this test case
		"/var/lib/docker/discovery_certs/cert.pem",
		"/var/lib/docker/volumes/ucp-controller-client-certs/_data/cert.pem",
		//"/var/lib/docker/volumes/ucp-controller-server-certs/_data/cert.pem", // Omitted since this is a BYO scenario
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

	log.Info("Regenerating certs on this controller")
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"regen-certs", "-D", "--external-server-cert", "--id", id, "--san", externalIP}, []string{}, []string{}))

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

	log.Info("Verifying old bundle does work")
	version, err := oldAdminClient.Version()
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

	log.Info("Verifying node and controller counts")
	err = utils.VerifyControllerCount(oldAdminClient, 1, 1)
	assert.Nil(s.T(), err)
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(nodes))

	log.Info("Everything looks good on the recertified system")
}

func (s *TestSuite) TestVerifyControllerCount() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	err = utils.VerifyControllerCount(client, 1, 1)
	require.Nil(s.T(), err)
}

func (s *TestSuite) TestNodeCount() {
	expected := 2 // TODO Someday driven by the suite node count...
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Equal(s.T(), expected, len(nodes))
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
