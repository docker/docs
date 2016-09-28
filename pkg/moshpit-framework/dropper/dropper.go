package dropper

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/dockerutil"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/flags"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/util"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/go-connections/nat"
)

type Dropper struct {
	ctx    context.Context
	log    logrus.FieldLogger
	config moshpit.Config
}

var serverInternalPort = 1337
var serverContainerName = "moshpit-server"
var clientContainerNamePrefix = "moshpit-client-"

func (d *Dropper) Execute() error {
	dc, err := d.login()
	if err != nil {
		return err
	}

	err = d.pullImages(dc)
	if err != nil {
		return err
	}

	d.log.WithField("node", d.config.Dropper.ServerNode).Info("starting server on ucp node")
	serverID, err := d.runServer(dc)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(d.config.Dropper.NumClients)

	errsLock := sync.Mutex{}
	errs := []error{}
	clientIDsLock := sync.Mutex{}
	clientIDs := []string{}
	clientIDsToNames := map[string]string{}
	resultsLock := sync.Mutex{}
	results := map[string]int{}
	for i := 0; i < d.config.Dropper.NumClients; i++ {
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("client%d", i)
			d.log.WithField("name", name).Info("starting client")
			clientID, err := d.runClient(dc, name)
			if err != nil {
				errsLock.Lock()
				errs = append(errs, err)
				errsLock.Unlock()
			}
			clientIDsLock.Lock()
			clientIDs = append(clientIDs, clientID)
			clientIDsToNames[clientID] = name
			clientIDsLock.Unlock()
			d.log.WithField("name", name).Info("started client")

			// TODO: set up a timeout after which we force remove it?
			res, err := dc.ContainerWait(d.ctx, clientID)
			if err != nil {
				errsLock.Lock()
				errs = append(errs, err)
				errsLock.Unlock()
			}
			resultsLock.Lock()
			results[clientID] = res
			resultsLock.Unlock()
		}(i)
	}

	d.log.Info("waiting for clients to finish")
	wg.Wait()
	if len(errs) > 0 {
		d.log.WithField("errors", errs).Error("errors running some of the clients")
	}

	if !d.config.Dropper.NoCleanup {
		d.log.Info("removing clients in parallel")
		rmWG := sync.WaitGroup{}
		rmWG.Add(d.config.Dropper.NumClients)
		errsLock := sync.Mutex{}
		errs := []error{}
		for _, clientID := range clientIDs {
			go func(clientID string) {
				defer rmWG.Done()
				clientInfo := map[string]interface{}{"clientID": clientID, "clientName": clientIDsToNames[clientID]}
				if results[clientID] != 0 {
					d.log.WithFields(clientInfo).Warn("not cleaning up failed client")
				}
				d.log.WithFields(clientInfo).Info("removing client")
				err = dockerutil.SafeContainerRemove(d.ctx, dc, clientID)
				if err != nil {
					errsLock.Lock()
					errs = append(errs, err)
					errsLock.Unlock()
				}
			}(clientID)
		}
		rmWG.Wait()
		if len(errs) > 0 {
			d.log.WithField("errors", errs).Error("errors removing some of the clients")
		}

		d.log.WithField("serverID", serverID).Info("removing server")
		err = dockerutil.SafeContainerRemove(d.ctx, dc, serverID)
		if err != nil {
			return err
		}
	}
	return nil
}

func New(ctx context.Context, config moshpit.Config) *Dropper {
	d := &Dropper{
		ctx:    ctx,
		log:    moshpit.LoggerFromCtx(ctx),
		config: config,
	}
	return d
}

func (d *Dropper) genServerConfig() ([]byte, error) {
	// we do this to make sure the server INSIDE the container is configured how we want it
	tmpDropper := d.config.Dropper
	tmpIP := d.config.Server.ListenIP
	tmpPort := d.config.Server.ListenPort
	d.config.Server.ListenIP = "0.0.0.0"
	d.config.Server.ListenPort = serverInternalPort
	d.config.Dropper = moshpit.DropperConfig{}
	configBytes, err := yaml.Marshal(d.config)
	if err != nil {
		return nil, err
	}
	d.config.Server.ListenIP = tmpIP
	d.config.Server.ListenPort = tmpPort
	d.config.Dropper = tmpDropper
	return configBytes, nil
}

