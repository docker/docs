// This test verifies https://docs.docker.com/accounts/create-account/#create-a-docker-id

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test("dockerAccountCreateID", async ({ page }) => {
  // Go to the Docker WWW and select Sign In
  await page.goto("https://www.docker.com/");
  await page.getByRole("link", { name: "Sign In" }).click();

  // Select Sign Up
  await page.getByRole("link", { name: "Sign Up" }).click();

  // Enter email, username, and password
  await page.getByLabel("Email").click();
  await page.getByLabel("Email").fill("test4328798473298@customer.com");

  await page.getByLabel("Username").click();
  await page.getByLabel("Username").fill("test934832094830");

  const passwordInput = page.getByLabel("Password", { exact: true });
  await passwordInput.waitFor({ state: "visible", timeout: 5000 });
  await passwordInput.click();
  await passwordInput.fill("password543878432");

  // Verify Sign up button is present
  await expect(page.getByRole("button", { name: "Sign up" })).toBeVisible();
});