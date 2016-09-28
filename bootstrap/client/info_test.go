package client

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
)

func TestGetLocalIDHappy(t *testing.T) {
	expected := "bogusname"
	shim := testShim{
		infoOut: types.Info{
			Name: expected,
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if id, err := c.GetLocalID(); err != nil {
		t.Error("Didn't pass")
	} else if id != expected {
		t.Errorf("Wrong ID (expected %s): %s", expected, id)
	}
}
func TestGetLocalIDError(t *testing.T) {
	// Clear the environment
	os.Setenv("UCP_LOCAL_NAME", "")
	expected := "foo"
	shim := testShim{
		infoErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if output, err := c.GetLocalID(); err == nil {
		t.Error("Didn't fail: %s", output)
	} else if !strings.Contains(err.Error(), expected) {
		t.Error("Wrong error returned (expected %s): %s", expected, err)
	}
}

func TestGetHostnamesHappy(t *testing.T) {
	shim := testShim{
		infoOut: types.Info{
			Name: "bogusname",
		},
	}
	c := EngineClient{}
	c.bootstrapper = &types.ContainerJSON{
		NetworkSettings: &types.NetworkSettings{
			DefaultNetworkSettings: types.DefaultNetworkSettings{
				Gateway: "1.1.1.1",
			},
		},
	}
	c.client = clientShim(shim)
	if res, err := c.GetHostnames(); err != nil {
		t.Error("Didn't pass")
	} else if len(res) != 3 {
		t.Errorf("Wrong number of results (expected 3): %d %v", len(res), res)
	}
}

func TestGetHostnamesError(t *testing.T) {
	expected := "foo"
	shim := testShim{
		infoErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if _, err := c.GetHostnames(); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (shoudl contain %s): %s", expected, err)
	}
}
