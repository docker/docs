package client

import (
	"bytes"
	"errors"
	"github.com/docker/engine-api/types"
	"os"
	"strings"
	"testing"
)

func TestGetHostAddressHappyPath(t *testing.T) {
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		logsOut: &closingBuffer{bytes.NewBufferString(`8.8.8.8 via 192.168.104.1 dev eth0  src 192.168.104.65 
    cache  users 3 age 719sec`)},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	addr, err := c.GetHostAddress()
	if err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
	if addr != "192.168.104.65" {
		t.Errorf("Wrong address (expected 192.168.104.65): %s", addr)
	}
}

func TestGetHostAddressEnvBypass(t *testing.T) {
	shim := testShim{
		createErr: errors.New("Foo"),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	expectedAddr := "1.2.3.4"
	os.Setenv("UCP_HOST_ADDRESS", expectedAddr)
	addr, err := c.GetHostAddress()
	os.Unsetenv("UCP_HOST_ADDRESS")
	if err != nil {
		t.Errorf("Didn't pass: %s", err)
	}
	if addr != expectedAddr {
		t.Errorf("Wrong address (expected %s): %s", expectedAddr, addr)
	}
}

func TestGetHostAddressBadCreate(t *testing.T) {
	expected := "Failed to create "
	shim := testShim{
		createErr: errors.New("Foo"),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if _, err := c.GetHostAddress(); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestGetHostAddressBadStart(t *testing.T) {
	expected := "Foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		startErr:  errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if _, err := c.GetHostAddress(); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestGetHostAddressBadLogs(t *testing.T) {
	expected := "Foo"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		logsErr:   errors.New(expected),
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if _, err := c.GetHostAddress(); err == nil {
		t.Error("Didn't fail as expected")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestGetHostAddressBadOutput(t *testing.T) {
	expected := "Failed to determine"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		logsOut:   &closingBuffer{bytes.NewBufferString("Some bad output")},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if _, err := c.GetHostAddress(); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}

func TestGetHostAddressBadOutputWithNewline(t *testing.T) {
	expected := "Failed to determine"
	shim := testShim{
		createOut: types.ContainerCreateResponse{ID: "123"},
		logsOut:   &closingBuffer{bytes.NewBufferString("Some bad \noutput")},
	}
	c := EngineClient{}
	c.client = clientShim(shim)
	c.bootstrapper = &types.ContainerJSON{}
	if _, err := c.GetHostAddress(); err == nil {
		t.Error("Didn't fail")
	} else if !strings.Contains(err.Error(), expected) {
		t.Errorf("Wrong error returned (should contain %s): %s", expected, err)
	}
}
