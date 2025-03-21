// This test verifies https://docs.docker.com/admin/organization/members/#add-a-member-to-a-team

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubAddTeamMember", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select organizations, your org, then Members
  await page2.getByRole("button", { name: "Accept All Cookies" }).click();
  await page2
    .getByTestId("layout-sidebar")
    .getByLabel("open context switcher")
    .click();
  await page2
    .getByTestId("layout-sidebar")
    .locator("a")
    .filter({ hasText: "sarahdatDocker Business" })
    .click();
  await page2
    .getByTestId("layout-sidebar")
    .getByTestId("org-page-tab-members")
    .click();

  // Select actions menu
  await page2
    .getByRole("row", { name: "sarahstestaccount Guest sarah" })
    .getByTestId("member-actions-menu-open")
    .click();

  // Select Add to team and choose team
  await page2.getByRole("menuitem", { name: "Add to team" }).click();
  await page2.getByLabel("Open", { exact: true }).click();
  await page2.getByRole("option", { name: "testeam" }).click();

  // Verify Add button exists
  await expect(page2.getByRole("button", { name: "Add" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});