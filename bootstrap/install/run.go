// Package install implements the high-level installation flow for orca
package install

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/howeyc/gopass"
	"golang.org/x/net/context"

	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
)

// Run the installation flow
func Run(c *cli.Context) {
	if code, err := install(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}

func install(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	if c.Bool("binpack") && c.Bool("random") {
		return 1, fmt.Errorf("You can only pick one scheduler strategy.  The default spread (if unspecified), '--binpack' or '--random'")
	}
	if c.Bool("binpack") {
		config.SwarmBinpack = true
	}
	if c.Bool("random") {
		config.SwarmRandom = true
	}
	if c.Bool("swarm-experimental") {
		config.SwarmExperimental = true
	}
	if c.IsSet("kv-timeout") {
		config.KVTimeout = c.Int("kv-timeout")
	}

	// Wire up ports
	orcaconfig.SwarmPort = c.Int("swarm-port")
	orcaconfig.OrcaPort = c.Int("controller-port")
	orcaconfig.SwarmGRPCPort = c.Int("swarm-grpc-port")

	if !config.InPhase2 {
		config.GetNewID()
		log.Info("Verifying your system is compatible with UCP")
	} else {
		log.Debugf("Beginning phase 2 install for instance %s", config.OrcaInstanceID)
	}

	ec, err := client.New()
	if err != nil {
		return 1, err
	}

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

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	if containers != nil {
		if !c.Bool("fresh-install") {
			return 1, errors.New("Existing UCP containers detected.  Either use --fresh-install or upgrade")
		}
	}

	if !config.InPhase2 {
		return runPhase1(c, ec)
	}

	// We're in phase 2
	return runPhase2(c, ec, containers, initNewSwarm)
}

func runPhase1(c *cli.Context, ec *client.EngineClient) (int, error) {
	if err := ec.SystemValidation(); err != nil {
		return 1, err
	}

	if c.Bool("interactive") && !promptForPassword() {
		return 1, fmt.Errorf("Password configuration error, giving up")
	}

	// Make sure we've got all our images lined up first, since we need some to run our early tests
	if err := ec.VerifyOrPullImages(c.Bool("interactive")); err != nil {
		log.Error("We were unable to pull one or more required images.  Please set REGISTRY_USERNAME, REGISTRY_PASSWORD, and REGISTRY_EMAIL environment variables for your Docker Hub account on this container with -e flags to run.")
		return 1, err
	}

	// Figure out addresses/hostnames next, in case they fail, so we can tell the user how to workaround it
	if err := ec.GatherHostnames(c.Bool("interactive")); err != nil {
		return 1, err
	}

	// host address must currently be an IP to advertise in an engine cluster
	if net.ParseIP(config.OrcaHostAddress) == nil {
		return 1, fmt.Errorf("Unsupported host address (%s). Please specify an IPv4 address using the '--host-address' flag, and pass the hostname using the '--san' flag.", config.OrcaHostAddress)
	}

	// Show the user what we think the host address is
	log.Infof("Installing UCP with host address %s - If this is incorrect, please specify an alternative address with the '--host-address' flag", config.OrcaHostAddress)

	// Prepare all the volumes we're going to need
	if err := ec.PrepareAllVolumes(); err != nil {
		return 1, fmt.Errorf("Failed to set up required storage volumes: %s", err)
	}

	// Wire up the extra volume mounts we want for the install flow
	config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
		fmt.Sprintf("%s:%s", config.OrcaRootCAVolumeName, config.OrcaRootCAVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmRootCAVolumeName, config.SwarmRootCAVolumeMount),
		fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, config.OrcaServerCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmNodeCertVolumeName, config.SwarmNodeCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmKvCertVolumeName, config.SwarmKvCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmControllerCertVolumeName, config.SwarmControllerCertVolumeMount),
		fmt.Sprintf("%s:%s", config.OrcaKVVolumeName, config.OrcaKVVolumeMount),
		fmt.Sprintf("%s:%s", config.EngineConfigDir, config.EngineConfigDir),
		fmt.Sprintf("%s:%s", config.EnginePidDir, config.EnginePidDir),
		// Auth Service Related Volume Mounts.
		fmt.Sprintf("%s:%s", config.AuthStoreDataVolumeName, config.AuthStoreDataVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthStoreCertsVolumeName, config.AuthStoreCertsVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthAPICertsVolumeName, config.AuthAPICertsVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthWorkerDataVolumeName, config.AuthWorkerDataVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthWorkerCertsVolumeName, config.AuthWorkerCertsVolumeMount),
	)

	return ec.StartPhase2(os.Args[1:], true)
}

func promptForPassword() bool {
	for range []int{1, 2, 3} {
		fmt.Print("Please choose your initial UCP admin password: ")
		pass1 := string(gopass.GetPasswd())
		if strings.TrimSpace(pass1) == "" {
			log.Error("Password cannot be empty, please try again")
			continue
		}

		fmt.Print("Confirm your initial password: ")
		pass2 := string(gopass.GetPasswd())
		if pass1 == pass2 {
			os.Setenv("UCP_ADMIN_PASSWORD", pass1)
			return true
		}

		log.Error("Passwords don't match, please try again")
	}

	return false
}

