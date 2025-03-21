// This test verifies https://docs.docker.com/admin/company/users/#invite-members-via-docker-id-or-email-address

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

test("adminConsoleInviteCompanyMembers", async ({ page }) => {
  // Select Admin Console and choose company
  const dashboardAdmin = page.getByTestId("dashboard-card-admin");
  await expect(dashboardAdmin).toBeVisible();
  await dashboardAdmin.click();

  const companyMenu = page.getByRole("menuitem", {
    name: "sarahscompany Company",
  });
  await expect(companyMenu).toBeVisible();
  await companyMenu.click();

  // Select Users
  const userManagementMenu = page.getByRole("menuitem", {
    name: "User management",
  });
  const usersMenuItem = page.getByRole("menuitem", { name: "Users" });
  if (!(await usersMenuItem.isVisible().catch(() => false))) {
    await userManagementMenu.click();
  } else await expect(usersMenuItem).toBeVisible();
  await usersMenuItem.click();

  // Select Invite and choose Emails or usernames
  const inviteButton = page.getByRole("button", { name: "Invite" });
  await expect(inviteButton).toBeVisible();
  await inviteButton.click();

  const inviteViaEmail = page.getByLabel("Invite members via email");
  await expect(inviteViaEmail).toBeVisible();
  await inviteViaEmail.click();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
