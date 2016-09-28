package utils

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
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
	// Check to see if the target already has the matching image
	log.Infof("Transfering %s to %s", imageName, destination.URL.Host)
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

	// Silly that dockerclient doesn't have a first class API for this...
	uri := fmt.Sprintf("/%s/images/get?names=%s", dockerclient.APIVersion, imageName)
	req, err := http.NewRequest("GET", source.URL.String()+uri, nil)
	if err != nil {
		return fmt.Errorf("Failed to located source image %s: %s", imageName, err)
	}
	//req.Header.Add("Content-Type", "application/json")
	resp, err := source.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to load image %sfrom source: %s", imageName, err)
	}

	if resp.StatusCode != 200 {
		payload, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("Failed to load image from source: HTTP error %d: %s", resp.StatusCode, string(payload))
	}

	// XXX Should we wrap a buffering reader?

	// Now stream the body to the other side
	if err := destination.LoadImage(resp.Body); err != nil {
		return fmt.Errorf("Failed to save image to destination: %s: %s", imageName, err)
	}
	resp.Body.Close()
	return nil
}

func PullImages(destination *dockerclient.DockerClient, imageNames []string) error {
	for _, imageName := range imageNames {
		if _, err := destination.InspectImage(imageName); err != nil {
			log.Infof("Pulling %s with user %s", imageName, os.Getenv("REGISTRY_USERNAME"))
			if err := destination.PullImage(imageName, &dockerclient.AuthConfig{
				Username: os.Getenv("REGISTRY_USERNAME"),
				Password: os.Getenv("REGISTRY_PASSWORD"),
				Email:    os.Getenv("REGISTRY_EMAIL"),
			}); err != nil {
				return err
			}
		} else {
			log.Infof("Skipping existing image %s", imageName)
		}
	}
	return nil

}

// Build a simple little "hello world" busybox container to generate some load
func BuildImage(c *dockerclient.DockerClient, repo string) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// TODO - might want to consider making this a little toy web server
	//        so we can use it to exercise other aspects of the system...
	dockerFile := `FROM busybox:latest
LABEL orcatest=1
RUN echo "$(date) hello" > /hello.txt
`
	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(dockerFile)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(dockerFile)); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	// Now build it
	buildCfg := &dockerclient.BuildImage{
		DockerfileName: "Dockerfile",
		Context:        buf,
		RepoName:       repo,
		Remove:         true,
		NoCache:        true, // Make sure we get fresh image layers to stress it some more
	}

	out, err := c.BuildImage(buildCfg)
	if err != nil {
		log.Fatalln("Build failed: ", err)
	}

	data := make([]byte, 1024)
	for {
		n, err := out.Read(data)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}
		// Might want to suppress this to reduce chatter... but we need to pay attention to errors
		fmt.Print(string(data[:n]))
	}
	out.Close()
	return nil
}
