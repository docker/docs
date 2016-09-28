// Package images implements the image verification command
package images

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

// Run the images flow
func images(c *cli.Context) (int, error) {
	// If we're just listing, do that first to eliminate log spew
	if c.Bool("list") {
		orcaconfig.ImageVersion = c.String("image-version")
		// Use a map to eliminate duplicates from the output
		imageMap := make(map[string]interface{})
		for name := range orcaconfig.Images {
			imageName, err := orcaconfig.GetContainerImage(name)
			if err != nil {
				log.Errorf("Unexpected error: %s", err)
			} else {
				imageMap[imageName] = struct{}{}
			}
		}
		for imageName := range imageMap {
			fmt.Println(imageName)
		}
		return 0, nil
	}

	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}

	if err := ec.VerifyOrPullImages(c.Bool("interactive")); err != nil {
		log.Error("We were unable to pull one or more required images.  Please set REGISTRY_USERNAME, REGISTRY_PASSWORD, and REGISTRY_EMAIL environment variables for your Docker Hub account on this container with -e flags to run.")
		return 1, err
	}
	return 0, nil
}

// Run the images flow
func Run(c *cli.Context) {
	if code, err := images(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
