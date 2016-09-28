package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
)

func TestStartOrcaServerHappy(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if err := c.StartOrcaServer("", "", ""); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestStartOrcaServerFailCreate(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartOrcaServer("", "", ""); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartOrcaServerFailStart(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartOrcaServer("", "", ""); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
