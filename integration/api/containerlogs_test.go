package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"

	"github.com/docker/orca/integration/utils"
)

type ContainerLogsResponse struct {
	ContainerId string `json:"container_id,omitempty"`
	Token       string `json:"token,omitempty"`
}

// Test that we can get logs from our websocket
func (s *APITestSuite) TestStdoutLogsContainer() {
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
		Cmd:          []string{"echo", "hello"},
		AttachStdout: true,
		AttachStderr: true,
	}
	hostConfig := &container.HostConfig{}

	container, err := client.ContainerCreate(context.TODO(), containerCfg, hostConfig, nil, "logs-test-container")
	require.Nil(err)
	containerID := container.ID

	client.ContainerStart(context.TODO(), containerID, types.ContainerStartOptions{
		CheckpointID: "",
	})

	// Wait for container to start and run command before attempting to get logs
	time.Sleep(1 * time.Second)

	defer client.ContainerRemove(context.TODO(), containerID, types.ContainerRemoveOptions{Force: true})

	log.Debug(container.Warnings)

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

	// Get a container logs authentication token
	tokenUrl := fmt.Sprintf("%s/api/containerlogs/%s", serverURL, containerID)
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBufferString(""))
	require.Nil(err)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Debug(err)
	}
	require.Nil(err)
	defer resp.Body.Close()

	require.Equal(http.StatusOK, resp.StatusCode)

	tokenJson := ContainerLogsResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenJson)

	// Create websocket config
	serverHost := strings.Replace(serverURL, "https://", "", -1)
	wsUrl := fmt.Sprintf("wss://%s/containerlogs?token=%s&id=%s", serverHost, tokenJson.Token, containerID)
	wsConfig, err := websocket.NewConfig(wsUrl, serverURL)
	require.Nil(err)

	// Add TLS certs to websocket config
	wsConfig.TlsConfig = tlsConfig

	// Connect to logs websocket
	ws, err := websocket.DialConfig(wsConfig)
	require.Nil(err)

	// Grab the first message
	var msg = make([]byte, 512)
	len, _ := ws.Read(msg)

	assert.Equal("hello\n", string(msg[:len]))
}
