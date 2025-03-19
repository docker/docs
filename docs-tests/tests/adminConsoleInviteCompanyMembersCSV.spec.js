// This test verifies https://docs.docker.com/admin/company/users/#invite-members-via-csv-file

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");

  const acceptCookies = page.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible().catch(() => false)) {
    await acceptCookies.click();
  }
});

test("adminConsoleInviteCompanyMembersCSV", async ({ page }) => {
  // Select Admin Console and choose company
  const adminConsole = page.getByTestId("dashboard-card-admin");
  await expect(adminConsole).toBeVisible();
  await adminConsole.click();

  const companyMenuItem = page.getByRole("menuitem", {
    name: "sarahscompany Company",
  });
  await expect(companyMenuItem).toBeVisible();
  await companyMenuItem.click();

  // Select Users
  const userManagementMenu = page.getByRole("menuitem", {
    name: "User management",
  });
  const usersMenuItem = page.getByRole("menuitem", { name: "Users" });
  if (!(await usersMenuItem.isVisible().catch(() => false))) {
    await userManagementMenu.click();
  } else
  await expect(usersMenuItem).toBeVisible();
  await usersMenuItem.click();

  // Select Invite and verify CSV upload is present
  const inviteButton = page.getByRole("button", { name: "Invite" });
  await expect(inviteButton).toBeVisible();
  await inviteButton.click();
  await expect(page.getByLabel("Invite members via a CSV file")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});
