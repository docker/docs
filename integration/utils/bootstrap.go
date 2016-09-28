package utils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/samalba/dockerclient"

	// Partial conversion to engine-api
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/strslice"
	"golang.org/x/net/context"
)

// Process the base name with a prefix and TAG
func BuildImageString(name string) string {
	return orcaImagesMap[name]
}

func BuildBootstrapImage() string {
	prefix := os.Getenv("ORCA_ORG")
	if prefix == "" {
		prefix = "docker"
	}
	tag := os.Getenv("TAG")
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", prefix, "ucp", tag)
}

// Figure out if we need a non-standard image version for the bootstrapper
func buildImageVersionArgs() []string {
	prefix := os.Getenv("ORCA_ORG")
	tag := os.Getenv("TAG")
	if strings.Contains(prefix, "dev") {
		if tag == "latest" {
			// Drop the latest tag since we want the bootstrapper to figure it out
			tag = ""
		}
		// Always let the bootstrapper pin to its version for dev images
		// Note: if we try to do upgrades from dockerorcadev based images, this assumption will break
		return []string{"--image-version", "dev:"}
	} else if tag != "" && tag != "latest" {
		return []string{"--image-version", tag}
	}
	// Else rely on the bootstrapper to pull the right images
	return []string{}
}

var (
	BootstrapImage = BuildBootstrapImage()
	BootstrapName  = "ucp"

	// TODO - we really need a centralized place for all these images instead of scattered all over the place

	// Replicated from ucp-bootstrap so we can avoid importing it
	// TODO - support dockerorcadev perhaps?
	staticOrcaImages = []string{
		BootstrapImage,
		"busybox:latest",
	}
	// populated by LoadAllLocalOrcaImages
	orcaImages       = []string{}
	orcaImagesMap    = map[string]string{}
	orcaImagesLock   = sync.Mutex{}
	ImageVersionArgs = buildImageVersionArgs()
	OrcaImageTag     = ""
)

// Return a list of images to either pull or transfer
func getAllOrcaImages(client *dockerclient.DockerClient) error {
	args := []string{"images", "--list"}
	args = append(args, ImageVersionArgs...)

	if os.Getenv("PULL_IMAGES") != "" {
		// Make sure we actually have the bootstrapper... (this is a no-op if it exists)
		err := PullImages(client, []string{BootstrapImage})
		if err != nil {
			return err
		}
	} else {
		localClient, err := GetClientFromEnv()
		if err != nil {
			return err
		}
		if err := TransferImages(localClient, client, []string{BootstrapImage}); err != nil {
			return err
		}
	}
	log.Infof("Running bootstraper to find versions with args: %s", args)

	// Prevent collisions during image list
	oldBootstrapName := BootstrapName
	BootstrapName = ""
	defer func() { BootstrapName = oldBootstrapName }()

	imageList, stderrOutput, err := RunBootstrapperWithIO(client, args, []string{}, []string{}, nil)
	if err != nil {
		log.Error(stderrOutput)
		log.Error(err)
		return err
	}
	orcaImages = append(orcaImages, staticOrcaImages...)
	orcaImages = append(orcaImages, strings.Split(strings.TrimSpace(imageList), "\n")...)

	for _, imageName := range orcaImages {

		imageParts := strings.Split(imageName, ":")
		key := imageParts[0]
		tag := imageParts[1]
		if tag != "latest" && OrcaImageTag == "" {
			OrcaImageTag = tag
			log.Debugf("setting orcaimagetag to %s", OrcaImageTag)
		}
		if strings.Contains(key, "/") {
			key = strings.Split(key, "/")[1]
		}
		orcaImagesMap[key] = imageName
	}

	log.Infof("Orca Images and friends: %s", orcaImages)
	return nil
}

