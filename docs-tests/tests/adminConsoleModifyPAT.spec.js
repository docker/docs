// This test verifies https://docs.docker.com/security/for-developers/access-tokens/#modify-existing-tokens

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

test("adminConsoleModifyPAT", async ({ page }) => {
  // Select avatar and choose Account settings
  const avatarButton = page.getByLabel("user menu sarahsanders720");
  await expect(avatarButton).toBeVisible();
  await avatarButton.click();

  const accountSettings = page.getByText("Account settings");
  await expect(accountSettings).toBeVisible();
  await accountSettings.click();

  // Select Personal access tokens
  await page.getByRole("menuitem", { name: "Personal access tokens" }).click();

  // Select actions menu and Edit
  await page
    .getByTestId("pat-menu-button-83f8d7e2-65a0-401e-b836-0362cddf794b")
    .click();
  await page.getByRole("menuitem", { name: "Edit" }).click();
  await page.getByLabel("Access token description").click();
  await page.getByLabel("Access token description").fill("newpat1");

  // Verify Generate button exists
  await expect(page.getByRole("button", { name: "Save token" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});