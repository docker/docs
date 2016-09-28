package utils

import (
	"os"
	"strconv"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Basic test rigging that spins up a controller machine but does NOT install
type TestSuite struct {
	suite.Suite
	ControllerMachine Machine
	PrimaryVersion    *version.Version
	preImageLoadHook  func(s *TestSuite)
}

func (s *TestSuite) SetupSuite() {
	machine, err := CreateTestMachine()
	require.Nil(s.T(), err)
	s.ControllerMachine = machine

	client, err := machine.GetClient()
	require.Nil(s.T(), err)

	// Stash away this machine's version.
	// XXX eventually we'll want a map for tracking this on a per-node basis
	s.PrimaryVersion, err = getEngineVersionInformation(client)
	require.Nil(s.T(), err)
	log.Infof("Primary engine version: %s", s.PrimaryVersion)

	if s.preImageLoadHook != nil {
		s.preImageLoadHook(s)
	}
	if os.Getenv("MACHINE_LOCAL") != "" {
		log.Info("Skipping image load for MACHINE_LOCAL run")
	} else {
		require.Nil(s.T(), LoadAllLocalOrcaImages(client))
	}
}

func GetEngineVersion(m Machine) (primaryVersion *version.Version, err error) {
	client, err := m.GetClient()
	if err != nil {
		return nil, err
	}
	return getEngineVersionInformation(client)
}

func getEngineVersionInformation(client *dockerclient.DockerClient) (primaryVersion *version.Version, err error) {
	verInfo, err := client.Version()
	if err != nil {
		return primaryVersion, err
	}

	// Chop off the trailing CS indicator
	// XXX This leaves us unable to tell apart RCs and regular releases!
	versionParts := strings.Split(verInfo.Version, "-")
	primaryVersion, err = version.NewVersion(versionParts[0])
	return primaryVersion, err
}

func (s *TestSuite) TearDownSuite() {
	if s.T().Failed() {
		log.Errorf("Test failed, not deleting machine %s", s.ControllerMachine.GetName())
		if os.Getenv("PRESERVE_TEST_MACHINE") == "" {
			// If PRESERVE_TEST_MACHINE is not set, stop the machine to avoid taking up too many resources
			s.ControllerMachine.Stop()
		}
		return
	}
	if os.Getenv("PRESERVE_TEST_MACHINE") == "" {
		require.Nil(s.T(), s.ControllerMachine.Remove())
	}
}

func HandleTestArgs(t *testing.T) {
	// Our jenkins is being very janky and is breaking when we have a lot of
	// log output, so we had to disable verbose logging. If we set logging to
	// debug as a default, then we will get debug-level logging whenever a test
	// fails, which should still work there and provide us more output.
	log.SetLevel(log.DebugLevel)

	if os.Getenv("MACHINE_DRIVER") == "" && os.Getenv("MACHINE_LOCAL") == "" {
		t.Skip("skipping integration test without $MACHINE_DRIVER or $MACHINE_LOCAL set.")
	}
}

// Perhaps this belongs someplace else...
func GetStressObjectCount(defaultCount int) int {
	v := os.Getenv("STRESS_OBJECT_COUNT")
	if v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			return val
		} else {
			log.Error("STRESS_OBJECT_COUNT must be an int")
		}
	}
	return defaultCount
}

// Test rigging for BYO Orca
type OrcaClientTestSuite struct {
	suite.Suite
	Client *dockerclient.DockerClient
}

func (s *OrcaClientTestSuite) SetupSuite() {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		s.T().Skip("skipping OrcaClient based test as DOCKER_HOST is not set.")
		return
	} else {
		if testing.Verbose() {
			log.SetLevel(log.DebugLevel)
		}
		// Kinda overkill, but leveraging what's already there...
		localMachine, err := NewLocalMachine()
		if err != nil {
			s.T().Skip("Unable to connect to " + dockerHost)
			return
		}

		client, err := localMachine.GetClient()
		if err != nil {
			s.T().Skip("Unable to connect to " + dockerHost)
			return
		}

		// Now verify it's actually UCP
		version, err := client.Version()
		if err != nil {
			s.T().Skip("Failed to get version from " + dockerHost)
			return
		}

		if !strings.Contains(version.Version, "ucp") {
			s.T().Skip("DOCKER_HOST does not appear to point to UCP " + version.Version)
			return
		}
		s.Client = client
	}
}
