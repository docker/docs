package orca

import (
	"time"

	"github.com/docker/engine-api/types"
)

type Event struct {
	Type          string               `json:"type,omitempty"`
	ContainerInfo *types.ContainerJSON `json:"container_info,omitempty"`
	Time          time.Time            `json:"time,omitempty"`
	Message       string               `json:"message,omitempty"`
	Username      string               `json:"username,omitempty"`
	RemoteAddr    string               `json:"remote_addr,omitempty"`
	Tags          []string             `json:"tags,omitempty"`
}
