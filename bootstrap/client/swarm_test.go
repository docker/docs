package client

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
)

func TestStartProxyHappy(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if err := c.StartProxy(); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestStartProxyFailCreate(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartProxy(); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartProxyFailStart(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartProxy(); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartSwarmManagerHappy(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if err := c.StartSwarmManager("", "", []string{}); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestStartSwarmManagerFailCreate(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartSwarmManager("", "", []string{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartSwarmManagerFailStart(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartSwarmManager("", "", []string{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartSwarmJoinHappy(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if err := c.StartSwarmJoin("", "", []string{}); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestStartSwarmJoinFailCreate(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createErr: errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartSwarmJoin("", "", []string{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestStartSwarmJoinFailStart(t *testing.T) {
	expected := "foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.StartSwarmJoin("", "", []string{}); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestFilterSwarmArgsNoFilter(t *testing.T) {
	inputArgs := []string{"a", "b", "c"}
	outputArgs := FilterSwarmArgs(inputArgs)
	if !reflect.DeepEqual(inputArgs, outputArgs) {
		t.Errorf("Expected: %v, got %v", inputArgs, outputArgs)
	}
}

func TestFilterSwarmArgsAdvertiseEqual(t *testing.T) {
	inputArgs := []string{"a", "--advertise=foo", "c"}
	expected := []string{"a", "c"}
	outputArgs := FilterSwarmArgs(inputArgs)
	if !reflect.DeepEqual(expected, outputArgs) {
		t.Errorf("Expected: %v, got %v", expected, outputArgs)
	}
}

func TestFilterSwarmArgsAdvertise(t *testing.T) {
	inputArgs := []string{"a", "--advertise", "foo", "c"}
	expected := []string{"a", "c"}
	outputArgs := FilterSwarmArgs(inputArgs)
	if !reflect.DeepEqual(expected, outputArgs) {
		t.Errorf("Expected: %v, got %v", expected, outputArgs)
	}
}

func TestFilterSwarmArgsAdvertiseManageJoin(t *testing.T) {
	inputArgs := []string{"a", "join", "c", "manage"}
	expected := []string{"a", "c"}
	outputArgs := FilterSwarmArgs(inputArgs)
	if !reflect.DeepEqual(expected, outputArgs) {
		t.Errorf("Expected: %v, got %v", expected, outputArgs)
	}
}
