package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/docker/orca/controller/mock_test"
	"github.com/stretchr/testify/assert"
)

func TestApiGetConsoleSession(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.consoleSession)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	cs := &orca.ConsoleSession{}
	if err := json.NewDecoder(res.Body).Decode(&cs); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, cs.ContainerID, nil, "expected console session; received nil")

	assert.Equal(t, cs.ContainerID, mock_test.TestConsoleSession.ContainerID, fmt.Sprintf("expected container id %s; got %s", mock_test.TestConsoleSession.ContainerID, cs.ContainerID))
}

func TestApiPostConsoleSessions(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.createConsoleSession)
	defer ts.Close()

	data := []byte(`{"container_id": "abcdefg", "token": "12345"}`)

	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
}

func TestApiDeleteConsoleSession(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	transport := &http.Transport{}
	client := &http.Client{Transport: transport}

	ts := pipeline.MockTestServer(api.removeConsoleSession)
	defer ts.Close()

	data := []byte(`{"id": "0"}`)

	req, err := http.NewRequest("DELETE", ts.URL, bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
}
