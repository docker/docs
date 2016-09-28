package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"

	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	orcautils "github.com/docker/orca/utils"
)

func DeployKVContainer(ec *client.EngineClient, existingStore string, chDone chan bool) (kvStoreURL *url.URL, err error) {
	kvCfg := &config.KVDeployCfg{
		Timeout: config.KVTimeout,
	}
	if len(existingStore) > 0 {
		// Try to get the configuration from the existing KV store
		existingUrl, err := url.Parse(existingStore)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse URL for existing KV store to load options: %s: %s")
		}
		kvCfg, err = config.GetKVDeployConfig(existingUrl)
		if err != nil {
			return nil, fmt.Errorf("Unable to contact existing KV store to load options: %s")
		}

		// Take a lock for the KV store so we don't try to add multiple at the same time and cause races
		kvStoreURL, err := url.Parse(existingStore)
		if err != nil {
			log.Debug("Failed to parse KV store url")
			return nil, err
		}
		kv, err := config.GetKV(kvStoreURL)
		if err != nil {
			log.Debug("Failed to connect to KV store")
			return nil, err
		}
		log.Info("Obtaining KV Lock")
		lockPath := path.Join(config.OrcaPrefix, "joinlock")
		if err := orcautils.GetKVLock(kv, lockPath, &kvstore.LockOptions{TTL: time.Duration(90) * time.Second}, chDone); err != nil {

			log.Debug("Failed to get new lock")
			return nil, err
		}
	}

	// Deploy KV, with mapped non-standard ports
	if err := ec.StartKv(existingStore, kvCfg); err != nil {
		log.Errorf(`Failed to start KV store.  Run "docker logs %s" for more details`, orcaconfig.OrcaKvContainerName)
		return nil, err
	}

	return &url.URL{
		Scheme: config.KvType,
		Host:   fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.KvPort),
	}, nil
}

func DeployProxyContainer(ec *client.EngineClient) (proxyEndpoint string, err error) {
	// Deploy the proxy with random exposed port
	if err := ec.StartProxy(); err != nil {
		log.Errorf(`Failed to start proxy.  Run "docker logs %s" for more details`, orcaconfig.OrcaProxyContainerName)
		return "", err
	}

	proxyEndpoint = fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.ProxyPort)
	log.Debugf("Proxy started on %s", proxyEndpoint)

	// Verify we can see the proxy
	err = client.WaitForEndpoint(fmt.Sprintf("https://%s", proxyEndpoint), config.SwarmNodeCertVolumeMount, config.AliveCheckTimeout)
	if err != nil {
		// As this is the first time we're trying to use the host address, give some more hints on failure.
		log.Errorf(`We were unable to communicate with proxy we just started at address %s.  Did you forget to open your firewall ports?  Did you forget to specify an alternate DNS server with the '--dns' flag?  If this address is incorrect, re-run the install using the '--host-address' option.    Run "docker logs %s" for more details from the proxy`, config.OrcaHostAddress, orcaconfig.OrcaProxyContainerName)
	}

	return proxyEndpoint, err
}

func DeploySwarmManagerContainer(ec *client.EngineClient, kvEndpoint string) (managerEndpoint string, err error) {
	managerEndpoint = fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.SwarmPort)
	if err := ec.StartSwarmManager(kvEndpoint, managerEndpoint, []string{}); err != nil {
		log.Errorf(`Failed to start swarm manager.  Run "docker logs %s" for more details`, orcaconfig.OrcaSwarmManagerContainerName)
		return "", err
	}
	log.Debugf("Swarm manager started on %s", managerEndpoint)

	// Verify we can see the swarm manager
	if err := client.WaitForEndpoint(fmt.Sprintf("https://%s", managerEndpoint), config.SwarmNodeCertVolumeMount, config.AliveCheckTimeout); err != nil {
		log.Errorf(`Swarm didn't come up within %v.  Run "docker logs %s" for more details`, config.AliveCheckTimeout, orcaconfig.OrcaSwarmManagerContainerName)
		return "", err
	}

	return managerEndpoint, nil
}

func DeployCAContainers(ec *client.EngineClient) error {
	if err := ec.StartCA(orcaconfig.OrcaCAContainerName, config.OrcaRootCAVolumeName, orcaconfig.OrcaCAPort); err != nil {
		log.Errorf(`Failed to start UCP certificate service.  Run "docker logs %s" for more details`, orcaconfig.OrcaCAContainerName)
		return err
	}
	if err := ec.StartCA(orcaconfig.OrcaSwarmCAContainerName, config.SwarmRootCAVolumeName, orcaconfig.SwarmCAPort); err != nil {
		log.Errorf(`Failed to start UCP swarm certificate service.  Run "docker logs %s" for more details`, orcaconfig.OrcaSwarmCAContainerName)
		return err
	}

	return nil
}

