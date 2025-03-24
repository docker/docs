// This test verifies https://docs.docker.com/admin/company/owners/#add-a-company-owner

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

test("adminConsoleAddCompanyOwner", async ({ page }) => {
  // Select Admin Console and choose company
  await page.getByTestId("dashboard-card-admin").click();
  await page.getByRole("menuitem", { name: "sarahscompany Company" }).click();

  // Select Company owners
  await page.getByRole("menuitem", { name: "Company owners" }).click();

  // Select Add owner and search for owner by Docker ID
  await page.getByRole("button", { name: "Add owner" }).click();
  await page.getByLabel("Search by Docker ID").click();
  await page.getByLabel("Search by Docker ID").fill("sarahstestaccount");

  // Verify Add company owner button exists
  await expect(
    page.getByRole("button", { name: "Add company owner" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});