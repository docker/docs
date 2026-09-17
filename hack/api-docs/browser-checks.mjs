// Run with a Playwright Page against a built site; no product API calls.
export default async function verify(page, base = "http://localhost:1314") {
  const results = [];
  const assert = (condition, message) => {
    if (!condition) throw Error(message);
    results.push(message);
  };
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.goto(base + "/reference/api/");
  assert(
    (await page.locator(".api-card").count()) === 5,
    "Five API catalog entries",
  );
  assert(
    (await page.locator("nav.navbar-font").count()) === 1,
    "Catalog retains Reference sidebar",
  );
  await page.goto(base + "/reference/api/dvp/latest/");
  assert(
    (await page.locator(".api-nav").count()) === 1 &&
      (await page.locator("nav.navbar-font").count()) === 0,
    "Reference uses local navigation",
  );
  const first = page.locator("[data-api-filter-item]").first();
  const operationURL = await first.getAttribute("href");
  await page.goto(base + "/reference/api/hub/dvp/");
  await page.waitForURL("**/reference/api/dvp/latest/");
  assert(
    page.url().endsWith("/reference/api/dvp/latest/"),
    "DVP alias reaches overview",
  );
  await page.goto(base + operationURL);
  await page.context().grantPermissions(["clipboard-read", "clipboard-write"]);
  const copy = page.locator("[data-api-copy]");
  await copy.focus();
  await page.keyboard.press("Enter");
  await page.waitForFunction(
    () => document.querySelector("[data-api-copy]").textContent === "Copied",
  );
  assert(
    (await page.evaluate(() => navigator.clipboard.readText())).includes(
      "curl",
    ),
    "Keyboard copies a request",
  );
  await page.locator(".api-nav-back").focus();
  await page.keyboard.press("Enter");
  await page.waitForURL("**/reference/api/");
  assert(
    (await page.locator("nav.navbar-font").count()) === 1,
    "Keyboard back link restores catalog",
  );
  await page.goto(base + "/reference/api/hub/latest/");
  await page.locator("[data-api-filter]").fill("token");
  const rows = page.locator("[data-api-filter-item]:visible");
  assert(
    (await rows.count()) > 0 &&
      (await rows.allTextContents()).every((x) =>
        x.toLowerCase().includes("token"),
      ),
    "Operation filtering",
  );
  const schemaURL = await page
    .locator(".api-schema-links a")
    .first()
    .getAttribute("href");
  await page.goto(base + schemaURL);
  assert(
    (await page.locator('[data-api-view="schema"]').count()) === 1,
    "Linked schema page",
  );
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto(base + operationURL);
  assert(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth + 1,
    ),
    "No narrow-screen document overflow",
  );
  await page.goto(base + "/reference/api/engine/version/v1.56/");
  assert(
    (await page.locator("redoc").count()) === 1,
    "Latest Engine retains ReDoc",
  );
  await page.goto(base + "/reference/api/ai-governance/");
  assert(
    (await page.locator("[data-api-view]").count()) === 0 &&
      (await page.locator("h1").count()) > 0,
    "Governance retains existing renderer",
  );
  const context = await page
    .context()
    .browser()
    .newContext({ javaScriptEnabled: false });
  try {
    const staticPage = await context.newPage();
    await staticPage.goto(base + operationURL);
    assert(
      (await staticPage.locator("[data-api-copy-source]").count()) === 1,
      "Request content exists without JavaScript",
    );
  } finally {
    await context.close();
  }
  return results;
}
