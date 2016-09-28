package agent

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"

	"github.com/docker/orca/agent/agent/components"
	"github.com/docker/orca/agent/agent/utils"
	"github.com/docker/orca/config"
	"github.com/docker/orca/config/discover"
	orcatypes "github.com/docker/orca/types"
)

// StartAgent is invoked for the `agent` command
func StartAgent(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	// Check if the docker socket is actually bound to the container
	dockerSocket := c.String("d")
	if _, err := os.Stat(dockerSocket); os.IsNotExist(err) {
		return fmt.Errorf("Unable to locate the Docker Socket at %s", dockerSocket)
	}

	// Extract the expected Cluster Config from Environment Variables
	clcfg := orcatypes.ClusterConfig{
		ControllerPort:     c.String("controller_port"),
		SwarmPort:          c.String("swarm_port"),
		TrustedRegistryCAs: c.StringSlice("trusted_registry_cas"),
		ImageVersion:       c.String("image_version"),
		Secret:             c.String("secret"),
		UCPInstanceID:      c.String("ucp_instance_id"),
		DNS:                c.StringSlice("dns"),
		DNSOpt:             c.StringSlice("dns_opt"),
		DNSSearch:          c.StringSlice("dns_search"),
	}

	// Create an engine-api client to the local engine
	dclient, err := client.NewClient(fmt.Sprintf("unix://%s", dockerSocket), "", nil, nil)
	if err != nil {
		return fmt.Errorf("could not create a docker engine-api client: %s", err)
	}
	log.Infof("Initialized engine-api client, version %s", dclient.ClientVersion())

	// Main reconciliation loop
	for {
		// Get the expected config from the agent's env and from engine-api
		log.Info("Inspecting expected and current node config")
		expectedConfig, err := config.GetExpectedNodeConfig(dclient, clcfg)
		if err != nil {
			return err
		}
		log.Debugf("Expected Node Config: %#v", expectedConfig)

		// Extract the current config from the UCP components
		currentConfig, err := discover.DiscoverNodeConfig(dclient)
		if err != nil {
			return err
		}
		log.Debugf("Current Node Config: %#v", currentConfig)

		// Perform state reconciliation if there is a mismatch in the two configs
		differs, err := stateDiffers(expectedConfig, currentConfig)
		if err != nil {
			return err
		}

		if differs {
			log.Info("Beginning State Reconciliation")
			// TODO: in a fresh install, fail hard and do not reconcile
			log.Error(reconcileState(dclient, expectedConfig, currentConfig))
		} else {
			log.Info("Expected and Current states match, no changes required")
		}

		time.Sleep(10 * time.Second) // TODO: fine-tune the reconciliation interval
	}
}

func stateDiffers(expected *orcatypes.NodeConfig, current *orcatypes.NodeConfig) (bool, error) {
	for _, component := range components.ComponentList {
		differs, err := component.RequiresReconciliation(expected, current)
		if err != nil {
			return false, err
		}
		if differs {
			log.Debugf("Component %#v requires reconciling", component)
			return true, nil
		}
	}
	return false, nil
}

// reconcileState creates the ucp-reconcile container and blocks until it is terminated
func reconcileState(dclient *client.Client, expected *orcatypes.NodeConfig, current *orcatypes.NodeConfig) error {
	// Look for an existing `ucp-reconcile` container
	_, err := dclient.ContainerInspect(context.TODO(), config.OrcaReconcileContainerName)
	if err == nil {
		// The `ucp-reconcile` container already exists, force remove it
		log.Infof("Detected existing %s container, removing it", config.OrcaReconcileContainerName)
		err2 := dclient.ContainerRemove(context.TODO(), config.OrcaReconcileContainerName, types.ContainerRemoveOptions{
			Force: true,
		})
		if err2 != nil {
			return fmt.Errorf("Unable to remove existing %s container: %s", config.OrcaReconcileContainerName, err2)
		}
	}

	// Determine the exact image to use for Reconcile from the imageVersion
	config.ImageVersion = expected.ClusterConfig.ImageVersion
	imageName, err := config.GetContainerImage(config.OrcaReconcileContainerName)
	if err != nil {
		return err
	}

	payload, err := utils.SerializeReconcileArgs(expected, current)
	if err != nil {
		return err
	}

	// Launch the Reconcile container for the current node
	log.Infof("Launching %s container with image %s", config.OrcaReconcileContainerName, imageName)
	containerID, err := launchReconcile(dclient, imageName, payload)
	if err != nil {
		return fmt.Errorf("unable to launch %s container: %s", config.OrcaReconcileContainerName, err)
	}

	errChan := make(chan error)
	// Goroutine: Wait for the ucp-reconcile container to exit
	go func(c chan<- error) {
		status, err := dclient.ContainerWait(context.TODO(), containerID)
		if err != nil {
			c <- err
		}
		if status != 1 {
			c <- fmt.Errorf("%s container exited with status code: %d", config.OrcaReconcileContainerName, status)
		}
		c <- nil
	}(errChan)

	// Block until ucp-reconcile is finished
	return <-errChan
}
