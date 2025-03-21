// This test verifies https://docs.docker.com/admin/organization/manage-a-team/#delete-a-team

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");

  // Accept cookies if visible
  const acceptCookies = page.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible().catch(() => false)) {
    await acceptCookies.click();
    await expect(acceptCookies).toBeHidden();
  }
});

test("adminConsoleDeleteTeam", async ({ page }) => {
  // Select Admin Console and choose organization
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Teams
  const teamsMenuItem = page.getByRole("menuitem", { name: "Teams" });
  await expect(teamsMenuItem).toBeVisible();
  await teamsMenuItem.click();

  // Select actions menu
  const teamActionsButton = page.getByTestId("team-actions-button-teamteam");
  await expect(teamActionsButton).toBeVisible();
  await teamActionsButton.click();

  // Verify the Delete button is present
  const deleteMenuItem = page.getByRole("menuitem", { name: "Delete" });
  await expect(deleteMenuItem).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});