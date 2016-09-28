package orca

import (
	"time"
)

type Node struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Addr           string    `json:"addr,omitempty"`
	Containers     string    `json:"containers,omitempty"`
	ReservedCPUs   string    `json:"reserved_cpus,omitempty"`
	ReservedMemory string    `json:"reserved_memory,omitempty"`
	Labels         []string  `json:"labels,omitempty"`
	ResponseTime   float64   `json:"response_time"`
	Error          string    `json:"error,omitempty"`
	Status         string    `json:"status,omitempty"`
	ServerVersion  string    `json:"server_version,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// EOL - only used by 1.0.X bootstrappers
// If we detect this, we'll always error and tell the user to use a new bootstrapper
type V1NodeRequest struct {
	CertificateRequest string `json:"certificate_request,omitempty"`
	Name               string `json:"name,omitempty"`

	// Indicate if this request is for a controller replica, or just a swarm signed cert
	Replica bool `json:"replica,omitempty"`
}

// Legacy - Combined request.
// Only used by 1.1.x bootstrappers
type NodeRequest struct {
	// Requests for signatures by the cluster CA
	// The key is used for the output map to allow the caller to keep the key/cert consistent
	ClusterCertificateRequests map[string]string `json:"cluster_certificate_requests"`
	// Requests for signatures by the user CA
	// The key is used for the output map to allow the caller to keep the key/cert consistent
	UserCertificateRequests map[string]string `json:"user_certificate_requests"`

	// Set to true to indicate this request is to set up a new intermediate CA for a replica node
	Replica bool `json:"replica"`
}

// Response for a Join operation
type NodeConfiguration struct {
	ClusterCertificateChain string            `json:"cluster_certificate_chain"`
	ClusterCertificates     map[string]string `json:"cluster_certificates"` // Keys match request
	UserCertificateChain    string            `json:"user_certificate_chain"`
	UserCertificates        map[string]string `json:"user_certificates"` // Keys match request
	SwarmArgs               []string          `json:"swarm_args"`
	OrcaID                  string            `json:"orca_id"`
	KvStore                 []string          `json:"kv_store"`
	Warnings                string            `json:"warnings"` // Any messages to display as a result of the join
}