// Load all the local orca images
func LoadAllLocalOrcaImages(client *dockerclient.DockerClient) error {
	orcaImagesLock.Lock()
	if len(orcaImages) == 0 {
		err := getAllOrcaImages(client)
		if err != nil {
			orcaImagesLock.Unlock()
			return err
		}
	}
	orcaImagesLock.Unlock()

	if os.Getenv("PULL_IMAGES") != "" {
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
		if os.Getenv("PULL_IMAGES") != "" {
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

// Run the bootstrapper with the following args, output will be displayed to the log
func RunBootstrapper(dclient *dockerclient.DockerClient, args, env, extraBinds []string) error {

	// We want to use the new engine-api so we have better control of I/O
	client, err := ConvertToEngineAPI(dclient)
	if err != nil {
		return fmt.Errorf("Failed to convert clients: %s", err)
	}

	// This is a little magical - maybe we should have a different approach?
	var finalArgs []string
	if args[0] == "install" || args[0] == "upgrade" || args[0] == "join" || args[0] == "regen-certs" || args[0] == "uninstall" || args[0] == "engine-discovery" {
		hasImage := false
		for _, arg := range args {
			if arg == "--image-version" {
				hasImage = true
				break
			}
		}
		if !hasImage {
			finalArgs = append(args, ImageVersionArgs...)
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}

	cfg := &container.Config{
		Image:        BootstrapImage,
		Cmd:          strslice.StrSlice(finalArgs),
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    true,
		Env:          appendRegistryEnv(env),
	}
	hostConfig := &container.HostConfig{
		Binds: append(extraBinds, "/var/run/docker.sock:/var/run/docker.sock"),
	}

	resp, err := client.ContainerCreate(context.TODO(), cfg, hostConfig, nil, BootstrapName)
	if err != nil {
		return fmt.Errorf("Failed to create bootstrap container: %s", err)
	}
	containerId := resp.ID
	defer client.ContainerRemove(context.TODO(), containerId, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})

	attachResp, err := client.ContainerAttach(context.TODO(), containerId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return fmt.Errorf("Failed to attach: %s", err)
	}
	rd := bufio.NewReader(attachResp.Reader)

	log.Debugf("Launching %s with args:%v env:%v", BootstrapImage, finalArgs, env)
	if err := client.ContainerStart(context.TODO(), containerId, types.ContainerStartOptions{
		CheckpointID: "",
	}); err != nil {
		return fmt.Errorf("Failed to start bootstrap container: %s", err)
	}
	timeout := 5 * time.Second
	defer client.ContainerStop(context.TODO(), containerId, &timeout)

	doneChan := make(chan int)
	go func() {
		res, err := client.ContainerWait(context.TODO(), containerId)
		if err != nil {
			log.Errorf("Failed to wait for phase 2 container: %s", err)
		}
		doneChan <- res
	}()

	// Wire this up for input
	/*
	   // Send all remaining input to the second phase
	   go func() {
	           if _, err := io.Copy(attachResp.Conn, os.Stdin); err != nil {
	                   log.Warnf("Stdin copy interrupted: %s", err)
	           }
	   }()
	*/

	// If we supported tty, we'd do this differently...

	// stdCopy is really chatty in debug mode
	go func() {
		oldLevel := log.GetLevel()
		log.SetLevel(log.InfoLevel)
		defer log.SetLevel(oldLevel)
		if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, rd); err != nil {
			log.Errorf("Failed to read logs from bootstrap container: %s", err)
			return
		}
	}()

	exitCode := <-doneChan

	if exitCode != 0 {
		return fmt.Errorf("Container exited with %d", exitCode)
	}

	return nil
}

// Check the env for registry flags, append if not already in the list
func appendRegistryEnv(env []string) []string {
	res := []string{}
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	email := os.Getenv("REGISTRY_EMAIL")

	appendUsername := username != ""
	appendPassword := password != ""
	appendEmail := email != ""

	for _, item := range env {
		res = append(res, item)
		if strings.HasPrefix(item, "REGISTRY_USERNAME") {
			appendUsername = false
		} else if strings.HasPrefix(item, "REGISTRY_PASSWORD") {
			appendPassword = false
		} else if strings.HasPrefix(item, "REGISTRY_EMAIL") {
			appendEmail = false
		}
	}

	if appendUsername {
		res = append(res, fmt.Sprintf("REGISTRY_USERNAME=%s", username))
	}
	if appendPassword {
		res = append(res, fmt.Sprintf("REGISTRY_PASSWORD=%s", password))
	}
	if appendEmail {
		res = append(res, fmt.Sprintf("REGISTRY_EMAIL=%s", email))
	}
	return res
}

// Run the bootstrapper with the following args, output will be returned as a string
func RunBootstrapperWithIO(dclient *dockerclient.DockerClient, args, env, extraBinds []string, input io.Reader) (string, string, error) {
	// We want to use the new engine-api so we have better control of I/O
	client, err := ConvertToEngineAPI(dclient)
	if err != nil {
		return "", "", fmt.Errorf("Failed to convert clients: %s", err)
	}

	// This is a little magical - maybe we should have a different approach?
	var finalArgs []string
	if args[0] == "install" || args[0] == "upgrade" || args[0] == "join" || args[0] == "regen-certs" || args[0] == "uninstall" || args[0] == "engine-discovery" {
		hasImage := false
		for _, arg := range args {
			if arg == "--image-version" {
				hasImage = true
				break
			}
		}
		if !hasImage {
			finalArgs = append(args, ImageVersionArgs...)
		} else {
			finalArgs = args
		}
	} else {
		finalArgs = args
	}
	cfg := &container.Config{
		Image:        BootstrapImage,
		Cmd:          finalArgs,
		OpenStdin:    input != nil,
		AttachStdin:  input != nil,
		AttachStdout: true,
		AttachStderr: true,
		StdinOnce:    input != nil,
		Env:          appendRegistryEnv(env),
	}
	hostConfig := &container.HostConfig{
		Binds: append(extraBinds, "/var/run/docker.sock:/var/run/docker.sock"),
	}
	resp, err := client.ContainerCreate(context.TODO(), cfg, hostConfig, nil, BootstrapName)
	if err != nil {
		return "", "", fmt.Errorf("Failed to create bootstrap container %s: %s", BootstrapImage, err)
	}
	containerId := resp.ID
	defer client.ContainerRemove(context.TODO(), containerId, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})

	attachResp, err := client.ContainerAttach(context.TODO(), containerId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  input != nil,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("Failed to attach: %s", err)
	}
	rd := bufio.NewReader(attachResp.Reader)

	log.Debugf("Launching %s with args:%v env:%v", BootstrapImage, finalArgs, env)

	if err := client.ContainerStart(context.TODO(), containerId, types.ContainerStartOptions{
		CheckpointID: "",
	}); err != nil {
		return "", "", fmt.Errorf("Failed to start bootstrap container: %s", err)
	}
	timeout := 5 * time.Second
	defer client.ContainerStop(context.TODO(), containerId, &timeout)

	doneChan := make(chan int)
	go func() {
		res, err := client.ContainerWait(context.TODO(), containerId)
		if err != nil {
			log.Errorf("Failed to wait for phase 2 container: %s", err)
		}
		doneChan <- res
	}()

	// Wire this up for input
	if input != nil {
		// Send all remaining input to the second phase
		go func() {
			if _, err := io.Copy(attachResp.Conn, input); err != nil {
				log.Warnf("input copy interrupted: %s", err)
			}
		}()
	}

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	go func() {
		oldLevel := log.GetLevel()
		log.SetLevel(log.InfoLevel)
		defer log.SetLevel(oldLevel)
		if _, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, rd); err != nil {
			log.Errorf("Failed to read output from bootstrap container: %s", err)
			return
		}
	}()

	exitCode := <-doneChan

	if exitCode != 0 {
		return "", "", fmt.Errorf("Container exited with %d", exitCode)
	}

	return stdoutBuffer.String(), stderrBuffer.String(), nil
}

