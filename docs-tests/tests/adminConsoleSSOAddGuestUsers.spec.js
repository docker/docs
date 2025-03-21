// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/manage/#add-guest-users-when-sso-is-enabled

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

test("adminConsoleSSOAddGuestUsers", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Domain management
  await page.getByRole("menuitem", { name: "Domain management" }).click();

  // Select Add a domain
  await page.getByRole("button", { name: "Add a domain" }).click();

  // Enter domain in text box and verify Add domain button exists
  await page.getByLabel("Domain", { exact: true }).click();
  await page.getByLabel("Domain", { exact: true }).fill("mydomain.com");
  await expect(page.getByTestId("add-domain-submit-button")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});