// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Microsoft/go-winio"
	"github.com/docker/pinata/v1/docker_proxy"
)

func init() {
	log.SetOutput(os.Stdout)
	proxy.SetVerboseLogging(false)
}

func main() {
	vm := flag.String("VM", "", "GUID of the Moby VM")
	proxyGUID := flag.String("proxyGUID", "23a432c2-537a-4291-bcb5-d62504644739", "Proxy GUID")
	daemonNamedPipe := flag.String("daemonNamedPipe", "", "Redirect api proxy output to the specifed named pipe")
	flag.Parse()

	var network string
	var underlyingPath string
	var envRewriter proxy.EnvRewriter
	var mountRewriter proxy.MountRewriter

	if *daemonNamedPipe != "" {
		network = "npipe"
		underlyingPath = *daemonNamedPipe
	} else {
		network = "hvsock"
		underlyingPath = fmt.Sprintf("%s:%s", *vm, *proxyGUID)

		envRewriter = proxy.NewHTTPProxyRewriter()
		mountRewriter = &proxy.WindowsMountRewriter{}
	}

	backendDialer, err := proxy.NewBackendDialer(network, underlyingPath)
	if err != nil {
		log.Fatal(err)
	}

	go startProxy(backendDialer, mountRewriter, envRewriter)
	go startProxyDeprecatedDefaultTCPPort(backendDialer, mountRewriter, envRewriter)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%s) received, stopping\n", s.String())
		}
	}
}

func startProxy(backendDialer proxy.BackendDialer, mountRewriter proxy.MountRewriter, envRewriter proxy.EnvRewriter) {
	log.Println("docker proxy: ready")

	listener, err := winio.ListenPipe(`//./pipe/docker_engine`, &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)", // Open to everyone
		MessageMode:        true,              // Use message mode so that CloseWrite() is supported
		InputBufferSize:    65536,             // Use 64KB buffers to improve performance
		OutputBufferSize:   65536,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := proxy.ServeBackend(listener, backendDialer, nil, mountRewriter, envRewriter); err != nil {
		log.Fatal(err)
	}
}

func startProxyDeprecatedDefaultTCPPort(backendDialer proxy.BackendDialer, mountRewriter proxy.MountRewriter, envRewriter proxy.EnvRewriter) {
	log.Println("docker proxy (on deprecated port): ready")

	if err := proxy.ServeWindows("localhost:2375", backendDialer, nil, mountRewriter, envRewriter); err != nil {
		log.Fatal(err)
	}
}
