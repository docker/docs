package ha

import (
	"testing"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
)

var (
	// Stress count limits - set a bit low so by default the test doesn't take too long
	UserCount      = utils.GetStressObjectCount(200)
	ImageCount     = utils.GetStressObjectCount(200)
	ContainerCount = utils.GetStressObjectCount(200)
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 3
	workerCount = 0
	return
}

func (s *TestSuite) TestBasicHealth() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])

	expected := 3
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Equal(s.T(), expected, len(nodes))
}

func (s *TestSuite) TestBuildImages() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestBuildImages(s.T(), serverURL, ImageCount)
}

func (s *TestSuite) TestAddUsers() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestAddUsers(s.T(), serverURL, UserCount)
}

func (s *TestSuite) TestCreateContainers() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestCreateContainers(s.T(), serverURL, ContainerCount)
}

func TestUninstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
