/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package osxTasks

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/go-p9p"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"github.com/docker/pinata/v1/pinataSockets"
	"golang.org/x/net/context"

	"runtime/debug"
)

// Typed API functions (for the unmarshalling and unmarshalling, see the
// AddRoute calls below)

const (
	driver = "com.docker.driver.amd64-linux"

	actionRestartVM            = "restartvm"
	actionSetVMSettings        = "setvmsettings"
	actionGetVMSettings        = "getvmsettings"
	actionVMStateEvent         = "vmstateevent"
	actionGetSharedDirectories = "getshareddirectories"
	actionSetSharedDirectories = "setshareddirectories"
	actionSetCertificates      = "setcertificates"
	statusOK                   = "ok"
	statusError                = "error"
)

// ErrorResponse describes an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (r ErrorResponse) Write(w http.ResponseWriter) {
	writeJSON(w, r)
}

// NewErrorResponse returns a new ErrorResponse with status set to error
func NewErrorResponse(message string) Response {
	return ErrorResponse{
		Status:  statusError,
		Message: message,
	}
}

// NoError is returned when no error occurred
var NoError = ErrorResponse{Status: statusOK}

// InternalError is returned when some unknown error has occurred
var InternalError = NewErrorResponse("internal error")

// VMSettings describes the settings that can be applied to a VM
type VMSettings struct {
	Memory               uint64 `json:"memory"`
	Cpus                 uint64 `json:"cpus"`
	DaemonJSON           string `json:"daemonjson"`
	SystemProxyHTTP      string `json:"systemProxyHttp"`
	SystemProxyHTTPS     string `json:"systemProxyHttps"`
	SystemProxyExclude   string `json:"systemProxyExclude"`
	OverrideProxyHTTP    string `json:"overrideProxyHttp"`
	OverrideProxyHTTPS   string `json:"overrideProxyHttps"`
	OverrideProxyExclude string `json:"overrideProxyExclude"`
}

// APIRequest represents an incoming json message (request from the OS X app)
type APIRequest struct {
	Action string `json:"action"`
}

// APIRestartVMRequest is a VM restart request from the OSX app
type APIRestartVMRequest struct {
	*APIRequest
}

// APISetVMSettingsRequest is a setvmsettings request from the OSX app
type APISetVMSettingsRequest struct {
	*APIRequest
	Args *VMSettings `json:"args"`
}

// APIGetVMSettingsRequest is a getvmsettings request from the OSX app
type APIGetVMSettingsRequest struct {
	*APIRequest
}

// GetVMSettingsResponse describes response to "getvmsettings" request
type GetVMSettingsResponse struct {
	Status string `json:"status"`
	*VMSettings
}

func (r GetVMSettingsResponse) Write(w http.ResponseWriter) {
	writeJSON(w, r)
}

// GetSharedDirectoriesRequest is a getshareddirectories request from the OSX app
type GetSharedDirectoriesRequest struct {
	*APIRequest
}

// SetSharedDirectoriesRequest is a setshareddirectories request from the OSX app
type SetSharedDirectoriesRequest struct {
	*APIRequest
	Args *SharedDirectoryArgs `json:"args"`
}

// SharedDirectoryArgs contains a list of directories
type SharedDirectoryArgs struct {
	Directories []string `json:"directories"`
}

// SharedDirectoriesResponse is a response to a get/setshareddirectories request
type SharedDirectoriesResponse struct {
	Status      string   `json:"status"`
	Directories []string `json:"directories"`
}

func (r SharedDirectoriesResponse) Write(w http.ResponseWriter) {
	writeJSON(w, r)
}

// VMStateEventRequest is a request for VM state from the OSX app
type VMStateEventRequest struct {
	*APIRequest
	Args *VMStateEvent
}

// VMStateEvent is a VM state
type VMStateEvent struct {
	VMState string `json:"vmstate"`
}

// VMStateEventResponse is a response to a VM State Event request
type VMStateEventResponse struct {
	Status  string `json:"status"`
	VMState string `json:"vmstate"`
}

// SetCertificatesRequest is a setcertificates request from the OSX app
type SetCertificatesRequest struct {
	*APIRequest
	Value string `json:"value"`
}

func (r VMStateEventResponse) Write(w http.ResponseWriter) {
	writeJSON(w, r)
}

// Response interface encapsulates a "Write" behaviour that
// writes a response to a supplied http.ResponseWriter
type Response interface {
	Write(w http.ResponseWriter)
}

