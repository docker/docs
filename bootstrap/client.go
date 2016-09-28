package bootstrap

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
	version "github.com/docker/engine-api/types/versions"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

func DockerClientFromBundle(host, dockerCertPath, hubUsername, hubPassword string) (*dc.Client, error) {
	var client *http.Client
	caPath := filepath.Join(dockerCertPath, "ca.pem")
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", caPath, err)
	}
	trans, err := dtrutil.HTTPTransport(false, []string{string(ca)}, filepath.Join(dockerCertPath, "cert.pem"), filepath.Join(dockerCertPath, "key.pem"))
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %s", err)
	}
	client = &http.Client{
		Transport: trans,
	}

	hostPort := strings.Split(host, ":")
	if len(hostPort) == 1 {
		host = fmt.Sprintf("tcp://%s:443", host)
	} else {
		host = fmt.Sprintf("tcp://%s:%s", hostPort[0], hostPort[1])
	}
	headers := map[string]string{}
	if hubUsername != "" && hubPassword != "" {
		headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword)
	}
	return dc.NewClient(host, "v1.22", client, headers)
}

func DockerClientFromJWT(host, jwt string, client *http.Client, hubUsername, hubPassword string) (*dc.Client, error) {
	hostPort := strings.Split(host, ":")
	if len(hostPort) == 1 {
		host = fmt.Sprintf("tcp://%s:443", host)
	} else {
		host = fmt.Sprintf("tcp://%s:%s", hostPort[0], hostPort[1])
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", jwt),
	}
	if hubUsername != "" && hubPassword != "" {
		headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword)
	}
	return dc.NewClient(host, "v1.22", client, headers)
}

func DockerClientFromSocket(hubUsername, hubPassword string) (*dc.Client, error) {
	headers := map[string]string{}
	if hubUsername != "" && hubPassword != "" {
		headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword)
	}
	return dc.NewClient("unix:///var/run/docker.sock", "v1.22", nil, headers)
}

func (b BootstrapConfig) ImageExists(imageID string) (bool, error) {
	_, _, err := b.Client.ImageInspectWithRaw(context.Background(), imageID, false)
	if err != nil {
		if dc.IsErrImageNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("Couldn't inspect image '%s': %s", imageID, err)
	}
	return true, nil
}

// TODO: maybe use the real struct: https://godoc.org/github.com/docker/docker/pkg/jsonmessage#JSONMessage
type PullResponse struct {
	Error       string `json:"error"`
	ErrorDetail struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errorDetail"`
	ID     string `json:"id"`
	Status string `json:"status"`
	//ProgressDetail map[string]interface{} `json:"progressDetail"`
}

func (b BootstrapConfig) GetDockerClient() *dc.Client {
	return b.Client
}

func (bs BootstrapConfig) ConfirmNodeIsOkay() error {
	if flags.Unsafe {
		log.Info("Skipping safety check")
		return nil
	}
	log.Info("Checking if the node is okay to install on")
	ver, err := bs.DaemonVersion()
	if err != nil {
		return err
	}
	csParts := strings.Split(ver, "-cs")
	var csVersion int64
	if len(csParts) > 1 {
		// some clients have "special brews" of docker, such as 1.10.3-cs2-paypal-oauth2.
		// we need to parse _just_ the cs version here, so match via a regex
		re := regexp.MustCompile("^([\\d]+)")
		match := re.FindStringSubmatch(csParts[1])
		if len(match) != 2 {
			csVersion = 1
		}
		if csVersion, err = strconv.ParseInt(match[1], 10, 64); err != nil {
			return fmt.Errorf("Failed to parse cs version number %s", csParts[1])
		}
	}
	// if we don't use `csParts[0]` instead of `ver` for the first comparison, 1.11.2-cs is considered lower than 1.11.2 because it's treated as a pre-release
	if version.LessThan(csParts[0], "1.11.2") && !(version.Equal(csParts[0], "1.11.1") && csVersion >= 2) {
		// This version of docker has problems with running dtr and ucp on the same machine
		// because docker uses overlay networking and ucp does a hacky thing with etcd where
		// it tries to connect to an etcd that runs on the same engine
		// See https://github.com/docker/dhe-deploy/issues/1811 and https://github.com/docker/docker/issues/22486

		// check if ucp is running on the target node
		filter, err := filters.ParseFlag("name=ucp-controller", filters.NewArgs())
		if err != nil {
			return err
		}
		cList, err := bs.ContainerList(types.ContainerListOptions{All: true, Filter: filter})
		if err != nil {
			return err
		}
		for _, container := range cList {
			for _, name := range container.Names {
				log.Debugf("UCP controller name: %s", name)
				// the name starts with a slash, so we skip that slash before getting the first component
				node := strings.Split(name[1:], "/")[0]
				if node == bs.Node {
					return fmt.Errorf("UCP controller found on the DTR node. Refusing to install DTR because it might cause issues upon restarting. See https://github.com/docker/docker/issues/22486 and the DTR documentation. Use --%s to skip this test. This issue is fixed in versions of Docker greater than or equal to 1.11.2 or 1.11.1-cs2. Your version is %s.", flags.UnsafeFlag.Name, ver)
				}
			}
		}
	}
	log.Debugf("Node is okay because the version is %s", ver)
	return nil
}

