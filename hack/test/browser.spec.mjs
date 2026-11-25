import { test, expect } from "@playwright/test";
import axe from "axe-core";

test.beforeEach(async ({ context, baseURL }) => {
  await context.route("**/*", (route) =>
    new URL(route.request().url()).origin === new URL(baseURL).origin
      ? route.continue()
      : route.abort(),
  );
});

for (const route of [
  "/",
  "/build/",
  "/get-docker/desktop/",
  "/reference/cli/docker/buildx/build/",
]) {
  test(`page and accessibility: ${route}`, async ({ page }, testInfo) => {
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    const response = await page.goto(route);
    expect(response.status()).toBe(200);
    await page.waitForFunction(() => Boolean(window.Alpine));
    await expect(page.getByRole("heading", { level: 1 })).toHaveCount(1);
    await expect(page.locator('main a[href=""]')).toHaveCount(0);
    expect(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth + 1,
      ),
    ).toBe(true);
    await page.evaluate(axe.source);
    const { violations, incomplete } = await page.evaluate(() =>
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
    if (incomplete.length)
      await testInfo.attach("axe-incomplete", {
        body: JSON.stringify(incomplete, null, 2),
        contentType: "application/json",
      });
    expect(
      violations.map(({ id, nodes }) => ({
        id,
        targets: nodes.map(({ target }) => target),
      })),
    ).toEqual([]);
    expect(errors).toEqual([]);
  });
}

test("sidebar navigation remains usable after resizing", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 900 });
  await page.goto("/build/");
  const sidebar = page.locator(".sidebar-panel");
  const toggle = sidebar.locator(".sidebar-toggle");
  await expect(sidebar).toHaveJSProperty("open", false);
  await toggle.press("Enter");
  await expect(sidebar.locator('a[href="/build/checks/"]')).toBeVisible();
  await toggle.press("Enter");
  await expect(sidebar).toHaveJSProperty("open", false);
  await page.setViewportSize({ width: 1440, height: 900 });
  await expect(sidebar.locator(".sidebar-scroll")).toBeVisible();
});

test("search returns indexed results and restores keyboard focus", async ({
  page,
}) => {
  await page.goto("/build/");
  const opener = page.locator("[data-open-search]").first();
  await opener.focus();
  await opener.press("Enter");
  await expect(page.locator("#search-query")).toBeFocused();
  await page.locator("#search-query").fill("build cache");
  const results = page.locator("#docs-search .search-results li a");
  await expect(results.first()).toBeVisible();
  await page.keyboard.press("Escape");
  await expect(opener).toBeFocused();
});

test("search failure retains the option to ask Docker", async ({ page }) => {
  await page.route("**/pagefind/**", (route) => route.abort());
  await page.goto("/build/");
  await page.locator("[data-open-search]").first().click();
  await page.locator("#search-query").fill("build cache");
  await expect(page.locator("#docs-search [role=status]")).toContainText(
    "unavailable",
  );
  await expect(page.locator("#docs-search .search-ask")).toBeEnabled();
});

test("installation tabs work with a keyboard", async ({ page }) => {
  await page.goto("/get-docker/desktop/");
  const mac = page.getByRole("tab", { name: "Mac", exact: true });
  const windows = page.getByRole("tab", { name: "Windows", exact: true });
  await mac.focus();
  await mac.press("ArrowRight");
  await expect(windows).toBeFocused();
  await expect(windows).toHaveAttribute("aria-selected", "true");
  await expect(page.getByRole("tabpanel")).toBeVisible();
});

test("code copying preserves focus while pending and reports failures", async ({
  page,
}) => {
  await page.goto("/get-docker/engine/");
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
  const status = button.locator("..").getByRole("status");
  await button.focus();
  await button.press("Enter");
  await expect(button).toHaveAttribute("aria-busy", "true");
  await button.press("Enter");
  expect(await page.evaluate(() => window.clipboardCalls)).toBe(1);
  await expect(button).toBeFocused();
  await page.evaluate(() => window.finishCopy());
  await expect(status).toHaveText("Code copied.");
  await expect(button).toBeFocused();
  await page.evaluate(() => {
    Object.defineProperty(navigator.clipboard, "writeText", {
      value: async () => {
        throw Error("Clipboard denied");
      },
    });
  });
  await button.press("Enter");
  await expect(status).toContainText("Could not copy");
});

