package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

type RegistryToken struct {
	Token string `json:"token"`
}

var (
	PrivateImage = "dockerorcadev/ucp"
)

// Given an Orca client, attempt to pull a private image with a token from hub
func TestPrivatePullWithToken(t *testing.T, dclient *dockerclient.DockerClient) {
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	require := require.New(t)

	if username == "" || password == "" {
		t.Skip("To run the token pull test you must set $REGISTRY_USERNAME and $REGISTRY_PASSWORD")
		return
	}

	// Get the token
	client := &http.Client{}
	reqURL := fmt.Sprintf("https://auth.docker.io/token?scope=repository:%s:pull&service=registry.docker.io", PrivateImage)
	log.Info("Getting a registry token for %s", reqURL)
	req, err := http.NewRequest("GET", reqURL, nil)
	require.Nil(err)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	require.Nil(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(err)

	var token RegistryToken
	err = json.Unmarshal(body, &token)
	require.Nil(err)

	log.Debugf("Got token %s", token.Token)
	authConfig := &dockerclient.AuthConfig{RegistryToken: token.Token}
	err = dclient.PullImage(PrivateImage+":latest", authConfig)
	require.Nil(err, "Failed to pull a private hub image with a token")
	log.Info("Succesfully pulled private hub image with a token")
}
