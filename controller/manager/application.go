package manager

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
)

var (
	ErrApplicationDoesNotExist = errors.New("application does not exist")
)

type ComposeConfig struct {
	project        string
	configHash     string
	composeVersion string
	serviceName    string
}

func getComposeConfig(labels map[string]string) (*ComposeConfig, error) {
	project := labels["com.docker.compose.project"]
	configHash := labels["com.docker.compose.config-hash"]
	composeVersion := labels["com.docker.compose.version"]
	serviceName := labels["com.docker.compose.service"]

	return &ComposeConfig{
		project:        project,
		configHash:     configHash,
		composeVersion: composeVersion,
		serviceName:    serviceName,
	}, nil
}

func (m DefaultManager) getApplication(c *types.Container, containers []types.Container) (*orca.Application, error) {
	cfg, err := getComposeConfig(c.Labels)
	if err != nil {
		return nil, err
	}

	app := &orca.Application{
		Name:           cfg.project,
		Id:             cfg.configHash,
		ConfigHash:     cfg.configHash,
		ComposeVersion: cfg.composeVersion,
	}

	services, err := m.getAppServices(app, containers)
	if err != nil {
		return nil, err
	}

	app.Services = services

	return app, nil
}

func (m DefaultManager) getServiceContainers(name, serviceName string, containers []types.Container) ([]types.Container, error) {
	existing := make(map[string]bool)
	var serviceContainers []types.Container
	for _, cnt := range containers {

		cfg, err := getComposeConfig(cnt.Labels)
		if err != nil {
			return nil, err
		}

		if cfg.project != name {
			continue
		}

		if cfg.serviceName == serviceName {
			if _, exists := existing[cnt.ID]; !exists {
				log.Debugf("adding %s for service %s", cnt.ID, serviceName)
				existing[cnt.ID] = true
				serviceContainers = append(serviceContainers, cnt)
			}
		}
	}

	return serviceContainers, nil
}

func (m DefaultManager) getAppServices(app *orca.Application, containers []types.Container) ([]*orca.Service, error) {
	svcs := []*orca.Service{}
	existingServices := make(map[string]bool)

	for _, cnt := range containers {

		cfg, err := getComposeConfig(cnt.Labels)
		if err != nil {
			return nil, err
		}

		if cfg.project != app.Name {
			continue
		}

		if _, ok := existingServices[cfg.serviceName]; ok {
			continue
		}

		existingServices[cfg.serviceName] = true

		s := &orca.Service{
			Name: cfg.serviceName,
		}

		svcs = append(svcs, s)
	}

	for _, s := range svcs {
		svcContainers, err := m.getServiceContainers(app.Name, s.Name, containers)
		if err != nil {
			return nil, err
		}

		s.Containers = svcContainers
	}

	return svcs, nil
}

func (m DefaultManager) Applications(ctx *auth.Context) ([]*orca.Application, error) {
	// load compose containers
	containers, err := m.ListUserContainers(ctx, true, false, "")
	if err != nil {
		return nil, err
	}

	applications := []*orca.Application{}
	existingApps := make(map[string]bool)

	for _, cnt := range containers {
		// check for compose label
		if _, ok := cnt.Labels["com.docker.compose.project"]; !ok {
			continue
		}

		app, err := m.getApplication(&cnt, containers)
		if err != nil {
			log.Warn("error getting application config")
			continue
		}

		if _, exists := existingApps[app.Name]; !exists {
			existingApps[app.Name] = true
			applications = append(applications, app)
		}
	}

	return applications, nil
}

func (m DefaultManager) Application(ctx *auth.Context, project string) (*orca.Application, error) {
	var application *orca.Application
	apps, err := m.Applications(ctx)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.Name == project {
			application = app
			break
		}
	}
	if application == nil {
		return nil, ErrApplicationDoesNotExist
	}

	return application, nil
}

func (m DefaultManager) RemoveApplication(application *orca.Application) error {
	if application == nil {
		return ErrApplicationDoesNotExist
	}
	for _, s := range application.Services {
		for _, c := range s.Containers {
			if err := m.client.ContainerRemove(context.TODO(), c.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m DefaultManager) ApplicationContainers(application *orca.Application) ([]types.Container, error) {
	containers := []types.Container{}
	if application == nil {
		return containers, ErrApplicationDoesNotExist
	}
	for _, svc := range application.Services {
		containers = append(containers, svc.Containers...)
	}

	return containers, nil
}

func (m DefaultManager) RestartApplication(application *orca.Application) error {
	if application == nil {
		return ErrApplicationDoesNotExist
	}
	for _, s := range application.Services {
		for _, c := range s.Containers {
			timeout := 5 * time.Second
			if err := m.client.ContainerRestart(context.TODO(), c.ID, &timeout); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m DefaultManager) StopApplication(application *orca.Application) error {
	if application == nil {
		return ErrApplicationDoesNotExist
	}

	for _, s := range application.Services {
		for _, c := range s.Containers {
			timeout := 5 * time.Second
			if err := m.client.ContainerStop(context.TODO(), c.ID, &timeout); err != nil {
				return err
			}
		}
	}
	return nil
}
