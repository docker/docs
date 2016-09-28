package manager

import (
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
)

func (m DefaultManager) ContainerLogs(id string, opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	responseBody, err := m.client.ContainerLogs(context.TODO(), id, opts)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (m DefaultManager) ScaleContainer(id string, numInstances int) ScaleResult {
	var (
		errChan = make(chan (error))
		resChan = make(chan (string))
		result  = ScaleResult{Scaled: make([]string, 0), Errors: make([]string, 0)}
	)

	containerInfo, err := m.Container(id)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result
	}

	for i := 0; i < numInstances; i++ {
		go func(instance int) {
			log.Debugf("scaling: id=%s #=%d", containerInfo.ID, instance)
			config := containerInfo.Config
			// clear hostname to get a newly generated
			config.Hostname = ""
			hostConfig := containerInfo.HostConfig
			resp, err := m.client.ContainerCreate(context.TODO(), config, hostConfig, nil, "")
			if err != nil {
				errChan <- err
				return
			}
			if err := m.client.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{
				CheckpointID: "",
			}); err != nil {
				errChan <- err
				return
			}
			resChan <- resp.ID
		}(i)
	}

	for i := 0; i < numInstances; i++ {
		select {
		case id := <-resChan:
			result.Scaled = append(result.Scaled, id)
		case err := <-errChan:
			log.Errorf("error scaling container: err=%s", strings.TrimSpace(err.Error()))
			result.Errors = append(result.Errors, strings.TrimSpace(err.Error()))
		}
	}

	return result
}

func (m DefaultManager) ListContainers(all bool, size bool, filterArgs string) ([]types.Container, error) {
	// todo: filter system containers
	args, err := filters.FromParam(filterArgs)
	if err != nil {
		return nil, err
	}
	return m.client.ContainerList(context.TODO(), types.ContainerListOptions{All: all, Size: size, Filter: args})
}

func (m DefaultManager) ListUserContainers(ctx *auth.Context, all bool, size bool, filters string) ([]types.Container, error) {
	allContainers, err := m.ListContainers(all, size, filters)
	if err != nil {
		return nil, err
	}

	acct := ctx.User

	// return immediately for admins
	if acct.Admin {
		return allContainers, nil
	}

	access, err := m.GetAccess(ctx)
	if err != nil {
		return nil, err
	}

	userContainers := []types.Container{}

	for _, c := range allContainers {
		// check access labels
		matched := false
		for k, v := range c.Labels {
			if k != orca.UCPAccessLabel {
				continue
			}

			// check for label access
			if k == orca.UCPAccessLabel {
				if lvl, ok := access[v]; ok && lvl > auth.None {
					log.Debugf("user container (access): id=%s level=%d", c.ID, v)
					userContainers = append(userContainers, c)
					matched = true
					break
				}
			}
		}

		if matched {
			continue
		}

		for k, v := range c.Labels {
			if k != orca.UCPOwnerLabel {
				continue
			}

			// check for ownership
			if k == orca.UCPOwnerLabel && v == acct.Username {
				log.Debugf("user container (owner): id=%s user=%s", c.ID, v)
				userContainers = append(userContainers, c)
				break
			}
		}
	}

	log.Debugf("user containers: %v", userContainers)

	return userContainers, nil
}
