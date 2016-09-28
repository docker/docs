package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestApiGetAuthTeams(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	ts := pipeline.MockTestServer(api.teams)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	teams := []*auth.Team{}

	if err := json.NewDecoder(res.Body).Decode(&teams); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, len(teams), 0, "expected teams; received none")
}