// in the version of swarm used by ucp 1.1.0 there's no easy way to get the engine version
// for a particcular machine. This function hacks around that by using the docker socket to
// get the engine version.
// Note: to use this, the docker socket must be mounted at /var/run/docker.sock
func (bs BootstrapConfig) DaemonVersionHack() (string, error) {
	client, err := DockerClientFromSocket("", "")
	if err != nil {
		return "", err
	}
	version, err := client.ServerVersion(context.Background())
	if err != nil {
		return "", err
	}

	return version.Version, nil
}

func (bs BootstrapConfig) DaemonVersion() (string, error) {
	log.Debugf("Checking the docker version for %s", bs.Node)
	dc := bs.GetDockerClient()
	info, err := dc.Info(context.Background())
	if err != nil {
		return "", err
	}
	// example output that we need to parse
	//Filters: health, port, dependency, affinity, constraint
	//Nodes: 1
	// compooter: 172.17.0.1:12376
	//  └ Status: Healthy
	//  └ Containers: 15
	//  └ Reserved CPUs: 0 / 4
	//  └ Reserved Memory: 0 B / 8.094 GiB
	//  └ Labels: executiondriver=, kernelversion=4.4.5-1-docker-aufs, operatingsystem=Arch Linux, storagedriver=aufs
	//  └ Error: (none)
	//  └ UpdatedAt: 2016-04-27T03:50:41Z
	//  └ ServerVersion: 1.11.0-rc5
	//Cluster Managers: 1
	// 172.17.0.1: Healthy
	startedNodesList := false
	scanningNodeData := false
	numNodes := 0
	// currNode will become 0 the first time we reach a node
	currNode := -1
	currNodeName := ""
	for _, line := range info.SystemStatus {
		log.Debugf("Parsing UCP status line: %v", line)
		key := line[0]
		value := line[1]
		if scanningNodeData {
			if len(key) > 1 && key[0] == ' ' && key[1] == ' ' {
				// check the key and current node to see if it's the one we are looking for
				if currNodeName == bs.Node && key == "  └ ServerVersion" {
					return value, nil
				}
				continue
			} else {
				// we are on the first line after we finished scanning the node
				scanningNodeData = false
				// if that was the last node, we are out of nodes, but we didn't find it
				if currNode == numNodes-1 {
					log.Warn("Failed to determine docker version the easy way. Trying the hard way...")
					return bs.DaemonVersionHack()
				}
			}
		}
		// detect the start of a node's data
		if startedNodesList && len(key) > 1 && key[0] == ' ' && key[1] != ' ' {
			currNodeName = strings.TrimPrefix(key, " ")
			currNode += 1
			scanningNodeData = true
			continue
		}
		if key == "Nodes" {
			nodes, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return "", err
			}
			startedNodesList = true
			numNodes = int(nodes)
			continue
		}
	}
	return "", fmt.Errorf("Failed to find docker version for node %s", bs.Node)
}

func (b BootstrapConfig) VersionCheck() error {
	log.Debug("Checking for compatible engine version")

	ver, err := dtrutil.DockerVersionCheck(b.Client, deploy.DockerVersionRequired, deploy.UCPVersionRequired)
	if err != nil {
		return err
	}

	log.Debugf("Engine version %s is compatible", ver)
	return nil
}

func (bs BootstrapConfig) PullImages() error {
	for _, container := range containers.AllContainers {
		// we have to always try to pull the image because it might not exist everywhere or might be out of date
		err := bs.PullImageWithChecks(container.Image)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bs BootstrapConfig) PullImageWithChecks(image string) error {
	// if the image is dirty and already exists, we don't try to pull at all, otherwise we would pull an unknown image from another dev
	exists := false
	var err error
	if strings.HasSuffix(deploy.Version, "-dirty") {
		exists, err = bs.ImageExists(image)
		if err != nil {
			return err
		}
	}
	if !exists {
		// we have to always try to pull the image because it might not exist everywhere or might be out of date
		err := bs.PullImage(image)
		if err != nil {
			log.Warnf("Problem pulling '%s': %s", image, err)
		}
	}
	// check if we have the image now, whether we pulled or not
	exists, err = bs.ImageExists(image)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Failed to find pulled image %s locally; refusing to execute.", image)
	}
	return nil
}

