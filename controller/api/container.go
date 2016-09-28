package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	dockerfilters "github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/config"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/resources"
)

func (a *Api) execContainer(ws *websocket.Conn) {
	qry := ws.Request().URL.Query()
	containerId := qry.Get("id")
	command := qry.Get("cmd")
	ttyWidth := qry.Get("w")
	ttyHeight := qry.Get("h")
	token := qry.Get("token")
	cmd := strings.Split(command, ",")

	if !a.manager.ValidateConsoleSessionToken(containerId, token) {
		ws.Write([]byte("unauthorized"))
		ws.Close()
		return
	}

	log.Debugf("starting exec session: container=%s cmd=%s", containerId, command)
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          cmd,
		Detach:       true,
	}

	resp, err := a.manager.DockerClient().ContainerExecCreate(context.TODO(), containerId, execConfig)
	if err != nil {
		log.Errorf("error calling exec: %s", err)
		return
	}
	execId := resp.ID

	url, err := url.Parse(a.swarmClassicURL)
	if err != nil {
		log.Errorf("error parsing swarm URL for hijack: %s", err)
		return
	}

	if err := a.hijack(url.Host, "POST", "/exec/"+execId+"/start", true, ws, ws, ws, nil, nil); err != nil {
		log.Errorf("error during hijack: %s", err)
		return
	}

	// resize
	w, err := strconv.Atoi(ttyWidth)
	if err != nil {
		log.Error(err)
		return
	}

	h, err := strconv.Atoi(ttyHeight)
	if err != nil {
		log.Error(err)
		return
	}

	if err := a.manager.DockerClient().ContainerExecResize(context.TODO(), execId, types.ResizeOptions{Width: w, Height: h}); err != nil {
		log.Errorf("error resizing exec tty: %s", err)
		return
	}

}

