// This test verifies https://docs.docker.com/admin/organization/orgs/#view-an-organization

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubViewOrganizations", async ({ page }) => {
  // Navigate to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  const hubCard = page.getByTestId("dashboard-card-hub");
  await expect(hubCard).toBeVisible();
  await hubCard.click();
  const page2 = await page2Promise;
  await page2.waitForLoadState("load");

  // Wait for the organization menu item to be visible and click it
  await page2.getByTestId('layout-sidebar').getByLabel('open context switcher').click();
  await page2
    .getByTestId("layout-sidebar")
    .locator("a")
    .filter({ hasText: "sarahdatDocker Business" })
    .click();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
