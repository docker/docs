package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

type KVMember struct {
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	PeerURLs   []string `json:"peerURLs,omitempty"`
	ClientURLs []string `json:"clientURLs,omitempty"`
}
type KVMembers struct {
	Members []KVMember `json:"members"`
}

// NOTE: The given url should use HTTPS.
func GetKVMembers(client *http.Client, kvStoreURL *url.URL) *KVMembers {
	resp, err := client.Get(kvStoreURL.String() + "/v2/members")
	var members KVMembers
	if err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == 200 {
				if err := json.Unmarshal(body, &members); err == nil {
					return &members
				} else {
					log.Warnf("Failed to deserialize existing members from KV store: %s", err)
				}
			} else {
				log.Debugf("KV responded with error: %s", body)
			}
		} else {
			log.Debugf("KV body read error: %s", err)
		}
	} else {
		log.Debugf("KV connect failure: %s", err)
	}
	log.Warn("Failed to read members from KV store")
	return nil
}

// Start the specified container
func (c *EngineClient) StartKv(existingStore string, kvCfg *config.KVDeployCfg) error {
	// TODO - Take certs for the setup...

	log.Debug("Starting KV container")

	bindingMap := nat.PortMap{
		nat.Port("2379/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", orcaconfig.KvPort),
			},
		},
		nat.Port(fmt.Sprintf("%d/tcp", orcaconfig.KvPortPeer)): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", orcaconfig.KvPortPeer),
			},
		},
	}
	portMap := make(map[nat.Port]struct{})
	portMap["2379/tcp"] = struct{}{}
	portMap[nat.Port(fmt.Sprintf("%d/tcp", orcaconfig.KvPortPeer))] = struct{}{}

	// The endpoint of this local member
	localKvEndpoint := fmt.Sprintf("https://%s:%d", config.OrcaHostAddress, orcaconfig.KvPortPeer)

	initialClusterState := "new"

	peers := []string{}
	if existingStore != "" {
		initialClusterState = "existing"
		log.Debugf("Attemping to join etcd to an existing cluster at %s", existingStore)
		primaryKvURL, err := url.Parse(existingStore)
		if err != nil {
			return fmt.Errorf("Malformed KV url returned from Orca %s", existingStore)
		}

		kvURL := &url.URL{
			Scheme: "https",
			Host:   primaryKvURL.Host,
		}

		// get the peers, and append to peers
		// Note: this is very etcd specific
		// TODO - consider migrating this elsewhere into KV store specific routine
		tlsConfig, err := LoadCerts(config.SwarmKvCertVolumeMount)
		if err != nil {
			log.Debugf("Failed to load certs: %s", err)
			return err
		}
		client := &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
		members := GetKVMembers(client, kvURL)
		if members == nil {
			return fmt.Errorf("Failed to deserialize existing members from KV store")
		}
		// Process the member list and wire up our peers
		for _, m := range members.Members {
			if len(m.PeerURLs) == 0 {
				log.Warnf("Malformed KV store member list - empty peerURLs for %s", m.Name)
				continue
			}
			// TODO - What if there are multiple peerURLs for a member?
			//        Perhaps this should be a nested loop over all the URLs
			peers = append(peers, fmt.Sprintf("orca-kv-%s=%s", m.Name, m.PeerURLs[0]))
		}

		// call API to add this node as a new member
		log.Debugf("Adding new member %s to KV store %s", localKvEndpoint, kvURL)
		reqJson, err := json.Marshal(KVMember{PeerURLs: []string{localKvEndpoint}})
		if err != nil {
			log.Debug("Failed to generate member update payload")
			return err
		}
		resp, err := client.Post(kvURL.String()+"/v2/members", "application/json", bytes.NewBuffer(reqJson))
		if err != nil {
			log.Debug("Failed to send request")
			return err
		}
		if resp.StatusCode > 299 || resp.StatusCode < 200 { // Any 2xx response code is good!
			log.Debugf("Response code: %d", resp.StatusCode)
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				log.Errorf("Server response: %s", string(body))
				return fmt.Errorf("Failed to add member to KV store: %s", string(body))
			}
			return fmt.Errorf("Failed to add member to KV store")
		}
		// At this point the new member has been added, but the KV store needs us to come online ASAP
		// If we fail, and it was a single node, the KV store may be wedged since a dual-node
		// cluster requires both nodes to have quorum.
	}

	peers = append(peers, fmt.Sprintf("orca-kv-%s=https://%s:%d", config.OrcaHostAddress, config.OrcaHostAddress, orcaconfig.KvPortPeer))

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaKvContainerName)
	if err != nil {
		return err
	}
	cmd := strslice.StrSlice{
		"--data-dir", "/data",
		"--name", fmt.Sprintf("orca-kv-%s", config.OrcaHostAddress),
		"--listen-peer-urls", fmt.Sprintf("https://0.0.0.0:%d", orcaconfig.KvPortPeer),
		"--listen-client-urls", "https://0.0.0.0:2379",
		"--advertise-client-urls", fmt.Sprintf("https://%s:%d", config.OrcaHostAddress, orcaconfig.KvPort),
		"--initial-advertise-peer-urls", localKvEndpoint,
		"--initial-cluster", strings.Join(peers, ","),
		"--initial-cluster-state", initialClusterState,
		"--trusted-ca-file", filepath.Join(config.CertDir, config.CAFilename),
		"--cert-file", filepath.Join(config.CertDir, config.CertFilename),
		"--key-file", filepath.Join(config.CertDir, config.KeyFilename),
		"--client-cert-auth",
		"--peer-trusted-ca-file", filepath.Join(config.CertDir, config.CAFilename),
		"--peer-cert-file", filepath.Join(config.CertDir, config.CertFilename),
		"--peer-key-file", filepath.Join(config.CertDir, config.KeyFilename),
		"--peer-client-cert-auth",
		"--heartbeat-interval", fmt.Sprintf("%d", kvCfg.Timeout/10),
		"--election-timeout", fmt.Sprintf("%d", kvCfg.Timeout),
	}

	cfg := &container.Config{
		Hostname:     config.OrcaLocalName,
		Image:        imageName,
		ExposedPorts: portMap,
		Cmd:          cmd,
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaKvContainerName,
		},
		// TODO - probably want more...
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/data", config.OrcaKVVolumeName),
			fmt.Sprintf("%s:%s:ro", config.SwarmKvCertVolumeName, config.CertDir),
		},
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
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaKvContainerName)
	if err != nil {
		return err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch KV: %s", err)

		// TODO - If we were enabling replication, and this is the second node
		//        we may be leaving the cluster in a wedged state.  Consider trying
		//        to delete the new member (ourself) here
		return err
	}

	// check the KV store's actual health endpoint
	tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:   filepath.Join(config.SwarmNodeCertVolumeMount, config.CAFilename),
		CertFile: filepath.Join(config.SwarmNodeCertVolumeMount, config.CertFilename),
		KeyFile:  filepath.Join(config.SwarmNodeCertVolumeMount, config.KeyFilename),
	})
	if err != nil {
		return fmt.Errorf("Error making client to check KV store health: %s", err)
	}
	etcdclient := &http.Client{
		Timeout: config.AliveCheckTimeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	kvServer := fmt.Sprintf("https://%s:%d", config.OrcaHostAddress, orcaconfig.KvPort)
	if err := WaitForHealthyKv(etcdclient, kvServer, config.AliveCheckTimeout); err != nil {
		log.Errorf("error waiting for KV endpoint to be healthy: %s", err)
		return err
	}

	return nil
}

