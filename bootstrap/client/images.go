package client

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	version "github.com/hashicorp/go-version"

	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/utils/registry"
	orcaconfig "github.com/docker/orca/config"
)

func (c *EngineClient) processImage(imageName string, messageDisplayed *bool) error {
	var progress io.ReadCloser
	var err error
	switch config.PullBehavior {
	case "always":
		log.Debugf("Pulling %s", imageName)
		progress, err = c.client.PullImage(imageName, types.ImagePullOptions{PrivilegeFunc: registry.RequestPrivilegeFunc})
		if err != nil {
			return err
		}
		defer progress.Close()
	case "missing":
		if _, err := c.client.InspectImage(imageName); err != nil {
			// TODO - Might be nice to detect different classes of pull errors and change our output...
			if !*messageDisplayed {
				*messageDisplayed = true
				// Dump out the safe registry user settings to help troubleshooting
				log.Debugf("REGISTRY_USERNAME=%s", os.Getenv("REGISTRY_USERNAME"))
				log.Debugf("REGISTRY_EMAIL=%s", os.Getenv("REGISTRY_EMAIL"))
				log.Info("Pulling required images... (this may take a while)")
			}
			log.Debugf("Pulling %s", imageName)
			progress, err = c.client.PullImage(imageName, types.ImagePullOptions{PrivilegeFunc: registry.RequestPrivilegeFunc})
			if err != nil {
				return err
			}
			defer progress.Close()

		}
	case "never":
		if _, err := c.client.InspectImage(imageName); err != nil {
			return fmt.Errorf("Image %s - %s", imageName, err)
		}
	}

	// TODO - this'll need some more work
	if progress != nil {
		data, err := ioutil.ReadAll(progress)
		if err != nil {
			return err
		}
		// This is really ugly, but trying to pull in the docker/pkg/jsonmessage
		// code breaks our build currently (winds up with native code for tty stuff)
		// so we'll probably have to implement it from scratch when/if we want to
		// give fancy progress reporting on download status
		// If we do it ourselves, we'll want to rip off the json structure for
		// the messages from that package, and then read it line-by-line, and
		// do something fancy with the output
		log.Debug(string(data))
	}

	return nil
}

// Make sure we have the images we need
func (c *EngineClient) VerifyOrPullImages(interactive bool) error {
	prompt := func() bool {
		config.InteractivePrompt("REGISTRY_USERNAME")
		config.InteractivePrompt("REGISTRY_PASSWORD")
		config.InteractivePrompt("REGISTRY_EMAIL")
		log.Info("Pulling required images... (this may take a while)")
		return true
	}

	promptedOnce := false
	messageDisplayed := false

	// Try not to prompt if we've already got the images we need.
	log.Debug("Checking for images")
	for name := range orcaconfig.Images {
		imageName, err := orcaconfig.GetContainerImage(name)
		if err != nil { // Shouldn't happen
			return err
		}
		if err := c.processImage(imageName, &messageDisplayed); err != nil {
			if config.PullBehavior != "never" && interactive && !promptedOnce {
				promptedOnce = prompt()
				if err := c.processImage(imageName, &messageDisplayed); err != nil {
					return err
				}
				continue
			}
			return err
		}
	}
	if !messageDisplayed {
		log.Info("All required images are present")
	}
	return nil
}

// Remove the staged images so we can grab them fresh after an uninstall
func (c *EngineClient) RemoveImages() error {
	anyFailed := false
	for name := range orcaconfig.Images {
		imageName, err := orcaconfig.GetContainerImage(name)
		if err != nil { // Shouldn't happen
			return err
		}
		if _, err := c.client.RemoveImage(imageName, false); !client.IsErrImageNotFound(err) && err != nil {
			// Workaround - seems sometimes IsErrImageNotFound fails but it was that error
			if strings.Contains(err.Error(), "No such image") {
				continue
			}
			log.Debugf("Unable to remove image %s: %s", name, err)
			anyFailed = true
		}
	}
	if anyFailed {
		return errors.New("One or more Orca images were still in use and couldn't be removed")
	}
	return nil
}

// Given two image IDs or names, see if we can safely upgrade from the old to the new image based on label metadata
func (c *EngineClient) CheckUpgradeCompatible(oldImage, newImage string) bool {
	oldVersion := ""
	upgradesFrom := ""
	newVersion := ""
	oldInfo, err := c.client.InspectImage(oldImage)
	if err != nil {
		log.Errorf("Failed to gather image %s info %s", oldImage, err)
		return false
	}
	if oldInfo.Config != nil {
		oldVersion = oldInfo.Config.Labels["com.docker.ucp.version"]
	}

	newInfo, err := c.client.InspectImage(newImage)
	if err != nil {
		log.Errorf("Failed to gather image %s info %s", newImage, err)
		return false
	}
	if newInfo.Config != nil {
		upgradesFrom = newInfo.Config.Labels["com.docker.ucp.upgrades_from"]
		newVersion = newInfo.Config.Labels["com.docker.ucp.version"]
	}

	if upgradesFrom == "" {
		log.Debugf("New image %s can upgrade from any version", newImage)
		return true
	}

	// special case this so we can report epoch/unknown as a compatible version...
	if oldVersion == "" {
		oldVersion = "0.0.0"
	}

	// TODO We may need to refine this approach for version ranges, wildcards, etc. but for now we'll
	//      use static version lists.
	compatibleVersions := strings.Split(upgradesFrom, ",")

	oldVer, err := version.NewVersion(strings.Split(oldVersion, " ")[0])
	if err != nil {
		log.Errorf("Failed to parse the old version string: %s - %s", oldVersion, err)
		return false
	}

	for _, compat := range compatibleVersions {
		compatVer, err := version.NewVersion(compat)
		if err != nil {
			log.Errorf("Failed to parse the compatible version string: %s - %s", compat, err)
			return false
		}
		if compatVer.Equal(oldVer) {
			log.Debugf("The new version %s is upgrade compatible from your current version (%s) (%s)", newVersion, oldVer, newImage)
			return true
		}
	}

	// No upgrade for you!
	log.Errorf("The new version %s is unable to upgrade from your current version %v (%s)", newVersion, oldVer, newImage)
	return false
}
