package utils

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/types"
)

func SerializeReconcileArgs(expected *types.NodeConfig, current *types.NodeConfig) (string, error) {
	cfg := types.ReconcileConfig{
		Expected: expected,
		Current:  current,
	}
	s, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func DeserializeReconcileArgs(argsString string) (*types.NodeConfig, *types.NodeConfig, error) {
	var cfg types.ReconcileConfig
	err := json.Unmarshal([]byte(argsString), &cfg)
	if err != nil {
		return nil, nil, err
	}
	return cfg.Expected, cfg.Current, nil
}

// GetHTTPClient returns an *http.Client that uses the swarm-mode node cert as a client cert
func GetHTTPClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(orcaconfig.SwarmModeNodeCertPath, orcaconfig.SwarmModeNodeKeyPath)
	if err != nil {
		return nil, err
	}
	// TODO: controllerCA in caCertPool

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true, // TODO: Remove
			},
		},
	}
	return httpClient, nil
}

func AreServerCertsPresent() bool {
	if _, err := os.Stat(config.OrcaServerCertVolumeMount + "/ca.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.OrcaServerCertVolumeMount + "/cert.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.OrcaServerCertVolumeMount + "/key.pem"); os.IsNotExist(err) {
		return false
	}
	// TODO: validate/verify the certs
	return true
}

func IsRootMaterialPresent() bool {
	if _, err := os.Stat(config.SwarmRootCAVolumeMount + "/cert.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.SwarmRootCAVolumeMount + "/key.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.OrcaRootCAVolumeMount + "/cert.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(config.OrcaRootCAVolumeMount + "/key.pem"); os.IsNotExist(err) {
		return false
	}
	return true
}

// isFreshInstall checks whether the ucp-instance-key.pem file can be found
// at the expected location of the ucp-cluster-root-ca volume
func IsFreshInstall() bool {
	if _, err := os.Stat(config.ReconcileOrcaKeyFileMount); os.IsNotExist(err) {
		return false
	}
	return true
}

func CleanupFreshInstallFiles() error {
	err := os.Remove(config.ReconcileOrcaKeyFileMount)
	if err != nil {
		return err
	}
	return os.Remove(config.ReconcileOrcaCredFileMount)
}

// getUCPInstanceKey returns the private key stored in the ucp-instance-key.pem file
func GetUCPInstanceKey() (string, error) {
	file, err := os.Open(config.ReconcileOrcaKeyFileMount)
	if err != nil {
		return "", err
	}
	defer file.Close()
	res, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// getUCPCredentials returns the username/password provided during install time
func GetUCPCredentials() (string, string, error) {
	file, err := os.Open(config.ReconcileOrcaCredFileMount)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	res, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", err
	}
	split := strings.SplitN(string(res), "\n", 2)
	return split[0], split[1], nil
}
