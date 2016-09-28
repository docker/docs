package simplejoin

import (
	"strings"
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

func (s *TestSuite) TestUninstallNode() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachine)
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)

	/* TODO - shift to this
	expected := 3
	nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)
	require.Equal(expected, len(nodes))
	*/
	// Verify we're a 3 node cluster first
	// Retry for ~1 minute to give Swarm a chance to detect the node
	threeNode := false
	for i := 0; i < 60 && !threeNode; i++ {
		// verify the system has two nodes
		info, err := client.Info()
		require.Nil(s.T(), err)
		log.Debugf("Info Driver Status: %v", info.DriverStatus)
		for _, driver := range info.DriverStatus {
			if strings.Contains(driver[0], "Nodes") {
				log.Debugf("Nodes line: %v", driver)
				if strings.TrimSpace(driver[1]) == "3" {
					log.Debug("Three node detected")
					threeNode = true
					break
				}
			}
		}
		time.Sleep(time.Second)
	}
	require.True(s.T(), threeNode)

	// Now proceed to uninstall one of the secondary nodes
	client, err = s.SecondMachine.GetClient()
	require.Nil(s.T(), err)

	id, err := utils.GetOrcaID(client)
	require.Nil(s.T(), err)

	log.Infof("Uninstalling secondary node: %s", id)
	require.Nil(s.T(), utils.RunBootstrapper(client, []string{"uninstall", "-D", "--id", id}, []string{}, []string{}))

	// And finally check cluster has 2 nodes now
	client, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	err = utils.VerifyControllerCount(client, 2, 120)
	require.Nil(s.T(), err)
}

func TestUninstallTestSuite(t *testing.T) {
	t.Skip("XXX This is currently non-functional with seattle - re-enable once we have an uninstall flow working")
	return
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
