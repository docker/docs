package single

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
	controllerCount = 1
	workerCount = 0
	return
}

func (s *TestSuite) TestBackupAndRestore() {
	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(err)
	id, err := utils.GetOrcaID(client)
	require.Nil(err)

	log.Info("Taking full backup of the controller")
	backup, stderrOutput, err := utils.RunBootstrapperWithIO(client, []string{"backup", "-D", "--id", id}, []string{}, []string{}, nil)
	log.Debug(stderrOutput)
	require.Nil(err)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))

	log.Info("Restoring backup")
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", id}, []string{}, []string{}, bytes.NewBuffer([]byte(backup)))
	log.Debug(stdout)
	log.Debug(stderr)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 100))

	// Try to run a container using a client bundle
	require.Nil(utils.WaitForCLIReadiness(serverURL,
		utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
	log.Info("Looks good")
}

func (s *TestSuite) TestBackupAndRestoreWithCrypto() {
	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(err)
	id, err := utils.GetOrcaID(client)
	require.Nil(err)

	secret := "secret1234"

	log.Info("Taking full backup of the controller")
	backup, stderrOutput, err := utils.RunBootstrapperWithIO(client, []string{"backup", "-D", "--id", id, "--passphrase", secret}, []string{}, []string{}, nil)
	log.Debug(stderrOutput)
	require.Nil(err)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))

	log.Info("Restoring backup")
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", id, "--passphrase", secret}, []string{}, []string{}, bytes.NewBuffer([]byte(backup)))
	log.Debug(stdout)
	log.Debug(stderr)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 100))

	require.Nil(utils.WaitForCLIReadiness(serverURL,
		utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
}

func (s *TestSuite) TestBackupAndRestoreFreshInstall() {
	s.T().Skip("Fresh install is not currently supported so this test case doesn't work.")
	return

	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(err)
	id, err := utils.GetOrcaID(client)
	require.Nil(err)
	machineIP, err := s.ControllerMachines[0].GetIP()
	require.Nil(err)

	secret := "secret1234"

	log.Info("Taking full backup of the controller")
	backup, stderrOutput, err := utils.RunBootstrapperWithIO(client, []string{"backup", "-D", "--id", id, "--passphrase", secret}, []string{}, []string{}, nil)
	log.Debug(stderrOutput)
	require.Nil(err)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))

	log.Info("Performing a fresh install of the controller")
	require.Nil(utils.RunBootstrapper(client, []string{
		"install", "--swarm-port", "3376", "--fresh-install", "-D",
		"--disable-tracking", "--disable-usage", "--san", machineIP,
	}, []string{}, []string{}))

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))
	require.Nil(utils.WaitForCLIReadiness(serverURL,
		utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))

	log.Info("Restoring backup")
	stdout, stderr, err := utils.RunBootstrapperWithIO(client, []string{"restore", "-D", "--id", id, "--passphrase", secret}, []string{}, []string{}, bytes.NewBuffer([]byte(backup)))
	log.Debug(stdout)
	log.Debug(stderr)

	log.Info("Verifying the system recovers")
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 100))
	require.Nil(utils.ValidateClusterHealthy(serverURL, 100))
	require.Nil(utils.WaitForCLIReadiness(serverURL,
		utils.GetAdminUser(), utils.GetAdminPassword(), time.Minute))
}

func TestBackupRestore(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
