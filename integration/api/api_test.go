package api

import (
	"crypto/tls"
	"net/http"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/controller/api"
	"github.com/docker/orca/controller/middleware/pipeline"
	"github.com/docker/orca/controller/mock_test"
	"github.com/docker/orca/integration/utils"
)

// TEST_VERSION is the value of the version URL argument, satisfying [0-9.]+
// The current implementation assumes that TEST_VERSION is not a compatible version
const TEST_VERSION = "0.9"

// TIMEOUT is the maximum expected duration of time for any API endpoint.
// Should cover the expected runtime of the /api/support endpoint on any test environment
const TIMEOUT = 25 * time.Second

type APITestSuite struct {
	utils.OrcaTestSuite

	pipeline TestPipeline

	APIPaths map[string]int
}

func (s *APITestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 1
	workerCount = 0
	return
}

// Define a middleware pipeline that holds a list of all defined routes and methods
// TODO(alexmavr): use the pipeline to perform assertions on a per-method basis
type TestPipeline struct {
	pipeline.MiddlewarePipeline
	Routes map[string]map[string]bool
}

func (p TestPipeline) Route(path string, method string, parser pipeline.Parser, handler pipeline.Handler) {
	if _, ok := p.Routes[path]; !ok {
		// path does not exist in Routes
		p.Routes[path] = make(map[string]bool)
	}
	pathRoute, _ := p.Routes[path]
	if _, ok := pathRoute[method]; !ok {
		// path does not exist in Routes
		pathRoute[method] = true
	}
	p.MiddlewarePipeline.Route(path, method, parser, handler)
}

func NewTestPipeline(r *mux.Router) TestPipeline {
	return TestPipeline{
		MiddlewarePipeline: pipeline.New(r),
		Routes:             make(map[string]map[string]bool),
	}
}

func (s *APITestSuite) PreInstall(m utils.Machine) error {
	require := require.New(s.T())

	// Create an ApiConfig with a mock manager
	apicfg := api.ApiConfig{
		Manager: new(mock_test.MockManager),
	}

	s.APIPaths = make(map[string]int)

	// Create a Mock Api object
	a, err := api.NewApi(apicfg)
	require.Nil(err)
	require.NotNil(a)

	// Generate all Routes
	s.pipeline = NewTestPipeline(a.Router)
	a.Initialize(s.pipeline)

	return a.Router.Walk(s.APIWalker)
}

// APIWalker parses a route to add a path to s.APIPaths or increase the count of an existing path.
func (s *APITestSuite) APIWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// Inject dummy URL parameters to generate a final URL path
	url, err := route.URLPath("username", "username",
		"id", "id",
		"teamId", "teamId",
		"subsystem", "subsystem",
		"name", "name",
		"key", "key",
		"token", "token",
		"file", "file",
		"version", TEST_VERSION)
	require.Nil(s.T(), err)
	require.NotNil(s.T(), url)

	if _, ok := s.APIPaths[url.Path]; !ok {
		s.APIPaths[url.Path] = 0
	}
	s.APIPaths[url.Path] += 1

	return err
}

// TestAPINoAuth hits all API paths with every HTTP method and no Auth.
func (s *APITestSuite) TestAPINoAuth() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	require := require.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	tr := &http.Transport{
		// Skip cert verification
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	client := &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}

	for path, _ := range s.APIPaths {
		fullPath := serverURL + path
		log.Debug(fullPath)

		// Create a slice of dummy requests
		reqGet, err := http.NewRequest("GET", fullPath, nil)
		require.Nil(err)
		reqPost, err := http.NewRequest("POST", fullPath, nil)
		require.Nil(err)
		reqDelete, err := http.NewRequest("DELETE", fullPath, nil)
		require.Nil(err)
		reqPut, err := http.NewRequest("PUT", fullPath, nil)
		require.Nil(err)
		requests := []*http.Request{reqGet, reqPost, reqDelete, reqPut}

		for _, req := range requests {
			time.Sleep(15 * time.Millisecond)
			req.Close = true
			tr.CloseIdleConnections()
			resp, err := client.Do(req)
			if err != nil {
				log.Debug(err)
			}
			require.Nil(err)
			require.NotNil(resp)

			s.requireUnauthorized(resp, req)
		}
	}
}

