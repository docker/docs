package orca

import (
	"github.com/docker/engine-api/types"
)

type Service struct {
	Name       string            `json:"name,omitempty"`
	Containers []types.Container `json:"containers,omitempty"`
}

type Application struct {
	Name           string     `json:"name,omitempty"`
	Id             string     `json:"id,omitempty"`
	ConfigHash     string     `json:"config_hash,omitempty"`
	Services       []*Service `json:"services,omitempty"`
	ComposeVersion string     `json:"compose_version,omitempty"`
}
