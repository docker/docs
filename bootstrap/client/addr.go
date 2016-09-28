package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"

	"github.com/docker/orca/config"
)

var (
	ErrNoHostAddress = errors.New("Failed to determine host IP address.  Try setting UCP_HOST_ADDRESS")
)

// Figure out the host address
func (c *EngineClient) GetHostAddress() (string, error) {

	// Allow external caller to specify what their real address is for "unusual" routing setups
	addr := os.Getenv("UCP_HOST_ADDRESS")
	if addr != "" {
		return addr, nil
	}

	imageName, err := config.GetContainerImage("addr")
	if err != nil {
		return "", err
	}
	log.Debugf("Trying to determine host address using image %s", imageName)

	// Start up a dummy container with the specified ports punched through so we can fail fast
	cfg := &container.Config{
		Image:        imageName,
		AttachStdout: true,
		AttachStderr: true,
		Entrypoint:   strslice.StrSlice{"ip", "route", "get", "8.8.8.8"},
	}
	hostConfig := &container.HostConfig{
		NetworkMode: "host",
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	resp, err := c.client.CreateContainer(cfg, hostConfig, "")
	if err != nil {
		return "", fmt.Errorf("Failed to create test for host address %s", err)
	}
	containerId := resp.ID

	// Start the container
	if err := c.client.StartContainer(containerId); err != nil {
		log.Debugf("Failed to launch port test container: %s", err)
		return "", err
	}

	// XXX Some sort of race - if we try to get the logs too quickly, we get EOF without
	//     actually getting any output at all
	time.Sleep(100 * time.Millisecond)

	// TODO - this needs to be refactored next
	// Gather the output
	reader, err := c.client.ContainerLogs(containerId, types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}

	addr = ""
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		//log.Debug(line)
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

	// Cleanup
	c.client.StopContainer(containerId, 5) // Do we care about errors?
	c.client.RemoveContainer(containerId, true, true)
	if addr == "" {
		return "", ErrNoHostAddress
	}
	return addr, nil
}
