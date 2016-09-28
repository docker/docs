package client

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/orca/bootstrap/config"
)

var (
	p1    = 1
	p2    = 2
	ports = []*int{
		&p1,
		&p2,
	}
)

func TestFindAvailable(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		attachOut: types.HijackedResponse{
			Reader: bufio.NewReader(bytes.NewBufferString("")),
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Image:      "",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{
			Env: []string{""},
		},
	}
	config.DockerTimeout = 0
	// We'll expect a connect failure because we aren't actually running the container
	err := c.CheckPorts(ports)
	if !strings.Contains(err.Error(), "following required ports are blocked") {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestCreateFailure(t *testing.T) {
	message := "XXX Create failed XXX"
	createErr := errors.New(message)
	shim := testShim{
		createErr: createErr,
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Image:      "",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{
			Env: []string{""},
		},
	}
	if err := c.CheckPorts(ports); err == nil {
		t.Error("Didn't report create failure")
	} else if !strings.Contains(err.Error(), "ports are already in use") {
		t.Errorf("Wrong error returned (should contain %s): %s", message, err)
	}
}

func TestInUse(t *testing.T) {
	expected := "The following required ports"
	startErr := errors.New("error message")
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  startErr,
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Image:      "",
			HostConfig: &container.HostConfig{},
		},
		Config: &container.Config{
			Env: []string{""},
		},
	}
	if err := c.CheckPorts(ports); err == nil {
		t.Error("Didn't report start failure")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
