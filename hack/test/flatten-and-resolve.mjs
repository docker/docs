import { test } from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { execFileSync } from "node:child_process";

const script = path.resolve(import.meta.dirname, "../flatten-and-resolve.js");

test("flattens generated pages unchanged while resolving handwritten reference links", () => {
  const site = fs.mkdtempSync(path.join(os.tmpdir(), "docs-flatten-"));
  const generated =
    "<!-- link-rewriting: off -->\n\n# API\n\n" +
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
