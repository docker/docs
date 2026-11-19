import assert from "node:assert/strict";

export default async function verifyLayout(page, base) {
  await page.setViewportSize({ width: 390, height: 900 });
  await page.goto(base + "/build/");
  await page.waitForFunction(() => Boolean(window.Alpine));
  const sidebar = page.locator(".sidebar-panel");
  const toggle = sidebar.locator(".sidebar-toggle");
  assert.equal(await sidebar.evaluate((panel) => panel.open), false);
  assert.equal(await toggle.isVisible(), true);

  await page.setViewportSize({ width: 700, height: 900 });
  await page.waitForFunction(
    () =>
      document.querySelector(".sidebar-panel").open &&
      getComputedStyle(document.querySelector(".sidebar-toggle")).display ===
        "none",
  );
  assert.equal(await toggle.isVisible(), false);
  assert.equal(await sidebar.locator(".sidebar-scroll").isVisible(), true);

  await page.setViewportSize({ width: 699, height: 900 });
  await page.waitForFunction(
    () => !document.querySelector(".sidebar-panel").open,
  );
  await toggle.click();
  assert.equal(await sidebar.evaluate((panel) => panel.open), true);
  await toggle.click();
  assert.equal(await sidebar.evaluate((panel) => panel.open), false);

  await page.setViewportSize({ width: 1200, height: 900 });
  await page.waitForFunction(
    () => document.querySelector(".sidebar-panel").open,
  );
  assert.equal(
    await sidebar.locator(".sidebar-scroll").isVisible(),
    true,
    "Widening the viewport restores a manually collapsed sidebar",
  );

  for (const route of [
    "/reference/cli/docker/buildx/build/",
    "/reference/cli/sbx/",
  ]) {
    await page.goto(base + route);
    const toc = page.locator(".article-layout > .reading-toc");
    assert.equal(await toc.count(), 1, `${route}: shared TOC wrapper`);
    assert.equal(await toc.getAttribute("data-pagefind-ignore"), "");
    assert.ok(await toc.locator("nav a").count(), `${route}: TOC has links`);
    for (const width of [1440, 1200, 1199, 1100, 1024, 820, 390]) {
      await page.setViewportSize({ width, height: 900 });
      await page.waitForFunction(
        (visible) =>
          (getComputedStyle(document.querySelector(".reading-toc")).display !==
            "none") ===
          visible,
        width >= 1200,
      );
      assert.equal(
        await toc.isVisible(),
        width >= 1200,
        `${route}: TOC visibility at ${width}px`,
      );
      const layout = await page.locator(".article-layout").boundingBox();
      const article = await page
        .locator(".article-content > article")
        .boundingBox();
      if (width >= 1200) {
        const box = await toc.boundingBox();
        const aside = await toc.locator("aside").boundingBox();
        assert.ok(box.width <= 181, `${route}: TOC fits its 180px column`);
        assert.ok(aside.width <= box.width + 1, `${route}: aside fits TOC`);
        assert.ok(
          box.x >= article.x + article.width,
          `${route}: TOC beside article`,
        );
        assert.ok(box.x + box.width <= layout.x + layout.width + 1);
        assert.equal(
          await toc.evaluate(
            (element) => element.scrollWidth > element.clientWidth + 1,
          ),
          false,
          `${route}: TOC content stays within its column`,
        );
      } else {
        assert.ok(
          Math.abs(article.width - layout.width) <= 1,
          `${route}: article fills the single column at ${width}px`,
        );
      }
    }
  }
}
