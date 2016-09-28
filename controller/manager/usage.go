package manager

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/mssola/user_agent"
	"github.com/timehop/go-mixpanel"
	"golang.org/x/net/context"

	"github.com/docker/orca/version"
)

func (m DefaultManager) reportUsage() {
	err := m.reportUsageInner()
	if err != nil {
		log.Infof("Usage update err: %s", err)
	}
}

func (m DefaultManager) reportUsageInner() error {
	dclient := m.DockerClient()
	engineClient := m.ProxyClient()

	var numContainers, numImages, numServices, numNetworks, numVolumes int

	if containers, err := dclient.ContainerList(context.TODO(), types.ContainerListOptions{All: true, Size: false}); err == nil {
		numContainers = len(containers)
	}

	if images, err := engineClient.ImageList(context.TODO(), types.ImageListOptions{All: true}); err == nil {
		numImages = len(images)
	}

	if services, err := engineClient.ServiceList(context.TODO(), types.ServiceListOptions{}); err == nil {
		numServices = len(services)
	}

	if networks, err := engineClient.NetworkList(context.TODO(), types.NetworkListOptions{}); err == nil {
		numNetworks = len(networks)
	}

	if volumes, err := engineClient.VolumeList(context.TODO(), filters.Args{}); err == nil {
		numVolumes = len(volumes.Volumes)
	}

	nodes, err := m.Nodes()
	if err != nil {
		return err
	}
	numNodes := len(nodes)
	numManagers := len(m.GetManagers())

	// Hash the ID for anonymity
	id := fmt.Sprintf("%x", sha1.Sum([]byte(m.ID())))

	// Detect unlicensed state
	is_licensed := m.GetLicenseKeyID() != UnlicensedID
	license_tier, license_id := "", "anonymous"
	if is_licensed {
		license_tier = m.GetLicenseTier()
		if !m.AnonymizeTracking() {
			license_id = m.GetLicenseKeyID()
		}
	}

	cpus := 0.0
	memory := 0.0

	for _, n := range nodes {
		// TODO (ehazlett): we should add an endpoint
		// to swarm to return this via the api
		// swarm only gives back reservations in string
		// format; we must parse
		c := strings.Split(n.ReservedCPUs, "/")
		ms := strings.Split(n.ReservedMemory, "/")
		m := strings.Split(ms[1], " ")

		cpu := strings.TrimSpace(c[1])
		mem := strings.TrimSpace(m[1])

		nodeCpu, err := strconv.ParseFloat(cpu, 64)
		if err != nil {
			log.Warnf("error parsing cpu reservation from node: %s", err)
			continue
		}
		nodeMem, err := strconv.ParseFloat(mem, 64)
		if err != nil {
			log.Warnf("error parsing memory reservation from node: %s", err)
			continue
		}

		cpus += nodeCpu
		memory += nodeMem
	}

	versionInfo, err := dclient.ServerVersion(context.TODO())
	if err != nil {
		return err
	}

	swarmInfo, err := engineClient.SwarmInspect(context.TODO())
	if err != nil {
		return err
	}
	// Sanitize secrets in the policies
	swarmInfo.JoinTokens.Worker = ""
	swarmInfo.JoinTokens.Manager = ""

	// Content trust data
	dtrTrust := m.RequireContentTrustForDTR()
	hubTrust := m.RequireContentTrustForHub()

	usage := map[string]interface{}{
		"swarm_version":     versionInfo.Version,
		"container_count":   numContainers,
		"image_count":       numImages,
		"service_count":     numServices,
		"network_count":     numNetworks,
		"volume_count":      numVolumes,
		"cpus":              cpus,
		"memory":            memory,
		"node_count":        numNodes,
		"controller_count":  numManagers,
		"orca_version":      version.FullVersion(),
		"is_licensed":       is_licensed,
		"license_tier":      license_tier,
		"license_id":        license_id,
		"content_trust_dtr": dtrTrust,
		"content_trust_hub": hubTrust,
	}

	mp := &mixpanel.Mixpanel{
		Token:   MixpanelToken,
		BaseUrl: MixpanelUrl,
	}

	if err = mp.Track(id, "usage", usage); err != nil {
		return err
	}

	return nil
}

func sanitizeUA(useragent string, info map[string]interface{}) {
	ua := user_agent.New(useragent)
	info["Operating System"] = ua.OS()
	name, version := ua.Browser()
	versionNum := strings.Split(version, ".")[0]
	info["Browser"] = name
	info["Browser Version"] = versionNum
}

func (m DefaultManager) TrackClientInfo(req *http.Request) {
	if *m.disableTracking {
		return
	}

	go func(m DefaultManager, req *http.Request) {

		// Hash the ID for anonymity
		id := fmt.Sprintf("%x", sha1.Sum([]byte(m.ID())))
		info := make(map[string]interface{})
		sanitizeUA(req.UserAgent(), info)
		info["ip"] = req.RemoteAddr

		mp := &mixpanel.Mixpanel{
			Token:   MixpanelToken,
			BaseUrl: MixpanelUrl,
		}

		if err := mp.Track(id, "client info", info); err != nil {
			log.WithField("error", err).Error("Can't report client analytics")
		}
	}(m, req)
}
