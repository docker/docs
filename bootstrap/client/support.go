package client

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func (c *EngineClient) SupportDump() (int, error) {
	imageName, err := orcaconfig.GetContainerImage("dsinfo")
	if err != nil {
		return 1, err
	}
	containerConfig := &container.Config{
		Image: imageName,
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			"/boot:/boot",
			"/var/run/docker.sock:/var/run/docker.sock",
			"/var/lib/docker:/var/lib/docker",
			"/var/log:/var/log",
			"/etc/sysconfig:/etc/sysconfig",
			"/etc/default:/etc/default",
		},
	}

	id := ""
	if resp, err := c.client.CreateContainer(containerConfig, hostConfig, ""); err != nil {
		return 1, fmt.Errorf("Failed to create container: %s", err)
	} else {
		id = resp.ID
	}

	// Also wire up a signal handler in case we're killed
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func() {
		s := <-ch
		log.Errorf("Signal %v received, aborting", s)
		c.client.RemoveContainer(id, true, true)
		os.Exit(1)
	}()

	resp, err := c.client.ContainerAttach(id, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return 1, fmt.Errorf("Failed to attach to container: %s", err)
	}
	stream := resp.Reader

	if err := c.client.StartContainer(id); err != nil {
		return 1, fmt.Errorf("Failed to start container: %s", err)
	}
	defer c.client.RemoveContainer(id, true, false)

	doneChan := make(chan int)
	go func() {
		res, err := c.client.ContainerWait(id)
		if err != nil {
			log.Errorf("Failed to wait for support container: %s", err)
		}
		doneChan <- res
	}()

	go func() {
		// Temporarily change the log level since stdcopy can sometimes spew
		// TODO - fix stdcopy upstream so there's a public log ref that can be quieted down
		oldLevel := log.GetLevel()
		log.SetLevel(log.InfoLevel)
		defer log.SetLevel(oldLevel)
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream); err != nil {
			log.Errorf("Failed to stream logs: %s\n", err)
			return
		}
	}()

	select {
	case exitCode := <-doneChan:
		log.Debugf("Container '%s' exited correctly", id)
		return exitCode, nil
	case <-time.After(time.Duration(config.AliveCheckTimeout) * time.Second):
		return 1, fmt.Errorf("Container '%s' timed out while getting the support log", id)
	}

}
