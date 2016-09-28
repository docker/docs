package manager

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/orca"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func (m DefaultManager) SupportDump(w io.Writer) error {
	nodes, err := m.Nodes()
	if err != nil {
		return err
	}

	// Pull on swarm will retrieve the image on all nodes, which is exactly what we want
	//
	// Warning: With the current swarm, a down node can cause this to hang for a *very* long time
	//          however this bug is supposed to be fixed by the time 1.9 ships.
	//          In my testing, curl doesn't timeout, and correctly detects the down node
	//          it just takes a really long time.
	if _, _, err := m.client.ImageInspectWithRaw(context.TODO(), m.supportImage, false); client.IsErrImageNotFound(err) {
		log.Infof("Pulling image '%s'", m.supportImage)
		if out, err := m.client.ImagePull(context.TODO(), m.supportImage, types.ImagePullOptions{}); err != nil {
			return fmt.Errorf("Couldn't pull image '%s': %s", m.supportImage, err)
		} else {
			out.Close() // TODO - could do progress reporting....
		}
	} else if err != nil {
		return fmt.Errorf("Couldn't find image '%s': %s", m.supportImage, err)
	}

	z := zip.NewWriter(w)
	defer z.Close()

	// Note: If the size of the dumps gets large, we may want to
	//       refactor this and stage them in the local filesystem
	//       then read them back in one at a time
	payloadsMutex := sync.Mutex{}
	payloads := make(map[string][]byte)

	var wg sync.WaitGroup
	log.Debug("Starting support dumps in parallel")
	for _, node := range nodes {
		wg.Add(1)
		go func(node *orca.Node) {
			defer wg.Done()
			containerConfig := &container.Config{
				Image: m.supportImage,
				// Wire up constraints so it must run on the designated node
				Env: []string{
					fmt.Sprintf("constraint:node==%s", node.Name),
				},
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

			fail := func(msg string) {
				log.Errorf(msg)
				payloadsMutex.Lock()
				defer payloadsMutex.Unlock()
				payloads[node.Name] = []byte(msg)
				return
			}

			id := ""
			if resp, err := m.client.ContainerCreate(context.TODO(), containerConfig, hostConfig, nil, ""); err != nil {
				fail(fmt.Sprintf("Failed to create container: %s", err))
				return
			} else {
				id = resp.ID
			}

			resp, err := m.client.ContainerAttach(context.TODO(), id, types.ContainerAttachOptions{
				Stream: true,
				Stdout: true,
				Stderr: true,
			})
			if err != nil {
				fail(fmt.Sprintf("Failed to attach to container: %s", err))
				return
			}
			stream := resp.Reader

			if err := m.client.ContainerStart(context.TODO(), id, types.ContainerStartOptions{
				CheckpointID: "",
			}); err != nil {
				fail(fmt.Sprintf("Failed to start container: %s", err))
				return
			}
			defer m.client.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{Force: true, RemoveVolumes: false})

			// TODO - might want to contribute this up to the engine-api...
			ch := make(chan int)
			go func() {
				res, err := m.client.ContainerWait(context.TODO(), id)
				if err != nil {
					log.Errorf("Failed to wait for support container: %s", err)
				}
				ch <- res
			}()

			select {
			case <-ch:
				log.Debugf("Container '%s' exited correctly", id)
			case <-time.After(time.Duration(m.supportTimeout) * time.Second):
				fail(fmt.Sprintf("Container '%s' timed out while getting the support log", id))
				return
			}

			stdoutBuffer := new(bytes.Buffer)
			stderrBuffer := new(bytes.Buffer)
			if _, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, stream); err != nil {
				fail(fmt.Sprintf("Error copying support log from '%s': %s", node.Name, err))
				return
			}

			payloadsMutex.Lock()
			defer payloadsMutex.Unlock()
			payloads[node.Name+".tgz"] = stdoutBuffer.Bytes()

		}(node)
	}
	log.Debug("Gathering logs from infrastructure containers")
	containers, err := m.ListContainers(true, false, fmt.Sprintf(`{"label":{"com.docker.ucp.InstanceID=%s":true}}`, m.ID()))
	if err != nil {
		log.Warnf("Support dump failed to fetch containers: %s", err)
	}

	// XXX Should we throttle for a large cluster, or let the network I/O be our throttle?
	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()
			log.Debugf("Fetching log data for infrastructure container: %v", container.Names)

			// Clean up the name - these should be unique across the cluster
			name := ""
			if len(container.Names) > 0 {
				name = container.Names[0]
			} else {
				name = container.ID
			}
			if strings.HasPrefix(name, "/") {
				name = "." + name
			}
			fail := func(msg string) {
				log.Errorf(msg)
				payloadsMutex.Lock()
				defer payloadsMutex.Unlock()
				payloads[name] = []byte(msg)
				return
			}

			// Gather the logs
			reader, err := m.client.ContainerLogs(context.TODO(), container.ID, types.ContainerLogsOptions{
				Follow:     false,
				ShowStdout: true,
				ShowStderr: true,
			})
			if err != nil {
				fail(fmt.Sprintf("Failed to gather logs for infrastructure container: %s %s", name, err))
				return
			}

			out := new(bytes.Buffer)
			if _, err = stdcopy.StdCopy(out, out, reader); err != nil {
				fail(fmt.Sprintf("Error copying log from '%s': %s", name, err))
				return
			}
			payloadsMutex.Lock()
			defer payloadsMutex.Unlock()
			payloads[name+".log"] = out.Bytes()
		}(container)
	}
	log.Debugf("Waiting for support dumps to complete")
	wg.Wait()
	log.Debugf("Streaming dumps to client")
	for name, payload := range payloads {
		h := &zip.FileHeader{Name: name}
		h.SetModTime(time.Now())
		f, err := z.CreateHeader(h)
		if err != nil {
			msg := fmt.Sprintf("Failed to create %s in support dump: %s", name, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		_, err = f.Write(payload)
		if err != nil {
			msg := fmt.Sprintf("Failed to write payload %s in support dump: %s", name, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	log.Debugf("Streaming complete")

	return nil
}
