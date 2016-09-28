// Package restore will take a tar file on stdin and replace the local nodes state
package restore

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	version "github.com/hashicorp/go-version"

	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/utils"
	"github.com/docker/orca/bootstrap/utils/restarter"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

func restore(c *cli.Context) (int, error) {
	// Make sure all logging goes to stderr so stdout can be our tar file
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	if ec.IsTty() {
		return 1, fmt.Errorf("Restoring will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	}

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	if !config.InPhase2 {
		nodeType := ec.DetectNodeType()
		if nodeType != client.Controller {
			return 1, fmt.Errorf("This node is not a controller.  In the future if you encounter problems on this node, re-run the 'join' command to reconnect to the cluster.")
		}

		// Check if Stdin is attached to the bootstrapper container
		useStdin, err := isStdinAttached(ec)
		if err != nil {
			return 1, fmt.Errorf("error determining whether stdin is attached: %s", err)
		}
		if !useStdin {
			return 1, fmt.Errorf("Stdin is not attached to the \"ucp\" container. Please run this container with the \"-i\" docker flag.")
		}

		// Get the UCP Instance ID of the running instance
		ids := client.GetInstanceIDs(containers)
		if len(ids) == 0 {
			return 1, fmt.Errorf("No installed UCP instances detected on this engine")
		} else if len(ids) > 1 {
			log.Warnf("Multiple UCP instances detected: %v", ids)
		}
		id := ids[0]

		if c.Bool("interactive") {
			if c.Bool("root-ca-only") {
				log.Infof("We're about to copy Root CA material to the local controller for UCP ID: %s", id)
			} else {
				log.Infof("We're about to restore the state of the UCP cluster with ID: %s", id)
				log.Infof("All other controllers of your cluster need to be stopped before proceeding, or this operation might corrupt your cluster state.")
			}
			fmt.Fprintf(os.Stderr, "Do you want proceed with the restore operation? (y/n): ")

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
				return 1, fmt.Errorf("Aborting restore by user request")
			}
		} else {
			flagID := c.String("id")
			if flagID == "" {
				if c.Bool("root-ca-only") {
					log.Warnf("This restore operation will discard the Root CA material of this node and replace it with those found in the backup file")
				} else {
					log.Warnf("This restore operation will discard the state of the cluster and replace it with the state stored in the backup file.")
				}
				log.Infof("We detected local components of UCP instance %s", id)

				return 1, fmt.Errorf(`Re-run the command with either "--interactive" or "--id %s" to confirm you want to perform a restore operation on this UCP node.`, id)
			}

			if flagID != id {
				return 1, fmt.Errorf("The provided UCP instance ID argument %s is not equal to the ID of the local UCP instance: %s", flagID, id)
			}
		}
		config.OrcaInstanceID = id

		// Check whether the version of the local UCP installation is less than the minimum supported version for restore
		// TODO: change this to check for version equality or support cross-version restores
		localVersion, err := utils.LocalUCPVersion(containers)
		if err != nil {
			return 1, fmt.Errorf("Failed to determine version of the local UCP installation: %s", err)
		}
		targetVersion, err := version.NewVersion("1.1.0")
		if err != nil {
			return 1, fmt.Errorf("Internal error: %s", err)
		}

		if localVersion.LessThan(targetVersion) {
			return 1, fmt.Errorf("The version of the local UCP installation (%s) is less than the minimum version required for the restore operation (%s)",
				localVersion.String(), targetVersion.String())
		}

		// Figure out all the local volumes that we need to mount
		volumes, err := ec.ListExistingOrcaVolumes()
		if err != nil {
			return 1, fmt.Errorf("Failed to lookup volumes: %s", err)
		}
		if c.Bool("root-ca-only") {
			// Verify this node actually has both CA volumes
			found := 0
			for _, volume := range volumes {
				if volume == config.OrcaRootCAVolumeName || volume == config.SwarmRootCAVolumeName {
					found++
				}
			}
			if found != 2 {
				return 1, fmt.Errorf("This system does not appear to be running the UCP CA containers")
			}
			volumes = []string{config.OrcaRootCAVolumeName, config.SwarmRootCAVolumeName}
			// Reduce the set of containers to just the CA containers
			caContainers := []string{}
			for _, info := range containers {
				if strings.Contains(info.Name, orcaconfig.OrcaCAContainerName) || strings.Contains(info.Name, orcaconfig.OrcaSwarmCAContainerName) {
					caContainers = append(caContainers, id)
				}
			}
			if len(caContainers) == 0 {
				return 1, fmt.Errorf("This system does not appear to be running the UCP CA containers")
			}
		}

		for _, volume := range volumes {
			config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
				fmt.Sprintf("%s:%s", volume,
					path.Join(config.Phase2VolMountDir, volume)))
		}

		// Containers are stopped in phase 2
		return ec.StartPhase2(os.Args[1:], false)

	} else { // This is phase 2
		log.Debug("Entered Phase 2")
		log.Debug("Rediscovering containers")
		containerIDs, err := ec.FindContainerIDsByOrcaInstanceID(config.OrcaInstanceID)
		if err != nil {
			log.Debug("Failed to find specified UCP instances")
			return 1, err
		}

		var reader io.Reader
		if c.Bool("interactive") {
			reader, err = os.Open(config.BackupFile)
			if err != nil {
				return 1, fmt.Errorf("Unable to open backup file from container bind: %s", err)
			}
		} else {
			reader = os.Stdin
		}

		log.Info("Checking whether the backup file is readable")

		// Create a tar reader from the selected io.Reader
		tr, err := getTarReader(reader, c.String("passphrase"))
		if err != nil {
			return 1, err
		}

		// Try to parse the .json files at the top of the tar asynchronously
		ch := make(chan string)
		go func(ch chan string) {
			ip, err := parseInspectDumps(ec, tr, c.Bool("root-ca-only"))
			if err != nil {
				close(ch)
				log.Errorf("Unable to read from the backup file: %s", err)
				return
			}
			ch <- ip
		}(ch)

		// The master goroutine waits for completion of the inspect dump processing or a timeout
		var ip string
		select {
		case readIP, ok := <-ch:
			if !ok {
				return 1, fmt.Errorf("Backup file was not readable")
			}
			ip = readIP
		// Arbitrary timeout of 3 seconds
		case <-time.After(3 * time.Second):
			return 1, fmt.Errorf("Timed out while attempting to read metadata from the backup file. Is your backup file corrupt?")
		}

		log.Info("Backup file validated, proceeding with restore operation")

		// The restarter object guarantees that any stopped containers will be restarted once
		restarter := restarter.NewRestarter(ec)
		defer restarter.RestartAll()
		if !c.Bool("root-ca-only") {
			log.Debug("Obtaining KV URL")
			etcdURL, err := ec.FindKv()
			if err != nil {
				return 1, err
			}
			log.Debugf("KV found at %s", etcdURL.Host)
			log.Debug("Locating all UCP Controllers on the cluster")
			controllers, controllerErr := config.GetControllers(etcdURL)
			if controllerErr != nil {
				log.Errorf("Unable to connect to the KV store on this host. Please ensure that all other UCP controllers of the cluster are stopped")
				log.Warnf("The restore operation will continue automatically after 20 seconds. You may cancel this operation now with no risk of corruption")
				time.Sleep(20 * time.Second)
			}

			// Trigger the restarter to attempt to restart all UCP runtime containers
			// Only a subset of these containers will be actually present and stopped
			restarter.SetContainerTargets(orcaconfig.RuntimeContainerNames)
			log.Infof("Stopping local containers before restoring")
			if err := ec.StopContainers(containerIDs); err != nil {
				return 1, fmt.Errorf("Failed to stop containers (%s)", err)
			}

			if controllerErr == nil {
				// No controllers should respond at this point
				log.Infof("Ensuring no controllers in the cluster are still running")
				tlsConfig := &tls.Config{
					InsecureSkipVerify: true,
				}
				httpClient := &http.Client{
					Timeout: time.Second,
					Transport: &http.Transport{
						TLSClientConfig: tlsConfig,
					},
				}

				for _, controller := range controllers {
					_, err := httpClient.Get("https://" + controller.Controller)
					if err != nil {
						// A connection refused error is expected
						if strings.Contains(err.Error(), "connection refused") {
							continue
						}
						return 1, err
					} else {
						return 1, fmt.Errorf("The UCP controller at %s is still responding. Please ensure all containers are stopped on other UCP controller nodes", controller.Controller)
					}
				}
			}
		} else {
			caContainers := []string{}
			for _, info := range containers {
				if strings.Contains(info.Name, orcaconfig.OrcaCAContainerName) || strings.Contains(info.Name, orcaconfig.OrcaSwarmCAContainerName) {
					caContainers = append(caContainers, info.ID)
				}
			}
			// Trigger the restarter to restart only the CA containers in Root CA Only mode
			restarter.SetContainerTargets([]string{
				orcaconfig.OrcaSwarmCAContainerName,
				orcaconfig.OrcaCAContainerName,
			})

			if err := ec.StopContainers(caContainers); err != nil {
				return 1, fmt.Errorf("Failed to stop containers: %s", err)
			}
		}

		// We are ready to restore, start removing old state.
		// Once we start removing old state, lets proceed in the face of partial failures
		// and try to recover as much as possible with the hopes that we might get the
		// system closer to a working state
		matches, err := filepath.Glob(filepath.Join(config.Phase2VolMountDir, "*", "*"))
		if err != nil {
			// Should not happen, error out here if it does
			return 1, err
		}

		log.Info("Cleaning up old state to prepare for the restore")
		log.Debugf("Purging old state from UCP volumes mounted in %s", config.Phase2VolMountDir)

		purgeErrors := purgeVolumes(matches)

		log.Info("Beginning restore operation")
		log.Debugf("Beginning tar restore of UCP volumes mounted in %s", config.Phase2VolMountDir)
		restoreErrors := restoreVolumes(tr)

		errors := append(purgeErrors, restoreErrors...)

		if !c.Bool("root-ca-only") {
			// Start the etcd fixup container
			defer ec.RemoveContainer(orcaconfig.OrcaKvRestoreContainerName, true, false)
			cmd := []string{
				"--data-dir",
				"/data",
				"--name",
				"orca-kv-" + ip,
				"--force-new-cluster",
				"--listen-client-urls",
				"http://0.0.0.0:2379",
				"--advertise-client-urls",
				"http://127.0.0.1:2379",
				"--listen-peer-urls",
				"http://0.0.0.0:2380",
			}
			id, err := ec.StartEtcdFixupContainer(cmd)
			if err != nil {
				errors = append(errors, err)
			}

			log.Debug("KV restore container started")

			// Create the http client that we'll use to curl the recovery container

			etcdClient := &http.Client{
				Timeout: 10 * time.Second,
			}

			timeout := 5 * time.Minute
			startTime := time.Now()

			// Wait until the kv restore container is reachable
			log.Debug("Waiting for the Restore KV to become healthy")
			for startTime.Add(timeout).After(time.Now()) {
				time.Sleep(3 * time.Second)

				info, err := ec.InspectContainer(id)
				if err != nil {
					errors = append(errors, fmt.Errorf("unable to inspect the %s container: %s",
						orcaconfig.OrcaKvRestoreContainerName, err))
				}

				// The recovery container exited,
				if info.State.Status == "exited" {
					log.Info(ec.GetContainerLogs(id))
					errors = append(errors, fmt.Errorf("the recovery etcd container exited with an error"))
					break
				}
				etcdHost := fmt.Sprintf("http://%s:2379", info.NetworkSettings.IPAddress)
				healthy, err := client.IsEtcdHealthy(etcdClient, etcdHost)
				if err != nil {
					log.Debug(err.Error())
					if strings.Contains(err.Error(), "connection refused") {
						continue
					}
					errors = append(errors, err)
				}
				if healthy {
					log.Debug("KV Restore seems healthy")
					break
				}
			}
			ec.RemoveContainer(orcaconfig.OrcaKvRestoreContainerName, true, false)

			// Fix up kv file permissions
			err = os.Chown(filepath.Join(config.Phase2VolMountDir, config.OrcaKVVolumeName), 65534, 65534)
			if err != nil {
				return 1, err
			}
			err = filepath.Walk(filepath.Join(config.Phase2VolMountDir, config.OrcaKVVolumeName), func(path string, info os.FileInfo, err error) error {
				err = os.Chown(path, 65534, 65534)
				if err != nil {
					log.Warn("Failed to update permissions: %s %s", path, err)
					return err
				}
				return nil
			})
			if err != nil {
				return 1, err
			}

			// Restart the ucp-kv container
			err = restarter.RestartContainer(orcaconfig.OrcaKvContainerName)
			if err != nil {
				errors = append(errors, err)
			}

			// Make an http.Client for the original ucp-kv container
			tlsConfig, err := client.LoadCerts(config.SwarmKvCertVolumeMount)
			originalKVClient := &http.Client{
				Timeout: time.Second,
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
			}

			etcdHost := "https://" + ip + ":12379"
			startTime = time.Now()
			err = client.WaitForHealthyKv(originalKVClient, etcdHost, timeout)
			if err != nil {
				errors = append(errors, err)
			}

			err = EtcdRestoreMembers(originalKVClient, ip)
			if err != nil {
				errors = append(errors, err)
			}

			kvStoreURL, err := url.Parse("etcd://" + ip + ":12379")
			if err != nil {
				errors = append(errors, fmt.Errorf("Unable to generate a URL for  the restarted KV store: %s ", err))
			}

			// Try to get the set of controllers from the KV store
			startTime = time.Now()
			var controllers []types.Controller
			for startTime.Add(timeout).After(time.Now()) {
				time.Sleep(3 * time.Second)
				controllers, err = config.GetControllers(kvStoreURL)
				if err != nil {
					log.Debug(err)
					continue
				} else {
					break
				}
			}

			if err != nil {
				errors = append(errors, fmt.Errorf("Unable to locate the previous set of controllers from the restarted KV store: %s ", err))
			}

			log.Info("Removing old controllers from the KV Store")
			for _, controller := range controllers {
				controllerIP, _, err := net.SplitHostPort(controller.Controller)
				if err != nil {
					errors = append(errors, fmt.Errorf("Unable to detect controller IP: %s", err))
					continue
				}
				if controllerIP == ip {
					continue
				}
				if err := config.RemoveSwarmManager(kvStoreURL, controllerIP); err != nil {
					log.Warningf("Failed to remove swarm manager registration from KV store: %s", err)
				}
			}
			err = restarter.RestartContainer(orcaconfig.AuthStoreContainerName)
			if err != nil {
				errors = append(errors, fmt.Errorf("Unable to restart \"%s\" container: %s",
					orcaconfig.AuthStoreContainerName, err))
			}

			authStoreAddr := net.JoinHostPort(ip,
				fmt.Sprintf("%d", orcaconfig.AuthStorePort))

			// Perform a SyncDB --emergency-repair operation, up to five times
			// TODO: add a health check for rethinkdb. Retries should be good for now
			retries := 5
			for {
				if retries <= 0 {
					break
				}
				time.Sleep(10 * time.Second)
				err = ec.RunAuthSyncDB(true, authStoreAddr)
				if err != nil {
					retries--
					log.Debugf("Failed to restore auth store, retrying : %s", err)
				} else {
					break
				}
				ec.RemoveContainer(orcaconfig.AuthSyncDBContainerName, true, false)
			}
			if err != nil {
				errors = append(errors, err)
			}

			// Restart the remaining containers
			log.Info("Restarting UCP containers")
			restarter.RestartAll()

			// Wait for the UCP controller to come alive
			info, err := ec.InspectContainer(orcaconfig.OrcaControllerContainerName)
			if err != nil {
				errors = append(errors, fmt.Errorf("unable to inspect the %s container: %s",
					orcaconfig.OrcaControllerContainerName, err))
			}
			ctrlPortMap, ok := info.NetworkSettings.Ports["8080/tcp"]
			if !ok {
				errors = append(errors, fmt.Errorf("unable to extract port mapping information for the ucp controller"))
			}
			ucpURL := "https://" + net.JoinHostPort(ip, ctrlPortMap[0].HostPort)

			err = client.WaitForOrca(ucpURL, time.Minute)

			if err != nil {
				errors = append(errors, err)
			}

			log.Warnf("Cluster membership was restored. Please ensure that this controller is operational by accessing the User Interface")
			log.Warnf("IMPORTANT: The restore operation has reset your cluster to have a single non-HA controller.")
			log.Warnf("In order to re-join your old controllers to this cluster, first 'uninstall' those controllers, then perform a 'join'. Failing to do so might leave your cluster in a corrupted state")
		}

		if len(errors) > 0 {
			for _, err = range errors {
				log.Warnf("%s", err)
			}
			return 1, fmt.Errorf("errors were encountered, the restored instance may be corrupted")
		}
		log.Info("Restore completed successfully")
	}
	return 0, nil
}

// Run the restore flow
func Run(c *cli.Context) {
	if code, err := restore(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
