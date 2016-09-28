package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mostly used for little CLI tools
func GetClientFromEnv() (*dockerclient.DockerClient, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	if dockerHost == "" {
		return dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	}
	caFilename := path.Join(dockerCertPath, "ca.pem")
	certFilename := path.Join(dockerCertPath, "cert.pem")
	keyFilename := path.Join(dockerCertPath, "key.pem")
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
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return dockerclient.NewDockerClient(dockerHost, tlsConfig)
}

func GetServerURLFromClient(client *dockerclient.DockerClient) string {
	return fmt.Sprintf("https://%s", client.URL.Host)
}

type CLIHandle struct {
	DirPath    string
	DockerHost string
	Binaries   []string
}

func CLIPresent() bool {
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	}
	return true
}

func (h *CLIHandle) Cleanup() {
	// Consider preserving if PRESERVE_TEST_MACHINE is set
	if err := os.RemoveAll(h.DirPath); err != nil {
		log.Errorf("Failed to remove %s %s", h.DirPath, err)
	}
}

func (h *CLIHandle) RunCLICommand(args []string) (string, error) {
	log.Debugf("Running %v", args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = h.DirPath
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CERT_PATH=%s", h.DirPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_HOST=tcp://%s:443", h.DockerHost))
	cmd.Env = append(cmd.Env, "DOCKER_TLS_VERIFY=1")
	//log.Debugf("Running %#v", cmd) // Mostly useless noise, but if something strange happens this might be useful...
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(string(out))
		return "", fmt.Errorf("%s - %s", err, string(out))
	}
	return string(out), nil
}

func NewCLIHandle(serverURL, username, password string) (*CLIHandle, error) {
	zr, err := GetBundle(serverURL, username, password)
	if err != nil {
		return nil, err
	}

	// Make the temp dir and expand the cert bundle there
	dirname, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	u, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	// Figure out what docker binaries we have available
	primary, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}
	commands := []string{primary}
	matches, err := filepath.Glob(primary + "-[1-9]*")
	if err == nil && len(matches) > 0 {
		commands = append(commands, matches...)
	}

	log.Infof("Detected the following docker binaries: %v", commands)

	h := &CLIHandle{
		DirPath:    dirname,
		DockerHost: strings.Split(u.Host, ":")[0],
		Binaries:   commands,
	}

	log.Debugf("Writing bundle out to %s", h.DirPath)
	for _, zf := range zr.File {

		src, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		var data []byte
		data, err = ioutil.ReadAll(src)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(h.DirPath, zf.Name), data, 0600); err != nil {
			return nil, err
		}
	}

	return h, err
}

// The following routines are test routines that can be used in different contexts

func TestSimpleCLI(t *testing.T, serverURL, username, password string, primaryVersion *version.Version) {
	// Might want to soften and skip if this is a problem...
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	type scenario struct {
		e string   //Expected output
		c []string //command
	}
	scenarios := []scenario{
		scenario{
			e: "it worked",
			c: []string{"run", "--rm", "busybox", "echo", "it worked"},
		},
		scenario{
			e: "it worked",
			c: []string{"run", "--rm", "-i", "busybox", "echo", "it worked"},
		},
		// Can't do "-it" -- it errors out with "cannot enable tty mode on non tty input"

		scenario{
			e: "nobody",
			c: []string{"run", "--rm", "-u", "nobody", "busybox", "whoami"},
		},
		// TODO - Add more scenarios here as we bump into corner cases and bugs
	}

	// Inconsistently fails on 1.10 on some AWS AMIs, lock to 1.11+
	minimumTtyVersion, _ := version.NewVersion("1.11.0")
	if !primaryVersion.LessThan(minimumTtyVersion) {
		scenarios = append(scenarios, scenario{
			e: "it worked",
			c: []string{"run", "--rm", "-t", "busybox", "echo", "it worked"},
		})
	} else {
		log.Infof("Skipping tty test, version %s too low", primaryVersion)
	}
	for _, cmd := range ch.Binaries {
		for _, s := range scenarios {
			fullCmd := append([]string{cmd}, s.c...)
			log.Infof("Running %v", fullCmd)

			output, err := ch.RunCLICommand(fullCmd)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(output, s.e), "Failed: %s -> %s NOT expected %s", fullCmd, output, s.e)
		}
	}
}

func TestPrivilegedCLI(t *testing.T, serverURL, username, password string) {
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	expected := "it worked"

	for _, cmd := range ch.Binaries {

		log.Debug("running docker CLI with privileged flag")
		output, err := ch.RunCLICommand([]string{cmd, "run", "--rm", "--privileged", "busybox", "echo", expected})
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, expected))
		log.Debugf("It contained expected output: %s", output)
	}
}

