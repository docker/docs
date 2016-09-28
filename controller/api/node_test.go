package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/orca"
	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/docker/orca/controller/mock_test"
	"github.com/stretchr/testify/assert"
)

func TestApiGetNodes(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.nodes)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, res.StatusCode, "expected response code 200")
	nodes := []*orca.Node{}
	if err := json.NewDecoder(res.Body).Decode(&nodes); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, len(nodes), 0, "expected nodes; received none")

	node := nodes[0]

	assert.Equal(t, node.Name, mock_test.TestNode.Name, fmt.Sprintf("expected name %s; got %s", mock_test.TestNode.Name, node.Name))
}

func TestApiGetNode(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.node)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, res.StatusCode, "expected response code 200")
	node := &orca.Node{}
	if err := json.NewDecoder(res.Body).Decode(&node); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, node.ID, nil, "expected node; received nil")

	assert.Equal(t, node.Name, mock_test.TestNode.Name, fmt.Sprintf("expected name %s; got %s", mock_test.TestNode.Name, node.Name))
}

func TestApiPostAddNode(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(api.authorizeNodeRequest))
	defer ts.Close()

	data := []byte(`{"certificate_request": "-----REAL CSR WOULD GO HERE-----", "callback_url": "http://localhost/foo", "name": "foo-node"}`)

	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, res.StatusCode, "expected response code 200")
}
