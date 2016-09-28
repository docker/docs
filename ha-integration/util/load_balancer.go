package util

import (
	"archive/tar"
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/gomega"
)

const HAProxyInternalPort = "443"

var HAProxyCfgTemplate = fmt.Sprintf(`
frontend localnodes
	bind *:%s
	mode tcp
	default_backend nodes
	timeout client 1m

backend nodes
	mode tcp
	balance roundrobin
	{{range .}}
	server clone{{.Index}} {{.IP}}:{{.Port}} check
	{{end}}
	timeout connect 10s
	timeout server 1m
`, HAProxyInternalPort)

type Node struct {
	IP string

	Index int
	Port  int
}

// XXX: the name must not start with dtr- or it'll be deleted by the dtr purge
const DTRLoadBalancerContainerName = "dtr_haproxy"
const UCPLoadBalancerContainerName = "ucp_haproxy"
const DTRLoadBalancerPort = 4443
const DefaultDTRNodePort = 443

// Idempotent removal of the container
func TearDownContainer(client *dc.Client, name string) {
	err := client.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		Expect(err.Error()).To(Equal(fmt.Sprintf("Error response from daemon: Container %s not found", name)))
	} else {
		Expect(err).To(BeNil())
	}
}

func RunContainerWithConfigFile(client *dc.Client, image, name, internalPort, externalPort, config, filePath string, constraints []string) string {
	// prepare the port config
	portBindings := make(nat.PortMap)
	portBindings[nat.Port(internalPort)] = append(portBindings[nat.Port(internalPort)], nat.PortBinding{HostIP: "0.0.0.0", HostPort: externalPort})

	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		Env:          constraints,
		Image:        image,
		ExposedPorts: map[nat.Port]struct{}{nat.Port(internalPort): {}},
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		PortBindings:  portBindings,
	}, &network.NetworkingConfig{}, name)
	Expect(err).To(BeNil())

	id := resp.ID

	if filePath != "" {
		configBuffer := strings.NewReader(config)
		filename := path.Base(filePath)
		dir := path.Dir(filePath)

		// prepare the file
		tarBuffer := new(bytes.Buffer)
		tarWriter := tar.NewWriter(tarBuffer)
		tarWriter.WriteHeader(&tar.Header{
			Name: filename,
			Mode: 0666,
			Size: int64(configBuffer.Len()),
		})
		_, err = configBuffer.WriteTo(tarWriter)
		Expect(err).To(BeNil())

		err = client.CopyToContainer(context.Background(), id, dir, tarBuffer, types.CopyToContainerOptions{})
		Expect(err).To(BeNil())
	}

	err = client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	Expect(err).To(BeNil())

	inspect, err := client.ContainerInspect(context.Background(), id)
	Expect(err).To(BeNil())

	node := inspect.Node.Name
	log.Infof("Your container is now set up on %s", node)
	return node

}

func SetupLoadBalancer(client *dc.Client, nodeDescriptions []Node, constraints []string, lbPort int, name string) string {
	configBuffer := new(bytes.Buffer)

	haproxyCfgTemplate := template.Must(template.New("haproxy.cfg").Parse(HAProxyCfgTemplate))

	err := haproxyCfgTemplate.Execute(configBuffer, nodeDescriptions)
	Expect(err).To(BeNil())

	config := string(configBuffer.Bytes())

	return RunContainerWithConfigFile(client, "haproxy", name, HAProxyInternalPort, fmt.Sprintf("%d", lbPort), config, "/usr/local/etc/haproxy/haproxy.cfg", constraints)
}

// schedule based on the given constraints, and if none given, put it on the first node
func SetupDefaultLoadBalancer(machines []ha_utils.Machine, client *dc.Client, constraints []string, lbPort, dstPort int, name string) string {
	log.Infof("Setting up HAProxy %s", name)

	TearDownContainer(client, name)

	IPs, err := getAllNodeIPs(machines)
	Expect(err).To(BeNil())

	nodes := []Node{}
	for i, IP := range IPs {
		nodes = append(nodes, Node{IP, i + 1, dstPort})
	}

	log.Infof("deploying load balancer with %d constraints: %v", len(constraints), constraints)
	if len(constraints) == 0 {
		nodeName := machines[0].GetName()
		log.Infof("No constraints given. Putting the load balancer on %s.", nodeName)
		constraints = []string{fmt.Sprintf("constraint:node==%s", nodeName)}
	}

	return SetupLoadBalancer(client, nodes, constraints, lbPort, name)
}

func getAllNodeIPs(machines []ha_utils.Machine) (IPs []string, err error) {
	for _, machine := range machines {
		ip, err := machine.GetIP()
		if err != nil {
			return nil, err
		}

		IPs = append(IPs, ip)
	}
	return
}

func getUCPControllerMachine(machines []ha_utils.Machine) ha_utils.Machine {
	return machines[0]
}
