// This test verifies https://docs.docker.com/admin/company/organizations/#add-organizations-to-a-company

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

test("adminConsoleAddOrganizationsToCompany", async ({ page }) => {
  // Select Admin Console and choose company
  await page.getByTestId("dashboard-card-admin").click();
  await page.getByRole("menuitem", { name: "sarahscompany Company" }).click();

  // Select Add organization
  await page.getByRole("button", { name: "Add organization" }).click();

  // Choose organization to add from menu
  await page.getByLabel("Open", { exact: true }).click();
  await page.getByTestId("add-sarahdat-to-co-menu-item").click();

  // Verify Submit button present
  await expect(page.getByTestId("add-org-submit-button")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});