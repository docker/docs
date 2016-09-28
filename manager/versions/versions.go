package versions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"

	log "github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"github.com/samalba/dockerclient"
)

type Checker interface {
	NewestVersion(*dockerclient.AuthConfig) (string, error)
	VersionList(*dockerclient.AuthConfig) (ManagerVersionList, error)
}

type versionChecker struct {
	settingsStore hubconfig.SettingsStore
}

func NewChecker(settingsStore hubconfig.SettingsStore) Checker {
	return &versionChecker{settingsStore: settingsStore}
}

func (v *versionChecker) NewestVersion(authConfig *dockerclient.AuthConfig) (string, error) {
	versionList, err := v.VersionList(authConfig)
	if err != nil {
		return "", err
	}
	return versionList[len(versionList)-1], nil
}

func (v *versionChecker) VersionList(authConfig *dockerclient.AuthConfig) (ManagerVersionList, error) {
	var err error
	if authConfig == nil {
		if authConfig, err = v.settingsStore.HubCredentials(); err != nil {
			if !deploy.IsProduction() {
				log.WithField("error", err).Error("Failed to retrieve stored hub credentials")
			}
		}

		if authConfig == nil {
			// Use empty auth config if still nil.
			authConfig = new(dockerclient.AuthConfig)
		}
	}

	managerRepoName := deploy.BootstrapRepo.Name()

	hubConfig, err := v.settingsStore.UserHubConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve hub configuration for release channel")
	}
	if hubConfig != nil && hubConfig.ReleaseChannel != "" {
		managerRepoName = deploy.ParseReleaseChannel(hubConfig.ReleaseChannel).ManagerRepoName()
	}

	repoTags, err := getRemoteTagList(authConfig, managerRepoName, "")
	if err != nil {
		return nil, err
	}
	var managerVersions []string
	for _, repoTag := range repoTags {
		if repoTag != "latest" {
			managerVersions = append(managerVersions, repoTag)
		}
	}
	if len(managerVersions) == 0 {
		log.Warn("No manager versions found")
		return nil, fmt.Errorf("No manager versions found")
	}
	sort.Sort(ManagerVersionList(managerVersions))
	log.WithField("versions", managerVersions).Info("Found available manager versions")
	return managerVersions, nil
}

func TagToSemver(v string) (semver.Version, error) {
	return semver.New(strings.Replace(v, "_", "+", 1))
}

func Less(v1, v2 string) bool {
	version1, err1 := TagToSemver(v1)
	version2, err2 := TagToSemver(v2)
	// TODO for GA we should be able to panic if we encounter any errors
	if err1 == nil && err2 == nil {
		return version1.LT(version2)
	} else if err1 == nil && err2 != nil {
		return false
	} else if err1 != nil && err2 == nil {
		return true
	} else {
		return v1 < v2
	}
}

type ManagerVersionList []string

func (l ManagerVersionList) Len() int {
	return len(l)
}
func (l ManagerVersionList) Less(i, j int) bool {
	return Less(l[i], l[j])
}
func (l ManagerVersionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ManagerVersionList) PatchFor(baseVersion string) (newestPatch string) {
	newestPatch = baseVersion
	candidateSemVer, err := semver.New(strings.Replace(baseVersion, "_", "+", 1))
	if err != nil {
		return
	}
	for _, v := range l {
		vSem, err := semver.New(strings.Replace(v, "_", "+", 1))
		if err != nil {
			continue
		}
		if candidateSemVer.Major == vSem.Major && candidateSemVer.Minor == vSem.Minor && candidateSemVer.Patch < vSem.Patch {
			candidateSemVer = vSem
			newestPatch = v
		}
	}
	return
}