func (d *Dropper) pullImages(client *client.Client) error {
	parts := strings.Split(d.config.Dropper.MoshpitImage, ":")
	imageID := parts[0]
	tag := parts[1]
	readCloser, err := client.ImagePull(d.ctx, types.ImagePullOptions{ImageID: imageID, Tag: tag, RegistryAuth: dockerutil.MakeRegistryAuth(d.config.Dropper.HubUsername, d.config.Dropper.HubPassword, d.config.Dropper.HubRefreshToken)}, nil)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	decoder := json.NewDecoder(readCloser)
	for err == nil {
		response := jsonmessage.JSONMessage{}
		err = decoder.Decode(&response)
		if response.Error == nil {
			d.log.Info(response.Status)
		} else {
			return fmt.Errorf("Daemon error: %s", response.Error)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("Failed to parse response: %s", err)
	}

	return nil
}

func (d *Dropper) runServer(client *client.Client) (string, error) {
	configBytes, err := d.genServerConfig()
	if err != nil {
		return "", err
	}

	envs := []string{}
	envs = append(envs, fmt.Sprintf("constraint:node==%s", d.config.Dropper.ServerNode))

	cmd := []string{}
	if d.config.Dropper.Debug {
		cmd = append(cmd, "--debug")
	}
	cmd = append(cmd,
		"server",
		"--"+flags.ConfigDataFlag.Name, string(configBytes),
	)

	pubPort := nat.Port(fmt.Sprint(d.config.Server.ListenPort))
	portBindings := make(nat.PortMap)
	portBindings[pubPort] = append(portBindings[pubPort], nat.PortBinding{HostIP: d.config.Server.ListenIP, HostPort: fmt.Sprint(serverInternalPort)})

	resp, err := client.ContainerCreate(d.ctx, &container.Config{
		Env:          envs,
		Image:        d.config.Dropper.MoshpitImage,
		Cmd:          cmd,
		ExposedPorts: map[nat.Port]struct{}{pubPort: {}},
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, serverContainerName)
	if err != nil {
		return "", fmt.Errorf("Couldn't create container '%s' from image '%s': %s", serverContainerName, d.config.Dropper.MoshpitImage, err)
	}
	containerID := resp.ID

	err = client.ContainerStart(d.ctx, containerID)
	if err != nil {
		return containerID, err
	}

	return containerID, nil
}

func (d *Dropper) runClient(client *client.Client, name string) (string, error) {
	cmd := []string{}
	if d.config.Dropper.Debug {
		cmd = append(cmd, "--debug")
	}
	cmd = append(cmd,
		"client",
		"--"+flags.ServerFlag.Name, fmt.Sprintf("%s:%d", d.config.Dropper.ServerHost, d.config.Server.ListenPort),
		"--"+flags.NameFlag.Name, name,
	)
	envs := []string{}
	if d.config.Dropper.Spread {
		envs = append(envs, "affinity:container!=moshpit*")
	}
	clientConstraints := strings.Split(d.config.Dropper.ClientConstraints, "|")
	if len(clientConstraints) == 1 && clientConstraints[0] == "" {
		clientConstraints = []string{}
	}
	envs = append(envs, clientConstraints...)

	cfg := &container.Config{
		Env:   envs,
		Image: d.config.Dropper.MoshpitImage,
		Cmd:   cmd,
	}

	containerName := clientContainerNamePrefix + name
	resp, err := client.ContainerCreate(d.ctx, cfg, &container.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("Couldn't create container '%s' from image '%s': %s", containerName, d.config.Dropper.MoshpitImage, err)
	}

	containerID := resp.ID
	err = client.ContainerStart(d.ctx, containerID)
	if err != nil {
		return containerID, err
	}

	return containerID, nil
}

func (d *Dropper) login() (*client.Client, error) {
	httpClient, err := util.HTTPClient(d.config.Dropper.UCPInsecureTLS, d.config.Dropper.UCPCA)
	if err != nil {
		return nil, err
	}

	jwt, err := getJWT(d.config.Dropper.UCPURL, httpClient, d.config.Dropper.Username, d.config.Dropper.Password)
	if err != nil {
		if strings.Contains(err.Error(), "x509: certificate signed by unknown authority") {
			return nil, fmt.Errorf(fmt.Sprintf(`Certificate validation for UCP failed. You can get the UCP CA from https://%s/ca and then use it in your config.`, d.config.Dropper.UCPURL))
		}
		return nil, err
	}

	client, err := dockerutil.MakeDockerClient(d.config.Dropper.UCPURL, jwt, httpClient, d.config.Dropper.HubUsername, d.config.Dropper.HubPassword, d.config.Dropper.HubRefreshToken)
	if err != nil {
		return nil, err
	}
	return client, nil
}
