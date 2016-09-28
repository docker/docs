package utils

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/client"
	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	kvstore "github.com/docker/libkv/store"
)

const (
	// DefaultHostname is the default built-in hostname
	DefaultHostname = "docker.io"
	// LegacyDefaultHostname is automatically converted to DefaultHostname
	LegacyDefaultHostname = "index.docker.io"
	// DefaultRepoPrefix is the prefix used for default repositories in default host
	DefaultRepoPrefix = "library/"
)

func FromUnixTimestamp(timestamp int64) (*time.Time, error) {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		return nil, err
	}

	t := time.Unix(i, 0)
	return &t, nil
}

func GetTLSConfig(caCert, cert, key []byte, allowInsecure bool) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}
	if allowInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	return &tlsConfig, nil
}

func GetClient(dockerUrl string, caPem, certPem, keyPem []byte, allowInsecure bool) (*dockerclient.Client, *http.Transport, error) {
	// only load env vars if no args
	// check environment for docker client config
	envDockerHost := os.Getenv("DOCKER_HOST")
	if dockerUrl == "" && envDockerHost != "" {
		dockerUrl = envDockerHost
	}

	// load tlsconfig
	var tlsConfig *tls.Config
	if len(caPem) != 0 && len(certPem) != 0 && len(keyPem) != 0 {
		log.Debug("using tls for communication with docker")

		cfg, err := GetTLSConfig(caPem, certPem, keyPem, allowInsecure)
		if err != nil {
			log.Fatalf("error configuring tls: %s", err)
		}
		tlsConfig = cfg
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
	}
	version := "" // Always take the latest swarm API

	client, err := dockerclient.NewClient(dockerUrl, version, httpClient, nil)
	if err != nil {
		return nil, nil, err
	}

	return client, transport, nil
}

func GenerateID(n int) string {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[:n]
}

// Preserve ordering, but de-dup, clean up white space, etc.
func JoinCerts(certs ...string) string {
	rawJoin := strings.Join(certs, "\n")
	outList := []string{}
	certMap := map[string]interface{}{}
	for _, certRaw := range strings.SplitAfter(rawJoin, "-----END CERTIFICATE-----") {
		cert := strings.TrimSpace(certRaw)
		if len(cert) == 0 {
			continue
		}
		if _, ok := certMap[cert]; !ok {
			certMap[cert] = struct{}{}
			outList = append(outList, cert)
		}
	}
	return strings.Join(outList, "\n") + "\n"
}

// MaybeWrapEtcdClusterErr is used to add extra detail to etcd client erros.
// This is extremely useful because the default error string contains no detail
// on what the actual error is and only says:
//     etcd cluster is unavailable or misconfigured
func MaybeWrapEtcdClusterErr(err error) error {
	if err == nil {
		return nil
	}

	clusterErr, ok := err.(*etcdclient.ClusterError)
	if !ok {
		// Return the error unaltered.
		return err
	}

	return fmt.Errorf("%s: %s", clusterErr.Error(), clusterErr.Detail())
}

// ErrTimedOut is used to indicate that the function passed to RunWithTimeout
// has timed out.
var ErrTimedOut = errors.New("timed out")

// RunWithTimeout runs the given function and returns either it's result or a
// timeout error after the given timeout duration.
func RunWithTimeout(f func() error, timeout time.Duration) (err error) {
	errC := make(chan error, 1)

	go func() {
		errC <- f()
	}()

	select {
	case err = <-errC:
	case <-time.After(timeout):
		return ErrTimedOut
	}

	return err
}

// splitHostname is copied straight from the docker/docker/reference package
// however, we cannot add docker/docker as a dependency
// splitHostname splits a repository name to hostname and remotename string.
// If no valid hostname is found, the default hostname is used. Repository name
// needs to be already validated before.
func splitHostname(name string) (hostname, remoteName string) {
	i := strings.IndexRune(name, '/')
	if i == -1 || (!strings.ContainsAny(name[:i], ".:") && name[:i] != "localhost") {
		hostname, remoteName = DefaultHostname, name
	} else {
		hostname, remoteName = name[:i], name[i+1:]
	}
	if hostname == LegacyDefaultHostname {
		hostname = DefaultHostname
	}
	if hostname == DefaultHostname && !strings.ContainsRune(remoteName, '/') {
		remoteName = DefaultRepoPrefix + remoteName
	}
	return
}

