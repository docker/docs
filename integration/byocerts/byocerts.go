package byocerts

// Common upgrade test logic for upgrades, based on acceptance tests, with a few additional upgrade specific tests

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/integration/acceptance"
	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/require"
)

type BYOCertsSuite struct {
	acceptance.AcceptanceSuite
	BaseVersion string // TODO - make this a list so we can do chained upgrades
}

func (s *BYOCertsSuite) PreInstall(m utils.Machine) error {
	if m == s.ControllerMachines[0] {
		log.Debug("Attempting to BYO certs for controller")
		return utils.BYOServerCertInit(m)
	} else {
		// Only do certs for replica controllers
		for _, test := range s.ControllerMachines[1:] {
			if m == test {
				err := utils.BYOServerCertGenSecondary(s.ControllerMachines[0], m)
				if err != nil {
					return err
				}

				// TODO - this really belongs in it's own customizaion routine that can be mixed-and-matched
				// it's not specific to the general BYO cert scenario
				log.Info("Replicating CA material to secondary controller nodes")
				client, err := s.ControllerMachines[0].GetClient()
				if err != nil {
					return err
				}
				id, err := utils.GetOrcaID(client)
				if err != nil {
					return err
				}

				backup, stderrOutput, err := utils.RunBootstrapperWithIO(client, []string{"backup", "--id", id, "--root-ca-only"}, []string{}, []string{}, nil)
				log.Debug(stderrOutput)
				if err != nil {
					return err
				}
				log.Info("Got backup, now copying to the secondary nodes")
				client, err = m.GetClient()
				if err != nil {
					return err
				}
				err = utils.LoadFileInVolume(client, "test_root", "backup.tar", backup)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

func (s *BYOCertsSuite) InstallArgs(m utils.Machine) []string {
	log.Debug("Wiring install up for external-ucp-ca")
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "3376", "--external-ucp-ca", "--san", externalIP}
}

func (s *BYOCertsSuite) JoinReplicaArgs(m utils.Machine, serverURL, fingerprint string) [3][]string {
	args := s.JoinWorkerArgs(m, serverURL, fingerprint)
	args[0] = append(args[0], "--replica")
	return args
}
func (s *BYOCertsSuite) JoinWorkerArgs(m utils.Machine, serverURL, fingerprint string) [3][]string {
	log.Debug("Wiring join up for external-ucp-ca")
	res := s.AcceptanceSuite.JoinWorkerArgs(m, serverURL, fingerprint)
	res[0] = append(res[0], "--external-ucp-ca")
	res[2] = append(res[2], "/var/lib/docker/volumes/test_root/_data/backup.tar:/backup.tar")
	return res
}

func (s *BYOCertsSuite) doTestCADownN(count int) {
	if count > len(s.ControllerMachines) {
		log.Infof("Not enough controllers to shut down %d CAs", count)
		return
	}
	for i := 0; i < count; i++ {
		client, err := s.ControllerMachines[i].GetClient()
		require.Nil(s.T(), err)

		log.Infof("Stopping CA containers on %s", s.ControllerMachines[i].GetName())
		require.Nil(s.T(), client.StopContainer("ucp-client-root-ca", 20))
		defer client.StartContainer("ucp-client-root-ca", nil)
		require.Nil(s.T(), client.StopContainer("ucp-cluster-root-ca", 20))
		defer client.StartContainer("ucp-cluster-root-ca", nil)
		time.Sleep(2 * time.Second)
	}

	if count >= len(s.ControllerMachines) {
		// Total outage. Try to get a bundle from the all controllers and make sure it fails
		for _, serverURL := range s.GetServerURLs() {
			_, err := utils.GetBundle(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
			require.NotNil(s.T(), err, "Still able to get bundle from %s but it should fail", serverURL)
		}
	} else {
		// Try to get a bundle from the all controllers
		nodeURLs := s.GetServerURLs()
		utils.TestAdminHasSwarm(s.T(), nodeURLs)
		utils.TestNonAdminUserNoSwarmForYou(s.T(), nodeURLs)
	}
}

// Try to get a bundle when the primary controller is CA is down
func (s *BYOCertsSuite) TestHACAsDown() {
	for i := 1; i <= len(s.ControllerMachines); i++ {
		s.doTestCADownN(i)
	}
}
