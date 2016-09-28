package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/mock_test"
	"github.com/gorilla/mux"
	"github.com/mailgun/oxy/forward"
	"github.com/stretchr/testify/assert"
)

func insertUser(fn func(http.ResponseWriter, *ctx.OrcaRequestContext), username string) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		// Mock the provided username as an admin user
		rc := ctx.MockAdmin(r)
		rc.Auth.User.Username = username
		fn(w, rc)
	})
}

func getForwardServer(api *Api, router *mux.Router) (*httptest.Server, error) {
	var err error

	// Our forwarder hasn't been set up since api.Run() is never called, so let's create a
	// new forwarder
	api.fwd, err = forward.New()
	if err != nil {
		return nil, err
	}

	ts := httptest.NewServer(router)
	return ts, nil
}

func TestApiImagePush(t *testing.T) {
	var buf io.Reader
	var outHeaders http.Header
	testResponse := "hello there\n"

	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	gs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outHeaders = r.Header
		fmt.Fprint(w, testResponse)
	}))
	defer gs.Close()
	api.swarmClassicURL = gs.URL

	r := mux.NewRouter()
	r.HandleFunc(
		"/v1.22/images/{name:.*}/push",
		http.HandlerFunc(insertUser(api.swarmRegistryImagesPush, mock_test.TestRepositoryUser)))

	ts, err := getForwardServer(api, r)
	defer ts.Close()
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("%s/v1.22/images/%s/push", ts.URL, mock_test.TestRepositoryName)

	res, err := http.Post(url, "application/x-www-form-urlencoded", buf)
	if err != nil {
		t.Fatal(err)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	assert.Equal(t, string(d), testResponse, "wrong response from server")

	// Check that the token was inserted properly
	regAuth := outHeaders.Get("X-Registry-Auth")
	data, err := base64.StdEncoding.DecodeString(regAuth)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(data), `{"registrytoken": "foo"}`, "unexpected X-Registry-Auth header")
}

func TestApiImagePushExistingHeader(t *testing.T) {
	var outHeaders http.Header
	testResponse := "hello there\n"

	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	gs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outHeaders = r.Header
		fmt.Fprint(w, testResponse)
	}))
	defer gs.Close()
	api.swarmClassicURL = gs.URL

	r := mux.NewRouter()
	r.HandleFunc(
		"/v1.22/images/{name:.*}/push",
		http.HandlerFunc(insertUser(api.swarmRegistryImagesPush, mock_test.TestRepositoryUser)))

	ts, err := getForwardServer(api, r)
	defer ts.Close()
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("%s/v1.22/images/%s/push", ts.URL, mock_test.TestRepositoryName)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	b64Auth := base64.URLEncoding.EncodeToString([]byte(mock_test.TestRegistryAuth))
	req.Header.Add("X-Registry-Auth", b64Auth)
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, fmt.Sprintf("expected response code 200, got: %s", res.Status))

	// Check that the token was not inserted
	regAuth := outHeaders.Get("X-Registry-Auth")
	assert.Equal(t, regAuth, b64Auth)
}

func TestApiImageCreate(t *testing.T) {
	var buf io.Reader
	var outHeaders http.Header
	testResponse := "hello there\n"

	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	gs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outHeaders = r.Header
		fmt.Fprint(w, testResponse)
	}))
	defer gs.Close()
	api.swarmClassicURL = gs.URL

	r := mux.NewRouter()
	r.HandleFunc(
		"/v1.22/images/create",
		http.HandlerFunc(insertUser(api.swarmRegistryImagesCreate, mock_test.TestRepositoryUser)))

	ts, err := getForwardServer(api, r)
	defer ts.Close()
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("%s/v1.22/images/create?fromImage=%s", ts.URL, mock_test.TestRepositoryName)

	res, err := http.Post(url, "application/x-www-form-urlencoded", buf)
	if err != nil {
		t.Fatal(err)
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, "expected response code 200")
	assert.Equal(t, string(d), testResponse, "wrong response from server")

	// Check that the token was inserted properly
	regAuth := outHeaders.Get("X-Registry-Auth")
	data, err := base64.StdEncoding.DecodeString(regAuth)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(data), `{"registrytoken": "foo"}`, "unexpected X-Registry-Auth header")
}

func TestApiImageCreateExistingHeader(t *testing.T) {
	var outHeaders http.Header
	testResponse := "hello there\n"

	api, err := getTestApi()
	if err != nil {
		t.Fatal(err)
	}

	gs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outHeaders = r.Header
		fmt.Fprint(w, testResponse)
	}))
	defer gs.Close()
	api.swarmClassicURL = gs.URL

	r := mux.NewRouter()
	r.HandleFunc(
		"/v1.22/images/create",
		http.HandlerFunc(insertUser(api.swarmRegistryImagesCreate, mock_test.TestRepositoryUser)))

	ts, err := getForwardServer(api, r)
	defer ts.Close()
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("%s/v1.22/images/create?fromImage=%s", ts.URL, mock_test.TestRepositoryName)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	b64Auth := base64.URLEncoding.EncodeToString([]byte(mock_test.TestRegistryAuth))
	req.Header.Add("X-Registry-Auth", b64Auth)
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.StatusCode, 200, fmt.Sprintf("expected response code 200, got: %s", res.Status))

	// Check that the token was not inserted
	regAuth := outHeaders.Get("X-Registry-Auth")
	assert.Equal(t, regAuth, b64Auth, "unexpected X-Registry-Auth header")
}
