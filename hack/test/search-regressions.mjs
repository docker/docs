import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

const source = await readFile(
  new URL("../../assets/js/docs-search.js", import.meta.url),
  "utf8",
);
let moduleId = 0;

// Supply Pagefind at the module boundary; exercise the production search state.
async function createSearch(api) {
  const injected = source.replace(
    'import("/pagefind/pagefind.js")',
    "Promise.resolve(globalThis.searchRegressionApi)",
  );
  globalThis.searchRegressionApi = { options: async () => {}, ...api };
  const { default: create } = await import(
    `data:text/javascript;base64,${Buffer.from(injected + `\n// ${++moduleId}`).toString("base64")}`
  );
  return create();
}
const match = (
  index,
  data = async () => ({ url: `/${index}`, meta: { title: String(index) } }),
) => ({ data });

test("typing waits for a pause and searches only the final query", async (t) => {
  t.mock.timers.enable({ apis: ["setTimeout"] });
  const search = await createSearch({});
  const queries = [];
  search.search = () => queries.push(search.query);
  for (const query of ["d", "do", "doc", "docker"]) {
    search.queueSearch(query);
    t.mock.timers.tick(250);
    assert.deepEqual(queries, []);
  }
  t.mock.timers.tick(150);
  assert.deepEqual(queries, ["docker"]);
  search.queueSearch("discard");
  search.queueSearch("");
  t.mock.timers.tick(400);
  assert.deepEqual(queries, ["docker"]);
});

test("Enter searches immediately without repeating the queued search", async (t) => {
  t.mock.timers.enable({ apis: ["setTimeout"] });
  const queries = [];
  const search = await createSearch({
    search: async (query) => {
      queries.push(query);
      return { results: [] };
    },
  });
  search.queueSearch("docker");
  search.onEnter({ preventDefault() {}, stopPropagation() {} });
  // Wait for the asynchronous Pagefind initialization and search to settle.
  while (search.loading) await Promise.resolve();
  assert.deepEqual(queries, ["docker"]);
  t.mock.timers.tick(400);
  assert.deepEqual(queries, ["docker"]);
});

test("composition waits for committed text and preserves the IME Enter key", async (t) => {
  t.mock.timers.enable({ apis: ["setTimeout"] });
  const search = await createSearch({});
  const queries = [];
  search.search = () => queries.push(search.query);
  search.queueSearch("previous");
  search.queueSearch("正在输入", true);
  search.onEnter({ isComposing: true });
  search.onEnter({ keyCode: 229 });
  t.mock.timers.tick(1000);
  assert.deepEqual(queries, []);
  search.queueSearch("容器");
  t.mock.timers.tick(400);
  assert.deepEqual(queries, ["容器"]);
});

test("Ctrl+Enter and Command+Enter keep the Ask Docker shortcut", async () => {
  const search = await createSearch({});
  let asked = 0;
  search.ask = () => ++asked;
  for (const modifier of ["ctrlKey", "metaKey"]) {
    search.onEnter({
      [modifier]: true,
      preventDefault() {},
      stopPropagation() {},
    });
  }
  assert.equal(asked, 2);
});

test("editing while Pagefind initializes skips the superseded search", async () => {
  let release;
  const queries = [];
  const search = await createSearch({
    options: () => new Promise((resolve) => (release = resolve)),
    search: async (query) => {
      queries.push(query);
      return { results: [] };
    },
  });
  search.query = "original";
  const pending = search.search();
  while (!release) await Promise.resolve();
  search.queueSearch("replacement");
  release();
  await pending;
  assert.deepEqual(queries, []);
  search.destroy();
});

test("all matches can be loaded in batches without duplicates", async () => {
  const search = await createSearch({
    search: async () => ({
      results: Array.from({ length: 25 }, (_, i) => match(i)),
    }),
  });
  search.query = "docker";
  await search.search();
  assert.equal(search.results.length, 12);
  assert.equal(search.status, "Showing 12 of 25 results");
  await search.loadMore();
  assert.equal(search.results.length, 24);
  await search.loadMore();
  assert.equal(search.results.length, 25);
  assert.equal(new Set(search.results.map((result) => result.url)).size, 25);
  assert.equal(search.status, "Showing 25 of 25 results");
});

test("a failed additional batch preserves loaded results and can be retried", async () => {
  let fail = true;
  const search = await createSearch({
    search: async () => ({
      results: Array.from({ length: 13 }, (_, i) =>
        match(i, async () => {
          if (i === 12 && fail) throw Error("Network failure");
          return { url: `/${i}` };
        }),
      ),
    }),
  });
  search.query = "docker";
  await search.search();
  await search.loadMore();
  assert.equal(search.results.length, 12);
  assert.equal(search.loading, false);
  assert.match(search.status, /Could not load/);
  fail = false;
  await search.loadMore();
  assert.equal(search.results.length, 13);
});

test("editing or clearing invalidates results before the debounce expires", async () => {
  for (const query of ["", "replacement"]) {
    let release;
    const search = await createSearch({
      search: async () => ({
        results: [
          match(
            0,
            () =>
              new Promise((resolve) => {
                release = resolve;
              }),
          ),
        ],
      }),
    });
    search.query = "original";
    const pending = search.search();
    while (!release) await Promise.resolve();
    search.queueSearch(query);
    release({ url: "/old" });
    await pending;
    assert.deepEqual(search.results, []);
    assert.equal(search.status, query.trim() ? "Searching…" : "");
    assert.equal(search.loading, Boolean(query.trim()));
    search.destroy();
  }
});

test("visible results stay in place until a replacement batch is ready", async () => {
  let release;
  const search = await createSearch({
    search: async (query) => ({
      results: [
        match(0, () =>
          query === "original"
            ? Promise.resolve({ url: "/original" })
            : new Promise((resolve) => (release = resolve)),
        ),
      ],
    }),
  });
  search.query = "original";
  await search.search();
  const originalResults = search.results;
  const originalMatches = search.matches;
  search.queueSearch("replacement");
  assert.equal(search.results, originalResults);
  assert.equal(search.matches, originalMatches);
  assert.equal(search.status, "Showing 1 of 1 results");
  assert.equal(search.loading, true);
  const pending = search.search();
  while (!release) await Promise.resolve();
  assert.equal(search.results, originalResults);
  assert.equal(search.matches, originalMatches);
  release({ url: "/replacement" });
  await pending;
  assert.deepEqual(search.results, [{ url: "/replacement" }]);
  assert.notEqual(search.matches, originalMatches);
  assert.equal(search.loading, false);
  search.queueSearch("   ");
  assert.deepEqual(search.results, []);
  assert.deepEqual(search.matches, []);
  assert.equal(search.status, "");
  assert.equal(search.loading, false);
});

test("editing also invalidates an in-flight additional batch", async () => {
  let release;
  const search = await createSearch({
    search: async () => ({
      results: Array.from({ length: 13 }, (_, i) =>
        match(i, () =>
          i === 12
            ? new Promise((resolve) => {
                release = resolve;
              })
            : Promise.resolve({ url: `/${i}` }),
        ),
      ),
    }),
  });
  search.query = "docker";
  await search.search();
  const pending = search.loadMore();
  search.queueSearch("");
  release({ url: "/old" });
  await pending;
  assert.deepEqual(search.results, []);
  assert.deepEqual(search.matches, []);
});
