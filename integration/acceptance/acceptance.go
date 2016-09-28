package acceptance

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/integration/utils"
	"github.com/docker/orca/integration/utils/ui"
)

type AcceptanceSuite struct {
	utils.OrcaTestSuite
}

func (s *AcceptanceSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

func (s *AcceptanceSuite) TestVerifyControllerCount() {
	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)
	err = utils.VerifyControllerCount(client, len(s.ControllerMachines), 1)
	require.Nil(err)
}

func (s *AcceptanceSuite) TestFingerprintTool() {
	// TODO: retest on AWS
	s.T().Skip("XXX This fingerprint test is broken on AWS - disabling until I can figure out why...")
	return

	sidebandFingerprint, err := utils.GetOrcaFingerprint(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(s.T(), err)
	cmdFingerprint, stderrOutput, err := utils.RunBootstrapperWithIO(client, []string{"fingerprint"}, []string{}, []string{}, nil)
	log.Debug(stderrOutput)
	require.Nil(s.T(), err)

	require.Equal(s.T(), sidebandFingerprint, strings.TrimSpace(cmdFingerprint))
}

func (s *AcceptanceSuite) TestBadUninstall() {
	require := require.New(s.T())
	// If this test fails and actually uninstalls, other tests in this suite may blow up!

	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
	client, err := s.ControllerMachines[0].GetClient()
	require.Nil(err)
	// Without the ID, it should fail to uninstall
	require.NotNil(utils.RunBootstrapper(client, []string{"uninstall"}, []string{}, []string{}))

	// And orca should still be running
	require.Nil(utils.ValidateOrcaServerRunning(s.ControllerMachines[0], 30))
}

func (s *AcceptanceSuite) TestAddUserEmptyPass() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestAddUserEmptyPass(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestAddUser() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestAddUser(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestDisableUser() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestDisableUser(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestManagedMembers() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestManagedMembers(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestMultipleCerts() {
	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])

	// Now try to get a client as that user, and verify the system is accessible
	log.Debug("getting docker client as admin user")
	client1, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword()) // Gets a bundle each time
	require.Nil(err)
	client2, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword()) // Gets a bundle each time
	require.Nil(err)
	log.Debug("getting version from client1")
	ver, err := client1.Version()
	require.Nil(err)
	require.True(strings.Contains(ver.Version, "ucp"))

	log.Debug("getting version from client2")
	ver, err = client2.Version()
	require.Nil(err)
	require.True(strings.Contains(ver.Version, "ucp"))
}

func (s *AcceptanceSuite) TestNonAdminUserNoSwarmForYou() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(s.T(), err)
	utils.TestNonAdminUserNoSwarmForYou(s.T(), urls)
}

func (s *AcceptanceSuite) TestNonAdminUserNoProxyForYou() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(s.T(), err)
	utils.TestNonAdminUserNoProxyForYou(s.T(), urls)
}

func (s *AcceptanceSuite) TestAdminHasSwarm() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(s.T(), err)
	utils.TestAdminHasSwarm(s.T(), urls)
}

func (s *AcceptanceSuite) TestEscalationsBlocked() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	// TODO - fix it so this can run across controllers
	utils.TestEscalationsBlocked(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestBuildImage() {
	require := require.New(s.T())
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(err)

	for _, serverURL := range urls {
		client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(err)

		log.Debug("Building image")
		require.Nil(utils.BuildImage(client, "testimage:latest"))
		log.Debug("Inspecting image")
		_, err = client.InspectImage("testimage:latest")
		require.Nil(err)
		log.Debug("PASS")
	}
}

func (s *AcceptanceSuite) TestSimpleCompose() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	// TODO - fix it so this can run across controllers
	utils.TestSimpleCompose(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
}

func (s *AcceptanceSuite) TestSimpleDockerCLI() {
	require := require.New(s.T())
	for _, m := range s.ControllerMachines {
		serverURL, err := utils.GetOrcaURL(m)
		require.Nil(err)
		engineVersion, err := utils.GetEngineVersion(m)
		require.Nil(err)
		utils.TestSimpleCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), engineVersion)
	}
}

