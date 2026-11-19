import assert from "node:assert/strict";
import { pathToFileURL } from "node:url";
import { chromium } from "playwright";

const event = (data) => `data: ${JSON.stringify(data)}\r\n\r\n`;
const answerStream = (content) =>
  event({ thread_id: "test-thread" }) +
  event({ question_answer_id: "test-answer" }) +
  event({ choices: [{ delta: { content } }] }) +
  "data: [DONE]\r\n\r\n";

async function mockAnswer(page, content, split = false) {
  await page.evaluate(
    ({ stream, split }) => {
      const original = window.fetch;
      window.chatRequests = [];
      window.fetch = async (url, options) => {
        if (!String(url).endsWith("/public/ask")) return original(url, options);
        window.chatRequests.push(JSON.parse(options.body));
        const bytes = new TextEncoder().encode(stream);
        return new Response(
          new ReadableStream({
            start(controller) {
              if (split) {
                // Include splits inside UTF-8 characters, JSON, CRLF and metadata IDs.
                for (const byte of bytes)
                  controller.enqueue(new Uint8Array([byte]));
              } else controller.enqueue(bytes);
              controller.close();
            },
          }),
          { headers: { "Content-Type": "text/event-stream" } },
        );
      };
    },
    { stream: answerStream(content), split },
  );
}

async function ask(page, question = "Explain Docker") {
  await page.evaluate(() => Alpine.store("gordon").open());
  const input = page.locator("#question");
  await input.fill(question);
  await input.press("Enter");
  await page.waitForFunction(
    () =>
      !Alpine.$data(document.querySelector("#gordon-chat").parentElement)
        .isLoading,
  );
}

