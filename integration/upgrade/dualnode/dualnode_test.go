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

func (s *TestSuite) TestUpgrade() {
	expected := len(s.WorkerMachines) + len(s.ControllerMachines)
	// First make sure we're a proper dual node rig
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])

	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Equal(s.T(), expected, len(nodes))

	require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))

	controllerClient, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	id, err := utils.GetOrcaID(controllerClient)
	require.Nil(s.T(), err)

	log.Infof("Upgrading controller: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(controllerClient, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))

	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Now upgrade the secondary node
	secondClient, err := s.WorkerMachines[0].GetClient()
	require.Nil(s.T(), err)
	log.Infof("Upgrading secondary: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(secondClient, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))

	// Make sure the swarm is still dual node
	time.Sleep(time.Second) // Might want to wait a little longer
	dualNode := false
	for i := 0; i < 60 && !dualNode; i++ {
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), err)
		if expected == len(nodes) {
			dualNode = true
			break
		}
		time.Sleep(time.Second)
	}

	// TODO Should we try to run a hello world container with constraints to the two nodes just to be paranoid?

	require.True(s.T(), dualNode, "Didn't detect dual node after upgrade")
}

func TestUninstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