// Write will write a JSON response to the ResponseWriter
func writeJSON(w http.ResponseWriter, d interface{}) {
	data, err := json.Marshal(d)
	if err != nil {
		logrus.Errorf("Failed to unmarshal response object %+v", d)
		res := NewErrorResponse("failed to marshal response object")
		data, _ = json.Marshal(res)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(data)
}

const whitespace = " \r\n\t"

var enoent = p9p.MessageRerror{Ename: "No such file or directory"}

func vmFilesPath() string {
	return filepath.Join(appleutil.GetContainerPath(), "database", "com.docker.driver.amd64-linux")
}

func newWatch(ctx context.Context, client *datakit.Client) *datakit.Watch {
	watch, err := datakit.NewWatch(ctx, client, "master", []string{driver, "state"})
	// NOTE(aduermael): I had to add this work around, because watching on a path
	// that doesn't exist returns an error. It should be in the 9p API
	for err != nil {
		logrus.Errorf("Failed to watch the %s/state/ directory in the database: %s", driver, err)
		time.Sleep(time.Second)
		watch, err = datakit.NewWatch(ctx, client, "master", []string{driver, "state"})
	}
	return watch
}

// watch for startup/shutdown events
func watchForStateEvents(ctx context.Context, client *datakit.Client, stateEventsChannel chan<- string) {
	watch := newWatch(ctx, client)

	lastStartTime := time.Unix(0, 0)
	lastShutdownTime := time.Unix(0, 0)
	for {
		snapshot, err := watch.Next(ctx)
		if err != nil {
			logrus.Printf("Failed to watch for updates: %s", err)
			time.Sleep(time.Second)
			watch.Close(ctx)
			watch = newWatch(ctx, client)
			continue
		}
		newStartTime, err := readUnixTime(ctx, snapshot, []string{"last-start-time"})
		if err != nil {
			logrus.Printf("Failed to read the last-start-time key: %s", err)
		}
		newShutdownTime, err := readUnixTime(ctx, snapshot, []string{"last-shutdown-time"})
		if err != nil {
			logrus.Printf("Failed to read the last-shutdown-time key: %s", err)
		}

		// one or the other has changed
		if newStartTime != lastStartTime {
			lastStartTime = newStartTime
		}

		if newShutdownTime != lastShutdownTime {
			lastShutdownTime = newShutdownTime
		}

		if lastShutdownTime.Before(lastStartTime) {
			logrus.Printf("VM started at %s", newStartTime.String())
			stateEventsChannel <- "running"
		} else {
			logrus.Printf("VM shutdown at %s", newShutdownTime.String())
			stateEventsChannel <- "starting"
		}
	}
}

func readUnixTime(ctx context.Context, snapshot *datakit.Snapshot, path []string) (time.Time, error) {
	s, err := snapshot.Read(ctx, path)
	if err != nil {
		return time.Unix(0, 0), err
	}
	i, err := strconv.ParseInt(strings.Trim(s, " \r\t\n"), 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(i, 0), nil
}

// APIServer reexec entry point
func APIServer() {
	logrus.Println("API server starting")

	var err error

	// First, construct the path to the socket
	// deleting the socket if it exists
	var socketPath = filepath.Join(appleutil.GetContainerPath(), "s20")
	logrus.Println("🍀 socket path is:", socketPath)
	if _, err = os.Stat(socketPath); err == nil { // path to socket exists
		err = os.Remove(socketPath)
		if err != nil {
			logrus.Fatal("error while removing socket:", err)
		}
		logrus.Println("API socket removed")
	}

	apiListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		logrus.Fatal(err)
	}

	// share one database connection
	ctx := context.TODO()
	var client *datakit.Client
	attemptsTotal := 100
	attemptsRemaining := attemptsTotal
	for attemptsRemaining > 0 {
		client, err = datakit.Dial(ctx, "unix", pinataSockets.GetDBSocketPath())
		if err == nil {
			break
		}
		logrus.Printf("Try %d/%d failed to connect to db: %s\n",
			attemptsTotal-attemptsRemaining, attemptsTotal, err)
		time.Sleep(100 * time.Millisecond)
		attemptsRemaining = attemptsRemaining - 1
	}
	if client == nil {
		logrus.Fatalf("Failed to connect to the database after %d attempts", attemptsTotal)
	}

	stateEvents := make(chan string)
	go watchForStateEvents(ctx, client, stateEvents)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var apiRequest APIRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Error("Unable to read body of API request")
		}
		err = json.Unmarshal(body, &apiRequest)
		if err != nil {
			logrus.Error("can't read API request:", err)
		}

		var response Response

		defer func() {
			exception := recover()
			if exception != nil {
				logrus.Errorf("API handler panicked: %#v\n%s",
					exception, debug.Stack())
				response = NewErrorResponse(fmt.Sprintf("API handler panicked: %#v", exception))
				response.Write(w)
			}
		}()

		switch apiRequest.Action {
		case actionVMStateEvent:
			var clientClosed <-chan bool
			cn, ok := w.(http.CloseNotifier)
			if ok {
				clientClosed = cn.CloseNotify() // returns a chan bool, written to when/if this http request is closed
				var r VMStateEventRequest
				err := json.Unmarshal(body, &r)
				if err != nil {
					logrus.Errorf("can't decode VMStateEventRequest: %s", err)
					response = NewErrorResponse("can't decode vmstateevent request")
					break
				}
				response = apiRequestVMStateEvent(r, ctx, client, stateEvents, clientClosed)
			}
		case actionGetVMSettings:
			var r APIGetVMSettingsRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Error("can't decode GetVMSettingsRequest")
				response = NewErrorResponse("can't decode getvmsettings request")
			}
			response = apiRequestGetVMSettings(r, ctx, client)
		case actionSetVMSettings:
			var r APISetVMSettingsRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Errorf("can't decode SetVMSettingsRequest: %s", err)
				response = NewErrorResponse("can't decode setvmsettings request")
				break
			}
			response = apiRequestSetVMSettings(r, ctx, client)
		case actionRestartVM:
			var r APIRestartVMRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Error("can't decode RestartVMRequest")
				response = NewErrorResponse("can't decode restartvm request")
				break
			}
			response = apiRequestRestartVM(r, ctx, client)
		case actionGetSharedDirectories:
			var r GetSharedDirectoriesRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Error("can't decode GetSharedDirectoriesRequest")
				response = NewErrorResponse("can't decode getshareddirectories request")
				break
			}
			response = apiRequestGetSharedDirectories(r, ctx, client)
		case actionSetSharedDirectories:
			var r SetSharedDirectoriesRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Error("can't decode SetSharedDirectoriesRequest")
				response = NewErrorResponse("can't decode setshareddirectories request")
				break
			}
			response = apiRequestSetSharedDirectories(r, ctx, client)
		case actionSetCertificates:
			var r SetCertificatesRequest
			err := json.Unmarshal(body, &r)
			if err != nil {
				logrus.Error("can't decode SetCertificatesRequest")
				response = NewErrorResponse("can't decode setcertificates request")
				break
			}
			response = apiRequestSetCertificates(r, ctx, client)
		default:
			logrus.Error("api: unknown action")
			response = NewErrorResponse("unknown action")
		}
		response.Write(w)
	})

	err = http.Serve(apiListener, nil)
	if err != nil {
		logrus.Error("backend api:", err)
	}
}

