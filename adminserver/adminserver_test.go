package adminserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/garant/authn"
	"github.com/docker/dhe-deploy/garant/authz"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/memory"
	"github.com/docker/dhe-deploy/licensing"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/manager/versions"

	"github.com/docker/garant/auth/common"
	garantconfig "github.com/docker/garant/config"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	newestVersion = "0.2.0-alpha-000697_g23f938c"
	containerID   = "8d716166e5cb4a627c0d1755f4dc2776f740007673a594a17426022b2ae4ba20"
)

var dockerVersion = dockerclient.Version{
	ApiVersion: "1.20",
	Version:    "1.9.0",
}

func TestNewAdminServer(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	assert.NotNil(t, setup.adminServer)
}

func TestIndexHandler(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/", nil)
	assert.Nil(t, err)

	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
	respBody := resp.Body.String()
	assert.Contains(t, respBody, "Docker Trusted Registry")
}

func TestVersionHandler(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/version", nil)
	assert.Nil(t, err)

	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
	assert.Equal(t, deploy.Version, resp.Body.String(), "Response body doesn't contain version number")
}

// TODO maybe bring this back later???
// func TestCsrfTokenHandlerNoCookie(t *testing.T) {
// 	setup, deferFunc := buildTestingSetup(t)
// 	defer deferFunc()
// 	setup.ProvideAdminAuthentication()
// 	setup.handler.ServeHTTP(resp, req)
// 	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
// 	assert.Contains(t, resp.Body.String(), `<meta name='csrf' content='' />`, "Expected CSRF tag to have empty CSRF token")
// }

// we can't run tests against the go-restful api here because it's not wired up
// TODO: wire up the go-restful api here and/or move this test to the api dir
//func TestGetSettingsHandler(t *testing.T) {
//	setup, deferFunc := buildTestingSetup(t)
//	defer deferFunc()
//	setup.ProvideAdminAuthentication()
//	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
//	resp := httptest.NewRecorder()
//	req, err := http.NewRequest("GET", "/api/v0/meta/settings", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	setup.handler.ServeHTTP(resp, req)
//	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
//}

// we can't run tests against the go-restful api here because it's not wired up
// TODO: wire up the go-restful api here and/or move this test to the api dir
//func TestGetSettingsHandlerAuthBackend(t *testing.T) {
//	setup, deferFunc := buildTestingSetup(t)
//	defer deferFunc()
//	setup.ProvideAdminAuthentication()
//	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
//	resp := httptest.NewRecorder()
//	req, err := http.NewRequest("GET", "/api/v0/meta/settings", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = setup.adminServer.settingsStore.SetAuthConfig(&garantconfig.Configuration{})
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer func() {
//		setup.adminServer.settingsStore.SetAuthConfig(&garantconfig.Configuration{})
//	}()
//
//	setup.handler.ServeHTTP(resp, req)
//	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
//}

func TestGetLoginHandler(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()

	type testCase struct {
		authenticated    bool
		hasAuthErr       bool
		expectedStatus   int
		expectedRedirect string
	}

	for _, test := range []testCase{
		{
			authenticated:    true,
			hasAuthErr:       false,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/",
		},
		{
			authenticated:    false,
			hasAuthErr:       false,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/api/v0/openid/begin",
		},
		{
			authenticated:    false,
			hasAuthErr:       false,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/api/v0/openid/begin",
		},
		{
			authenticated:    false,
			hasAuthErr:       true,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/api/v0/openid/begin",
		},
	} {
		enziAccount := &enziresponses.Account{
			Name:    "testuser",
			IsAdmin: &[]bool{true}[0],
		}

		var user *authn.User
		if test.authenticated {
			user = &authn.User{
				EnziSession: nil,
				Account:     enziAccount,
			}
		} else {
			user = &authn.User{
				IsAnonymous: true,
			}
		}

		var authErr error
		if test.hasAuthErr {
			authErr = &common.ClientError{Err: fmt.Errorf("some auth error"), Code: http.StatusUnauthorized}
		}
		setup.mockAuthorizer.On("AuthenticateRequestUser", mock.Anything, mock.AnythingOfType("*http.Request")).Return(user, authErr).Once()

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/admin/login", nil)
		assert.Nil(t, err)

		setup.handler.ServeHTTP(resp, req)
		assert.Equal(t, test.expectedStatus, resp.Code, "Unexpected status code")

		if test.expectedStatus == http.StatusFound {
			assert.Equal(t, test.expectedRedirect, resp.Header().Get("Location"), "Unexpected redirect location")
		}
		setup.mockAuthorizer.AssertExpectations(t)
	}
}

func TestGetLoginHandlerAuthorized(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusFound, resp.Code, "Unexpected status code")
	assert.Equal(t, resp.HeaderMap.Get("Location"), "/", "Unexpected redirect location")
}

