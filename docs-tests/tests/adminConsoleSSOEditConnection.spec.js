// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/manage/#connect-an-organization
// and https://docs.docker.com/security/for-admins/single-sign-on/manage/#remove-an-organization
// and https://docs.docker.com/security/for-admins/single-sign-on/manage/#edit-a-connection

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

test("adminConsoleSSOEditConnection", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();
  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select SSO and SCIM
  await page
    .getByRole("menuitem", { name: "SSO and SCIM" })
    .scrollIntoViewIfNeeded();
  await page.getByRole("menuitem", { name: "SSO and SCIM" }).click();

  // Select action menu and then Edit connection
  await page
    .getByTestId("sso-connection-button-bd5e11ac-76b6-43bb-8046-4ef1a01bf587")
    .waitFor({ state: "visible", timeout: 30000 });
  await page
    .getByTestId("sso-connection-button-bd5e11ac-76b6-43bb-8046-4ef1a01bf587")
    .click();
  await page.getByRole("menuitem", { name: "Edit connection" }).click();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});