// Package join implements the high-level host add flow for orca
package join

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/types"
	"github.com/docker/orca"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/orcaclient"
	"github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
	orcatypes "github.com/docker/orca/types"
	"golang.org/x/net/context"
)

// TODO - There's a lot of commonality with install/run.go - this should be refactored
//        into some common routines that each can leverage

func Run(c *cli.Context) {
	if code, err := join(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}

// Run the join flow
func join(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	// Wire up ports
	orcaconfig.SwarmPort = c.Int("swarm-port")
	orcaconfig.OrcaPort = c.Int("controller-port")

	log.Debug("Starting join")

	ec, err := client.New()
	if err != nil {
		return 1, err
	}

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	if containers != nil {
		if !c.Bool("fresh-install") {
			return 1, errors.New("Existing UCP containers detected.  Either use --fresh-install or upgrade")
		}
	}

	if !config.InPhase2 {
		return runPhase1(c, ec)
	}

	return runPhase2(c, ec, containers)
}

func runPhase1(c *cli.Context, ec *client.EngineClient) (int, error) {
	if err := ec.SystemValidation(); err != nil {
		return 1, err
	}
	// Figure out the connection URL before we launch phase 2, so we can do an interactive TOFU
	if c.Bool("interactive") {
		if err := config.InteractivePrompt("UCP_URL"); err != nil {
			return 1, errors.New("Failed to gather interactive input")
		}
	} else if c.String("url") != "" {
		os.Setenv("UCP_URL", c.String("url"))
	}
	orcaEnv := os.Getenv("UCP_URL")
	if orcaEnv == "" {
		return 1, errors.New("You must specify an UCP server to join.  Either use --interactive, or --url")
	}

	orcaURL, err := url.Parse(orcaEnv)
	if err != nil {
		return 1, fmt.Errorf("Malformed UCP connection URL: %s", err)
	}
	if orcaURL.Scheme != "https" {
		return 1, fmt.Errorf(`We only support https connection URLs - you specified "%s"`, orcaURL.String())
	}

	// Get optionns like fingerprint, username, password, etc (setting the
	// command line arg or env variable in the process).
	if err := gatherOpts(c, orcaURL); err != nil {
		return 1, err
	}

	// TODO - Connect to the UCP instance and check the version - reject if the X.Y don't match (let .Z drift)

	// Make sure we've got all our images lined up first, since we need some to run our early tests
	// Technically this is more than we need, but lets load them all anyway for now...
	if err := ec.VerifyOrPullImages(c.Bool("interactive")); err != nil {
		log.Error("Unable to pull image.  Please set REGISTRY_USERNAME, REGISTRY_PASSWORD, REGISTRY_EMAIL in your environment")
		return 1, err
	}

	config.OrcaInstanceID = "unknown" // This will get filled in correctly after we call the join API
	if err := ec.GatherHostnames(c.Bool("interactive")); err != nil {
		return 1, err
	}

	// Prepare the volumes we need
	if err := ec.PrepareVolume(config.SwarmNodeCertVolumeName); err != nil {
		return 1, fmt.Errorf("Failed to set up required storage volume: %s", err)
	}

	config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
		fmt.Sprintf("%s:%s", config.OrcaKVVolumeName, config.OrcaKVVolumeMount),
		fmt.Sprintf("%s:%s", config.OrcaRootCAVolumeName, config.OrcaRootCAVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmRootCAVolumeName, config.SwarmRootCAVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmNodeCertVolumeName, config.SwarmNodeCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmKvCertVolumeName, config.SwarmKvCertVolumeMount),
		fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, config.OrcaServerCertVolumeMount),
		fmt.Sprintf("%s:%s", config.SwarmControllerCertVolumeName, config.SwarmControllerCertVolumeMount),
		fmt.Sprintf("%s:%s", config.EngineConfigDir, config.EngineConfigDir),
		fmt.Sprintf("%s:%s", config.EnginePidDir, config.EnginePidDir),
		// Auth Service Related Volume Mounts.
		fmt.Sprintf("%s:%s", config.AuthStoreDataVolumeName, config.AuthStoreDataVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthStoreCertsVolumeName, config.AuthStoreCertsVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthAPICertsVolumeName, config.AuthAPICertsVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthWorkerDataVolumeName, config.AuthWorkerDataVolumeMount),
		fmt.Sprintf("%s:%s", config.AuthWorkerCertsVolumeName, config.AuthWorkerCertsVolumeMount),
	)

	// TODO XXX - only add the extra mounts if this is a replica so we don't polute a non-replica node with extra volumes

	log.Debug("Launching phase 2")
	return ec.StartPhase2(os.Args[1:], true)
}

