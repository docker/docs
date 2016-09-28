package ha_utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	. "github.com/onsi/gomega"
	"github.com/samalba/dockerclient"
)

var (
	// Image name of the UCP bootstrapper
	BootstrapName = "ucp"
	LogFile       = io.Writer(os.Stdout)

	// Replicated from ucp-bootstrap so we can avoid importing it
	// populated by LoadAllLocalOrcaImages
	orcaImages     = []string{}
	orcaImagesLock = sync.Mutex{}
)

// Process the base name with a prefix and TAG
func BuildImageString(name string) string {
	prefix := os.Getenv("UCP_REPO")
	if prefix == "" {
		prefix = "docker"
	}
	tag := os.Getenv("UCP_TAG")
	if tag == "" {
		tag = "latest"
	}

	// Omit tag for images that are pinned
	if strings.Contains(name, ":") {
		return fmt.Sprintf("%s/%s", prefix, name)
	} else {
		return fmt.Sprintf("%s/%s:%s", prefix, name, tag)
	}
}

// Figure out if we need a non-standard image version for the bootstrapper
func buildImageVersionArgs() []string {
	prefix := os.Getenv("UCP_REPO")
	if prefix == "" {
		prefix = "docker"
	}
	tag := os.Getenv("UCP_TAG")
	if tag == "" {
		tag = "latest"
	}
	if strings.Contains(prefix, "dev") {
		return []string{"--image-version", fmt.Sprintf("dev:%s", tag)}
	} else {
		return []string{"--image-version", tag}
	}
}

func GetBootstrapImage() string {
	return BuildImageString("ucp")
}

func GetStaticOrcaImages() []string {
	return []string{
		GetBootstrapImage(),
		"busybox:latest",
	}
}

// Return a list of images to either pull or transfer
func getAllOrcaImages() error {
	localClient, err := GetClientFromEnv()
	if err != nil {
		return err
	}
	args := []string{"images", "--list"}
	args = append(args, buildImageVersionArgs()...)

	pullImages := os.Getenv("UCP_PULL_IMAGES")
	// If the UCP bootstrapper is already present
	if pullImages != "" {
		// Make sure we actually have the bootstrapper... (this is a no-op if it exists)
		err = PullImages(localClient, []string{GetBootstrapImage()})
		if err != nil {
			return err
		}
	}
	log.Infof("Running bootstraper to find versions with args: %s", args)

	// Prevent collisions during image list
	oldBootstrapName := BootstrapName
	BootstrapName = ""
	defer func() { BootstrapName = oldBootstrapName }()

	imageList, err := RunUCPBootstrapperWithOutput(localClient, args, []string{})
	if err != nil {
		return err
	}
	orcaImages = append(orcaImages, GetStaticOrcaImages()...)
	orcaImages = append(orcaImages, strings.Split(strings.TrimSpace(imageList), "\n")...)

	log.Infof("Orca Images and friends: %s", orcaImages)
	return nil
}

// Load all the local orca images
func LoadAllLocalOrcaImages(client *dockerclient.DockerClient) error {
	orcaImagesLock.Lock()
	if len(orcaImages) == 0 {
		err := getAllOrcaImages()
		if err != nil {
			orcaImagesLock.Unlock()
			return err
		}
	}
	orcaImagesLock.Unlock()

	if os.Getenv("UCP_PULL_IMAGES") != "" {
		if err := PullImages(client, orcaImages); err != nil {
			return err
		}
	} else {
		localClient, err := GetClientFromEnv()
		if err != nil {
			return err
		}
		if err := TransferImages(localClient, client, orcaImages); err != nil {
			return err
		}
	}

	// Allow hijacking swarm to test specific images
	swarmImage := os.Getenv("SWARM_IMAGE")
	if swarmImage != "" {
		log.Infof("Using alternative swarm image %s", swarmImage)
		if os.Getenv("UCP_PULL_IMAGES") != "" {
			if err := PullImages(client, []string{swarmImage}); err != nil {
				return err
			}
		} else {
			localClient, err := GetClientFromEnv()
			if err != nil {
				return err
			}
			if err := TransferImages(localClient, client, []string{swarmImage}); err != nil {
				return err
			}
		}
		imageTag := strings.Split(BuildImageString("ucp-swarm"), ":")
		log.Infof("Tagging %s -> %v", swarmImage, imageTag)

		return client.TagImage(swarmImage, imageTag[0], imageTag[1], true)
	}
	return nil
}