func (b BootstrapConfig) PullImage(imageID string) error {
	options := types.ImagePullOptions{
		RegistryAuth: b.AuthConfig,
	}

	buf, err := b.Client.ImagePull(b.Context, imageID, options)
	if err != nil {
		return fmt.Errorf("Failed to pull image: %s", err)
	}
	defer buf.Close()
	// TODO: remove garbage comments
	//out, err := ioutil.ReadAll(buf)
	//if err != nil {
	//	return fmt.Errorf("Failed to read response from pull attempt: %s", err)
	//}
	//log.Warn("out: %s", string(out))
	decoder := json.NewDecoder(buf)
	for err == nil {
		response := PullResponse{}
		err = decoder.Decode(&response)
		if response.Error == "" {
			log.Info(response.Status)
		} else {
			return fmt.Errorf("Daemon error: %s", response.Error)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("Failed to parse response: %s", err)
	}
	return nil
}

func (b *BootstrapConfig) SetUCPClient(apiClient *ucpclient.APIClient) {
	b.UCPClient = apiClient
}

func (b BootstrapConfig) GetUCPClient() *ucpclient.APIClient {
	return b.UCPClient
}

func (b BootstrapConfig) RunContainer(cfg *container.Config, hostCfg *container.HostConfig, netCfg *network.NetworkingConfig, name string) (string, error) {
	resp, err := b.Client.ContainerCreate(context.Background(), cfg, hostCfg, netCfg, name)
	if err != nil {
		return "", err
	}
	if err = b.Client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return resp.ID, err
	}
	return resp.ID, err
}

func (b BootstrapConfig) ContainerCreate(cfg *container.Config, hostCfg *container.HostConfig, netCfg *network.NetworkingConfig, name string) (types.ContainerCreateResponse, error) {
	return b.Client.ContainerCreate(context.Background(), cfg, hostCfg, netCfg, name)
}

func (b BootstrapConfig) ContainerRemove(id string, options types.ContainerRemoveOptions) error {
	// we keep trying to delete the container until we are told it doesn't exist
	// it's not safe to block until delete returns when using swarm
	err := dtrutil.Poll(time.Second, 10, func() error {
		log.Debugf("Deleting container %s", id)
		err := b.Client.ContainerRemove(context.Background(), id, options)
		if err == nil {
			return fmt.Errorf("Timeout trying to delete container %s", id)
		}
		log.Debugf("Failed to remove container: %s", err)
		// ignore aufs race errors
		if (strings.Contains(err.Error(), "500 Internal Server Error: Driver") && strings.Contains(err.Error(), "failed to remove root filesystem")) ||
			strings.Contains(err.Error(), "500 Internal Server Error: Unable to remove filesystem") {
			log.Debugf("Ignored error, but retrying %s...", err)
			return err
		}
		// accept not found as success
		if IsNoSuchImageErr(err.Error()) {
			log.Debugf("Ignored error %s...", err)
			return nil
		}
		// if it's any other error, don't bother retrying
		panic(fmt.Errorf("Unknown error: %s", err))
	})

	return err
}

func (b BootstrapConfig) ContainerStart(containerID string) error {
	return b.Client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
}

func (b BootstrapConfig) ContainerStop(containerID string) error {
	sec := time.Second
	return b.Client.ContainerStop(context.Background(), containerID, &sec)
}

func (b BootstrapConfig) ContainerWait(containerID string) (int, error) {
	ctx := context.Background()
	return b.Client.ContainerWait(ctx, containerID)
}

type byString []string

func (a byString) Len() int           { return len(a) }
func (a byString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byString) Less(i, j int) bool { return a[i] < a[j] }

