package utils

import (
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"

	bootUtils "github.com/docker/orca/bootstrap/client"
)

func (s *OrcaTestSuite) Init(ss OTSC) {
	s.self = ss // Wire up the reference so we can use derived implementations
	s.Data = map[string][]byte{}
}

func (s *OrcaTestSuite) SetupSuite() {
	// This routine leverages the internal "self" interface reference so we
	// can get subtype overridden behavior when applicable
	require := require.New(s.T())
	controllerCount, workerCount := s.self.GetNodeCounts()
	// If we're in inventory mode, make sure we actually have enough potential machines, else skip

	/* TODO - uncomment once inventory model is merged
	inventoryCount, err := GetTotalInventoryCount()
	require.Nil(err)

	log.Debugf("Total machine inventory: %d", inventoryCount)
	log.Infof("Initializing suite with %d controllers and %d workers", controllerCount, workerCount)

	if os.Getenv("MACHINE_INVENTORY") != "" && inventoryCount < controllerCount+workerCount {
		s.T().Skip("Not enough total machines in the inventory for this test")
		return
	}
	*/

	machines, err := GetTestMachines(controllerCount + workerCount)
	require.Nil(err)
	s.ControllerMachines = machines[0:controllerCount]
	s.WorkerMachines = machines[controllerCount:]

	require.NotEqual(len(s.ControllerMachines), 0, "You must call Init(...) before starting suite")
	if os.Getenv("MACHINE_LOCAL") != "" && len(s.ControllerMachines)+len(s.WorkerMachines) > 1 {
		s.T().Skip("Multinode test skipping for MACHINE_LOCAL")
		return
	}

	require.Nil(s.self.PreImageLoad())

	m := s.ControllerMachines[0]
	client, err := m.GetClient()
	require.Nil(err)

	require.Nil(LoadAllLocalOrcaImages(client))

	require.Nil(s.self.PreInstall(m))

	args := s.self.InstallArgs(m)

	require.Nil(RunBootstrapper(client, args, []string{}, []string{}))
	require.Nil(s.self.PostInstall(m))

	// TODO - consider trying to parallelize some aspects of this...

	// Now join all controllers
	require.Nil(err)
	serverURL, err := GetOrcaInternalURL(m)
	require.Nil(err)
	for _, m := range s.ControllerMachines[1:] {
		client, err = m.GetClient()
		require.Nil(err)
		require.Nil(LoadAllLocalOrcaImages(client)) // TODO parallel?
		require.Nil(s.self.PreInstall(m))

		log.Debugf("Joining replica %s to: %s", m.GetName(), serverURL)
		require.Nil(JoinNode(s.ControllerMachines[0], m, true))
		// Wait for the UCP controller to come up on that node
		addr, err := GetOrcaInternalURL(m)
		require.Nil(err)
		require.Nil(bootUtils.WaitForOrca(addr, 5*time.Minute))
		require.Nil(s.self.PostInstall(m))
	}

	// Now join all workers
	for _, m := range s.WorkerMachines {
		client, err = m.GetClient()
		require.Nil(err)
		require.Nil(LoadAllLocalOrcaImages(client)) // TODO parallel?
		require.Nil(s.self.PreInstall(m))

		log.Debugf("Joining non-replica %s to: %s", m.GetName(), serverURL)
		require.Nil(JoinNode(s.ControllerMachines[0], m, false))
		// Check for the node in a loop before joining the next one
		detected := false
		for i := 0; i < 60 && !detected; i++ {
			time.Sleep(2 * time.Second)
			nodes, err := GetNodes(serverURL, GetAdminUser(), GetAdminPassword())
			require.Nil(err)
			for _, node := range nodes {
				if node.Description.Hostname == m.GetName() {
					detected = true
					break
				}
			}
		}
		require.True(detected, "Failed to detect new joined node")
		require.Nil(s.self.PostInstall(m))
	}
	log.Info("Running post-install configuration tasks")
	require.Nil(s.self.InitialConfig())
	log.Info("All machines installed and joined for test suite")
}

func (s *OrcaTestSuite) TearDownSuite() {
	log.Info("Tearing down suite")
	var wg sync.WaitGroup
	for _, m := range s.ControllerMachines {
		wg.Add(1)
		go func(m Machine) {
			defer wg.Done()
			err := m.Remove()
			// err := m.Finished(s.T().Failed()) // TODO
			if err != nil {
				log.Info("Failed to finish machine: %s: %s", m.GetName(), err)
			}
		}(m)

	}
	for _, m := range s.WorkerMachines {
		wg.Add(1)
		go func(m Machine) {
			defer wg.Done()
			err := m.Remove()
			// err := m.Finished(s.T().Failed()) // TODO
			if err != nil {
				log.Info("Failed to finish machine: %s: %s", m.GetName(), err)
			}
		}(m)
	}
	wg.Wait()
}

