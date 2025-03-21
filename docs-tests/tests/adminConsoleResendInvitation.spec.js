// This test verifies https://docs.docker.com/admin/organization/members/#resend-an-invitation

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

test("adminConsoleResendInvitation", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select actions button and select Resend
  await page.getByRole("menuitem", { name: "Members" }).click();
  await page
    .getByRole("row", { name: "-- -- sarahsanderstestinvite@" })
    .getByLabel("Member actions menu")
    .click();
  await page.getByRole("menuitem", { name: "Resend" }).click();

  // Verify Remove button is there
  await expect(page.getByRole("button", { name: "Invite" })).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});