package client

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// Verify that the ports we need are available
func (c *EngineClient) CheckPorts(ports []*int) error {

	log.Info("Checking that required ports are available and accessible")
	conflicts := make(chan int, len(ports))
	blocked := make(chan int, len(ports))

	var wg sync.WaitGroup

	// Note: this test relies on the fact that if we run the proxy without a docker.sock
	//       it fails silently and always returns a 200 with an empty body (great for testing accessibility)
	imageName, err := orcaconfig.GetContainerImage(orcaconfig.OrcaAgentContainerName)
	if err != nil {
		return err
	}

	for i, p := range ports {
		wg.Add(1)
		// Spawn the port test containers in parallel to make this faster
		go func(i int, p *int) {
			defer wg.Done()

			log.Debugf("Checking for available and accessible port %d", *p)
			portMap := map[nat.Port]struct{}{
				nat.Port("2376/tcp"): struct{}{},
			}
			bindingMap := nat.PortMap{
				nat.Port("2376/tcp"): []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", *p),
					},
				},
			}

			// Start up a simple server container with the specified ports punched through so we can fail fast
			cfg := &container.Config{
				Image:        imageName,
				ExposedPorts: portMap,
				Cmd:          strslice.StrSlice{"test-server"},
			}
			hostConfig := &container.HostConfig{
				PortBindings: bindingMap,
				Resources: container.Resources{
					MemorySwap: -1,
				},
			}
			resp, err := c.client.CreateContainer(cfg, hostConfig, "")
			if err != nil {
				log.Errorf("Failed to create test for available ports %s", err)
				conflicts <- *p
				return
			}
			containerId := resp.ID
			defer c.client.RemoveContainer(containerId, true, true)

			// Start the container
			if err := c.client.StartContainer(containerId); err != nil {
				log.Debugf("Failed to launch port test container: %s", err)
				conflicts <- *p
				return
			}
			defer c.client.StopContainer(containerId, 5)
			err = waitForServer(
				fmt.Sprintf("http://%s:%d/", config.OrcaHostAddress, *p),
				200,
				&http.Client{Timeout: 500 * time.Millisecond},
				config.DockerTimeout,
			)
			if err != nil {
				// No need to log a special error message here as we'll consolidate the results
				blocked <- *p
			}
		}(i, p)
	}
	wg.Wait()
	close(conflicts) // No more input expected, so make sure the loop below will terminate
	close(blocked)   // No more input expected, so make sure the loop below will terminate
	if len(conflicts) == 0 && len(blocked) == 0 {
		log.Debug("All ports are open and available")
		return nil
	}

	// Favor the conflict messages, and display blocked message only if all are available
	if len(conflicts) == 0 {
		portList := make([]string, 0, len(blocked))
		for p := range blocked {
			portList = append(portList, fmt.Sprintf("%d", p))
		}
		return fmt.Errorf("The following required ports are blocked on your host: %s.  Check your firewall settings.", strings.Join(portList, ", "))
	}

	// Special case the Swarm port so we can give the swarm port flag message
	msg := ""
	portList := make([]string, 0, len(conflicts))
	for p := range conflicts {
		portList = append(portList, fmt.Sprintf("%d", p))
		if p == orcaconfig.SwarmPort {
			msg = fmt.Sprintf("  You may specify an alternative port number to %d with the --swarm-port argument.", orcaconfig.SwarmPort)
		}
	}

	return fmt.Errorf("The following required ports are already in use on your host - %s.%s", strings.Join(portList, ", "), msg)
}