func apiRequestGetVMSettings(request APIGetVMSettingsRequest, ctx context.Context, client *datakit.Client) Response {

	sha, err := datakit.Head(ctx, client, "master")
	if err != nil {
		logrus.Error("getvmsettings: db access error:", err)
		return NewErrorResponse("database access error")
	}
	snap := datakit.NewSnapshot(ctx, client, datakit.COMMIT, sha)
	// read "memory" setting value
	memoryW, err := snap.Read(ctx, []string{driver, "memory"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (memory):", err)
		return NewErrorResponse("database access error")
	}
	memoryStr := strings.Trim(memoryW, whitespace)
	// read "ncpu" setting value
	cpusW, err := snap.Read(ctx, []string{driver, "ncpu"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (ncpu):", err)
		return NewErrorResponse("database access error")
	}
	cpusStr := strings.Trim(cpusW, whitespace)
	// read "daemon.json" setting value
	daemonjsonW, err := snap.Read(ctx, []string{driver, "etc", "docker", "daemon.json"})
	if err != nil && err != enoent {
		logrus.Fatalf("Failed to read daemon.json from snapshot: %s", err)
	}
	daemonjsonStr := strings.Trim(daemonjsonW, whitespace)
	// read "proxy-system/http" setting value
	systemProxyHTTPW, err := snap.Read(ctx, []string{driver, "proxy-system", "http"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-system/http):", err)
		return NewErrorResponse("database access error")
	}
	systemProxyHTTPStr := strings.Trim(systemProxyHTTPW, whitespace)
	// read "proxy-system/https" setting value
	systemProxyHTTPSW, err := snap.Read(ctx, []string{driver, "proxy-system", "https"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-system/https):", err)
		return NewErrorResponse("database access error")
	}
	systemProxyHTTPSStr := strings.Trim(systemProxyHTTPSW, whitespace)
	// read "proxy-system/exclude" setting value
	systemProxyExcludeW, err := snap.Read(ctx, []string{driver, "proxy-system", "exclude"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-system/exclude):", err)
		return NewErrorResponse("database access error")
	}
	systemProxyExcludeStr := strings.Trim(systemProxyExcludeW, whitespace)
	// read "proxy-override/http" setting value
	overrideProxyHTTPW, err := snap.Read(ctx, []string{driver, "proxy-override", "http"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-override/http):", err)
		return NewErrorResponse("database access error")
	}
	overrideProxyHTTPStr := strings.Trim(overrideProxyHTTPW, whitespace)
	// read "proxy-override/https" setting value
	overrideProxyHTTPSW, err := snap.Read(ctx, []string{driver, "proxy-override", "https"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-override/https):", err)
		return NewErrorResponse("database access error")
	}
	overrideProxyHTTPSStr := strings.Trim(overrideProxyHTTPSW, whitespace)
	// read "proxy-override/exclude" setting value
	overrideProxyExcludeW, err := snap.Read(ctx, []string{driver, "proxy-override", "exclude"})
	if err != nil && err != enoent {
		logrus.Error("getvmsettings: db access error (proxy-override/exclude):", err)
		return NewErrorResponse("database access error")
	}
	overrideProxyExcludeStr := strings.Trim(overrideProxyExcludeW, whitespace)

	// validate that memoryStr is a uint64
	memory, err := strconv.ParseUint(memoryStr, 10, 64)
	if err != nil {
		// in beta 24 we accidentally released a version which wrote single bytes
		memoryB := []byte(memoryW)
		if len(memoryB) == 1 {
			memory = uint64(memoryB[0])
			logrus.Warn("getvmsettings: single byte memory setting detected (%d)", memory)
			// rewrite the broken database value
			t, err := datakit.NewTransaction(ctx, client, "master", "fix broken memory setting")
			if err != nil {
				logrus.Error("getvmsettings: failed to create transaction to fix broken database value")
				return NewErrorResponse(err.Error())
			}
			if err = t.Write(ctx, []string{driver, "memory"}, fmt.Sprintf("%d", memory)); err != nil {
				logrus.Error("getvmsettings: failed to write new value to fix broken database value")
				return NewErrorResponse(err.Error())
			}
			err = t.Commit(ctx, "fix broken memory setting")
			if err != nil {
				logrus.Error("getvmsettings: failed to commit transaction to fix broken database value")
				return NewErrorResponse(err.Error())
			}
		} else {
			logrus.Error("getvmsettings: failed to convert memory string into uint64:", err)
			return NewErrorResponse("failed to get memory setting")
		}
	}
	// validate that cpus is a uint64
	cpus, err := strconv.ParseUint(cpusStr, 10, 64)
	if err != nil {
		// in beta 24 we accidentally released a version which wrote single bytes
		cpusB := []byte(cpusW)
		if len(cpusB) == 1 {
			cpus = uint64(cpusB[0])
			logrus.Warn("getvmsettings: single byte cpus setting detected (%d)", cpus)
			// rewrite the broken database value
			t, err := datakit.NewTransaction(ctx, client, "master", "fix broken cpu setting")
			if err != nil {
				logrus.Error("getvmsettings: failed to create transaction to fix broken database value")
				return NewErrorResponse(err.Error())
			}
			if err = t.Write(ctx, []string{driver, "ncpu"}, fmt.Sprintf("%d", cpus)); err != nil {
				logrus.Error("getvmsettings: failed to write new value to fix broken database value")
				return NewErrorResponse(err.Error())
			}
			err = t.Commit(ctx, "fix broken cpus setting")
			if err != nil {
				logrus.Error("getvmsettings: failed to commit transaction to fix broken database value")
				return NewErrorResponse(err.Error())
			}
		} else {
			logrus.Error("getvmsettings: failed to convert cpus string into uint64:", err)
			return NewErrorResponse("failed to get cpus setting")
		}
	}

	response := GetVMSettingsResponse{
		Status: statusOK,
		VMSettings: &VMSettings{
			Memory:               memory,
			Cpus:                 cpus,
			DaemonJSON:           daemonjsonStr,
			SystemProxyHTTP:      systemProxyHTTPStr,
			SystemProxyHTTPS:     systemProxyHTTPSStr,
			SystemProxyExclude:   systemProxyExcludeStr,
			OverrideProxyHTTP:    overrideProxyHTTPStr,
			OverrideProxyHTTPS:   overrideProxyHTTPSStr,
			OverrideProxyExclude: overrideProxyExcludeStr,
		},
	}

	return response
}

func apiRequestSetVMSettings(request APISetVMSettingsRequest, ctx context.Context, client *datakit.Client) Response {
	// input: {"action":"setvmsettings", "args": {"memory":"2","cpus":"2","network":"native"}}
	// output: {"status":"ok"}

	// TODO(djs55): improve the protocol here: how do we reliably tell when the
	// operation has finished? Perhaps we need to watch a "live" branch which is
	// updated by xhyve or Moby

	// TODO(dave-tucker): Frontend is only sending values that have to be changed
	// We should probably make them send the entire config and make this idempotent

	t, err := datakit.NewTransaction(ctx, client, "master", "setvmsettings")
	if err != nil {
		return NewErrorResponse(err.Error())
	}

	if request.Args.Memory > 0 {
		memory := fmt.Sprintf("%d", request.Args.Memory)
		if err = t.Write(ctx, []string{driver, "memory"}, memory); err != nil {
			return NewErrorResponse("Can't set memory: " + err.Error())
		}
	}

	if request.Args.Cpus > 0 {
		cpus := fmt.Sprintf("%d", request.Args.Cpus)
		if err = t.Write(ctx, []string{driver, "ncpu"}, cpus); err != nil {
			return NewErrorResponse("Can't set cpus: " + err.Error())
		}
	}

	if request.Args.DaemonJSON != "" {
		if err = t.Write(ctx, []string{driver, "etc", "docker", "daemon.json"}, request.Args.DaemonJSON); err != nil {
			return NewErrorResponse("Can't set daemon opts: " + err.Error())
		}
	}

	// if overrideProxyHTTP value is provided, we apply it
	// if overrideProxyHTTP value is provided, we apply it
	// if overrideProxyExclude value is provided, we apply it
	// if they are all empty, we delete the override fields in DB

	if request.Args.OverrideProxyHTTP == "" && request.Args.OverrideProxyHTTPS == "" && request.Args.OverrideProxyExclude == "" {
		if err := t.Remove(ctx, []string{driver, "proxy-override", "http"}); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
		if err := t.Remove(ctx, []string{driver, "proxy-override", "https"}); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
		if err := t.Remove(ctx, []string{driver, "proxy-override", "exclude"}); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
	} else { // if at least one of them is NOT empty, we write the three values in DB
		if err := t.Write(ctx, []string{driver, "proxy-override", "http"}, request.Args.OverrideProxyHTTP); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
		if err := t.Write(ctx, []string{driver, "proxy-override", "https"}, request.Args.OverrideProxyHTTPS); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
		if err := t.Write(ctx, []string{driver, "proxy-override", "exclude"}, request.Args.OverrideProxyExclude); err != nil {
			return NewErrorResponse("Can't set proxy override settings: " + err.Error())
		}
	}

	if err = t.Commit(ctx, "setvmsettings"); err != nil {
		return NewErrorResponse("Can't apply settings: " + err.Error())
	}

	// changes commited
	// we don't have to wait for VM to reboot,
	// the UI application is tracking start / stop events on the side
	return NoError
}

// SetMemory sets the amount of memory available to the backend
func SetMemory(ctx context.Context, client *datakit.Client, memory string) error {
	t, err := datakit.NewTransaction(ctx, client, "master", "setvmsettings")
	if err != nil {
		return err
	}

	if err = t.Write(ctx, []string{driver, "memory"}, memory); err != nil {
		return err
	}

	return t.Commit(ctx, "setvmsettings")
}

// UI sends to vmstateevent a state that is considered as current state (for the UI)
// if the actual state is different or as soon as it changes, we answer with state name
func apiRequestVMStateEvent(request VMStateEventRequest, ctx context.Context, client *datakit.Client, stateEventsChannel <-chan string, clientClosed <-chan bool) Response {
	// input: {"action":"vmstateevent" ,"args" : { vmstate":"running" }} // possible values: starting, running
	// output: {"status":"ok","vmstate":"starting"} // possible values: starting, running
	currentVMstate := request.Args.VMState

	vmstate := getvmstate(ctx, client)

	if vmstate != currentVMstate {
		var response = VMStateEventResponse{
			Status:  statusOK,
			VMState: vmstate,
		}
		return response
	}

	// wait for state event and send response only if the new state is different
	for {
		select {
		// stop watching if client closed connection
		case <-clientClosed:
			return NewErrorResponse("client closed")
		case state := <-stateEventsChannel:
			if state != currentVMstate {
				var response = &VMStateEventResponse{
					Status:  statusOK,
					VMState: state,
				}
				return response
			}
		}
	}
}

//
func apiRequestRestartVM(request APIRestartVMRequest, ctx context.Context, client *datakit.Client) Response {
	err := restartVM(ctx, client)
	if err != nil {
		return NewErrorResponse(err.Error())
	}
	return NoError
}

///////////////////////////////// FUNCTIONS THAT ACTUALLY DO THE JOB ////////////////////////////////////////////

// gets list of mounted directories
func apiRequestGetSharedDirectories(request GetSharedDirectoriesRequest, ctx context.Context, client *datakit.Client) Response {
	// the "directories" entry is an "array of strings"
	// like this: "directories":["/a/path", "/another/path"]
	sha, err := datakit.Head(ctx, client, "master")
	if err != nil {
		logrus.Fatalf("Failed to discover the HEAD of master: %#v", err)
	}
	snap := datakit.NewSnapshot(ctx, client, datakit.COMMIT, sha)

	mounts, err := snap.Read(ctx, []string{driver, "mounts"})
	if err != nil && err != enoent {
		logrus.Fatalf("Failed to read memory from snapshot: %#v", err)
	}
	mounts = strings.Trim(mounts, "\n")
	var lines = make([]string, 0)
	if len(mounts) > 0 {
		// split string on the \n char (LF)
		lines = strings.Split(mounts, "\n")
	}
	// split each line on ':' char
	var directories = make([]string, 0)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			if parts[1] != "/host_docker_app" {
				directories = append(directories, parts[0])
			}
		}
	}
	response := SharedDirectoriesResponse{
		Status:      statusOK,
		Directories: directories,
	}

	return response
}

