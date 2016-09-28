package ui

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	UCPTestApplicationName = "seleniumtestapp"
	UCPTestComposeYml      = "redis:\\n    image: redis\\n"
)

func TestCreateApp(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestCreateApp")
	wdt.Get(baseURL + "/#/applications")
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Applications']")

	createApp := wdt.FindElement(selenium.ByXPATH, "//div[contains(@ng-click,'showComposeDialog')]")
	ClickWhenVisible(createApp)

	WaitUntilVisible(wdt.FindElement(selenium.ById, "compose-modal"))
	projectName := wdt.FindElement(selenium.ByName, "projectName")
	WaitForTransition()
	projectName.SendKeys(UCPTestApplicationName)

	wdt.ExecuteScript("var editor = $('.CodeMirror')[0].CodeMirror;editor.setValue('"+UCPTestComposeYml+"');", nil)

	submit := wdt.Q("#compose-modal .ui.submit.button")
	submit.Click()

	// Add a generous wait, since deploying an app could take a while depending on network connectivity
	// Since the implicit wait timeout is set to 10 seconds, each iteration of the loop will pause for a maximum
	// of 10 seconds, so we only need to loop 6 times to increase this to a 60 second wait
	wd := wdt.WebDriver()
	for i := 0; i < 6; i++ {
		if _, err := wd.FindElement(selenium.ByXPATH, "//pre[contains(.,'Successfully deployed')]"); err == nil {
			break
		}
	}

	wdt.FindElement(selenium.ByXPATH, "//pre[contains(.,'Successfully deployed')]")
	wdt.FindElement(selenium.ByXPATH, "//pre[contains(.,'Removed container')]")

	ClickWhenVisibleByXPATH(wdt, "//div[contains(@ng-click,'closeComposeDialog') and contains(., 'Done')]")
	WaitForPageLoad(wdt)

	userFilter := wdt.FindElement(selenium.ByXPATH, "//*[contains(@ng-model, 'vm.filter')]")
	userFilter.SendKeys(UCPTestApplicationName)
	WaitForTransition()

	result := wdt.Q("span.name")
	if result.Text() != UCPTestApplicationName {
		t.Fatalf("expected to see %s in applications list, found %s", UCPTestApplicationName, result.Text())
	}
}
