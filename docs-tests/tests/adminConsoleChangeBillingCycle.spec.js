// This test verifies https://docs.docker.com/billing/cycle/#personal-account
// and https://docs.docker.com/billing/cycle/#organization

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

test("adminConsoleChangeBillingCycle", async ({ page }) => {
  // Select Billing Console and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "SS Sarah Sanders Docker Pro" })
    .click();

  // Verify switch to annual billing is present
  await expect(
    page.getByRole("link", { name: "Switch to annual billing" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});