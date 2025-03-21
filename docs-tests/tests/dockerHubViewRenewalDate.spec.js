// This test verifies https://docs.docker.com/billing/history/#view-renewal-date

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubViewRenewalDate", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  const hubCard = page.getByTestId("dashboard-card-hub");
  await expect(hubCard).toBeVisible();
  await hubCard.click();
  const page2 = await page2Promise;
  await page2.waitForLoadState("load");

  // Accept cookies if the banner is present
  const acceptCookies = page2.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible({ timeout: 5000 }).catch(() => false)) {
    await acceptCookies.click();
    await expect(acceptCookies).toBeHidden();
  }

  // Navigate to Billing tab
  const billingTab = page2
    .getByTestId("layout-sidebar")
    .getByTestId("org-page-tab-billing");
  await expect(billingTab).toBeVisible();
  await billingTab.click();

  // Open Docker Billing page
  const page3Promise = page2.waitForEvent("popup");
  const billingLink = page2.getByRole("link", { name: "Docker Billing⁠" });
  await expect(billingLink).toBeVisible();
  await billingLink.click();
  const page3 = await page3Promise;
  await page3.waitForLoadState("load");

  // Verify the next bill date is visible
  const renewalDateLocator = page3
    .locator("div")
    .filter({ hasText: /^Next bill/ })
    .getByRole("paragraph");
  await expect(renewalDateLocator).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
