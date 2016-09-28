package restarter

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/utils/registry"
)

const testImage = "busybox:latest"

var testCmd = []string{"tail", "-f", "/etc/hosts"}

type TestSuite struct {
	suite.Suite

	// Containers is a map from containerID to a PID
	containers map[string]int

	// sampleID is the ID of one of the containers
	sampleID string

	ec *client.EngineClient
}

func (s *TestSuite) SetupSuite() {
	require := require.New(s.T())
	ec, err := client.Mock()
	require.Nil(err)
	s.ec = ec

	progress, err := s.ec.PullImage(testImage, types.ImagePullOptions{
		PrivilegeFunc: registry.RequestPrivilegeFunc,
	})
	require.Nil(err)
	require.NotNil(progress)
	defer progress.Close()

	_, err = ioutil.ReadAll(progress)
	require.Nil(err)

	s.containers = make(map[string]int)

	for i := 1; i <= 3; i++ {
		id, err := containerFactory(s.ec)
		require.Nil(err)
		s.containers[id] = 0
		s.sampleID = id
	}
}

func (s *TestSuite) TearDownSuite() {
	for ctr, _ := range s.containers {
		s.ec.RemoveContainer(ctr, true, false)
	}
}

// containerFactory starts a test container and returns a containerID and error
func containerFactory(ec *client.EngineClient) (string, error) {
	cfg := &container.Config{
		Image: testImage,
		Cmd:   testCmd,
	}
	hostCfg := &container.HostConfig{}
	res, err := ec.CreateContainer(cfg, hostCfg, "")
	if err != nil {
		return "", err
	}
	err = ec.StartContainer(res.ID)
	return res.ID, err
}

func (s *TestSuite) TestRestarterNoSuchID() {
	require := require.New(s.T())
	restarter := NewRestarter(s.ec)

	// Try to restart a container that is not configured for this restarter
	err := restarter.RestartContainer("nosuchcontainerid")
	require.NotNil(err)
	require.True(strings.Contains(err.Error(), "is not configured to restart"))
}

func (s *TestSuite) TestRestarterNoContainers() {
	require := require.New(s.T())
	restarter := NewRestarter(s.ec)

	// Try to restart a container that is not configured for this restarter
	err := restarter.RestartAll()
	require.NotNil(err)
	require.True(strings.Contains(err.Error(), "No containers to be restarted"))
}

func (s *TestSuite) TestRestarterOneContainer() {
	require := require.New(s.T())
	restarter := NewRestarter(s.ec)

	// Start a container and register it in the restarter
	restarter.SetContainerTargets([]string{s.sampleID})

	// Get the current restart count of the container
	info, err := s.ec.InspectContainer(s.sampleID)
	prevPid := info.State.Pid

	// Restart the container
	err = restarter.RestartContainer(s.sampleID)
	require.Nil(err)

	time.Sleep(200 * time.Millisecond)

	// Require that the running process Pid is not the same, indicating a restart
	info, err = s.ec.InspectContainer(s.sampleID)
	newPid := info.State.Pid

	require.NotEqual(prevPid, newPid)
}

func (s *TestSuite) TestRestarterOneContainerTwice() {
	require := require.New(s.T())
	restarter := NewRestarter(s.ec)

	// Register a container in the restarter
	restarter.SetContainerTargets([]string{s.sampleID})

	// Get the current restart count of the container
	info, err := s.ec.InspectContainer(s.sampleID)
	prevPid := info.State.Pid

	// Restart the container
	err = restarter.RestartContainer(s.sampleID)
	require.Nil(err)

	time.Sleep(200 * time.Millisecond)

	// Require that the running process Pid is not the same, indicating a restart
	info, err = s.ec.InspectContainer(s.sampleID)
	newPid := info.State.Pid

	require.NotEqual(prevPid, newPid)

	// Restart the container again, require a failure
	err = restarter.RestartContainer(s.sampleID)
	require.NotNil(err)
	require.True(strings.Contains(err.Error(), "this container has been already restarted"))

	time.Sleep(200 * time.Millisecond)

	// Require that the running process Pid is the same
	info, err = s.ec.InspectContainer(s.sampleID)
	require.Equal(newPid, info.State.Pid)
}

func (s *TestSuite) TestRestartAll() {
	require := require.New(s.T())
	restarter := NewRestarter(s.ec)

	containerIDs := []string{}
	for id, _ := range s.containers {
		info, err := s.ec.InspectContainer(id)
		require.Nil(err)

		s.containers[id] = info.State.Pid
		containerIDs = append(containerIDs, id)
	}

	// Register all containers in the restarter
	restarter.SetContainerTargets(containerIDs)

	// Restart one of the containers manually
	err := restarter.RestartContainer(s.sampleID)
	require.Nil(err)

	// Restart all containers with RestartAll
	err = restarter.RestartAll()
	require.Nil(err)

	time.Sleep(200 * time.Millisecond)

	// Ensure all the registered containers have been restarted
	for id, prevPid := range s.containers {
		info, err := s.ec.InspectContainer(id)
		require.Nil(err)
		require.NotEqual(prevPid, info.State.Pid)
	}
}

func TestRestarter(t *testing.T) {
	t.Skip("Skipping for now - need to refactor this to use mock's so this doesn't tickle the engine")
	return
	suite.Run(t, new(TestSuite))
}
