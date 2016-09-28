package integration

import (
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var globalFramework *framework.IntegrationFramework
var setupOnce sync.Once

func setupFramework(suite suite.TestingSuite) (*framework.IntegrationFramework, *util.Util) {
	setupOnce.Do(func() {
		newFramework := framework.BuildFramework(suite.T())
		newUtil := util.MakeUtil(suite, newFramework)

		err := newFramework.API.Login(newFramework.Config.AdminUsername, newFramework.Config.AdminPassword)
		if urlError, ok := err.(*url.Error); ok && urlError.Err == apiclient.RedirectAttemptedError {
			err = nil
		}
		require.Nil(suite.T(), err, "%s", err)

		if err := newUtil.PollAvailable(); err != nil {
			suite.T().Fatalf("DTR is not running: %v", err)
		}

		if licenseSettings, err := newFramework.API.GetLicenseSettings(); err == nil && licenseSettings.IsValid == false {
			licenseConfig, err := util.GetOnlineLicense()
			if err != nil {
				suite.T().Fatalf("Failed to get default offline license: %s", err)
			}

			if _, err := newFramework.API.SetLicenseSettings(licenseConfig); err != nil {
				suite.T().Fatalf("Failed to set license: %s", err)
			}

			if err := newUtil.PollAvailable(); err != nil {
				suite.T().Fatalf("DTR failed to restart in a reasonable amount of time")
			}

			licenseSettings, err := newFramework.API.GetLicenseSettings()
			if err != nil {
				suite.T().Fatalf("Failed to get license settings")
			}
			if !assert.True(suite.T(), licenseSettings.IsValid, "Default Offline License is not valid or failed to set properly") {
				suite.T().FailNow()
			}
		} else if err != nil {
			suite.T().Fatalf("There was an error retrieving your license settings: %s", err)
		}

		upgradeVersion, upgradeChannel := os.Getenv("UPGRADE_TO_VERSION"), os.Getenv("UPGRADE_TO_CHANNEL")

		if upgradeVersion != "" {
			if upgradeChannel != "" {
				var newSettings forms.Settings
				newSettings.ReleaseChannel = &upgradeChannel
				if err := newFramework.API.SetHTTPSettings(&newSettings); err != nil {
					suite.T().Fatalf("Failed to set release channel before upgrading: %s", err)
				}

				if err := newUtil.PollAvailable(); err != nil {
					suite.T().Fatalf("DTR failed to restart in a reasonable amount of time after updating release channel")
				}
				suite.T().Logf("Upgrading to version: %s (%s) prior to running tests", upgradeVersion, upgradeChannel)
			} else {
				suite.T().Logf("Upgrading to version: %s prior to running tests", upgradeVersion)
			}
			if err := newFramework.API.Upgrade(upgradeVersion); err != nil {
				suite.T().Fatal("Failed to initiate upgrade prior to running tests")
			}
			// TODO(bbland): smarter polling rather than this long fixed wait
			time.Sleep(2 * time.Minute)
			if err := newUtil.PollAvailable(); err != nil {
				suite.T().Fatalf("DTR not available after upgrade attempt: %s", err)
			}
			version, err := newFramework.API.Version()
			if err != nil {
				suite.T().Fatalf("Failed to retrieve new DTR version: %s", err)
			}
			// There's currently no reasonable check for the "latest" version, as it will resolve to a more specific version
			if upgradeVersion != "latest" {
				if !assert.Equal(suite.T(), upgradeVersion, version, "DTR failed to upgrade prior to test") {
					suite.T().FailNow()
				}
			}
		}
		globalFramework = newFramework

		util.AppendDTRIgnorableLoggedErrors(util.KnownDTRInstallLogErrors)
		util.AppendDockerIgnorableLoggedErrors(util.KnownDockerInstallLogErrors)
		// we need to create an instance of u to run TestLogs
		u := util.MakeUtil(suite, newFramework)
		u.TestLogs()
	})

	if globalFramework == nil {
		suite.T().Fatal("Failed to set up the integration framework")
	}

	u := util.MakeUtil(suite, globalFramework)
	if u.IsSuiteRunningInLDAPMode() {
		u.RestartLDAPContainer()
	}

	return globalFramework, u
}
