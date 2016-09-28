package client

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func (c *EngineClient) stopPhase2(id string) {
	// Ignore errors, just try to clean up
	c.client.StopContainer(id, 5)
	c.client.RemoveContainer(id, true, true)
}

// Start the phase 2 bootstrap container
func (c *EngineClient) StartPhase2(args []string, privileged bool) (int, error) {
	//log.Debugf("Phase 1 looks like: %v", c.bootstrapper)
	//log.Debugf("Volumes: %v", c.bootstrapper.Volumes)
	var pidMode container.PidMode
	if privileged {
		pidMode = "host"
	}

	isTty := c.bootstrapper.Config.Tty // Mimic phase 1 setup

	// If there was a license file, pass it through
	for _, mount := range c.bootstrapper.HostConfig.Binds {
		if strings.Contains(mount, config.LicenseFile) || strings.Contains(mount, config.BackupFile) {
			config.Phase2VolumeMounts = append(config.Phase2VolumeMounts, mount)
			break
		}
	}

	cfg := &container.Config{
		Image:        c.bootstrapper.Image,
		Cmd:          strslice.StrSlice(args),
		Tty:          isTty,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,

		Env: append(c.bootstrapper.Config.Env, []string{
			// TODO - This could probably stand to have some refactoring so these values aren't
			//        scattered all over the code base.

			// Command line args
			fmt.Sprintf("UCP_URL=%s", os.Getenv("UCP_URL")),
			fmt.Sprintf("UCP_FINGERPRINT=%s", os.Getenv("UCP_FINGERPRINT")),
			fmt.Sprintf("UCP_HOST_ADDRESS=%s", os.Getenv("UCP_HOST_ADDRESS")),

			// Grey area... these should be doc'd but we don't want them passed as CLI args for security
			fmt.Sprintf("REGISTRY_USERNAME=%s", os.Getenv("REGISTRY_USERNAME")),
			fmt.Sprintf("REGISTRY_PASSWORD=%s", os.Getenv("REGISTRY_PASSWORD")),
			fmt.Sprintf("REGISTRY_EMAIL=%s", os.Getenv("REGISTRY_EMAIL")),
			fmt.Sprintf("UCP_ADMIN_USER=%s", os.Getenv("UCP_ADMIN_USER")),
			fmt.Sprintf("UCP_ADMIN_PASSWORD=%s", os.Getenv("UCP_ADMIN_PASSWORD")),
			fmt.Sprintf("UCP_PASSPHRASE=%s", os.Getenv("UCP_PASSPHRASE")),

			// Internal env vars (not command line args)
			fmt.Sprintf("UCP_BOOTSTRAP_PHASE2=%s", config.OrcaInstanceID),
			fmt.Sprintf("UCP_INSTANCE_KEY=%s", config.OrcaInstanceKey),
			fmt.Sprintf("UCP_LOCAL_NAME=%s", os.Getenv("UCP_LOCAL_NAME")),
			fmt.Sprintf("UCP_HOSTNAMES=%s", os.Getenv("UCP_HOSTNAMES")),
		}...),
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
		},
		// TODO Any other things we want to match up?
	}
	hostConfig := &container.HostConfig{
		Binds:      config.Phase2VolumeMounts,
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Privileged: privileged,
		PidMode:    pidMode,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}
	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.BootstrapPhase2ContainerName)
	if err != nil {
		return 1, err
	}
	containerId := resp.ID

	defer c.client.RemoveContainer(containerId, true, true)
	defer c.stopPhase2(containerId) // XXX is this redundant?

	// Also wire up a signal handler in case we're killed
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func() {
		s := <-ch
		log.Errorf("Signal %v received, aborting", s)
		c.client.RemoveContainer(containerId, true, true)
		os.Exit(1)
	}()

	// Wait for it to finish and pass through all the output...
	attachResp, err := c.client.ContainerAttach(containerId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return 1, err
	}

	log.Debugf("Launching phase 2 with: %v (%s)", args, containerId)
	if err := c.client.StartContainer(containerId); err != nil {
		return 1, err
	}

	doneChan := make(chan int)
	go func() {
		res, err := c.client.ContainerWait(containerId)
		if err != nil {
			log.Errorf("Failed to wait for phase 2 container: %s", err)
		}
		doneChan <- res
	}()

	rd := bufio.NewReader(attachResp.Reader)
	// Send all remaining input to the second phase
	go func() {
		if _, err := io.Copy(attachResp.Conn, os.Stdin); err != nil {
			log.Warnf("Stdin copy interrupted: %s", err)
		}
	}()

	if isTty {
		go func() {
			if _, err := io.Copy(os.Stdout, rd); err != nil {
				log.Warnf("TTY output copy interrupted: %s", err)
			}
		}()
	} else {
		// Keep draining the output while we wait
		go func() {
			// Temporarily change the log level since stdcopy can sometimes spew
			// TODO - fix stdcopy upstream so there's a public log ref that can be quieted down
			oldLevel := log.GetLevel()
			log.SetLevel(log.InfoLevel)
			defer log.SetLevel(oldLevel)
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, rd); err != nil {
				log.Errorf("Failed to stream logs: %s\n", err)
				return
			}
		}()
	}

	// Wait for the exit
	exitCode := <-doneChan

	return exitCode, nil
}
