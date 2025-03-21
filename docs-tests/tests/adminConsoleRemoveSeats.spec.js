// This test verifies https://docs.docker.com/subscription/manage-seats/#remove-seats

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");

  // Accept cookies if visible
  const acceptCookies = page.getByRole("button", {
    name: "Accept All Cookies",
  });
  if (await acceptCookies.isVisible().catch(() => false)) {
    await acceptCookies.click();
    await expect(acceptCookies).toBeHidden();
  }
});

test("adminConsoleRemoveSeats", async ({ page }) => {
  // Select Billing Console and choose organization
  await page.getByTestId("dashboard-card-billing-account-center").click();
  await page
    .getByRole("menuitem", { name: "docs dat Docker Business" })
    .click();

  // Select action menu and choose Remove seats
  const addSeatsButton = page.getByTestId("add-seats-action-button");
  await expect(addSeatsButton).toBeVisible();
  await addSeatsButton.click();
  await page.getByTestId("remove-seats-button").click();

  // Select number of seats to remove
  const seatTextbox = page.getByRole("textbox");
  await expect(seatTextbox).toBeVisible();
  await seatTextbox.click();
  await seatTextbox.fill("1");
  await page.goto(
    "https://app-stage.docker.com/billing/sarahdat/update/quantity/plan?remove=1",
    { waitUntil: "networkidle", timeout: 60000 }
  );

  // Verify Update subscription button exists
  await expect(
    page.getByRole("button", { name: "Update subscription" })
  ).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});