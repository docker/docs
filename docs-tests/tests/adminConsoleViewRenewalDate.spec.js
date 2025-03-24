// This test verifies https://docs.docker.com/billing/history/#view-renewal-date

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

test("adminConsoleViewRenewalDate", async ({ page }) => {
  // Select Billing Console and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Verify next bill date is present
  await expect(
    page
      .locator("div")
      .filter({ hasText: /^Next bill/ })
      .getByRole("paragraph")
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});