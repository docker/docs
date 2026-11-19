import assert from "node:assert/strict";

export default async function verifyWashes(page, base) {
  for (const width of [1440, 390]) {
    await page.setViewportSize({ width, height: 900 });
    const paints = {};
    for (const [category, route] of [
      ["home", "/"],
      ["ai", "/c/ai/"],
      ["dev", "/c/dev/"],
      ["admin", "/c/admin/"],
    ]) {
      await page.goto(base + route);
      const samples = {};
      for (const theme of ["light", "dark"]) {
        await page.evaluate((theme) => {
          document.documentElement.classList.toggle("dark", theme === "dark");
          document.documentElement.classList.toggle("light", theme === "light");
        }, theme);
        samples[theme] = await page.evaluate(() => ({
          category: document.body.dataset.category,
          header: getComputedStyle(document.body, "::before").backgroundImage,
          footer: getComputedStyle(document.querySelector(".site-footer"))
            .backgroundImage,
          pointerEvents: getComputedStyle(document.body, "::before")
            .pointerEvents,
        }));
        assert.equal(samples[theme].category, category);
        assert.notEqual(samples[theme].header, "none");
        assert.notEqual(samples[theme].footer, "none");
        assert.notEqual(samples[theme].header, samples[theme].footer);
        assert.equal(samples[theme].pointerEvents, "none");
      }
      paints[category] = samples;
    }
    assert.equal(
      new Set(Object.values(paints).map((p) => p.light.header)).size,
      4,
    );
    assert.equal(
      new Set(Object.values(paints).map((p) => p.dark.header)).size,
      4,
    );
    assert.equal(
      new Set(Object.values(paints).map((p) => p.dark.footer)).size,
      3, // Home and application development share the quiet right-hand footer.
    );
    for (const [category, route] of [
      ["ai", "/ai/sandboxes/"],
      ["dev", "/build/"],
      ["admin", "/accounts/"],
      ["dev", "/reference/dockerfile/"],
      ["dev", "/reference/cli/docker/compose/"],
      ["ai", "/reference/cli/docker/mcp/"],
      ["ai", "/reference/cli/docker/model/"],
      ["ai", "/reference/api/ai-governance/latest/"],
      ["ai", "/reference/api/sandboxes/latest/"],
    ]) {
      await page.goto(base + route);
      await page.evaluate(() => {
        document.documentElement.classList.remove("dark");
        document.documentElement.classList.add("light");
      });
      const paint = await page.evaluate(
        () => getComputedStyle(document.body, "::before").backgroundImage,
      );
      assert.equal(
        paint,
        paints[category].light.header,
        `${route}: category wash`,
      );
      await page.goto(base + route + "?check=stable#example");
      await page.evaluate(() => {
        document.documentElement.classList.remove("dark");
        document.documentElement.classList.add("light");
      });
      assert.equal(
        await page.evaluate(
          () => getComputedStyle(document.body, "::before").backgroundImage,
        ),
        paint,
        `${route}: stable category wash`,
      );
    }
  }
  const noJS = await page
    .context()
    .browser()
    .newPage({ javaScriptEnabled: false });
  try {
    await noJS.goto(base + "/c/ai/");
    assert.notEqual(
      await noJS.evaluate(
        () => getComputedStyle(document.body, "::before").backgroundImage,
      ),
      "none",
    );
    assert.notEqual(
      await noJS
        .locator(".site-footer")
        .evaluate((e) => getComputedStyle(e).backgroundImage),
      "none",
    );
  } finally {
    await noJS.close();
  }
}
