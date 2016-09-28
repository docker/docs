package types

// serializable in env vars & json
type ClusterConfig struct {
	ControllerPort     string   `json:"controller_port"`
	SwarmPort          string   `json:"swarm_port"`
	TrustedRegistryCAs []string `json:"trusted_registry_cas"`
	ImageVersion       string   `json:"image_version"`
	Secret             string   `json:"secret"` // Go away you ugly beast!
	UCPInstanceID      string   `json:"ucp_instance_id"`
	DNS                []string `json:"dns"`
	DNSOpt             []string `json:"dns_opt"`
	DNSSearch          []string `json:"dns_search"`
}

// serializable in json
type NodeConfig struct {
	ClusterConfig
	IsManager     bool                      `json:"is_manager"`
	CertsExpiring bool                      `json:"certs_expiring"`
	HostAddress   string                    `json:"host_address"` // IP or Hostname of the node
	Managers      []string                  `json:"managers"`     // Set of manager IPs
	Containers    map[string]*OrcaContainer `json:"containers"`   // Map from container name to image/state
}

type OrcaContainer struct {
	Image   string `json:"image"`
	Running bool   `json:"running"`
}

type ReconcileConfig struct {
	Expected *NodeConfig
	Current  *NodeConfig
}

// TODO: Create a golang env serializer for the ClusterConfig?
