import assert from "node:assert/strict";
import { mkdir, writeFile } from "node:fs/promises";
import axe from "axe-core";
import { chromium } from "playwright";

const base = process.env.DOCS_TEST_URL || "http://localhost:1313";
const browser = await chromium.launch({
  executablePath: process.env.CHROMIUM_PATH || undefined,
});
const routes = [
  "/",
  "/c/dev/",
  "/build/",
  "/engine/",
  "/testcontainers/",
  "/get-docker/engine/",
  "/desktop/setup/install/mac-install/",
  "/desktop/use-desktop/builds/",
  "/reference/dockerfile/",
  "/reference/cli/docker/buildx/build/",
  "/ai/sandboxes/architecture/",
  "/ai/sandboxes/security/",
  "/404.html",
];
let scans = 0;
const incomplete = [];
async function scan(page, label) {
  await page.evaluate(axe.source);
  const results = await page.evaluate(() =>
    axe.run(document, {
      runOnly: {
        type: "tag",
        values: [
          "wcag2a",
          "wcag2aa",
          "wcag21a",
          "wcag21aa",
          "wcag22aa",
          "best-practice",
        ],
      },
    }),
  );
  if (results.incomplete.length) {
    incomplete.push({
      label,
      route: new URL(page.url()).pathname,
      theme: await page.evaluate(() =>
        document.documentElement.classList.contains("dark") ? "dark" : "light",
      ),
      viewport: page.viewportSize(),
      findings: results.incomplete.map(({ id, nodes }) => ({
        id,
        nodes: nodes.map(({ target, failureSummary, any, all, none }) => ({
          target,
          failureSummary,
          reasons: [...any, ...all, ...none].map(({ message }) => message),
        })),
      })),
    });
  }
  assert.deepEqual(
    results.violations.map(({ id, nodes }) => ({
      id,
      targets: nodes.map(({ target }) => target),
    })),
    [],
    label,
  );
  // Gradient contrast and simulated screen-reader behavior still require manual review.
  scans++;
}
try {
  for (const theme of ["light", "dark"]) {
    const context = await browser.newContext({
      colorScheme: theme,
      reducedMotion: "reduce",
    });
    await context.addInitScript(
      (theme) => localStorage.setItem("theme-preference", theme),
      theme,
    );
    // Keep the audit repeatable and never send questions to the live AI service.
    await context.route("**/*", (route) =>
      new URL(route.request().url()).origin === new URL(base).origin
        ? route.continue()
        : route.abort(),
    );
    const page = await context.newPage();
    for (const width of [1440, 320]) {
      await page.setViewportSize({ width, height: 900 });
      for (const route of routes) {
        const response = await page.goto(base + route);
        assert.equal(response.status(), 200, route);
        await page.waitForFunction(() => Boolean(window.Alpine));
        await scan(page, `${route}, ${theme}, ${width}`);
        assert.ok(
          await page.evaluate(
            () => document.documentElement.scrollWidth <= innerWidth,
          ),
          `${route} reflows at ${width}px`,
        );
      }
    }

    // A 720px CSS viewport models a 1440px window at 200% browser zoom.
    await page.setViewportSize({ width: 720, height: 450 });
    for (const route of [
      "/",
      "/desktop/setup/install/mac-install/",
      "/reference/dockerfile/",
    ]) {
      await page.goto(base + route);
      assert.ok(
        await page.evaluate(
          () => document.documentElement.scrollWidth <= innerWidth,
        ),
        `${route} at 200% zoom`,
      );
    }
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.goto(base + "/reference/api/sandboxes/latest/");
    assert.equal(
      await page.locator("[data-api-filter-status]").getAttribute("role"),
      "status",
    );
    const operationURL = await page
      .locator("[data-api-filter-item]")
      .first()
      .getAttribute("href");
    await page
      .locator("[data-api-filter]")
      .fill("no-matching-operation-for-this-test");
    await page.waitForFunction(
      () =>
        document.querySelector("[data-api-filter-status]").textContent ===
        "0 operations found.",
    );
    await page.goto(base + operationURL);
    assert.equal(
      await page.locator("[data-api-copy-status]").getAttribute("role"),
      "status",
    );
    await page.evaluate(() => {
      navigator.clipboard.writeText = async () => {};
    });
    const requestCopy = page.locator("[data-api-copy]");
    await requestCopy.focus();
    await page.keyboard.press("Enter");
    await page.waitForFunction(
      () =>
        document.querySelector("[data-api-copy-status]").textContent ===
        "Request copied.",
    );
    await page.evaluate(() => {
      navigator.clipboard.writeText = async () => {
        throw new Error("Clipboard denied");
      };
    });
    await page.keyboard.press("Enter");
    await page.waitForFunction(() =>
      document
        .querySelector("[data-api-copy-status]")
        .textContent.includes("Could not copy"),
    );
    await page.goto(base + "/desktop/setup/install/mac-install/");
    await page.keyboard.press("Tab");
    assert.equal(await page.locator(":focus").textContent(), "Skip to content");
    await page.keyboard.press("Enter");
    assert.equal(
      await page.locator(":focus").getAttribute("id"),
      "main-content",
    );

    const tabs = page.getByRole("tablist").first();
    const tabCount = await tabs.getByRole("tab").count();
    assert.ok(tabCount >= 2, "Keyboard navigation needs multiple tabs");
    const first = tabs.getByRole("tab").first();
    await first.focus();
    for (const [key, index] of [
      ["ArrowRight", 1],
      ["End", tabCount - 1],
      ["ArrowLeft", tabCount - 2],
      ["Home", 0],
    ]) {
      const expected = tabs.getByRole("tab").nth(index);
      const tabId = await expected.getAttribute("id");
      const panelId = await expected.getAttribute("aria-controls");
      await page.keyboard.press(key);
      await page.waitForFunction(
        (id) =>
          document.activeElement.id === id &&
          document.activeElement.getAttribute("aria-selected") === "true",
        tabId,
      );
      await page.locator(`[id="${panelId}"]`).waitFor({ state: "visible" });
      assert.equal(await tabs.locator('[aria-selected="true"]').count(), 1);
      assert.equal(await page.getByRole("tabpanel").count(), 1);
      assert.equal(
        await page.getByRole("tabpanel").getAttribute("id"),
        panelId,
      );
      assert.equal(await page.getByRole("tabpanel").isVisible(), true);
    }
    await scan(page, `${theme} tabs after keyboard navigation`);
    await page.goto(base + "/reference/api/engine/sdk/examples/");
    await page.evaluate(() =>
      window.dispatchEvent(
        new CustomEvent("tab-select", {
          detail: { group: "lang", name: "unavailable-language" },
        }),
      ),
    );
    for (const group of await page.locator(".tabs").all()) {
      assert.equal(
        await group.getByRole("tabpanel").count(),
        1,
        "Unavailable synchronized tabs keep their selection",
      );
    }
    await page.goto(base + "/desktop/setup/install/mac-install/");

    const disclosure = page.locator(".nav-branch summary").first();
    await disclosure.focus();
    const wasOpen = await disclosure.evaluate((e) => e.parentElement.open);
    await page.keyboard.press("Space");
    assert.equal(
      await disclosure.evaluate((e) => e.parentElement.open),
      !wasOpen,
    );

    await page.goto(base + "/reference/cli/docker/buildx/build/");
    const tooltipButton = page.locator("[data-tooltip-button]").first();
    await tooltipButton.focus();
    const tooltip = page.locator(
      `#${await tooltipButton.getAttribute("aria-describedby")}`,
    );
    await tooltip.waitFor({ state: "visible" });
    await page.keyboard.press("Escape");
    await tooltip.waitFor({ state: "hidden" });
    await tooltipButton.evaluate((el) => el.blur());
    await page.mouse.move(0, 0);
    await tooltipButton.hover();
    await tooltip.hover();
    await tooltip.waitFor({ state: "visible" });

    await page.goto(base + "/desktop/use-desktop/builds/");
    const image = page.locator(".image-zoom-trigger").first();
    await image.focus();
    await page.keyboard.press("Enter");
    await scan(page, `${theme} enlarged image`);
    await page.keyboard.press("Control+k");
    assert.equal(
      await page.locator(".image-zoom-dialog[open]").count(),
      0,
      "Search closes image zoom before opening",
    );
    await page.keyboard.press("Escape");
    assert.equal(
      await image.evaluate((e) => e === document.activeElement),
      true,
    );

    await page.goto(base + "/ai/docker-agent/features/tui/");
    const animation = page.locator(".animated-image").first();
    assert.equal(
      await animation.evaluate((el) => el.open),
      false,
      "Animations start hidden",
    );
    await animation.locator("summary").focus();
    await page.keyboard.press("Enter");
    assert.equal(await animation.evaluate((el) => el.open), true);
    await page.keyboard.press("Enter");
    assert.equal(await animation.evaluate((el) => el.open), false);
    await page.goto(base + "/desktop/use-desktop/builds/");
    const search = page.locator("[data-open-search]");
    await search.focus();
    await page.keyboard.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.id === "search-query",
    );
    await page.locator("#search-query").fill("build cache");
    await page.waitForFunction(() =>
      document.querySelector("#docs-search .search-results li"),
    );
    await scan(page, `${theme} populated search`);
    await page.setViewportSize({ width: 320, height: 900 });
    await scan(page, `${theme} populated mobile search`);
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.keyboard.press("Escape");
    assert.equal(
      await search.evaluate((e) => e === document.activeElement),
      true,
    );

    // Mock both response and feedback failures to inspect announcement semantics.
    await page.route("**/public/ask", (route) =>
      route.fulfill({
        contentType: "text/event-stream",
        body:
          [
            { question_answer_id: "a11y-test" },
            {
              choices: [
                {
                  delta: {
                    content:
                      "## Run containers\n\nUse Docker to build and run containers.\n\n| Command | Purpose |\n| --- | --- |\n| docker run | Start a container |\n\n```console\ndocker run hello-world\n```",
                  },
                },
              ],
            },
          ]
            .map((data) => `data: ${JSON.stringify(data)}\n\n`)
            .join("") + "data: [DONE]\n\n",
      }),
    );
    await page.route("**/feedback", (route) => route.fulfill({ status: 500 }));
    const ask = page
      .locator(".page-actions")
      .getByRole("button", { name: "Ask Docker", exact: true });
    await ask.focus();
    await page.keyboard.press("Enter");
    const chat = page.getByRole("dialog", { name: "Ask Docker", exact: true });
    const input = chat.getByRole("textbox", {
      name: "Ask a question about Docker",
    });
    await page.waitForFunction(() => document.activeElement.id === "question");
    for (let i = 0; i < 12; i++) {
      await page.keyboard.press("Tab");
      assert.ok(
        await chat.evaluate((el) => el.contains(document.activeElement)),
        "Chat traps forward focus",
      );
    }
    for (let i = 0; i < 12; i++) {
      await page.keyboard.press("Shift+Tab");
      assert.ok(
        await chat.evaluate((el) => el.contains(document.activeElement)),
        "Chat traps backward focus",
      );
    }
    await input.fill("How do I use Docker?");
    await page.keyboard.press("Enter");
    await page.waitForFunction(() =>
      document
        .querySelector("#gordon-chat > [role=status]")
        .textContent.includes("Use Docker"),
    );
    await scan(page, `${theme} completed chat`);
    await page.setViewportSize({ width: 320, height: 900 });
    await scan(page, `${theme} completed mobile chat`);
    assert.ok(
      await chat.evaluate((el) => el.scrollWidth <= el.clientWidth),
      "Chat reflows on mobile",
    );
    await page.setViewportSize({ width: 1440, height: 900 });
    await chat.getByRole("button", { name: "Helpful", exact: true }).click();
    await chat.getByRole("alert").waitFor();
    await scan(page, `${theme} chat feedback error`);
    await page.keyboard.press("Escape");
    await page.waitForFunction(() => !Alpine.store("gordon").isOpen);
    // Alpine's focus trap restores the opener after its closing transition.
    await page.waitForFunction(() =>
      document.activeElement.textContent.includes("Ask Docker"),
    );

    await page.route("**/public/ask", (route) =>
      route.fulfill({ status: 500 }),
    );
    await ask.click();
    await input.fill("Try again");
    await page.keyboard.press("Enter");
    await page.waitForFunction(() =>
      document
        .querySelector("#gordon-chat > [role=status]")
        .textContent.includes("Failed to get response"),
    );
    await scan(page, `${theme} chat request error`);
    await page.keyboard.press("Escape");
    await context.close();
  }
  console.log(
    `Passed ${scans} accessibility scans, keyboard flows, and reflow checks.`,
  );
} finally {
  try {
    await mkdir("tmp/redesign", { recursive: true });
    await writeFile(
      "tmp/redesign/accessibility-incomplete.json",
      JSON.stringify(incomplete, null, 2) + "\n",
    );
  } finally {
    await browser.close();
  }
}
