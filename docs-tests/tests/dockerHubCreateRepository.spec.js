// This test verifies https://docs.docker.com/get-started/introduction/build-and-push-first-image/#create-an-image-repository,
// https://docs.docker.com/get-started/docker-concepts/the-basics/what-is-a-registry/#create-your-first-repository
// and https://docs.docker.com/get-started/workshop/04_sharing_app/#create-a-repository

import { test, expect } from "@playwright/test";
test.use({ storageState: "app-auth.json" });

test.beforeEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});

test("dockerHubCreateRepository", async ({ page }) => {
  // Go to Docker Hub
  await page.goto("https://app-stage.docker.com/");
  const page4Promise = page.waitForEvent("popup");
  await page.getByTestId("dashboard-card-hub").click();
  const page4 = await page4Promise;

  // Select Create repository
  await page4.getByTestId("createRepoBtn").click();

  // Enter a repository name, description, and select Public visibility
  await page4.getByTestId("repoNameField-input").click();
  await page4
    .getByTestId("repoNameField-input")
    .fill("getting-started-todo-app");
  await page4.getByTestId("repoDescriptionField-input").click();
  await page4.getByTestId("repoDescriptionField-input").fill("description");
  await page4.getByLabel("PublicAppears in Docker Hub").check();

  // Verify Create button exists
  await expect(page4.getByTestId("submit")).toBeVisible();
});

test.afterEach(async ({ page }) => {
  await page.goto("https://app-stage.docker.com/");
  await page.waitForLoadState("networkidle");
});