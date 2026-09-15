import { test } from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { execFileSync } from "node:child_process";

const script = path.resolve(import.meta.dirname, "../flatten-and-resolve.js");

test("flattens all pages while preserving examples and resolving documentation links", () => {
  const site = fs.mkdtempSync(path.join(os.tmpdir(), "docs-flatten-"));
  const generated =
    "# API\n\n" +
    "[Guide](/reference/api/registry/auth/)\n\n" +
    '```json\n{"path":"/manuals/example/index.md","text":"[sample](file.md)"}\n```\n';
  const write = (name, value) => {
    const file = path.join(site, name);
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, value);
  };
  try {
    write("reference/api/index.md", generated);
    write("reference/api/registry/latest/index.md", generated);
    write(
      "reference/api/registry/auth/index.md",
      "# Guide\n\n[Example](example.md)\n",
    );
    write("reference/api/registry/auth/example.md", "# Example\n");
    for (let run = 0; run < 2; run++) {
      execFileSync(process.execPath, [script, site], { stdio: "pipe" });
      assert.equal(
        fs.readFileSync(path.join(site, "reference/api.md"), "utf8"),
        generated,
      );
      assert.equal(
        fs.readFileSync(
          path.join(site, "reference/api/registry/latest.md"),
          "utf8",
        ),
        generated,
      );
      assert.ok(
        !fs.existsSync(
          path.join(site, "reference/api/registry/latest/index.md"),
        ),
      );
      assert.match(
        fs.readFileSync(
          path.join(site, "reference/api/registry/auth.md"),
          "utf8",
        ),
        /\[Example\]\(\/reference\/api\/registry\/auth\/example\/\)/,
      );
    }
  } finally {
    fs.rmSync(site, { recursive: true, force: true });
  }
});

test("rewrites only destinations and preserves Markdown syntax", () => {
  const site = fs.mkdtempSync(path.join(os.tmpdir(), "docs-links-"));
  const cases = [
    ["[Guide](/manuals/engine/_index.md#install)", "[Guide](/engine/#install)"],
    ["[Parent](../index.md)", "[Parent](/guide/)"],
    ["[Root](/index.md)", "[Root](/)"],
    [
      '[Title](<file.md> "A title")',
      '[Title](</guide/section/file/> "A title")',
    ],
    [
      "[Nested **label**](file(test).md?q=1#part)",
      "[Nested **label**](/guide/section/file(test)/?q=1#part)",
    ],
    ["[`odd ] bracket`](file.md)", "[`odd ] bracket`](/guide/section/file/)"],
    [
      "![Image [label]](picture.png)",
      "![Image [label]](/guide/section/picture.png)",
    ],
    [
      '[ref]: /manuals/engine/index.md "Title"\n\n[Reference][ref]',
      '[ref]: /engine/ "Title"\n\n[Reference][ref]',
    ],
    ["> - [Guide](file.md)", "> - [Guide](/guide/section/file/)"],
    [
      "[External](https://example.com/manuals/index.md)",
      "[External](https://example.com/manuals/index.md)",
    ],
    [
      "[CDN](//example.com/manuals/index.md)",
      "[CDN](//example.com/manuals/index.md)",
    ],
    ["[Mail](mailto:docs@example.com)", "[Mail](mailto:docs@example.com)"],
    ["[Anchor](#example)", "[Anchor](#example)"],
    ["`[Example](/manuals/index.md)`", "`[Example](/manuals/index.md)`"],
    ["    [Example](/manuals/index.md)", "    [Example](/manuals/index.md)"],
    [
      '~~~json\n{"file":"/manuals/index.md"}\n~~~',
      '~~~json\n{"file":"/manuals/index.md"}\n~~~',
    ],
    [
      "A literal /manuals/example/index.md filename.",
      "A literal /manuals/example/index.md filename.",
    ],
  ];
  try {
    fs.mkdirSync(path.join(site, "guide/section"), { recursive: true });
    fs.writeFileSync(
      path.join(site, "guide/section/index.md"),
      cases.map(([before]) => before).join("\n\n"),
    );
    execFileSync(process.execPath, [script, site]);
    assert.equal(
      fs.readFileSync(path.join(site, "guide/section.md"), "utf8"),
      cases.map(([, after]) => after).join("\n\n"),
    );
  } finally {
    fs.rmSync(site, { recursive: true, force: true });
  }
});
