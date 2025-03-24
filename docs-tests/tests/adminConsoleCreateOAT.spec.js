// This test verifies https://docs.docker.com/security/for-admins/access-tokens/#create-an-organization-access-token

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");

  const acceptCookies = page.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible({ timeout: 5000 }).catch(() => false)) {
    await acceptCookies.click();
  }
});

test("adminConsoleCreateOAT", async ({ page }) => {
  // Select Admin Console
  await page.getByTestId("dashboard-card-admin").click();

  // Wait for the organization menu item to be visible
  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Wait for and click the "Access tokens" menu item
  const accessTokensMenu = page.getByRole("menuitem", {
    name: /Access tokens/,
  });
  await expect(accessTokensMenu).toBeVisible();
  await accessTokensMenu.click();

  // Verify that the "Generate access token" button is visible
  const generateTokenLink = page.getByRole("link", {
    name: "Generate access token",
  });
  await expect(generateTokenLink).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});