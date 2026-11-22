import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import {
  engineOutputProblems,
  localLinkProblems,
  markdownLinks,
} from "./output-checks.mjs";

function fixture(t) {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "api-output-"));
  t.after(() => fs.rmSync(root, { recursive: true, force: true }));
  const write = (name, text) => {
    const file = path.join(root, name);
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, text);
    return file;
  };
  return { root, write };
}

test("Markdown fragments use published HTML anchors, including through redirects", (t) => {
  const { root, write } = fixture(t);
  const file = write(
    "reference/api/example/operations/get/index.html",
    '<article id="parameters"></article>',
  );
  write(
    "reference/api/example/index.html",
    "<section id=connection></section>",
  );
  const markdown = `
[Connection](/reference/api/example/#connection)
[Alias][overview]
[Parameters](#parameters)
[Missing fragment](/reference/api/example/#authentication)
[Missing page](/reference/api/missing/)

[overview]: https://docs.docker.com/old-example/#connection

\`[Not a link](/inline-code/)\`

\`\`\`json
{"description":"[Not a link](/code-example/)"}
\`\`\`
`;
  const hrefs = markdownLinks(markdown);
  assert.equal(hrefs.length, 5);
  const problems = localLinkProblems(root, file, hrefs, {
    "/old-example/": "/reference/api/example/",
  });
  assert.equal(problems.length, 2);
  assert.match(problems[0], /Broken API anchor: .*#authentication/);
  assert.match(problems[1], /Broken API link: .*missing/);
});

test("historical Engine output is expected from source pages even if its entire output directory is missing", (t) => {
  const { root, write } = fixture(t);
  const sources = path.dirname(
    write("source/v1.56.md", "---\ntitle: Engine\n---\n"),
  );
  write("source/v1.57.md", "---\ntitle: Engine\n---\n");
  write("source/_index.md", "Versions");
  const legacy = "/reference/api/engine/version/v1.56/";
  const generated = "/reference/api/engine/version/v1.57/";
  const legacyHTML = `<redoc spec-url="${legacy.slice(0, -1)}.yaml"></redoc><script src="https://cdn.redoc.ly/redoc/v2.5.1/bundles/redoc.standalone.js"></script>`;
  const legacyFile = write(legacy + "index.html", legacyHTML);
  write(legacy.slice(0, -1) + ".yaml", "swagger: '2.0'");
  const generatedFile = write(
    generated + "index.html",
    '<article data-api-view="overview"></article>',
  );
  write(generated.slice(0, -1) + ".yaml", "openapi: 3.2.0");
  const apis = [{ url: generated }];
  assert.deepEqual(engineOutputProblems(root, sources, apis), []);
  fs.rmSync(path.dirname(legacyFile), { recursive: true });
  assert.deepEqual(engineOutputProblems(root, sources, apis), [
    `Missing Engine page: ${legacy}`,
  ]);
  write(legacy + "index.html", legacyHTML);
  fs.writeFileSync(generatedFile, legacyHTML);
  assert.deepEqual(engineOutputProblems(root, sources, apis), [
    `Generated Engine renderer missing: ${generated}`,
  ]);
});
