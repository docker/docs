package specific_upgrade

import (
	"github.com/docker/orca/integration/upgrade"
	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestSuite struct {
	upgrade.UpgradeSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 3
	workerCount = 0
	return
}

func TestInstallTestSuite(t *testing.T) {
	t.Skip("XXX This is currently non-functional with seattle - re-enable once we have an upgrade HA flow working")
	return

	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.UpgradeSuite.BaseVersion = "1.1.0"
	s.Init(s)
	suite.Run(t, s)
}