func ValidateOrcaControllerCerts(machine Machine) error {
	tarBytes, err := machine.TarHostDir("/etc/docker/ssl/orca")
	if err != nil {
		return err
	}
	paths := make(map[string]bool)
	tr := tar.NewReader(bytes.NewReader(tarBytes))
	failed := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Bad tar file returned from the machine: %s", err)
		}
		paths[hdr.Name] = true
		// Could verify contents looks plausible, but we'll call it good enough with the filenames
	}
	for _, requiredPath := range []string{
		// Not quite an exhastive list, but good enough
		"./config.json",
		"./ucp_ca_chain.pem",
		// "./orca_controller.pem", // TODO Not merged upstream yet
		// "./orca_controller_key.pem",
		"./swarm_ca.pem",
		"./ucp_ca.pem",
		"./ucp_ca_key.pem",
		"./swarm_ca_key.pem",
	} {
		if !paths[requiredPath] {
			log.Errorf("Missing: %s", requiredPath)
			failed = true
		}
	}
	if failed {
		return errors.New("Failed to find one or more required cert files")
	} else {
		return nil
	}
}

func PingOrcaServer(machine Machine) (int, error) {
	ip, err := machine.GetIP()
	if err != nil {
		return 0, fmt.Errorf("Failed to lookup machines IP: %s", err)
	}
	uri := fmt.Sprintf("https://%s:443/_ping", ip)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(uri)
	if err != nil {
		log.Debug("Failed to connect to orca at %s - %s", uri, err)
		return 0, err
	} else {
		return resp.StatusCode, nil
	}
}

