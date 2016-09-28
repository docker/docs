package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/integration/utils"
)

// TestScaleContainer attempts to scale a container at the API level
func (s *APITestSuite) TestScaleContainer() {
	require := require.New(s.T())
	assert := assert.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	// Obtain admin docker client
	dclient, err := utils.GetUserDockerClient(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	client, err := utils.ConvertToEngineAPI(dclient)
	require.Nil(err)

	// Create a container
	containerCfg := &container.Config{
		Image:        "busybox",
		Cmd:          []string{"tail", "--follow", "--retry"},
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,
		// NOTE: swarm affinity and constraint environment variables are NOT inheritted to scaled replicas
		// NOTE: a large set of dummy environment variables are needed to reduce flaky behavior in older engines
		Env: []string{"TESTENV", "DOCKER_FIX=true", "DUMMY_VAR=1234512345",
			"ANOTHER_VAR=1234567891234567789", "YET_ONE_MORE=123456789", "OK_LAST_ONE=123123123123123123123123"},
		Labels: map[string]string{"com.docker.test": "set"},
	}

	hostConfig := &container.HostConfig{
		Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		Privileged:  true,
		DNS:         []string{"8.8.8.8", "8.8.4.4"},
		NetworkMode: "none",
	}

	cresp, err := client.ContainerCreate(context.TODO(), containerCfg, hostConfig, nil, "scale-test-container")
	require.Nil(err)
	containerID := cresp.ID

	defer client.ContainerRemove(context.TODO(), containerID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})

	log.Debug(cresp.Warnings)

	// Obtain admin TLS certificates
	tlsConfig, err := utils.GetUserTLSConfig(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	httpClient := &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}

	// Scale four replicas
	path := fmt.Sprintf("%s/api/containers/%s/scale?n=4", serverURL, containerID)
	req, err := http.NewRequest("POST", path, bytes.NewBufferString(""))
	require.Nil(err)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Debug(err)
	}
	log.Debug(resp)
	require.Nil(err)
	require.Equal(resp.StatusCode, http.StatusOK)

	// Decode the response
	var scaleResult manager.ScaleResult
	err = json.NewDecoder(resp.Body).Decode(&scaleResult)
	require.Nil(err)

	require.Equal(0, len(scaleResult.Errors))

	// Make sure every replica has the same container and host configs
	for _, replicaID := range scaleResult.Scaled {
		ctr, err := client.ContainerInspect(context.TODO(), replicaID)
		require.Nil(err)

		// We cannot simply use reflect.DeepEqual as the structs slightly differ
		assert.Equal(ctr.Config.Image, containerCfg.Image)
		assert.Equal(ctr.Config.Cmd, containerCfg.Cmd)
		assert.Equal(ctr.Config.OpenStdin, containerCfg.OpenStdin)
		assert.Equal(ctr.Config.AttachStdin, containerCfg.AttachStdin)
		assert.Equal(ctr.Config.AttachStdout, containerCfg.AttachStdout)
		assert.Equal(ctr.Config.AttachStderr, containerCfg.AttachStderr)
		assert.Equal(ctr.Config.StdinOnce, containerCfg.StdinOnce)
		assert.Equal(ctr.Config.Env, containerCfg.Env)
		for label, value := range containerCfg.Labels {
			assert.Equal(value, ctr.Config.Labels[label])
		}

		assert.Equal(ctr.HostConfig.Binds, hostConfig.Binds)
		assert.Equal(ctr.HostConfig.Privileged, hostConfig.Privileged)
		assert.Equal(ctr.HostConfig.DNS, hostConfig.DNS)
		assert.Equal(ctr.HostConfig.NetworkMode, hostConfig.NetworkMode)
	}
}