func (s *AcceptanceSuite) TestContainerListDockerCLI() {
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(s.T(), err)
	for _, serverURL := range urls {
		utils.TestContainerListCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	}
}

func (s *AcceptanceSuite) TestPrivilegedCLI() {
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(s.T(), err)
	for _, serverURL := range urls {
		utils.TestPrivilegedCLI(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	}
}

func (s *AcceptanceSuite) TestBridgeNetwork() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	// TODO - fix this so it can run across all controllers
	utils.TestBridgeNetwork(s.T(), serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
}

func (s *AcceptanceSuite) TestPingHealthWithKVStopped() {
	// TODO - All the stop tests are currently incompatible with the agent reconciliation - need to revisit
	s.T().Skip("Component stopping tests are not currently possible with agent reconciliation.")
	return

	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	log.Debug("Testing _ping endpoint with KV store shut down")

	for _, m := range s.ControllerMachines {
		client, err := m.GetClient()
		require.Nil(s.T(), err)

		status, err := utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 200)

		log.Debug("Stopping KV store")
		require.Nil(s.T(), client.StopContainer("ucp-kv", 20))

		time.Sleep(2 * time.Second)
		status, err = utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 500)

		log.Debug("Starting KV store")
		require.Nil(s.T(), client.StartContainer("ucp-kv", nil))

		time.Sleep(2 * time.Second)
		status, err = utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 200)

		serverURL, err := utils.GetOrcaURL(m)
		client, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), utils.WaitForEngineRecovery(client, 30))
	}
}

func (s *AcceptanceSuite) TestPingHealthWithSwarmStopped() {
	// TODO - All the stop tests are currently incompatible with the agent reconciliation - need to revisit
	s.T().Skip("Component stopping tests are not currently possible with agent reconciliation.")
	return

	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	log.Debug("Testing _ping endpoint with Swarm Manager shut down")

	for _, m := range s.ControllerMachines {
		client, err := m.GetClient()
		require.Nil(s.T(), err)

		status, err := utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 200)

		log.Debug("Stopping Swarm Manager")
		require.Nil(s.T(), client.StopContainer("ucp-swarm-manager", 20))

		time.Sleep(2 * time.Second)
		status, err = utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 500)

		log.Debug("Starting Swarm Manager")
		require.Nil(s.T(), client.StartContainer("ucp-swarm-manager", nil))

		time.Sleep(2 * time.Second)
		status, err = utils.PingOrcaServer(m)
		require.Nil(s.T(), err)
		require.Equal(s.T(), status, 200)

		serverURL, err := utils.GetOrcaURL(m)
		client, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(s.T(), utils.WaitForEngineRecovery(client, 30))

		// Make sure the cluster is healty too so we get swarm to recover
		require.Nil(s.T(), utils.ValidateClusterHealthy(serverURL, 30))
	}
}

func (s *AcceptanceSuite) TestStartOrderKV() {
	s.T().Skip("Start Order tests are no longer relevant with the ucp reconciler")
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	log.Debug("Testing start ordering with KV delayed")

	for _, m := range s.ControllerMachines {
		client, err := m.GetClient()
		require.Nil(s.T(), err)

		// Stop KV and controller, restart controller, wait a little while, then start kv
		log.Debug("Stopping KV")
		require.Nil(s.T(), client.StopContainer("ucp-kv", 20))
		log.Debug("Stopping controller")
		require.Nil(s.T(), client.StopContainer("ucp-controller", 20))
		log.Debug("Starting controller")
		require.Nil(s.T(), client.StartContainer("ucp-controller", nil))
		log.Debug("Delaying a moment")
		time.Sleep(2 * time.Second)
		log.Debug("Verifying controller isn't alive")
		require.NotNil(s.T(), utils.ValidateOrcaServerRunning(m, 1))
		log.Debug("Starting KV (expect a few errors below until it finishes coming up)")
		require.Nil(s.T(), client.StartContainer("ucp-kv", nil))
		for i := 0; i < 300; i++ {
			err = utils.ValidateOrcaServerRunning(m, 1)
			if err != nil {
				log.Debugf("Not ready yet: %s", err)
				time.Sleep(1 * time.Second)
			} else {
				log.Debug("Controller came up")
				break
			}
		}
		require.Nil(s.T(), err)
	}
}

