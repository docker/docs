package ha_utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/samalba/dockerclient"
)

// Transfer multiple images from the source to the destination
func TransferImages(source, destination *dockerclient.DockerClient, imageNames []string) error {
	for _, imageName := range imageNames {
		if err := TransferImage(source, destination, imageName); err != nil {
			return err
		}
	}
	return nil
}

// Transfer an image from source to destination, being smart not to copy if it's alreaddy there
func TransferImage(source, destination *dockerclient.DockerClient, imageName string) error {
	log.Infof("Beginning to transfer %s", imageName)
	// Check to see if the target already has the matching image
	if sourceInfo, err := source.InspectImage(imageName); err != nil {
		return fmt.Errorf("Unable to locate the source image %s: %s", imageName, err)
	} else {
		if destInfo, err := destination.InspectImage(imageName); err == nil {
			if sourceInfo.Id == destInfo.Id && sourceInfo.Parent == destInfo.Parent {
				log.Infof("Image %s already synchronized", imageName)
				return nil
			} else {
				log.Debugf("Source %v - Dest %v", sourceInfo, destInfo)
			}
		}
	}
	log.Infof("Beginning to transfer %s", imageName)

	// Silly that dockerclient doesn't have a first class API for this...
	uri := fmt.Sprintf("/%s/images/get?names=%s", dockerclient.APIVersion, imageName)
	req, err := http.NewRequest("GET", source.URL.String()+uri, nil)
	if err != nil {
		return fmt.Errorf("Failed to located source image %s: %s", imageName, err)
	}
	//req.Header.Add("Content-Type", "application/json")
	resp, err := source.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to load image from source: %s: %s", imageName, err)
	}

	if resp.StatusCode != 200 {
		payload, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("Failed to load image from source: %d: %s", resp.StatusCode, string(payload))
	}

	// XXX Should we wrap a buffering reader?

	log.Infof("Transfering %s to %s", imageName, destination.URL.Host)
	// Now stream the body to the other side
	if err := destination.LoadImage(resp.Body); err != nil {
		return fmt.Errorf("Failed to save image to destination: %s: %s", imageName, err)
	}
	resp.Body.Close()
	return nil
}

func PullImages(destination *dockerclient.DockerClient, imageNames []string) error {
	// Following go pipeline guidelines: https://blog.golang.org/pipelines
	doneChan := make(chan struct{})
	for _, imageName := range imageNames {
		go func(imageName string) {
			defer GinkgoRecover()
			if _, err := destination.InspectImage(imageName); err != nil {
				log.Infof("Pulling %s with user %s", imageName, os.Getenv("REGISTRY_USERNAME"))
				err := destination.PullImage(imageName, &dockerclient.AuthConfig{
					Username: os.Getenv("REGISTRY_USERNAME"),
					Password: os.Getenv("REGISTRY_PASSWORD"),
					Email:    os.Getenv("REGISTRY_EMAIL"),
				})
				doneChan <- struct{}{}
				Expect(err).To(BeNil())
			} else {
				log.Infof("Skipping existing image %s", imageName)
				doneChan <- struct{}{}
			}
		}(imageName)
	}

	for i := 0; i < len(imageNames); i++ {
		<-doneChan
	}
	return nil
}
