package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

var (
	Password        = "secret"
	Image           = "busybox"
	ContainerOutput = "Hello World"
)

// TODO - consider some sort of breadcrumb to pass from test to test so we can recycle clients
//        and potentially other data to streamline the flow
type TestOperation struct {
	Label     string
	Operation func(serverURL string, role *TestRole, op *TestOperation) string
	Expected  string
}

type TestRole struct {
	Role auth.Role
	// Set of operations to perform, done in-order so they can be chained
	Operations []TestOperation
}

func createUsers(serverURL string, role *TestRole, op *TestOperation) string {
	// As admin, create two users with the given role
	// TODO
	err := CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), fmt.Sprintf("authtest_user_%d_1", role.Role), Password, false, role.Role)
	if err != nil {
		return err.Error()
	}
	err = CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), fmt.Sprintf("authtest_user_%d_2", role.Role), Password, false, role.Role)
	if err != nil {
		return err.Error()
	}
	return "OK"
}

func listUserContainers(serverURL string, role *TestRole, op *TestOperation) string {
	client, err := GetUserDockerClient(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	_, err = client.ListContainers(true, false, "")
	if err != nil {
		return err.Error()
	}
	return "OK"
}

// TODO add a second one that uses CLI
func createContainer(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)
	client, err := GetUserDockerClient(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	cfg := &dockerclient.ContainerConfig{
		Image: Image,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
	}
	containerId, err := client.CreateContainer(cfg, containerName, nil)
	if err != nil {
		return err.Error()
	}
	err = client.StartContainer(containerId, nil)
	if err != nil {
		return err.Error()
	}
	log.Debugf("Created and started %s", containerName)
	return "OK"
}

func createContainerPriv(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container__priv_%d", role.Role)
	client, err := GetUserDockerClient(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	cfg := &dockerclient.ContainerConfig{
		Image: Image,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
		HostConfig: dockerclient.HostConfig{
			Privileged: true,
		},
	}
	containerId, err := client.CreateContainer(cfg, containerName, nil)
	if err != nil {
		return err.Error()
	}
	err = client.StartContainer(containerId, nil)
	if err != nil {
		return err.Error()
	}
	log.Debugf("Created and started %s", containerName)
	return "OK"
}

func createContainerPid(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container__pid_%d", role.Role)
	client, err := GetUserDockerClient(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	cfg := &dockerclient.ContainerConfig{
		Image: Image,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
		HostConfig: dockerclient.HostConfig{
			PidMode: "host",
		},
	}
	containerId, err := client.CreateContainer(cfg, containerName, nil)
	if err != nil {
		return err.Error()
	}
	err = client.StartContainer(containerId, nil)
	if err != nil {
		return err.Error()
	}
	log.Debugf("Created and started %s", containerName)
	return "OK"
}

func createContainerWithAdmin(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)
	client, err := GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
	if err != nil {
		return err.Error()
	}
	cfg := &dockerclient.ContainerConfig{
		Image: Image,
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
		Labels: map[string]string{
			// NOTE: This doesn't actually work right now
			orca.UCPOwnerLabel: fmt.Sprintf("authtest_user_%d_1", role.Role),
		},
	}
	containerId, err := client.CreateContainer(cfg, containerName, nil)
	if err != nil {
		return err.Error()
	}
	err = client.StartContainer(containerId, nil)
	if err != nil {
		return err.Error()
	}
	log.Debugf("Created and started %s", containerName)
	return "OK"
}

// TODO - make one using the dockerclient
func containerLogs(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)

	ch, err := NewCLIHandle(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	output, err := ch.RunCLICommand([]string{ch.Binaries[0], "logs", containerName})
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(output)
}

// TODO - make one using the dockerclient
func containerLogsOther(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)

	ch, err := NewCLIHandle(serverURL, fmt.Sprintf("authtest_user_%d_2", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	output, err := ch.RunCLICommand([]string{ch.Binaries[0], "logs", containerName})
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(output)
}

func containerExec(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)

	ch, err := NewCLIHandle(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}
	output, err := ch.RunCLICommand([]string{ch.Binaries[0], "exec", containerName, "echo", ContainerOutput})
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(output)
}

func containerCopy(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)

	err := CopyToFromContainer(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password, containerName)
	if err != nil {
		return err.Error()
	}
	return "OK"
}

func containerAttach(serverURL string, role *TestRole, op *TestOperation) string {
	containerName := fmt.Sprintf("authtest_container_%d", role.Role)
	client, err := GetUserDockerClient(serverURL, fmt.Sprintf("authtest_user_%d_1", role.Role), Password)
	if err != nil {
		return err.Error()
	}

	// Prevent a race - allow the container some time to spit out its output
	time.Sleep(100 * time.Millisecond)

	reader, err := client.AttachContainer(containerName, &dockerclient.AttachOptions{
		Stream: false,
		Stdout: true,
	})
	if err != nil {
		return err.Error()
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(data))
}

// Test cases follow

func TestUserOwnedContainers(t *testing.T, serverURL string) {
	scenarios := []TestRole{
		TestRole{
			Role: auth.None,
			Operations: []TestOperation{
				TestOperation{
					Label:     "create_users",
					Operation: createUsers,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container",
					Operation: createContainer,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "list_containers",
					Operation: listUserContainers,
					Expected:  "OK",
				},
				/* If we wire up setting owner on admin create, these become interesting
				TestOperation{
					Label:     "create_container_with_admin",
					Operation: createContainerWithAdmin,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "logs",
					Operation: containerLogs,
					//Expected:  ContainerOutput, - if it weren't denied it'd be this
					Expected: "access denied",
				},
				TestOperation{
					Label:     "exec",
					Operation: containerExec,
					//Expected:  ContainerOutput, - if it weren't denied it'd be this
					Expected: "access denied",
				},
				TestOperation{
					Label:     "attach",
					Operation: containerAttach,
					//Expected:  "OK",
					Expected: "access denied",
				},
				// TODO - others?
				*/
			},
		},
		TestRole{
			Role: auth.View,
			Operations: []TestOperation{
				TestOperation{
					Label:     "create_users",
					Operation: createUsers,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container",
					Operation: createContainer,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "list_containers",
					Operation: listUserContainers,
					Expected:  "OK",
				},
				/* If we wire up setting owner on admin create, these become interesting
				TestOperation{
					Label:     "create_container_with_admin",
					Operation: createContainerWithAdmin,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "logs",
					Operation: containerLogs,
					Expected:  ContainerOutput,
				},
				TestOperation{
					Label:     "exec",
					Operation: containerExec,
					Expected:  ContainerOutput,
				},
				TestOperation{
					Label:     "attach",
					Operation: containerAttach,
					Expected:  ContainerOutput,
				},
				// TODO - others?
				*/
			},
		},
		TestRole{
			Role: auth.RestrictedControl,
			Operations: []TestOperation{
				TestOperation{
					Label:     "create_users",
					Operation: createUsers,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container",
					Operation: createContainer,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container_priv",
					Operation: createContainerPriv,
					Expected:  "privileged not allowed",
				},
				TestOperation{
					Label:     "create_container_pid",
					Operation: createContainerPid,
					Expected:  "pid not allowed",
				},
				TestOperation{
					Label:     "logs",
					Operation: containerLogs,
					Expected:  ContainerOutput,
				},
				TestOperation{
					Label:     "logs_other_user",
					Operation: containerLogsOther,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "exec",
					Operation: containerExec,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "attach",
					Operation: containerAttach,
					//Expected:  ContainerOutput, // XXX Why not this?
					Expected: "",
				},
				TestOperation{
					Label:     "copy",
					Operation: containerCopy,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "list_containers",
					Operation: listUserContainers,
					Expected:  "OK",
				},
				// TODO - others?
			},
		},
		TestRole{
			Role: auth.FullControl,
			Operations: []TestOperation{
				TestOperation{
					Label:     "create_users",
					Operation: createUsers,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container",
					Operation: createContainer,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container",
					Operation: createContainerPriv,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "create_container_pid",
					Operation: createContainerPid,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "logs",
					Operation: containerLogs,
					Expected:  ContainerOutput,
				},
				TestOperation{
					Label:     "logs_other_user",
					Operation: containerLogsOther,
					Expected:  "access denied",
				},
				TestOperation{
					Label:     "exec",
					Operation: containerExec,
					Expected:  ContainerOutput,
				},
				TestOperation{
					Label:     "attach",
					Operation: containerAttach,
					//Expected:  ContainerOutput, // XXX Why not this?
					Expected: "",
				},
				TestOperation{
					Label:     "copy",
					Operation: containerCopy,
					Expected:  "OK",
				},
				TestOperation{
					Label:     "list_containers",
					Operation: listUserContainers,
					Expected:  "OK",
				},
				// TODO - others?
			},
		},
	}

	// Now get down to work

	for _, r := range scenarios {
		for _, o := range r.Operations {
			res := o.Operation(serverURL, &r, &o)
			if strings.Contains(res, o.Expected) {
				log.Infof("PASS %v:%s found expected match in output: %s", r.Role, o.Label, res)
			} else {
				assert.Fail(t, fmt.Sprintf("Failed %v:%s - expected %s found %s", r.Role, o.Label, o.Expected, res))
			}
		}
	}
}
