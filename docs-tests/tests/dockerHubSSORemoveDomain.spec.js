// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/manage/#remove-a-domain-from-an-sso-connection

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubSSORemoveDomain", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select organizations, your org, then Security
  await page2.getByRole("button", { name: "Accept All Cookies" }).click();
  await page2
    .getByTestId("layout-sidebar")
    .getByLabel("open context switcher")
    .click();
  await page2
    .getByTestId("layout-sidebar")
    .locator("a")
    .filter({ hasText: "sarahdatDocker Business" })
    .click();
  await page2.getByRole("menuitem", { name: "Security" }).click();

  // Verify Remove domain button exists
  await expect(page2.getByLabel("Delete testverifydomain.com")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});