func (c *EngineClient) StartEtcdFixupContainer(cmd strslice.StrSlice) (string, error) {
	log.Debug("Starting KV Restore container")

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaKvContainerName)
	if err != nil {
		return "", err
	}

	cfg := &container.Config{
		Image: imageName,
		Cmd:   cmd,
		User:  "root", // We'll fix up permissions after the restore
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaKvRestoreContainerName,
		},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/data", config.OrcaKVVolumeName),
			fmt.Sprintf("%s:%s:ro", config.SwarmKvCertVolumeName, config.CertDir),
		},
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaKvRestoreContainerName)
	if err != nil {
		return "", err
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch KV Restore container: %s", err)
		return "", err
	}
	return containerId, nil
}

func (c *EngineClient) TakeKVBackup() error {
	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaKvBackupContainerName)
	if err != nil {
		return err
	}

	cfg := &container.Config{
		Hostname: config.OrcaLocalName,
		Image:    imageName,
		User:     "root", // We'll fix up permissions after the backup
		Entrypoint: strslice.StrSlice{"etcdctl", "backup", "--data-dir",
			"/data", "--backup-dir", "/backup-data"},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
			"com.docker.compose.container-number":                "1",
			"com.docker.compose.oneoff":                          "False",
			"com.docker.compose.project":                         "Docker Universal Control Plane " + config.OrcaInstanceID,
			"com.docker.compose.service":                         orcaconfig.OrcaKvBackupContainerName,
		},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/data", config.OrcaKVVolumeName),
			fmt.Sprintf("%s:%s:ro", config.SwarmKvCertVolumeName, config.CertDir),
			fmt.Sprintf("%s:/backup-data", config.OrcaKVBackupVolumeName),
		},
		AutoRemove: true,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.OrcaKvBackupContainerName)
	if err != nil {
		return err
	}
	defer c.client.RemoveContainer(resp.ID, true, true)

	// Start the container
	if err := c.client.StartContainer(resp.ID); err != nil {
		log.Debugf("Failed to launch KV container for etcdctl: %s", err)
		return err
	}

	// Wait for the etcd backup operation
	timeout := 5 * time.Minute
	started := time.Now()
	success := false

	log.Debug("Waiting for the etcd backup to complete")
	for !success && started.Add(timeout).After(time.Now()) {
		time.Sleep(2 * time.Second)
		cinfo, err := c.client.InspectContainer(resp.ID)
		if err != nil {
			log.Debug(err)
			return err
		}
		if cinfo.State.Status == "exited" {
			if cinfo.State.ExitCode != 0 {
				logs, err := c.GetContainerLogs(resp.ID)
				if err != nil {
					return fmt.Errorf("etcd backup was not successful: %s", err.Error())
				}
				log.Errorf(logs)
				return fmt.Errorf("etcd backup container exited with status code %d",
					cinfo.State.ExitCode)
			}
			success = true
		}
	}

	if !success {
		logs, err := c.GetContainerLogs(resp.ID)
		if err != nil {
			return fmt.Errorf("etcd backup was not successful: %s", err.Error())
		}

		log.Errorf(logs)
		return fmt.Errorf("etcd backup was not successful")
	}

	return nil
}

