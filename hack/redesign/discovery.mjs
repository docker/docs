import assert from "node:assert/strict";

export default async function verifyDiscovery(page, base) {
  for (const [route, ancestors] of [
    ["/c/ai/", ["/"]],
    ["/engine/", ["/", "/c/dev/"]],
    ["/desktop/", ["/", "/c/dev/"]],
    ["/engine/install/", ["/", "/c/dev/", "/engine/"]],
    ["/ai/mcp-catalog-and-toolkit/", ["/", "/c/ai/"]],
    [
      "/reference/cli/docker/compose/",
      ["/", "/c/dev/", "/reference/", "/reference/cli/docker/"],
    ],
    [
      "/reference/api/hub/latest/",
      ["/", "/c/dev/", "/reference/", "/reference/api/"],
    ],
  ]) {
    await page.goto(base + route);
    const result = await page.evaluate(() => {
      const canonical = document.querySelector('link[rel="canonical"]').href;
      const visible = Array.from(
        document.querySelectorAll(
          '.breadcrumbs > a, .breadcrumbs > [aria-current="page"]',
        ),
        (element) => ({
          name: element.textContent.trim(),
          url: new URL(element.getAttribute("href") || canonical, canonical)
            .href,
        }),
      );
      const schemas = Array.from(
        document.querySelectorAll('script[type="application/ld+json"]'),
        (element) => JSON.parse(element.textContent),
      );
      return {
        visible,
        structured: schemas.find(
          (schema) => schema["@type"] === "BreadcrumbList",
        )?.itemListElement,
        parent: schemas.find((schema) => schema["@type"] === "TechArticle")
          ?.isPartOf,
        pagefind: document
          .querySelector('meta[data-pagefind-meta^="breadcrumbs:"]')
          ?.getAttribute("data-pagefind-meta"),
      };
    });
    assert.ok(result.visible.length >= 2, `${route}: visible breadcrumbs`);
    assert.deepEqual(
      result.visible.map(({ url }) => new URL(url).pathname),
      [...ancestors, route],
      `${route}: published category and collection ancestry`,
    );
    assert.deepEqual(
      result.structured?.map(({ position, item }, index) => {
        assert.equal(position, index + 1, `${route}: sequential positions`);
        return { name: item.name, url: item["@id"] };
      }),
      result.visible,
      `${route}: structured and visible breadcrumbs agree`,
    );
    assert.equal(
      new Set(result.visible.map(({ url }) => url)).size,
      result.visible.length,
      `${route}: no duplicate category or current page`,
    );
    const searchAncestors = result.visible.slice(1, -1);
    assert.equal(
      result.pagefind,
      searchAncestors.length
        ? `breadcrumbs:${searchAncestors.map(({ name }) => name).join(" / ")}`
        : undefined,
      `${route}: Pagefind contains only published ancestors`,
    );
    assert.ok(
      result.visible.every(({ url }) => !url.endsWith("/manuals/")),
      `${route}: no retired Manuals ancestor`,
    );
    if (searchAncestors.length) {
      assert.equal(result.parent?.["@id"], searchAncestors.at(-1).url);
    }
  }

  await page.goto(base + "/reference/");
  const references = await page
    .locator(".discovery-list a")
    .evaluateAll((links) =>
      links.map((link) => ({
        label: link.querySelector("strong").textContent.trim(),
        path: new URL(link.href).pathname,
      })),
    );
  assert.ok(references.length > 0);
  for (const path of ["/reference.md", "/llms.txt"]) {
    const response = await page.request.get(base + path);
    assert.ok(response.ok(), path);
    const text = await response.text();
    const section =
      path === "/llms.txt"
        ? text
            .split("## Reference entry points")[1]
            ?.split("## Retrieval guidance")[0]
        : text;
    assert.ok(section, `${path}: reference section`);
    const entries = Array.from(
      section.matchAll(/^- \[([^\]]+)\]\(([^)]+)\)/gm),
      ([, label, url]) => ({ label, path: new URL(url, base).pathname }),
    );
    assert.deepEqual(
      entries,
      references,
      `${path}: registry reference inventory`,
    );
  }
}
