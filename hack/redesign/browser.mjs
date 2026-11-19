// Run against a built artifact with Pagefind and flattened Markdown, served over HTTP.
import assert from "node:assert/strict";
import verifyInteractions from "./interactions.mjs";
import verifyWashes from "./washes.mjs";
import verifyAskDock from "./ask-dock.mjs";
import verifyChatLayout from "./chat-layout.mjs";
import verifyLayout from "./layout.mjs";
import verifyClipboardFocus from "./clipboard.mjs";
import verifyDiscovery from "./discovery.mjs";
import verifySearchChatRegressions from "./search-chat-regressions.mjs";
import verifyNotFound from "./not-found.mjs";
import { mkdir } from "node:fs/promises";
import { chromium } from "playwright";
const base = process.env.DOCS_TEST_URL || "http://localhost:1313";
const browser = await chromium.launch({
  executablePath: process.env.CHROMIUM_PATH || undefined,
});
const context = await browser.newContext({
  permissions: ["clipboard-read", "clipboard-write"],
});
await context.route("**/*", (route) =>
  new URL(route.request().url()).origin === new URL(base).origin
    ? route.continue()
    : route.abort(),
);
const page = await context.newPage();
const errors = [];
page.on("pageerror", (error) => errors.push(error.message));
const routes = [
  "/",
  "/c/ai/",
  "/c/dev/",
  "/c/admin/",
  "/get-docker/",
  "/get-docker/desktop/",
  "/get-docker/sandboxes/",
  "/get-docker/engine/",
  "/build/",
  "/compose/",
  "/ai/model-runner/",
  "/ai/sandboxes-api/",
  "/ai/governance/",
  "/accounts/",
  "/testcontainers/",
  "/release-lifecycle/",
  "/engine/",
  "/desktop/",
  "/reference/",
  "/c/",
  "/reference/cli/docker/buildx/build/",
  "/reference/cli/sbx/",
  "/reference/api/sandboxes/latest/",
  "/reference/dockerfile/",
  "/reference/compose-file/",
  "/404.html",
];
try {
  for (const width of [1440, 820, 390])
    for (const theme of ["light", "dark"]) {
      await page.setViewportSize({ width, height: 900 });
      await page.addInitScript(
        (theme) => localStorage.setItem("theme-preference", theme),
        theme,
      );
      for (const route of routes) {
        const response = await page.goto(base + route);
        await page.waitForFunction(() => Boolean(window.Alpine));
        assert.equal(response.status(), 200, route);
        if (route === "/reference/dockerfile/" && width === 1440) {
          const failures = await page.evaluate(() => {
            const canvas = document.createElement("canvas");
            canvas.width = canvas.height = 1;
            const context = canvas.getContext("2d");
            const luminance = (color) => {
              context.fillStyle = color;
              context.fillRect(0, 0, 1, 1);
              const rgb = [...context.getImageData(0, 0, 1, 1).data]
                .slice(0, 3)
                .map((value) => {
                  value /= 255;
                  return value <= 0.04045
                    ? value / 12.92
                    : ((value + 0.055) / 1.055) ** 2.4;
                });
              return rgb[0] * 0.2126 + rgb[1] * 0.7152 + rgb[2] * 0.0722;
            };
            const fixture = document.createElement("div");
            fixture.className = "copy-container code-block";
            fixture.style.cssText = "position:fixed;visibility:hidden";
            document.body.append(fixture);
            const failures = [];
            try {
              for (const highlighted of [false, true]) {
                const line = document.createElement("span");
                line.className = highlighted ? "hl" : "line";
                fixture.append(line);
                const background = luminance(
                  getComputedStyle(highlighted ? line : fixture)
                    .backgroundColor,
                );
                for (const token of [
                  "c",
                  "gp",
                  "k",
                  "kc",
                  "n",
                  "nf",
                  "s2",
                  "mi",
                  "err",
                  "p",
                  "ln",
                  "x",
                ]) {
                  const span = document.createElement("span");
                  span.className = token;
                  line.append(span);
                  const foreground = luminance(getComputedStyle(span).color);
                  const ratio =
                    (Math.max(foreground, background) + 0.05) /
                    (Math.min(foreground, background) + 0.05);
                  if (ratio < 4.5) failures.push({ token, highlighted, ratio });
                }
                line.remove();
              }
            } finally {
              fixture.remove();
            }
            return failures;
          });
          assert.deepEqual(
            failures,
            [],
            `${theme}: code text meets 4.5:1 on raised and highlighted surfaces`,
          );
        }
        assert.deepEqual(
          await page.locator(".prose").evaluateAll((elements) =>
            elements.flatMap((element) => {
              const style = getComputedStyle(element);
              return [
                ["--tw-prose-body", "--color-surface-foreground"],
                ["--tw-prose-captions", "--color-muted-foreground"],
                ["--tw-prose-th-borders", "--color-edge"],
              ].filter(
                ([prose, token]) =>
                  style.getPropertyValue(prose).trim() !==
                  style.getPropertyValue(token).trim(),
              );
            }),
          ),
          [],
          `${route}: prose uses Trident tokens in ${theme} mode`,
        );
        assert.deepEqual(
          await page
            .locator(".site-logo path")
            .evaluateAll((paths) => [
              ...new Set(paths.map((path) => getComputedStyle(path).fill)),
            ]),
          ["rgb(16, 16, 16)"],
          `${route}: content SVG styles do not recolor the header logo`,
        );
        for (const card of await page.locator(".card:has(.card-icon)").all()) {
          const offset = await card.evaluate((card) => {
            const icon = card
              .querySelector(".card-icon svg")
              .getBoundingClientRect();
            const title = card.querySelector(".card-title");
            const line = parseFloat(getComputedStyle(title).lineHeight);
            return Math.abs(
              icon.y +
                icon.height / 2 -
                title.getBoundingClientRect().y -
                line / 2,
            );
          });
          assert.ok(
            offset < 1,
            `${route}: card icon aligns with the title's first line`,
          );
        }
        assert.equal(
          await page.locator("h1").count(),
          1,
          `${route}: one page title`,
        );
        assert.equal(
          await page.evaluate(
            () => document.documentElement.scrollWidth > innerWidth + 1,
          ),
          false,
          `${route}: overflow at ${width}`,
        );
        assert.equal(
          await page
            .locator('main a[href=""], .site-footer a[href=""]')
            .count(),
          0,
          `${route}: empty navigation link`,
        );
        assert.deepEqual(errors, [], `${route}: browser errors`);
      }
    }
  let releaseScripts;
  const scriptsReady = new Promise((resolve) => {
    releaseScripts = resolve;
  });
  await page.route("**/scripts.js", async (route) => {
    await scriptsReady;
    await route.continue();
  });
  await page.goto(base + "/404.html?before-scripts", { waitUntil: "commit" });
  await page.locator(".not-found-links").waitFor();
  await page.evaluate(() => document.fonts.ready);
  const recoveryBeforeScripts = await page
    .locator(".not-found-links")
    .boundingBox();
  assert.equal(await page.locator(".not-found-search").isVisible(), true);
  releaseScripts();
  await page.waitForFunction(() => Boolean(window.Alpine));
  assert.deepEqual(
    await page.locator(".not-found-links").boundingBox(),
    recoveryBeforeScripts,
    "404 recovery links stay in place when scripts initialize",
  );
  await page.unroute("**/scripts.js");
  await page.goto(base + "/404.html?from=broken#missing");
  assert.equal(
    await page.locator(".breadcrumbs, .page-actions, .ask-dock").count(),
    0,
  );
  assert.equal(
    await page.locator("main[data-pagefind-body], main img").count(),
    0,
  );
  assert.equal(
    await page.locator('meta[name="robots"]').getAttribute("content"),
    "noindex",
  );
  assert.deepEqual(
    await page
      .locator(".not-found-links a")
      .evaluateAll((links) => links.map((link) => new URL(link.href).pathname)),
    ["/", "/c/ai/", "/c/dev/", "/c/admin/"],
  );
  const report = new URL(
    await page
      .getByRole("link", { name: "Report a broken link" })
      .getAttribute("href"),
  );
  assert.equal(
    report.searchParams.get("target"),
    base + "/404.html?from=broken#missing",
  );
  assert.equal(report.searchParams.get("template"), "broken_link.yml");
  const recoverySearch = page.getByRole("button", {
    name: "Search documentation",
    exact: true,
  });
  await recoverySearch.focus();
  await page.keyboard.press("Enter");
  assert.equal(
    await page.locator("#docs-search").evaluate((el) => el.open),
    true,
  );
  await page.keyboard.press("Escape");
  for (const [route, targets] of [
    ["/ai/sandboxes-api/", ["sandboxes"]],
    ["/ai/governance/", ["ai-governance", "ai-governance-audit-logs"]],
  ]) {
    await page.goto(base + route);
    for (const target of targets) {
      const selector = `a[href="/reference/api/${target}/latest/"]`;
      assert.equal(await page.locator(`.sidebar-panel ${selector}`).count(), 1);
      assert.equal(
        await page.locator(`.related-reference ${selector}`).count(),
        1,
      );
    }
  }
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.goto(base + "/build/");
  assert.match(
    await page.locator(".sidebar-scroll").innerText(),
    /Related reference/,
  );
  assert.doesNotMatch(
    await page.locator(".sidebar-scroll").innerText(),
    /Docker Desktop/,
  );
  await page.goto(base + "/reference/cli/docker/buildx/build/");
  assert.equal(
    await page.locator("[data-sidebar-key]").getAttribute("data-sidebar-key"),
    "/reference/cli/docker",
  );
  await page.goto(base + "/c/dev/");
  assert.equal(await page.locator(".sidebar-panel").count(), 0);
  await page.goto(base + "/get-docker/desktop/");
  await page.getByRole("tab", { name: "Mac", exact: true }).focus();
  await page.keyboard.press("ArrowRight");
  await page.waitForFunction(
    () =>
      document
        .querySelector("[role=tab][aria-selected=true]")
        ?.textContent.trim() === "Windows",
  );
  assert.equal(
    await page
      .getByRole("tab", { name: "Windows", exact: true })
      .getAttribute("aria-selected"),
    "true",
  );
  await page.goto(base + "/build/?test=markdown#overview");
  await page
    .getByRole("button", { name: "Copy Markdown", exact: true })
    .click();
  await page.waitForFunction(
    () =>
      document.querySelector(".page-actions [role=status]").textContent ===
      "Markdown copied",
  );
  assert.match(
    await page.evaluate(() => navigator.clipboard.readText()),
    /Docker Build/,
  );
  const ask = page
    .locator(".page-actions")
    .getByRole("button", { name: "Ask Docker", exact: true });
  await ask.click();
  await page.locator("#gordon-chat").waitFor({ state: "visible" });
  await page.keyboard.press("Escape");
  await page.waitForTimeout(350);
  assert.equal(
    await ask.evaluate((el) => el === document.activeElement),
    true,
    "Chat restores focus",
  );
  await page.keyboard.press("Control+k");
  await page.locator("#search-query").fill("build cache");
  await page.waitForFunction(() =>
    document.querySelector("#docs-search .search-results li"),
  );
  await page.keyboard.press("Escape");
  await page.route("**/pagefind/**", (route) => route.abort());
  await page.reload();
  await page.keyboard.press("Control+k");
  await page.locator("#search-query").fill("build cache");
  await page.waitForFunction(() =>
    document
      .querySelector("#docs-search .search-results")
      .textContent.includes("unavailable"),
  );
  assert.equal(
    await page.locator("#docs-search .search-ask").isEnabled(),
    true,
    "Search failure retains Ask",
  );
  await page.keyboard.press("Escape");
  await page.unroute("**/pagefind/**");
  await page.route("**/build.md", (route) => route.abort());
  await page
    .getByRole("button", { name: "Copy Markdown", exact: true })
    .click();
  await page.waitForFunction(() =>
    document
      .querySelector(".page-actions [role=status]")
      .textContent.includes("Could not copy"),
  );
  await page.unroute("**/build.md");
  const noJS = await browser.newPage({
    javaScriptEnabled: false,
    viewport: { width: 390, height: 844 },
  });
  await noJS.goto(base + "/build/");
  assert.equal(await noJS.locator(".sidebar-panel").isVisible(), true);
  await noJS.locator(".collection-tree summary").first().click();
  assert.equal(
    await noJS.locator(".collection-tree details").first().getAttribute("open"),
    "",
  );
  await noJS.goto(base + "/");
  assert.equal(await noJS.locator(".home-install pre:visible").count(), 0);
  await noJS.goto(base + "/404.html");
  assert.equal(await noJS.locator(".not-found-links a").count(), 4);
  assert.equal(await noJS.locator(".not-found-search").isVisible(), false);
  assert.match(
    await noJS.locator(".not-found-report").getAttribute("href"),
    /template=broken_link.yml/,
  );
  await noJS.close();
  await verifyInteractions(page, base);
  await verifyWashes(page, base);
  await verifyAskDock(page, base);
  await verifyChatLayout(page, base);
  await verifyLayout(page, base);
  await verifyClipboardFocus(page, base);
  await verifyDiscovery(page, base);
  await verifySearchChatRegressions(page, base);
  await verifyNotFound(page, base);
  await mkdir("/tmp/docs-redesign-captures", { recursive: true });
  for (const route of [
    "/",
    "/c/ai/",
    "/ai/sandboxes/",
    "/reference/cli/sbx/",
  ]) {
    await page.goto(base + route);
    await page.screenshot({
      path: `/tmp/docs-redesign-captures/${route.replaceAll("/", "_") || "home"}.png`,
      fullPage: true,
    });
  }
  console.log(
    `Passed ${routes.length * 6} page/theme/viewport checks and search, clipboard, tabs, focus, CLI context, no-JavaScript, shared code copying, alignment, installation links, and pinned Ask checks.`,
  );
} finally {
  await browser.close();
}