func TestDockerHubLoginHandler(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
	authConfig := "{\"username\": \"dummy\", \"password\": \"dummy\"}"
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v0/admin/dhlogin", strings.NewReader(authConfig))
	if err != nil {
		t.Fatal(err)
	}

	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
}

func TestDockerHubLoginHandlerInvalidPayload(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
	authConfig := "finalesfunkeln"
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v0/admin/dhlogin", strings.NewReader(authConfig))
	if err != nil {
		t.Fatal(err)
	}

	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code, "Unexpected status code")
}

// we can't run tests against the go-restful api here because it's not wired up
// TODO: wire up the go-restful api here and/or move this test to the api dir
//func TestUpdateSettingsInvalidPayload(t *testing.T) {
//	setup, deferFunc := buildTestingSetup(t)
//	defer deferFunc()
//	setup.ProvideAdminAuthentication()
//	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
//	data := "perfectcherryblossom"
//	resp := httptest.NewRecorder()
//	req, err := http.NewRequest("PUT", "/api/v0/meta/settings", strings.NewReader(data))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	setup.handler.ServeHTTP(resp, req)
//	assert.Equal(t, http.StatusBadRequest, resp.Code, "Unexpected status code")
//}

func TestUpdateSettingsLicense(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
	settings := &hubconfig.LicenseConfig{
		KeyID:         "marisa",
		PrivateKey:    "kirisame",
		Authorization: "ifidonthavetherightilltaketheleft",
	}
	resp := makeRequestWithJSON(t, setup, "PUT", "/api/v0/admin/settings/license", settings)
	assert.Equal(t, http.StatusAccepted, resp.Code, "Unexpected status code")
}

