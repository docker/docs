import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { spawnSync } from "node:child_process";
import { verifyMarkdown } from "../verify-markdown.mjs";

function fixture(t) {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "docs-md-check-"));
  t.after(() => fs.rmSync(root, { recursive: true, force: true }));
  const write = (name, text) => {
    const file = path.join(root, name);
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, text);
  };
  write(
    "redirects.json",
    JSON.stringify({ "/old/": "/target/", "/learning/": "/learn/example/" }),
  );
  write("llms-full.txt", "Markdown: https://docs.docker.com/index.md\n");
  write("target/index.html", "<h1>Target</h1>");
  write("images/figure.png", "image");
  return { root, write };
}

test("checks links, reference images, redirects, and advertised Markdown without matching examples", async (t) => {
  const { root, write } = fixture(t);
  write(
    "index.md",
    `[Target](/target/?q=1#section)
[Alias](/old/)
![Figure][figure]
<img src="/images/figure.png">
[Local service](http://localhost:9090/targets)
[External](https://example.com/missing/)

[figure]: /images/figure.png

\`src=missing.png\`

\`\`\`html
<a href="/missing/"><img src="/missing.png"></a>
\`\`\`
`,
  );
  assert.deepEqual(await verifyMarkdown(root), {
    files: 1,
    remote: 0,
    problems: [],
  });
});

test("a broken link, image, or advertised Markdown URL fails the CLI", (t) => {
  const { root, write } = fixture(t);
  write("index.md", "[Broken](/missing/)\n![Broken](/missing.png)\n");
  write("llms-full.txt", "Markdown: https://docs.docker.com/missing.md\n");
  const result = spawnSync(
    process.execPath,
    [path.resolve(import.meta.dirname, "../verify-markdown.mjs"), root],
    { encoding: "utf8" },
  );
  assert.notEqual(result.status, 0);
  assert.match(
    result.stderr,
    /Broken Markdown link: index.md -> .*\/missing\//,
  );
  assert.match(
    result.stderr,
    /Broken Markdown link: index.md -> .*\/missing.png/,
  );
  assert.match(result.stderr, /Missing advertised Markdown: .*\/missing.md/);
});

test("separately published Learn URLs are checked and deduplicated after redirects", async (t) => {
  const { root, write } = fixture(t);
  write(
    "index.md",
    "[Learn](/learning/)\n[Same lesson](/learn/example/#section)\n",
  );
  const requested = [];
  const result = await verifyMarkdown(root, {
    request: async (url) => {
      requested.push(url);
      return new Response("missing", { status: 404 });
    },
  });
  assert.deepEqual(requested, ["https://docs.docker.com/learn/example/"]);
  assert.equal(result.problems.length, 1);
  assert.match(
    result.problems[0],
    /index.md -> .*\/learn\/example\/ \(HTTP 404\)/,
  );
});

test("preview URLs are checked locally while localhost examples remain external", async (t) => {
  const { root, write } = fixture(t);
  write("llms-full.txt", "Markdown: https://preview.example.com/index.md\n");
  write(
    "index.md",
    "[Target](https://preview.example.com/target/)\n[Broken](https://preview.example.com/missing/)\n[Service](http://localhost:8080/)\n",
  );
  const result = await verifyMarkdown(root);
  assert.equal(result.problems.length, 1);
  assert.match(result.problems[0], /preview.example.com\/missing/);
});
