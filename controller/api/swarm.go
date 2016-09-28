package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	goversion "github.com/hashicorp/go-version"

	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/version"
)

// Supported Engine API Versions
var maxVersion, _ = goversion.NewVersion("1.24")
var minVersion, _ = goversion.NewVersion("1.14")

func isSupported(v string) bool {
	requestedVersion, err := goversion.NewVersion(v)
	if err != nil {
		log.Errorf("Unable to parse version string: %s", v)
		return false
	}

	if requestedVersion.GreaterThan(maxVersion) ||
		requestedVersion.LessThan(minVersion) {
		return false
	}

	return true
}

// engineRedirect forwards a request to the engine proxy
// As of UCP 1.2, Engine Swarm APIs use engineRedirect instead of swarmRedirect
func (a *Api) engineRedirect(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var err error
	rc.Request.URL, err = url.ParseRequestURI(a.engineProxyURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.fwd.ServeHTTP(w, rc.Request)
}

// swarmRedirect forwards a request to the Swarm v1 cluster
func (a *Api) swarmRedirect(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	var err error
	rc.Request.URL, err = url.ParseRequestURI(a.swarmClassicURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.fwd.ServeHTTP(w, rc.Request)
}

type proxyWriter struct {
	Body       *bytes.Buffer
	Headers    *map[string][]string
	StatusCode *int
}

func (p proxyWriter) Header() http.Header {
	return *p.Headers
}
func (p proxyWriter) Write(data []byte) (int, error) {
	return p.Body.Write(data)
}
func (p proxyWriter) WriteHeader(code int) {
	*p.StatusCode = code
}

// The following routines are split out to allow profiling, and potential future optimizations

// Introspect the /info call so we can augment it with Orca details
func (a *Api) swarmInfo(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// TODO - if we find we need to do version handling elsewhere, lets refactor
	//        this - use semantic versioning, and make a helper routine
	apiVersion := rc.PathVars["version"]
	old := false // Default to new API if we can't figure it out for some reason
	if apiVersion != "" {
		v := strings.Split(apiVersion, ".")
		if len(v) >= 2 && v[0] == "1" {
			if v2, err := strconv.Atoi(v[1]); err == nil && v2 < 22 {
				old = true
			}
		}
	}

	if old {
		a.swarmAdapter(w, rc, a.fixupInfoOld)
	} else {
		a.swarmAdapter(w, rc, a.fixupInfoNew)
	}
}

func (a *Api) fixupInfoOld(body *bytes.Buffer) ([]byte, error) {
	return a.fixupInfo(body, true)
}
func (a *Api) fixupInfoNew(body *bytes.Buffer) ([]byte, error) {
	return a.fixupInfo(body, false)
}

func (a *Api) fixupInfo(body *bytes.Buffer, oldBehavior bool) ([]byte, error) {
	// From time to time the /info API endpoint changes, and to prevent
	// us from dropping content on the floor, lets parse this generically
	// (once we migrate to engine-api, we may be able to stay in sync
	//  with the engine more closely, but for now, dockerclient is lagging)
	var info map[string]interface{}
	d := json.NewDecoder(body)
	d.UseNumber() // Make sure numbers don't become floats (breaks MemTotal and others without this)
	if err := d.Decode(&info); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal the swarm info: %s", err)
	}
	// Replace a few fields to be Orca specific
	info["ID"] = a.manager.ID()
	if hostname, err := os.Hostname(); err == nil {
		info["Name"] = hostname
	} else {
		// Leave swarms Name if we can't get our hostname for some reason
		log.Errorf("Failed to lookup Orca hostname: %s", err)
	}
	info["ServerVersion"] = fmt.Sprintf("ucp/%s", version.TagVersion())

	managers := a.manager.GetManagers()
	pad1 := ""
	pad2 := " "
	keyName := "SystemStatus"
	if oldBehavior {
		pad1 = "\u0008"
		pad2 = ""
		keyName = "DriverStatus"
	}
	systemStatusRaw := []interface{}{}
	switch v := info[keyName].(type) {
	case []interface{}:
		systemStatusRaw = v
	case nil:
		// no driver status provided
	default:
		return nil, fmt.Errorf("Malformed info DriverStatus response from backend: %v - %T", v, v)
	}
	systemStatusRaw = append(systemStatusRaw, []string{pad1 + "Cluster Managers", fmt.Sprintf("%d", len(managers))})
	for _, mn := range managers {
		status := a.manager.GetStatus(mn)
		systemStatusRaw = append(systemStatusRaw, []string{pad2 + fmt.Sprintf("%s", mn.Name), status})
		systemStatusRaw = append(systemStatusRaw, []string{pad2 + " └ Orca Controller", mn.ControllerURL})
		systemStatusRaw = append(systemStatusRaw, []string{pad2 + " └ Classic Swarm Manager", mn.SwarmClassicManagerURL})
		systemStatusRaw = append(systemStatusRaw, []string{pad2 + " └ Engine Swarm Manager", mn.EngineProxyURL})
		systemStatusRaw = append(systemStatusRaw, []string{pad2 + " └ KV", mn.KVURL})
	}
	info[keyName] = systemStatusRaw

	licenseCfg := a.manager.GetLicense()
	labelsRaw := []interface{}{}
	switch v := info["Labels"].(type) {
	case []interface{}:
		labelsRaw = v
	case nil:
		// no driver status provided
	default:
		return nil, fmt.Errorf("Malformed info Label response from backend: %v - %T", v, v)
	}
	labelsRaw = append(labelsRaw, fmt.Sprintf("com.docker.ucp.license_key=%s", licenseCfg.License.KeyID))
	labelsRaw = append(labelsRaw, fmt.Sprintf("com.docker.ucp.license_max_engines=%d", licenseCfg.Details.MaxEngines))

	now := time.Now().UTC()
	if now.After(licenseCfg.Details.Expiration) {
		labelsRaw = append(labelsRaw, fmt.Sprintf("com.docker.ucp.license_expires=EXPIRED"))
	} else {
		labelsRaw = append(labelsRaw, fmt.Sprintf("com.docker.ucp.license_expires=%s", licenseCfg.Details.Expiration.String()))
	}
	info["Labels"] = labelsRaw

	// TODO - Add more orca specific stuff here...

	return json.Marshal(&info)
}

// Introspect the /version call so we can augment it with Orca version information
func (a *Api) swarmVersion(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmAdapter(w, rc, a.fixupVersion)
}
func (a *Api) fixupVersion(body *bytes.Buffer) ([]byte, error) {
	var ver types.Version
	if err := json.Unmarshal(body.Bytes(), &ver); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal the swarm version: %s", err)
	}
	ver.Version = fmt.Sprintf("ucp/%s", version.TagVersion())
	ver.GitCommit = version.GitCommit
	ver.GoVersion = runtime.Version()
	ver.BuildTime = version.BuildTime

	return json.Marshal(&ver)
}

