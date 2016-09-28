package ha_utils

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

	log "github.com/Sirupsen/logrus"
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
	log.Debugf("Running docker %v", args)
	cmd := exec.Command("docker", args...)
	cmd.Dir = h.DirPath
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CERT_PATH=%s", h.DirPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_HOST=tcp://%s:444", h.DockerHost))
	cmd.Env = append(cmd.Env, "DOCKER_TLS_VERIFY=1")
	log.Debugf("Running %#v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(string(out))
		return "", fmt.Errorf("%s - %s", err, string(out))
	}
	return string(out), nil
}

func NewCLIHandle(serverURL, username, password string) (*CLIHandle, error) {
	zr, err := GetUCPBundle(serverURL, username, password)
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
	h := &CLIHandle{
		DirPath:    dirname,
		DockerHost: strings.Split(u.Host, ":")[0],
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

func TestSimpleCLI(t require.TestingT, serverURL, username, password string) {
	// Might want to soften and skip if this is a problem...
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	expected := "it worked"

	log.Debug("running docker CLI")
	output, err := ch.RunCLICommand([]string{"run", "--rm", "busybox", "echo", expected})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)

	log.Debug("running docker CLI with -t")
	output, err = ch.RunCLICommand([]string{"run", "--rm", "-t", "busybox", "echo", expected})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)

	log.Debug("running docker CLI with -i")
	output, err = ch.RunCLICommand([]string{"run", "--rm", "-i", "busybox", "echo", expected})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)

	// Can't do "-it" -- it errors out with "cannot enable tty mode on non tty input"

	log.Debug("running docker CLI with -u nobody")
	output, err = ch.RunCLICommand([]string{"run", "--rm", "-u", "nobody", "busybox", "whoami"})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, "nobody"))
}

func TestPrivilegedCLI(t require.TestingT, serverURL, username, password string) {
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	expected := "it worked"

	log.Debug("running docker CLI with privileged flag")
	output, err := ch.RunCLICommand([]string{"run", "--rm", "--privileged", "busybox", "echo", expected})
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)
}

// CopyToFromContainer performs a copyFrom and a copyTo operation on a specific container.
// Both operations are performed even if copyFrom fails.
func CopyToFromContainer(serverURL, username, password, containerID string) error {
	var globalErr error
	ch, err := NewCLIHandle(serverURL, username, password)
	defer ch.Cleanup()

	// Perform a CopyFromContainer operation: copy /bin/rm to the host
	_, err = ch.RunCLICommand([]string{"cp", containerID + ":/bin/rm", ch.DirPath})
	if err != nil {
		globalErr = err
	}

	// Perform a CopyToContainer operation: copy the host's rm binary to container:/rm-copy
	_, err = ch.RunCLICommand([]string{"cp", "/bin/rm", containerID + ":/rm-copy"})
	if err != nil {
		globalErr = err
	}

	return globalErr
}

func TestBridgeNetwork(t require.TestingT, serverURL, username, password string) {
	log.Info("Attempting to create and run with a user defined bridge network")
	require.True(t, CLIPresent())

	ch, err := NewCLIHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	networkName := "testNetwork"
	expected := "it worked"

	log.Debugf("creating network %s", networkName)
	output, err := ch.RunCLICommand([]string{"network", "create", "--driver", "bridge", networkName})
	assert.Nil(t, err)
	networkID := strings.TrimSpace(output)
	log.Debugf("Created with ID %s", networkID)
	// Ignore failures attempting to clean up
	defer ch.RunCLICommand([]string{"network", "rm", networkID}) // Use the ID to prevent possible clashes

	listedName := ""
	output, err = ch.RunCLICommand([]string{"network", "ls"})
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
	output, err = ch.RunCLICommand([]string{"run", "--rm", "--net", listedName, "busybox", "echo", expected})
	assert.Nil(t, err)
	require.True(t, strings.Contains(output, expected))
	log.Debugf("It contained expected output: %s", output)
}
