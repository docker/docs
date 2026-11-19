import assert from "node:assert/strict";

export default async function verifyNotFound(page, base) {
  const response = await page.request.get(base + "/404.html");
  const html = await response.text();
  const sample = "/guides/retired-guide/";
  const paths = [
    [sample, true],
    ["/guides", true],
    ["/guides/", true],
    ["/guides-old/", false],
    ["/unrelated/?next=/guides/", false],
  ];
  const fixture = async (route) => {
    await route.fulfill({ status: 404, contentType: "text/html", body: html });
  };
  const fixtureURL = (url) => paths.some(([path]) => url.href === base + path);
  await page.route(fixtureURL, fixture);
  try {
    for (const width of [1440, 390]) {
      await page.setViewportSize({ width, height: 900 });
      for (const theme of ["light", "dark"]) {
        await page.addInitScript(
          (value) => localStorage.setItem("theme-preference", value),
          theme,
        );
        for (const [path, visible] of paths) {
          const response = await page.goto(base + path);
          assert.equal(response.status(), 404);
          const callout = page.locator(".not-found-guides");
          assert.equal(await callout.isVisible(), visible, path);
          assert.equal(
            await callout.locator("a").getAttribute("href"),
            "https://docs.docker.com/learn/",
          );
          assert.equal(
            await page.evaluate(
              () => document.documentElement.scrollWidth > innerWidth,
            ),
            false,
          );
        }
      }
    }
    let release;
    const scriptsReady = new Promise((resolve) => {
      release = resolve;
    });
    await page.route("**/scripts.js", async (route) => {
      await scriptsReady;
      await route.continue();
    });
    try {
      await page.goto(base + sample, { waitUntil: "commit" });
      await page.locator(".not-found-guides").waitFor({ state: "visible" });
      await page.evaluate(() => document.fonts.ready);
      const before = await page.locator(".not-found-links").boundingBox();
      release();
      await page.waitForFunction(() => Boolean(window.Alpine));
      assert.deepEqual(
        await page.locator(".not-found-links").boundingBox(),
        before,
        "The guides callout is present before the shared bundle initializes",
      );
    } finally {
      release();
      await page.unroute("**/scripts.js");
    }
  } finally {
    await page.unroute(fixtureURL, fixture);
  }
}