func TestAuthenticationForAPI(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	for _, test := range []struct {
		authenticated  bool
		hasAuthErr     bool
		expectedStatus int
	}{
		{
			authenticated:  false,
			hasAuthErr:     false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			authenticated:  false,
			hasAuthErr:     true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			authenticated:  true,
			hasAuthErr:     false,
			expectedStatus: http.StatusOK,
		},
	} {
		var authErr error
		if test.hasAuthErr {
			authErr = &common.ClientError{Err: fmt.Errorf("some auth error"), Code: http.StatusUnauthorized}
		}

		enziAccount := &enziresponses.Account{
			Name:    "testuser",
			IsAdmin: &[]bool{true}[0],
		}

		var user *authn.User
		if test.authenticated {
			user = &authn.User{
				EnziSession: nil,
				Account:     enziAccount,
			}
		} else {
			user = &authn.User{
				IsAnonymous: true,
			}
		}

		setup.mockAuthorizer.On("AuthenticateRequestUser", mock.Anything, mock.AnythingOfType("*http.Request")).Return(user, authErr).Once()
		setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
		setup.mockVersionChecker.On("VersionList", mock.AnythingOfType("*dockerclient.AuthConfig")).Return(versions.ManagerVersionList{newestVersion}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v0/admin/upgrade", nil)
		assert.Nil(t, err)
		setup.handler.ServeHTTP(resp, req)

		assert.Equal(t, test.expectedStatus, resp.Code, "Unexpected status code")
		setup.mockAuthorizer.AssertExpectations(t)
	}
}

type testingSetup struct {
	adminServer        *AdminServer
	handler            http.Handler
	mockSyslogWriter   *MockSyslogWriter
	mockAuthorizer     *authz.MockAuthorizer
	mockVersionChecker *versions.MockChecker
	alertsDir          string
}

func buildTestingSetup(t *testing.T) (testingSetup, func()) {
	var (
		settingsStore = memory.NewSettingsStore()
		kvStore       = hubconfig.NewMockKeyValueStore()
		userHubConfig = &hubconfig.UserHubConfig{
			DTRHost: "gensokyo.jp",
		}
		haConfig           = &defaultconfigs.DefaultHAConfig
		garantCfg          = &garantconfig.Configuration{}
		mockVersionChecker = new(versions.MockChecker)
		lChecker           = new(licensing.MockChecker)
		mockAuthorizer     = new(authz.MockAuthorizer)
	)

	if err := settingsStore.SetUserHubConfig(userHubConfig); err != nil {
		t.Fatal(err)
	}
	if err := settingsStore.SetHAConfig(haConfig); err != nil {
		t.Fatal(err)
	}
	if err := settingsStore.SetAuthConfig(garantCfg); err != nil {
		t.Fatal(err)
	}

	session, _, deferFunc, err := schema.TestSetup()
	if err != nil {
		t.Fatal(err)
	}

	tmpAlertsDir, err := ioutil.TempDir("", "alerts")
	assert.Nil(t, err)

	mockVersionChecker.On("NewestVersion", mock.AnythingOfType("*dockerclient.AuthConfig")).Return(newestVersion, nil)
	mockVersionChecker.On("VersionList", mock.AnythingOfType("*dockerclient.AuthConfig")).Return(versions.ManagerVersionList{newestVersion}, nil)
	lChecker.On("LicensingEnforced").Return(true)
	lChecker.On("IsValid").Return(true)
	lChecker.On("LicenseTier").Return("")
	lChecker.On("LicenseType").Return("")
	lChecker.On("GetLicenseID").Return("")
	lChecker.On("IsExpired").Return(false)
	lChecker.On("LoadLicenseFromConfig", mock.AnythingOfType("*hubconfig.LicenseConfig"), mock.AnythingOfType("bool")).Return(nil)

	mockSyslogWriter := new(MockSyslogWriter)

	adminServer, err := New(settingsStore, kvStore, http.DefaultClient, http.DefaultClient, lChecker, mockVersionChecker, mockSyslogWriter, mockAuthorizer, session)
	assert.Nil(t, err)
	setup := testingSetup{
		adminServer:        adminServer,
		handler:            adminServer.wrapHandler(adminServer.buildRouter()),
		mockSyslogWriter:   mockSyslogWriter,
		alertsDir:          tmpAlertsDir,
		mockAuthorizer:     mockAuthorizer,
		mockVersionChecker: mockVersionChecker,
	}

	TemplatesDir = "./ui/src"
	return setup, deferFunc
}

func makeRequestWithJSON(t *testing.T, setup testingSetup, method, path string, object interface{}) *httptest.ResponseRecorder {
	data, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	setup.handler.ServeHTTP(resp, req)
	return resp
}

func uintPtr(i uint16) *uint16 {
	return &i
}

func (ts testingSetup) ProvideAdminAuthentication() {
	enziAccount := &enziresponses.Account{
		Name:     "mrpoopybutthole",
		ID:       "ididid",
		FullName: "Mr. Poopy Butthole",
		IsAdmin:  &[]bool{true}[0],
		IsActive: &[]bool{true}[0],
	}
	ts.mockAuthorizer.On("AuthenticateRequestUser", mock.Anything, mock.Anything).Return(&authn.User{
		EnziSession: nil,
		Account:     enziAccount,
	}, nil)
}
