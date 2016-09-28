package client

import (
	"strings"
	"testing"

	"github.com/docker/orca/bootstrap/config"
)

func TestCheckSocketDir(t *testing.T) {
	config.DockerSock = "unix:///tmp/"
	c := EngineClient{}
	if err := c.CheckSocket(); !strings.Contains(err.Error(), "is not a socket") {
		t.Error("Didn't detect a directory")
	}
}

func TestCheckMissingFile(t *testing.T) {
	config.DockerSock = "unix:///foobar/baz/bif"
	c := EngineClient{}
	if err := c.CheckSocket(); !strings.Contains(err.Error(), "Missing docker.sock.") {
		t.Error("Didn't handle missing file correctly")
	}
}
