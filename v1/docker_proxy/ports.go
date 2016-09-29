package proxy

import (
	"log"

	"github.com/docker/engine-api/types"
)

// PortRewriter rewrites ports
type PortRewriter interface {
	// Rewrite(containerID, ports) is called when the client queries the exposed
	// ports, allowing us to change the addresses or set up port forwards
	RewritePorts(containerID string, ports []types.Port) []types.Port
}

// Adder adds ?
type Adder interface {
	// Add(containerID) is called asynchronously when a given container has started, and
	// needs resources setting up.
	Add(string) error
}

// Remover removes ?
type Remover interface {
	// Remove(containerID) is called asynchronously when a given container has died, and
	// needs resources cleaning up.
	Remove(string) error
}

// IPRewriter rewrites IP Addresses
type IPRewriter struct {
	dockerIP string
	lookupIP func() string
}

var _ PortRewriter = &IPRewriter{}

// NewIPRewriter returns a new IPRewriter
func NewIPRewriter(lookupIP func() string) *IPRewriter {
	return &IPRewriter{"", lookupIP}
}

// RewritePorts rewrites ports
func (r *IPRewriter) RewritePorts(containerID string, ports []types.Port) []types.Port {
	log.Printf("trying to rewrite IP addresses for port mapping for %s %#v\n", containerID, ports)
	if r.dockerIP == "" {
		r.dockerIP = r.lookupIP()
	}
	for i := range ports {
		if ports[i].IP == "0.0.0.0" {
			ports[i].IP = r.dockerIP
		}
	}
	return ports
}
