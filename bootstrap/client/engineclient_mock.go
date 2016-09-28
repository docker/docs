package client

import (
	log "github.com/Sirupsen/logrus"
	dockerclient "github.com/docker/engine-api/client"

	"github.com/docker/orca/bootstrap/config"
)

func Mock() (*EngineClient, error) {
	ret := &EngineClient{}

	log.Debugf("Connecting to docker %s", config.DockerSock)
	version := ""
	docker, err := dockerclient.NewClient(config.DockerSock, version, nil, nil)
	if err != nil {
		log.Error("Unable to connect to the docker engine")
		return nil, err
	}

	shim := realShim{client: docker}
	ret.client = clientShim(shim)

	// ret.bootstrapper is nil in the engineclient mock
	return ret, nil
}
