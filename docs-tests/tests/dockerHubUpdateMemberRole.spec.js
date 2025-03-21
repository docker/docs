// This test verifies https://docs.docker.com/admin/organization/members/#update-a-member-role

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubUpdateMemberRole", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select organizations, your org, then members
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

  // Select actions menu and Edit role
  await page2
    .getByRole("row", { name: "sarahstestaccount Guest sarah" })
    .getByTestId("member-actions-menu-open")
    .click();
  await page2.getByRole("menuitem", { name: "Edit role" }).click();

  // Choose new role
  await page2.getByLabel("Full administrative access to").check();

  // Verify Save button exists
  await expect(page2.getByRole("button", { name: "Save" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});