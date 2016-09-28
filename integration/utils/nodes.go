package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/orca"
	"golang.org/x/net/context"
)

func GetNodes(serverURL, username, password string) ([]swarm.Node, error) {
	client, err := GetUserEngineAPI(serverURL, username, password)
	if err != nil {
		return nil, err
	}
	return client.NodeList(context.TODO(), types.NodeListOptions{})
}

// Keep around for querying pre 1.12 based releases
func OldGetNodes(client *http.Client, serverURL, adminUser, adminPassword string) ([]orca.Node, error) {

	nodes := []orca.Node{}
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return nodes, err
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	// Login and get a token for the admin user
	token, err := GetOrcaToken(client, serverURL, adminUser, adminPassword)
	if err != nil {
		return nodes, err
	}

	orcaURL.Path = "/api/nodes"

	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return nodes, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return nodes, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nodes, fmt.Errorf(string(body))
	}

	err = json.Unmarshal(body, &nodes)
	return nodes, err
}
