import { chromium, FullConfig } from '@playwright/test';
import fs from 'fs';
import { appLogin } from './loginHelper';

async function globalSetup(config: FullConfig) {
  const browser = await chromium.launch();

  console.log("Logging into Docker...");
  const appContext = await browser.newContext({ ignoreHTTPSErrors: true });
  const appPage = await appContext.newPage();
  await appLogin(appPage);

  try {
    await appPage.waitForURL('https://app-stage.docker.com/', { timeout: 60000 });
    console.log("App stage found, login confirmed.");
  } catch (error) {
    console.error("App stage not found. Login may not have completed successfully:", error);
  }

  try {
    await appPage.goto('https://app-stage.docker.com/', { waitUntil: 'networkidle', timeout: 60000 });
  } catch (error) {
    console.error("Error navigating to Docker Home after login:", error);
  }
  await appPage.waitForLoadState('networkidle');

  const storageState = await appContext.storageState();
  console.log("Original cookies:", storageState.cookies);

  const updatedCookies = storageState.cookies.map(cookie => {
    if (cookie.domain.includes('login-stage.docker.com')) {
      console.log(`Updating cookie ${cookie.name} domain from ${cookie.domain} to app-stage.docker.com`);
      return { ...cookie, domain: 'app-stage.docker.com' };
    }
    return cookie;
  });
  const newStorageState = { ...storageState, cookies: updatedCookies };

  fs.writeFileSync('app-auth.json', JSON.stringify(newStorageState, null, 2));
  console.log("Updated storage state saved to app-auth.json");

  await appContext.close();
  await browser.close();
}

export default globalSetup;