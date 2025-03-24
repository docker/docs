// This test verifies https://docs.docker.com/security/for-admins/access-tokens/#modify-existing-tokens

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

test("adminConsoleModifyOAT", async ({ page }) => {
  // Select Admin Console
  await page.getByTestId("dashboard-card-admin").click();

  // Select organization
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Select Access tokens from left-nav
  await page
    .getByRole("menuitem", { name: "Access tokens New" })
    .scrollIntoViewIfNeeded();
  await page.getByRole("menuitem", { name: "Access tokens New" }).click();

  // Select actions button
  await page.getByRole('cell', { name: 'token-0-actions-menu' }).click();

  // Verify Deactivate button is there
  await expect(
    page.getByRole("menuitem", { name: "Deactivate" })
  ).toBeVisible();

  // Verify Edit button is there
  await expect(page.getByRole("menuitem", { name: "Edit" })).toBeVisible();

  // Verify Delete button is there
  await expect(page.getByRole("menuitem", { name: "Delete" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
