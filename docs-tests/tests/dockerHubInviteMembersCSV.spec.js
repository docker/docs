// This test verifies https://docs.docker.com/admin/organization/members/#invite-members-via-csv-file

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubInviteMembersCSV", async ({ page }) => {
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

  // Select Invite members
  await page2.getByTestId("openAddMemberModal").click();

  // Select Invite by CSV
  await page2.getByLabel("Invite members via a CSV file").click();

  // Select Browse files and select a role
  await page2.getByRole("button", { name: "Browse files" }).click();
  await page2.getByLabel("Select a role").click();
  await page2.getByTestId("select-role-option-member").click();

  // Verify Review button exists
  await expect(page2.getByLabel("Select a team and pick a file")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});