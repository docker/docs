package proxy

import (
	"log"
	"net"
	"sync"

	"github.com/docker/engine-api/client"
)

// BackendDialer knows how to connect to the Docker backend
// either directly or with a DockerClient.
type BackendDialer interface {
	Dial() (ReadWriteWriteCloser, error)

	DockerClient() (*client.Client, error)
}

// NewBackendDialer returns a new BackendDialer.
func NewBackendDialer(network, underlyingPath string) (BackendDialer, error) {
	return &socketBackendDialer{
		network:        network,
		underlyingPath: underlyingPath,
	}, nil
}

type socketBackendDialer struct {
	network        string
	underlyingPath string

	sync.Mutex
	docker *client.Client
}

func (d *socketBackendDialer) Dial() (ReadWriteWriteCloser, error) {
	log.Println("Dial Socket", d.underlyingPath)
	conn, err := net.Dial(d.network, d.underlyingPath)
	if err != nil {
		return nil, err
	}

	dialer, _ := conn.(ReadWriteWriteCloser)
	return dialer, nil
}

func (d *socketBackendDialer) DockerClient() (*client.Client, error) {
	d.Lock()
	defer d.Unlock()

	if d.docker == nil {
		host := d.network + "://" + d.underlyingPath

		docker, err := client.NewClient(host, "", nil, nil)
		if err != nil {
			log.Printf("Failed to create client: %s", err)
			return nil, err
		}

		d.docker = docker
	}

	return d.docker, nil
}
