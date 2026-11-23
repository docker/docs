import assert from "node:assert/strict";
import { pathToFileURL } from "node:url";
import { chromium } from "playwright";

export default async function verifyChatLayout(parentPage, base) {
  const context = await parentPage.context().browser().newContext();
  await context.route("**/*", (route) =>
    new URL(route.request().url()).origin === new URL(base).origin
      ? route.continue()
      : route.abort(),
  );
  const page = await context.newPage();
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  try {
    for (const [width, panelWidth] of [
      [390, 374],
      [820, 440],
      [1280, 480],
      [1535, 480],
      [1536, 480],
      [1920, 576],
      [2560, 608],
    ]) {
      await page.setViewportSize({ width, height: 900 });
      await page.goto(base + "/build/");
      const opener = page.locator(".page-actions button[marlin-action]");
      await opener.click();
      await page.waitForFunction(
        () => document.activeElement.id === "question",
      );
      const panel = page.locator("#gordon-chat");
      await page.waitForFunction(
        () =>
          Math.abs(
            document.querySelector("#gordon-chat").getBoundingClientRect()
              .right -
              (innerWidth - (innerWidth >= 1536 ? 0 : 8)),
          ) < 0.1,
      );
      const bounds = await panel.boundingBox();
      assert.equal(Math.round(bounds.width), panelWidth);
      assert.equal(
        await page.evaluate(
          () => document.documentElement.scrollWidth > innerWidth,
        ),
        false,
      );
      if (width >= 1536) {
        assert.equal(bounds.y, 0);
        assert.equal(bounds.height, 900);
        assert.equal(await panel.getAttribute("aria-modal"), null);
        assert.equal(await page.locator(".gordon-backdrop").isVisible(), false);
        for (const selector of [".docs-main", ".site-header", ".site-footer"]) {
          const box = await page.locator(selector).boundingBox();
          assert.ok(
            box.x + box.width <= bounds.x + 0.5,
            selector + " clears chat",
          );
        }
        const shell = await page.locator(".docs-shell").boundingBox();
        assert.ok(
          Math.abs(shell.x + shell.width - bounds.x) < 0.5,
          "Dock joins the page",
        );
        await page.locator(".page-actions a").first().focus();
        assert.equal(
          await panel.evaluate((el) => el.contains(document.activeElement)),
          false,
        );
        await page.keyboard.press("Escape");
        assert.equal(
          await panel.isVisible(),
          true,
          "Escape in the docs leaves docked chat open",
        );
        await page.evaluate(() => window.scrollTo(0, 300));
        assert.ok(await page.evaluate(() => scrollY > 0));
      } else {
        assert.equal(await panel.getAttribute("aria-modal"), "true");
        assert.equal(await page.locator(".gordon-backdrop").isVisible(), true);
        assert.equal(
          await page.evaluate(
            () => getComputedStyle(document.documentElement).overflow,
          ),
          "hidden",
        );
        for (let i = 0; i < 12; i++) {
          await page.keyboard.press("Tab");
          assert.equal(
            await panel.evaluate((el) => el.contains(document.activeElement)),
            true,
          );
        }
      }
      await panel
        .getByRole("button", { name: "Close chat", exact: true })
        .click();
      await page.waitForFunction(() => !Alpine.store("gordon").isOpen);
      await page.waitForFunction(() =>
        document.activeElement.matches(".page-actions button[marlin-action]"),
      );
    }

    // Crossing the breakpoint preserves the draft and releases modal restrictions.
    await page.setViewportSize({ width: 1535, height: 900 });
    await page.goto(base + "/build/");
    await page.locator(".page-actions button[marlin-action]").click();
    await page.locator("#question").fill("Explain build caching");
    await page.setViewportSize({ width: 1536, height: 900 });
    await page.waitForFunction(
      () => !document.querySelector("#gordon-chat").hasAttribute("aria-modal"),
    );
    await page.waitForFunction(
      () => getComputedStyle(document.documentElement).overflow !== "hidden",
    );
    assert.equal(
      await page.locator("#question").inputValue(),
      "Explain build caching",
    );
    assert.equal(
      await page.locator(".docs-shell").getAttribute("aria-hidden"),
      null,
    );
    await page.waitForFunction(() => document.activeElement.id === "question");
    await page.locator(".page-actions a").first().focus();
    await page.setViewportSize({ width: 1535, height: 900 });
    await page.waitForFunction(() => document.activeElement.id === "question");
    assert.equal(
      await page.locator("#question").inputValue(),
      "Explain build caching",
    );
    await page.keyboard.press("Escape");
    await page.waitForFunction(() => !Alpine.store("gordon").isOpen);

    assert.deepEqual(errors, []);
  } finally {
    await context.close();
  }
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  const browser = await chromium.launch({
    executablePath: process.env.CHROMIUM_PATH || undefined,
  });
  try {
    await verifyChatLayout(
      await browser.newPage(),
      process.env.DOCS_TEST_URL || "http://localhost:1313",
    );
    console.log("Chat layout and keyboard checks passed");
  } finally {
    await browser.close();
  }
}