func (a *Api) scaleContainer(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	w.Header().Set("content-type", "application/json")

	n := rc.QueryVars["n"]

	if len(n) == 0 {
		http.Error(w, "you must enter a number of instances (param: n)", http.StatusBadRequest)
		return
	}

	numInstances, err := strconv.Atoi(n[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if numInstances <= 0 {
		http.Error(w, "you must enter a positive value", http.StatusBadRequest)
		return
	}

	result := a.manager.ScaleContainer(rc.PathVars["id"], numInstances)
	// If we received any errors, continue to write result to the writer, but return a 500
	if len(result.Errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) listContainers(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	allParam := rc.QueryVars.Get("all")
	all := false
	if allParam != "" {
		all = true
	}
	s := rc.QueryVars.Get("size")
	size := false
	if s != "" {
		size = true
	}

	args, err := dockerfilters.FromParam(rc.QueryVars.Get("filters"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filters, err := dockerfilters.ToParam(args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allContainers, err := a.manager.ListUserContainers(rc.Auth, all, size, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	containers := []types.Container{}

	// filter out system containers
	if all {
		containers = allContainers
	} else {
		for _, c := range allContainers {
			if _, ok := c.Labels["com.docker.ucp.InstanceID"]; !ok {
				containers = append(containers, c)
			}
		}
	}

	if err := json.NewEncoder(w).Encode(containers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Api) createContainer(w http.ResponseWriter, rc *ctx.OrcaRequestContext) {
	// Identify the Container Resource within the request context
	var containerReq *resources.ContainerResourceRequest
	var found bool
	var effectiveRole auth.Role

	// Attempt to cast the MainResource as a CRUDResourceRequest
	crudReq, foundCrud := rc.MainResource.(*resources.CRUDResourceRequest)
	if !foundCrud {
		http.Error(w, fmt.Sprintf("internal error: unable to locate CRUD resource in context"), http.StatusInternalServerError)
	}

	// Check if the underlying Labelled Resource can be cast to a container resource
	containerReq, found = crudReq.LabelledResource.(*resources.ContainerResourceRequest)
	if !found {
		http.Error(w, fmt.Sprintf("internal error: unable to locate container resource in context"), http.StatusInternalServerError)
		return
	}

	container := containerReq.ContainerCreateRequest
	effectiveRole = crudReq.GetEffectiveRole()

	// TODO(alexmavr): Set the admin user's role to auth.Admin globally
	if rc.Auth.User.Admin {
		effectiveRole = auth.Admin
	}

	// Make sure we've got a label for the userid
	if container.Config.Labels == nil {
		container.Config.Labels = make(map[string]string)
	}
	container.Config.Labels[orca.UCPOwnerLabel] = rc.Auth.User.Username

	// TODO - do any other mutations of the request here (notary, etc.)

	rejectedMessages := []string{}
	privileged := false
	capAdd := false
	ipcMode := false
	pidMode := false
	hostBindMounts := false

	log.Debugf("checking container config for user: %s", rc.Auth.User.Username)

	// check for permissions beyond restricted

	// NOTE: all of the following checks would be a better fit for the container resource.
	// However, the check for allowUCPScheduling further below still requires the effectiveRole check
	// TODO(alexmavr): address this in a separate pass
	if container.HostConfig.Privileged {
		privileged = true
		rejectedMessages = append(rejectedMessages, "privileged not allowed")
	}

	if len(container.HostConfig.CapAdd) > 0 {
		capAdd = true
		rejectedMessages = append(rejectedMessages, "adding capabilities not allowed")
	}

	if container.HostConfig.IpcMode != "" {
		ipcMode = true
		rejectedMessages = append(rejectedMessages, "ipc not allowed")
	}

	if container.HostConfig.PidMode != "" {
		pidMode = true
		rejectedMessages = append(rejectedMessages, "pid not allowed")
	}

	for _, bind := range container.HostConfig.Binds {
		// TODO: WINDOWS support
		parts := strings.Split(bind, ":")
		vol := parts[0]
		if strings.Index(vol, "/") == 0 {
			hostBindMounts = true
			rejectedMessages = append(rejectedMessages, "host mounted volumes not allowed")
			break
		}
	}

	if effectiveRole <= auth.RestrictedControl {
		if privileged || capAdd || ipcMode || pidMode || hostBindMounts {
			log.Warnf("permission denied to restricted: user=%s", rc.Auth.User.Username)
			log.Debugf("user role for restricted: %d", effectiveRole)

			http.Error(w, "Permission denied: "+strings.Join(rejectedMessages, ", ")+".  You must specify a group label with Full Access or request Full Access from your administrator.", http.StatusBadRequest)
			return
		}
	}

	var allowUCPScheduling bool
	if effectiveRole < auth.Admin {
		allowUCPScheduling = a.manager.EnableUserUCPScheduling()
	} else {
		allowUCPScheduling = a.manager.EnableAdminUCPScheduling()
	}

	// check to allow scheduling anywhere
	if !allowUCPScheduling {
		ucpControllerName := config.OrcaControllerContainerName
		container.Config.Env = append(container.Config.Env,
			fmt.Sprintf("affinity:container!=%s", ucpControllerName),
			"affinity:container!=dtr-api-*",
		)
	}

	data, err := json.Marshal(container)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Re-forge the request
	rc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	rc.Request.ContentLength = int64(len(data))
	a.swarmRegistryRedirect(w, rc, container.Config.Image, "pull")
}

func (a *Api) containerLogs(ws *websocket.Conn) {
	qry := ws.Request().URL.Query()
	containerId := qry.Get("id")
	token := qry.Get("token")

	opts := types.ContainerLogsOptions{
		Details:    paramBoolValue(qry.Get("details")),
		Follow:     paramBoolValue(qry.Get("follow")),
		ShowStdout: paramBoolValue(qry.Get("stdout")),
		ShowStderr: paramBoolValue(qry.Get("stderr")),
		Since:      qry.Get("since"),
		Timestamps: paramBoolValue(qry.Get("timestamps")),
		Tail:       qry.Get("tail"),
	}

	if !a.manager.ValidateContainerLogsToken(containerId, token) {
		ws.Write([]byte("unauthorized"))
		ws.Close()
		return
	}

	log.Debugf("starting logs session: container=%s token=%s", containerId, token)

	c, err := a.manager.Container(containerId)

	responseBody, err := a.manager.ContainerLogs(containerId, opts)
	if err != nil {
		ws.Write([]byte(err.Error()))
		ws.Close()
		return
	}
	defer responseBody.Close()

	if c.Config.Tty {
		_, err = io.Copy(ws, responseBody)
	} else {
		_, err = stdcopy.StdCopy(ws, ws, responseBody)
	}

}
