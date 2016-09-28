package util

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/samalba/dockerclient"

	log "github.com/Sirupsen/logrus"
)

func RunContainerWithOutput(client *dockerclient.DockerClient, image string, args, env []string) (string, error) {
	hostConfig := dockerclient.HostConfig{}

	cfg := &dockerclient.ContainerConfig{
		Image:      image,
		Cmd:        args,
		Env:        env,
		HostConfig: hostConfig,
	}
	containerId, err := client.CreateContainer(cfg, "", nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create container %s: %s", image, err)
	}
	defer client.RemoveContainer(containerId, true, true)

	log.Debugf("Launching %s with args:%v env:%v", image, args, env)
	if err := client.StartContainer(containerId, &hostConfig); err != nil {
		return "", fmt.Errorf("Failed to start container: %s", err)
	}
	defer client.StopContainer(containerId, 5)

	reader, err := client.ContainerLogs(containerId, &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	})

	buffer := new(bytes.Buffer)

	if _, err = stdcopy.StdCopy(buffer, buffer, reader); err != nil {
		log.Debug("cannot read logs from logs reader")
		return "", err
	}

	info, err := client.InspectContainer(containerId)
	if err != nil {
		return buffer.String(), fmt.Errorf("Failed to inspect container after completion: %s", err)
	}
	if info.State == nil {
		return buffer.String(), fmt.Errorf("Container didn't finish!")
	}

	if info.State.ExitCode != 0 {
		return buffer.String(), fmt.Errorf("Container exited with %d", info.State.ExitCode)
	}

	return buffer.String(), nil
}