// Run the UCP bootstrapper with the following args, output will be displayed to the log
func RunUCPBootstrapper(client *dockerclient.DockerClient, args, env []string) error {
	// This is a little magical - maybe we should have a different approach?
	var finalArgs []string
	if args[0] == "install" || args[0] == "upgrade" || args[0] == "join" || args[0] == "regen-certs" {
		hasImage := false
		for _, arg := range args {
			if arg == "--image-version" {
				hasImage = true
				break
			}
		}
		if !hasImage {
			finalArgs = append(args, buildImageVersionArgs()...)
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}

	hostConfig := dockerclient.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}
	cfg := &dockerclient.ContainerConfig{
		Image: GetBootstrapImage(),
		Cmd:   finalArgs,
		Env: append(env,
			[]string{ // Add in the registry account if set so it can pull if it needs to
				fmt.Sprintf("REGISTRY_USERNAME=%s", os.Getenv("REGISTRY_USERNAME")),
				fmt.Sprintf("REGISTRY_PASSWORD=%s", os.Getenv("REGISTRY_PASSWORD")),
				fmt.Sprintf("REGISTRY_EMAIL=%s", os.Getenv("REGISTRY_EMAIL")),
			}...),
		HostConfig: hostConfig,
	}
	containerId, err := client.CreateContainer(cfg, BootstrapName, nil)
	if err != nil {
		return fmt.Errorf("Failed to create bootstrap container from %s: %s; cfg: %v", GetBootstrapImage(), err, cfg)
	}
	defer client.RemoveContainer(containerId, true, true)

	log.Debugf("Launching %s with args:%v env:%v", GetBootstrapImage(), finalArgs, env)
	if err := client.StartContainer(containerId, &hostConfig); err != nil {
		return fmt.Errorf("Failed to start bootstrap container: %s", err)
	}
	defer client.StopContainer(containerId, 5)

	reader, err := client.ContainerLogs(containerId, &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	})

	if _, err = stdcopy.StdCopy(LogFile, LogFile, reader); err != nil {
		return fmt.Errorf("Failed to read logs from bootstrap container: %s", err)
	}

	info, err := client.InspectContainer(containerId)
	if err != nil {
		return fmt.Errorf("Failed to inspect container after completion: %s", err)
	}
	if info.State == nil {
		return fmt.Errorf("Container didn't finish!")
	}

	if info.State.ExitCode != 0 {
		return fmt.Errorf("Container exited with %d", info.State.ExitCode)
	}

	return nil
}

// Run the UCP bootstrapper with the following args, output will be returned as a string
func RunUCPBootstrapperWithOutput(client *dockerclient.DockerClient, args, env []string) (string, error) {
	// This is a little magical - maybe we should have a different approach?
	var finalArgs []string
	if args[0] == "install" || args[0] == "upgrade" || args[0] == "join" || args[0] == "regen-certs" {
		hasImage := false
		for _, arg := range args {
			if arg == "--image-version" {
				hasImage = true
				break
			}
		}
		if !hasImage {
			finalArgs = append(args, buildImageVersionArgs()...)
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}
	hostConfig := dockerclient.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}

	cfg := &dockerclient.ContainerConfig{
		Image:      GetBootstrapImage(),
		Cmd:        finalArgs,
		Env:        env,
		HostConfig: hostConfig,
	}
	containerId, err := client.CreateContainer(cfg, BootstrapName, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create bootstrap container %s: %s", GetBootstrapImage(), err)
	}
	defer client.RemoveContainer(containerId, true, true)

	log.Debugf("Launching %s with args:%v env:%v", GetBootstrapImage(), finalArgs, env)
	if err := client.StartContainer(containerId, &hostConfig); err != nil {
		return "", fmt.Errorf("Failed to start bootstrap container: %s", err)
	}
	defer client.StopContainer(containerId, 5)

	reader, err := client.ContainerLogs(containerId, &dockerclient.LogOptions{
		Follow: true,
		Stdout: true,
		Stderr: true,
	})

	buffer := new(bytes.Buffer)

	if _, err = stdcopy.StdCopy(buffer, buffer, reader); err != nil {
		log.Debug("cannot read logs from logs reader")
		return "", err
	}

	info, err := client.InspectContainer(containerId)
	if err != nil {
		return buffer.String(), fmt.Errorf("Failed to inspect container after completion: %s", err)
	}
	if info.State == nil {
		return buffer.String(), fmt.Errorf("Container didn't finish!")
	}

	if info.State.ExitCode != 0 {
		return buffer.String(), fmt.Errorf("Container exited with %d", info.State.ExitCode)
	}

	return buffer.String(), nil
}

