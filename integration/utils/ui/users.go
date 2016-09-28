package ui

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	UCPTestUsername      = "ba_baracus"
	UCPTestPassword      = "secret123"
	UCPTestUserFirstName = "Mr"
	UCPTestUserLastName  = "T"
	UCPTestTeamName      = "the-a-team"
)

func TestCreateUser(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestCreateUser")
	wdt.Get(baseURL + "/#/users")

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	createUser := wdt.FindElement(selenium.ByXPATH, "//div[contains(@ng-click,'showCreateUser')]/i")
	ClickWhenVisible(createUser)

	WaitUntilVisible(wdt.FindElement(selenium.ById, "create-user-modal"))
	username := wdt.FindElement(selenium.ByName, "username")
	WaitForTransition()
	username.SendKeys(UCPTestUsername)

	password := wdt.FindElement(selenium.ByName, "password")
	WaitForTransition()
	password.SendKeys(UCPTestPassword)

	firstName := wdt.FindElement(selenium.ByName, "firstName")
	WaitForTransition()
	firstName.SendKeys(UCPTestUserFirstName)

	lastName := wdt.FindElement(selenium.ByName, "lastName")
	WaitForTransition()
	lastName.SendKeys(UCPTestUserLastName)

	submit := wdt.Q("#create-user-modal .ui.positive.button")
	submit.Click()

	userFilter := wdt.FindElement(selenium.ByCSSSelector, "#content div.right.aligned.floated.column input")
	userFilter.SendKeys(UCPTestUsername)

	result := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td[1]")
	if result.Text() != UCPTestUsername {
		t.Fatalf("expected to see %s in users list", UCPTestUsername)
	}
}

func TestCreateTeam(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestCreateTeam")
	wdt.Get(baseURL + "/#/users")
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	createTeam := wdt.FindElement(selenium.ByXPATH, "//a[contains(@ng-click, 'showCreateTeam')]")
	createTeam.Click()

	WaitForTransition()
	WaitUntilVisible(wdt.FindElement(selenium.ById, "create-team-modal"))
	teamName := wdt.FindElement(selenium.ByName, "teamName")
	teamName.SendKeys(UCPTestTeamName)

	submit := wdt.Q("#create-team-modal .ui.positive.button")
	submit.Click()

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Team: "+UCPTestTeamName+"']")
}

func TestAddUserToTeam(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestAddUserToTeam")
	wdt.Get(baseURL + "/#/users")
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	team := wdt.FindElement(selenium.ByXPATH, "//a[contains(., '"+UCPTestTeamName+"')]")
	team.Click()

	addUser := wdt.FindElement(selenium.ByXPATH, "//*[contains(@ng-click, 'showAddUser')]")
	addUser.Click()
	WaitForTransition()

	WaitUntilVisible(wdt.FindElement(selenium.ById, "add-user-to-team-modal"))
	filter := wdt.Q("#add-user-to-team-modal input")
	filter.SendKeys(UCPTestUsername)
	username := wdt.Q("#add-user-to-team-modal table > tbody > tr > td:nth-child(1)")
	if username.Text() != UCPTestUsername {
		t.Fatalf("expected to see %s in users list", UCPTestUsername)
	}
	addToTeam := wdt.Q("#add-user-to-team-modal table > tbody > tr > td:nth-child(4) div.button")
	addToTeam.Click()
	done := wdt.Q("#add-user-to-team-modal .ui.positive.button")
	done.Click()

	userFilter := wdt.FindElement(selenium.ByCSSSelector, "#content div.right.aligned.floated.column input")
	userFilter.SendKeys(UCPTestUsername)

	WaitForTransition()
	WaitForPageLoad(wdt)

	// TODO: Super ugly, need to add better IDs here to the HTML
	username = wdt.Q("#accountsBase-container > div > div.ng-scope > div.ui.grid.ng-scope > div:nth-child(4) > div > table > tbody > tr > td:nth-child(1)")
	if username.Text() != UCPTestUsername {
		t.Fatalf("expected to see %s in users list", UCPTestUsername)
	}
}

func TestRemoveUserFromTeam(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestRemoveUserFromTeam")

	username := wdt.Q("#accountsBase-container > div > div.ng-scope > div.ui.grid.ng-scope > div:nth-child(4) > div > table > tbody > tr > td:nth-child(1)")
	if username.Text() != UCPTestUsername {
		t.Fatalf("expected to see %s in users list", UCPTestUsername)
	}

	removeButton := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td/div[contains(.,'Remove')]")
	removeButton.Click()

	// Confirm image removal in modal
	WaitForTransition()

	WaitUntilVisible(wdt.Q("#remove-from-team-modal"))
	remove := wdt.FindElement(selenium.ByCSSSelector, "#remove-from-team-modal div.ui.positive.button")
	ClickWhenVisible(remove)
}

func TestRemoveTeam(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestRemoveTeam")
	wdt.Get(baseURL + "/#/users")
	WaitForPageLoad(wdt)

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	team := wdt.FindElement(selenium.ByXPATH, "//a[contains(., '"+UCPTestTeamName+"')]")
	team.Click()
	WaitForPageLoad(wdt)

	settings := wdt.Q("#accountsBase-container div div.ui.secondary.pointing.compact.menu > a:nth-child(3)")
	settings.Click()
	WaitForPageLoad(wdt)

	deleteTeam := wdt.FindElement(selenium.ByXPATH, "//div[contains(@ng-click,'showDeleteTeam')]")
	deleteTeam.Click()

	WaitForTransition()
	WaitUntilVisible(wdt.Q("#delete-team-modal"))
	remove := wdt.FindElement(selenium.ByCSSSelector, "#delete-team-modal div.ui.positive.button")
	remove.Click()
	WaitForTransition()

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	if !strings.HasSuffix(wdt.CurrentURL(), "/#/users") {
		t.Fatalf("expected to be on users page after removing a team")
	}
}

func TestRemoveUser(t *testing.T, wdt selenium.WebDriverT, baseURL string) {
	log.Info("Starting TestRemoveUser")
	wdt.Get(baseURL + "/#/users")

	wdt.FindElement(selenium.ByXPATH, "//span[text()='Users & Teams']")
	userFilter := wdt.FindElement(selenium.ByCSSSelector, "#content div.right.aligned.floated.column input")
	userFilter.SendKeys(UCPTestUsername)

	username := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td[1]")
	if username.Text() != UCPTestUsername {
		t.Fatalf("expected to see %s in users list", UCPTestUsername)
	}

	removeButton := wdt.FindElement(selenium.ByXPATH, "//table/tbody/tr/td/div/div/i[@class='remove icon']")
	ClickWhenVisible(removeButton)
	WaitForTransition()

	// Confirm image removal in modal
	WaitUntilVisible(wdt.Q("#remove-modal"))
	remove := wdt.FindElement(selenium.ByCSSSelector, "#remove-modal div.ui.positive.button")
	ClickWhenVisible(remove)
}