// CopyToFromContainer performs a copyFrom and a copyTo operation on a specific container.
// Both operations are performed even if copyFrom fails.
func CopyToFromContainer(serverURL, username, password, containerID string) error {
	var globalErr error
	ch, err := NewCLIHandle(serverURL, username, password)
	defer ch.Cleanup()

	// Perform a CopyFromContainer operation: copy /bin/rm to the host
	// TODO: find a way to clean the copied rm binary from the host
	_, err = ch.RunCLICommand([]string{ch.Binaries[0], "cp", containerID + ":/bin/rm", ch.DirPath})
	if err != nil {
		globalErr = err
	}

	// Perform a CopyToContainer operation: copy the host's rm binary to container:/rm-copy
	_, err = ch.RunCLICommand([]string{ch.Binaries[0], "cp", "/bin/rm", containerID + ":/rm-copy"})
	if err != nil {
		globalErr = err
	}

	return globalErr
}

func TestBridgeNetwork(t *testing.T, serverURL, username, password string) {
	log.Info("Attempting to create and run with a user defined bridge network")
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	networkName := "testNetwork"
	expected := "it worked"

	log.Debugf("creating network %s", networkName)
	output, err := ch.RunCLICommand([]string{ch.Binaries[0], "network", "create", "--driver", "bridge", networkName})
	assert.Nil(t, err)
	networkID := strings.TrimSpace(output)
	log.Debugf("Created with ID %s", networkID)
	// Ignore failures attempting to clean up
	defer ch.RunCLICommand([]string{"network", "rm", networkID}) // Use the ID to prevent possible clashes

	listedName := ""
	output, err = ch.RunCLICommand([]string{ch.Binaries[0], "network", "ls"})
	assert.Nil(t, err)
	re := regexp.MustCompile(`\S+\s+(\S+)\s+(\S+)`)
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			if matches[2] == "bridge" && strings.Contains(matches[1], networkName) {
				listedName = matches[1]
				break
			}
		}
	}
	require.NotEqual(t, listedName, "")

	// Now run a container with that network
	// Note: we don't run this over all commands since some older CLIs don't support `--net` in this way
	output, err = ch.RunCLICommand([]string{ch.Binaries[0], "run", "--rm", "--net", listedName, "busybox", "echo", expected})
	assert.Nil(t, err)
	require.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)
}

func TestContainerListCLI(t *testing.T, serverURL, username, password string) {
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	for _, cmd := range ch.Binaries {
		// start a container that exits instantly to compare against
		output, err := ch.RunCLICommand([]string{cmd, "run", "-d", "busybox", "echo", "hello world"})
		assert.Nil(t, err)
		exitedContainerId := strings.TrimSpace(output)
		log.Debugf("Started container %s", exitedContainerId)
		defer ch.RunCLICommand([]string{ch.Binaries[0], "rm", "-f", exitedContainerId})

		// start a long-lived container
		output, err = ch.RunCLICommand([]string{cmd, "run", "-d", "busybox", "sleep", "60"})
		assert.Nil(t, err)
		containerId := strings.TrimSpace(output)
		log.Debugf("Started container %s", containerId)
		defer ch.RunCLICommand([]string{ch.Binaries[0], "rm", "-f", containerId})

		// look for the exited one
		shortId := exitedContainerId[:6]
		output, err = ch.RunCLICommand([]string{cmd, "ps", "-a", "-f", "status=exited"})
		assert.Nil(t, err)
		log.Debugf("Looking for container %s %s", exitedContainerId, shortId)
		assert.True(t, strings.Contains(output, shortId))
		assert.False(t, strings.Contains(output, containerId[:6]))

		// look for the running one
		shortId = containerId[:6]
		output, err = ch.RunCLICommand([]string{cmd, "ps", "-a", "-f", "status=running"})
		assert.Nil(t, err)
		log.Debugf("Looking for container %s %s", containerId, shortId)
		assert.True(t, strings.Contains(output, shortId))
		assert.False(t, strings.Contains(output, exitedContainerId[:6]))
	}
}

func WaitForCLIReadiness(serverURL, username, password string, timeout time.Duration) error {
	startTime := time.Now()
	for startTime.Add(timeout).After(time.Now()) {
		time.Sleep(2 * time.Second)
		ch, err := NewCLIHandle(serverURL, username, password)
		if err != nil {
			return err
		}
		output, err := ch.RunCLICommand([]string{ch.Binaries[0], "run", "--rm", "busybox", "echo", "it worked"})
		if err != nil {
			if strings.Contains(err.Error(), "No healthy nodes available in the cluster") {
				continue
			}
			return err
		}
		if strings.TrimSpace(output) != "it worked" {
			return fmt.Errorf("unexpected output: %s", output)
		}
		return nil
	}
	return fmt.Errorf("Timeout while waiting for CLI to be ready")
}
