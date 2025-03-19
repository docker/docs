// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/troubleshoot/#view-sso-and-scim-error-logs

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

test("adminConsoleSSOViewErrorLogs", async ({ page }) => {
  // Select Admin Console
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  // Select organization and wait for the navigation to complete
  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible({ timeout: 10000 });
  await Promise.all([
    page.waitForNavigation({ waitUntil: "networkidle" }),
    orgMenuItem.click(),
  ]);

  // Wait for and select "SSO and SCIM"
  const ssoScimMenuItem = page.getByRole("menuitem", { name: "SSO and SCIM" });
  await expect(ssoScimMenuItem).toBeVisible({ timeout: 10000 });
  await ssoScimMenuItem.click();

  // Open the action menu and click "View error logs"
  const ssoConnectionButton = page.getByTestId(
    "sso-connection-button-bd5e11ac-76b6-43bb-8046-4ef1a01bf587"
  );
  await expect(ssoConnectionButton).toBeVisible({ timeout: 10000 });
  await ssoConnectionButton.click();

  const viewErrorLogsItem = page.getByRole("menuitem", {
    name: "View error logs",
  });
  await expect(viewErrorLogsItem).toBeVisible({ timeout: 10000 });
  await viewErrorLogsItem.click();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});