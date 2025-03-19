import { defineConfig } from '@playwright/test';

export default defineConfig({
  workers: 1,
  globalSetup: 'tests/globalSetup.ts',
  use: {
    storageState: 'app-auth.json',
  },
});