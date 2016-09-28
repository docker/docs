package orca

type (
	Usage struct {
		ID             string  `json:"id,omitempty"`
		Addr           string  `json:"addr,omitempty"`
		OrcaVersion    string  `json:"orca_version,omitempty"`
		SwarmVersion   string  `json:"swarm_version,omitempty"`
		ContainerCount int     `json:"container_count,omitempty"`
		ImageCount     int     `json:"image_count,omitempty"`
		Cpus           float64 `json:"cpus,omitempty"`
		Memory         float64 `json:"memory,omitempty"`
		NodeCount      int     `json:"node_count,omitempty"`
	}
)
