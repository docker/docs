// This test verifies https://docs.docker.com/billing/cycle/#personal-account
// and https://docs.docker.com/billing/cycle/#organization

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubChangeBillingCycle", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select Billing
  await page2.getByRole("button", { name: "Accept All Cookies" }).click();
  await page2
    .getByTestId("layout-sidebar")
    .getByTestId("org-page-tab-billing")
    .click();

  // Select Docker Billing
  const page3Promise = page2.waitForEvent("popup");
  await page2.getByRole("link", { name: "Docker Billing⁠" }).click();
  const page3 = await page3Promise;

  // Verify switch to annual billing is present
  await expect(
    page3.getByRole("link", { name: "Switch to annual billing" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});