func gatherOpts(c *cli.Context, orcaURL *url.URL) error {
	if err := getFingerprint(c, orcaURL); err != nil {
		return err
	}

	// Also prompt for the admin username and password if in interactive
	// mode, otherwise, assume they are specified as env vars.
	if !c.Bool("interactive") {
		return nil
	}

	if err := config.InteractivePrompt("UCP_ADMIN_USER"); err != nil {
		return errors.New("Failed to gather interactive input")
	}
	if err := config.InteractivePrompt("UCP_ADMIN_PASSWORD"); err != nil {
		return errors.New("Failed to gather interactive input")
	}

	return nil
}

func getFingerprint(c *cli.Context, orcaURL *url.URL) error {
	fingerprint := c.String("fingerprint")
	if fingerprint != "" {
		// Fingerprint was specified on the command line, so we can
		// pass the fingerprint command line arg through to phase 2.
		return nil
	}

	lines, err := certs.Tofu(orcaURL.Host)
	if err != nil {
		return fmt.Errorf("Failed to connect to UCP URL: %s", err)
	}

	if len(lines) == 0 {
		// The CA for the orcaURL is already trusted by the system.
		return nil
	}

	if !c.Bool("interactive") && !c.Bool("insecure-fingerprint") {
		// We must be in interactive mode to prompt for acceptance of
		// the server cert fingerprint.
		log.Infof("UCP server %s", orcaURL.String())
		for _, line := range lines {
			log.Info(line)
			if strings.Contains(line, "Fingerprint=") {
				fingerprint = line
			}
		}
		return errors.New("Repeat with --interactive and accept the fingerprint, or give the --fingerprint argument eactly matching a fingerprint above to proceed")
	}

	fmt.Printf("UCP server %s\n", orcaURL.String())
	for _, line := range lines {
		// only print in interactive mode
		if c.Bool("interactive") {
			fmt.Println(line)
		}
		if strings.Contains(line, "Fingerprint=") {
			fingerprint = line
		}
	}

	if c.Bool("insecure-fingerprint") {
		// Accept fingerprint on behalf of user here
		log.Debugf(`User used --insecure-fingerprint, automatically accepting "%s"`, fingerprint)
		os.Setenv("UCP_FINGERPRINT", fingerprint)
		return nil
	}

	fmt.Printf("Do you want to trust this server and proceed with the join? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		log.Debugf("Failed to read input: %s", err)
		return err
	}

	value = strings.TrimSpace(strings.ToLower(value))
	if !(value == "y" || value == "yes") {
		return errors.New("Repeat with --interactive and accept the fingerprint, or give the --fingerprint argument eactly matching a fingerprint above to proceed")
	}

	log.Debugf(`User accepted TLS Fingerprint "%s", proceeding`, fingerprint)
	os.Setenv("UCP_FINGERPRINT", fingerprint)

	return nil
}