func (c *EngineClient) GetContainerLogs(containerID string) (string, error) {
	logOutput, err := c.client.ContainerLogs(containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		return "", fmt.Errorf("unable to retrieve logs: %s", err)
	}
	logBuf := new(bytes.Buffer)
	logBuf.ReadFrom(logOutput)
	return logBuf.String(), nil
}

// FindKvNonController finds a KV store without being on a controller node by
// checking to what swarm is connnected
func (c *EngineClient) FindKvNonController() (*url.URL, error) {
	// First look for the swarm-join container
	info, err := c.client.InspectContainer(orcaconfig.OrcaSwarmJoinContainerName)
	if err != nil {
		log.Debugf("Unable to inspect local %s: %s", orcaconfig.OrcaSwarmJoinContainerName, err)
		return nil, err
	}

	// The swarm-join container's command line ends with "etcd://ip:port", not
	// preceded by any flag (e.g. it's positional), so just grab the last flag
	// TODO: make this less fragile by using env vars to pass this to swarm
	cmdSlice := info.Config.Cmd
	return url.Parse(cmdSlice[len(cmdSlice)-1])
}

// GetKVStoreURL searches the given command list from a container config for
// a `--discovery` flag and parses its value as a URL.
func GetKVStoreURL(cmd strslice.StrSlice) (*url.URL, error) {
	// Check for missing
	for i, arg := range cmd {
		if !strings.HasPrefix(arg, "--discovery") {
			continue
		}

		var rawURL string
		if strings.Contains(arg, "=") {
			rawURL = strings.SplitN(arg, "=", 2)[1]
		} else {
			// The value is the next element in the command list.
			if i == len(cmd)-1 {
				// There is no next element.
				break
			}

			rawURL = cmd[i+1]
		}

		return url.Parse(rawURL)
	}

	return nil, fmt.Errorf("Failed to detect --discovery flag")
}