func runPhase2(c *cli.Context, ec *client.EngineClient, containers []*types.ContainerJSON, initNewSwarm bool) (int, error) {
	// Reload hostnames from the environment
	if err := ec.GatherHostnames(false); err != nil {
		return 1, err
	}

	if initNewSwarm {
		orcaconfig.RequiredPorts = append(orcaconfig.RequiredPorts, &orcaconfig.SwarmGRPCPort)
	}

	if containers != nil {
		log.Info("Removing old UCP containers")
		if err := ec.RemoveOrcaContainers(containers); err != nil {
			return 1, fmt.Errorf("Failed to remove an existing UCP container %s", err)
		}
	}

	// Check ports *after* we've shut down any running orca conatiners
	if err := ec.CheckPorts(orcaconfig.RequiredPorts); err != nil {
		return 1, err
	}

	// Purge volumes if necessary.
	if c.Bool("fresh-install") {
		if err := utils.CleanupVolumes(c); err != nil {
			return 1, err
		}
	}

	// Verify filesystem permissions
	if err := utils.VerifyPermissions(false); err != nil {
		return 1, err
	}

	if initNewSwarm {
		swarmListenAddr := config.OrcaHostAddress
		if c.String("host-address") == "" {
			swarmListenAddr = "0.0.0.0"
		}

		listenAddr := fmt.Sprintf("%s:%d", swarmListenAddr, orcaconfig.SwarmGRPCPort)
		nodeName, err := ec.CreateSwarmV2Swarm(listenAddr)
		if err != nil {
			return 1, err
		}
		log.Debugf("Created new swarm, this node is %s", nodeName)
	}

	// Set the hostnames as node labels since we are a swarm manager
	if err := ec.SetHostnamesAsLabels(); err != nil {
		return 1, err
	}

	nodeName, err := ec.FindSwarmV2NodeID()
	if err != nil {
		return 1, err
	}
	// Re-grab the host-address from the node APIs.
	log.Debugf("This node is known as %s", nodeName)
	host, err := ec.FindIPFromNode(nodeName)
	if err != nil {
		return 1, err
	}
	log.Debugf("Got node IP %s from swarm", host)
	config.OrcaHostAddress = host

	if err := setupCerts(c, ec); err != nil {
		return 1, err
	}

	orcaEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.OrcaPort)

	log.Infof("Deploying UCP Service")
	if err := utils.DeployUCPAgentService(ec); err != nil {
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
					if service.Spec.Annotations.Name == "ucp-agent" {
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
		utils.SpewContainerLogs(ec, "ucp-reconcile")

		log.Errorf(`UCP didn't come up within %v.  Run "docker logs %s" for more details`, serviceWaitTimeout, orcaconfig.OrcaReconcileContainerName)
		return 1, err
	}

	// Report the Orca cert fingerprint
	fingerprint, err := certs.GetFingerprint(filepath.Join(config.OrcaServerCertVolumeMount, "cert.pem"))
	if err != nil {
		log.Error("Failed to get fingerprint")
		return 1, err
	}
	log.Infof("UCP Instance ID: %s", config.OrcaInstanceID)
	log.Infof("UCP Server SSL: %s", fingerprint)

	userMsg := "\"admin\""
	if os.Getenv("UCP_ADMIN_USER") != "" {
		userMsg = fmt.Sprintf("\"%s\"", os.Getenv("UCP_ADMIN_USER"))
	}

	passMsg := "\"orca\""
	if os.Getenv("UCP_ADMIN_PASSWORD") != "" {
		passMsg = "(your admin password)"
	}

	// Report the Orca URL
	log.Infof("Login as %s/%s to UCP at https://%s", userMsg, passMsg, orcaEndpoint)

	return 0, nil
}

func setupCerts(c *cli.Context, ec *client.EngineClient) error {
	// Copy swarm's CA root key and certificate to our volume
	if err := certs.CopySwarmRootCA(ec.GetClient()); err != nil {
		log.Error("Unable to copy Swarm Root CA material")
		return err
	}

	// Get local node info which contains its swarm node ID.
	info, err := ec.GetClient().Info(context.Background())
	if err != nil {
		log.Error("unable to get local node info")
		return err
	}

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
			return err
		}
	}

	// Orca CA
	if err := certs.InitCA(config.ControllerOrcaCACN, config.OrcaRootCAVolumeMount); err != nil {
		log.Error("Failed to initialize UCP CA")
		return err
	}

	// Do not generate a server cert for the controller if the user has
	// brought their own.
	if !(c.Bool("external-ucp-ca") || c.Bool("external-server-cert")) {
		if err := certs.InitLocalNode(config.OrcaRootCAVolumeMount, config.OrcaServerCertVolumeMount, "ucp", "", config.OrcaLocalName, config.OrcaHostnames, 65534, 65534); err != nil {
			return err
		}
	}

	// Make sure orca trusts swarm too
	if err := certs.SetupOrcaTrust(); err != nil {
		return err
	}

	return nil
}
