package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"golang.org/x/net/context"
	"time"
)

// Adder is an interface for adding containers
type Adder interface {
	// Add(containerID) is called asynchronously when a given container has started, and
	// needs resources setting up.
	Add(string) error
}

// Remover is an interface for removing containers
type Remover interface {
	// Remove(containerID) is called asynchronously when a given container has died, and
	// needs resources cleaning up.
	Remove(string) error
}

var docker *dockerclient.DockerClient
var ctx context.Context
var adders []Adder
var removers []Remover

func eventCallback(e *dockerclient.Event, ec chan error, args ...interface{}) {
	if e.Type == "container" && e.Action == "start" {
		for _, r := range adders {
			err := r.Add(e.Actor.ID)
			if err != nil {
				logrus.Printf("Failed to add resource on container start: %#v\n", err)
			}
		}
		return
	}
	if e.Type == "container" && e.Action == "die" {
		for _, r := range removers {
			err := r.Remove(e.Actor.ID)
			if err != nil {
				logrus.Printf("Failed to remove resource on container die: %#v\n", err)
			}
		}
		return
	}
}

func startWatchingEvents(onAdd []Adder, onRemove []Remover) {
	adders = onAdd
	removers = onRemove
	ctx = context.Background()

	waitForDockerUp()
	client, err := dockerclient.NewDockerClient("unix://"+getDockerLocalPath(), nil)
	if err != nil {
		logrus.Fatalf("Failed to create event stream docker client: %#v", err)
	}
	docker = client

	docker.StartMonitorEvents(eventCallback, nil)

	go func() {
		// FIXME: there's a race here between the event registration and listing
		// the containers. It's possible for events to be lost. Close the window a
		// little by giving the event thread 1s to register.
		time.Sleep(time.Second)
		containers, err := docker.ListContainers(false, false, "")
		if err != nil {
			logrus.Fatalf("Failed to connect client to event stream: %#v", err)
		}
		for _, container := range containers {
			logrus.Printf("startWatchingEvents existing (autostarted?) container %s\n", container.Id)
			for _, r := range onAdd {
				err := r.Add(container.Id)
				if err != nil {
					logrus.Printf("Failed to add resource on container start: %#v\n", err)
				}
			}
		}

	}()
}
