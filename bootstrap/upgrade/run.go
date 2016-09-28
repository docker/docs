// Package upgrade implements the high-level upgrade flow for orca
package upgrade

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/discovery"
	"github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/config/discover"
	"golang.org/x/net/context"
)

type sortedConfigKeys struct {
	Keys []string
	Cfg  map[string]*types.ContainerJSON
}

func (s sortedConfigKeys) Len() int {
	return len(s.Keys)
}

func (s sortedConfigKeys) Swap(i, j int) {
	s.Keys[i], s.Keys[j] = s.Keys[j], s.Keys[i]
}

func (s sortedConfigKeys) Less(i, j int) bool {
	// orca-controller is always last
	iInfo := s.Cfg[s.Keys[i]]
	jInfo := s.Cfg[s.Keys[j]]
	if iInfo.Name == "orca-controller" {
		return false
	} else if jInfo.Name == "orca-controller" {
		return true
	}

	// WARNING - this algorithm doesn't walk back links so it only works
	//           for single level links (all we use in UCP today)
	//           If we start to do multi-level linking then this will
	//           need adjusting to ensure a proper sort order

	// If these two depend on eachother, simple comparison
	if iInfo.HostConfig != nil {
		for _, link := range iInfo.HostConfig.Links {
			linkName := path.Base(strings.Split(link, ":")[0])
			if linkName == jInfo.Name {
				return false
			}
		}
	}
	return true
}

func verifyCAMaterialValid() error {
	mounts := [2]string{config.SwarmRootCAVolumeMount, config.OrcaRootCAVolumeMount}
	for _, mount := range mounts {
		filename := filepath.Join(mount, config.CertFilename)
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Failed to load %s: %s", filename, err)
		}
		for len(data) > 0 {
			der, rest := pem.Decode(data)
			if der == nil {
				return fmt.Errorf("Failed to decode pem %s", filename)
			}
			c, err := x509.ParseCertificate(der.Bytes)
			if err != nil {
				return fmt.Errorf("Failed to parse cert %s: %s", filename, err)
			}
			if strings.Contains(c.Subject.CommonName, "Placeholder invalid") {
				return fmt.Errorf("This node contains a placeholder root CA, not a valid root CA.  Please re-run the command on the controller which contains the valid root, or backup and restore the root CA material onto this node.")
			}
			data = rest
		}
	}
	return nil
}

