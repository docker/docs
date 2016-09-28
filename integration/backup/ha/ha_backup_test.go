package ha

import (
	"bytes"
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

func (s *TestSuite) InitialConfig() error {
	// TODO if we want to support upgrade flows, then we'd call the upgrade InitialConfig here
	s.TakeBackup(s.ControllerMachines[0])
	return nil
}

// TODO
// Should we refactor this as two distinct suites derived from acceptance so we can do the backup/restore
// during setup, then run a full acceptance suite to verify nothing is broken after the backup/restore?

// Run the restore on the controller that was backed up and rejoin the other controllers
func (s *TestSuite) TestRestoreControllerAndRejoin() {
	backup := bytes.NewBuffer(s.Data[s.ControllerMachines[0].GetName()+":backup"])
	require := require.New(s.T())
	secondClient, err := s.ControllerMachines[1].GetClient()
	require.Nil(err)
	require.Nil(utils.RunBootstrapper(secondClient, []string{"stop"}, []string{}, []string{}))
	thirdClient, err := s.ControllerMachines[2].GetClient()
	require.Nil(err)
	require.Nil(utils.RunBootstrapper(thirdClient, []string{"stop"}, []string{}, []string{}))

	m := s.ControllerMachines[0]
	serverURL, err := utils.GetOrcaURL(m)
	log.Infof("Restoring on %s", m.GetName())
	client, err := m.GetClient()
	require.Nil(err)
	ClusterID, err := utils.GetOrcaID(client)
	require.Nil(err)
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", ClusterID}, []string{}, []string{}, backup)
	log.Debug(stdout)
	log.Debug(stderr)
	require.Nil(err)

	require.Nil(utils.ValidateOrcaServerRunning(m, 60))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 60))
	require.Nil(utils.WaitForCLIReadiness(serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
	log.Infof("Restoring complete and cluster recovered")

	fingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])

	args := []string{
		"join", "-D", "--fresh-install", "--replica", "--fingerprint", fingerprint,
		"--url", serverURL, "--swarm-port", "3376",
	}
	sip, err := s.ControllerMachines[1].GetIP()
	require.Nil(err)
	secargs := append(args, []string{"--san", sip}...)

	tip, err := s.ControllerMachines[2].GetIP()
	require.Nil(err)
	thirdargs := append(args, []string{"--san", tip}...)

	env := []string{"UCP_ADMIN_USER=" + utils.GetAdminUser(), "UCP_ADMIN_PASSWORD=" + utils.GetAdminPassword()}
	misc := []string{}

	log.Debug("re-joining second node")
	require.Nil(utils.RunBootstrapper(secondClient, secargs, env, misc))

	log.Debug("Waiting for cluster to be ready again")
	require.Nil(utils.ValidateOrcaServerRunning(m, 60))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 60))
	require.Nil(utils.WaitForCLIReadiness(serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[1], 60))

	log.Debug("re-joining third node")
	require.Nil(utils.RunBootstrapper(thirdClient, thirdargs, env, misc))

	log.Debug("Waiting for cluster to be ready again")
	require.Nil(utils.ValidateOrcaServerRunning(m, 100))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 100))
	require.Nil(utils.WaitForCLIReadiness(serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[1], 100))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[2], 100))
}

// Run restore on a controller while the other controllers of the cluster are still running.
// The restore operation should fail in this scenario
func (s *TestSuite) TestRestoreControllerWhileOthersStillRunning() {
	backup := bytes.NewBuffer(s.Data[s.ControllerMachines[0].GetName()+":backup"])
	require := require.New(s.T())
	m := s.ControllerMachines[0]
	serverURL, err := utils.GetOrcaURL(m)
	log.Infof("Restoring on %s", m.GetName())
	client, err := m.GetClient()
	require.Nil(err)
	ClusterID, err := utils.GetOrcaID(client)
	require.Nil(err)
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", ClusterID}, []string{}, []string{}, backup)
	if err != nil {
		log.Debug(err.Error())
	}
	log.Debug(stdout)
	log.Debug(stderr)
	require.NotNil(err)

	require.Nil(utils.ValidateOrcaServerRunning(m, 60))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[1], 60))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[2], 60))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 60))
	require.Nil(utils.WaitForCLIReadiness(serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
}

// Run restore on a controller that was not backed up.
// The restore operation should fail in this scenario
func (s *TestSuite) TestRestoreWrongController() {
	backup := bytes.NewBuffer(s.Data[s.ControllerMachines[0].GetName()+":backup"])
	require := require.New(s.T())

	// Stop the other controllers
	secondClient, err := s.ControllerMachines[1].GetClient()
	require.Nil(err)
	require.Nil(utils.RunBootstrapper(secondClient, []string{"stop"}, []string{}, []string{}))
	thirdClient, err := s.ControllerMachines[2].GetClient()
	require.Nil(err)
	require.Nil(utils.RunBootstrapper(thirdClient, []string{"stop"}, []string{}, []string{}))

	m := s.ControllerMachines[1]
	serverURL, err := utils.GetOrcaURL(m)
	log.Infof("Restoring on %s", m.GetName())
	client, err := m.GetClient()
	require.Nil(err)
	ClusterID, err := utils.GetOrcaID(client)
	require.Nil(err)
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", ClusterID}, []string{}, []string{}, backup)
	log.Debug(stdout)
	log.Debug(stderr)
	if err != nil {
		log.Debug(err.Error())
	}
	require.NotNil(err)

	// Restart the other controllers
	require.Nil(utils.RunBootstrapper(secondClient, []string{"restart"}, []string{}, []string{}))
	require.Nil(utils.RunBootstrapper(thirdClient, []string{"restart"}, []string{}, []string{}))

	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[1], 100))
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[2], 100))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 100))
	require.Nil(utils.WaitForCLIReadiness(serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
}

func TestHABackupRestoreSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
