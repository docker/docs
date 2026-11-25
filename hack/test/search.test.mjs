import assert from "node:assert/strict";
import test from "node:test";

import createSearchState from "../../assets/js/docs-search.js";

function createSearch(api) {
  return createSearchState(async () => api);
}
const match = (
  index,
  data = async () => ({ url: `/${index}`, meta: { title: String(index) } }),
) => ({ data });

test("all matches can be loaded in batches without duplicates", async () => {
  const search = createSearch({
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
  const search = createSearch({
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
    const { promise: started, resolve: signalStarted } =
      Promise.withResolvers();
    let release;
    const search = createSearch({
      search: async () => ({
        results: [
          match(
            0,
            () =>
              new Promise((resolve) => {
                release = resolve;
                signalStarted();
              }),
          ),
        ],
      }),
    });
    search.query = "original";
    const pending = search.search();
    await started;
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
  const { promise: started, resolve: signalStarted } = Promise.withResolvers();
  let release;
  const search = createSearch({
    search: async (query) => ({
      results: [
        match(0, () =>
          query === "original"
            ? Promise.resolve({ url: "/original" })
            : new Promise((resolve) => {
                release = resolve;
                signalStarted();
              }),
        ),
      ],
    }),
  });
  search.query = "original";
  await search.search();
  search.queueSearch("replacement");
  assert.deepEqual(search.results, [{ url: "/original" }]);
  const pending = search.search();
  await started;
  assert.deepEqual(search.results, [{ url: "/original" }]);
  release({ url: "/replacement" });
  await pending;
  assert.deepEqual(search.results, [{ url: "/replacement" }]);
});

test("editing also invalidates an in-flight additional batch", async () => {
  let release;
  const search = createSearch({
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
