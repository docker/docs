package framework

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/sshclient"
	"github.com/samalba/dockerclient"
)

type IntegrationFramework struct {
	API    apiclient.APIClient
	SSH    sshclient.SSHClient
	Docker *dockerclient.DockerClient
	Config TestConfig
}

type TestConfig struct {
	DTRHost           string
	EnziHost          string
	DaemonHost        string
	AuthMethodLDAP    bool
	CleanDataDump     bool
	AdminUsername     string
	AdminPassword     string
	RetryAttempts     int
	SSHUsername       string
	SSHPort           int
	SSHPassword       string
	SSHPrivateKeyPath string
	DindHost          string
	DockerAPIVersion  string
	// used for pushing
	AdminAuthConfig *dockerclient.AuthConfig
}

func (t TestConfig) AdminNamespace() string {
	return t.DTRHost + "/" + t.AdminUsername
}

func BuildFramework(t *testing.T) *IntegrationFramework {
	host := os.Getenv("DTR_HOST")
	if host == "" {
		host = "172.17.42.1"
	}

	SSHPrivateKeyPath := os.Getenv("SSH_PRIVATE_KEY_PATH")
	SSHPassword := os.Getenv("SSH_PASSWORD")
	SSHPort, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	if err != nil {
		SSHPort = 22
	}

	SSHUsername := os.Getenv("SSH_USERNAME")
	if SSHUsername == "" {
		SSHUsername = "ubuntu"
	}

	daemonHost := os.Getenv("DAEMON_HOST")
	if daemonHost == "" {
		daemonHost = "unix:///var/run/docker.sock"
	}

	retryAttempts, err := strconv.Atoi(os.Getenv("RETRY_ATTEMPTS"))
	if err != nil {
		t.Fatal("Please provide a valid integer for retry attempts")
	}

	authMethodLDAP := os.Getenv("AUTH_METHOD") == "ldap"

	cleanDataDump := os.Getenv("CLEAN_DATA_DUMP") != "false"

	adminUsername := os.Getenv("ENZI_ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("ENZI_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "password"
	}
	dindHost := os.Getenv("DIND_HOST")

	config := TestConfig{
		DTRHost:        host,
		AuthMethodLDAP: authMethodLDAP,
		CleanDataDump:  cleanDataDump,
		AdminUsername:  adminUsername,
		AdminPassword:  adminPassword,
		AdminAuthConfig: &dockerclient.AuthConfig{
			Username: adminUsername,
			Password: adminPassword,
			Email:    "a@a.a", // Classic
		},
		SSHUsername:       SSHUsername,
		SSHPassword:       SSHPassword,
		SSHPrivateKeyPath: SSHPrivateKeyPath,
		SSHPort:           SSHPort,
		RetryAttempts:     retryAttempts,
		DaemonHost:        daemonHost,
		DindHost:          dindHost,
		DockerAPIVersion:  "v1.22",
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			MaxIdleConnsPerHost:   5,
		},
	}

	api := apiclient.New(config.DTRHost, config.RetryAttempts, httpClient)

	ssh := sshclient.New(strings.Split(config.DTRHost, ":")[0], config.SSHUsername,
		config.SSHPassword, config.SSHPrivateKeyPath, config.SSHPort)

	docker, err := dockerclient.NewDockerClient(config.DaemonHost, nil)
	if err != nil {
		t.Fatal("Failed to construct docker client", err)
	}

	return &IntegrationFramework{api, ssh, docker, config}
}
