// This test verifies https://docs.docker.com/subscription/manage-seats/#add-seats
// and https://docs.docker.com/subscription/manage-seats/#add-seats

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

test("adminConsoleAddSeats", async ({ page }) => {
  // Select Billing and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Select Add seats
  await page.getByRole("link", { name: "Add seats" }).click();

  // Specify number of seats to add and select Continue to billing
  await page.getByRole("textbox").fill("10");
  await page.goto(
    "https://app-stage.docker.com/billing/sarahdat/update/quantity/plan?add=10"
  );
  await page.getByRole("link", { name: "Continue to billing" }).click();

  // Select Continue to payment
  await page.getByRole("link", { name: "Continue to payment" }).click();

  // Verify Update subscription button is there
  await expect(
    page.getByRole("button", { name: "Update subscription" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});