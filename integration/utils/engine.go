package utils

// Someday this might go away, but for now, we need both docker clients
// in our integration tests so we don't have to re-write the world.

import (
	"fmt"

	"github.com/docker/engine-api/client"
	"github.com/samalba/dockerclient"
)

func ConvertToEngineAPI(dclient *dockerclient.DockerClient) (*client.Client, error) {
	return client.NewClient(fmt.Sprintf("tcp://%s", dclient.URL.Host), "", dclient.HTTPClient, nil)
}
