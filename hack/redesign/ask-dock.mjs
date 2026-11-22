import assert from "node:assert/strict";

export default async function verifyAskDock(page, base) {
  for (const width of [1440, 390]) {
    await page.setViewportSize({ width, height: 900 });
    await page.goto(base + "/build/");
    const dock = page.locator(".ask-dock");
    await dock.waitFor({ state: "visible" });
    const bounds = await dock.boundingBox();
    assert.ok(bounds.y >= 0 && bounds.y + bounds.height <= 900);
    if (width > 700) {
      const input = dock.getByRole("textbox");
      const send = dock.getByRole("button", {
        name: "Ask Docker",
        exact: true,
      });
      await page.evaluate(() => window.scrollTo(0, 600));
      const scrollY = await page.evaluate(() => window.scrollY);
      await input.click();
      assert.equal(
        await page.evaluate(() => window.scrollY),
        scrollY,
        "Focusing the dock leaves the page in place",
      );
      assert.equal(await send.isDisabled(), true);
      await page.evaluate(() => {
        const store = Alpine.store("gordon");
        const open = store.open.bind(store);
        store.open = (query, send) => {
          window.dockRequest = { query, send };
          open(query, false);
        };
      });
      for (const character of "  Explain build caching  ") {
        await input.pressSequentially(character);
        await page.evaluate(() => new Promise(requestAnimationFrame));
        assert.equal(
          await page.evaluate(() => window.scrollY),
          scrollY,
          "Typing in the dock leaves the page in place",
        );
      }
      await input.press("Enter");
      await page.waitForFunction(() => Alpine.store("gordon").isOpen);
      assert.deepEqual(await page.evaluate(() => window.dockRequest), {
        query: "Explain build caching",
        send: true,
      });
      assert.equal(await dock.locator("input").inputValue(), "");
    } else {
      assert.equal(await dock.getByRole("textbox").isVisible(), false);
      await dock
        .getByRole("button", { name: "Ask Docker", exact: true })
        .click();
      await page.waitForFunction(() => Alpine.store("gordon").isOpen);
    }
    assert.equal(await dock.isVisible(), false);
    await page.getByRole("button", { name: "Close chat", exact: true }).click();
    await dock.waitFor({ state: "visible" });
    await page.locator("[data-open-search]").click();
    assert.equal(await dock.isVisible(), false);
    await page.keyboard.press("Escape");
    await dock.waitFor({ state: "visible" });
    await page.locator(".site-footer").scrollIntoViewIfNeeded();
    const end = await dock.boundingBox();
    const footer = await page.locator(".site-footer").boundingBox();
    assert.ok(end.y + end.height <= footer.y, "Dock clears the footer");
  }
  for (const route of [
    "/",
    "/c/dev/",
    "/reference/api/engine/version/v1.56/",
  ]) {
    await page.goto(base + route);
    assert.equal(await page.locator(".ask-dock").count(), 0, route);
  }
  const noJS = await page
    .context()
    .browser()
    .newPage({ javaScriptEnabled: false });
  try {
    await noJS.goto(base + "/build/");
    assert.equal(await noJS.locator(".ask-dock").isVisible(), false);
  } finally {
    await noJS.close();
  }
}