// Run the join flow
func upgrade(c *cli.Context) (int, error) {
	log.Debug("Starting upgrade")

	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	dclient := ec.GetClient()

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	if !config.InPhase2 {
		// Since upgrade can be potentially disruptive, get confirmation before proceeding
		ids := client.GetInstanceIDs(containers)
		if len(ids) == 0 {
			return 1, fmt.Errorf("No running UCP instances detected on this engine")
		} else if len(ids) > 1 {
			log.Warnf("Multiple UCP instances detected: %v", ids)
		}

		id := c.String("id")
		if id != "" {
			config.OrcaInstanceID = id
		} else {
			// Normally this will be a single item, and we're not optimizing for multiple orca's
			id = ids[0]
			if c.Bool("interactive") {
				log.Infof("Upgrade UCP %s containers on this engine to %s for UCP ID: %s", utils.GetUCPVersionString(containers), orcaconfig.ImageVersion, id)
				fmt.Printf("Do you want proceed with the upgrade? (y/n): ")

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
					return 1, fmt.Errorf("Not upgrading per user request")
				}

			} else {
				log.Infof("We detected local components of UCP instance %s", id)
				return 1, fmt.Errorf(`Re-run the command with "--id %s" or --interactive to confirm you want to upgrade this UCP instance.`, id)
			}
		}

		// Start by pulling all the updated images
		if err := ec.VerifyOrPullImages(c.Bool("interactive")); err != nil {
			log.Error("We were unable to pull one or more required images.  Please set REGISTRY_USERNAME, REGISTRY_PASSWORD, and REGISTRY_EMAIL environment variables for your Docker Hub account on this container with -e flags to run.")
			return 1, err
		}
		config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
			fmt.Sprintf("%s:%s", config.EngineConfigDir, config.EngineConfigDir),
			fmt.Sprintf("%s:%s", config.SwarmKvCertVolumeName, config.SwarmKvCertVolumeMount),
			fmt.Sprintf("%s:%s", config.OrcaKVVolumeName, config.OrcaKVVolumeMount),
			fmt.Sprintf("%s:%s", config.SwarmNodeCertVolumeName, config.SwarmNodeCertVolumeMount),
			fmt.Sprintf("%s:%s", config.SwarmControllerCertVolumeName, config.SwarmControllerCertVolumeMount),
			fmt.Sprintf("%s:%s", config.OrcaRootCAVolumeName, config.OrcaRootCAVolumeMount),
			fmt.Sprintf("%s:%s", config.SwarmRootCAVolumeName, config.SwarmRootCAVolumeMount),
			fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, config.OrcaServerCertVolumeMount),
			// Auth Service Certificate Volume Mounts.
			fmt.Sprintf("%s:%s", config.AuthStoreCertsVolumeName, config.AuthStoreCertsVolumeMount),
			fmt.Sprintf("%s:%s", config.AuthStoreDataVolumeName, config.AuthStoreDataVolumeMount),
			fmt.Sprintf("%s:%s", config.AuthAPICertsVolumeName, config.AuthAPICertsVolumeMount),
			fmt.Sprintf("%s:%s", config.AuthWorkerCertsVolumeName, config.AuthWorkerCertsVolumeMount),
			fmt.Sprintf("%s:%s", config.AuthWorkerDataVolumeName, config.AuthWorkerDataVolumeMount),
		)
		return ec.StartPhase2(os.Args[1:], false)
	} else {
		containers, err := dclient.ContainerList(context.TODO(), types.ContainerListOptions{All: true})
		if err != nil {
			log.Debug("Failed to find specified UCP instances")
			return 1, err
		}
		found := true
		log.Info("Checking for version compatibility")
		for _, container := range containers {
			name := path.Base(container.Names[0])
			// Filter out so we're only looking at relevant containers
			check := false
			for _, n := range orcaconfig.RuntimeContainerNames {
				if name == n {
					check = true
				}
			}
			if !check {
				continue
			}
			log.Debugf("Checking container %s for upgrade compatibility", name)
			if id, found := container.Labels[fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)]; found {
				if id != config.OrcaInstanceID {
					continue
				}
			} else {
				continue
			}
			// Assume first name is what we want
			// Make sure the upgrade path is compatible
			newImage, err := orcaconfig.GetContainerImage(name)
			if err != nil {
				log.Debugf("Failed to lookup %s: %s", name, err)
				return 1, fmt.Errorf("Unable to perform the requested upgrade.  Please refer to UCP documentation for a supported upgrade path")
			}
			if !ec.CheckUpgradeCompatible(container.Image, newImage) {
				return 1, fmt.Errorf("Unable to perform the requested upgrade.  Please refer to UCP documentation for a supported upgrade path")
			}
			// Upgrade possible for this container
		}
		if !found {
			return 1, fmt.Errorf("No matching UCP containers detected for ID: %s", config.OrcaInstanceID)
		}

		// Upgrade is safe to proceed based on version compatibility check

		// Extract the current cluster/node configuration
		nodeConfig, err := discover.DiscoverNodeConfig(ec.GetClient())
		if err != nil {
			return 1, err
		}

		// XXX This shouldn't be necessary once we get everything wired up properly
		config.OrcaHostAddress = nodeConfig.HostAddress
		orcaconfig.SwarmPort, err = strconv.Atoi(nodeConfig.SwarmPort)
		if err != nil {
			return 1, err
		}
		orcaconfig.OrcaPort, err = strconv.Atoi(nodeConfig.ControllerPort)
		if err != nil {
			return 1, err
		}

		// Determine if we're already in swarm mode, or need to possibly unwind engine-discovery
		swarmV2, err := ec.FindSwarmV2Cluster()
		if err != nil {
			return 1, fmt.Errorf("Error when checking for Swarm cluster: %s", err)
		}
		if swarmV2 != "manager" && swarmV2 != "" {
			// Either worker or unknown role
			return 1, fmt.Errorf("This node is a %s in an existing swarm. Either promote it to a manager or leave the swarm before installing UCP.", swarmV2)
		}
		// Create a new swarm if there isn't one.
		initNewSwarm := swarmV2 == ""

		// Verify that the root CA certs on this controller node are good
		if err := verifyCAMaterialValid(); err != nil {
			return 1, err
		}

		if initNewSwarm {
			// Load current discovery configuration (if present)
			engineCfg, err := discovery.LoadCurrentConfiguration()
			if err != nil {
				log.Info("Failed to lookup existing engine discovery configuration")
				return 1, err
			}
			needsRestart := false
			if _, found := engineCfg[discovery.ClusterAdvertise]; found {
				needsRestart = true
				delete(engineCfg, discovery.ClusterAdvertise)
			}
			if _, found := engineCfg[discovery.ClusterStore]; found {
				needsRestart = true
				delete(engineCfg, discovery.ClusterStore)
			}
			if _, found := engineCfg[discovery.ClusterStoreOpts]; found {
				needsRestart = true
				delete(engineCfg, discovery.ClusterStoreOpts)
			}
			if needsRestart {
				if err := discovery.ApplyConfigChanges(ec, engineCfg, 1, false); err != nil {
					log.Info("Failed to update engine configuration")
					return 1, err
				}
				return 2, fmt.Errorf("Engine discovery mode has been unconfigured.  To proceed with the upgrade, you must restart your daemon, then re-run this upgrade tool")
			}
			// Note: there doesn't appear to be a good way to detect if the user forgot to bounce, so we may still fail the swarm init

			// TODO - should we somehow extract the old overlay network definitions? How do we get that upgraded?

			swarmListenAddr := config.OrcaHostAddress

			listenAddr := fmt.Sprintf("%s:%d", swarmListenAddr, orcaconfig.SwarmGRPCPort)
			nodeName, err := ec.CreateSwarmV2Swarm(listenAddr)
			if err != nil {
				return 1, err
			}
			log.Debugf("Created new swarm-mode cluster, this node is %s", nodeName)
		}

		orcaEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.OrcaPort)

		// Stop the controller so we make sure when a ping comes back it's post upgrade
		// TODO - what if it's already stopped?  Will this error out?
		if err := ec.StopContainer(orcaconfig.OrcaControllerContainerName); err != nil {
			return 1, err
		}

		log.Infof("Deploying UCP service")
		// TODO - theoretically this should take a cluster config, and we can just pass what we found above...
		if err := ec.StartUCPAgentService(); err != nil {
			return 1, err
		}
		// Attempt to gather output from the beachhead so we can see if something goes wrong
		// TODO - if we can wire up debug vs. non-debug to the service, then we can do this unconditionally
		if c.Bool("debug") {
			go func() {
				client := ec.GetClient()
				done := false
				for !done {
					time.Sleep(1 * time.Second)
					services, err := client.ServiceList(context.TODO(), types.ServiceListOptions{})
					if err != nil {
						log.Debugf("Failed to lookup services: %s", err)
						continue
					}
					// Try to find the beachhead
					serviceID := ""
					for _, service := range services {
						if service.Spec.Annotations.Name == "ucp-beachhead" {
							serviceID = service.ID
						}

					}
					if serviceID == "" {
						continue
					}
					filter := filters.NewArgs()
					filter.Add("service", serviceID)
					// TODO - consider adding node filter
					tasks, err := client.TaskList(context.TODO(), types.TaskListOptions{Filter: filter})
					if err != nil {
						log.Debugf("Failed to lookup tasks: %s", err)
						continue
					}
					for _, task := range tasks {
						id := task.Status.ContainerStatus.ContainerID
						if id == "" {
							// The container probably hasn't started yet...
							continue
						}
						rd, err := client.ContainerLogs(context.TODO(), id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
						if err != nil {

							// TODO - consider trying to get the log from the proxy so we can pick up remote ones
							log.Debugf("Failed to get logs for %s: %s", id, err)
							continue
						}
						done = true
						go func() {
							oldLevel := log.GetLevel()
							log.SetLevel(log.InfoLevel)
							defer log.SetLevel(oldLevel)
							if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, rd); err != nil {
								log.Errorf("Failed to stream logs: %s\n", err)
								return
							}
						}()
					}
				}
			}()
		}

		serviceWaitTimeout := 3 * config.AliveCheckTimeout
		if err := client.WaitForOrca(fmt.Sprintf("https://%s", orcaEndpoint), serviceWaitTimeout); err != nil {
			// XXX: should we keep this log statement?
			utils.SpewContainerLogs(ec, "ucp-proxy-phase2")

			log.Errorf(`UCP didn't come up within %v.  Run "docker logs %s" for more details`, serviceWaitTimeout, orcaconfig.OrcaProxyContainerName)
			return 1, err
		}

		log.Info("Success!  Please log in to the UCP console to verify your system before proceeding to upgrade additional nodes.")
	}
	return 0, nil
}

// Run the upgrade flow
func Run(c *cli.Context) {
	if code, err := upgrade(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
