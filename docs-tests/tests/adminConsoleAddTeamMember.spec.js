// This test verifies https://docs.docker.com/admin/organization/members/#add-a-member-to-a-team

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");

  // Accept cookies if visible
  const acceptCookies = page.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible({ timeout: 5000 }).catch(() => false)) {
    await acceptCookies.click();
    await expect(acceptCookies).toBeHidden();
  }
});

test("adminConsoleAddTeamMember", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select team tab and Select team name
  await page.getByRole("menuitem", { name: "Teams" }).click();
  await page.getByRole("link", { name: "teamteam" }).click();

  // Select Add member
  await page.getByRole("button", { name: "Add member" }).click();

  // Select a member from the drop down list
  await page.getByLabel("Open", { exact: true }).click();
  await page.getByTestId("member-option-sarahsanders720").click();

  // Verify Add button is there
  await expect(page.getByRole("button", { name: "Add" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});