// TestStatic explicitly tests whether /bundle.js is accessible
func (s *APITestSuite) TestStatic() {
	require := require.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	tr := &http.Transport{
		// Skip cert verification
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}

	resp, err := client.Get(serverURL + "/bundle.js")
	if err != nil {
		log.Debug(err)
	}
	require.Nil(err)
	require.Equal(200, resp.StatusCode)
}

// TestAPIAuthTLS hits all API paths with every HTTP method and TLS auth as admin.
func (s *APITestSuite) TestAPIAuthTLS() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	require := require.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	// Obtain user certs and generate a tls config
	tlsConfig, err := utils.GetUserTLSConfig(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)
	tlsConfig.InsecureSkipVerify = true

	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	client := &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}

	for path, _ := range s.APIPaths {
		fullPath := serverURL + path
		log.Debug(fullPath)

		// Create a slice of dummy requests
		reqGet, err := http.NewRequest("GET", fullPath, nil)
		require.Nil(err)
		reqPost, err := http.NewRequest("POST", fullPath, nil)
		require.Nil(err)
		reqDelete, err := http.NewRequest("DELETE", fullPath, nil)
		require.Nil(err)
		reqPut, err := http.NewRequest("PUT", fullPath, nil)
		require.Nil(err)
		requests := []*http.Request{reqGet, reqPost, reqDelete, reqPut}

		for _, req := range requests {
			time.Sleep(15 * time.Millisecond)
			tr.CloseIdleConnections()
			req.Close = true
			resp, err := client.Do(req)
			if err != nil {
				log.Debug(err)
			}
			require.Nil(err)

			s.requireAuthorized(resp, req)
		}
	}
}

// TestAPIAuthJWT hits all API paths with every HTTP method while authenticated as admin with JWT.
func (s *APITestSuite) TestAPIAuthJWT() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}
	require := require.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	// Skip cert verification
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	client := &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}

	// Authenticate with UCP as Admin and obtain JWT token
	token, err := utils.GetOrcaToken(client, serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	require.Nil(err)

	for path, _ := range s.APIPaths {
		// Skip the login/logout endpoints to avoid messing with the auth token
		if strings.Contains(path, "auth/log") {
			continue
		}
		fullPath := serverURL + path
		log.Debug(fullPath)

		// Create a slice of dummy requests
		reqGet, err := http.NewRequest("GET", fullPath, nil)
		require.Nil(err)
		reqPost, err := http.NewRequest("POST", fullPath, nil)
		require.Nil(err)
		reqDelete, err := http.NewRequest("DELETE", fullPath, nil)
		require.Nil(err)
		reqPut, err := http.NewRequest("PUT", fullPath, nil)
		require.Nil(err)
		requests := []*http.Request{reqGet, reqPost, reqDelete, reqPut}

		for _, req := range requests {
			time.Sleep(15 * time.Millisecond)
			req.Header.Set(utils.GetTokenHeader(token))
			req.Close = true
			tr.CloseIdleConnections()
			resp, err := client.Do(req)
			if err != nil {
				log.Debug(err)
			}
			require.Nil(err)

			s.requireAuthorized(resp, req)
		}
	}
}

// requireUnauthorized performs assertions for the behavior of unauthorized requests
func (s *APITestSuite) requireUnauthorized(resp *http.Response, req *http.Request) {
	require := require.New(s.T())
	path := req.URL.Path
	if path == "/" ||
		strings.Contains(path, "/_ping") ||
		path == "/openid_keys" ||
		path == "/ca" {
		require.True(200 == resp.StatusCode || 404 == resp.StatusCode)
	} else {
		require.True(400 <= resp.StatusCode && 500 >= resp.StatusCode)
	}
}

// requireAuthorized performs assertions for the behavior of authorized requests
func (s *APITestSuite) requireAuthorized(resp *http.Response, req *http.Request) {
	require := require.New(s.T())
	path := req.URL.Path

	if path == "/" {
		require.Equal(200, resp.StatusCode)
	} else {
		// Everything else should be one of the following status codes
		// The important assertion is no 401 or 403 response
		require.True(200 == resp.StatusCode ||
			201 == resp.StatusCode ||
			204 == resp.StatusCode ||
			400 == resp.StatusCode ||
			404 == resp.StatusCode ||
			405 == resp.StatusCode ||
			500 == resp.StatusCode ||
			501 == resp.StatusCode)
	}
}

func TestAPITestSuite(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &APITestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
