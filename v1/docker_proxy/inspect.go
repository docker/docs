package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/docker/engine-api/types"
	"github.com/docker/go-connections/nat"
)

type inspectRewriter struct {
	portRewriter PortRewriter
}

// NewInspectRewriter returns a new Rewriter
func NewInspectRewriter(portRewriter PortRewriter) Rewriter {
	if portRewriter == nil {
		return &nopRewriter{}
	}

	return &inspectRewriter{
		portRewriter: portRewriter,
	}
}

func (p *inspectRewriter) RewritePortMap(ID string, input nat.PortMap) nat.PortMap {
	var ports []types.Port
	for natPort, natPortBindings := range input {
		for _, natPortBinding := range natPortBindings {
			publicPort, err := strconv.Atoi(natPortBinding.HostPort)
			if err != nil {
				log.Printf("Failed to parse HostPort: %s\n", natPortBinding.HostPort)
				return input // fall back to pass-through
			}
			port := types.Port{
				IP:          natPortBinding.HostIP,
				PrivatePort: natPort.Int(),
				PublicPort:  publicPort,
				Type:        natPort.Proto(),
			}
			ports = append(ports, port)
		}
	}
	ports = p.portRewriter.RewritePorts(ID, ports)
	results := make(nat.PortMap, 0)
	for _, port := range ports {
		private, err := nat.NewPort(port.Type, fmt.Sprintf("%d", port.PrivatePort))
		if err != nil {
			log.Printf("Failed to rewrite nat.PortMap: %#v\n", input)
			return input // fall back to pass-through
		}
		natPortBindings, ok := results[private]
		if !ok {
			natPortBindings = make([]nat.PortBinding, 0)
		}
		natPortBindings = append(natPortBindings, nat.PortBinding{
			HostIP:   port.IP,
			HostPort: fmt.Sprintf("%d", port.PublicPort),
		})
		results[private] = natPortBindings
	}
	return results
}

func (p *inspectRewriter) Rewrite(body io.ReadCloser) (int, io.ReadCloser) {
	buffer := bytes.NewBuffer([]byte{})

	_, err := io.Copy(buffer, body)
	if err != nil {
		log.Printf("Failed to read inspect body: %#v\n", err)
		panic(err)
	}

	var c types.ContainerJSON
	decoder := json.NewDecoder(buffer)
	if err := decoder.Decode(&c); err != nil {
		log.Printf("Failed to decode types.ContainerJSON: %#v\n", err)
		// return the data unmodified if we can't parse it
		return buffer.Len(), ioutil.NopCloser(buffer)
	}

	if p.portRewriter != nil {
		// ContainerJSON has a nat.PortMap rather than a []types.Port
		if c.NetworkSettings != nil {
			c.NetworkSettings.Ports = p.RewritePortMap(c.ID, c.NetworkSettings.Ports)
		}
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&c); err != nil {
		log.Printf("Failed to re-encode json: %#v\n", err)
		panic(err)
	}

	return b.Len(), ioutil.NopCloser(&b)
}
