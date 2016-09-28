package restart

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
)

func restart(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)
	ec, err := client.New()
	if err != nil {
		return 1, err
	}

	instanceID, err := ec.FindOrcaInstanceID()
	if err != nil {
		return 1, err
	}

	containerIDs, err := ec.FindContainerIDsByOrcaInstanceID(instanceID)
	if err != nil {
		log.Debug("Failed to find specified UCP instances")
		return 1, err
	}

	if containerIDs == nil {
		log.Infof("No matching UCP containers detected for ID: %s", instanceID)
		return 1, nil
	}

	log.Debugf("Found UCP ID %s, containers %s", instanceID, containerIDs)

	good := true
	for _, containerID := range containerIDs {
		log.Debugf("Starting %s", containerID)
		err := ec.ContainerRestart(containerID)
		if err != nil {
			log.Warnf("Unable to start container %s: %s", containerID, err)
			good = false
		}
	}

	if good {
		return 0, nil
	} else {
		return 1, fmt.Errorf("One or more containers failed to start.")
	}
}

// Run the installation flow
func Run(c *cli.Context) {
	if code, err := restart(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
