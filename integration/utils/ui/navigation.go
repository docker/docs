package ui

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

func TestSidebarNavToApplications(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToApplications")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Applications')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Applications')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/applications") {
		t.Fatalf("expected to be on applications page")
	}
}

func TestSidebarNavToContainers(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToContainers")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Containers')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Containers')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/containers") {
		t.Fatalf("expected to be on containers page")
	}
}

func TestSidebarNavToNodes(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToNodes")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Nodes')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Nodes')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/nodes") {
		t.Fatalf("expected to be on nodes page")
	}
}

func TestSidebarNavToVolumes(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToVolumes")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Volumes')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Volumes')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/volumes") {
		t.Fatalf("expected to be on volumes page")
	}
}

func TestSidebarNavToNetworks(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToNetworks")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Networks')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Networks')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/networks") {
		t.Fatalf("expected to be on networks page")
	}
}

func TestSidebarNavToImages(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToImages")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Images')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Images')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/images") {
		t.Fatalf("expected to be on images page")
	}
}

func TestSidebarNavToUsersTeams(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToUsersTeams")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/*/a[contains(.,'Users')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Users')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/users") {
		t.Fatalf("expected to be on users and teams page")
	}
}

func TestSidebarNavToSettings(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToSettings")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/*/a[contains(.,'Settings')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Settings')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/settings/logging") {
		t.Fatalf("expected to be on settings page")
	}
}

func TestSidebarNavToDashboard(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestSidebarNavToDashboard")
	wdt.Get(baseURL + "/#/containers")
	WaitForPageLoad(wdt)

	wdt.Q(".sidenavbutton").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='sidebar']/a[contains(.,'Dashboard')]").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Dashboard')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/dashboard") {
		t.Fatalf("expected to be on dashboard page")
	}
}

func TestNavToUserProfile(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestNavToUserProfile")
	wdt.Get(baseURL + "/#/dashboard")
	WaitForPageLoad(wdt)

	wdt.Q("#topnav .right.floated.dropdown.item").Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//*[@ui-sref='dashboard.user']").Click()
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[contains(text(), 'Profile')]")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/user") {
		t.Fatalf("expected to be on user profile page")
	}
}
