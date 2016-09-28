package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

// Build lots of images to stress the system
// count specifies how many images to build
// nameFormat is an optional format string for naming (default "testimg%d:latest")
// keepTrying if false, aborts on first failure
// returns total created, and the last known error if anything went wrong
func BuildImages(client *dockerclient.DockerClient, count int, nameFormat string, keepTrying bool) (int, error) {
	if nameFormat == "" {
		nameFormat = "testimg%d:latest"
	}
	var retError error
	total := 0

	log.Debugf("Building %d images with %s", count, nameFormat)
	for i := 0; i < count; i++ {
		err := BuildImage(client, fmt.Sprintf(nameFormat, i))
		if err != nil {
			retError = err
			if !keepTrying {
				break
			}
		} else {
			total += 1
		}
	}
	log.Debugf("Completed run after %d images", total)
	return total, retError
}

// Create lots of users to stress the system
// count specifies how many images to build
// nameFormat is an optional format string for naming (default "newuser%d")
// passwordFormat is an optional format string for naming (default "secret%d")
// roles is an optional roles (deafult readwrite)
// keepTrying if false, aborts on first failure
// returns total created, and the last known error if anything went wrong
func AddUsers(client *http.Client, serverURL string, count int, nameFormat, passwordFormat, adminUser, adminPassword string, isAdmin bool, keepTrying bool) (int, error) {
	if nameFormat == "" {
		nameFormat = "newuser%d"
	}
	if passwordFormat == "" {
		passwordFormat = "secret%d"
	}
	log.Debugf("Building %d images with %s", count, nameFormat)
	var retError error
	total := 0

	for i := 0; i < count; i++ {
		err := CreateNewUser(client, serverURL, adminUser,
			adminPassword, fmt.Sprintf(nameFormat, i), fmt.Sprintf(passwordFormat, i), isAdmin, auth.RestrictedControl)
		if err != nil {
			retError = err
			if !keepTrying {
				break
			}
		} else {
			total += 1
		}
	}
	return total, retError
}

// Create a bunch of containers, using busybox, or images created by the BuildImage(s) routines
func CreateContainers(client *dockerclient.DockerClient, count int, keepTrying bool) (int, error) {

	// List all the images available, and just loop through that set for the base image
	log.Infof("Listing images")
	availableImages, err := client.ListImages(false)
	if err != nil {
		return 0, err
	}
	// Filter the list down to our test images, or use busybox
	imageList := []string{"busybox"}
	for _, image := range availableImages {
		if image.Labels["orcatest"] != "" {
			imageList = append(imageList, image.Id)
		}
	}

	log.Infof("Creating containers based on %d applicable images", len(imageList))
	var retError error
	total := 0

	for i := 0; i < count; i++ {
		cfg := &dockerclient.ContainerConfig{
			Image: imageList[i%len(imageList)],
			Cmd:   []string{"sleep", "96h"}, // Do nothing for a long time so we can see "running" containers
		}
		containerId, err := client.CreateContainer(cfg, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "cannot specify 64-byte hexadecimal strings") {
				log.Warnf("Unable to create image based on id, falling back to busybox")
				cfg = &dockerclient.ContainerConfig{
					Image: "busybox",
					Cmd:   []string{"sleep", "96h"}, // Do nothing for a long time so we can see "running" containers
				}
				containerId, err = client.CreateContainer(cfg, "", nil)
			} else {
				log.Warnf("Failed to start container %d - %s", i, err)

				retError = err
				if keepTrying {
					continue
				}
				break
			}
		}
		err = client.StartContainer(containerId, nil)
		if err != nil {
			retError = err
			if keepTrying {
				continue
			}
			break
		}
		log.Debugf("Created container and started %d - %s", i, containerId)
		total += 1
	}
	return total, retError
}

// The following routines are test routines that can be used in different contexts

func TestAddUsers(t *testing.T, serverURL string, UserCount int) {
	log.Debugf("Creating %d users", UserCount)

	// Reuse a client so it can recycle connections
	client := &http.Client{
		Transport: &http.Transport{
			// Sloppy for testing only - don't copy this into production code!
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
	}

	_, err := AddUsers(client, serverURL, UserCount, "", "", GetAdminUser(), GetAdminPassword(), false, true)
	require.Nil(t, err)

	log.Debug("Spot check access with a few of them")
	for i := 0; i < UserCount; i += 100 {
		newUsername := fmt.Sprintf("newuser%d", i)
		newPassword := fmt.Sprintf("secret%d", i)
		log.Debug("getting docker client as new user")
		client, err := GetUserDockerClient(serverURL, newUsername, newPassword)
		require.Nil(t, err)
		log.Debug("getting version")
		version, err := client.Version()
		require.Nil(t, err)
		require.True(t, strings.Contains(version.Version, "ucp"))
	}
}

func TestBuildImages(t *testing.T, serverURL string, ImageCount int) {
	client, err := GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
	require.Nil(t, err)

	// Keep trying, but fail test if we don't succeed completely
	_, err = BuildImages(client, ImageCount, "", true)
	require.Nil(t, err)

	log.Debug("Spot check access with a few of them")
	for i := 0; i < ImageCount; i += 100 {
		_, err := client.InspectImage(fmt.Sprintf("testimg%d:latest", i))
		require.Nil(t, err)
		// TODO - should we check anything on the inspect results?
	}
}

func TestCreateContainers(t *testing.T, serverURL string, ContainerCount int) {
	client, err := GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
	require.Nil(t, err)

	count, err := CreateContainers(client, ContainerCount, true)
	log.Infof("Created %d containers out of %d goal", count, ContainerCount)
	require.Nil(t, err)

	// Sanity test by listing containers at the end
	containerList, err := client.ListContainers(true, false, "")
	require.Nil(t, err)
	require.True(t, len(containerList) >= ContainerCount)
}

func TestComposeApps(t *testing.T, serverURL string, AppCount int) {
	log.Debugf("Creating %d apps", AppCount)
	for i := 0; i < AppCount; i++ {
		TestSimpleCompose(t, serverURL, GetAdminUser(), GetAdminPassword())
	}
}
