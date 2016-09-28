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

func (s *TestSuite) InitialConfig() error {
	// Do the LDAP config first
	err := utils.OrcaSuiteEnableLDAPLegacy(&s.UpgradeSuite.OrcaTestSuite)
	if err != nil {
		return err
	}

	// Then do the upgrade initial config, which does the upgrade
	return s.UpgradeSuite.InitialConfig()
}

func TestInstallTestSuite(t *testing.T) {
	if !utils.LDAPTestsEnabled(t) {
		// Skip the whole suite if no env vars are present
		return
	}
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.UpgradeSuite.BaseVersion = "1.0.4"
	s.Init(s)
	suite.Run(t, s)
}
