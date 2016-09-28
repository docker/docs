package client

import (
	"errors"
	"testing"

	"github.com/docker/engine-api/types"
)

func TestFindPresent(t *testing.T) {
	shim := testShim{
		inspectContainerOut: types.ContainerJSON{},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if res := c.FindContainers([]string{"foo"}); len(res) != 1 {
		t.Errorf("Didn't pass: %v", res)
	}
}

func TestFindMissing(t *testing.T) {
	shim := testShim{
		inspectContainerErr: errors.New("Missing"),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if res := c.FindContainers([]string{"foo"}); len(res) != 0 {
		t.Errorf("Didn't pass: %v", res)
	}
}

func TestFindContainerIDsByOrcaInstanceIDHappy(t *testing.T) {
	expected := "id123"
	instanceID := "123foo"
	shim := testShim{
		listOut: []types.Container{
			types.Container{
				Names: []string{"bar"},
				ID:    expected,
				Labels: map[string]string{
					"com.docker.ucp.InstanceID": instanceID,
				},
			},
			types.Container{
				Names: []string{"bar"},
				ID:    "id3",
				Labels: map[string]string{
					"com.docker.ucp.InstanceID": "someotherid",
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
	res, err := c.FindContainerIDsByOrcaInstanceID(instanceID)
	if err != nil {
		t.Errorf("Didn't pass: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Wrong number of results (expected 1): %d", len(res))
	}
	if res[0] != expected {
		t.Errorf("Wrong ID returned (expected %s): %s", expected, res[0])
	}
}

func TestFindContainerIDsByOrcaInstanceIDListFail(t *testing.T) {
	instanceID := "123foo"
	shim := testShim{
		listErr: errors.New("foo"),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if _, err := c.FindContainerIDsByOrcaInstanceID(instanceID); err == nil {
		t.Error("Didn't fail as expected")
	}
}
