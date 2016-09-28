package ui

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

func TestLogin(t *testing.T, wdt selenium.WebDriverT, baseURL string, name, pass string, licensed bool) {
	log.Info("Starting TestLogin")
	wdt.Get(baseURL + "/#/login")

	username := wdt.FindElement(selenium.ByName, "username")
	username.SendKeys(name)

	password := wdt.FindElement(selenium.ByName, "password")
	password.SendKeys(pass)

	loginButton := wdt.FindElement(selenium.ById, "login-button")
	loginButton.Click()

	if !licensed {
		// Unlicensed and admin login should show license reminder screen
		wdt.FindElement(selenium.ByXPATH, "//p[text()='Your system is unlicensed.']")

		// Allow user to skip license upload
		skipUpload := wdt.FindElement(selenium.ByXPATH, "//a[text()='» Skip for now']")
		skipUpload.Click()
	}

	// Dashboard should be displayed
	wdt.FindElement(selenium.ByXPATH, "//span[text()='Dashboard']")
}
