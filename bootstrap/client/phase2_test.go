package client

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func TestStartPhase2HappyPath(t *testing.T) {
	expected := 456
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		logsOut:   &closingBuffer{bytes.NewBufferString("\n")},
		waitOut:   expected,
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
	if ret, err := c.StartPhase2([]string{""}, false); err != nil {
		t.Errorf("Didn't pass: %s", err)
	} else if ret != expected {
		t.Errorf("Wrong exit status (expected %d): %d", expected, ret)
	}
}

func TestStartPhase2CreateFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
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
	if _, err := c.StartPhase2([]string{""}, false); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartPhase2StartFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
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
	if _, err := c.StartPhase2([]string{""}, false); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartPhase2AttachFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		attachErr: errors.New(expected),
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
	if _, err := c.StartPhase2([]string{""}, false); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
