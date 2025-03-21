import { Page } from '@playwright/test';

export async function appLogin(page: Page) {
  if (!process.env.ENV_APP_USERNAME || !process.env.ENV_APP_PASSWORD) {
    throw new Error("Missing required environment variables: ENV_APP_USERNAME or ENV_APP_PASSWORD");
  }

  await page.goto('https://app-stage.docker.com/');
  await page.getByLabel('Username or email address').waitFor({ timeout: 10000 });
  await page.getByLabel('Username or email address').click();
  await page.getByLabel('Username or email address').fill(process.env.ENV_APP_USERNAME);
  await page.getByRole('button', { name: 'Continue', exact: true }).click();
  await page.getByRole('link', { name: 'Sign in with Okta FastPass' }).click();
  await page.waitForTimeout(3000);
  await page.getByLabel('Password').click();
  await page.getByLabel('Password').fill(process.env.ENV_APP_PASSWORD);
  await page.getByRole('button', { name: 'Verify' }).click();
  await page.waitForLoadState('networkidle');
}