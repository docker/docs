package client

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	version "github.com/hashicorp/go-version"

	"github.com/docker/orca/bootstrap/config"
)

// Verify the docker engine is compatible with this version of orca
func (c *EngineClient) CheckDockerVersion() error {
	log.Debug("Checking for compatible engine version")

	verInfo, err := c.client.Version()
	if err != nil {
		return err
	}

	minVer, _ := version.NewVersion(config.MinVersion)
	// Drop the portion after the "-" since CS doesn't follow semantic versioning
	engineVer, _ := version.NewVersion(strings.Split(verInfo.Version, "-")[0])
	if engineVer.LessThan(minVer) {
		return fmt.Errorf("Your engine version %v is too old.  UCP requires at least version %v.", engineVer, minVer)
	}
	if !config.InPhase2 {
		log.Infof("Your engine version %s, build %s (%s) is compatible", verInfo.Version, verInfo.GitCommit, verInfo.KernelVersion)
	}
	return nil
}

// Determine whether the engine can be signaled for configuration changes in
// engine-discovery or not. This feature was added in CS 1.10 and 1.11 OSS, but
// CS 1.10 is buggy here. 1.10 OSS will do nothing with SIGHUP.
// Note that UCP's minimum engine is 1.10 for both CS and OSS so assume that as
// a floor.
func (c *EngineClient) EngineSupportsSignal() bool {
	verInfo, err := c.client.Version()
	// We already checked the version once when constructing this engine
	// client, so ignore version failures here
	if err != nil {
		return false
	}

	minVer, _ := version.NewVersion("1.11.0")
	// Drop the portion after the "-" since CS doesn't follow semantic versioning
	versionParts := strings.Split(verInfo.Version, "-")
	engineVer, _ := version.NewVersion(versionParts[0])
	if engineVer.LessThan(minVer) {
		return false
	}
	return true
}
