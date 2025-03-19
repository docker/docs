// This test verifies https://docs.docker.com/accounts/create-account/#sign-up-with-google-or-github

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test("dockerAccountSignupGitHub", async ({ page }) => {
  // Go to the Docker WWW and select Sign In
  await page.goto("https://www.docker.com/");
  await page.getByRole("link", { name: "Sign In" }).click();

  // Select Sign Up
  await page.getByRole("link", { name: "Sign Up" }).click();

  // Verify Sign up with GitHub button is present
  await expect(
    page.getByRole("link", { name: "Continue with GitHub" })
  ).toBeVisible();
});