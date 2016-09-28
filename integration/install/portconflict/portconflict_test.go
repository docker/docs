package portconflict

import (
	"testing"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
)

type TestSuite struct {
	utils.TestSuite
}

func (s *TestSuite) TestBasicPortConflictInstall() {
	client, err := s.ControllerMachine.GetClient()
	require.Nil(s.T(), err)
	// Relying on the fact that machine claims the default port
	require.NotNil(s.T(), utils.RunBootstrapper(client, []string{"install", "-D"}, []string{}, []string{}))

	// The Orca server shouldn't be running
	require.NotNil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachine, 1))
}

func TestPortConflictInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	suite.Run(t, new(TestSuite))
}
