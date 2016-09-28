package ha

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
	controllerCount = 2
	workerCount = 0
	return
}

const KVTimeout = "1500"

func (s *TestSuite) InstallArgs(m utils.Machine) []string {
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--kv-timeout", KVTimeout, "--san", externalIP}
}

func (s *TestSuite) TestTimeoutIsSet() {
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)

	timeout, err := utils.GetContainerFlagValue(client, "ucp-kv", "--election-timeout")
	require.Nil(s.T(), err)
	require.Equal(s.T(), timeout, KVTimeout)

	// check secondary machine
	client, err = s.ControllerMachines[1].GetClient()
	require.Nil(s.T(), err)

	timeout, err = utils.GetContainerFlagValue(client, "ucp-kv", "--election-timeout")
	require.Nil(s.T(), err)
	require.Equal(s.T(), timeout, KVTimeout)
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
