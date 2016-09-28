package simplejoin

import (
	"testing"
	"time"

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
	workerCount = 1
	return
}

func (s *TestSuite) TestVerifyControllerCount() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	err = utils.VerifyControllerCount(client, 1, 1) // Only expect 1 since this isn't an HA scenario
	require.Nil(s.T(), err)
}
func (s *TestSuite) TestSimpleDockerCLI() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	engineVersion, err := utils.GetEngineVersion(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestSimpleCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), engineVersion)
}

// Make sure we can join twice without failure
func (s *TestSuite) TestRejoin() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	fingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)

	// Record the certs so we can verify they change
	discoCert, err := s.WorkerMachines[0].CatHostFile("/var/lib/docker/discovery_certs/cert.pem")
	require.Nil(s.T(), err)
	nodeCert, err := s.WorkerMachines[0].CatHostFile("/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem")
	require.Nil(s.T(), err)

	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))
	secondClient, err := s.WorkerMachines[0].GetClient()
	require.Nil(s.T(), err)

	log.Info("Attempting to re-join the existing node")
	require.Nil(s.T(), utils.RunBootstrapper(secondClient,
		[]string{"join", "-D", "--fingerprint", fingerprint, "--fresh-install", "--url", serverURL},
		[]string{"UCP_ADMIN_USER=" + utils.GetAdminUser(), "UCP_ADMIN_PASSWORD=" + utils.GetAdminPassword()}, []string{}))
	require.Nil(s.T(), err)
	engineVersion, err := utils.GetEngineVersion(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestSimpleCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), engineVersion)

	log.Info("Verifying the certs have changed")
	newDiscoCert, err := s.WorkerMachines[0].CatHostFile("/var/lib/docker/discovery_certs/cert.pem")
	require.Nil(s.T(), err)
	newNodeCert, err := s.WorkerMachines[0].CatHostFile("/var/lib/docker/volumes/ucp-node-certs/_data/cert.pem")
	require.Nil(s.T(), err)
	require.NotEqual(s.T(), string(discoCert), string(newDiscoCert))
	require.NotEqual(s.T(), string(nodeCert), string(newNodeCert))

	//  Make sure we still have 2 nodes
	// Retry for ~1 minute to give Swarm a chance to detect the node
	for i := 0; i < 60; i++ {
		// verify the system has two nodes
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), err)
		if len(nodes) == 2 {
			log.Info("Dual node (still) detected: PASS")
			return
		}
		time.Sleep(time.Second)
	}
	require.FailNow(s.T(), "Didn't detect dual node")
}

func (s *TestSuite) TestLDAPBasic() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestLDAPBasic(s.T(), serverURL)
}

func (s *TestSuite) TestNodeCount() {
	expected := len(s.WorkerMachines) + len(s.ControllerMachines)
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Equal(s.T(), expected, len(nodes))
}

func TestUninstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
