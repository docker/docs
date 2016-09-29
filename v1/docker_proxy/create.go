package proxy

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/runconfig"
)

type proxyCreate struct {
	passthru      *passthru
	mountRewriter MountRewriter
	envRewriter   EnvRewriter
}

func (c *proxyCreate) HandleHTTP(writer http.ResponseWriter, r *http.Request) error {
	log.Printf("proxy >> %s %s [rewriteBinds]\n", r.Method, r.URL)

	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	if err := httputils.CheckForJSON(r); err != nil {
		// error means there was no json
		return c.passthru.HandleHTTP(writer, r)
	}

	var w runconfig.ContainerConfigWrapper
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&w); err != nil {
		return err
	}

	if c.mountRewriter != nil {
		err := c.mountRewriter.RewriteMounts(w.HostConfig)
		if err != nil {
			return err
		}
	}

	if c.envRewriter != nil {
		c.envRewriter.RewriteEnv(w.Config)
	}

	body, err := encodeData(w)
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
