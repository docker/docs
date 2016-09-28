package ui

import (
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"sourcegraph.com/sourcegraph/go-selenium"

	"github.com/docker/orca/integration/utils"
)

const (
	TIMEOUT_SETTING     = uint(10000)
	TRANSITION_DURATION = 1 * time.Second
)

var (
	DefaultWindowSize = selenium.Size{Width: 1024.0, Height: 768.0}
)

func WaitForTransition() {
	time.Sleep(TRANSITION_DURATION)
}

func WaitForPageLoad(wdt selenium.WebDriverT) {
	dimmer := wdt.Q(".ui.dimmer")
	for i := 0; i < 10; i++ {
		if !dimmer.IsDisplayed() {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func WaitUntilVisible(elem selenium.WebElementT) {
	for i := 0; i < 10; i++ {
		if elem.IsDisplayed() {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ClickWhenVisible(elem selenium.WebElementT) {
	WaitUntilVisible(elem)
	elem.Click()
}

func ClickWhenVisibleByXPATH(wdt selenium.WebDriverT, xpath string) {
	for i := 0; i < 10; i++ {
		elem := wdt.FindElement(selenium.ByXPATH, xpath)
		if elem.IsDisplayed() {
			elem.Click()
			break
		}
		time.Sleep(1 * time.Second)

	}
}

func NewWebDriver(seleniumURL, browser string) (selenium.WebDriver, error) {
	// disable debug logging
	selenium.Log = nil

	caps := selenium.Capabilities(map[string]interface{}{
		"browserName":    browser,
		"platform":       "linux",
		"acceptSslCerts": true,
	})

	driver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Errorf("Failed to create selenium remote driver: %s", err)
		return nil, err
	}

	driver.SetAsyncScriptTimeout(TIMEOUT_SETTING)
	driver.SetImplicitWaitTimeout(TIMEOUT_SETTING)

	// resize window
	curWindow, err := driver.CurrentWindowHandle()
	if err != nil {
		log.Errorf("Failed to get current window handle: %s", err)
		return nil, err
	}
	if err := driver.ResizeWindow(curWindow, DefaultWindowSize); err != nil {
		log.Errorf("Failed to resize window: %s", err)
		return nil, err
	}

	return driver, nil
}

func TestUI(t *testing.T, serverURL string) {
	seleniumURL := os.Getenv("SELENIUM_URL")
	if seleniumURL == "" {
		t.Skip("Skipping UI tests - set SELENIUM_URL to enable")
		return
	}

	browsers := []string{"chrome", "firefox"}
	for _, browser := range browsers {
		drv, err := NewWebDriver(seleniumURL, browser)
		require.Nil(t, err)
		defer drv.Quit()
		wdt := drv.T(t)

		log.Infof("[%s] UI tests started", browser)

		TestLogin(t, wdt, serverURL, "admin", utils.GetAdminPassword(), false)

		// Applications
		TestCreateApp(t, wdt, serverURL)

		// Users & Teams
		TestCreateUser(t, wdt, serverURL)
		TestCreateTeam(t, wdt, serverURL)
		TestAddUserToTeam(t, wdt, serverURL)
		TestRemoveUserFromTeam(t, wdt, serverURL)
		TestRemoveTeam(t, wdt, serverURL)
		TestRemoveUser(t, wdt, serverURL)

		// Navigation Tests
		TestSidebarNavToDashboard(t, wdt, serverURL)
		TestSidebarNavToApplications(t, wdt, serverURL)
		TestSidebarNavToContainers(t, wdt, serverURL)
		TestSidebarNavToNodes(t, wdt, serverURL)
		TestSidebarNavToImages(t, wdt, serverURL)
		TestSidebarNavToSettings(t, wdt, serverURL)
		TestSidebarNavToVolumes(t, wdt, serverURL)
		TestSidebarNavToNetworks(t, wdt, serverURL)
		TestSidebarNavToUsersTeams(t, wdt, serverURL)

		TestNavToUserProfile(t, wdt, serverURL)

		// Containers
		TestContainerDeploy(t, wdt, serverURL)
		TestContainerRemove(t, wdt, serverURL)

		// Dashboard
		TestDashboardApplications(t, wdt, serverURL)
		TestDashboardContainers(t, wdt, serverURL)
		TestDashboardNodes(t, wdt, serverURL)

		// Images
		TestPullImage(t, wdt, serverURL)
		TestRemoveImage(t, wdt, serverURL)

		log.Infof("[%s] UI tests finished", browser)
	}
}
