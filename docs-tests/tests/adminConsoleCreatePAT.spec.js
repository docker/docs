// This test verifies https://docs.docker.com/security/for-developers/access-tokens/#create-an-access-token

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

test("adminConsoleCreatePAT", async ({ page }) => {
  // Select avatar and choose Account settings
  await page.getByLabel('user menu sarahsanders720').click();

  const accountSettings = page.getByText("Account settings");
  await expect(accountSettings).toBeVisible();
  await accountSettings.click();

  // Select Personal access tokens and Generate new token
  await page.getByRole("menuitem", { name: "Personal access tokens" }).click();

  // Verify Generate button exists
  await expect(
    page.getByRole("link", { name: "Generate new token" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});