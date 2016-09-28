package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func TestRemoveError(t *testing.T) {
	msg := "XXX 123"
	shim := testShim{
		removeErr: errors.New(msg),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.RemoveContainers([]*types.ContainerJSON{&types.ContainerJSON{ContainerJSONBase: &types.ContainerJSONBase{ID: "1234"}}}); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), msg) {
		t.Errorf("Didn't fail with expected message (expected %s): %s", msg, err)
	}
}

func TestStopError(t *testing.T) {
	msg := "XXX 456"
	shim := testShim{
		stopErr: errors.New(msg),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	// Expect stop failure to be ignored
	if err := c.RemoveContainers([]*types.ContainerJSON{&types.ContainerJSON{ContainerJSONBase: &types.ContainerJSONBase{ID: "1234"}}}); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestRemoveOrcaContainersHappySimple(t *testing.T) {
	shim := testShim{}
	c := EngineClient{}
	c.client = clientShim(shim)
	input := []*types.ContainerJSON{
		&types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: "1234",
			},
		},
	}
	if err := c.RemoveOrcaContainers(input); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestRemoveOrcaContainersHappy(t *testing.T) {
	instanceID := "instance123"
	shim := testShim{
		listOut: []types.Container{
			types.Container{
				Names: []string{"bar"},
				ID:    "id2",
				Labels: map[string]string{
					"com.docker.orca.InstanceID": instanceID,
				},
			},
			types.Container{
				Names: []string{"bar"},
				ID:    "id3",
				Labels: map[string]string{
					"com.docker.orca.InstanceID": "someotherid",
				},
			},
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			Image: "docker/ucp",
		},
	}
	input := []*types.ContainerJSON{
		&types.ContainerJSON{
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.orca.InstanceID": instanceID,
				},
			},
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: "id1",
			},
		},
	}
	if err := c.RemoveOrcaContainers(input); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestRemoveOrcaContainersHappyStopFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		stopErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	input := []*types.ContainerJSON{
		&types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: "1234",
			},
		},
	}
	if err := c.RemoveOrcaContainers(input); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Didn't fail with expected message (expected %s): %s", expected, err)
	}
}

func TestRemoveOrcaContainersHappyRemoveFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		removeErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	input := []*types.ContainerJSON{
		&types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{
				ID: "1234",
			},
		},
	}
	if err := c.RemoveOrcaContainers(input); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Didn't fail with expected message (expected %s): %s", expected, err)
	}
}
