package util

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var PrometheusServerName = "prometheus-server"
var PrometheusExporterPrefix = "prometheus-exporter"
var PrometheusServerExternalPort = "9090"
var PrometheusServerInternalPort = "9090"
var PrometheusExporterInternalPort = "9100"
var PrometheusExporterExternalPort = "9100"
var PrometheusExporterExternalPortInt = 9100

const PrometheusCfgTemplate = `
scrape_configs:
  - job_name: "node"
    scrape_interval: "15s"
    static_configs:
      - targets: [{{range .}}'{{.IP}}:{{.Port}}',{{end}}]
`

func ContainerIsPrometheusExporter(container types.Container) bool {
	for _, name := range container.Names {
		cleanName := strings.Split(name, "/")[2]
		if strings.HasPrefix(cleanName, PrometheusExporterPrefix) {
			return true
		}
	}
	return false
}

// return values: purged, server node
// if there's an existing prometheus and it's in good shape, we return purged false and the node the server is on
// "in good shape" means it has the right number of exporters and we are not told to purge it anyway
// if there's no prometheus left at the end of this operation, we return purged true (whether we removed it or it didn't exist)
func MaybePurgePrometheus(client *dc.Client, forcePurge bool, numExporters int) (bool, string) {
	list, err := client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	Expect(err).To(BeNil())

	broken := forcePurge
	serverInspect, err := client.ContainerInspect(context.Background(), PrometheusServerName)
	if err != nil && strings.Contains(err.Error(), "Error: No such container:") {
		log.Debugf("no prometheus server!")
		broken = true
	} else {
		Expect(err).To(BeNil())
	}

	exporterContainers := []types.Container{}
	for _, container := range list {
		if ContainerIsPrometheusExporter(container) {
			log.Debug("found prometheus exporter")
			exporterContainers = append(exporterContainers, container)
			if container.State != "running" {
				log.Debug("prometheus exporter is not running! It's %s", container.State)
				broken = true
			}
		}
	}
	log.Debugf("prometheus broken? %t %t, exporters: %d of %d", broken, numExporters != len(exporterContainers), len(exporterContainers), numExporters)
	if broken || numExporters != len(exporterContainers) {
		for _, container := range exporterContainers {
			log.Debug("removing prometheus exporter")
			client.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{Force: true})
		}
		client.ContainerRemove(context.Background(), PrometheusServerName, types.ContainerRemoveOptions{Force: true})
		return true, ""
	} else {
		return false, serverInspect.Node.Name
	}
}

// we can assume no prometheus exists at this point
func DeployPrometheus(client *dc.Client, prometheusServerConstraints, prometheusExporterConstraints []string, numPrometheusExporters int) string {
	nodes := []Node{}
	for i := 0; i < numPrometheusExporters; i++ {
		log.Infof("deploying prometheus exporter %d", i)
		node := RunContainerWithConfigFile(client, "prom/node-exporter", fmt.Sprintf("%s-%d", PrometheusExporterPrefix, i), PrometheusExporterInternalPort, PrometheusExporterExternalPort, "", "", prometheusExporterConstraints)

		IP, err := ha_utils.GetIP(node)
		Expect(err).To(BeNil())
		nodes = append(nodes, Node{IP: IP, Port: PrometheusExporterExternalPortInt})
	}

	configBuffer := new(bytes.Buffer)
	prometheusCfgTemplate := template.Must(template.New("prometheus.yml").Parse(PrometheusCfgTemplate))
	err := prometheusCfgTemplate.Execute(configBuffer, nodes)
	Expect(err).To(BeNil())

	log.Infof("deploying prometheus server")
	node := RunContainerWithConfigFile(client, "prom/prometheus", PrometheusServerName, PrometheusServerInternalPort, PrometheusServerExternalPort, configBuffer.String(), "/etc/prometheus/prometheus.yml", prometheusServerConstraints)
	return node
}