func IsHub(hostname string) bool {
	return hostname == DefaultHostname
}

func GetRepoFullName(named string) string {
	hostname, remoteName := splitHostname(named)
	return hostname + "/" + remoteName
}

// GetRepoFromFormVars parses the repository information from the
// http request
func GetRepoFromFormVars(req *http.Request) (string, error) {
	// Fill in the form variables for the request so we can extract the image name,
	// but also get rid of the body so that we're not parsing it (it could be huge).

	newReq := *req
	newReq.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	err := newReq.ParseForm()
	if err != nil {
		return "", err
	}
	return newReq.Form.Get("fromImage"), nil
}

func GetKVLock(kv kvstore.Store, p string, opts *kvstore.LockOptions, chDone chan bool) error {
	lock, err := kv.NewLock(p, opts)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		chLost, err := lock.Lock(nil)
		if err != nil {
			if err == kvstore.ErrKeyNotFound {
				time.Sleep(time.Duration(100) * time.Millisecond)
				continue
			}
		}

		// Unlock if we don't lose the lock
		go func() {
			select {
			case <-chDone:
				log.Debugf("unlocking: path=%s", p)
				lock.Unlock()
			case <-chLost:
				log.Warn("kv lock lost during access list update; probably took too long")
			}
		}()

		return nil
	}

	return fmt.Errorf("unable to get lock: path=%s", p)
}

// GetHostAddressFromContainerJSON extracts the current node's host address from the
// --discovery argument of a container
func GetHostAddressFromContainerJSON(containerJSON types.ContainerJSON, sentinel string) (string, error) {
	for i := 0; i < len(containerJSON.Config.Cmd); i++ {
		flag := containerJSON.Config.Cmd[i]
		if strings.Contains(flag, sentinel) {
			log.Debug("extracting original host-address")
			discoveryURLString := ""
			if strings.Contains(flag, "=") {
				splitFlag := strings.SplitN(flag, "=", 2)
				discoveryURLString = splitFlag[1]
			} else {
				discoveryURLString = containerJSON.Config.Cmd[i+1]
			}
			discoveryURL, err := url.Parse(discoveryURLString)
			if err != nil {
				return "", fmt.Errorf("Failed to extract existing host-address from the running ucp-controller container: %s", err)
			}
			if discoveryURL.Scheme == "" {
				// Not a URL, just parse the arg directly
				address, _, _ := net.SplitHostPort(discoveryURLString)
				return address, nil
			} else {
				address, _, _ := net.SplitHostPort(discoveryURL.Host)
				return address, nil
			}
		}
	}
	return "", fmt.Errorf("the %s flag was not found", sentinel)
}

// GetHostPortsFromContainerJSON returns a list of all the host ports bound to a container
func GetHostPortsFromContainerJSON(containerJSON types.ContainerJSON) ([]string, error) {
	// Get the container name and strip the "/" prefix, if any
	containerName := strings.TrimPrefix(containerJSON.Name, "/")

	exposedPorts := []string{}
	if len(containerJSON.HostConfig.PortBindings) != 1 {
		return exposedPorts, fmt.Errorf("Container %s has %d port bindings, instead of 1",
			containerName, len(containerJSON.HostConfig.PortBindings))
	}

	for port, binding := range containerJSON.HostConfig.PortBindings {
		if len(binding) != 1 {
			return exposedPorts, fmt.Errorf("Container %s has a port binding for port %s with %d entries",
				containerName, port, len(binding))
		}
		exposedPorts = append(exposedPorts, binding[0].HostPort)
	}
	return exposedPorts, nil
}

// equalSlices returns true if two slices of N unique elements actually contain
// the same set of elements. This is done in O(N^2) with no extra space, as we espect
// the slices to have very few elements
func equalSlices(first []string, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for _, n := range first {
		found := false
		for _, m := range second {
			if n == m {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
