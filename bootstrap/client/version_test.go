package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/orca/bootstrap/config"
)

func TestSameVersion(t *testing.T) {
	config.MinVersion = "1.5.0"

	shim := testShim{
		versionOut: types.Version{
			Version: config.MinVersion,
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.CheckDockerVersion(); err != nil {
		t.Errorf("Didn't pass matching version: %s", err)
	}
}

func TestGoodVersion(t *testing.T) {
	config.MinVersion = "1.5.0"

	shim := testShim{
		versionOut: types.Version{
			Version: "1.6.0",
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.CheckDockerVersion(); err != nil {
		t.Errorf("Didn't pass matching version: %s", err)
	}
}

func TestTooOld(t *testing.T) {
	expected := "is too old"
	config.MinVersion = "1.5.0"

	shim := testShim{
		versionOut: types.Version{
			Version: "1.4.0",
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.CheckDockerVersion(); err == nil {
		t.Error("Passed too old version")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Invalid error for too old (should contain %s): %s", expected, err)
	}
}

func TestVersionFail(t *testing.T) {
	expected := "foo"
	shim := testShim{
		versionErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.CheckDockerVersion(); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestCSVersion(t *testing.T) {
	config.MinVersion = "1.9.0"

	shim := testShim{
		versionOut: types.Version{
			Version: "1.9.0-cs1",
		},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.CheckDockerVersion(); err != nil {
		t.Errorf("Didn't pass matching version: %s", err)
	}
}
