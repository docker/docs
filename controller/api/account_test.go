package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/mock_test"
	"github.com/stretchr/testify/assert"
)

func TestApiGetAccounts(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}
	wrap := func(w http.ResponseWriter, req *http.Request) {
		// Inject bogus admin request context
		rc := ctx.MockAdmin(req)
		api.accounts(w, rc)
	}

	ts := httptest.NewServer(http.HandlerFunc(wrap))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	accts := []*auth.Account{}
	if err := json.NewDecoder(res.Body).Decode(&accts); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, len(accts), 0, "expected accounts; received none")

	acct := accts[0]

	assert.Equal(t, acct.Username, mock_test.TestAccount.Username, fmt.Sprintf("expected username %s; got %s", mock_test.TestAccount.Username, acct.Username))
}

func TestApiPostAccounts(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	wrap := func(w http.ResponseWriter, req *http.Request) {
		// Inject bogus admin request context
		rc := ctx.MockAdmin(req)
		api.saveAccount(w, rc)
	}

	ts := httptest.NewServer(http.HandlerFunc(wrap))
	defer ts.Close()

	data := []byte(`{"username": "testuser", "password": "foo"}`)

	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 204, "expected response code 204")
}

func TestApiDeleteAccount(t *testing.T) {
	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	transport := &http.Transport{}
	client := &http.Client{Transport: transport}

	wrap := func(w http.ResponseWriter, req *http.Request) {
		rc := ctx.MockAdmin(req)
		rc.PathVars["username"] = mock_test.TestOrcaUser
		api.deleteAccount(w, rc)
	}

	ts := httptest.NewServer(http.HandlerFunc(wrap))
	defer ts.Close()

	// Delete the mock user
	req, err := http.NewRequest("DELETE", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 204, "expected response code 204")
}