// sets list of mounted directories
func apiRequestSetSharedDirectories(request SetSharedDirectoriesRequest, ctx context.Context, client *datakit.Client) Response {

	stringSlice := request.Args.Directories

	// check proposed mounts using osxfs
	sharedDirectories, err := osxfsControlCheckProposedMounts(stringSlice)
	if err != nil {
		return NewErrorResponse(err.Error())
	}

	// generate string to be written in the DB
	var dbvalue string
	for _, element := range sharedDirectories {
		dbvalue += element + ":" + element + "\n"
	}
	dbvalue += appleutil.GetContainerPath() + ":/host_docker_app\n"

	response := SharedDirectoriesResponse{
		Status:      statusOK,
		Directories: sharedDirectories,
	}

	t, err := datakit.NewTransaction(ctx, client, "master", "setshareddirectories")
	if err != nil {
		return NewErrorResponse(err.Error())
	}
	if err = t.Write(ctx, []string{driver, "mounts"}, dbvalue); err != nil {
		return NewErrorResponse(err.Error())
	}
	if err = t.Commit(ctx, "setshareddirectories"); err != nil {
		return NewErrorResponse(err.Error())
	}

	// because we have already committed the changes and the error
	// response is only a message and not a variant right now, an error
	// during VM restart only gets logged
	// TODO: rollback, retry, or fail but don't revert GUI list
	if err := restartVM(ctx, client); err != nil {
		logrus.Errorf("Failed to restart VM after FS export change: %s", err)
	}
	return response
}

