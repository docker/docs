// This test verifies https://docs.docker.com/billing/history/#view-renewal-date

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubDowngradeSubscription", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;
  await page2.waitForLoadState("load");

  // Accept cookies if present
  const acceptCookies = page2.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible({ timeout: 5000 }).catch(() => false)) {
    await acceptCookies.click();
  }

  // Select organization
  await page2
    .getByTestId("layout-sidebar")
    .getByLabel("open context switcher")
    .click();
  await page2
    .getByTestId("layout-sidebar")
    .locator("a")
    .filter({ hasText: "sarahdatDocker Business" })
    .click();

  // Navigate to Billing tab
  const billingTab = page2
    .getByTestId("layout-sidebar")
    .getByTestId("org-page-tab-billing");
  await billingTab.waitFor({ state: "visible", timeout: 10000 });
  await billingTab.click();

  // Open Docker Billing page
  const page3Promise = page2.waitForEvent("popup");
  await page2.getByRole("link", { name: "Docker Billing⁠" }).click();
  const page3 = await page3Promise;
  await page3.waitForLoadState("load");

  // Ensure correct organization selection in Billing
  await page3
    .getByTestId("layout-sidebar")
    .getByLabel("open context switcher")
    .click();
  await page3
    .getByTestId("layout-sidebar")
    .locator("a")
    .filter({ hasText: "sarahdat" })
    .click();

  // Wait for the actions menu button to appear
  const actionsMenu = page3.getByRole("group").getByRole("button");
  await actionsMenu.waitFor({ state: "visible", timeout: 10000 });
  console.log(await actionsMenu.isVisible());

  // Click actions menu
  await actionsMenu.click();

  // Verify Cancel Subscription button is present
  await expect(
    page3.getByRole("link", { name: "Cancel subscription" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});