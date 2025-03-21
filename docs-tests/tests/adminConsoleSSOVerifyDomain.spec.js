// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/configure/#step-two-verify-your-domain

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

test("adminConsoleSSOVerifyDomain", async ({ page }) => {
  // Select Admin Console and choose organization
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Domain management
  const securityAndAccessMenu = page.getByRole("menuitem", {
    name: "Security and access",
  });
  const domainManagementMenu = page.getByRole("menuitem", {
    name: "Domain management",
  });
  if (!(await domainManagementMenu.isVisible().catch(() => false))) {
    await securityAndAccessMenu.click();
  } else await expect(domainManagementMenu).toBeVisible();
  await domainManagementMenu.click();

  // Verify the Verify button is present
  const verifyButton = page.getByTestId("verify-testverifydomain.com");
  await expect(verifyButton).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});