// Introspect a swarm call so we can augment it with Orca details
func (a *Api) swarmAdapter(w http.ResponseWriter, rc *ctx.OrcaRequestContext, fixup func(body *bytes.Buffer) ([]byte, error)) {
	body := bytes.Buffer{}
	Headers := make(map[string][]string)
	var StatusCode int
	p := proxyWriter{
		Body:       &body,
		Headers:    &Headers,
		StatusCode: &StatusCode,
	}
	a.engineRedirect(p, rc)
	if StatusCode == 200 {
		data, err := fixup(&body)
		if err != nil {
			log.Errorf("Failed to fixup swarm API: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the headers over
		dest := w.Header()
		for k, vv := range Headers {
			for _, v := range vv {
				dest.Add(k, v)
			}
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(StatusCode)
		if _, err := w.Write(data); err != nil {
			log.Errorf("Failed to write response: %s", err)
		}
	} else {
		// Something went wrong, pass along the error without modification
		dest := w.Header()
		for k, vv := range Headers {
			for _, v := range vv {
				dest.Add(k, v)
			}
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", body.Len()))
		w.WriteHeader(StatusCode)
		if _, err := w.Write(body.Bytes()); err != nil {
			log.Errorf("Failed to write response: %s", err)
		}
	}
}

func (a *Api) swarmImages(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRedirect(w, rc)
}

func (a *Api) swarmImagesList(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRedirect(w, rc)
}

func (a *Api) swarmImagesInfo(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRedirect(w, rc)
}

func (a *Api) swarmContainersInfo(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRedirect(w, rc)
}

func (a *Api) swarmContainers(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	a.swarmRedirect(w, rc)
}

func (a *Api) swarmHijack(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	t := a.manager.DockerClientTransport()
	addr := a.swarmClassicURL

	// TODO: use a url.Url object
	if parts := strings.SplitN(addr, "://", 2); len(parts) == 2 {
		addr = parts[1]
	}

	var (
		d   net.Conn
		err error
	)

	if t.TLSClientConfig != nil {
		d, err = tls.Dial("tcp", addr, t.TLSClientConfig)
	} else {
		d, err = net.Dial("tcp", addr)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer nc.Close()
	defer d.Close()

	err = rc.Request.Write(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cp := func(dst io.Writer, src io.Reader, chDone chan struct{}) {
		io.Copy(dst, src)
		if conn, ok := dst.(interface {
			CloseWrite() error
		}); ok {
			conn.CloseWrite()
		}
		close(chDone)
	}

	inDone := make(chan struct{})
	outDone := make(chan struct{})
	go cp(d, nc, inDone)
	go cp(nc, d, outDone)

	// 1. When stdin is done, wait for stdout always
	// 2. When stdout is done, close the stream and wait for stdin to finish
	//
	// On 2, stdin copy should return immediately now since the out stream is closed.
	// Note that we probably don't actually even need to wait here.
	//
	// If we don't close the stream when stdout is done, in some cases stdin will hange
	select {
	case <-inDone:
		// wait for out to be done
		<-outDone
	case <-outDone:
		// close the conn and wait for stdin
		nc.Close()
		<-inDone
	}
}

func (a *Api) swarmModeInit(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	http.Error(w, "This UCP Controller is already configured as a swarm-mode cluster", http.StatusBadRequest)
}

// Might want to tune the error message...
func (a *Api) swarmModeJoin(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	http.Error(w, "This UCP Controller is already configured as a swarm-mode cluster", http.StatusBadRequest)
}

func (a *Api) swarmModeLeave(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	http.Error(w, "To perform a leave on a UCP node, please run the leave command directly against the node", http.StatusBadRequest)
}
