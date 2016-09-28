package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/orca/bootstrap/config"
)

func TestStartKvHappy(t *testing.T) {
	t.Skip("TODO: fix this test; uses a mock that fails when trying to connect confirm service is running")

	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if err := c.StartKv("", &config.KVDeployCfg{}); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestStartKvFailCreate(t *testing.T) {
	expected := "Foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartKv("", &config.KVDeployCfg{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartKvFailStart(t *testing.T) {
	expected := "Foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartKv("", &config.KVDeployCfg{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
