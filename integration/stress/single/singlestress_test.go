package simple

import (
	"testing"

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

var (
	// Stress count limits - set a bit low so by default the test doesn't take too long
	UserCount      = utils.GetStressObjectCount(200)
	ImageCount     = utils.GetStressObjectCount(200)
	ContainerCount = utils.GetStressObjectCount(200)
)

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

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