func strSliceContains(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

// ContainerDiff returns true if the container running with the given name was not started with this
// ContainerConfig
// Note that this does NOT cover every possible change
func (b BootstrapConfig) ContainerDiff(container containers.ContainerConfig) (bool, error) {
	state, err := b.ContainerInspect(container.Name)
	if err != nil {
		if IsNoSuchImageErr(err.Error()) {
			return true, nil
		}
		return false, err
	}

	// diff the environment variables (position indepedant)
	// NOTE: we can't detect that an env variable should be removed, so we shouldn't conditionally include env variables
	envs := b.genEnvs(container)
	for _, env := range envs {
		if env[:11] == "constraint:" {
			continue
		}
		envName := strings.Split(env, "=")[0]
		if envName == "ETCD_INITIAL_CLUSTER_STATE" || envName == "ETCD_INITIAL_CLUSTER" {
			continue
		}
		if !strSliceContains(env, state.Config.Env) {
			log.Debugf("comparing %s - missing env: %v in %v", container.Name, env, state.Config.Env)
			return true, nil
		}
	}

	// diff the exposed ports only
	portBindings := b.genPorts(container)

	allPorts := state.NetworkSettings.Ports
	exposedPorts := nat.PortMap{}
	for srcPort, bindings := range allPorts {
		if bindings != nil && len(bindings) != 0 {
			newBindings := []nat.PortBinding{}
			for _, binding := range bindings {
				binding.HostIP = "0.0.0.0"
				newBindings = append(newBindings, binding)
			}
			exposedPorts[srcPort] = newBindings
		}
	}
	if !reflect.DeepEqual(exposedPorts, portBindings) {
		log.Debugf("comparing %s - unequal port configs: %v and %v", container.Name, exposedPorts, portBindings)
		return true, nil
	}

	// diff the log configs
	if state.HostConfig.LogConfig.Type != container.LogConfig.Type || state.HostConfig.LogConfig.Config["syslog-address"] != container.LogConfig.Config["syslog-address"] {
		log.Debugf("comparing %s - unequal log configs: %v and %v", container.Name, state.HostConfig.LogConfig, container.LogConfig)
		return true, nil
	}

	return false, nil
}

func (b BootstrapConfig) genEnvs(container containers.ContainerConfig) []string {
	var envs []string
	for k, v := range container.Environment {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	envs = append(envs, container.Constraints...)

	if b.ReplicaID != "" {
		envs = append(envs, fmt.Sprintf("%s=%s", flags.ReplicaIDFlag.EnvVar, b.ReplicaID))
	}

	if container.Node != "" {
		envs = append(envs, fmt.Sprintf("constraint:node==%s", container.Node))
	}

	if len(container.ExcludeNodes) > 0 {
		for _, n := range container.ExcludeNodes {
			envs = append(envs, fmt.Sprintf("constraint:node!=%s", n))
		}
	}

	log.Debug("envs:")
	for _, env := range envs {
		log.Debugf("env: %s", env)
	}
	return envs
}

func (b BootstrapConfig) genPorts(container containers.ContainerConfig) nat.PortMap {
	portBindings := make(nat.PortMap)
	for _, port := range container.Ports {
		p := strings.Split(port, ":")
		pubPort := nat.Port(p[0])
		portBindings[pubPort] = append(portBindings[pubPort], nat.PortBinding{HostIP: "0.0.0.0", HostPort: p[1]})
	}
	return portBindings
}

func (b BootstrapConfig) ContainerCreateFromContainerConfig(containerConfig containers.ContainerConfig) (*types.ContainerCreateResponse, error) {
	envs := b.genEnvs(containerConfig)
	portBindings := b.genPorts(containerConfig)

	cfg := &container.Config{
		AttachStdin:  containerConfig.AttachStdin,
		AttachStdout: containerConfig.AttachStdout,
		AttachStderr: containerConfig.AttachStderr,
		OpenStdin:    containerConfig.OpenStdin,
		StdinOnce:    containerConfig.StdinOnce,
		Image:        containerConfig.Image,
		Entrypoint:   containerConfig.Entrypoint,
		Env:          envs,
		Labels:       map[string]string{},
		Tty:          containerConfig.Tty,
	}

	// Don't make the phase2 bootstrap look like an application in UCP
	if containerConfig.Name != deploy.BootstrapPhase2ContainerName && containerConfig.Name != deploy.BootstrapHelperContainerName && containerConfig.Name != deploy.OverlayTestContainer1Name && containerConfig.Name != deploy.OverlayTestContainer2Name {
		cfg.Labels["com.docker.compose.project"] = fmt.Sprintf("Docker Trusted Registry %s - (Replica %s)", deploy.ShortVersion, b.ReplicaID)
		cfg.Labels["com.docker.compose.service"] = strings.TrimSuffix(strings.TrimPrefix(containerConfig.Name, "dtr-"), fmt.Sprintf("-%s", b.ReplicaID))

		// One off is false, since this is not a one time use, burnable container
		cfg.Labels["com.docker.compose.oneoff"] = "False"

		// The reason this will always be 1 is that we have one container of each type per Replica and in this model each replica is one project
		cfg.Labels["com.docker.compose.container-number"] = "1"

		// These are usually the compose version used to deploy and a hash made from the state of the compose file (to signify whether the configuration needs to be restarted).
		// Leaving these commented out since we don't have equivalents.
		// cfg.Labels["com.docker.compose.version"] = ""
		// cfg.Labels["com.docker.compose.config-hash"] = ""
	}

	volumes := containerConfig.DumbVolumes
	for _, volume := range containerConfig.Volumes {
		volumes = append(volumes, volume.FormatForReplica(b.ReplicaID))
	}

	hostCfg := &container.HostConfig{
		Binds:         volumes,
		PortBindings:  portBindings,
		RestartPolicy: container.RestartPolicy{Name: containerConfig.Restart},
		LogConfig:     containerConfig.LogConfig,
	}

	// XXX: maybe we can just pass nil
	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	resp, err := b.ContainerCreate(cfg, hostCfg, netCfg, containerConfig.Name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create container '%s' from image '%s': %s", containerConfig.Name, containerConfig.Image, err)
	}
	containerId := resp.ID

	// XXX: consider adding the overlay network after starting the containers, not before.
	// that way we can guarantee that no extra bridge is created for host port mapping
	// This is not actually relevant right now because the only container with ports mapped
	// is nginx and it's not on the overlay network
	nodeName := b.GetNodeName()
	if containerConfig.Node != "" {
		nodeName = containerConfig.Node
	}
	for _, net := range containerConfig.Networks {
		networkName := net.Name
		if networkName == deploy.BridgeNetworkName && nodeName != "" {
			networkName = fmt.Sprintf("%s/%s", nodeName, networkName)
		}
		log.Debugf("Network set to: %s", networkName)
		endpointCfg := &network.EndpointSettings{
			Aliases: net.Aliases,
		}

		netId, err := b.NetworkID(networkName)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get network id for %s; err: %s", networkName, err)
		}
		err = b.NetworkConnect(netId, containerId, endpointCfg)
		if err != nil {
			return nil, fmt.Errorf("Couldn't attach %s to %s; err: %s", netId, containerId, err)
		}
	}

	log.Debugf("Checking for bridge network. Have node: %s", containerConfig.Node)
	// remove the default network if we added a network
	if len(containerConfig.Networks) > 0 {
		bridgeName := "bridge"
		if containerConfig.Node != "" {
			bridgeName = fmt.Sprintf("%s/%s", containerConfig.Node, bridgeName)
		}
		log.Debugf("Bridge name is: %s", bridgeName)
		netId, err := b.NetworkID(bridgeName)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get network id for %s; err: %s", bridgeName, err)
		}
		// XXX: do we want to force?
		err = b.Client.NetworkDisconnect(context.Background(), netId, containerId, false)
		if err != nil {
			return nil, fmt.Errorf("Couldn't detach %s from %s", netId, containerId)
		}
	}

	return &resp, nil
}

