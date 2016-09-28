package adminserver

import (
	"html/template"
	"net/http"
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/memory"
	"github.com/docker/dhe-deploy/licensing"

	"gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type AlertsSuite struct {
	req           *http.Request
	alerts        *alerts
	settingsStore hubconfig.SettingsStore
	mockChecker   *licensing.MockChecker
}

var _ = check.Suite(&AlertsSuite{})

func (s *AlertsSuite) SetUpSuite(c *check.C) {
	settingsStore := memory.NewSettingsStore()
	s.settingsStore = settingsStore
}

func (s *AlertsSuite) SetUpTest(c *check.C) {
	req, err := http.NewRequest("GET", "http://gensokyo.jp/admin", nil)
	c.Assert(err, check.IsNil)
	s.req = req
	s.mockChecker = new(licensing.MockChecker)
	s.alerts = &alerts{
		storageDir:     c.MkDir(),
		settingsStore:  s.settingsStore,
		licenseChecker: s.mockChecker,
		kvStore:        &hubconfig.MockKeyValueStore{},
	}
}

func (s *AlertsSuite) TestnewAlert(c *check.C) {
	var (
		alertMsg   = "sakuya izayoi"
		alertClass = "alert"
		alertURL   = "/marisa/kirisame"
	)
	alert := newAlert(alertMsg, alertClass, alertURL)
	c.Assert(alert.Message, check.Equals, alertMsg)
	c.Assert(alert.Class, check.Equals, alertClass)
	c.Assert(alert.URL, check.Equals, template.URL(alertURL))
}

func (s *AlertsSuite) TestContextAlerts(c *check.C) {
	alert := mockAlert("")
	s.alerts.addContextAlert(s.req, alert)
	data := s.alerts.getAlertsFromContext(s.req)
	c.Assert(len(data), check.Equals, 1)
	c.Assert(data[0], check.Equals, alert)

	s.alerts.addContextAlert(s.req, alert)
	data = s.alerts.getAlertsFromContext(s.req)
	c.Assert(len(data), check.Equals, 2)
}

func (s *AlertsSuite) TestEmptyContextAlerts(c *check.C) {
	data := s.alerts.getAlertsFromContext(s.req)
	c.Assert(data, check.DeepEquals, []alert{})
}

func (s *AlertsSuite) TestGlobalAlerts(c *check.C) {
	s.mockChecker.On("LicensingEnforced").Return(true)
	s.mockChecker.On("IsValid").Return(true)
	s.mockChecker.On("LicenseTier").Return("")
	s.mockChecker.On("IsExpired").Return(false)
	data := s.alerts.globalAlerts(s.req)
	c.Assert(data, check.DeepEquals, []alert{})

	var (
		alert2 = mockAlert("subterranean animism")
	)

	s.alerts.addContextAlert(s.req, alert2)

	data = s.alerts.globalAlerts(s.req)
	c.Assert(len(data), check.Equals, 1)
	c.Assert(isIn(alert2, data), check.Equals, true)
}

func isIn(item alert, slice []alert) bool {
	for _, a := range slice {
		if item == a {
			return true
		}
	}
	return false
}

func mockAlert(message string) alert {
	if message == "" {
		message = "sakuya izayoi"
	}
	return newAlert(message, "alert", "/marisa/kirisame")
}