func ValidateOrcaServerRunning(machine Machine, retryCount int) error {
	ip, err := machine.GetIP()
	if err != nil {
		return fmt.Errorf("Failed to lookup machines IP: %s", err)
	}
	uri := fmt.Sprintf("https://%s:443/_ping", ip)
	expected := 200

	// For this scenario, we're not going to bother validating the certificate
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	// Retry a bit incase things are coming up
	var lastError error
	for i := 0; i < retryCount; i += 1 {
		resp, err := client.Get(uri)
		if err != nil {
			log.Infof("Failed to connect to orca at %s - %s (retrying...)", uri, err)
			lastError = err
		} else if resp.StatusCode != expected {
			body, _ := ioutil.ReadAll(resp.Body)
			lastError = fmt.Errorf("Unexpected status code: %d - Payload %s", resp.StatusCode, body)
		} else {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return lastError
}

// Try to create a container in a loop to detect when the cluster recovers (at least one healthy node required)
func ValidateClusterHealthy(serverURL string, retryCount int) error {
	log.Infof("Verifying cluster is healthy: %s", serverURL)
	var adminClient *dockerclient.DockerClient
	var err error
	for i := 0; i < 60; i++ {
		if adminClient == nil {
			adminClient, err = GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
			if err != nil {
				log.Debugf("Not ready yet: %s", err)
				time.Sleep(2 * time.Second)
				continue
			}
		}
		_, err = CreateContainers(adminClient, 1, false)
		if err == nil {
			log.Debug("Cluster recovered")
			return nil
		}
		log.Debugf("Not ready yet: %s", err)
		time.Sleep(2 * time.Second)
	}
	return err
}

func GetOrcaID(client *dockerclient.DockerClient) (string, error) {
	containers, err := client.ListContainers(true, false, "")
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		if id, found := container.Labels["com.docker.ucp.InstanceID"]; found {
			// Just assume a single instance...
			return id, nil
		}
	}
	return "", fmt.Errorf("Failed to locate an Orca instance on this engine")
}

// GetOrcaAddr returns the concatenation of the Host-Address and Orca-Port for
// the given machine.
func GetOrcaAddr(machine Machine) (string, error) {
	ip, err := machine.GetIP()
	if err != nil {
		return "", err
	}

	// TODO: make the port configurable.
	return net.JoinHostPort(ip, "443"), nil
}

// get the public URL - should be accessible from the host running the test
func GetOrcaURL(machine Machine) (string, error) {
	addr, err := GetOrcaAddr(machine)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s", addr), nil
}

func GetOrcaURLs(s *OrcaTestSuite) ([]string, error) {
	addrs := []string{}
	for _, m := range s.ControllerMachines {
		addr, err := GetOrcaAddr(m)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, fmt.Sprintf("https://%s", addr))
	}
	return addrs, nil
}

// Get the internal orca URL - should be accessible from other nodes in the same network
func GetOrcaInternalURL(machine Machine) (string, error) {
	ip, err := machine.GetInternalIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s:443", ip), nil
}

func GetOrcaFingerprint(machine Machine) (string, error) {
	data, err := machine.CatHostFile(ServerCertHostPath)
	if err != nil {
		return "", err
	}
	return GetFingerprint(data)
}

func WaitForEngineRecovery(client *dockerclient.DockerClient, maxAttempts int) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		version, err := client.Version()
		if err == nil {
			log.Infof("Engine is back: %s", version)
			return nil
		}
		log.Debugf("Waiting for engine, got error %s", err)
		time.Sleep(time.Second)
	}
	return err
}

// Try to find specific command arguments for a running container on the system
func GetContainerFlagValue(client *dockerclient.DockerClient, container string, argName string) (string, error) {
	info, err := client.InspectContainer(container)
	if err != nil {
		return "", err
	}

	for i, arg := range info.Args {
		if arg == argName {
			return info.Args[i+1], nil
		}
	}
	return "", fmt.Errorf("Argument %s not found for container %s", container, argName)
}
