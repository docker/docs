// This test verifies https://docs.docker.com/subscription/change/#upgrade-your-subscription

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubUpgradeSubscription", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  const hubCard = page.getByTestId("dashboard-card-hub");
  await expect(hubCard).toBeVisible();
  await hubCard.click();
  const page2 = await page2Promise;
  await page2.waitForLoadState("load");

  // Accept Cookies if present on the Docker Hub popup
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

  // Ensure the correct organization is selected in Billing
  const sidebarBilling = page3.getByTestId("layout-sidebar");
  const contextSwitcherBilling = sidebarBilling.getByLabel(
    "open context switcher"
  );
  await expect(contextSwitcherBilling).toBeVisible();
  await contextSwitcherBilling.click();

  const orgBilling = sidebarBilling
    .locator("a")
    .filter({ hasText: "vrunopayed" });
  await expect(orgBilling).toBeVisible();
  await orgBilling.click();

  // Wait for the "Upgrade" button to appear
  const upgradeButton = page3.getByRole("link", { name: "Upgrade" });
  await expect(upgradeButton).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
