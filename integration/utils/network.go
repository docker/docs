package utils

import (
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

func TestOverlayNetwork(t *testing.T, client dockerclient.Client) {
	// Now create an overlay network to make sure everything is happy on the back-end
	log.Infof("Attempting to create a network")

	// HACK - try to find a better way for this...

	var err error
	for i := 0; i < 60; i++ {
		resp, err := client.CreateNetwork(&dockerclient.NetworkCreate{
			Name:   "orca-test-net",
			Driver: "overlay",
		})

		if err == nil {
			log.Infof("Network create: %s - %s", resp.ID, resp.Warning)
			return // we're good now
		} else if strings.Contains(err.Error(), "No healthy node available in the cluster") {
			log.Info("Swarm complaining about no healthy nodes, retrying in 2 secs...")
		} else {
			log.Warn(err) // XXX Short circuit and die?
		}
		time.Sleep(2 * time.Second)
	}
	// failed if we got here
	require.Nil(t, err)
}
