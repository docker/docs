// Package regen_certs will regenerate the local certs for this node
package regen_certs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/utils"
	"golang.org/x/net/context"
)

var (
	BufSize = 1024 * 10 * 10
)

func isSameRoot() (bool, error) {
	rootCA, err := ioutil.ReadFile(filepath.Join(config.SwarmRootCAVolumeMount, config.CertFilename))
	if err != nil {
		log.Error("Failed to load cluster CA")
		return false, err
	}
	expectedRootCA, err := ioutil.ReadFile(filepath.Join(config.SwarmNodeCertVolumeMount, config.CAFilename))
	if err != nil {
		log.Error("Failed to load cluster CA")
		return false, err
	}
	sameRoot := strings.TrimSpace(string(rootCA)) == strings.TrimSpace(string(expectedRootCA))
	log.Debugf("Same root CAs: %v", sameRoot)
	return sameRoot, nil
}

func recertify(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	interactive := c.Bool("interactive")

	ec, err := client.New()
	if err != nil {
		return 1, err
	}

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	rootCAOnly := c.Bool("root-ca-only")

	nodeType := ec.DetectNodeType()
	if nodeType != client.Controller {
		return 1, fmt.Errorf("This node is not a controller.  To regenerate certs for a non-controller cluster member, run `join --fresh-install`")
	}

	if !config.InPhase2 {
		ids := client.GetInstanceIDs(containers)
		if len(ids) == 0 {
			return 1, fmt.Errorf("No running UCP instances detected on this engine")
		} else if len(ids) > 1 {
			log.Warnf("Multiple UCP instances detected: %v", ids)
		}
		msg := ""
		// Add the relevant volume mounts we need
		if rootCAOnly {
			msg = "regenerate the local UCP Root CAs"
			config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
				fmt.Sprintf("%s:%s", config.OrcaRootCAVolumeName, config.OrcaRootCAVolumeMount),
				fmt.Sprintf("%s:%s", config.SwarmRootCAVolumeName, config.SwarmRootCAVolumeMount),
			)
		} else {
			msg = "regenerate the controller certs"
			config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
				fmt.Sprintf("%s:%s", config.OrcaRootCAVolumeName, config.OrcaRootCAVolumeMount),
				fmt.Sprintf("%s:%s", config.SwarmRootCAVolumeName, config.SwarmRootCAVolumeMount),
				fmt.Sprintf("%s:%s", config.SwarmNodeCertVolumeName, config.SwarmNodeCertVolumeMount),
				fmt.Sprintf("%s:%s", config.SwarmKvCertVolumeName, config.SwarmKvCertVolumeMount),
				fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, config.OrcaServerCertVolumeMount),
				fmt.Sprintf("%s:%s", config.SwarmControllerCertVolumeName, config.SwarmControllerCertVolumeMount),
				// Auth Service Related Volume Mounts.
				fmt.Sprintf("%s:%s", config.AuthStoreCertsVolumeName, config.AuthStoreCertsVolumeMount),
				fmt.Sprintf("%s:%s", config.AuthAPICertsVolumeName, config.AuthAPICertsVolumeMount),
				fmt.Sprintf("%s:%s", config.AuthWorkerCertsVolumeName, config.AuthWorkerCertsVolumeMount),
			)
		}

		id := c.String("id")
		if id != "" {
			config.OrcaInstanceID = id
		} else {
			// Normally this will be a single item, and we're not optimizing for multiple orca's
			id = ids[0]
			if interactive {
				log.Warnf("We're about to %s for UCP ID: %s", msg, id)
				log.Warn("Regenerating certs is a multi-step process and will involve some downtime of the cluster.")
				if rootCAOnly {
					log.Warn("As step (1), you should run this command on ONE CONTROLLER ONLY in your cluster")
					log.Warn("After completing the process, users will need to download new cert bundles.")
				} else {
					log.Warn("If you intend to regenerate the Root CAs as well, and have not already done so, you should cancel and re-run the command with '--root-ca-only' on one controller in the cluster.")
				}
				fmt.Printf("Do you want proceed with regenerating certs? (y/n): ")

				reader := bufio.NewReader(os.Stdin)
				value, err := reader.ReadString('\n')
				if err != nil {
					log.Debugf("Failed to read input: %s", err)
					return 1, err
				}
				value = strings.TrimSpace(strings.ToLower(value))
				if value == "y" || value == "yes" {
					config.OrcaInstanceID = id
				} else {
					return 1, fmt.Errorf("Not regenerating certs per user request")
				}

			} else {
				log.Infof("We detected local components of UCP instance %s", id)
				return 1, fmt.Errorf(`Re-run the command with "--id %s" or --interactive to confirm you want to regenerate certs for this UCP instance.`, id)
			}
		}
		return ec.StartPhase2(os.Args[1:], false)
	}

	containerIDs, err := ec.FindContainerIDsByOrcaInstanceID(config.OrcaInstanceID)
	if err != nil {
		log.Debug("Failed to find specified UCP instances")
		return 1, err
	}

	// Need to get list of controllers from UCP
	// TODO - Adrian's changing this to be structured instead of list of strings
	controllers, err := ec.AutodetectControllers()
	if err != nil {
		return 1, fmt.Errorf("We were unable to locate the UCP controllers and can not proceed: %s", err)
	}
	haDetected := false
	if len(controllers) > 1 { // HA deployment
		haDetected = true
	}
	if rootCAOnly {
		// Reduce the set of containers to just the CA containers
		caContainers := []string{}
		for _, id := range containerIDs {
			info, err := ec.InspectContainer(id)
			if err != nil {
				return 1, fmt.Errorf("Failed to lookup container %s name", id)
			}
			if strings.Contains(info.Name, orcaconfig.OrcaCAContainerName) || strings.Contains(info.Name, orcaconfig.OrcaSwarmCAContainerName) {
				caContainers = append(caContainers, id)
			}
		}
		if len(caContainers) == 0 {
			return 1, fmt.Errorf("This system does not appear to be running the UCP CA containers")
		}
		containerIDs = caContainers
		if err := ec.StopContainers(containerIDs); err != nil {
			return 1, fmt.Errorf("Failed to stop containers (%s)", err)
		}

		// TODO DRY THIS OUT
		// Move the old certs aside for posterity
		when := time.Now().Format(time.RFC3339)
		for _, filename := range []string{config.CertFilename, config.KeyFilename} {
			for _, mount := range []string{config.SwarmRootCAVolumeMount, config.OrcaRootCAVolumeMount} {
				old := filepath.Join(mount, filename)
				backup := filepath.Join(mount, fmt.Sprintf("%s_replaced_%s", filename, when))
				log.Debugf("Archiving %s -> %s", old, backup)
				if err := os.Rename(old, backup); err != nil {
					// TODO - consider refining this to handle specific errors differently
					log.Warn("Failed to rename %s -> %s : %s", old, backup, err)
				}
			}
		}
		// TODO END DRY BLOCK

		if err := certs.InitCA(config.ControllerSwarmCACN, config.SwarmRootCAVolumeMount); err != nil {
			log.Error("Failed to initialize cluster CA")
			ec.RestartContainers(containerIDs) // Ignoring failures, we're alreaddy busted...
			return 1, err
		}
		if err := certs.InitCA(config.ControllerOrcaCACN, config.OrcaRootCAVolumeMount); err != nil {
			log.Error("Failed to initialize UCP client CA")
			ec.RestartContainers(containerIDs) // Ignoring failures, we're alreaddy busted...
			return 1, err
		}

		// Make sure to restart before spitting out the final message
		err = ec.RestartContainers(containerIDs)
		if err != nil {
			log.Error(err)
			// Don't bail out, since there might be a failure mode due to the new certs, so let the user know what to do next
		}

		// Display the final message based on type of cluster
		if haDetected {
			log.Warn("The local Root CAs have been regenerated, and will not yet be trusted by the rest of this cluster!")
			log.Warn("To complete the process, proceed with the following steps:")
			log.Warn("Step 2) You should now run the 'backup --root-ca-only' command on THIS node")
			log.Warn("Step 3) Then run 'restore --root-ca-only' with that backup on ALL THE OTHER CONTROLLERs")
			log.Warn("Step 4) Then re-run 'regen-certs' without the '--root-ca-only' flag on ALL CONTROLLERs")
			log.Warn("Step 5) Then restart the docker daemons on each controller, one at a time.")
			log.Warn("Step 6) Then you can run 'join --fresh-install' and restart the daemons on all the other non-controller nodes in the cluster")

		} else {
			log.Warn("The local Root CAs have been regenerated!")
			log.Warn("To complete the process, proceed with the following steps:")
			log.Warn("Step 2) You should now run the 'regen-certs' command without the '--root-ca-only' flag on THIS controller")
			log.Warn("Step 3) Then restart the docker daemon on this node (this may take many minutes to recover)")
			log.Warn("Step 4) Then you can run 'join --fresh-install' and restart the daemons on all the other non-controller nodes in the cluster")
		}

	} else {
		if err := ec.GatherHostnames(interactive); err != nil {
			return 1, err
		}

		// TODO - Detect if this is an HA cluster, if so try to catch pebkacs
		//          - Check all CAs, make sure they're consistent, if not, tell user to fix that first
		//          - other checks?

		sameRoot, err := isSameRoot()
		if err != nil {
			return 1, err
		}

		// TODO	- Detect if the local certs align with the existing root
		//	    - Check all CAs and verify they're either placeholders or valid roots
		//	    - If they don't align, warn the user that there will be disruption

		// TODO - We could implement this to support signing by a remote CA, but that
		//        adds complexity.  For now, we'll assume a local CA to simplify the algorithm
		//        Note: Once we use swarm v2 remote CA this will need to be adjusted

		log.Info("Temporarily stopping local controller services")
		if err := ec.StopContainers(containerIDs); err != nil {
			return 1, fmt.Errorf("Failed to stop containers (%s)", err)
		}

		filesToClean := []string{config.CertFilename, config.KeyFilename, config.CAFilename}
		mountsToClean := []string{
			config.SwarmNodeCertVolumeMount,
			config.SwarmKvCertVolumeMount,
			config.SwarmControllerCertVolumeMount,
			config.AuthStoreCertsVolumeMount,
			config.AuthAPICertsVolumeMount,
			config.AuthWorkerCertsVolumeMount,
		}
		if !c.Bool("external-server-cert") {
			mountsToClean = append(mountsToClean, config.OrcaServerCertVolumeMount)
		}
		// TODO DRY THIS OUT
		// Move the old certs aside for posterity
		when := time.Now().Format(time.RFC3339)
		for _, filename := range filesToClean {
			for _, mount := range mountsToClean {
				old := filepath.Join(mount, filename)
				backup := filepath.Join(mount, fmt.Sprintf("%s_replaced_%s", filename, when))
				// TODO check for existence
				log.Debugf("Archiving %s -> %s", old, backup)
				if err := os.Rename(old, backup); err != nil {
					// TODO - consider refining this to handle specific errors differently
					log.Warnf("Failed to rename %s -> %s : %s", old, backup, err)
				}
			}
		}
		// TODO END DRY BLOCK

		// Get local node info which contains its swarm node ID.
		info, err := ec.GetClient().Info(context.Background())
		if err != nil {
			log.Error("unable to get local node info")
			return 1, err
		}

		log.Info("Regenerating local controller certs")
		// TODO - DRY this out with install:setupCerts
		for _, certConfig := range []struct {
			ou    string
			mount string
			uid   int
			gid   int
		}{
			// Local swarm manager.
			{"swarm", config.SwarmNodeCertVolumeMount, 65534, 65534},
			// KV store.
			{"kv", config.SwarmKvCertVolumeMount, 65534, 65534},
			// For UCP controller to talk to CA.
			{"ucp", config.SwarmControllerCertVolumeMount, 65534, 65534},
			// Auth datastore.
			{"auth-store", config.AuthStoreCertsVolumeMount, 65534, 65534},
			// Auth API server.
			{"auth-api", config.AuthAPICertsVolumeMount, 65534, 65534},
			// Auth worker server.
			{"auth-worker", config.AuthWorkerCertsVolumeMount, 65534, 65534},
		} {
			if err := certs.InitLocalNode(
				config.SwarmRootCAVolumeMount,
				certConfig.mount,
				info.Swarm.NodeID,
				certConfig.ou,
				config.OrcaLocalName,
				config.OrcaHostnames,
				certConfig.uid,
				certConfig.gid,
			); err != nil {
				log.Errorf("Failed to setup %s cert", certConfig.ou)
				// Partial failures probably result in a broken system, but what else can we do at this point?
				ec.RestartContainers(containerIDs) // Ignoring failures, we're alreaddy busted...
				return 1, err
			}
		}
		// Do not regenerate a server cert for the controller if the user has
		// brought their own.
		if !c.Bool("external-server-cert") {
			if err := certs.InitLocalNode(config.OrcaRootCAVolumeMount, config.OrcaServerCertVolumeMount, "ucp", "", config.OrcaLocalName, config.OrcaHostnames, 65534, 65534); err != nil {
				// Partial failures probably result in a broken system, but what else can we do at this point?
				ec.RestartContainers(containerIDs) // Ignoring failures, we're alreaddy busted...
				return 1, err
			}
		} else {
			log.Debug("Appending user provided chain of trust with root CA certs")
			caPath := filepath.Join(config.OrcaServerCertVolumeMount, "ca.pem")
			caData, err := ioutil.ReadFile(caPath)
			if err != nil {
				return 1, fmt.Errorf("Failed to process provided CA: %s", err)
			}
			clusterRootCAPath := filepath.Join(config.SwarmRootCAVolumeMount, "cert.pem")
			clusterRootCA, err := ioutil.ReadFile(clusterRootCAPath)
			if err != nil {
				return 1, fmt.Errorf("Failed to load root ca: %s", err)
			}
			clientRootCAPath := filepath.Join(config.OrcaRootCAVolumeMount, "cert.pem")
			clientRootCA, err := ioutil.ReadFile(clientRootCAPath)
			if err != nil {
				return 1, fmt.Errorf("Failed to load root ca: %s", err)
			}
			caChain := utils.JoinCerts(string(clusterRootCA), string(clientRootCA), string(caData))
			err = ioutil.WriteFile(filepath.Join(config.SwarmControllerCertVolumeMount, "ca.pem"), []byte(caChain), 0655)
			if err != nil {
				return 1, fmt.Errorf("Failed to write CA: %s", err)
			}
			// Note: we do *NOT* glue the chain of trust onto the server cert.pem, as it was signed externally
		}
		// TODO - end DRY block

		log.Info("Restarting local controller services")
		// Make sure to restart before spitting out the final message
		err = ec.RestartContainers(containerIDs)
		if err != nil {
			log.Error(err)
			// Don't bail out, since there might be a failure mode due to the new certs, so let the user know what to do next
		}

		// Display the final message based on type of cluster
		if haDetected {
			if sameRoot {
				log.Warn("The local UCP controller certs have been regenerated based on a pre-existing Root CAs!")
				log.Warn("To complete the process, proceed with the following steps:")
				log.Warn("Step 2) Continue running 'regen-certs' without the '--root-ca-only' flag on any remaining controllers")
				log.Warn("Step 3) Then restart the docker daemons on each controller, one at a time (this may take many minutes to recover)")
				log.Warn("Step 4) Then you may run 'join --fresh-install' and restart the daemons on other non-controller nodes in the cluster as needed")
			} else {
				log.Warn("The local UCP controller certs have been regenerated based on NEW Root CAs!")
				log.Warn("To complete the process, proceed with the following steps:")
				log.Warn("Step 4) Continue running 'regen-certs' without the '--root-ca-only' flag on any remaining controllers")
				log.Warn("Step 5) Then restart the docker daemons on each controller, one at a time (this may take many minutes to recover)")
				log.Warn("Step 6) Then run 'join --fresh-install' and restart the daemons on all the other non-controller nodes in the cluster")
			}

		} else {
			if sameRoot {
				log.Warn("The local UCP controller certs have been regenerated based on a pre-existing Root CAs!")
				log.Warn("To complete the process, proceed with the following steps:")
				log.Warn("Step 2) Restart the docker daemon on this node (this may take many minutes to recover)")
				log.Warn("Step 3) Then you may run 'join --fresh-install' and restart the daemons on other non-controller nodes in the cluster as needed")
			} else {
				log.Warn("The local UCP controller certs have been regenerated based on new Root CAs!")
				log.Warn("To complete the process, proceed with the following steps:")
				log.Warn("Step 3) Restart the docker daemon on this node (this may take many minutes to recover)")
				log.Warn("Step 4) Then run 'join --fresh-install' and restart the daemons on all the other non-controller nodes in the cluster")
			}
		}
	}

	if interactive {
		fmt.Printf("(Press enter to acknowledge)")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	} else {
		time.Sleep(4 * time.Second)
	}

	return 0, nil
}

// Run the recertify flow
func Run(c *cli.Context) {
	if code, err := recertify(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
