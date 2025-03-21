// This test verifies https://docs.docker.com/admin/organization/members/#remove-an-invitation

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

test("adminConsoleRemoveInvitation", async ({ page }) => {
  // Select Admin Console and choose organization
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Members
  const membersMenuItem = page.getByRole("menuitem", { name: "Members" });
  await expect(membersMenuItem).toBeVisible();
  await membersMenuItem.click();

  const inviteRow = page.getByRole("row", {
    name: "-- -- sarahsanderstestinvite@",
  });
  await expect(inviteRow).toBeVisible();

  // Open the actions menu
  const actionsMenu = inviteRow.getByLabel("Member actions menu");
  await expect(actionsMenu).toBeVisible();
  await actionsMenu.click();

  // Click Remove in the actions dropdown
  const removeMenuItem = page.getByRole("menuitem", { name: "Remove" });
  await expect(removeMenuItem).toBeVisible();
  await removeMenuItem.click();

  // Verify that the Remove button is visible
  const removeButton = page.getByRole("button", { name: "Remove" });
  await expect(removeButton).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});