func runPhase2(c *cli.Context, ec *client.EngineClient, containers []*types.ContainerJSON) (int, error) {
	// Reload hostnames from the environment
	if err := ec.GatherHostnames(false); err != nil {
		return 1, err
	}

	// host address must currently be an IP to advertise in an engine cluster
	if net.ParseIP(config.OrcaHostAddress) == nil {
		return 1, fmt.Errorf("Unsupported host address (%s). Please specify an IPv4 address using the '--host-address' flag, and pass the hostname using the '--san' flag.", config.OrcaHostAddress)
	}

	// Show the user what we think the host address is
	log.Infof("This engine will join UCP and advertise itself with host address %s - If this is incorrect, please specify an alternative address with the '--host-address' flag", config.OrcaHostAddress)

	orcaURL, err := url.Parse(c.String("url"))
	if err != nil {
		return 1, fmt.Errorf("Malformed UCP connection URL: %s", err)
	}

	fingerprint := c.String("fingerprint")
	httpClient, err := certs.GetTofuClient(fingerprint, orcaURL.Host)
	if err != nil {
		log.Debug("Failed to connect to UCP controller")
		return 1, err
	}

	// Get Admin credentials from environment
	orcaUser := os.Getenv("UCP_ADMIN_USER")
	orcaPass := os.Getenv("UCP_ADMIN_PASSWORD")

	if orcaUser == "" {
		return 1, errors.New("You must set UCP_ADMIN_USER and UCP_ADMIN_PASSWORD in the environment to proceed")
	}

	log.Info("Verifying your system is compatible with UCP")

	if containers != nil {
		log.Info("Removing old UCP containers")
		if err := ec.RemoveOrcaContainers(containers); err != nil {
			return 1, fmt.Errorf("Failed to remove an existing UCP container %s", err)
		}
	}

	// Check ports *after* we've shut down any running orca containers
	if c.Bool("replica") {
		if err := ec.CheckPorts(orcaconfig.RequiredPorts); err != nil {
			return 1, err
		}
	}

	if c.Bool("fresh-install") {
		// TODO - Check for controller replica, and if detected, try to
		// unwind it so we don't wedge the cluster
		if err := utils.CleanupVolumes(c); err != nil {
			return 1, err
		}
	}

	token, err := orcaclient.Login(httpClient, c.String("url"), orcaUser, orcaPass)
	if err != nil {
		log.Debug("Failed to login")
		return 1, err
	}

	// Get local node info which contains its swarm node ID.
	info, err := ec.GetClient().Info(context.Background())
	if err != nil {
		log.Error("unable to get local node info")
		return 1, err
	}

	nodeConfig, err := orcaclient.DoJoin(info.Swarm.NodeID, orcaUser, token, c.String("passphrase"), c.Bool("replica"), c.Bool("external-ucp-ca") || c.Bool("external-server-cert"), httpClient, orcaURL)
	if err != nil {
		log.Debug("Failed to get CSR signed by UCP")
		return 1, err
	}
	if len(nodeConfig.SwarmArgs) == 0 {
		return 1, fmt.Errorf("UCP server did not return swarm join arguments, unable to proceed")
	}

	if nodeConfig.Warnings != "" {
		// The controller needs to tell the user something urgent...
		log.Warning(nodeConfig.Warnings)
		time.Sleep(4 * time.Second)
	}

	log.Debug(nodeConfig.SwarmArgs)
	// Parse the experimental flag from SwarmArgs

	if nodeConfig.SwarmArgs[0] == "--experimental" {
		config.SwarmExperimental = true
		nodeConfig.SwarmArgs = nodeConfig.SwarmArgs[1:]
	}

	swarmArgs := client.FilterSwarmArgs(nodeConfig.SwarmArgs)
	log.Debugf("Applicable SwarmArgs: %v", swarmArgs)

	// TODO - Test that we can see swarm.  If not, then we've got some networking
	//        issue we need to warn the user about.  Detect the difference between non-accessible
	//        vs. cert failures to tease out classes of misconfiguration.

	if err := deployContainers(ec, nodeConfig, swarmArgs, c.Bool("replica"), orcaUser, orcaPass); err != nil {
		return 1, err
	}

	/* TODO - Either Need to get the swarm manager back from Orca, or Orca needs to allow swarm certs to connect

	// Now check to see if the node appeared in the orca swarm
	if err := client.WaitForNewNode(fmt.Sprintf("https://%s", XXX MISSING DATA XXX),
		proxyEndpoint, config.AliveCheckTimeout); err != nil {

		return 1, fmt.Errorf("Your new node (%s) didn't appear in the UCP Swarm.  Make sure your hosts have network connectivity", proxyEndpoint)
	}
	*/

	return 0, nil
}

