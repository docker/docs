package adminserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/distribution/configuration"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type configUpdate struct {
	name string
	// Data we PUT into the server
	data string

	expectedStatus int
	expectedConfig configuration.Storage
}

func TestUpdateRegistrySettingsViaFormHandler(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil)

	// Test updating registry settings
	fsConfig := configUpdate{
		name:           "Test updating filesystem settings",
		data:           `{"filesystem":{"rootdirectory":"/tmp"}}`,
		expectedStatus: http.StatusAccepted,
		expectedConfig: configuration.Storage{
			"filesystem": map[string]interface{}{
				"rootdirectory": "/storage",
			},
		},
	}
	s3Config := configUpdate{
		name:           "Test updating S3 settings",
		data:           `{"s3":{"region":"us-east-1", "bucket": "bucketnamelol", "accesskey": "myaccess", "secretkey": "inplaintext?!", "rootdirectory": "iamroot"}}`,
		expectedStatus: http.StatusAccepted,
		expectedConfig: configuration.Storage{
			"s3": map[string]interface{}{
				"region":        "us-east-1",
				"bucket":        "bucketnamelol",
				"accesskey":     "myaccess",
				"secretkey":     "inplaintext?!",
				"rootdirectory": "iamroot",
			},
		},
	}
	// This test always fails because it forces dtr to try to actually log in with that account. We can't mock this right now, so the test is disabled
	//{
	//	name:           "Test updating Azure settings",
	//	data:           `{"azure":{"accountname":"godzilla", "accountkey": "deadbeef", "container": "containername"}}`,
	//	expectedStatus: http.StatusBadRequest,
	//	expectedConfig: configuration.Storage{
	//		"azure": map[string]interface{}{
	//			"accountname": "godzilla",
	//			"accountkey":  "somekey",
	//			"container":   "containername",
	//		},
	//	},
	//},

	// first try with s3
	resp := updateConfig(t, setup, s3Config)

	// we should get invalid access key when we test write/read
	require.Contains(t, resp.Body.String(), "InvalidAccessKeyId")
	require.Equal(t, resp.Code, http.StatusBadRequest)

	// try with fs
	resp = updateConfig(t, setup, fsConfig)

	// Get the updated storage settings information
	config, err := setup.adminServer.settingsStore.RegistryConfig()
	require.Nil(t, err)
	require.NotNil(t, config)
	actual := config.Storage
	util.SetReadonlyModeJSON(&fsConfig.expectedConfig, util.GetReadonlyMode(&fsConfig.expectedConfig))
	require.Equal(t, fsConfig.expectedStatus, resp.Code, "Unexpected status code")
	require.Equal(t, fsConfig.expectedConfig, actual, "Unexpected configuration.Storage")
}

func updateConfig(t *testing.T, setup testingSetup, conf configUpdate) *httptest.ResponseRecorder {
	req, err := http.NewRequest("PUT", "/api/v0/admin/settings/registry/simple", strings.NewReader(conf.data))
	require.Nil(t, err)

	resp := httptest.NewRecorder()
	setup.handler.ServeHTTP(resp, req)

	return resp
}
