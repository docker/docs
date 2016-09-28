package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

var errNoAuthStoreAddrsGiven = errors.New("no auth store addresses given")

// RunAuthSyncDB runs the Auth SyncDB container.
func (c *EngineClient) RunAuthSyncDB(repair bool, authStoreAddrs ...string) error {
	log.Debug("Running Auth SyncDB")

	// There must be at least one db driver address given.
	if len(authStoreAddrs) == 0 {
		return errNoAuthStoreAddrsGiven
	}

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.AuthSyncDBContainerName)
	if err != nil {
		return err
	}

	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.AuthAPICertsVolumeName, "/tls"),
	}

	dbAddrArgs := make([]string, len(authStoreAddrs))
	for i, authStoreAddr := range authStoreAddrs {
		dbAddrArgs[i] = fmt.Sprintf("--db-addr=%s", authStoreAddr)
	}

	cmd := append(dbAddrArgs,
		fmt.Sprintf("--debug=%t", log.GetLevel() == log.DebugLevel),
		"--jsonlog",
		"sync-db",
	)

	if repair {
		cmd = append(cmd, "--emergency-repair")
	}

	cfg := &container.Config{
		Image: imageName,
		Tty:   c.IsTty(),
		Cmd:   cmd,
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
		},
	}

	hostConfig := &container.HostConfig{
		Binds:      mounts,
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	createResp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.AuthSyncDBContainerName)
	if err != nil {
		log.Debugf("Failed to create auth sync-db container: %s", err)
		return err
	}

	containerID := createResp.ID

	// Attach to the container and pass through all the container
	// output.
	if err := c.streamContainerOutput(orcaconfig.AuthSyncDBContainerName, c.IsTty()); err != nil {
		log.Errorf("Failed to stream container logs: %s", err)
		return err
	}

	if err := c.client.StartContainer(containerID); err != nil {
		log.Debugf("Failed to start auth sync-db container: %s", err)
		return err
	}

	// Wait for the sync-db task to complete (it may take up to 1 minute
	// to perform initial schema setup but we should time out after 2
	// minutes.)
	timeout := 2 * time.Minute
	done := make(chan error, 1)

	timer := time.AfterFunc(timeout, func() {
		done <- fmt.Errorf("auth sync-db timed out after %s", timeout)
	})

	go func() {
		exitStatus, err := c.client.ContainerWait(containerID)
		timer.Stop()

		if err != nil {
			err = fmt.Errorf("unable to wait for auth sync-db container: %s", err)
		} else if exitStatus != 0 {
			err = fmt.Errorf("auth sync-db returned non-zero exit status: %d", exitStatus)
		}

		done <- err
	}()

	if err := <-done; err != nil {
		return err
	}

	c.client.RemoveContainer(containerID, true, false)

	return nil
}

// streamContainerOutput attaches to the container and returns immediately if
// there is an error. Upon successful attach, a goroutine is spawned to copy
// output from the container to stdout and stderr of this process before
// returning immediately.
func (c *EngineClient) streamContainerOutput(containerNameOrID string, isTty bool) error {
	attachResp, err := c.client.ContainerAttach(containerNameOrID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}

	// Spawn a goroutine to copy stdout and stderr from the container.
	go func() {
		defer attachResp.Close()

		if isTty {
			// If the container uses a TTY, we can just copy
			// everything to stdout.
			if _, err := io.Copy(os.Stdout, attachResp.Reader); err != nil {
				log.Errorf("unable to copy output from container %s: %s", containerNameOrID, err)
			}

			return
		}

		// We need to temporarily bump the global log level because the
		// stdcopy package uses the global logger to log extraneous
		// debug info and we don't want to get that mixed up with our
		// own logs here.
		// TODO - fix stdcopy upstream so there's a public log ref that
		// can be quieted down
		defer log.SetLevel(log.GetLevel())
		log.SetLevel(log.InfoLevel)

		// Demultiplex to stdout and stderr.
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader); err != nil {
			log.Errorf("unable to copy output from container %s: %s", containerNameOrID, err)
		}
	}()

	return nil
}

// RunAuthCreateAdmin runs the auth setup container to create the initial
// 'admin' user with the given password. We do not bother streaming any logs
// from this container because there's not really any valuable output from this
// task. It has an interactive mode for getting the username and password which
// we don't use here. Other than that, the only thing it prints out is a simple
// success message (without log formatting) or an error message.
func (c *EngineClient) RunAuthCreateAdmin(adminUsername, adminPassword string) error {
	log.Debug("Running Auth Create Admin")

	imageName, err := orcaconfig.GetContainerImage(orcaconfig.AuthCreateAdminContainerName)
	if err != nil {
		return err
	}

	mounts := []string{
		fmt.Sprintf("%s:%s:ro", config.AuthAPICertsVolumeName, "/tls"),
	}

	authStoreAddr := fmt.Sprintf("%s:%d", config.OrcaHostAddress, orcaconfig.AuthStorePort)

	cfg := &container.Config{
		Image: imageName,
		Env: []string{
			fmt.Sprintf("USERNAME=%s", adminUsername),
			fmt.Sprintf("PASSWORD=%s", adminPassword),
		},
		Cmd: strslice.StrSlice{
			// Connect to the local auth store.
			fmt.Sprintf("--db-addr=%s", authStoreAddr),
			fmt.Sprintf("--debug=%t", log.GetLevel() == log.DebugLevel),
			"--jsonlog",
			"create-admin",
		},
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
		},
	}

	hostConfig := &container.HostConfig{
		Binds:      mounts,
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	resp, err := c.client.CreateContainer(cfg, hostConfig, orcaconfig.AuthCreateAdminContainerName)
	if err != nil {
		return err
	}

	containerID := resp.ID

	if err := c.client.StartContainer(containerID); err != nil {
		log.Debugf("Failed to launch Auth API: %s", err)
		return err
	}

	// Wait for the create-admin task to complete (it should not take more
	// than 1 minute to perform).
	timeout := time.Minute
	done := make(chan error, 1)

	timer := time.AfterFunc(timeout, func() {
		done <- fmt.Errorf("auth create-admin timed out after %s", timeout)
	})

	go func() {
		exitStatus, err := c.client.ContainerWait(containerID)
		timer.Stop()

		if err != nil {
			err = fmt.Errorf("unable to wait for auth create-admin container: %s", err)
		} else if exitStatus != 0 {
			err = fmt.Errorf("auth create-admin returned non-zero exit status: %d", exitStatus)
		}

		done <- err
	}()

	if err := <-done; err != nil {
		return err
	}

	c.client.RemoveContainer(containerID, true, false)

	return nil
}
