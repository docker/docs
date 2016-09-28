package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	dockerclient "github.com/docker/engine-api/client"
	"golang.org/x/net/context"

	"github.com/docker/orca/bootstrap/config"
)

var (
	ErrUnableToConnect = errors.New("Unable to connect to system")
)

func LoadCerts(mount string) (*tls.Config, error) {
	caFilename := filepath.Join(mount, "ca.pem")
	certFilename := filepath.Join(mount, "cert.pem")
	keyFilename := filepath.Join(mount, "key.pem")
	// Load up the certs
	cert, err := tls.LoadX509KeyPair(certFilename, keyFilename)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(caFilename)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return tlsConfig, nil

}

// Wait for a node to appear on a node (assumed swarm or orca)
func WaitForNewNode(url, mount, newNode string, timeout time.Duration) error {
	log.Debugf("Checking for new node %s on %s", newNode, url)

	tlsConfig, err := LoadCerts(mount)
	if err != nil {
		log.Debugf("Failed to load certs: %s", err)
		return err
	}
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	version := ""

	for start := time.Now(); time.Since(start) < timeout; {
		docker, err := dockerclient.NewClient(url, version, client, nil)
		if err == nil {
			info, err := docker.Info(context.TODO())
			if err == nil {
				// look for the node in the output
				for _, label := range info.Labels {
					if strings.Contains(label, newNode) {
						log.Debugf("Found the node")
						return nil
					}
				}
			}
		}
		// TODO - For some classes of failures, we might want to fail fast instead of timing out
		time.Sleep(500 * time.Millisecond)
	}
	return ErrUnableToConnect
}

// Wait for a remote docker endpoint to come up (always assumes swarm node certs)
func WaitForEndpoint(url, mount string, timeout time.Duration) error {
	log.Debugf("Checking for liveness of %s", url)

	tlsConfig, err := LoadCerts(mount)
	if err != nil {
		log.Debugf("Failed to load certs: %s", err)
		return err
	}
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	version := ""

	// Verify we can talk to the endpoint
	for start := time.Now(); time.Since(start) < timeout; {
		docker, err := dockerclient.NewClient(url, version, client, nil)
		if err == nil {
			info, err := docker.Info(context.TODO())
			if err == nil {
				log.Debugf("Connected to %s %s %s", url, info.Name, info.ID)
				return nil
			}
		}
		// TODO - For some classes of failures, we might want to fail fast instead of timing out
		time.Sleep(500 * time.Millisecond)
	}
	return ErrUnableToConnect
}

// Wait for a UCP instance to come up
func WaitForOrca(url string, timeout time.Duration) error {
	// For the ping we don't bother trying to validate TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return waitForServer(fmt.Sprintf("%s/_ping", url), 200, client, timeout)

}

func WaitForCfssl(url string, timeout time.Duration) error {
	return waitForControllerTLSService(url, timeout)
}

func WaitForKv(url string, timeout time.Duration) error {
	return waitForControllerTLSService(url, timeout)
}

type KVHealthResp struct {
	Health string `json:"health,omitempty"`
}

func WaitForHealthyKv(client *http.Client, url string, timeout time.Duration) error {
	startTime := time.Now()
	for startTime.Add(timeout).After(time.Now()) {
		time.Sleep(2 * time.Second)

		healthy, err := IsEtcdHealthy(client, url)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				continue
			}
			return err
		}
		if healthy {
			return nil
		}
	}
	return fmt.Errorf("timed out while waiting for KV store to be healthy")
}

func IsEtcdHealthy(client *http.Client, url string) (bool, error) {
	resp, err := client.Get(url + "/health")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var healthResp KVHealthResp
	err = json.Unmarshal(body, &healthResp)
	if err != nil {
		return false, err
	}

	if healthResp.Health == "true" {
		return true, nil
	}
	return false, nil
}

func waitForControllerTLSService(url string, timeout time.Duration) error {
	// Use the phase 2 mounted Orca server client certs so mutual TLS works properly
	caFilename := filepath.Join(config.SwarmControllerCertVolumeMount, config.CAFilename)
	certFilename := filepath.Join(config.SwarmControllerCertVolumeMount, config.CertFilename)
	keyFilename := filepath.Join(config.SwarmControllerCertVolumeMount, config.KeyFilename)

	cert, err := tls.LoadX509KeyPair(certFilename, keyFilename)
	if err != nil {
		return err
	}
	caCert, err := ioutil.ReadFile(caFilename)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return waitForServer(url, 200, client, timeout)

}

func waitForServer(url string, expectedStatusCode int, client *http.Client, timeout time.Duration) error {
	log.Debugf("Checking for liveness of %s", url)
	// Verify we can talk to the endpoint
	for start := time.Now(); time.Since(start) < timeout; {
		resp, err := client.Get(url)
		if err == nil {
			// TODO - Need to revisit this if we get an unauth'd ping API
			if resp.StatusCode == expectedStatusCode {
				log.Debugf("Connected to %s", url)
				return nil
			}
			//log.Debugf("Status: %d", resp.StatusCode)
		} else {
			//log.Debugf("Connection error: %s", err)
		}
		// TODO - For some classes of failures, we might want to fail fast instead of timing out
		time.Sleep(500 * time.Millisecond)
	}
	return ErrUnableToConnect
}

// WaitForTLSConn waits for the TLS server at the given addr to be ready to
// accept connections. TLS client credentials are loaded from the given
// certMount directory. If we are unable to connect to the server withith the
// given timeout then ErrUnableToConnect is returned.
func WaitForTLSConn(addr, certMount string, timeout time.Duration) error {
	log.Debugf("Checking for liveness of TLS server at %s", addr)

	tlsConfig, err := LoadCerts(certMount)
	if err != nil {
		log.Debugf("Failed to load certs: %s", err)
		return err
	}

	// Wait no longer than 500 milliseconds for each attempt.
	dialer := &net.Dialer{Timeout: 500 * time.Millisecond}

	// Verify that we can complete a TLS handshake with the given addr.
	for start := time.Now(); time.Since(start) < timeout; {
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err == nil {
			// Success!
			conn.Close()
			return nil
		}

		// TODO - For some classes of failures, we might want to fail fast instead of timing out
		time.Sleep(500 * time.Millisecond)
	}

	log.Debugf("Failed to make TLS connection to %s: %s", addr, err)

	return ErrUnableToConnect
}

// WaitForAuthServer waits for an instance of the auth API or Worker server to
// be ready and responding to HTTP requests at the given URL and that the
// response has the given expectedStatuscode. if we are unable to connect to
// the server within the given timeout then ErrUnableToConnect is returned.
func WaitForAuthServer(url, certMount string, expectedStatusCode int, timeout time.Duration) error {
	log.Debugf("Checking for liveness of Auth API/Worker HTTP server at %s", url)

	tlsConfig, err := LoadCerts(certMount)
	if err != nil {
		log.Debugf("Failed to load certs: %s", err)
		return err
	}

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Verify we can talk to the endpoint
	for start := time.Now(); time.Since(start) < timeout; {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == expectedStatusCode {
				// Success!
				return nil
			}
		}

		// TODO - For some classes of failures, we might want to fail fast instead of timing out
		time.Sleep(500 * time.Millisecond)
	}

	return ErrUnableToConnect
}
