package utils

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/samalba/dockerclient"
)

type LocalMachine struct {
	name       string
	dockerHost string
	tlsConfig  *tls.Config
	ip         string // Cache so we don't have to look it up so much
	internalip string // Cache so we don't have to look it up so much
}

// Wire up for local machine usage
func NewLocalMachine() (Machine, error) {
	// Load up the local environment settings, or the local sock

	host := os.Getenv("DOCKER_HOST")
	certPath := os.Getenv("DOCKER_CERT_PATH")
	m := &LocalMachine{}

	if host == "" {
		m.name = "localhost"
		m.dockerHost = "unix:///var/run/docker.sock"
	} else {
		m.name = host
		m.dockerHost = host
		caFilename := filepath.Join(certPath, "ca.pem")
		certFilename := filepath.Join(certPath, "cert.pem")
		keyFilename := filepath.Join(certPath, "key.pem")
		// Load up the certs
		cert, err := tls.LoadX509KeyPair(certFilename, keyFilename)
		if err != nil {
			return nil, err
		}
		caCert, err := ioutil.ReadFile(caFilename)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		m.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}
	c, err := m.GetClient()
	if err != nil {
		return nil, err
	}
	m.internalip, err = GetHostAddress(c)
	if err != nil {
		return nil, err
	}
	if machineURL, err := url.Parse(m.dockerHost); err == nil && machineURL.Host != "" {
		m.ip = strings.Split(machineURL.Host, ":")[0]
	} else {
		// Local connection just use the internal IP
		m.ip = m.internalip
	}
	log.Infof("%s external IP: %s", m.name, m.ip)
	log.Infof("%s internal IP: %s", m.name, m.internalip)
	return m, nil
}

func (m *LocalMachine) GetName() string {
	return m.name
}

func (m *LocalMachine) GetDockerHost() string {
	return m.dockerHost
}

// Return a docker client connected to the machine
func (m *LocalMachine) GetClient() (*dockerclient.DockerClient, error) {
	return dockerclient.NewDockerClientTimeout(m.dockerHost, m.tlsConfig, Timeout, nil)
}

func (m *LocalMachine) GetEngineAPI() (*client.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: m.tlsConfig,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   Timeout,
	}
	version := "" // Always take the latest swarm API
	return client.NewClient(m.dockerHost, version, httpClient, nil)
}

// Remove the machine after the tests have completed
func (m *LocalMachine) Remove() error {
	log.Info("Local Machine not deleted")

	// This is a little hacky, but to make sure we've cleaned up everything for subsequent tests
	// we go ahead and do an uninstall, then purge things explicitly
	c, err := m.GetClient()
	if err != nil {
		return err
	}

	info, err := c.InspectContainer("ucp-controller")
	if err != nil {
		log.Warn("Unable to find ucp-controller during machine remove - something may have gone wrong...")
		return nil
	}
	id := info.Config.Labels["com.docker.ucp.InstanceID"]
	log.Infof("Uninstalling %s", id)
	err = RunBootstrapper(c, []string{"uninstall", "--preserve-images", "--id", id}, []string{}, []string{})
	if err != nil {
		log.Warnf("Failure during uninstall - ignoring: %s", err)
	}

	_ = c.RemoveContainer("ucp", true, false)
	// TODO - probably belongs in the specific test case...
	_ = c.RemoveContainer("selenium-firefox", true, false)
	return nil
}

func (m *LocalMachine) Stop() error {
	log.Error("Attempt to stop local machine!  This test shoudl be skipped for local mode")
	return fmt.Errorf("Stopping a local machine is not supported")
}

func (m *LocalMachine) Start() error {
	log.Error("Attempt to start local machine!  This test shoudl be skipped for local mode")
	return fmt.Errorf("Starting a local machine is not supported")
}

// Return the public IP of the machine
func (m *LocalMachine) GetIP() (string, error) {
	return m.ip, nil
}

// Get the internal IP (useful for join operations)
func (m *LocalMachine) GetInternalIP() (string, error) {
	return m.internalip, nil
}

// Get the contents of a specific file on the engine
func (m *LocalMachine) CatHostFile(hostPath string) ([]byte, error) {
	return catHostFile(m, hostPath)
}

// Get the content of a directory as a tar file from the engine
func (m *LocalMachine) TarHostDir(hostPath string) ([]byte, error) {
	return tarHostDir(m, hostPath)
}

func (m *LocalMachine) MachineSSH(command string) (string, error) {
	// TODO - might be worth while to implement with a local command execution?
	return "", fmt.Errorf("LocalMachine doesn't support SSH commands")
}

// Borrowed from bootstrap code (TODO - refactor to a common utility!)
func GetHostAddress(client *dockerclient.DockerClient) (string, error) {
	// TODO pull busybox, or switch this to use one of our built-in images
	cfg := &dockerclient.ContainerConfig{
		Image:        "busybox",
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   []string{"ip", "route", "get", "8.8.8.8"},
	}
	containerId, err := client.CreateContainer(cfg, "", nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create test for host address %s", err)
	}
	defer client.RemoveContainer(containerId, true, true)
	hostConfig := &dockerclient.HostConfig{
		NetworkMode: "host",
	}
	// Start the container
	if err := client.StartContainer(containerId, hostConfig); err != nil {
		log.Debugf("Failed to launch port test container: %s", err)
		return "", err
	}
	defer client.StopContainer(containerId, 5)

	// XXX Some sort of race - if we try to get the logs too quickly, we get EOF without
	//     actually getting any output at all
	time.Sleep(100 * time.Millisecond)

	// Gather the output
	reader, err := client.ContainerLogs(containerId, &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return "", err
	}

	addr := ""
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Read Error:", err)
			return "", err
		}
		offset := strings.Index(line, "src ")
		if offset < 0 {
			continue
		}
		addr = strings.TrimSpace(strings.SplitAfter(line[offset+4:], " ")[0])
		break
	}

	if addr == "" {
		return "", fmt.Errorf("Unable to determine host address")
	}
	return addr, nil
}