func (b BootstrapConfig) ContainerKill(containerID, signal string) error {
	return b.Client.ContainerKill(context.Background(), containerID, signal)
}

func (b BootstrapConfig) ContainerRunFromContainerConfig(image containers.ContainerConfig) (string, error) {
	log.Debugf("Running container '%s'", image.Name)

	resp, err := b.ContainerCreateFromContainerConfig(image)
	if err != nil {
		log.Errorf("Couldn't create container: %s", err)
		return "", err
	}

	if err = b.ContainerStart(resp.ID); err != nil {
		log.Errorf("Couldn't start container: %v, %s", image.Entrypoint, err)
		return "", err
	}

	return resp.ID, nil
}

func (b BootstrapConfig) StartDTRContainers(ignoreContainers map[string]bool, replicaID string) (int, error) {
	log.Debug("Starting DTR containers")
	containers, err := b.GetRunningDTRContainers(replicaID)
	if err != nil {
		return 1, err
	}
	for name, id := range containers {
		if _, ignore := ignoreContainers[name]; ignore {
			continue
		}
		log.Debugf("Starting up container: %s", name)
		if err := b.ContainerStart(id); err != nil {
			log.Warnf("Couldn't start container %s: %s", id, err)
			return 1, err
		}
	}
	return 0, nil
}

// XXX: deduplicate this
func (b BootstrapConfig) StopDTRContainers(ignoreContainers map[string]bool, replicaID string) error {
	log.Debug("Stopping DTR containers")
	containers, err := b.GetRunningDTRContainers(replicaID)
	if err != nil {
		return err
	}
	for name, id := range containers {
		if _, ignore := ignoreContainers[name]; ignore {
			continue
		}
		log.Debugf("Stopping container: %s", name)
		if err := b.ContainerStop(id); err != nil {
			log.Warnf("Couldn't stop container: %s: %s", id, err)
			return err
		}
	}
	return nil
}

