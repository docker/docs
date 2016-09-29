package proxy

import (
	"errors"
	"log"
	"strings"

	"github.com/Microsoft/go-winio"
	"github.com/docker/engine-api/client"
	"github.com/rneugeba/virtsock/go/hvsock"
)

// BackendDialer knows how to connect to the Docker backend
// either directly or with a DockerClient.
type BackendDialer interface {
	Dial() (ReadWriteWriteCloser, error)

	DockerClient() (*client.Client, error)
}

// NewBackendDialer returns a new BackendDialer or an error if the provided network
// is not supported.
func NewBackendDialer(network string, underlyingPath string) (BackendDialer, error) {
	switch network {
	case "npipe":
		return &npipeBackendDialer{
			pipeName: underlyingPath,
		}, nil

	case "hvsock":
		vmid, svcid, err := parseHVsockPath(underlyingPath)
		if err != nil {
			return nil, err
		}

		return &hvsockBackendDialer{
			addr: hvsock.HypervAddr{
				VmId:      vmid,
				ServiceId: svcid,
			},
		}, nil
	}
	return nil, errors.New("Unsupported Protocol " + network)
}

type hvsockBackendDialer struct {
	addr hvsock.HypervAddr
}

func (d *hvsockBackendDialer) Dial() (ReadWriteWriteCloser, error) {
	log.Printf("Dial Hyper-V socket %s", d.addr.String())

	con, err := hvsock.Dial(d.addr)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully dialed Hyper-V socket %s", d.addr.String())
	return con, nil
}

func (d *hvsockBackendDialer) DockerClient() (*client.Client, error) {
	return nil, errors.New("DockerClient() not compatible with hvsocks")
}

func parseHVsockPath(underlyingPath string) (hvsock.GUID, hvsock.GUID, error) {
	ids := strings.Split(underlyingPath, ":")
	if len(ids) != 2 {
		return hvsock.GUID_ZERO, hvsock.GUID_ZERO, errors.New("Malformed path" + underlyingPath)
	}

	vmid, err := hvsock.GuidFromString(ids[0])
	if err != nil {
		return hvsock.GUID_ZERO, hvsock.GUID_ZERO, err
	}
	svcid, err := hvsock.GuidFromString(ids[1])
	if err != nil {
		return hvsock.GUID_ZERO, hvsock.GUID_ZERO, err
	}

	return vmid, svcid, nil
}

type npipeBackendDialer struct {
	pipeName string
}

func (d *npipeBackendDialer) Dial() (ReadWriteWriteCloser, error) {
	log.Println("Dial name pipe ", d.pipeName)
	conn, err := winio.DialPipe(d.pipeName, nil)
	if err != nil {
		return nil, err
	}

	dialer, _ := conn.(ReadWriteWriteCloser)
	log.Printf("Successfully dialed name pipe %s", d.pipeName)
	return dialer, nil
}

func (d *npipeBackendDialer) DockerClient() (*client.Client, error) {
	return nil, errors.New("DockerClient() not compatible with npipe")
}
