import { Marlin, Environment } from "@docker/marlin-sdk-web-public";

const config = window.__marlinConfig;

if (config && config.apiKey && config.endpoint) {
  const environment =
    config.environment === "staging" ? Environment.STAGING : Environment.PROD;
  try {
    window.marlin = new Marlin({
      endpoint: config.endpoint,
      apiKey: config.apiKey,
      site: "docs",
      environment,
    });
  } catch (err) {
    console.warn("Marlin SDK init failed:", err);
  }
}