test("Gordon cards preserve text and reject unsafe links", async ({ page }) => {
  await page.goto("/build/");
  await page.addScriptTag({ url: "/gordon-cards.js" });
  const cases = [
    [
      "Download https://download.docker.com/linux/static/stable/.",
      "https://download.docker.com/linux/static/stable/",
      "Download https://download.docker.com/linux/static/stable/.",
    ],
    [
      "[nested](https://example.com/a_(b_(c)))",
      "https://example.com/a_(b_(c))",
      "nested",
    ],
    [
      "[<img src=x onerror=alert(1)>](https://example.com/)",
      "https://example.com/",
      "<img src=x onerror=alert(1)>",
    ],
    ...[
      "<img src=x onerror=alert(1)>",
      "[bad](javascript:alert(1))",
      "[bad](data:text/html,abc)",
      "[bad](ftp://example.com/file)",
      "`curl https://docs.docker.com/`",
    ].map((text) => [text, null, text]),
  ];
  for (const [input, href, text] of cases) {
    const rendered = await page.evaluate((input) => {
      const card = GordonCards.renderGordonCard(
        `root = Card("Example", [], ${JSON.stringify(input)})`,
      );
      const description = card.querySelector(".gordon-card-description");
      return {
        text: description.textContent,
        href: description.querySelector("a")?.href ?? null,
        images: description.querySelectorAll("img").length,
      };
    }, input);
    expect(rendered).toEqual({ text, href, images: 0 });
  }
});

test("API operations can be filtered and requests copied", async ({
  page,
  context,
}) => {
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto("/reference/api/hub/latest/");
  await page.locator("[data-api-filter]").fill("token");
  const operations = page.locator("[data-api-filter-item]:visible");
  await expect(operations.first()).toBeVisible();
  for (const text of await operations.allTextContents())
    expect(text.toLowerCase()).toContain("token");
  await operations.first().click();
  const copy = page.locator("[data-api-copy]");
  await copy.focus();
  await copy.press("Enter");
  await expect(copy).toHaveText("Copied");
  expect(await page.evaluate(() => navigator.clipboard.readText())).toContain(
    "curl",
  );
});

test("chat exposes answer and error states and restores focus when closed", async ({
  page,
}) => {
  const answer = "Docker café 🐳";
  await page.goto("/build/");
  await page.evaluate((answer) => {
    const original = window.fetch;
    window.fetch = async (url, options) => {
      if (!String(url).endsWith("/public/ask")) return original(url, options);
      if (window.failAnswer) return new Response(null, { status: 500 });
      return new Response(
        new ReadableStream({
          async start(controller) {
            await new Promise((resolve) => (window.releaseAnswer = resolve));
            controller.enqueue(
              new TextEncoder().encode(
                `data: ${JSON.stringify({ choices: [{ delta: { content: answer } }] })}\r\n\r\ndata: [DONE]\r\n\r\n`,
              ),
            );
            controller.close();
          },
        }),
        { headers: { "Content-Type": "text/event-stream" } },
      );
    };
  }, answer);
  const opener = page
    .locator(".page-actions")
    .getByRole("button", { name: "Ask Docker", exact: true });
  await opener.focus();
  await opener.press("Enter");
  const chat = page.locator("#gordon-chat");
  const input = page.locator("#question");
  await expect(input).toBeFocused();
  await expect(chat).toHaveAttribute("data-chat-status", "idle");
  await input.fill("Explain Docker");
  await input.press("Enter");
  await page.waitForFunction(() => Boolean(window.releaseAnswer));
  await expect(chat).toHaveAttribute("data-chat-status", "answering");
  await expect(chat.locator("[data-message-role]")).toHaveCount(2);
  expect(
    await chat
      .locator("[data-message-role]")
      .evaluateAll((elements) => elements.map((el) => el.dataset.messageRole)),
  ).toEqual(["user", "assistant"]);
  await page.evaluate(() => window.releaseAnswer());
  await expect(chat).toHaveAttribute("data-chat-status", "idle");
  await expect(chat.locator('[data-message-role="assistant"]')).toContainText(
    answer,
  );
  await expect(chat.locator(":scope > [role=status]")).toContainText(answer);
  await page.evaluate(() => (window.failAnswer = true));
  await input.fill("Try again");
  await input.press("Enter");
  await expect(chat).toHaveAttribute("data-chat-status", "error");
  await expect(chat.locator('[data-message-role="assistant"]')).toHaveCount(1);
  await page.getByRole("button", { name: "Close chat", exact: true }).click();
  await expect(opener).toBeFocused();
});

test.describe("without JavaScript", () => {
  test.use({ javaScriptEnabled: false });
  test("navigation and installation links remain available", async ({
    page,
  }) => {
    await page.goto("/build/");
    await expect(
      page.locator('.sidebar-panel a[href="/build/checks/"]'),
    ).toBeVisible();
    await page.goto("/");
    await expect(
      page.locator('.home-install a[href="/get-docker/desktop/"]'),
    ).toBeVisible();
    await expect(page.locator(".home-install pre:visible")).toHaveCount(0);
  });
});
