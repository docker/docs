// Run with a Playwright Page connected to the local preview; no product API calls.
export default async function verify(page, base = "http://localhost:1314") {
  const results = [];
  const assert = (condition, message) => {
    if (!condition) throw Error(message);
    results.push(message);
  };
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.goto(base + "/api-prototype/");
  assert(
    (await page.locator(".api-card").count()) === 6,
    "Catalog exposes all six snapshots",
  );
  await page.goto(base + "/api-prototype/engine-1.56/");
  await page.locator("[data-api-filter]").fill("archive");
  const visible = page.locator("[data-api-filter-item]:visible");
  assert(
    (await visible.count()) > 0 && (await visible.count()) < 108,
    "Operation filter narrows Engine operations",
  );
  assert(
    (await visible.allTextContents()).every((t) =>
      t.toLowerCase().includes("archive"),
    ),
    "Filter results match their label",
  );
  await page.locator("[data-api-version]").focus();
  await page.keyboard.press("ArrowDown");
  await page.keyboard.press("Enter");
  await page.waitForURL("**/api-prototype/engine-1.55/");
  assert(
    page.url().endsWith("/engine-1.55/"),
    "Keyboard version selection changes the reference",
  );
  await page.goto(base + "/api-prototype/governance/operations/listPolicies/");
  await page.context().grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.locator("[data-api-copy]").focus();
  await page.keyboard.press("Enter");
  const copied = await page.evaluate(() => navigator.clipboard.readText());
  assert(
    copied.includes("${TOKEN}") &&
      !copied.includes("'Authorization: Bearer $TOKEN'"),
    "Copy preserves shell-expandable credentials",
  );
  const manual = await page
    .getByRole("link", { name: "Product manual", exact: false })
    .getAttribute("href");
  await page.goto(manual.startsWith("http") ? manual : base + manual);
  assert(
    (await page
      .getByRole("link", { name: /Explore the AI Governance API 1 prototype/ })
      .count()) === 1,
    "Product manual links back to the same API reference",
  );
  await page.goto(base + "/api-prototype/engine-1.56/operations/SystemPing/");
  const select = page.locator("[data-api-media-select]");
  const options = await select
    .locator("option")
    .evaluateAll((nodes) => nodes.map((n) => n.value).filter(Boolean));
  assert(options.length > 1, "Engine exposes multiple media types");
  await select.selectOption(options[0]);
  assert(
    (await page.locator("[data-api-media][hidden]").count()) > 0,
    "Media selection filters without removing source variants",
  );
  await page.goto(base + "/api-prototype/governance/operations/createPolicy/");
  const exampleSelect = page.locator("[data-api-example-select]").first();
  if (await exampleSelect.count()) {
    const n = await exampleSelect.locator("option").count();
    if (n > 1) {
      await exampleSelect.selectOption({ index: 1 });
      assert(
        (await page.locator("[data-api-example][hidden]").count()) > 0,
        "Example selection updates displayed examples",
      );
    }
  }
  const search = await page.evaluate(async () => {
    const pf = await import("/pagefind/pagefind.js");
    const found = await pf.search("GwPriority");
    return Promise.all(
      found.results.slice(0, 15).map(async (r) => {
        const d = await r.data();
        return { url: d.url, title: d.meta.title };
      }),
    );
  });
  assert(
    search.some(
      (r) =>
        r.url.includes("/api-prototype/engine-1.56/") &&
        r.title.includes("1.56"),
    ),
    "Search finds schema properties with API version context",
  );
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto(base + "/api-prototype/governance/operations/listPolicies/");
  assert(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth + 1,
    ),
    "Operation page fits a narrow viewport",
  );
  assert(
    await page.locator("[data-api-copy]").isVisible(),
    "Request example remains usable on a narrow screen",
  );
  await page.setViewportSize({ width: 1440, height: 1000 });
  const noJS = await page
    .context()
    .browser()
    .newContext({ javaScriptEnabled: false });
  const staticPage = await noJS.newPage();
  await staticPage.goto(
    base + "/api-prototype/governance/operations/listPolicies/",
  );
  assert(
    (await staticPage.locator("[data-api-variant]").count()) > 0 &&
      (await staticPage.locator(".api-signature").isVisible()),
    "Reference content is present without JavaScript",
  );
  await noJS.close();
  return { result: "pass", checks: results, search };
}
