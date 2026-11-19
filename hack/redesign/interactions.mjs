import assert from "node:assert/strict";

export default async function verifyInteractions(page, base) {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(base + "/");
  assert.deepEqual(
    await page
      .locator(".home-section-card")
      .evaluateAll((cards) => cards.map((card) => new URL(card.href).pathname)),
    ["/c/ai/", "/c/dev/", "/c/admin/"],
  );
  const updates = page.locator(".home-updates li");
  const updateToggle = page.locator(".home-updates-toggle");
  if (await updateToggle.count()) {
    const featured = await updates.filter({ visible: true }).count();
    await updateToggle.focus();
    await page.keyboard.press("Enter");
    await page.waitForFunction(() =>
      [...document.querySelectorAll(".home-updates li")].every(
        (item) => item.getBoundingClientRect().height > 0,
      ),
    );
    assert.equal(
      await updates.filter({ visible: true }).count(),
      await updates.count(),
    );
    const dates = await updates
      .locator("time")
      .evaluateAll((times) => times.map((time) => time.dateTime));
    assert.deepEqual(
      dates,
      [...dates].sort().reverse(),
      "Expanded updates stay chronological",
    );
    await page.keyboard.press("Enter");
    await page.waitForFunction(
      (count) =>
        [...document.querySelectorAll(".home-updates li")].filter(
          (item) => getComputedStyle(item).display !== "none",
        ).length === count,
      featured,
    );
  }
  for (const width of [1440, 390]) {
    await page.setViewportSize({ width, height: 900 });
    for (const route of ["/c/dev/", "/ai/sandboxes/", "/engine/"]) {
      await page.goto(base + route);
      const header = await page.locator(".reading-header").boundingBox();
      const nav = await page.locator(".site-header").boundingBox();
      assert.equal(header.x, nav.x, `${route}: header left edge`);
      assert.equal(header.width, nav.width, `${route}: header width`);
    }
    for (const route of ["/", "/get-docker/engine/"]) {
      await page.goto(base + route);
      const buttons = page.locator("[data-copy-text]");
      assert.ok(await buttons.count());
      for (const button of await buttons.all()) {
        await button.focus();
        assert.equal(
          await button.evaluate((e) => getComputedStyle(e).opacity),
          "1",
        );
        await page.keyboard.press("Enter");
        await page.waitForFunction(
          () =>
            document.activeElement
              .closest(".code-copy")
              .querySelector('[role="status"]').textContent === "Code copied.",
        );
        assert.equal(
          await page.evaluate(() => navigator.clipboard.readText()),
          await button.getAttribute("data-copy-text"),
        );
      }
      if (route === "/") {
        assert.ok(
          await page
            .locator(".home-install pre")
            .evaluateAll((es) =>
              es.every((e) => e.scrollWidth <= e.clientWidth),
            ),
        );
      }
      await page.evaluate(() =>
        Object.defineProperty(navigator.clipboard, "writeText", {
          configurable: true,
          value: async () => {
            throw Error("Clipboard denied");
          },
        }),
      );
      await buttons.first().click();
      await page.waitForFunction(() =>
        document
          .querySelector('.code-copy [role="status"]')
          .textContent.includes("Could not copy"),
      );
      assert.ok(
        await page.locator(".code-copy-status:not(.sr-only)").isVisible(),
      );
    }
  }
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(base + "/build/");
  await page.keyboard.press("Control+k");
  await page.locator("#search-query").fill("build cache");
  await page.waitForSelector("#docs-search .search-results li");
  const ask = page.locator("#docs-search .search-ask");
  const before = await ask.boundingBox();
  const results = page.locator("#docs-search .search-results");
  assert.ok(before.y + before.height <= (await results.boundingBox()).y);
  await results.evaluate((e) => (e.scrollTop = e.scrollHeight));
  assert.equal((await ask.boundingBox()).y, before.y, "Ask remains pinned");
  // Exercise the handoff without sending a real question to the backend.
  await page.evaluate(() => {
    const store = Alpine.store("gordon");
    const open = store.open.bind(store);
    store.open = (query, send) => {
      window.testAsk = { query, send };
      open(query, false);
    };
  });
  await ask.click();
  await page.waitForFunction(() => Alpine.store("gordon").isOpen);
  assert.deepEqual(await page.evaluate(() => window.testAsk), {
    query: "build cache",
    send: true,
  });
  assert.equal(
    await page.locator("#docs-search").evaluate((e) => e.open),
    false,
  );
  await page.getByRole("button", { name: "Close chat", exact: true }).click();
}
