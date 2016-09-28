package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/orca/registry/v2"
	"github.com/docker/orca/version"
	ver "github.com/hashicorp/go-version"
	"golang.org/x/net/context"
)

var updateBanner = ""

const (
	RegistryAuthURL = "https://auth.docker.io/token"
	RegistryService = "registry.docker.io"
	RegistryURL     = "https://registry-1.docker.io/v2"
	CheckUpdateRepo = "docker/ucp"
)

type TagList struct {
	Tags []string `json:"tags"`
}

func (m DefaultManager) periodicCheckForUpdates() {
	orcaVersion, err := ver.NewVersion(version.Version)
	if err != nil {
		log.Warnf("Failed to parse ucp version: %s", err)
		return
	}

	// Get an anonymous token (the repo is assumed to be public)
	tokenURL := fmt.Sprintf("%s?scope=repository:%s:pull&service=%s", RegistryAuthURL, CheckUpdateRepo, RegistryService)
	resp, err := http.Get(tokenURL)
	if err != nil {
		log.Warnf("Failed to check for updates: %s", err)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		log.Infof("Auth Failed to check for updates: %s", string(data))
	}

	var token v2.AuthToken
	if err := json.Unmarshal(data, &token); err != nil {
		log.Infof("Auth Failed to check for updates: %s", err)
	}

	// Get the list of available tags
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/tags/list", RegistryURL, CheckUpdateRepo), nil)
	if err != nil {
		log.Warnf("Failed to check for updates: %s", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Infof("Failed to check for updates: %s", err)
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		log.Infof("Failed to check for updates: %s", string(data))
	}
	var tags TagList
	if err := json.Unmarshal(data, &tags); err != nil {
		log.Infof("Failed to check for updates: %s", err)
	}

	// Build up a list of all available updates newer than the current version
	newer := []string{}
	for _, tag := range tags.Tags {
		if tag == "latest" {
			continue
		}
		// Skip pre-release versions
		if strings.Contains(tag, "-") {
			log.Debugf("Skipping pre-release %s", tag)
			continue
		}
		testVersion, err := ver.NewVersion(tag)
		if err != nil {
			log.Debugf("Failed to parse tag %s: %s", tag, err)
			continue
		}
		if orcaVersion.LessThan(testVersion) {
			newer = append(newer, tag)
		}
	}

	if len(newer) > 0 {
		updateBanner = fmt.Sprintf("New updates are available for UCP: %s.  For instructions visit https://docs.docker.com/ucp/installation/upgrade/", strings.Join(newer, ", "))
	}
}

func (m DefaultManager) getVersionSkewBanners() []Banner {
	// Check for all running ucp-controllers and make sure their version match
	args := filters.NewArgs()
	args.Add("name", "ucp-controller")
	controllers, err := m.client.ContainerList(context.TODO(), types.ContainerListOptions{All: false, Size: false, Filter: args})
	if err != nil {
		log.Infof("Failed to lookup UCP-controllers for version comparison check: %s", err)
		return []Banner{}
	}
	if len(controllers) == 0 {
		log.Info("Failed to find any UCP-controllers for version comparison check")
		return []Banner{}
	}
	found := map[string]interface{}{}
	for _, c := range controllers {
		ver, ok := c.Labels["com.docker.ucp.version"]
		if ok {
			found[ver] = struct{}{}
		}
	}
	if len(found) > 1 {
		log.Infof("Detected upgrade in progress: %v", found)
		return []Banner{{
			Level:   BannerWARN,
			Message: fmt.Sprintf("Multiple running UCP controller versions detected - please complete your upgrade as soon as possible"),
		}}
	}
	return []Banner{}
}

func (m DefaultManager) getUpdateBanners() []Banner {
	if updateBanner != "" {
		return []Banner{{
			Level:   BannerINFO,
			Message: updateBanner,
		}}
	}
	return []Banner{}
}
