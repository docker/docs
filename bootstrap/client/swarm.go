package client

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/version"
)

// TODO - Refactor these, since there's a fair bit of duplicate code...

// Start the docker proxy
func (c *EngineClient) StartProxy() error {
	log.Debug("Starting docker proxy")

	portMap := make(map[nat.Port]struct{})
	portMap["2376/tcp"] = struct{}{}
	bindingMap := nat.PortMap{
		nat.Port("2376/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", orcaconfig.ProxyPort),
			},
		},
	}
	// Swarm nodes
	mounts := []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		fmt.Sprintf("%s:%s:ro", config.SwarmNodeCertVolumeName, config.CertDir),
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaProxyContainerName)
	if err != nil {
		return err
	}
	cfg := &container.Config{
		Image:        imageName,
		ExposedPorts: portMap,
		Env: []string{
			fmt.Sprintf("SSL_CA=%s", filepath.Join(config.CertDir, config.CAFilename)),
			fmt.Sprintf("SSL_CERT=%s", filepath.Join(config.CertDir, config.CertFilename)),
			fmt.Sprintf("SSL_KEY=%s", filepath.Join(config.CertDir, config.KeyFilename)),
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaProxyContainerName,
		},
		Cmd: []string{
			"proxy",
		},
	}
	hostConfig := &container.HostConfig{
		Binds:        mounts,
		PortBindings: bindingMap,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},

		// TODO - probably want more...
	}

	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaProxyContainerName)
	if err != nil {
		return err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch proxy: %s", err)
		return err
	}

	return nil
}

// Remove extranious args from the swarm args returned by the orca server
func FilterSwarmArgs(args []string) []string {
	res := []string{}
	skipNext := false
	for _, entry := range args {
		if skipNext {
			skipNext = false
			continue
		} else if entry == "join" {
			continue
		} else if entry == "manage" {
			continue
		} else if strings.HasPrefix(entry, "--advertise") {
			if !strings.Contains(entry, "=") {
				skipNext = true
			}
			continue
		}
		res = append(res, entry)
	}
	return res
}

// TODO - Replace KV with "args" list so we can make this a little more generic once the proxy/join/manager are combined
// Start the swarm join
func (c *EngineClient) StartSwarmJoin(kvEndpoint, proxy string, extraArgs []string) error {

	// Swarm certs
	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.SwarmNodeCertVolumeName, config.CertDir),
	}

	// XXX Refactor common goop...
	log.Debug("Starting swarm join")
	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaSwarmJoinContainerName)
	if err != nil {
		return err
	}

	cmd := append([]string{
		"join",
		"--discovery-opt", fmt.Sprintf("kv.cacertfile=%s", filepath.Join(config.CertDir, config.CAFilename)),
		"--discovery-opt", fmt.Sprintf("kv.certfile=%s", filepath.Join(config.CertDir, config.CertFilename)),
		"--discovery-opt", fmt.Sprintf("kv.keyfile=%s", filepath.Join(config.CertDir, config.KeyFilename)),
		"--discovery-opt", fmt.Sprintf("kv.path=%s", config.KvPath),
		"--advertise", proxy,
		kvEndpoint,
	}, extraArgs...)

	if config.SwarmExperimental {
		cmd = append([]string{"--experimental"}, cmd...)
	}

	cfg := &container.Config{
		Image: imageName,
		Cmd:   strslice.StrSlice(cmd),
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			fmt.Sprintf("%s.node", config.OrcaLabelPrefix):       "true",
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaSwarmJoinContainerName,
		},
		// TODO - probably want more...
	}
	hostConfig := &container.HostConfig{
		Binds: mounts,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		// TODO - probably want more...
	}
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaSwarmJoinContainerName)
	if err != nil {
		return err
	}
	containerId := resp.ID
	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch swarm join: %s", err)
		return err
	}
	return nil
}

