// This test verifies https://docs.docker.com/subscription/scale/#add-docker-build-cloud-minutes

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

test("adminConsoleAddDBCMinutes", async ({ page }) => {
  // Select Billing and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Select View build minutes
  await page.getByRole("link", { name: "View build minutes" }).click();

  // Skip DBC pop-up
  await page.getByRole("button", { name: "Skip" }).click();

  // Select Purchase addiitonal minutes and choose amount
  await page.getByRole("link", { name: "Purchase additional minutes" }).click();
  await page.getByLabel("build minutes | $50 $40").check();

  // Select Continue to payment
  await page.getByRole("button", { name: "Continue to payment" }).click();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});