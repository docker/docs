// This test verifies https://docs.docker.com/admin/organization/manage-a-team/#configure-repository-permissions-for-a-team

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubConfigureRepoPermissions", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select organizations, your org, then Teams
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
    .getByTestId("org-page-tab-teams")
    .click();

  // Select a team
  await page2.getByRole("link", { name: "teamteam" }).click();

  // Select Permissions
  await page2.getByTestId("permissions").click();

  // Select repo and permission level
  await page2.getByLabel("Open", { exact: true }).click();
  await page2.getByRole("option", { name: "sarahdatrepo" }).click();
  await page2.getByLabel("Permission").click();
  await page2.getByRole("option", { name: "Admin" }).click();

  // Verify Remove button exists
  await expect(page2.getByRole("button", { name: "Add" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});