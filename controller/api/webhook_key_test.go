package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/docker/orca/dockerhub"
	"github.com/stretchr/testify/assert"
)

func TestApiGetWebhookKeys(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.webhookKeys)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	keys := []*dockerhub.WebhookKey{}

	if err := json.NewDecoder(res.Body).Decode(&keys); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, len(keys), 0, "expected keys; received none")

	key := "abcdefg"

	assert.Equal(t, keys[0].Key, key, "expected key %s; received %s", key, keys[0].Key)
}
