package config

import (
	"fmt"
	"strings"

	"github.com/docker/orca/version"
)

// Process an image based on the current configuration

func GetContainerImage(containerName string) (string, error) {
	devMode := false
	ver := ImageVersion
	targetImage := Images[containerName]
	if targetImage == "" {
		return "", fmt.Errorf("Unrecognized container %s", containerName)
	}

	// Process the requested version
	if strings.HasPrefix(ImageVersion, "dev:") {
		devMode = true
		verSplit := strings.SplitN(ImageVersion, ":", 2)
		if verSplit[1] != "" {
			ver = verSplit[1]
		} else {
			ver = strings.Split(version.FullVersion(), " ")[0]
		}
	} else if ImageVersion != "" {
		ver = ImageVersion
	}

	// Append user supplied version if not pinned
	if !strings.Contains(targetImage, ":") {
		targetImage = targetImage + ":" + ver
	}

	// Now handle dev mode
	if devMode && strings.Contains(targetImage, "/") {
		splitImage := strings.SplitN(targetImage, "/", 2)
		targetImage = "dockerorcadev/" + splitImage[1]
	}
	return targetImage, nil
}