func (s *AcceptanceSuite) TestStartOrderSwarmCA() {
	s.T().Skip("Start Order tests are no longer relevant with the ucp reconciler")
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	if len(s.ControllerMachines) > 1 {
		s.T().Skip("Multi-controller scenario - this test only works for single controllers")
		return
	}
	log.Debug("Testing start ordering with Swarm CA delayed")

	// TODO - if we replicate the CA, then this can be done across all controllers
	m := s.ControllerMachines[0]
	client, err := m.GetClient()
	require.Nil(s.T(), err)

	// Note: we use the proxy, since the linking gets messed up if we try to do the CA itself
	log.Debug("Stopping Swarm CA proxy")
	require.Nil(s.T(), client.StopContainer("ucp-cluster-root-ca", 20))
	log.Debug("Stopping controller")
	require.Nil(s.T(), client.StopContainer("ucp-controller", 20))
	log.Debug("Starting controller")
	require.Nil(s.T(), client.StartContainer("ucp-controller", nil))
	log.Debug("Delaying a moment")
	time.Sleep(2 * time.Second)
	log.Debug("Verifying controller is alive")
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(m, 20))
	log.Debug("Verifying we can't get an admin bundle")
	serverURL, err := utils.GetOrcaURL(m)
	require.Nil(s.T(), err)
	_, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword()) // Gets a bundle each time
	require.NotNil(s.T(), err)

	log.Debug("Starting Swarm CA proxy (expect a few errors below until it finishes coming up)")
	require.Nil(s.T(), client.StartContainer("ucp-cluster-root-ca", nil))
	for i := 0; i < 60; i++ {
		_, err = utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword()) // Gets a bundle each time
		if err != nil {
			log.Debugf("Not ready yet: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			log.Debug("Controller recovered")
			break
		}
	}
	require.Nil(s.T(), err)
}

