package util

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
)

var lettersAndDigits = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

// RandStringBytes generates a random string of length `n` containing only of letters and digits
func RandStringBytes(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersAndDigits[r.Intn(len(lettersAndDigits))]
	}
	return string(b)
}

// DockerTrustClient is a utility that shells out to the docker client binary in a DIND container
// in order to test the docker client and pushing and pulling images, because it's important to
// make sure that Docker's particular Notary+Registry+Auth handling works.  (engine-api does not
// include Notary integration).
//
// This utility spins up a DIND container using the engine-api so that it can use the docker client
// within it to run DCT-enabled docker commands.
type DockerTrustClient struct {
	// these are provided
	registry      string
	dindImageName string
	client        *client.Client

	// this is generated
	dindContainerName string
	dindContainerID   string
	configDir         string
}

// NewDockerTrustClient returns a DockerTrustClient object which spins up a new DIND container
// and execs into that container in order to connect to a registry and perform trusted
// pull/push operations
func NewDockerClientWithTrust(client *client.Client, registry, dindImageName string, registryCA []byte) (*DockerTrustClient, error) {
	d := DockerTrustClient{
		registry:      registry,
		client:        client,
		dindImageName: dindImageName,

		dindContainerName: RandStringBytes(12),
		configDir:         "/dockerconfig",
	}

	// Create a tar archive that contains the TLS cert for the registry
	tlsBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tlsBuf)
	header := &tar.Header{
		Name: filepath.Join(strings.TrimPrefix(d.configDir, "/"), "tls", d.registry, "ca.crt"),
		Mode: 0666,
		Size: int64(len(registryCA)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}

	if _, err := tw.Write(registryCA); err != nil {
		return nil, err
	}

	ctx := context.Background()

	// pull DIND image
	readcloser, err := client.ImagePull(ctx, d.dindImageName, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer readcloser.Close()
	// we need to give the pull time to download
	dec := json.NewDecoder(readcloser)
	for {
		var i interface{}
		if err := dec.Decode(&i); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	// start DIND container
	resp, err := client.ContainerCreate(
		ctx,
		&container.Config{
			Cmd:   []string{"--debug", "--insecure-registry", registry},
			Image: d.dindImageName,
			Env: []string{
				"DOCKER_CONTENT_TRUST=1",
				"DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE=password",
				"DOCKER_CONTENT_TRUST_ROOT_PASSPHRASE=password",
				"DOCKER_CONFIG=" + d.configDir,
			},
		},
		&container.HostConfig{Privileged: true},
		&network.NetworkingConfig{},
		d.dindContainerName,
	)
	if err != nil {
		d.Cleanup()
		return nil, err
	}

	d.dindContainerID = resp.ID
	if err := client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		d.Cleanup()
		return nil, err
	}

	if err := d.client.CopyToContainer(ctx, d.dindContainerName, "/", tlsBuf, types.CopyToContainerOptions{}); err != nil {
		d.Cleanup()
		return nil, err
	}

	// block until the container is ready
	for i := 0; ; i++ {
		if _, err := d.rawExecute("docker version"); err == nil {
			break
		} else {
			if i >= 60 {
				return nil, err
			}
			time.Sleep(1 * time.Second)
		}
	}

	return &d, nil
}

func (d *DockerTrustClient) rawExecute(cmd string) ([]byte, error) {
	ctx := context.Background()

	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", cmd},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}
	createResponse, err := d.client.ContainerExecCreate(ctx, d.dindContainerName, execConfig)
	if err != nil {
		return nil, err
	}

	// get output
	attachResp, err := d.client.ContainerExecAttach(ctx, createResponse.ID, execConfig)
	if err != nil {
		return nil, err
	}
	defer attachResp.Close()
	output, err := ioutil.ReadAll(attachResp.Reader)
	if err != nil {
		return nil, err
	}

	// now get exit code
	inspectResp, err := d.client.ContainerExecInspect(ctx, createResponse.ID)
	if err != nil {
		return nil, err
	}

	if inspectResp.Running {
		return nil, fmt.Errorf("exec command %s is still running", execConfig.Cmd)
	}

	if inspectResp.ExitCode != 0 {
		return output, fmt.Errorf("%s\n%s\nexit code=%d", cmd, string(output), inspectResp.ExitCode)
	}

	return output, nil
}

// Login executes a docker login in the DIND container
func (d *DockerTrustClient) Login(username, password string) ([]byte, error) {
	loginCmd := fmt.Sprintf("docker login -u %s -p %s %s", username, password, d.registry)
	return d.rawExecute(loginCmd)
}

// Pull executes a docker pull in the DIND container, with or without trust on
func (d *DockerTrustClient) Pull(image string, contentTrustOn bool) ([]byte, error) {
	pullCmd := "docker pull "
	if !contentTrustOn {
		pullCmd = pullCmd + "--disable-content-trust "
	}
	return d.rawExecute(pullCmd + image)
}

// Pull executes a docker push in the DIND container, with or without trust on
func (d *DockerTrustClient) Push(image string, contentTrustOn bool) ([]byte, error) {
	pushCmd := "docker push "
	if !contentTrustOn {
		pushCmd = pushCmd + "--disable-content-trust "
	}
	return d.rawExecute(pushCmd + image)
}

func (d *DockerTrustClient) Tag(oldImage, newImage string) ([]byte, error) {
	return d.rawExecute(fmt.Sprintf("docker tag %s %s", oldImage, newImage))
}

// ClearTrustData removes all trust data in the DIND container
func (d *DockerTrustClient) ClearConfigData() ([]byte, error) {
	clearCmd := fmt.Sprintf("rm -rf %s/trust %s/config.json", d.configDir, d.configDir)
	return d.rawExecute(clearCmd)
}

// Cleanup destroys the DIND container
func (d *DockerTrustClient) Cleanup() error {
	ctx := context.Background()
	if err := d.client.ContainerKill(ctx, d.dindContainerID, ""); err != nil {
		return err
	}
	return d.client.ContainerRemove(ctx, d.dindContainerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
}
