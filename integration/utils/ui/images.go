package ui

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	UCPTestImageName = "hello-world"
)

func TestPullImage(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestPullImage")
	wdt.Get(baseURL + "/#/images")

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Images']")

	pullImage := wdt.FindElement(selenium.ByXPATH, "//div[text()=' Pull Image']")
	ClickWhenVisible(pullImage)

	imageName := wdt.FindElement(selenium.ByName, "imageName")
	imageName.SendKeys(UCPTestImageName)

	pull := wdt.FindElement(selenium.ByCSSSelector, "#pull-modal div.ui.positive.button")
	ClickWhenVisible(pull)

	imageFilter := wdt.FindElement(selenium.ByCSSSelector, "#content div.right.aligned.floated.column input")
	imageFilter.SendKeys(UCPTestImageName)

	// After filtering, only element in table should be our test image
	result := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td[1]/div")
	if result.Text() != UCPTestImageName+":latest" {
		t.Fatalf("expected to see %s:latest in images list", UCPTestImageName)
	}
}

func TestRemoveImage(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestRemoveImage")

	imageName := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td[1]/div")
	if imageName.Text() != UCPTestImageName+":latest" {
		t.Fatalf("expected to see %s:latest in images list", UCPTestImageName)
	}

	ClickWhenVisibleByXPATH(wdt, "//table/tbody/tr/td/div/div/i[@class='remove icon']")

	// Confirm image removal in modal
	remove := wdt.FindElement(selenium.ByCSSSelector, "#remove-modal div.ui.positive.button")
	ClickWhenVisible(remove)
}
