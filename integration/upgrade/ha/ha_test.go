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
	controllerCount = 3
	workerCount = 0
	return
}

func (s *TestSuite) TestBasicAddHost() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])

	// Retry for ~1 minute to give Swarm a chance to detect the node
	triNode := false
	for i := 0; i < 60 && !triNode; i++ {
		// verify the system has three nodes
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), err)
		if len(nodes) == 3 {
			log.Debug("Three node detected: PASS")
			triNode = true
			break
		}
		time.Sleep(time.Second)
	}
	require.True(s.T(), triNode, "Didn't detect three nodes")

	// Now upgrade the nodes

	firstClient, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	id, err := utils.GetOrcaID(firstClient)
	require.Nil(s.T(), err)
	log.Infof("Upgrading first node: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(firstClient, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Now upgrade the second node
	secondClient, err := s.ControllerMachines[1].GetClient()
	require.Nil(s.T(), err)
	log.Infof("Upgrading second node: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(secondClient, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Now upgrade the second node
	thirdClient, err := s.ControllerMachines[2].GetClient()
	require.Nil(s.T(), err)
	log.Infof("Upgrading third node: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(thirdClient, []string{"upgrade", "-D", "--pull", "never", "--id", id}, []string{}, []string{}))
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Retry for ~1 minute to give Swarm a chance to detect the node
	triNode = false
	for i := 0; i < 60 && !triNode; i++ {
		// verify the system has three nodes
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), err)
		if len(nodes) == 3 {
			log.Debug("Three node detected: PASS")
			triNode = true
			break
		}
		time.Sleep(time.Second)
	}
	require.True(s.T(), triNode, "Didn't detect three nodes")
}

func TestUninstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
