import { defineConfig } from "@playwright/test";
import path from "node:path";

export default defineConfig({
  testDir: ".",
  testMatch: "browser.spec.mjs",
  fullyParallel: true,
  workers: 2,
  forbidOnly: Boolean(process.env.CI),
  outputDir: "../../tmp/test-results",
  use: {
    baseURL: process.env.DOCS_TEST_URL || "http://127.0.0.1:1313",
    reducedMotion: "reduce",
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    launchOptions: {
      executablePath: process.env.CHROMIUM_PATH || undefined,
    },
  },
  projects: [
    {
      name: "desktop",
      use: { viewport: { width: 1440, height: 900 }, colorScheme: "light" },
    },
    {
      name: "mobile",
      use: { viewport: { width: 390, height: 900 }, colorScheme: "dark" },
    },
  ],
  webServer: process.env.DOCS_TEST_URL
    ? undefined
    : {
        command:
          "python3 -m http.server 1313 --bind 127.0.0.1 --directory public",
        cwd: path.resolve(import.meta.dirname, "../.."),
        url: "http://127.0.0.1:1313/build/",
        stderr: "ignore",
      },
});
