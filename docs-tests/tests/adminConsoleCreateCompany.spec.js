// This test verifies https://docs.docker.com/admin/company/new-company/#create-a-company

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

test("adminConsoleCreateCompany", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Company management
  await page.getByRole("menuitem", { name: "Company management" }).click();

  // Select create company button and provide a company name
  await page.getByTestId("create-company-button").click();
  await page.getByLabel("Company name").click();
  await page.getByLabel("Company name").fill("sarahdatcompany");

  // Select Continue button
  await page.getByRole("button", { name: "Continue" }).click();

  // Verify billing flow
  await expect(
    page.getByRole("button", { name: "Create company" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});