export default async function verifySearchChatRegressions(parentPage, base) {
  const context = await parentPage.context().browser().newContext();
  await context.route("**/*", (route) =>
    new URL(route.request().url()).origin === new URL(base).origin
      ? route.continue()
      : route.abort(),
  );
  const page = await context.newPage();
  try {
    // Test pagination through the actual Alpine bindings and keyboard focus.
    await page.route("**/pagefind/pagefind.js", (route) =>
      route.fulfill({
        contentType: "text/javascript",
        body: `
      export async function options() {}
      export async function search(query) {
        return {results: Array.from({length: query === 'slow' ? 1 : 25}, (_, i) => ({data: async () => {
          if(query === 'slow') await new Promise(resolve => window.releaseSearch = resolve);
          if(query === 'retry' && i === 12 && !window.retried) await new Promise((resolve, reject) => window.rejectMore = () => reject(Error('Network failure')));
          return {url: '/result-' + i, meta: {title: query + ' ' + i, description: 'Description'}};
        }}))};
      }
    `,
      }),
    );
    // Editing must not collapse the homepage while replacement results load.
    await page.goto(base + "/");
    const homeQuery = page.locator("#home-query");
    const homeResults = page.locator(".home-hero .search-results li");
    await homeQuery.fill("docker");
    await page.waitForFunction(
      () =>
        document.querySelectorAll(".home-hero .search-results li").length ===
        12,
    );
    const installTop = await page
      .locator(".home-install")
      .evaluate((element) => element.getBoundingClientRect().top);
    await homeQuery.fill("slow");
    assert.equal(await homeResults.count(), 12);
    await page.waitForFunction(() => Boolean(window.releaseSearch));
    assert.equal(await homeResults.count(), 12);
    assert.equal(
      await page
        .locator(".home-install")
        .evaluate((element) => element.getBoundingClientRect().top),
      installTop,
    );
    assert.equal(
      await page
        .locator(".home-hero .search-results")
        .getAttribute("aria-busy"),
      "true",
    );
    await homeQuery.fill("");
    await page.evaluate(() => window.releaseSearch());
    assert.equal(await homeResults.count(), 0);

    // Homepage Enter keeps its Ask behavior after the search debounce changes.
    await mockAnswer(page, "Answer from the homepage");
    await homeQuery.fill("Explain build caching");
    await homeQuery.evaluate((input) =>
      input.dispatchEvent(
        new KeyboardEvent("keydown", {
          key: "Enter",
          bubbles: true,
          isComposing: true,
        }),
      ),
    );
    assert.equal(await page.locator("#gordon-chat").isVisible(), false);
    await homeQuery.press("Enter");
    await page.waitForFunction(() =>
      document
        .querySelector("#gordon-chat > [role=status]")
        .textContent.includes("Answer from the homepage"),
    );
    assert.equal(await page.evaluate(() => window.chatRequests.length), 1);

    await page.goto(base + "/build/");
    await page.locator("[data-open-search]").click();
    const query = page.locator("#search-query");
    const results = page.locator("#docs-search .search-results li");
    await query.fill("docker");
    await page.waitForFunction(
      () =>
        document.querySelectorAll("#docs-search .search-results li").length ===
        12,
    );
    const more = page.getByRole("button", { name: "Load more results" });
    await more.focus();
    await more.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.getAttribute("href") === "/result-12",
    );
    assert.equal(await results.count(), 24);
    await more.focus();
    await more.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.getAttribute("href") === "/result-24",
    );
    assert.equal(await results.count(), 25);
    assert.equal(await more.isVisible(), false);
    await query.fill("retry");
    await page.waitForFunction(
      () =>
        document.querySelectorAll("#docs-search .search-results li").length ===
        12,
    );
    await more.focus();
    await more.press("Enter");
    await page.waitForFunction(() => Boolean(window.rejectMore));
    assert.equal(
      await more.evaluate((element) => element === document.activeElement),
      true,
    );
    await page.evaluate(() => window.rejectMore());
    await page.waitForFunction(() =>
      document
        .querySelector("#docs-search [role=status]")
        .textContent.includes("Could not load"),
    );
    assert.equal(
      await more.evaluate((element) => element === document.activeElement),
      true,
    );
    await page.evaluate(() => (window.retried = true));
    await more.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.getAttribute("href") === "/result-12",
    );
    assert.equal(await results.count(), 24);
    await query.fill("slow");
    await page.waitForFunction(() => Boolean(window.releaseSearch));
    await query.fill("");
    await page.evaluate(() => window.releaseSearch());
    await page.waitForTimeout(220);
    assert.equal(await results.count(), 0);
    assert.equal(
      await page.locator("#docs-search [role=status]").textContent(),
      "",
    );

    // A denied sessionStorage getter must not break initialization or requests.
    for (const storageFailure of ["getter", "write"]) {
      const denied = await context.newPage();
      try {
        const errors = [];
        denied.on("pageerror", (error) => errors.push(error.message));
        await denied.addInitScript((failure) => {
          if (failure === "getter") {
            Object.defineProperty(window, "sessionStorage", {
              get() {
                throw new DOMException("Storage denied", "SecurityError");
              },
            });
          } else {
            Object.defineProperty(Storage.prototype, "setItem", {
              value() {
                throw new DOMException("Storage full", "QuotaExceededError");
              },
            });
          }
        }, storageFailure);
        await denied.goto(base + "/build/");
        await mockAnswer(denied, "Docker café 🐳", true);
        await ask(denied);
        const chat = await denied.evaluate(() => {
          const data = Alpine.$data(
            document.querySelector("#gordon-chat").parentElement,
          );
          return {
            content: data.messages.at(-1).content,
            threadId: data.threadId,
            answerId: data.messages.at(-1).questionAnswerId,
            announcement: data.announcement,
          };
        });
        assert.equal(chat.content, "Docker café 🐳");
        assert.equal(chat.threadId, "test-thread");
        assert.equal(chat.answerId, "test-answer");
        assert.equal(chat.announcement, "Ask Docker: Docker café 🐳");
        await ask(denied, "Follow up");
        const requests = await denied.evaluate(() => window.chatRequests);
        const sessionId = (request) =>
          request.messages.at(-1).copilot_references[0].data.uuid;
        assert.equal(sessionId(requests[0]), sessionId(requests[1]));
        assert.equal(requests[1].thread_uuid, "test-thread");
        assert.deepEqual(errors, []);
        await denied.evaluate(() =>
          Object.defineProperty(navigator.clipboard, "writeText", {
            configurable: true,
            value: async () => {
              throw Error("Clipboard denied");
            },
          }),
        );
        await denied.locator(".gordon-feedback-copy").last().click();
        await denied.waitForFunction(() =>
          [...document.querySelectorAll("#gordon-chat [role=status]")].some(
            (e) => e.textContent.includes("Could not copy"),
          ),
        );
      } finally {
        await denied.close();
      }
    }

    // Code blocks in answers, plain and inside cards, copy their code without prompts.
    const code = await context.newPage();
    try {
      await code
        .context()
        .grantPermissions(["clipboard-read", "clipboard-write"]);
      await code.goto(base + "/build/");
      await mockAnswer(
        code,
        "```console\n$ docker build .\n```\n\n" +
          '```openui-lang\nroot = Card("Install", [Steps([Step("Run", "On macOS.", "brew install docker/tap/sbx")])])\n```',
      );
      await ask(code);
      await code.waitForSelector("#gordon-chat .gordon-card");
      const buttons = code.locator(
        "#gordon-chat .gordon-code [data-copy-text]",
      );
      assert.deepEqual(
        await buttons.evaluateAll((es) => es.map((e) => e.dataset.copyText)),
        ["docker build .\n", "brew install docker/tap/sbx"],
      );
      for (const button of await buttons.all()) {
        assert.equal(
          await button.evaluate((e) => getComputedStyle(e).opacity),
          "1",
        );
        await button.focus();
        await code.keyboard.press("Enter");
        await code.waitForFunction(
          () =>
            document.activeElement
              .closest(".code-copy")
              .querySelector('[role="status"]').textContent === "Code copied.",
        );
        assert.equal(
          await code.evaluate(() => navigator.clipboard.readText()),
          await button.getAttribute("data-copy-text"),
        );
      }
    } finally {
      await code.close();
    }

    // A pending optional script must not hold the composer or completed announcement.
    const cards = await context.newPage();
    let pendingCardRoute;
    await cards.route("**/gordon-cards.js", (route) => {
      pendingCardRoute = route;
    });
    try {
      await cards.clock.install();
      await cards.goto(base + "/build/");
      await mockAnswer(
        cards,
        '```openui-lang\nroot = Card({title: "Example", children: []})\n```',
      );
      await ask(cards);
      await cards.waitForFunction(() =>
        document.querySelector('script[src$="gordon-cards.js"]'),
      );
      assert.equal(
        await cards.locator("#question").getAttribute("readonly"),
        null,
      );
      assert.match(
        await cards.locator("#gordon-chat > [role=status]").textContent(),
        /answer card is available/,
      );
      await ask(cards, "Another question while cards are pending");
      assert.equal(await cards.evaluate(() => window.chatRequests.length), 2);
      await cards.clock.fastForward(10001);
      await cards.waitForFunction(() => Alpine.store("gordonCards").failed);
      assert.ok(
        await cards
          .locator("#gordon-chat pre code.language-openui-lang")
          .count(),
      );
    } finally {
      if (pendingCardRoute) await pendingCardRoute.abort().catch(() => {});
      await cards.close();
    }

    // Restoring focus matters for dock input submission and its mobile button.
    for (const [width, pointer] of [
      [1440, false],
      [1440, true],
      [390, false],
    ]) {
      await page.setViewportSize({ width, height: 900 });
      await page.goto(base + "/build/");
      await mockAnswer(page, "Answer");
      const opener = page.locator(
        width > 700 ? ".ask-dock input" : ".ask-dock-mobile",
      );
      await opener.focus();
      if (width > 700) {
        await opener.fill("Dock question");
        if (pointer)
          await page.locator('.ask-dock-form button[type="submit"]').click();
        else await opener.press("Enter");
      } else await opener.press("Enter");
      await page.waitForFunction(
        () => document.activeElement.id === "question",
      );
      await page.keyboard.press("Escape");
      await page.waitForFunction(
        (selector) =>
          document.querySelector(selector) === document.activeElement,
        width > 700 ? ".ask-dock input" : ".ask-dock-mobile",
      );
    }

    await page.goto(base + "/engine/cli/otel/");
    const block = page
      .locator(".copy-container")
      .filter({
        has: page.locator('[x-ref="expand"]'),
      })
      .first();
    const expand = block.getByRole("button", {
      name: "Show more",
      exact: true,
    });
    await expand.focus();
    await expand.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.getAttribute("x-ref") === "code",
    );
    const hide = block.getByRole("button", { name: "Hide", exact: true });
    await hide.focus();
    await hide.press("Enter");
    await page.waitForFunction(
      () => document.activeElement.getAttribute("x-ref") === "expand",
    );
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
    const page = await browser.newPage();
    await verifySearchChatRegressions(
      page,
      process.env.DOCS_TEST_URL || "http://localhost:1313",
    );
    console.log("Search and chat regressions passed");
  } finally {
    await browser.close();
  }
}
