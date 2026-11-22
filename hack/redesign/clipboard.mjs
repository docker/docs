import assert from "node:assert/strict";

export default async function verifyClipboardFocus(page, base) {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(base + "/get-docker/engine/");
  await page.evaluate(() => {
    window.clipboardCalls = 0;
    Object.defineProperty(navigator.clipboard, "writeText", {
      configurable: true,
      value: () => {
        window.clipboardCalls++;
        return new Promise((resolve) => {
          window.finishCopy = resolve;
        });
      },
    });
  });
  const button = page.locator("[data-copy-text]").first();
  await button.focus();
  await page.keyboard.press("Enter");
  assert.equal(await button.getAttribute("aria-busy"), "true");
  assert.equal(
    await button.evaluate((e) => e === document.activeElement),
    true,
  );
  await page.keyboard.press("Enter");
  assert.equal(await page.evaluate(() => window.clipboardCalls), 1);
  await page.evaluate(() => window.finishCopy());
  await page.waitForFunction(
    () => !document.querySelector("[data-copy-text]").hasAttribute("aria-busy"),
  );
  assert.equal(
    await button.evaluate((e) => e === document.activeElement),
    true,
  );
  const status = button.locator("..").locator('[role="status"]');
  assert.equal(await status.textContent(), "Code copied.");
}
