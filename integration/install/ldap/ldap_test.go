package ldap

import (
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

func (s *TestSuite) PostInstall(m utils.Machine) error {
	serverURL, err := utils.GetOrcaURL(m)
	require.Nil(s.T(), err)
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	utils.SetupLDAPAuth(s.T(), client)
	// At this point the builtin admin user is unusable
	client, err = utils.GetUserDockerClient(serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	require.Nil(s.T(), err)
	return nil
}

func (s *TestSuite) TestBasicInstall() {
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
}

func (s *TestSuite) TestLDAPBasicClientConnect() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestLDAPBasicClientConnect(s.T(), serverURL)
}

func (s *TestSuite) TestComposeWithLDAPAdmin() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestSimpleCompose(s.T(), serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
}

func (s *TestSuite) TestStillWorksAfterRestart() {
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)

	// Stop ucp-controller and controller, restart controller, wait a little while, then start kv
	log.Debug("Stopping UCP")
	require.Nil(s.T(), client.StopContainer("ucp-controller", 20))
	time.Sleep(2 * time.Second)
	log.Debug("Starting UCP")
	require.Nil(s.T(), client.StartContainer("ucp-controller", nil))
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err = utils.GetUserDockerClient(serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	require.Nil(s.T(), err)
	version, err := client.Version()
	require.Nil(s.T(), err)
	require.True(s.T(), strings.Contains(version.Version, "ucp"))
}

// This relies on the fact that the LDAP default role is restricted
func (s *TestSuite) TestComposeWithLDAPUser() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestSimpleCompose(s.T(), serverURL, os.Getenv("LDAP_USER"), os.Getenv("LDAP_USER_PASSWORD"))
}

// Make sure that an explicit sync doesn't cause the role to change
func (s *TestSuite) TestUserRoleDoesntChangeOnSync() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestUserRoleDoesntChangeOnSync(s.T(), serverURL)
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	if !utils.LDAPTestsEnabled(t) {
		return
	}
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