// copy certificates in etc/ssl/certs/ca-certificates.crt
func apiRequestSetCertificates(request SetCertificatesRequest, ctx context.Context, client *datakit.Client) Response {
	// input: {"action":"setcertificates", "value": "-----BEGIN CERTIFICATE-----"}
	// output: {"status":"ok"}
	t, err := datakit.NewTransaction(ctx, client, "master", "setcertificates")
	if err != nil {
		return NewErrorResponse(err.Error())
	}

	if request.Value != "" {
		if err = t.Write(ctx, []string{driver, "etc", "ssl", "certs", "ca-certificates.crt"}, request.Value); err != nil {
			return NewErrorResponse("Can't set certificates: " + err.Error())
		}
	}

	if err = t.Commit(ctx, "setcertificates"); err != nil {
		return NewErrorResponse("Can't set certificates: " + err.Error())
	}

	return NoError
}

//
func getvmstate(ctx context.Context, client *datakit.Client) (vmstate string) {
	sha, err := datakit.Head(ctx, client, "master")
	if err != nil {
		logrus.Fatalf("Failed to discover the HEAD of master: %s", err)
	}
	snap := datakit.NewSnapshot(ctx, client, datakit.COMMIT, sha)
	startTime, err := readUnixTime(ctx, snap, []string{driver, "state", "last-start-time"})
	if err != nil {
		// can't read start time, so the vm never finish starting yet
		vmstate = "starting"
		return
	}
	shutdownTime, err := readUnixTime(ctx, snap, []string{driver, "state", "last-shutdown-time"})
	if err != nil {
		// can't read shutdown, because maybe it never happened
		// so app should be running
		vmstate = "running"
		return
	}
	if shutdownTime.Before(startTime) {
		vmstate = "running"
	} else {
		vmstate = "starting"
	}
	return
}

