package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/dhe-deploy/shared/containers"

	"github.com/docker/engine-api/types"
	"github.com/stretchr/stew/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RemoteTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u                 *util.Util
	dtrContainerNames []string
}

func (suite *RemoteTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)

	suite.dtrContainerNames = make([]string, 0)
	containers := containers.AllContainers
	for _, v := range containers {
		suite.dtrContainerNames = append(suite.dtrContainerNames, v.FullName())
	}
}

func (suite *RemoteTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *RemoteTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

// Access to exit status
func (suite *RemoteTestSuite) TestRemoteHostCommandDoesntExist() {
	_, _, exitcode := util.RawExecute(suite.T(), suite.SSH, "this_command_doesnt_exist")
	if exitcode != 127 {
		suite.T().Fatal("command should have returned exit code 127")
	}
}

func (suite *RemoteTestSuite) TestRemoteHostCommandExistsButFails() {
	_, _, exitcode := util.RawExecute(suite.T(), suite.SSH, "pidof command_that_exists_but_fails")
	if exitcode != 1 {
		suite.T().Fatal("command should have returned exit code 1")
	}
}

func (suite *RemoteTestSuite) TestRemoteHostDockerRunning() {
	_, _, exitcode := util.RawExecute(suite.T(), suite.SSH, "pidof docker")
	_, _, exitcode2 := util.RawExecute(suite.T(), suite.SSH, "pidof dockerd")
	if exitcode != 0 && exitcode2 != 0 {
		suite.T().Fatal("command should have returned zero exit code")
	}
}

func (suite *RemoteTestSuite) DockerCLIInspectRunningContainer(containerName string) *types.ContainerJSONBase {
	inspect := suite.DockerCLIInspect(containerName)
	if !inspect.State.Running {
		suite.T().Fatal("Unable to inspect running container , as the container is not running")
	}
	return inspect
}

// DockerCLIInspect uses the structs from docker engine API types to unmarshall JSON
// from `docker inspect container_name` which itself is the JSON from the response of Remote API:
// GET "/containers/{containername:.*}/json"
func (suite *RemoteTestSuite) DockerCLIInspect(containerName string) *types.ContainerJSONBase {
	inspectOutput := suite.DockerCLIInspectExecute(containerName)
	var containerInspect types.ContainerJSONBase
	err := json.Unmarshal([]byte(inspectOutput), &containerInspect)
	if err != nil {
		suite.T().Fatalf("Failed to parse JSON : %v", err)
	}
	return &containerInspect
}

// DockerCLIInspectMap is like DockerCLIInspect, except that it returns a simplified
// helper map to get any data from docker inspect. This can be used in testing multiple
// docker versions, as the API JSON format can and will vary across versions so it might
// be neccesary to just directly access a workable unmarshalled form of the inspect JSON.
func (suite *RemoteTestSuite) DockerCLIInspectMap(containerName string) *objects.Map {
	inspectOutput := suite.DockerCLIInspectExecute(containerName)
	m, err := objects.NewMapFromJSON(inspectOutput)
	if err != nil {
		suite.T().Fatalf("Failed to parse JSON : %v", err)
	}
	return &m
}

// DockerCLIInspectExecute wraps the `docker inspect container_name` call for both parsing methods
func (suite *RemoteTestSuite) DockerCLIInspectExecute(containerName string) string {
	output := util.Execute(suite.T(), suite.SSH, "sudo docker inspect "+containerName, false)
	// trim some output to avoid parse errors
	output = strings.TrimPrefix(output, "[")
	output = strings.TrimSuffix(output, "]")
	return output
}

// TestDockerDaemonRunning uses the IsProcessRunning for the process "docker" and tests whether that returns as true to see if the docker daemon is up.
func (suite *RemoteTestSuite) TestDockerDaemonRunning() {
	assert.True(suite.T(), suite.IsProcessRunning("docker") || suite.IsProcessRunning("dockerd"), "docker is not currently running")
}

// // TestDTRContainersRunning calls the ContainerStateRunning function to determine if all dtr containers are running.
// func (suite *RemoteTestSuite) TestDTRContainersRunning() {
// 	output := suite.ContainerStateRunning(suite.dtrContainerNames)
// 	for i, v := range strings.Split(output, "\n") {
// 		assert.Equal(suite.T(), v, "true", fmt.Sprintf("Container %s not currently running", suite.dtrContainerNames[i]))
// 	}
// }
//
// func (suite *RemoteTestSuite) TestForDTRLogErrors() {
// 	util.AppendDTRIgnorableLoggedErrors(util.KnownDTRInstallLogErrors)
// 	util.AppendDockerIgnorableLoggedErrors(util.KnownDockerInstallLogErrors)
// 	util.TestLogs(suite.T(), suite.ssh)
// }

// IsProcessRunning checks to see if the process is running by returning the exitcode of "pidof processName".
func (suite *RemoteTestSuite) IsProcessRunning(processName string) bool {
	_, _, exitcode := util.RawExecute(suite.T(), suite.SSH, fmt.Sprintf("pidof %s", processName))
	return exitcode == 0
}

// ContainerStateRunning gets the values from docker inspect for the state of all the passed in containers in terms of whether they are running or not.
func (suite *RemoteTestSuite) ContainerStateRunning(containers []string) string {
	names := strings.Join(containers, " ")
	output := util.Execute(suite.T(), suite.SSH, fmt.Sprintf("sudo docker inspect --format '{{ .State.Running }}' %s", names), false)
	return output
}

func TestRemoteSuite(t *testing.T) {
	suite.Run(t, new(RemoteTestSuite))
}
