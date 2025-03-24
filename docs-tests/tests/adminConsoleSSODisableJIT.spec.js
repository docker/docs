// This test verifies https://docs.docker.com/security/for-admins/provisioning/just-in-time/#disable-jit-provisioning

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

test("adminConsoleSSODisableJIT", async ({ page }) => {
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

  // Select action menu and select Verify Disable JIT provisioning is present
  await page
    .getByTestId("sso-connection-button-d628fde4-179e-4dd9-bb03-6da663f6c108")
    .click();
  await expect(
    page.getByRole("menuitem", { name: "Disable JIT provisioning" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});