func (c *EngineClient) StartSwarmManager(kvEndpoint, manager string, extraArgs []string) error {
	log.Debug("Starting swarm manager")

	// Swarm certs
	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.SwarmNodeCertVolumeName, config.CertDir),
	}

	portMap := make(map[nat.Port]struct{})
	portMap["2375/tcp"] = struct{}{}
	bindingMap := nat.PortMap{
		nat.Port("2375/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", orcaconfig.SwarmPort),
			},
		},
	}

	strategy := "spread"
	if config.SwarmBinpack {
		strategy = "binpack"
	} else if config.SwarmRandom {
		strategy = "random"
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaSwarmManagerContainerName)
	if err != nil {
		return err
	}

	cmd := append([]string{
		"manage",
		"--tlsverify",
		"--tlscacert", filepath.Join(config.CertDir, config.CAFilename),
		"--tlscert", filepath.Join(config.CertDir, config.CertFilename),
		"--tlskey", filepath.Join(config.CertDir, config.KeyFilename),
		"--replication",
		"--discovery-opt", fmt.Sprintf("kv.cacertfile=%s", filepath.Join(config.CertDir, config.CAFilename)),
		"--discovery-opt", fmt.Sprintf("kv.certfile=%s", filepath.Join(config.CertDir, config.CertFilename)),
		"--discovery-opt", fmt.Sprintf("kv.keyfile=%s", filepath.Join(config.CertDir, config.KeyFilename)),
		"--discovery-opt", fmt.Sprintf("kv.path=%s", config.KvPath),
		"--advertise", manager,
		"--strategy", strategy,
		kvEndpoint,
	}, extraArgs...)

	if config.SwarmExperimental {
		cmd = append([]string{"--experimental"}, cmd...)
	}

	cfg := &container.Config{
		Image:        imageName,
		ExposedPorts: portMap,
		Cmd:          strslice.StrSlice(cmd),
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaSwarmManagerContainerName,
		},
	}
	hostConfig := &container.HostConfig{
		Binds:        mounts,
		PortBindings: bindingMap,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
		// TODO - probably want more...
	}
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaSwarmManagerContainerName)
	if err != nil {
		return err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Errorf("Failed to launch proxy: %s", err)
		return err
	}

	return nil
}

func (c *EngineClient) StartUCPAgentService() error {
	agentImage, err := orcaconfig.GetContainerImage(orcaconfig.OrcaAgentContainerName)
	if err != nil {
		return fmt.Errorf("Can't find container image for %s: %s", orcaconfig.OrcaAgentContainerName, err)
	}

	ucpInstanceLabelKey := fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix)
	ucpVersionLabelKey := fmt.Sprintf("%s.version", config.OrcaLabelPrefix)
	env := []string{
		"IMAGE_VERSION=" + orcaconfig.ImageVersion,
		"UCP_INSTANCE_ID=" + config.OrcaInstanceID,
		fmt.Sprintf("SWARM_PORT=%d", orcaconfig.SwarmPort),
		fmt.Sprintf("CONTROLLER_PORT=%d", orcaconfig.OrcaPort),
		"DNS=" + strings.Join(config.DNS, ","),
		"DNS_OPT=" + strings.Join(config.DNSOpt, ","),
		"DNS_SEARCH=" + strings.Join(config.DNSSearch, ","),
	}
	if log.GetLevel() == log.DebugLevel {
		env = append(env, "DEBUG=1")
	}
	_, err = c.ServiceCreate(swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "ucp-agent",
			Labels: map[string]string{
				ucpInstanceLabelKey: config.OrcaInstanceID,
				ucpVersionLabelKey:  version.Version,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: agentImage,
				Labels: map[string]string{
					ucpInstanceLabelKey: config.OrcaInstanceID,
					ucpVersionLabelKey:  version.Version,
					// Note: compose labels omitted so this doesn't show up
					// Since tasks get recycled frequently, many exited containers
					// will show up and lead to unnecessary concern
				},
				Command: []string{"/ucp-agent", "agent"},
				Env:     env,
				Mounts: []swarm.Mount{
					{
						Type:     swarm.MountTypeBind,
						Source:   "/var/run/docker.sock",
						Target:   "/var/run/docker.sock",
						ReadOnly: false,
					},
				},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyConditionAny,
			},
		},
		Mode: swarm.ServiceMode{
			Global: &swarm.GlobalService{},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism: 1,
			Delay:       2 * time.Second, // TODO fine-tune
		},
	})

	return err
}

// XXX: below is stolen from swarm/create.go from docker/docker, wish that
// engine-api provided a higher-level implementation
func (c *EngineClient) CreateSwarmV2Swarm(listenAddr string) (string, error) {
	req := swarm.InitRequest{
		ListenAddr:      fmt.Sprintf("0.0.0.0:%d", orcaconfig.SwarmGRPCPort), // TODO consider exposing flag for listen during install
		AdvertiseAddr:   listenAddr,
		ForceNewCluster: false,
	}
	log.Debugf("Swarm init: %#v", req)
	name, err := c.client.SwarmInit(req)
	if err != nil {
		if strings.Contains(err.Error(), "must specify a listening address") {
			return name, fmt.Errorf("Your system has multiple addresses.  You must specify a --host-address to use")
		}
		return name, err
	}

	log.Info("Initializing a new swarm")
	return name, nil
}

func (c *EngineClient) LeaveSwarmV2Swarm(force bool) error {
	err := c.client.SwarmLeave(force)
	if err != nil {
		return err
	}

	log.Info("Left swarm.")
	return nil
}

func (c *EngineClient) ServiceCreate(spec swarm.ServiceSpec) (types.ServiceCreateResponse, error) {
	return c.client.ServiceCreate(spec)
}

func (c *EngineClient) ServiceList() ([]swarm.Service, error) {
	return c.client.ServiceList(types.ServiceListOptions{})
}

func (c *EngineClient) ServiceRemove(service string) error {
	return c.client.ServiceRemove(service)
}