func (s *AcceptanceSuite) TestStartOrderUCPCA() {
	// TODO - refine this test so we don't have to create a user!
	s.T().Skip("Start Order tests are no longer relevant with the ucp reconciler")
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	if len(s.ControllerMachines) > 1 {
		s.T().Skip("Multi-controller scenario - this test only works for single controllers")
		return
	}
	log.Debug("Testing start ordering with UCP CA delayed")

	// TODO - if we replicate the CA, then this can be done across all controllers
	m := s.ControllerMachines[0]
	client, err := m.GetClient()
	require.Nil(s.T(), err)

	log.Debug("Stopping UCP CA proxy")
	require.Nil(s.T(), client.StopContainer("ucp-client-root-ca", 20))
	log.Debug("Stopping controller")
	require.Nil(s.T(), client.StopContainer("ucp-controller", 20))
	log.Debug("Starting controller")
	require.Nil(s.T(), client.StartContainer("ucp-controller", nil))
	log.Debug("Delaying a moment")
	time.Sleep(2 * time.Second)
	log.Debug("Verifying controller is alive")
	require.Nil(s.T(), utils.ValidateOrcaServerRunning(m, 60))
	log.Debug("Verifying we can't get a user bundle")
	serverURL, err := utils.GetOrcaURL(m)
	require.Nil(s.T(), err)
	newUsername := "newrwuserN"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(s.T(), utils.CreateNewUser(nil, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))
	_, err = utils.GetUserDockerClient(serverURL, newUsername, newPassword)
	require.NotNil(s.T(), err)

	log.Debug("Starting UCP CA proxy (expect a few errors below until it finishes coming up)")
	require.Nil(s.T(), client.StartContainer("ucp-client-root-ca", nil))
	for i := 0; i < 60; i++ {
		_, err = utils.GetUserDockerClient(serverURL, newUsername, newPassword)
		if err != nil {
			log.Debugf("Not ready yet: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			log.Debug("Controller recovered")
			break
		}
	}
	require.Nil(s.T(), err)
}

func (s *AcceptanceSuite) TestUserOwnedContainers() {
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	t := s.T()
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(t, err)
	// TODO - fix it so this can run across controllers
	utils.TestUserOwnedContainers(t, serverURL)
}

func (s *AcceptanceSuite) TestTeamOwnedContainers() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	// Note: This isn't idempotent, so we can only run against one controller
	// TODO - fix it so this can run across controllers
	utils.TestTeamOwnedContainers(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestLDAPBasic() {
	// This test tries to set up ldap, which makes no sense if we're already LDAP mode
	if s.IsLDAP {
		s.T().Skip("LDAP enabled environment, skipping this test")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	utils.TestLDAPBasic(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestWebTests() {
	// TODO - this needs to be refactored with whatever new UI test framework we have...
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(s.T(), err)
	ui.TestUI(s.T(), serverURL)
}

func (s *AcceptanceSuite) TestNodeCount() {

	require := require.New(s.T())
	expected := len(s.ControllerMachines) + len(s.WorkerMachines)
	urls, err := utils.GetOrcaURLs(&s.OrcaTestSuite)
	require.Nil(err)

	for _, serverURL := range urls {
		nodes, err := utils.GetNodes(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(err)
		require.Equal(expected, len(nodes))

		// Also make sure the old classic swarm based node inventory matches
		classicNodes, err := utils.OldGetNodes(nil, serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
		require.Nil(err)
		require.Equal(expected, len(classicNodes))
	}
}

func (s *AcceptanceSuite) TestPrivatePullWithToken() {
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(s.T(), err)
	utils.TestPrivatePullWithToken(s.T(), client)
}

func (s *AcceptanceSuite) TestPullPrivateImage() {
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	if username == "" {
		s.T().Skip("Skipping private image pull test without REGISTRY_USERNAME and REGISTRY_PASSWORD set.")
		return
	}

	require := require.New(s.T())
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	require.Nil(client.PullImage("dockerorcadev/ucp:latest", &dockerclient.AuthConfig{
		Username: username,
		Password: password,
	}))
}
func (s *AcceptanceSuite) TestEngineDiscoveryAutomaticConfiguration() {
	s.T().Skip("Overlay networks not available until swarmkit is wired up")
	if s.IsDiscoveryMissing {
		s.T().Skip("Discovery not available, skipping related tests")
		return
	}
	require := require.New(s.T())
	// Engine discovery is only automatic on engine 1.11+
	minimumAutoDiscoveryVersion, _ := version.NewVersion("1.11.0")
	engineVersion, err := utils.GetEngineVersion(s.ControllerMachines[0])
	require.Nil(err)
	if engineVersion.LessThan(minimumAutoDiscoveryVersion) {
		s.T().Skip("skipping test -- needs engine 1.11 or above.")
		return
	}
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	client, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	// Now create an overlay network to make sure everything is happy on the back-end
	log.Infof("Attempting to create a network")

	// HACK - try to find a better way for this...
	for i := 0; i < 60; i++ {
		resp, err2 := client.CreateNetwork(&dockerclient.NetworkCreate{
			Name:   "orca-test-net",
			Driver: "overlay",
		})

		if err2 == nil {
			log.Infof("Network create: %s - %s", resp.ID, resp.Warning)
			err = nil // Clear it since we're good now
			break
		} else if strings.Contains(err2.Error(), "No healthy node available in the cluster") {
			log.Info("Swarm complaining about no healthy nodes, retrying in 2 secs...")
		} else {
			log.Warn(err2) // XXX Short circuit and die?
		}
		err = err2
		time.Sleep(2 * time.Second)
	}
	require.Nil(err)
	// TODO - try to create containers across the networks and verify they worked
}
func (s *AcceptanceSuite) TestRequireTrustWithHub() {
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	if username == "" {
		s.T().Skip("Skipping require trust with hub test without REGISTRY_USERNAME and REGISTRY_PASSWORD set.")
		return
	}

	require := require.New(s.T())

	// Set require trust from hub
	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)
	config := manager.TrustConfiguration{
		RequireContentTrustForDTR: true,
		RequireContentTrustForHub: true,
	}
	require.Nil(updateConfig(serverURL, "/api/config/trust", config))

	dockerClient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "signedImgService",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:   "alpine:latest",
				Command: []string{"sh"},
			},
		},
	}
	// create service
	err = createService(dockerClient, username, password, serviceSpec)
	// should have no errors since alpine latest image should always be signed
	require.Nil(err)

	// the latest tag will be inferred if we don't give a tag, and it is signed
	serviceSpecNoTag := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "signedImgServiceNoTag",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:   "alpine",
				Command: []string{"sh"},
			},
		},
	}
	// create service
	err = createService(dockerClient, username, password, serviceSpecNoTag)
	// should have no errors since alpine latest image should always be signed
	require.Nil(err)

	// create service with unsigned image with explicit latest tag
	serviceSpec.Annotations.Name = "unsignedImgService"
	serviceSpec.TaskTemplate.ContainerSpec.Image = "riyaz/unsigned-img:latest"
	err = createService(dockerClient, username, password, serviceSpec)
	// should error since trust is required but image is unsigned
	// if the error is different, make sure riyaz/unsigned-img:latest still exists and is unsigned
	require.Equal(err.Error(), "Error response from daemon: image or trust data does not exist for docker.io/riyaz/unsigned-img:latest")

	// create service with unsigned image, without tag - latest will be inferred
	serviceSpec.Annotations.Name = "unsignedImgServiceNoTag"
	serviceSpec.TaskTemplate.ContainerSpec.Image = "riyaz/unsigned-img"
	err = createService(dockerClient, username, password, serviceSpec)
	// should error since trust is required but image is unsigned
	// if the error is different, make sure riyaz/unsigned-img:latest still exists and is unsigned
	require.Equal(err.Error(), "Error response from daemon: image or trust data does not exist for docker.io/riyaz/unsigned-img:latest")

	// turn off content trust requirement
	config.RequireContentTrustForDTR = false
	config.RequireContentTrustForHub = false
	require.Nil(updateConfig(serverURL, "/api/config/trust", config))

	// create service with unsigned image, which should work
	err = createService(dockerClient, username, password, serviceSpec)
	require.Nil(err)
}

func updateConfig(serverURL, configSubPath string, config interface{}) error {
	tlsConfig, err := utils.GetUserTLSConfig(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	if err != nil {
		return err
	}

	tlsConfig.InsecureSkipVerify = true
	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}
	client := &http.Client{
		Transport: tr,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	fullConfigPath := serverURL + configSubPath
	req, err := http.NewRequest("POST", fullConfigPath, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Close = true
	_, err = client.Do(req)
	return err
}

func createService(client *dockerclient.DockerClient, username, password string, spec swarm.ServiceSpec) error {
	engineClient, err := utils.ConvertToEngineAPI(client)
	if err != nil {
		return err
	}

	createOptions := types.ServiceCreateOptions{}
	if username != "" {
		authConfig := types.AuthConfig{
			Username: username,
			Password: password,
		}

		encodedJson, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}

		createOptions.EncodedRegistryAuth = base64.StdEncoding.EncodeToString(encodedJson)
	}
	_, err = engineClient.ServiceCreate(context.TODO(), spec, createOptions)
	return err
}

func TestAcceptance(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &AcceptanceSuite{}
	s.Init(s)
	suite.Run(t, s)
}
