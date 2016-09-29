package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/runconfig"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types/container"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

type proxyStart struct {
	passthru      *passthru
	approver      Approver
	backendDialer BackendDialer
}

func (c *proxyStart) HandleHTTP(writer http.ResponseWriter, r *http.Request) error {
	log.Printf("proxy >> %s %s [start]\n", r.Method, r.URL)

	vars := mux.Vars(r)
	containerID := vars["name"]

	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	var hc *container.HostConfig
	var forwardConfig *runconfig.ContainerConfigWrapper
	if err := httputils.CheckForJSON(r); err == nil {
		// Either an old-style message with a JSON config,
		// or docker-compose passing {} instead of the empty body.
		var w runconfig.ContainerConfigWrapper
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&w); err != nil {
			if err != io.EOF {
				return err
			}
		} else {
			hc := w.HostConfig
			if hc != nil {
				// Old-style message.
				forwardConfig = &w
			}
		}
	}

	if hc == nil && c.approver != nil {
		cl, err := c.backendDialer.DockerClient()
		if err != nil {
			log.Printf("Failed to create client: %s", err)
			http.Error(writer, "ContainerInspect failed", 502)
			return nil // ??
		}

		// New-style message. Config was sent previously in start call.
		// Query Docker to get it. Note if the container doesn't exist then we
		// want to fall back to the passthru code below so the true error code
		// from docker is returned, rather than a 502 here.
		info, err := cl.ContainerInspect(context.Background(), containerID)
		if err != nil && !client.IsErrContainerNotFound(err) {
			log.Printf("Failed to inspect container %s: %s", containerID, err)
			http.Error(writer, "ContainerInspect failed", 502)
			return nil // ??
		}

		// Only approve mounts if the container exists
		if err == nil {
			err = approveMounts(c.approver, containerID, info.HostConfig)
			if err != nil {
				log.Printf("Mounts denied: %s", err)
				http.Error(writer, "Mounts denied: "+err.Error(), 502)
				return nil // ??
			}
		}
	}

	if forwardConfig == nil {
		// FIXME: if this call actually fails, what do we do about the
		// mounts and ports? I assume we don't get a container die?
		r.ContentLength = 0
		return c.passthru.HandleHTTP(writer, r)
	}

	body, err := encodeData(forwardConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(r.Method, r.RequestURI, body)
	if err != nil {
		return err
	}

	req.Header = r.Header
	req.Host = r.Host
	req.Trailer = r.Trailer
	req.RemoteAddr = r.RemoteAddr
	req.TLS = r.TLS
	req.Cancel = r.Cancel

	return c.passthru.HandleHTTP(writer, req)
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(params).Encode(data); err != nil {
			return nil, err
		}
	}
	return params, nil
}