// To restart the VM we look for the hypervisor's pid in its pid file
// and send a SIGTERM signal directly to it.
func restartVM(ctx context.Context, client *datakit.Client) error {
	// hypervisor's pid file is ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/pid
	pidFilePath := filepath.Join(appleutil.GetContainerPath(), "com.docker.driver.amd64-linux", "pid")
	if _, err := os.Stat(pidFilePath); os.IsNotExist(err) {
		// path does not exist
		return errors.New("could not locate hypervisor process")
	}
	// pid file exists. We read it.
	bytes, err := ioutil.ReadFile(pidFilePath)
	if err != nil {
		return err
	}
	// parse bytes and try converting them into an integer
	pid, err := strconv.ParseInt(string(bytes), 10, 32)
	if err != nil {
		return err
	}
	// try getting a handle on the process having this pid
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

const (
	osxfsCheckMountsOp = 0
)

type osxfsMessage struct {
	Len int32
	Op  int16
}

func osxfsReceiveStrings(osxfs net.Conn) ([]string, error) {
	var header osxfsMessage
	err := binary.Read(osxfs, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(osxfs)

	switch header.Op {
	case 0: // error
		errString, err := buf.ReadString('\000')
		if err != nil {
			return nil, err
		}
		logrus.Error(errString)
		return nil, errors.New(errString[0 : len(errString)-1])
	case 1: // strings
		goodPaths := []string{}
		for header.Len > 6 {
			path, err := buf.ReadString('\000')
			if err != nil {
				return nil, err
			}
			goodPaths = append(goodPaths, path[0:len(path)-1])
			header.Len = header.Len - int32(len(path))
		}
		return goodPaths, nil
	default:
		UnknownCodeErr := fmt.Sprintf("Unknown response code %d", header.Op)
		return nil, errors.New(UnknownCodeErr)
	}
}

func osxfsControlCheckProposedMounts(paths []string) ([]string, error) {
	osxfs, err := net.Dial("unix", pinataSockets.GetOsxfsControlSocketPath())
	if err != nil {
		logrus.Printf("Failed to connect to osxfs: %s\n", err)
	}

	length := 6 + len(paths)
	for _, path := range paths {
		length = length + len(path)
	}
	header := osxfsMessage{int32(length), osxfsCheckMountsOp}
	osxfsBuffered := bufio.NewWriterSize(osxfs, 4096) // yuck

	if err := binary.Write(osxfsBuffered, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	for _, path := range paths {
		_, err := osxfsBuffered.WriteString(path)
		if err != nil {
			return nil, err
		}
		err = osxfsBuffered.WriteByte(0)
		if err != nil {
			return nil, err
		}
	}

	if err := osxfsBuffered.Flush(); err != nil {
		return nil, err
	}

	return osxfsReceiveStrings(osxfs)
}
