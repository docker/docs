// This test verifies https://docs.docker.com/security/for-admins/hardened-desktop/settings-management/configure-admin-console/#create-a-settings-policy

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

test("adminConsoleCreateSettingsPolicy", async ({ page }) => {
  // Select Admin Console and choose organization
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await orgMenuItem.click();

  // Select Docker Desktop Settings Management
  await page.getByRole('menuitem', { name: 'Settings Management' }).click();

  // Select Create a settings policy
  const createPolicyLink = page.getByRole("link", {
    name: "Create a settings policy",
  });
  await expect(createPolicyLink).toBeVisible();
  await createPolicyLink.click();

  // Give your settings policy a name and description
  const policyNameField = page.getByLabel("Settings policy name");
  await expect(policyNameField).toBeVisible();
  await policyNameField.click();
  await policyNameField.fill("my policy");

  const descriptionField = page.getByLabel("Description (optional)");
  await expect(descriptionField).toBeVisible();
  await descriptionField.click();
  await descriptionField.fill("my description");

  // Assign the policy
  const allUsersCheckbox = page.getByLabel("All users");
  await expect(allUsersCheckbox).toBeVisible();
  await allUsersCheckbox.check();

  const selectedUsersCheckbox = page.getByLabel("Selected users");
  await expect(selectedUsersCheckbox).toBeVisible();
  await selectedUsersCheckbox.check();

  const userAssignmentTextField = page.getByTestId("user-assignment-textfield");
  await expect(userAssignmentTextField).toBeVisible();
  await userAssignmentTextField.click();

  const userOption = page
    .getByTestId("user-option-sarahstestaccount")
    .getByRole("checkbox");
  await expect(userOption).toBeVisible();
  await userOption.check();

  // Configure settings
  const usageStatsConfig = page
    .getByTestId("Send usage statistics")
    .getByLabel("Configuration");
  await expect(usageStatsConfig).toBeVisible();
  await usageStatsConfig.click();

  const alwaysEnabledOption = page.getByRole("option", {
    name: "Always enabled",
  });
  await expect(alwaysEnabledOption).toBeVisible();
  await alwaysEnabledOption.click();

  // Verify Create button exists
  const submitButton = page.getByTestId("submit-button");
  await expect(submitButton).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});