// This test verifies https://docs.docker.com/security/for-admins/provisioning/scim/#enable-scim-in-docker

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

test("adminConsoleSSOEnableSCIM", async ({ page }) => {
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

  // Select action menu and Enable SCIM
  await page
    .getByTestId("sso-connection-button-bd5e11ac-76b6-43bb-8046-4ef1a01bf587")
    .click();
  const enableScimMenuItem = page.getByRole("menuitem", {
    name: "Enable SCIM",
  });
  await expect(enableScimMenuItem).toBeVisible();
  await enableScimMenuItem.click();

  // Verify Enable button is present
  await expect(page.getByRole("button", { name: "Enable" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
