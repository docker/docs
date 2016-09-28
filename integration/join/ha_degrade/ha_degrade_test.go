package simplejoin

import (
	"io/ioutil"
	neturl "net/url"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/engine-api/types/swarm"
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

// TODO Once this works, replicate for controller outage, and ThirdMachine as well to verify no idiosyncracies
func (s *TestSuite) TestDegradedSecondary() {
	// Get a connection to the primary server
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)

	expected := len(s.WorkerMachines) + len(s.ControllerMachines)
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	require.Equal(s.T(), expected, len(nodes))

	// Check for no HA banners
	log.Debugf("Checking for HA banners")
	httpClient := client.HTTPClient
	orcaURL, err := neturl.Parse(serverURL)
	require.Nil(s.T(), err)
	orcaURL.Path = "/api/banner"
	resp, err := httpClient.Get(orcaURL.String())
	require.Nil(s.T(), err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(s.T(), string(body))
	}
	// Crude, but good enough for our test
	if strings.Contains(string(body), "HA Degraded") {
		require.FailNow(s.T(), "HA reported degraded"+string(body))
	}
	log.Debugf("Healthy cluster detected as expected")

	secondURL, err := utils.GetOrcaURL(s.ControllerMachines[1])
	require.Nil(s.T(), err)
	// Shutdown the secondary machine
	log.Info("Shuttind down daemon and looking for degraded state...")
	require.Nil(s.T(), utils.StopDockerDaemon(s.ControllerMachines[1]))

	// It takes a while to detect, so might as well wait a bit before we start checking
	time.Sleep(30 * time.Second)

	// Now query the system and wait for the node count to drop to 2

	// Retry for ~2 minutes to give Swarm a chance to detect the lost node
	degraded := false
	for i := 0; i < 30 && !degraded; i++ {
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), err)
		for _, node := range nodes {
			if node.Status.State != swarm.NodeStateReady {
				degraded = true
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
	require.True(s.T(), degraded)

	// Double check the banner too
	resp, err = httpClient.Get(orcaURL.String())
	require.Nil(s.T(), err)
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(s.T(), string(body))
	}
	// Crude, but good enough for our test
	if !strings.Contains(string(body), "HA Degraded") {
		require.FailNow(s.T(), "HA didn't report degraded"+string(body))
	}
	log.Debugf("Unhealthy cluster detected as expected")

	log.Info("Restarting daemon and looking for recovered state...")
	require.Nil(s.T(), utils.StartDockerDaemon(s.ControllerMachines[1]))
	require.Nil(s.T(), utils.ValidateClusterHealthy(secondURL, 90))
	log.Info("Cluster recovered")
}

func TestHADegradeTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
