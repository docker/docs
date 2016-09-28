package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
)

func GenerateImageArgs(dtrRepo, dtrTag string) (DTRImageArgs, error) {
	dtrBootstrapImage := ""
	if dtrRepo == "docker" {
		dtrBootstrapImage = fmt.Sprintf("%s/%s:%s", dtrRepo, "dtr", dtrTag)
	} else {
		dtrBootstrapImage = fmt.Sprintf("%s/%s:%s", dtrRepo, "dtr-dev", dtrTag)
	}

	dtrImages := []string{dtrBootstrapImage}

	localClient, err := ha_utils.GetClientFromEnv()
	if err != nil {
		return DTRImageArgs{}, err
	}

	pullImages := os.Getenv("DTR_PULL_IMAGES")
	// If the DTR bootstrapper is already present
	if pullImages != "" {
		// Make sure we actually have the bootstrapper... (this is a no-op if it exists)
		err = ha_utils.PullImages(localClient, []string{dtrBootstrapImage})
		if err != nil {
			return DTRImageArgs{}, err
		}
	}

	imageList, err := RunContainerWithOutput(localClient, dtrBootstrapImage, []string{"images"}, []string{})
	if err != nil {
		return DTRImageArgs{}, err
	}

	dtrImages = append(dtrImages, strings.Split(strings.TrimSpace(imageList), "\n")...)

	return DTRImageArgs{
		DTRRepo:           dtrRepo,
		DTRTag:            dtrTag,
		DTRBootstrapImage: dtrBootstrapImage,
		DTRImages:         dtrImages,
		SkipImages:        os.Getenv("DTR_SKIP_IMAGES") != "",
		PullImages:        os.Getenv("DTR_PULL_IMAGES") != "",
	}, nil
}
