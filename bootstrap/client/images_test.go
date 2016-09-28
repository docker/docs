package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/orca/bootstrap/config"
)

func TestAlreadyHave(t *testing.T) {
	shim := testShim{
		inspectImageOut: types.ImageInspect{},
		pullOut:         ioutil.NopCloser(bytes.NewBufferString("")),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.VerifyOrPullImages(false); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}

func TestGoodPull(t *testing.T) {
	shim := testShim{
		inspectImageErr: errors.New("Not found"),
		pullOut:         ioutil.NopCloser(bytes.NewBufferString("")),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.VerifyOrPullImages(false); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}
func TestBadPull(t *testing.T) {
	pullErr := errors.New("Failed to pull")
	shim := testShim{
		inspectImageErr: errors.New("Not found"),
		pullErr:         pullErr,
		pullOut:         ioutil.NopCloser(bytes.NewBufferString("")),
	}
	config.PullBehavior = "missing"
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.VerifyOrPullImages(false); err != pullErr {
		t.Errorf("Didn't fail as expected (expected %s): %s", pullErr, err)
	}
}

func TestGoodRemove(t *testing.T) {
	shim := testShim{}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.RemoveImages(); err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
}
func TestBadRemove(t *testing.T) {
	expected := "One or more Orca images were still"
	shim := testShim{
		removeImageErr: errors.New("Failed to remove"),
		pullOut:        ioutil.NopCloser(bytes.NewBufferString("")),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	if err := c.RemoveImages(); err == nil {
		t.Errorf("Expected a failure, but it succeeded")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Didn't fail as expected (expected %s): %s", expected, err)
	}
}
