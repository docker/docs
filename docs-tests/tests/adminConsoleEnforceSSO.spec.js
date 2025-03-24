// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/connect/#optional-enforce-sso

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

test("adminConsoleEnforceSSO", async ({ page }) => {
  // Select Admin Console and choose organization
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

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

  // Select SSO connection action button
  const ssoConnectionButton = page.getByTestId(
    "sso-connection-button-bd5e11ac-76b6-43bb-8046-4ef1a01bf587"
  );
  await expect(ssoConnectionButton).toBeVisible();
  await ssoConnectionButton.click();

  // Verify Enable enforcement button is present
  const enableEnforcement = page.getByRole("menuitem", {
    name: "Enable enforcement",
  });
  await expect(enableEnforcement).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});