package integration

import (
	"testing"

	"github.com/docker/dhe-deploy/adminserver"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SettingsAPITestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util
}

func (suite *SettingsAPITestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
}

func (suite *SettingsAPITestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *SettingsAPITestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *SettingsAPITestSuite) TestAdminSettingsHttpHttpProxyNotSet() {
	httpSettings := suite.getHTTPSettings()

	// Test: check if NoProxy is not set
	if httpSettings.HTTPProxy != "" {
		suite.T().Fatal("Proxy should not be set")
	}
}

// This feature was removed. Now it needs to be done through the CLI
// TestAdminSettingsSetHTTPPort changes the HTTP port for DTR, while deferring the changes back to the original settings it started with.
//func (suite *SettingsAPITestSuite) TestAdminSettingsSetHTTPPort() {
//	util.AppendDockerIgnorableLoggedErrors(util.KnownDockerRestartLogErrors)
//	util.AppendDTRIgnorableLoggedErrors(util.KnownDTRRestartLogErrors)
//
//	// grab existing settings
//	originalHttpSettings := suite.getHTTPSettings()
//	originalApiClientPort := suite.API.GetApiClientPort()
//	originalApiClientUrlScheme := suite.API.GetApiClientUrlScheme()
//	// make a copy of the settings that get modified
//	httpSettings := *originalHttpSettings
//
//	// defer changing back API Client port/scheme, then revert HTTPSettings settings
//	defer func() {
//		suite.API.SetApiClientPort(originalApiClientPort)
//		suite.API.SetApiClientUrlScheme(originalApiClientUrlScheme)
//		suite.setHTTPSettings(originalHttpSettings)
//		if err := suite.u.PollAvailable(); err != nil {
//			suite.T().Fatalf("DTR failed to restart in a reasonable amount of time")
//		}
//	}()
//
//	// change DTR HTTP port to 81
//	httpSettings.LoadBalancerHTTPPort = 81
//	suite.setHTTPSettings(&httpSettings)
//	if err := suite.u.PollAvailable(); err != nil {
//		suite.T().Fatalf("DTR failed to restart in a reasonable amount of time")
//	}
//
//	// wait 4 seconds to gaurentee nginx check interval has the required 2 rise counts (1 per 2 seconds)
//	time.Sleep(time.Second * 4)
//
//	// verify the change in the http settings
//	httpSettingsAfterPortChange := suite.getHTTPSettings()
//	require.EqualValues(suite.T(), 81, httpSettingsAfterPortChange.LoadBalancerHTTPPort, "port is not set to 81")
//
//	// request the load balancer status on port 81 as a way to verify that port 81 is now open
//	suite.API.SetApiClientPort(81)
//	suite.API.SetApiClientUrlScheme("http")
//	nginxLoadBalancerStatus := suite.loadBalancerStatus()
//
//	// assert some of the data in nginx load balancer status
//	assert.EqualValues(suite.T(), 4, nginxLoadBalancerStatus.NginxServers.Total, "number of servers not correct")
//	for _, server := range nginxLoadBalancerStatus.NginxServers.NginxServer {
//		assert.Equal(suite.T(), "up", server.Status, "server is not up", server)
//	}
//}

// TestSetAndCheckValidLicense verifies that it is possible to set a license and then verify that the license is valid,
// as determined by DTR itself.  Verifies both setting and getting license info.
func (suite *SettingsAPITestSuite) TestSetAndCheckValidLicense() {
	util.AppendDockerIgnorableLoggedErrors(util.KnownDockerRestartLogErrors)
	util.AppendDTRIgnorableLoggedErrors(util.KnownDTRRestartLogErrors)

	licenseConfig, err := util.GetOnlineLicense()
	if err != nil {
		suite.T().Fatalf("Failed to get default online license: %s", err)
	}

	suite.testSetAndCheckValidLicense(licenseConfig)

	// licenseConfig, err = util.GetOfflineLicense()
	// if err != nil {
	// 	suite.T().Fatalf("Failed to get default offline license: %s", err)
	// }
	//
	// suite.testSetAndCheckValidLicense(licenseConfig)
}

func (suite *SettingsAPITestSuite) testSetAndCheckValidLicense(licenseConfig *hubconfig.LicenseConfig) {
	// set the license
	licenseSettings := suite.setLicenseSettings(licenseConfig)

	if err := suite.u.PollAvailable(); err != nil {
		suite.T().Fatalf("DTR failed to restart in a reasonable amount of time")
	}

	assert.True(suite.T(), licenseSettings.IsValid, "License is not valid")
	assert.Equal(suite.T(), licenseConfig.KeyID, licenseSettings.KeyID, "Mis-match license KeyIDs")
	// full circle , check also retrieving the license settings without setting the license ... should be the same
	licenseSettingsFromGet := suite.getLicenseSettings()
	assert.True(suite.T(), licenseSettingsFromGet.IsValid, "License is not valid")
	assert.Equal(suite.T(), licenseConfig.KeyID, licenseSettingsFromGet.KeyID, "Mis-match license KeyIDs")
}

func (suite *SettingsAPITestSuite) TestAPISettingsHeadersServerPresent() {
	resp, err := suite.API.GetHTTPSettingsResponse()
	if err != nil {
		suite.T().Fatal("Failed Request")
	}

	// check for Server header content
	if resp.Header.Get("Server") != "nginx/1.8.1" {
		suite.T().Fatal("Nginx Header not found")
	}
}

func (suite *SettingsAPITestSuite) getLicenseSettings() *adminserver.LicenseSettings {
	licenseSettings, err := suite.API.GetLicenseSettings()
	if err != nil {
		suite.T().Fatalf("Error: Problem with retrieving the license settings: %s", err)
	}
	return licenseSettings
}

func (suite *SettingsAPITestSuite) setLicenseSettings(licenseConfig *hubconfig.LicenseConfig) *adminserver.LicenseSettings {
	require.Nil(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	licenseSettings, err := suite.API.SetLicenseSettings(licenseConfig)
	if err != nil {
		suite.T().Fatalf("Error: Problem with setting the license settings: %s", err)
	}
	// set actually returns back a response
	return licenseSettings
}

// getHTTPSettings will retrieve the HTTP settings ... wrapper for api.GetHTTPSettings
func (suite *SettingsAPITestSuite) getHTTPSettings() *responses.Settings {
	httpSettings, err := suite.API.GetHTTPSettings()
	if err != nil {
		suite.T().Fatalf("Error: Problem with retrieving the API settings: %s", err)
	}

	return httpSettings
}

// setHTTPSettings will set the HTTP settings ... wrapper for api.SetHTTPSettings
func (suite *SettingsAPITestSuite) setHTTPSettings(settings *forms.Settings) {
	err := suite.API.SetHTTPSettings(settings)
	if err != nil {
		suite.T().Fatalf("Error: Problem with setting the API settings: %s", err)
	}
}

func (suite *SettingsAPITestSuite) loadBalancerStatus() (status *apiclient.NginxLoadBalancerStatus) {
	status, err := suite.API.LoadBalancerStatus()
	if err != nil {
		suite.T().Fatalf("Error request nginx load balancer status : %s", err)
	}

	return status
}

func TestSettingsAPISuite(t *testing.T) {
	suite.Run(t, new(SettingsAPITestSuite))
}