// XXX: deduplicate this
func (b BootstrapConfig) RemoveDTRContainers(ignoreContainers map[string]bool, replicaID string) error {
	log.Debug("Removing DTR containers")
	containerz, err := b.GetRunningDTRContainers(replicaID)
	if err != nil {
		return err
	}
	for name, id := range containerz {
		if _, ignore := ignoreContainers[name]; ignore {
			continue
		}
		log.Debugf("Removing container: %s", name)
		options := types.ContainerRemoveOptions{
			Force: true,
		}
		if err := b.ContainerRemove(id, options); err != nil {
			if !IsNoSuchImageErr(err.Error()) {
				log.Warnf("Couldn't remove container: %s: %s", id, err)
				return err
			} else {
				log.Debugf("Couldn't remove container because it doesn't exist: %s: %s", id, err)
			}
		}
	}
	return nil
}

func (b BootstrapConfig) RemoveDTRVolumes(replicaID string, includingCA bool) error {
	log.Info("Removing replica volumes...")
	for _, vol := range containers.Volumes {
		if vol == containers.CAVolume && !includingCA {
			continue
		}
		volumeName := vol.ReplicaName(replicaID)
		if !flags.NoUCP {
			volumeName = fmt.Sprintf("%s/%s", b.GetNodeName(), volumeName)
		}
		if b.VolumeExists(volumeName) {
			log.Debugf("Removing volume '%s'", volumeName)
			if err := b.VolumeRemove(volumeName); err != nil {
				log.Errorf("Couldn't remove volume '%s'", volumeName)
				return err
			}
		}
	}
	return nil
}

func (b BootstrapConfig) GetRunningDTRContainers(replicaID string) (map[string]string, error) {
	runningContainers, err := b.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: false})
	if err != nil {
		log.Debugf("Can't get containers: %s", err)
		return nil, err
	}
	dtrContainers := map[string]string{}
	for _, runningContainer := range runningContainers {
		runningContainerName := runningContainer.Names[0]
		for _, containerConfig := range containers.AllContainers {
			dtrContainerName := containerConfig.ReplicaName(replicaID)
			if strings.Contains(runningContainerName, dtrContainerName) {
				dtrContainers[dtrContainerName] = runningContainer.ID
			}
		}
	}
	return dtrContainers, nil
}

func (b BootstrapConfig) VolumeExists(volumeID string) bool {
	_, err := b.Client.VolumeInspect(context.Background(), volumeID)
	if err != nil {
		return false
	}
	return true
}

func (b BootstrapConfig) VolumeCreate(name string) error {
	options := types.VolumeCreateRequest{
		Name: name,
	}

	_, err := b.Client.VolumeCreate(context.Background(), options)
	return err
}

func (b BootstrapConfig) VolumeRemove(volumeID string) error {
	return b.Client.VolumeRemove(context.Background(), volumeID)
}

func (b BootstrapConfig) NetworkInspect(netID string) (types.NetworkResource, error) {
	return b.Client.NetworkInspect(context.Background(), netID)
}

func (b BootstrapConfig) NetworkExists(networkID string) bool {
	log.Debugf("Looking for network: %s", networkID)
	_, err := b.Client.NetworkInspect(context.Background(), networkID)

	// XXX - client doesn't export the error, so we'll just assume that any error
	//       is a missing network
	if err != nil {
		log.Debugf("Error inspecting network: %s", err)
		return false
	}

	return true
}

func (b BootstrapConfig) NetworkID(name string) (string, error) {
	cfg, err := b.Client.NetworkInspect(context.Background(), name)
	if err != nil {
		return "", fmt.Errorf("Couldn't get network ID for network '%s'", name)
	}
	return cfg.ID, nil
}

func (b BootstrapConfig) CreateNetwork(name, driver string) (types.NetworkCreateResponse, error) {
	options := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         driver,
		EnableIPv6:     false,
	}
	// for overlay networks use a less common default subnet
	if driver == "overlay" {
		options.IPAM = network.IPAM{Config: []network.IPAMConfig{{Subnet: "10.1.0.0/24"}}}
	}

	resp, err := b.Client.NetworkCreate(context.Background(), name, options)
	if err != nil {
		if driver == "overlay" {
			if strings.Contains(err.Error(), "has conflicts in the host while running in host mode") {
				log.Warnf("Failed to create overlay network. The overlay network's subnet is conflicting with a subnet already allocated to an interface. We'll try again with a different subnet. If that fails too, check the documentation for how to manually create the dtr-ol network. The error was: %s", err)
				options.IPAM = network.IPAM{Config: []network.IPAMConfig{{Subnet: "10.2.0.0/24"}}}
				resp, err := b.Client.NetworkCreate(context.Background(), name, options)
				if err != nil {
					log.Errorf("Couldn't create network '%s': %s", name, err)
					return types.NetworkCreateResponse{}, err
				}
				log.Debugf("resp: %q", resp)
				return resp, nil
			} else if strings.Contains(err.Error(), "datastore for scope \"global\" is not initialized") {
				err = fmt.Errorf("Failed to create network '%s'. Overlay networking is not configured. Error: %s", name, err)
				return types.NetworkCreateResponse{}, err
			} else if strings.Contains(err.Error(), "could not get pools config from store: client: etcd cluster is unavailable or misconfigured") {
				err = fmt.Errorf("Failed to create network '%s'. Overlay networking is not configured. Did you forget to restart the daemon after re-installing UCP?. Error: %s", name, err)
				return types.NetworkCreateResponse{}, err
			}
		}
		err = fmt.Errorf("Couldn't create network '%s': %s", name, err)
		return types.NetworkCreateResponse{}, err
	}
	log.Debugf("resp: %q", resp)
	return resp, nil
}

