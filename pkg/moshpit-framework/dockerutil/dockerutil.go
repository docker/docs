package dockerutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
)

func MakeRegistryAuth(username, password, refreshToken string) string {
	authCfg := types.AuthConfig{}
	if refreshToken != "" {
		authCfg.IdentityToken = refreshToken
	} else {
		authCfg.Username = username
		authCfg.Password = password
	}
	buf, _ := json.Marshal(authCfg)
	return base64.URLEncoding.EncodeToString(buf)
}

func MakeDockerClient(host, jwt string, httpClient *http.Client, hubUsername, hubPassword, refreshToken string) (*client.Client, error) {
	if !strings.HasPrefix(host, "unix://") {
		hostPort := strings.Split(host, ":")
		if len(hostPort) == 1 {
			host = fmt.Sprintf("tcp://%s:443", host)
		} else {
			host = fmt.Sprintf("tcp://%s:%s", hostPort[0], hostPort[1])
		}
	}
	headers := map[string]string{}
	if jwt != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", jwt)
	}
	if hubUsername != "" && hubPassword != "" {
		headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword, refreshToken)
	}
	return client.NewClient(host, "v1.22", httpClient, headers)
}

func Poll(interval time.Duration, retries int, run func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	for i := 0; i < retries; i++ {
		err = run()
		if err != nil {
			time.Sleep(interval)
		} else {
			return nil
		}
	}
	return fmt.Errorf("Polling failed with %d attempts %s apart: %s", retries, interval, err)
}

func IsNoSuchImageErr(err string) bool {
	return strings.Contains(err, "No such container") || strings.Contains(err, "not found") || strings.Contains(err, "unable to detect a container with the given ID or name")
}

func SafeContainerRemove(ctx context.Context, dc *client.Client, id string) error {
	//log := moshpit.LoggerFromCtx(ctx)
	options := types.ContainerRemoveOptions{
		ContainerID: id,
		Force:       true,
	}
	// we keep trying to delete the container until we are told it doesn't exist
	// it's not safe to block until delete returns when using swarm
	err := Poll(time.Second, 10, func() error {
		//log.Debugf("Deleting container %s", options.ContainerID)
		err := dc.ContainerRemove(context.Background(), options)
		if err == nil {
			return fmt.Errorf("Timeout trying to delete container %s", options.ContainerID)
		}
		//log.Debugf("Failed to remove container: %s", err)
		// ignore aufs race errors
		if (strings.Contains(err.Error(), "500 Internal Server Error: Driver") && strings.Contains(err.Error(), "failed to remove root filesystem")) ||
			strings.Contains(err.Error(), "500 Internal Server Error: Unable to remove filesystem") {
			//log.Debugf("Ignored error, but retrying %s...", err)
			return err
		}
		// accept not found as success
		if IsNoSuchImageErr(err.Error()) {
			//log.Debugf("Ignored error %s...", err)
			return nil
		}
		// if it's any other error, don't bother retrying
		panic(fmt.Errorf("Unknown error: %s", err))
	})

	return err
}
