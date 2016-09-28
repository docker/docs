package ui

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

func TestDashboardApplications(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestDashboardApplications")
	wdt.Get(baseURL + "/#/dashboard")

	applications := wdt.FindElement(selenium.ByXPATH, "//a[@name='applications']")
	ClickWhenVisible(applications)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Applications']")
}

func TestDashboardContainers(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestDashboardContainers")
	wdt.Get(baseURL + "/#/dashboard")

	containers := wdt.FindElement(selenium.ByXPATH, "//a[@name='containers']")
	ClickWhenVisible(containers)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Containers']")
}

func TestDashboardImages(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestDashboardImages")
	wdt.Get(baseURL + "/#/dashboard")

	images := wdt.FindElement(selenium.ByXPATH, "//a[@name='images']")
	ClickWhenVisible(images)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Images']")
}

func TestDashboardNodes(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestDashboardNodes")
	wdt.Get(baseURL + "/#/dashboard")

	nodes := wdt.FindElement(selenium.ByXPATH, "//a[@name='nodes']")
	ClickWhenVisible(nodes)

	wdt.FindElement(selenium.ByXPATH, "//h3[text()='Nodes']")
	wdt.FindElement(selenium.ByXPATH, "//h3[text()='Cluster Controllers']")
}
