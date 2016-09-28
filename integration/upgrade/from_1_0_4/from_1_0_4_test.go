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

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.UpgradeSuite.BaseVersion = "1.0.4"
	s.Init(s)
	suite.Run(t, s)
}
