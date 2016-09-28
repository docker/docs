package client

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/bootstrap/config"
)

var (
	ErrMissingDockerSock = errors.New(
		"Missing docker.sock. You must run the bootstrap container with \"-v /var/run/docker.sock:/var/run/docker.sock\"")
	ErrDockerSockNotSock = errors.New(
		"docker.sock is not a socket. You must run with \"-v /var/run/docker.sock:/var/run/docker.sock\"")
)

// Verify that the local docker.sock looks valid
func (c *EngineClient) CheckSocket() error {
	log.Debug("Verifying docker.sock")
	fi, err := os.Stat(config.DockerSock[7:])
	if err != nil {
		log.Debugf("Failed to stat: %s", err)
		return ErrMissingDockerSock
	}
	mode := fi.Mode()
	if mode&os.ModeSocket != os.ModeSocket {
		log.Debugf("docker.sock is %v", mode)
		return ErrDockerSockNotSock
	}
	return nil
}
