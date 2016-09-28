// Package uninstall implements the high-level uninstallation flow for orca
package uninstall

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// Run the uninstall flow
func Run(c *cli.Context) {
	if code, err := uninstall(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}

func uninstall(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}

	if !config.InPhase2 {
		return runPhase1(c, ec)
	}

	// We're in phase 2
	return runPhase2(c, ec)
}

func runPhase1(c *cli.Context, ec *client.EngineClient) (int, error) {
	instanceID, err := getInstanceID(c, ec)
	if err != nil {
		return 1, err
	}

	config.OrcaInstanceID = instanceID

	// Mount the volume(s) we'll need to unwind an HA node if detected
	config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
		fmt.Sprintf("%s:%s", config.SwarmKvCertVolumeName, config.SwarmKvCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmControllerCertVolumeName, config.SwarmControllerCertVolumeMount),
		fmt.Sprintf("%s:%s", config.EngineConfigDir, config.EngineConfigDir),
	)

	ret, err := ec.StartPhase2(os.Args[1:], false)
	if err != nil {
		return ret, err
	}

	// Remove the volumes in phase1 so we don't have any mounted.
	log.Info("Removing UCP volumes")

	// We've seen some raciness where the containers exited, but the engine
	// would fail the volume removal claiming they were still in use,
	// however pausing for a moment gives the engine enough time to get
	// consistent and allow the removal.
	time.Sleep(2 * time.Second)

	if err := removeVolumes(c, ec); err != nil {
		return 1, err
	}

	return 0, nil
}

// The uninstall may be run with the `--id` option to specify which instance of
// UCP to uninstall. It may also be run with the `--interactive` flag to prompt
// for confirmation before uninstalling. If neither are specified, then we
// cannot continue.
func getInstanceID(c *cli.Context, ec *client.EngineClient) (string, error) {
	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)

	// The UCP containers will each have a label which identifies the
	// instance of UCP. This gets the set of distnict instance IDs.
	// Normally there will only be a single ID.
	ids := client.GetInstanceIDs(containers)
	if len(ids) == 0 {
		return "", fmt.Errorf("No running UCP instances detected on this engine")
	}
	if len(ids) > 1 {
		log.Warnf("Multiple UCP instances detected: %v", ids)
	}

	id := c.String("id")
	if id != "" {
		// Use the one specified on the command line.
		return id, nil
	}

	// Normally the list of IDs will be a single item, but we're not yet
	// optimizing for multiple UCP instances on the same node so just use
	// the first one.
	id = ids[0]

	if !c.Bool("interactive") {
		// The user did not specify a UCP instance ID to uninstall and
		// did not specify to be prompted.
		log.Infof("We detected local components of UCP instance %s", id)
		return "", fmt.Errorf(`Re-run the command with "--id %s" or --interactive to confirm you want to remove this UCP instance.`, id)
	}

	// Propmt the user to confirm the uninstall.
	log.Infof("We're about to uninstall the local components for UCP ID: %s", id)
	fmt.Printf("Do you want proceed with the uninstall? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		log.Debugf("Failed to read input: %s", err)
		return "", err
	}

	value = strings.TrimSpace(strings.ToLower(value))
	if value == "y" || value == "yes" {
		return id, nil
	}

	return "", fmt.Errorf("Not uninstalling per user request")
}

// removeVolumes actually *deletes* the volume from Docker, unlike the method
// in the utils package (CleanupVolumes) which just removes volume contents.
func removeVolumes(c *cli.Context, ec *client.EngineClient) error {
	cleanup := []string{
		// Always remove data volumes.
		config.OrcaKVVolumeName,
		config.AuthStoreDataVolumeName,
		config.AuthWorkerDataVolumeName,
	}

	// Do not remove certificate volumes if asked to preserve them.
	if !c.Bool("preserve-certs") {
		cleanup = append(cleanup,
			config.OrcaRootCAVolumeName,
			config.OrcaServerCertVolumeName,
			config.SwarmControllerCertVolumeName,
			config.SwarmKvCertVolumeName,
			config.SwarmNodeCertVolumeName,
			config.SwarmRootCAVolumeName,
			config.AuthStoreCertsVolumeName,
			config.AuthAPICertsVolumeName,
			config.AuthWorkerCertsVolumeName,
		)
	}

	var err error
	for _, volume := range cleanup {
		if !ec.VolumeExists(volume) {
			continue
		}

		if lastError := ec.RemoveVolume(volume); lastError != nil {
			err = lastError
		}
	}

	return err
}

func runPhase2(c *cli.Context, ec *client.EngineClient) (int, error) {
	if err := removeServices(ec); err != nil {
		return 1, err
	}

	// Give it a moment to remove the service's container
	time.Sleep(2 * time.Second)

	if err := removeContainers(ec); err != nil {
		return 1, err
	}

	if !c.Bool("preserve-images") {
		log.Info("Removing UCP images")
		if err := ec.RemoveImages(); err != nil {
			return 1, err
		}
	}

	return 0, nil
}

func removeServices(ec *client.EngineClient) error {
	services, err := ec.FindServiceIDsByOrcaInstanceID(config.OrcaInstanceID)
	if err != nil {
		return err
	}

	var anyErrors error
	log.Info("Removing UCP Services")
	for _, serviceID := range services {
		if err := ec.ServiceRemove(serviceID); err != nil {
			log.Error(err)
			anyErrors = err
		}
	}

	if anyErrors != nil {
		return fmt.Errorf("Failed to remove an existing UCP service: %s", anyErrors)
	}
	return nil
}

func removeContainers(ec *client.EngineClient) error {
	containerIDs, err := ec.FindContainerIDsByOrcaInstanceID(config.OrcaInstanceID)
	if err != nil {
		log.Debug("Failed to find specified UCP instances")
		return err
	}

	if containerIDs == nil {
		log.Infof("No matching UCP containers detected for ID: %s", config.OrcaInstanceID)
		return nil
	}

	// Detect if this is an HA node, regular capacity node, or the last
	// controller.
	if kvStoreURL, err := ec.FindKv(); err == nil {
		if err := ec.ScaleDownAuthStore(kvStoreURL); err != nil {
			log.Errorf("unable to scale down auth storage cluster: %s", err)
			// Proceed anyway?
		}

		if err := config.RemoveEnziProviderConfig(kvStoreURL); err != nil {
			log.Errorf("unable to remove auth provider config: %s", err)
			// Proceed anyway though, since we may be doing a full uninstall.
		}

		// Dig a little deeper to see if we're in HA mode
		if err := ec.DetectAndRemoveKVHANode(kvStoreURL); err != nil {
			log.Debugf("KV removal error: %s", err)
			log.Warningf("Unable to remove local node from HA cluster.  If you are not planning to fully shutdown this cluster, you should add a new replica before uninstalling any other replica nodes.")
			// Proceed anyway though, since we may be doing a full uninstall
		}
	} // Not a controller node

	var anyErrors error
	log.Info("Removing UCP Containers")
	for _, containerID := range containerIDs {
		if err := ec.RemoveContainerByID(containerID, 5, true); err != nil {
			anyErrors = err
		}
	}
	if anyErrors != nil {
		return fmt.Errorf("Failed to remove an existing UCP container %s", anyErrors)
	}

	return nil
}
