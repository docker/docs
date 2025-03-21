// This test verifies https://docs.docker.com/subscription/change/#downgrade-your-subscription

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

test("adminConsoleDowngradeSubscription", async ({ page }) => {
  // Select Billing and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Select Action menu
  await page.getByRole("group").getByRole("button").click();

  // Verify Cancel subscription button is there
  await expect(
    page.getByRole("link", { name: "Cancel subscription" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});