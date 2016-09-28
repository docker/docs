package dns

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
)

var (
	DNS1    = "4.2.2.1"
	DNS2    = "4.2.2.2"
	Opt1    = "timeout:10"
	Opt2    = "rotate"
	Search1 = "acme.com"
	Search2 = "foo.com"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

func (s *TestSuite) InstallArgs(m utils.Machine) []string {
	log.Debug("Wiring install up for non-standard ports")
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376",
		"--dns", DNS1,
		"--dns", DNS2,
		"--dns-opt", Opt1,
		"--dns-opt", Opt2,
		"--dns-search", Search1,
		"--dns-search", Search2,
		"--san", externalIP,
	}
}

func (s *TestSuite) TestBasicInstall() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	// Inspect a few containers and make sure their DNS settings look right
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)

	info, err := client.InspectContainer("ucp-controller")
	require.Nil(s.T(), err)

	require.Equal(s.T(), len(info.HostConfig.Dns), 2)
	require.Equal(s.T(), len(info.HostConfig.DnsSearch), 2)
	require.Equal(s.T(), len(info.HostConfig.DNSOptions), 2)

	require.True(s.T(), info.HostConfig.Dns[0] == DNS1 || info.HostConfig.Dns[1] == DNS1)
	require.True(s.T(), info.HostConfig.Dns[0] == DNS2 || info.HostConfig.Dns[1] == DNS2)

	require.True(s.T(), info.HostConfig.DNSOptions[0] == Opt1 || info.HostConfig.DNSOptions[1] == Opt1)
	require.True(s.T(), info.HostConfig.DNSOptions[0] == Opt2 || info.HostConfig.DNSOptions[1] == Opt2)

	require.True(s.T(), info.HostConfig.DnsSearch[0] == Search1 || info.HostConfig.DnsSearch[1] == Search1)
	require.True(s.T(), info.HostConfig.DnsSearch[0] == Search2 || info.HostConfig.DnsSearch[1] == Search2)
}

func TestInstallDnsTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
