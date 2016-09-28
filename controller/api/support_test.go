package api

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestApiPostSupport(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.supportDump)
	defer ts.Close()

	data := []byte(``)

	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
}
