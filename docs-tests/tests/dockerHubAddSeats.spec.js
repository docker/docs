import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubAddSeats", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page2Promise = page.waitForEvent("popup");
  const dashboardHub = page.getByTestId("dashboard-card-hub");
  await expect(dashboardHub).toBeVisible({ timeout: 5000 });
  await dashboardHub.click();
  const page2 = await page2Promise;
  await page2.waitForLoadState("networkidle");

  // Select Billing
  const acceptCookiesButton = page2.getByRole("button", {
    name: "Accept All Cookies",
  });
  await expect(acceptCookiesButton).toBeVisible({ timeout: 5000 });
  await acceptCookiesButton.click();

  const layoutSidebarBilling = page2.getByTestId("layout-sidebar");
  await expect(layoutSidebarBilling).toBeVisible({ timeout: 5000 });
  const billingTab = layoutSidebarBilling.getByTestId("org-page-tab-billing");
  await expect(billingTab).toBeVisible({ timeout: 5000 });
  await billingTab.click();

  // Select Docker Billing and choose organization
  const page3Promise = page2.waitForEvent("popup");
  const dockerBillingLink = page2.getByRole("link", {
    name: "Docker Billing⁠",
  });
  await expect(dockerBillingLink).toBeVisible({ timeout: 5000 });
  await dockerBillingLink.click();
  const page3 = await page3Promise;
  await page3.waitForLoadState("networkidle");

  // Wait for the sidebar and user label to become visible in page3
  const layoutSidebar = page3.getByTestId("layout-sidebar");
  await expect(layoutSidebar).toBeVisible({ timeout: 5000 });

  const userLabel = layoutSidebar.getByLabel("open context switcher");
  await expect(userLabel).toBeVisible({ timeout: 5000 });
  await userLabel.click();

  const orgLink = layoutSidebar.locator("a").filter({ hasText: "sarahdat" });
  await expect(orgLink).toBeVisible({ timeout: 5000 });
  await orgLink.click();

  // Verify the "Add seats" button exists
  const addSeatsLink = page3.getByRole("link", { name: "Add seats" });
  await expect(addSeatsLink).toBeVisible({ timeout: 5000 });
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});