func DeployControllerContainer(ec *client.EngineClient, kvEndpoint, managerEndpoint, engineEndpoint, orcaEndpoint string) error {
	err := ec.StartOrcaServer(kvEndpoint, managerEndpoint, engineEndpoint)
	if err != nil {
		log.Errorf("Unable to start UCP Controller")
		return err
	}

	// Verify Orca server is up
	if err := client.WaitForOrca(fmt.Sprintf("https://%s", orcaEndpoint), config.AliveCheckTimeout); err != nil {
		SpewContainerLogs(ec, orcaconfig.OrcaControllerContainerName)

		log.Errorf("UCP didn't come up within %v.  Run \"docker logs %s\" for more details", config.AliveCheckTimeout, orcaconfig.OrcaControllerContainerName)
		return err
	}

	return nil
}

func SpewContainerLogs(ec *client.EngineClient, containerName string) {
	logReader, err2 := ec.ContainerLogs(containerName,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Tail:       "15",
		})
	if err2 != nil {
		log.Error(err2)
	} else {
		_, err3 := io.Copy(os.Stderr, logReader)
		if err3 == nil {
			log.Error(err3)
		}
	}
}

func DeployAuthStoreContainer(ec *client.EngineClient, peerAddrs ...string) error {
	if err := ec.StartAuthStore(peerAddrs...); err != nil {
		log.Errorf(`Failed to start Auth Data Store.  Run "docker logs %s" for more details`, orcaconfig.AuthStoreContainerName)
		return err
	}

	peerAddr := net.JoinHostPort(config.OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthStorePeerPort))
	clientAddr := net.JoinHostPort(config.OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthStorePort))

	// Wait for the server to be accepting connections on both the client
	// driver and cluster peer ports.
	for _, addr := range []string{peerAddr, clientAddr} {
		if err := client.WaitForTLSConn(addr, config.AuthStoreCertsVolumeMount, config.AuthStoreAliveCheckTimeout); err != nil {
			log.Errorf(`Auth Data Store didn't come up within %v.  Run "docker logs %s" for more details`, config.AliveCheckTimeout, orcaconfig.AuthStoreContainerName)
			return err
		}
	}

	// Perform necessary migrations of the database to the current version.
	// RethinkDB is schemaless, but this job will create tables, indexes,
	// and configure replication for each table.
	if err := ec.RunAuthSyncDB(false, clientAddr); err != nil {
		log.Errorf(`Failed to sync Auth Data Store.  Run "docker logs %s" for more details`, orcaconfig.AuthSyncDBContainerName)
		return err
	}

	return nil
}

func DeployAuthAPIContainer(ec *client.EngineClient) error {
	if err := ec.StartAuthAPI(); err != nil {
		log.Errorf(`Failed to start Auth API server.  Run "docker logs %s" for more details`, orcaconfig.AuthAPIContainerName)
		return err
	}

	// Wait for the API server to be responsive to API calls. There is no
	// PING endpoint (yet) so just use the "List Accounts" endpoint.
	authAPIAddr := net.JoinHostPort(config.OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthAPIPort))
	authAPIEndpoint := fmt.Sprintf("https://%s/enzi/v0/accounts", authAPIAddr) // We should get an AUTHENTICATION_REQUIRED response.
	if err := client.WaitForAuthServer(authAPIEndpoint, config.AuthAPICertsVolumeMount, http.StatusUnauthorized, config.AliveCheckTimeout); err != nil {
		log.Errorf(`Auth API server didn't come up within %v.  Run "docker logs %s" for more details`, config.AliveCheckTimeout, orcaconfig.AuthAPIContainerName)
		return err
	}

	return nil
}

func DeployAuthWorkerContainer(ec *client.EngineClient) error {
	if err := ec.StartAuthWorker(); err != nil {
		log.Errorf(`Failed to start Auth Worker server.  Run "docker logs %s" for more details`, orcaconfig.AuthWorkerContainerName)
		return err
	}

	// Wait for the Worker server to be responsive to API calls. The ping
	// endpoint returns no content.
	authWorkerAddr := net.JoinHostPort(config.OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthWorkerPort))
	authWorkerEndpoint := fmt.Sprintf("https://%s/v0/ping", authWorkerAddr)
	if err := client.WaitForAuthServer(authWorkerEndpoint, config.AuthWorkerCertsVolumeMount, http.StatusNoContent, config.AliveCheckTimeout); err != nil {
		log.Errorf(`Auth Worker server didn't come up within %v.  Run "docker logs %s" for more details`, config.AliveCheckTimeout, orcaconfig.AuthWorkerContainerName)
		return err
	}

	return nil
}

func DeployUCPAgentService(ec *client.EngineClient) error {
	// copy the key to the volume
	certPath := filepath.Join(config.SwarmRootCAVolumeMount, "ucp-instance-key.pem")
	data := []byte(config.OrcaInstanceKey)
	if err := ioutil.WriteFile(certPath, data, os.FileMode(0600)); err != nil {
		return fmt.Errorf("error writing cert to %s: %s", certPath, err)
	}

	username := os.Getenv("UCP_ADMIN_USER")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("UCP_ADMIN_PASSWORD")
	if password == "" {
		password = "orca"
	}
	credentials := []byte(fmt.Sprintf("%s\n%s", username, password))
	credPath := filepath.Join(config.SwarmRootCAVolumeMount, "ucp-credentials")
	if err := ioutil.WriteFile(credPath, credentials, os.FileMode(0600)); err != nil {
		return fmt.Errorf("error writing credentials to %s: %s", credPath, err)
	}

	return ec.StartUCPAgentService()
}
