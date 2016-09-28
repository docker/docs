package adminserver

import (
	"html/template"
	"net/http"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/util"
	"github.com/docker/dhe-deploy/hubconfig"

	"github.com/docker/dhe-license-server/tiers"
	"github.com/gorilla/context"
)

const (
	alertContextKey = "alerts"
	gcInProgress    = "GC_IN_PROGRESS"
	noLicense       = "NO_LICENSE"
)

type alert struct {
	ID          string       `json:"id,omitempty"`
	Class       string       `json:"class,omitempty"`
	Message     string       `json:"message"`
	Img         string       `json:"img,omitempty"`
	URL         template.URL `json:"url,omitempty"`
	CloseAction template.URL `json:"onclose,omitempty"`
}

type alerts struct {
	storageDir     string
	kvStore        hubconfig.KeyValueStore
	settingsStore  hubconfig.SettingsStore
	licenseChecker hubconfig.LicenseChecker
	alertCache     map[string]alert
}

func newAlert(message, class, url string) alert {
	tURL := template.URL(url)
	return alert{
		Message: message, Class: class, URL: tURL,
	}
}

func (a *alerts) addContextAlert(r *http.Request, data alert) {
	context.Set(r, alertContextKey, append(a.getAlertsFromContext(r), data))
}

func (a *alerts) getAlertsFromContext(r *http.Request) []alert {
	if al := context.Get(r, alertContextKey); al != nil {
		return al.([]alert)
	}
	return []alert{}
}

func (a *alerts) globalAlerts(request *http.Request) []alert {
	result := a.getAlertsFromContext(request)

	keys, _ := a.kvStore.List(deploy.RegistryROStatePath)
	if len(keys) > 0 {
		result = append(result, alert{
			ID:      gcInProgress,
			Class:   "warning",
			Message: "Deleting images... No one can push images while we clean your storage. Go to Garbage Collection.",
			Img:     "/public/img/broom-white.png",
			URL:     "/admin/settings/gc",
		})
	}

	user := util.GetAuthenticatedUser(request)
	if user != nil && !user.IsAnonymous && !*user.Account.IsAdmin {
		return result
	}

	if !a.licenseChecker.LicensingEnforced() {
		result = append(result, alert{Class: "wat", Message: "LICENSE ENFORCING IS DISABLED"})
	}

	if !a.licenseChecker.IsValid() {
		if a.licenseChecker.LicenseType() == tiers.Hourly {
			result = append(result, alert{
				ID:      noLicense,
				Class:   "alert",
				Message: "Warning: Unlicensed copy. Please contact support if this problem persists.",
				URL:     "/admin/support",
			})
		} else if a.licenseChecker.IsExpired() {
			//TODO distinguish between Offline and Online licenses?
			result = append(result, alert{
				ID:      noLicense,
				Class:   "alert",
				Message: "Your license is expired. Please upload a new license in the Settings page.",
				URL:     "/admin/settings/general",
			})
		} else {
			result = append(result, alert{
				ID:      noLicense,
				Class:   "alert",
				Message: "Warning: Unlicensed copy. Please register your license on the Settings page.",
				URL:     "/admin/settings/general",
			})
		}
	}

	return result
}