// If the local node is running a KV, return the URL to connect to it
func (c *EngineClient) FindKv() (*url.URL, error) {
	// First look for orca-controller
	info, err := c.client.InspectContainer(orcaconfig.OrcaControllerContainerName)
	if err != nil {
		log.Debugf("Unable to inspect local %s: %s", orcaconfig.OrcaControllerContainerName, err)
		return nil, err
	}

	return GetKVStoreURL(info.Config.Cmd)
}

// Detect if the local KV is in HA mode, and if so, remove it
func (c *EngineClient) DetectAndRemoveKVHANode(kvStoreURL *url.URL) error {
	tlsConfig, err := LoadCerts(config.SwarmKvCertVolumeMount)
	if err != nil {
		log.Errorf("Failed to load KV certs: %s", err)
		return err
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make an HTTP URL.
	kvURL := &url.URL{
		Scheme: "https",
		Host:   kvStoreURL.Host,
	}

	members := GetKVMembers(client, kvURL)
	if members == nil {
		return nil
	}
	log.Debugf("Detected %d members in the KV cluster", len(members.Members))
	if len(members.Members) <= 1 {
		return nil
	}

	// Loop through and try to find ourself
	localIP := strings.Split(kvURL.Host, ":")[0]
	myID := ""
	anotherNode := ""
	names := []string{} // Just to help troubleshooting if something goes wrong
	for _, m := range members.Members {
		names = append(names, m.Name)
		if strings.Contains(m.Name, localIP) {
			myID = m.Id
			continue
		}
		if len(m.ClientURLs) > 0 {
			anotherNode = m.ClientURLs[0]
		}
	}
	if myID == "" {
		return fmt.Errorf("Something went wrong - unable to find our local node %s in the list of peers %v",
			localIP, names)
	} else if anotherNode == "" {
		return fmt.Errorf("Failed to find another peer to send delete request to")
	}

	// Now do the delete via the peer
	deleteURL := anotherNode + "/v2/members/" + myID
	log.Debugf("Deleting %s", deleteURL)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		log.Errorf("Failed to build request to remove this member from the KV store")
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to remove this member from the KV store")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		// Wait for the etcd cluster to be healthy again since we just removed a member.
		if err := WaitForHealthyKv(client, anotherNode, config.AliveCheckTimeout); err != nil {
			log.Errorf("error waiting for kv endpoint: %s", err)
			return err
		}

		log.Debug("Successfully deleted node from KV store cluster")
		if len(members.Members) == 3 {
			log.Warning("You have reduced your cluster to 2 nodes, which provides reduced high-availability protection.  You should add at least one more replica node, or uninstall one more controller to disable HA.")
		}
		u, err := url.Parse(anotherNode)
		if err != nil {
			return fmt.Errorf("Malformed URL for secondary KV url: %s %s", anotherNode, err)
		}
		u.Scheme = "etcd"
		if err := config.RemoveSwarmManager(u, localIP); err != nil {
			log.Warningf("Failed to remove swarm manager registration from KV store: %s", err)
		}
		return nil
	}

	// Something went wrong
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		return fmt.Errorf("Failed to remove this member from the KV store: %s", body)
	} else {
		return err
	}
}
