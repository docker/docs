package agent

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func launchReconcile(dclient *client.Client, imageName string, payload string) (string, error) {
	binds := []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		"/var/lib/docker/swarm/certificates:/var/lib/docker/swarm/certificates",
		config.SwarmNodeCertVolumeName + ":" + config.SwarmNodeCertVolumeMount,
		config.OrcaServerCertVolumeName + ":" + config.OrcaServerCertVolumeMount,
		config.SwarmControllerCertVolumeName + ":" + config.SwarmControllerCertVolumeMount,
		config.SwarmRootCAVolumeName + ":" + config.SwarmRootCAVolumeMount,
		config.SwarmKvCertVolumeName + ":" + config.SwarmKvCertVolumeMount,
		config.OrcaRootCAVolumeName + ":" + config.OrcaRootCAVolumeMount,
		config.OrcaKVVolumeName + ":" + config.OrcaKVVolumeMount,
		config.AuthAPICertsVolumeName + ":" + config.AuthAPICertsVolumeMount,
		config.AuthWorkerCertsVolumeName + ":" + config.AuthWorkerCertsVolumeMount,
		config.AuthStoreCertsVolumeName + ":" + config.AuthStoreCertsVolumeMount,
		config.AuthStoreDataVolumeName + ":" + config.AuthStoreDataVolumeMount,
		config.AuthWorkerDataVolumeName + ":" + config.AuthWorkerDataVolumeMount,
	}

	cfg := &container.Config{
		Image:        imageName,
		Cmd:          []string{"reconcile", "--payload", payload},
		Tty:          false,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,
		Labels: map[string]string{
			fmt.Sprintf("%s.InstanceID", config.OrcaLabelPrefix): config.OrcaInstanceID,
		},
		// TODO: set up UCP version label
	}

	hostConfig := &container.HostConfig{
		Binds:      binds,
		DNS:        config.DNS,
		DNSOptions: config.DNSOpt,
		DNSSearch:  config.DNSSearch,
		Resources: container.Resources{
			MemorySwap: -1,
		},
	}

	// NOTE: ucp-reconcile container will not be removed automatically
	resp, err := dclient.ContainerCreate(context.TODO(), cfg, hostConfig, nil, orcaconfig.OrcaReconcileContainerName)
	if err != nil {
		return "", err
	}

	// Stream all output from the ucp-reconcile container to the agent
	attachResp, err := dclient.ContainerAttach(context.TODO(), resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return "", err
	}

	// Start the ucp-reconcile container
	err = dclient.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return resp.ID, err
	}

	// Create a reader of the ucp-reconcile container's output
	rd := bufio.NewReader(attachResp.Reader)

	// Launch a goroutine that redirects the ucp-reconcile output to the agent's stdout and stderr
	go func() {
		// Temporarily change the log level since stdcopy can sometimes spew
		oldLevel := log.GetLevel()
		log.SetLevel(log.InfoLevel)
		defer log.SetLevel(oldLevel)

		// TODO: is this blocking forever on Reconcile exit?
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, rd); err != nil {
			log.Errorf("Failed to stream logs: %s\n", err)
			return
		}
	}()

	return resp.ID, nil
}
