package utils

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	version "github.com/hashicorp/go-version"

	"github.com/docker/orca/config"
)

func GetUCPVersionString(containers []*types.ContainerJSON) string {
	// First look for the controller
	for _, container := range containers {
		switch container.Name {
		case "/" + config.OrcaControllerContainerName:
			ver, ok := container.Config.Labels["com.docker.ucp.version"]
			if ok {
				return ver
			}
		}
	}

	// Fall back to the proxy
	for _, container := range containers {
		switch container.Name {
		case "/" + config.OrcaProxyContainerName:
			ver, ok := container.Config.Labels["com.docker.ucp.version"]
			if ok {
				return ver
			}
			image := container.Config.Image
			if image != "" {
				s := strings.Split(image, ":")
				if len(s) == 2 {
					return s[1]
				}
			}
		}
	}
	log.Warning("Unable to determine current version")
	return "(version unknown)"
}

func LocalUCPVersion(containers []*types.ContainerJSON) (*version.Version, error) {
	versionString := GetUCPVersionString(containers)
	if strings.Contains(versionString, "unknown") {
		return nil, fmt.Errorf("Unable to determine Local UCP Version")
	}
	// Sometimes the version might be followed by a space and a commit hash
	stripVersion := strings.Split(versionString, " ")[0]

	v, err := version.NewVersion(stripVersion)
	if err != nil {
		return nil, err
	}

	return v, nil
}
