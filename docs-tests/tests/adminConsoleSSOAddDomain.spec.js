// This test verifies https://docs.docker.com/security/for-admins/single-sign-on/configure/#step-one-add-your-domain

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

test("adminConsoleSSOAddDomain", async ({ page }) => {
  // Select Admin Console
  const adminCard = page.getByTestId("dashboard-card-admin");
  await expect(adminCard).toBeVisible();
  await adminCard.click();

  // Wait for and click the organization menu item
  const orgMenuItem = page.getByRole("menuitem", {
    name: "docs dat Docker Business",
  });
  await expect(orgMenuItem).toBeVisible();
  await Promise.all([
    page.waitForNavigation({ waitUntil: "networkidle" }),
    orgMenuItem.click(),
  ]);

  // Select Domain management
  const domainManagement = page.getByRole("menuitem", {
    name: "Domain management",
  });
  await expect(domainManagement).toBeVisible();
  await domainManagement.click();

  // Select Add a domain
  const addDomainButton = page.getByRole("button", { name: "Add a domain" });
  await expect(addDomainButton).toBeVisible();
  await addDomainButton.click();

  // Enter domain in text box and verify the Add domain submit button exists
  const domainInput = page.getByLabel("Domain", { exact: true });
  await expect(domainInput).toBeVisible();
  await domainInput.fill("mydomain.com");

  await expect(page.getByTestId("add-domain-submit-button")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});