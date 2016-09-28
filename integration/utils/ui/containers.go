package ui

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	UCPTestContainerName = "selenium-ucp-test-container"
)

func TestContainerDeploy(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestContainerDeploy")
	wdt.Get(baseURL + "/#/deploy")
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//h3[text()='Basic Settings']")
	wdt.FindElement(selenium.ByXPATH, "//span[text()='Deploy']")

	image := wdt.FindElement(selenium.ByName, "image")
	image.SendKeys("busybox")

	containerName := wdt.FindElement(selenium.ByName, "containerName")
	containerName.SendKeys(UCPTestContainerName)

	containerConfig := wdt.FindElement(selenium.ById, "deploy-container")
	ClickWhenVisible(containerConfig)

	entrypoint := wdt.FindElement(selenium.ByName, "entrypoint")
	entrypoint.SendKeys("ash")

	run := wdt.FindElement(selenium.ByCSSSelector, "div.submit")
	ClickWhenVisible(run)

	wdt.FindElement(selenium.ByXPATH, "//td[text()='"+UCPTestContainerName+"']")
	wdt.FindElement(selenium.ByXPATH, "//ol[@id='breadcrumbs']/li/span[text()='Containers']")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/containers") {
		t.Fatalf("expected to be on containers page after deploy")
	}
}

func TestContainerRemove(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestContainerRemove")
	wdt.Get(baseURL + "/#/containers")
	WaitForPageLoad(wdt)

	container := wdt.FindElement(selenium.ByXPATH, "//td[text()='"+UCPTestContainerName+"']")
	container.Click()

	wdt.FindElement(selenium.ByXPATH, "//h2[text()='"+UCPTestContainerName+"']")
	destroy := wdt.FindElement(selenium.ByXPATH, "//button[text()=' Remove']")
	ClickWhenVisible(destroy)

	wdt.FindElement(selenium.ByXPATH, "//p[text()='Are you sure you want to remove this container?']")

	confirm := wdt.FindElement(selenium.ByXPATH, "//div[contains(@ng-click,'removeContainer')]")
	confirm.Click()

	wdt.FindElement(selenium.ByXPATH, "//*[@id='breadcrumbs']/li/span[text()='Containers']")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/containers") {
		t.Fatalf("expected to be on containers page after removing container")
	}
}