func (s *OrcaTestSuite) GetServerURLs() []string {
	serverURLs := []string{}
	for _, m := range s.ControllerMachines {
		serverURL, err := GetOrcaURL(m)
		require.Nil(s.T(), err)
		serverURLs = append(serverURLs, serverURL)
	}
	return serverURLs
}
func (s *OrcaTestSuite) RestartAllControllers() error {
	log.Info("Restarting all ucp-controllers")
	clients := []*dockerclient.DockerClient{}
	for _, m := range s.ControllerMachines {
		client, err := m.GetClient()
		if err != nil {
			return err
		}
		clients = append(clients, client)
	}

	for _, client := range clients {
		err := client.RestartContainer("ucp-controller", 20)
		if err != nil {
			return err
		}
	}
	log.Info("Waiting for controllers to recover")
	for _, m := range s.ControllerMachines {
		err := ValidateOrcaServerRunning(m, 60)
		if err != nil {
			return err
		}
	}
	return nil
}

// Take a backup of machine m, storing as s.Data["<machine name>:backup"]
func (s *OrcaTestSuite) TakeBackup(m Machine) error {
	log.Info("Taking a backup of the HA cluster")
	serverURL, err := GetOrcaURL(m)
	if err != nil {
		return err
	}

	//Now to initial backup
	client, err := m.GetClient()
	if err != nil {
		return err
	}
	ClusterID, err := GetOrcaID(client)
	if err != nil {
		return err
	}
	log.Info("Taking full backup of the controller")
	backup, stderrOutput, err := RunBootstrapperWithIO(client, []string{"backup", "-D", "--id", ClusterID}, []string{}, []string{}, nil)
	log.Debug(stderrOutput)
	if err != nil {
		return err
	}
	s.Data[m.GetName()+":backup"] = []byte(backup)

	log.Info("Verifying the system recovers")
	for _, m = range s.ControllerMachines {
		err = ValidateOrcaServerRunning(m, 100)
		if err != nil {
			return err
		}
	}
	err = ValidateClusterHealthy(serverURL, 100)
	if err != nil {
		return err
	}
	err = WaitForCLIReadiness(serverURL, GetAdminUser(), GetAdminPassword(), time.Minute)
	if err != nil {
		return err
	}
	return nil
}

// Default routines
func (s *OrcaTestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	log.Errorf("Derived test should implement GetNodeCounts at a minimum!")
	return 0, 0
}
func (s *OrcaTestSuite) PreImageLoad() error {
	return nil
}

func (s *OrcaTestSuite) PreInstall(m Machine) error {
	return nil
}
func (s *OrcaTestSuite) InstallArgs(m Machine) []string {
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	internalIP, err := m.GetInternalIP()
	require.Nil(s.T(), err)
	return []string{
		"install",
		"--disable-tracking",
		"--disable-usage",
		"-D",
		"--swarm-port", "3376",
		"--san", externalIP,
		"--host-address", internalIP,
	}
}
func (s *OrcaTestSuite) JoinReplicaArgs(m Machine, serverURL, fingerprint string) [3][]string {
	args := s.JoinWorkerArgs(m, serverURL, fingerprint)
	args[0] = append(args[0], "--replica")
	return args
}
func (s *OrcaTestSuite) JoinWorkerArgs(m Machine, serverURL, fingerprint string) [3][]string {
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	internalIP, err := m.GetInternalIP()
	require.Nil(s.T(), err)
	return [3][]string{
		{
			"join", "-D",
			"--fingerprint", fingerprint,
			"--url", serverURL,
			"--swarm-port", "3376",
			"--san", externalIP,
			"--host-address", internalIP,
		},
		{
			"UCP_ADMIN_USER=" + GetAdminUser(),
			"UCP_ADMIN_PASSWORD=" + GetAdminPassword(),
		},
		{},
	}
}

// Default adds a license for first node
func (s *OrcaTestSuite) PostInstall(m Machine) error {
	if m == s.ControllerMachines[0] {
		serverPublicURL, err := GetOrcaURL(m)
		if err != nil {
			return err
		}
		AddValidLicense(s.T(), serverPublicURL, GetAdminUser(), GetAdminPassword())
	}
	return nil
}
func (s *OrcaTestSuite) InitialConfig() error {
	return nil
}
