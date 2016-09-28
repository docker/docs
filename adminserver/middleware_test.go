package adminserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpgradeMiddlewareWithUpdatesEnabled(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v0/admin/upgrade", nil)
	if err != nil {
		t.Fatal(err)
	}
	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Unexpected status code")
}

func TestUpgradeMiddlewareWithUpdatesDisabled(t *testing.T) {
	setup, deferFunc := buildTestingSetup(t)
	defer deferFunc()
	setup.ProvideAdminAuthentication()
	setup.mockSyslogWriter.On("Info", mock.Anything).Return(nil).Once()
	config := &hubconfig.UserHubConfig{
		DTRHost:         "gensokyo.jp",
		DisableUpgrades: true,
	}
	setup.adminServer.settingsStore.SetUserHubConfig(config)

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v0/admin/upgrade", nil)
	if err != nil {
		t.Fatal(err)
	}
	setup.handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code, "Unexpected status code")
}
