// This test verifies https://docs.docker.com/security/for-admins/hardened-desktop/registry-access-management/#configure-registry-access-management-permissions

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

test("adminConsoleConfigureRegistryAccessManagement", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Choose Registry access
  await page.getByRole("menuitem", { name: "Registry access" }).click();

  // Enable Registry access management
  const ramSwitch = page.getByTestId("ram-enabled-switch");
  if (!(await ramSwitch.isChecked())) {
    await ramSwitch.check();
  } else {
    console.log("Registry access management is already enabled.");
  }

  // Select Add registry and enter registry details
  await page.getByRole("button", { name: "Add registry" }).click();
  await page.getByLabel("Registry address *").click();
  await page.getByLabel("Registry address *").fill("docker");
  await page.getByLabel("Registry nickname *").click();
  await page.getByLabel("Registry nickname *").fill("test");

  // Verify Create button exists
  await expect(page.getByTestId("button-save")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});