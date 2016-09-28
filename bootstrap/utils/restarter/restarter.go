package restarter

import (
	"fmt"

	"github.com/docker/orca/bootstrap/client"
)

type Restarter struct {
	containerIDs map[string]bool
	engineclient *client.EngineClient
}

func (r *Restarter) RestartAll() error {
	var targets []string
	for id, restarted := range r.containerIDs {
		if restarted {
			continue
		}
		targets = append(targets, id)
	}
	if len(targets) == 0 {
		return fmt.Errorf("No containers to be restarted")
	}
	return r.engineclient.RestartContainers(targets)
}

func (r *Restarter) RestartContainer(id string) error {
	restarted, ok := r.containerIDs[id]
	if !ok {
		return fmt.Errorf("this restarter is not configured to restart container with id or name %s", id)
	}
	if restarted {
		return fmt.Errorf("this container has been already restarted")
	}
	r.containerIDs[id] = true

	return r.engineclient.RestartContainers([]string{id})
}

func (r *Restarter) SetContainerTargets(ids []string) {
	for _, id := range ids {
		r.containerIDs[id] = false
	}
}

func NewRestarter(ec *client.EngineClient) *Restarter {
	return &Restarter{
		containerIDs: make(map[string]bool),
		engineclient: ec,
	}
}
