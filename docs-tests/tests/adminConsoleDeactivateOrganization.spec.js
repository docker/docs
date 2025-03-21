// This test verifies https://docs.docker.com/admin/deactivate-account/#deactivate

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

test("adminConsoleDeactivateOrganization", async ({ page }) => {
  // Select avatar menu and choose Account settings
  await page.getByLabel("user menu sarahsanders720").click();
  await page.getByText("Account settings").click();

  // Verify Deactivate button is present
  await expect(
    page.getByRole("menuitem", { name: "Deactivate" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});