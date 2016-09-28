package ha_utils

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

// Copied from github.com/docker/orca/auth

type AuthToken struct {
	ID        string `json:"token_id,omitempty"`
	Token     string `json:"auth_token,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// Copied from github.com/docker/orca/controller/api
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func GetOrcaToken(client *http.Client, url, username, password string) (string, error) {
	// XXX Consider making this do tofu...
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	orcaURL, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}

	orcaURL.Path = "/auth/login"
	creds := Credentials{
		Username: username,
		Password: password,
	}

	reqJson, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	resp, err := client.Post(orcaURL.String(), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode == 200 {
		var authToken AuthToken
		if err := json.Unmarshal(body, &authToken); err != nil {
			return "", err
		} else {
			return authToken.Token, nil
		}
	} else {
		return "", fmt.Errorf("Unexpected error logging in to Orca: %s", string(body))
	}

}

func GetTokenHeader(token string) (string, string) {
	return "Authorization", fmt.Sprintf("Bearer %s", token)
}

func GetUCPBundle(serverURL, username, password string) (*zip.Reader, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Login and get a token for the user
	token, err := GetOrcaToken(client, serverURL, username, password)
	if err != nil {
		return nil, err
	}
	log.Debugf("Token was: %s", token)

	orcaURL.Path = "/api/clientbundle"

	// Now download the bundle
	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		if err == nil {
			return nil, fmt.Errorf("Failed to get client bundle (%d): %s", resp.StatusCode, string(body))
		} else {
			return nil, err
		}
	}
	return zip.NewReader(bytes.NewReader(body), int64(len(body)))
}

func GetUCPUserTLSConfig(serverURL, username, password string) (*tls.Config, error) {
	caData := ""
	certData := ""
	keyData := ""

	// Get a client bundle
	zr, err := GetUCPBundle(serverURL, username, password)
	if err != nil {
		return nil, err
	}
	// Decode bundle and extract ca/cert/key
	for _, zf := range zr.File {

		src, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		var data []byte
		if strings.HasSuffix(zf.Name, "ca.pem") {
			data, err = ioutil.ReadAll(src)
			caData = string(data)
		} else if strings.HasSuffix(zf.Name, "cert.pem") {
			data, err = ioutil.ReadAll(src)
			certData = string(data)
		} else if strings.HasSuffix(zf.Name, "key.pem") {
			data, err = ioutil.ReadAll(src)
			keyData = string(data)
		}
		if err != nil {
			return nil, err
		}
	}
	if caData == "" || certData == "" || keyData == "" {
		return nil, fmt.Errorf("Failed to decode zip bundle from Orca")
	}

	// Set up the tls cert for connection, and dump out in case we need to troubleshoot later
	// Uncomment for troubleshooting, otherwise this kinda spews in the logs
	//log.Debugf("Loading CA/CERT/KEY: \n%s\n%s\n%s", caData, certData, keyData)
	cert, err := tls.X509KeyPair([]byte(certData), []byte(keyData))
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caData))
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return tlsConfig, nil
}

// Given a orca connection URL, and cred's, return a dockerclient for that user
func GetUCPUserDockerClient(serverURL, username, password string) (*dockerclient.DockerClient, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	// Obtain TLS config
	tlsConfig, err := GetUCPUserTLSConfig(serverURL, username, password)
	if err != nil {
		return nil, err
	}

	// And connect to the docker API
	orcaURL.Scheme = "tcp"
	orcaURL.Path = ""
	log.Debugf("creating client to : %s", orcaURL.String())
	return dockerclient.NewDockerClient(orcaURL.String(), tlsConfig)
}