func (b BootstrapConfig) NetworkConnect(networkID, containerID string, endpointCfg *network.EndpointSettings) error {
	return b.Client.NetworkConnect(context.Background(), networkID, containerID, endpointCfg)
}

func (b BootstrapConfig) ContainerInspect(containerID string) (types.ContainerJSON, error) {
	return b.Client.ContainerInspect(context.Background(), containerID)
}

// returns the name of the node that a container is running on
func (b BootstrapConfig) ContainerNode(containerID string) (string, error) {
	log.Debugf("Inspecting '%s'", containerID)
	cData, err := b.ContainerInspect(containerID)
	if err != nil {
		log.Errorf("Couldn't inspect container '%s': %s", containerID, err)
		return "", err
	}
	log.Debugf("inspect data: %q", cData.Node)
	log.Debugf("node name is '%s'", cData.Node.Name)
	return cData.Node.Name, nil
}

func (b BootstrapConfig) ContainerList(options types.ContainerListOptions) ([]types.Container, error) {
	return b.Client.ContainerList(context.Background(), options)
}

func (b *BootstrapConfig) ValidateOrGetNewReplicaID(replicaID string) (string, error) {
	if replicaID == "" {
		return b.GetNewReplicaID()
	} else {
		_, err := hex.DecodeString(replicaID)
		if err != nil {
			return "", fmt.Errorf("Replica IDs must be hexadecimal: %s", err)
		}
		if len(replicaID) != deploy.ReplicaIDLen {
			return "", fmt.Errorf("Replica IDs must be %d characters long. This one is %d", deploy.ReplicaIDLen, len(replicaID))
		}
		return strings.ToLower(replicaID), nil
	}
}

func (b *BootstrapConfig) GetNewReplicaID() (string, error) {
	replicas, err := b.ListReplicas(true)
	if err != nil {
		return "", err
	}

	for {
		b := make([]byte, deploy.ReplicaIDLen/2)
		n, err := rand.Read(b)
		if n != deploy.ReplicaIDLen/2 {
			return "", fmt.Errorf("Failed to generate random replica id. Only generated %d random bytes", n)
		}
		if err != nil {
			return "", fmt.Errorf("Failed to generate random replica id: %s", err)
		}
		choice := hex.EncodeToString(b)

		if replicas.GetReplica(choice) == nil {
			return choice, nil
		}
	}
}

func (b *BootstrapConfig) ListReplicas(includingStopped bool) (Replicas, error) {
	opts := types.ContainerListOptions{
		All: includingStopped,
	}
	cList, err := b.ContainerList(opts)
	if err != nil {
		return nil, err
	}
	var replicas Replicas
	namePattern := regexp.MustCompile("dtr-registry-(.*)$")
	versionPattern := regexp.MustCompile(":(.*)$")
	for _, container := range cList {
		// XXX - this might be slow for thousands of containers
		for _, name := range container.Names {
			log.Debugf("Container name = %s", name)
			nameMatches := namePattern.FindStringSubmatch(name)
			if len(nameMatches) > 0 {
				containerJSON, err := b.ContainerInspect(container.ID)
				if err != nil {
					return nil, err
				}

				log.Debugf("Version name = %s", containerJSON.Config.Image)

				versionMatches := versionPattern.FindStringSubmatch(containerJSON.Config.Image)
				if len(versionMatches) == 0 {
					continue
				}
				replica := NewReplica(nameMatches[1], versionMatches[1])
				replicas = append(replicas, replica)
			}
		}
	}
	return replicas, nil
}

func (b *BootstrapConfig) ListReplicaIDs(includingStopped bool) ([]string, error) {
	replicas, err := b.ListReplicas(includingStopped)
	if err != nil {
		return nil, err
	}

	return replicas.ReplicaIDs(), nil
}

