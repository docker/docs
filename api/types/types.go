package types

import (
	"time"

	"github.com/docker/docker/daemon/network"
	"github.com/docker/docker/pkg/version"
	"github.com/docker/docker/runconfig"
)

// ContainerCreateResponse contains the information returned to a client on the
// creation of a new container.
type ContainerCreateResponse struct {
	// ID is the ID of the created container.
	ID string `json:"Id"`

	// Warnings are any warnings encountered during the creation of the container.
	Warnings []string `json:"Warnings"`
}

// POST /containers/{name:.*}/exec
type ContainerExecCreateResponse struct {
	// ID is the exec ID.
	ID string `json:"Id"`
}

// POST /auth
type AuthResponse struct {
	// Status is the authentication status
	Status string `json:"Status"`
}

// POST "/containers/"+containerID+"/wait"
type ContainerWaitResponse struct {
	// StatusCode is the status code of the wait job
	StatusCode int `json:"StatusCode"`
}

// POST "/commit?container="+containerID
type ContainerCommitResponse struct {
	ID string `json:"Id"`
}

// GET "/containers/{name:.*}/changes"
type ContainerChange struct {
	Kind int
	Path string
}

// GET "/images/{name:.*}/history"
type ImageHistory struct {
	ID        string `json:"Id"`
	Created   int64
	CreatedBy string
	Tags      []string
	Size      int64
	Comment   string
}

// DELETE "/images/{name:.*}"
type ImageDelete struct {
	Untagged string `json:",omitempty"`
	Deleted  string `json:",omitempty"`
}

// GET "/images/json"
type Image struct {
	ID          string `json:"Id"`
	ParentId    string
	RepoTags    []string
	RepoDigests []string
	Created     int
	Size        int
	VirtualSize int
	Labels      map[string]string
}

// GET "/images/{name:.*}/json"
type ImageInspect struct {
	Id              string
	Parent          string
	Comment         string
	Created         time.Time
	Container       string
	ContainerConfig *runconfig.Config
	DockerVersion   string
	Author          string
	Config          *runconfig.Config
	Architecture    string
	Os              string
	Size            int64
	VirtualSize     int64
}

type LegacyImage struct {
	ID          string `json:"Id"`
	Repository  string
	Tag         string
	Created     int
	Size        int
	VirtualSize int
}

// GET  "/containers/json"
type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

type Container struct {
	ID         string            `json:"Id"`
	Names      []string          `json:",omitempty"`
	Image      string            `json:",omitempty"`
	Command    string            `json:",omitempty"`
	Created    int               `json:",omitempty"`
	Ports      []Port            `json:",omitempty"`
	SizeRw     int               `json:",omitempty"`
	SizeRootFs int               `json:",omitempty"`
	Labels     map[string]string `json:",omitempty"`
	Status     string            `json:",omitempty"`
}

// POST "/containers/"+containerID+"/copy"
type CopyConfig struct {
	Resource string
}

// GET "/containers/{name:.*}/top"
type ContainerProcessList struct {
	Processes [][]string
	Titles    []string
}

type Version struct {
	Version       string
	ApiVersion    version.Version
	GitCommit     string
	GoVersion     string
	Os            string
	Arch          string
	KernelVersion string `json:",omitempty"`
}

// GET "/info"
type Info struct {
	ID                 string
	Containers         int
	Images             int
	Driver             string
	DriverStatus       [][2]string
	MemoryLimit        bool
	SwapLimit          bool
	CpuCfsPeriod       bool
	CpuCfsQuota        bool
	IPv4Forwarding     bool
	Debug              bool
	NFd                int
	OomKillDisable     bool
	NGoroutines        int
	SystemTime         string
	ExecutionDriver    string
	LoggingDriver      string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	IndexServerAddress string
	RegistryConfig     interface{}
	InitSha1           string
	InitPath           string
	NCPU               int
	MemTotal           int64
	DockerRootDir      string
	HttpProxy          string
	HttpsProxy         string
	NoProxy            string
	Name               string
	Labels             []string
}

// This struct is a temp struct used by execStart
// Config fields is part of ExecConfig in runconfig package
type ExecStartCheck struct {
	// ExecStart will first check if it's detached
	Detach bool
	// Check if there's a tty
	Tty bool
}

type ContainerState struct {
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  time.Time
	FinishedAt time.Time
}

// GET "/containers/{name:.*}/json"
type ContainerJSON struct {
	Id              string
	Created         time.Time
	Path            string
	Args            []string
	Config          *runconfig.Config
	State           *ContainerState
	Image           string
	NetworkSettings *network.Settings
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Name            string
	RestartCount    int
	Driver          string
	ExecDriver      string
	MountLabel      string
	ProcessLabel    string
	Volumes         map[string]string
	VolumesRW       map[string]bool
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      *runconfig.HostConfig
}