func deployContainers(ec *client.EngineClient, nodeConfig *orca.NodeConfiguration, swarmArgs []string, isReplica bool, adminUsername, adminPassword string) error {
	// Then spin up the swarm nodes for this host
	log.Info("Starting local swarm containers")

	proxyEndpoint, err := utils.DeployProxyContainer(ec)
	if err != nil {
		return err
	}

	// The first swarm argument should be a URL for the KV store.
	kvStoreURL, err := url.Parse(swarmArgs[0])
	if err != nil {
		log.Error("unable to parse first swarm argument as a kv store URL")
		return err
	}

	// Deploy Swarm join
	if err := ec.StartSwarmJoin(kvStoreURL.String(), proxyEndpoint, swarmArgs[1:]); err != nil {
		log.Error("Failed to start swarm")
		return err
	}

	if !isReplica {
		// We're done setting this up as a general-purpose node.
		return nil
	}

	// This node should be setup as a UCP controller replica.
	return deployControllerReplicaContainers(ec, nodeConfig, adminUsername, adminPassword, proxyEndpoint)
}

func deployControllerReplicaContainers(ec *client.EngineClient, nodeConfig *orca.NodeConfiguration, adminUsername, adminPassword, proxyEndpoint string) error {

	// Should not happen if the system is behaving properly
	if len(nodeConfig.KvStore) == 0 {
		return fmt.Errorf("Unable to set up replica without master KV store")
	}

	log.Info("Starting UCP Controller replica containers")

	chDone := make(chan bool)
	defer func() {
		chDone <- true
		log.Debug("closing channel for kv lock")
		close(chDone)
	}()
	kvStoreURL, err := utils.DeployKVContainer(ec, nodeConfig.KvStore[0], chDone)
	if err != nil {
		return err
	}

	log.Debugf("Secondary KV started at %s", kvStoreURL)

	// Deploy Swarm Manager
	managerEndpoint, err := utils.DeploySwarmManagerContainer(ec, kvStoreURL.String())
	if err != nil {
		return err
	}

	orcaEndpoint := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.OrcaPort)

	// Record the swarm manager in the KV store so Orca knows what all of them are
	if err := config.AddSwarmManager(kvStoreURL, orcaEndpoint, managerEndpoint, proxyEndpoint); err != nil {
		return err
	}

	// Deploy the two CAs
	if err := utils.DeployCAContainers(ec); err != nil {
		return err
	}

	// Deploy Auth service containers.
	// Need to get IP addresses of other controller nodes.
	controllers, err := config.GetControllers(kvStoreURL)
	if err != nil {
		log.Errorf("unable to list controllers: %s", err)
		return err
	}

	if len(controllers) == 0 {
		return fmt.Errorf(`Failed to discover other controller instances`)
	}

	peerAddrs := make([]string, len(controllers))
	for i, peerIP := range orcatypes.GetIPsFromControllers(controllers) {
		peerAddrs[i] = net.JoinHostPort(peerIP, fmt.Sprintf("%d", orcaconfig.AuthStorePeerPort))
	}

	if err := utils.DeployAuthStoreContainer(ec, peerAddrs...); err != nil {
		return err
	}
	if err := utils.DeployAuthAPIContainer(ec); err != nil {
		return err
	}
	if err := utils.DeployAuthWorkerContainer(ec); err != nil {
		return err
	}

	if err := config.UpdateUCPServiceAuthConfig(kvStoreURL); err != nil {
		return fmt.Errorf("unable to update auth configuration: %s", err)
	}

	// Finally, deploy the UCP server replica.
	return utils.DeployControllerContainer(ec, kvStoreURL.String(), managerEndpoint, proxyEndpoint, orcaEndpoint)
}