func (b BootstrapConfig) CreateNodeConstraintsFromPort(publicPorts []int) (*[]string, error) {
	nodeConstraints := make(map[string]bool)
	log.Debugf("public port constraints: %v", publicPorts)

	opts := types.ContainerListOptions{All: true}
	containers, err := b.ContainerList(opts)
	if err != nil {
		log.Errorf("Couldn't get a list of containers: %s", err)
		return nil, err
	}

	// iterate through each of the containers to see if there is a port conflict
	// if there is, inspect the container and extract the node name and add it to
	// our list of constraints
	for _, cont := range containers {
		for _, port := range cont.Ports {
			log.Debugf("found port: %d", port.PublicPort)
			for _, pubPort := range publicPorts {
				if pubPort == port.PublicPort {
					state, err := b.ContainerInspect(cont.ID)
					if err != nil {
						return nil, fmt.Errorf("Couldn't inspect container '%s' (using port %d): %s", cont.ID, pubPort, err)
					}
					if state.Node == nil {
						return nil, fmt.Errorf("Error inspecting container '%s': a container using port %d is running on an unknown node", cont.ID, pubPort)
					}
					if state.Node.Name != "" {
						log.Debugf("Adding constraint for node: %s", state.Node.Name)
						nodeConstraints[state.Node.Name] = true
					} else {
						log.Warnf("Couldn't for a node name for node '%s'", state.Node.ID)
					}
				}
			}
		}
	}
	nodes := []string{}
	for k := range nodeConstraints {
		nodes = append(nodes, k)
	}
	log.Debugf("node constraints: %s", nodes)
	return &nodes, nil
}

func (b BootstrapConfig) ExistingReplicaFlagPicker(prompt string, mustExist bool) (*Replica, error) {
	replica, err := b.ExistingReplicaPicker(flags.ExistingReplicaID, prompt, mustExist)
	if err != nil {
		return nil, err
	}
	flags.ExistingReplicaID = replica.ReplicaID
	return replica, nil
}

func (b BootstrapConfig) ExistingReplicaPicker(userInput string, prompt string, mustExist bool) (*Replica, error) {
	replicas, err := b.ListReplicas(false)
	if err != nil {
		return nil, err
	}

	// mustExist is set to false for the remove command so that a user can pick a non-existent replica as their "existing replica". They would want to do that if they cluster is broken and just want to destroy a replica.
	if len(replicas) == 0 && mustExist {
		log.Errorf("Did not find any potential Docker Trusted Registry replicas.")
		log.Errorf("You must use the 'install' command to first install a Docker Trusted Registry replica.")
		return nil, fmt.Errorf("Couldn't find any existing DTR replicas")
	}

	var pickedReplica *Replica
	replicaIDs := replicas.ReplicaIDs()
	pickedReplicaID := userInput
	if pickedReplicaID == "" {
		log.Infof("This cluster contains the replicas: %s", strings.Join(replicaIDs, " "))
		if len(replicaIDs) > 0 {
			pickedReplicaID = PromptString(fmt.Sprintf("%s [%s]: ", prompt, replicaIDs[0]), replicaIDs[0])
		} else {
			pickedReplicaID = PromptString(fmt.Sprintf("%s: ", prompt), "")
		}
		pickedReplica = replicas.GetReplica(pickedReplicaID)
		if pickedReplica == nil && mustExist {
			err = fmt.Errorf("Replica '%s' could not be found.", pickedReplicaID)
			return nil, err
		}
		if pickedReplicaID == "" {
			return nil, fmt.Errorf("Replica ID required.")
		}
	} else {
		pickedReplica = replicas.GetReplica(pickedReplicaID)
		if pickedReplica == nil && mustExist {
			err = fmt.Errorf("The replica '%s' you specified could not be found. Valid replicas: %s", pickedReplicaID, strings.Join(replicaIDs, " "))
			return nil, err
		}
	}
	return pickedReplica, nil
}

func (b *BootstrapConfig) SetReplicaID(id string) {
	b.ReplicaID = id
}

func (b *BootstrapConfig) GetReplicaID() string {
	return b.ReplicaID
}

func (b *BootstrapConfig) SetNodeName(nodeName string) {
	b.Node = nodeName
}

func (b *BootstrapConfig) GetNodeName() string {
	return b.Node
}

func (b BootstrapConfig) GetUsername() string {
	return b.Username
}

func (b BootstrapConfig) GetPassword() string {
	return b.Password
}

func (b BootstrapConfig) GetHost() string {
	return b.Host
}

func (b BootstrapConfig) GetTlsVerify() bool {
	return b.TlsVerify
}