func RunUCPBootstrapperWithInput(client *dc.Client, args, env []string, input string) (string, string, error) {
	// This is a little magical - maybe we should have a different approach?
	var finalArgs []string
	if args[0] == "install" || args[0] == "upgrade" || args[0] == "join" || args[0] == "regen-certs" {
		hasImage := false
		for _, arg := range args {
			if arg == "--image-version" {
				hasImage = true
				break
			}
		}
		if !hasImage {
			finalArgs = append(args, buildImageVersionArgs()...)
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}

	log.Debugf("Launching %s with args:%v env:%v", GetBootstrapImage(), finalArgs, env)
	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		StdinOnce:    true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Image:        GetBootstrapImage(),
		Cmd:          finalArgs,
		Env: append(env,
			[]string{ // Add in the registry account if set so it can pull if it needs to
				fmt.Sprintf("REGISTRY_USERNAME=%s", os.Getenv("REGISTRY_USERNAME")),
				fmt.Sprintf("REGISTRY_PASSWORD=%s", os.Getenv("REGISTRY_PASSWORD")),
				fmt.Sprintf("REGISTRY_EMAIL=%s", os.Getenv("REGISTRY_EMAIL")),
			}...),
	}, &container.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}, &network.NetworkingConfig{}, BootstrapName)
	Expect(err).To(BeNil())

	containerId := resp.ID
	keep := false
	defer func() {
		if !keep {
			err := client.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				log.Debugf("failed to remove ucp container %s", err)
			}
		} else {
			log.Debug("not deleting ucp container so it can be debugged")
		}
	}()

	err = client.ContainerStart(context.Background(), containerId, types.ContainerStartOptions{})
	Expect(err).To(BeNil())

	response, err := client.ContainerAttach(context.Background(), containerId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	Expect(err).To(BeNil())

	// it's important to start reading output before stream in the input
	errBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	go func() {
		errDebug := io.MultiWriter(errBuf, LogFile)
		outDebug := io.MultiWriter(outBuf, LogFile)
		_, err = stdcopy.StdCopy(outDebug, errDebug, response.Reader)
		if err != nil {
			log.Errorf("Failed to stream logs: %s\n", err)
		}
		log.Info("finished reading output")
	}()

	inBuf := strings.NewReader(input)
	go func() {
		log.Debug("about to iocopy the input")
		debugBuf := io.MultiWriter(response.Conn, LogFile)
		if _, err := io.Copy(debugBuf, inBuf); err != nil {
			log.Warnf("Stdin copy interruped: %s", err)
		}
		log.Debug("iocopied the input")
	}()

	log.Debug("container wait")
	rc, err := client.ContainerWait(context.Background(), containerId)
	if err != nil {
		log.Errorf("Failed to wait for phase 2 container: %s", err)
	}
	if rc != 0 {
		keep = true
		return "", "", fmt.Errorf("non-zero return code: %d; output: %s stderr: %s", rc, outBuf.String(), errBuf.String())
	}

	return outBuf.String(), errBuf.String(), nil
}

// get the public URL - should be accessible from the host running the test
func GetOrcaURL(machine Machine) (string, error) {
	ip, err := machine.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s:444", ip), nil
}

func GetOrcaFingerprint(machine Machine) (string, error) {
	data, err := machine.CatHostFile(ServerCertHostPath)
	if err != nil {
		return "", err
	}
	return GetFingerprint(data)
}
