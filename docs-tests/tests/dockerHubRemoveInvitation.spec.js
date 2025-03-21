// This test verifies https://docs.docker.com/admin/organization/members/#remove-an-invitation

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubRemoveInvitation", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page2 = await page2Promise;

  // Select organizations, your org, then Members
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
  await page2
    .getByTestId("layout-sidebar")
    .getByTestId("org-page-tab-members")
    .click();

  // Select actions menu
  await page2
    .getByRole("row", { name: "sarahsanderstestinvite@docker" })
    .getByTestId("member-actions-menu-open")
    .click();

  // Verify Resend invitation button exists
  await expect(
    page2.getByRole("menuitem", { name: "Remove invitee" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});