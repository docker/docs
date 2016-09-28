package simple

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
)

type TestSuite struct {
	utils.OrcaTestSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

func (s *TestSuite) InstallArgs(m utils.Machine) []string {
	log.Debug("Wiring install up for non-standard ports")
	externalIP, err := m.GetIP()
	require.Nil(s.T(), err)
	return []string{"install", "--disable-tracking", "--disable-usage", "-D", "--swarm-port", "4376", "--controller-port", "8443", "--san", externalIP}
}

// Override so we don't try to add a license
func (s *TestSuite) PostInstall(m utils.Machine) error {
	return nil
}

func (s *TestSuite) TestInstallWithAltControllerPort() {
	ip, err := s.ControllerMachines[0].GetIP()
	require.Nil(s.T(), err)
	uri := fmt.Sprintf("https://%s:8443/_ping", ip)
	expected := 200

	// For this scenario, we're not going to bother validating the certificate
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	// Retry a bit incase things are coming up
	var lastError error
	for i := 0; i < 60; i += 1 {
		resp, err := client.Get(uri)
		if err != nil {
			log.Debug("Failed to connect to orca at %s - %s", uri, err)
			lastError = err
		} else if resp.StatusCode != expected {
			body, _ := ioutil.ReadAll(resp.Body)
			lastError = fmt.Errorf("Unexpected status code: %d - Payload %s", resp.StatusCode, body)
		} else {
			lastError = nil
			break
		}
		time.Sleep(1 * time.Second)
	}
	require.Nil(s.T(), lastError)
}

func TestInstallTestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
