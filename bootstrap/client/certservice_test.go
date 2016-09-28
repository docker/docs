package client

import (
	"errors"
	"github.com/docker/engine-api/types"
	"github.com/docker/orca/config"
	"strings"
	"testing"
)

// TODO - Happy path requires mocking/stubbing/intercepting the WaitFor* routines somehow

func TestStartCertServiceFailCreate(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartCA(config.OrcaSwarmCAContainerName, "", 1234); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartCertServiceFailStart(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartCA(config.OrcaSwarmCAContainerName, "", 1234); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
