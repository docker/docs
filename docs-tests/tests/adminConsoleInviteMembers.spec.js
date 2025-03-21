// This test verifies https://docs.docker.com/admin/organization/members/#invite-members-via-docker-id-or-email-address
// and https://docs.docker.com/security/for-admins/single-sign-on/manage/#add-guest-users-when-sso-is-enabled

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

test("adminConsoleInviteMembers", async ({ page }) => {
  // Select Admin Console and choose organization
  await page.getByTestId("dashboard-card-admin").click();
  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Click Invite
  await page.getByRole("button", { name: "Invite" }).click();

  // Choose to invite by emails or usernames
  await page.getByLabel("Invite members via email").click();

  // Enter username or email and select role
  await page.getByLabel("Enter username or email").click();
  await page.getByLabel("Enter username or email").fill("invitee@test.com");
  await page.getByLabel("Select a role").click();
  await page.getByRole("option", { name: "Member Non-administrative" }).click();

  // Ensure dropdown is attached and stable before clicking
  const teamSelect = page.getByTestId("create-team-select");
  await teamSelect.waitFor({ state: "visible", timeout: 10000 });
  await page.waitForTimeout(500);

  // Click the team select dropdown
  await teamSelect.click();

  // Assert invite button is visible
  await expect(page.getByTestId("invite-member-modal-invite")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});