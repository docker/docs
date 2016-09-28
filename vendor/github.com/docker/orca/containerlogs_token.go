package orca

type ContainerLogsToken struct {
	ContainerID string `json:"container_id,omitempty"`
	Token       string `json:"token,omitempty"`
}
