import { test } from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { execFileSync, spawnSync } from "node:child_process";

const repo = path.resolve(import.meta.dirname, "../..");
const script = path.join(repo, "hack/flatten-and-resolve.js");

function fixture(t) {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "docs-markdown-"));
  t.after(() => fs.rmSync(root, { recursive: true, force: true }));
  const write = (name, content) => {
    const file = path.join(root, name);
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, content);
  };
  return { root, write };
}

test("Hugo resolves Markdown from source pages, bundles, and expanded shortcodes", (t) => {
  const { root, write } = fixture(t);
  for (const dir of ["_partials/markdown"]) {
    fs.cpSync(
      path.join(repo, "layouts", dir),
      path.join(root, "layouts", dir),
      { recursive: true },
    );
  }
  for (const name of ["render-link", "render-image"]) {
    write(
      `layouts/_markup/${name}.markdown.md`,
      fs.readFileSync(path.join(repo, `layouts/_markup/${name}.markdown.md`)),
    );
  }
  for (const name of ["button", "inline-image"]) {
    write(
      `layouts/_shortcodes/${name}.markdown.md`,
      fs.readFileSync(
        path.join(repo, `layouts/_shortcodes/${name}.markdown.md`),
      ),
    );
  }
  write(
    "content/manuals/_index.md",
    "---\ntitle: Manuals\nbuild:\n  render: never\n---\n",
  );
  write(
    "hugo.toml",
    `baseURL = "https://docs.docker.com/"
[outputs]
home = ["html", "markdown"]
page = ["html", "markdown"]
section = ["html", "markdown"]
[permalinks.page]
manuals = "/:sections[1:]/:slugorcontentbasename/"
[permalinks.section]
manuals = "/:sections[1:]/"
[markup.goldmark.renderer]
unsafe = true
`,
  );
  for (const name of ["single", "list", "home"]) {
    write(`layouts/${name}.html`, "<title>{{ .Title }}</title>");
    write(
      `layouts/${name}.markdown.md`,
      `{{- partial "markdown/render.md" (dict "page" . "content" (partial "markdown/body/${name}.md" .)) -}}
{{ define "_partials/markdown/body/${name}.md" -}}
# {{ .Title }}

{{ partial "markdown/content.html" . }}
{{ end -}}`,
    );
  }
  write(
    "layouts/cli.markdown.md",
    fs.readFileSync(path.join(repo, "layouts/cli.markdown.md")),
  );
  write(
    "data/cli/fixture/command.yaml",
    `short: Example command
usage: example
long: 'Read [the parent](/manuals/guide/parent.md).'
options:
  - option: target
    hidden: false
    description: 'Read [the target](/manuals/guide/nested/sibling.md).'
  - option: hidden
    hidden: true
    description: '[Unpublished](/does-not-exist.md)'
examples: |
  \`\`\`json
  {"text": "<example>", "path": "/manuals/example/index.md"}
  \`\`\`
`,
  );
  write(
    "data/cli/fixture/no-options.yaml",
    "short: No options\nusage: no-options\n",
  );
  write(
    "content/reference/no-options.md",
    "---\ntitle: No options\nlayout: cli\nparams:\n  datafolder: fixture\n  datafile: no-options\n---\n",
  );
  write(
    "content/reference/command.md",
    "---\ntitle: Command\nlayout: cli\nparams:\n  datafolder: fixture\n  datafile: command\n---\n",
  );
  write(
    "content/_index.md",
    '---\ntitle: Home\n---\n<img src="images/home.png">\n',
  );
  const page = (name, body, extra = "") =>
    write(
      `content/manuals/guide/${name}`,
      `---\ntitle: Example\n${extra}---\n${body}\n`,
    );
  page(
    "_index.md",
    "[Parent](parent.md)\n![Section image](images/section.png)",
  );
  page("nested/_index.md", "Nested");
  write("content/manuals/other/_index.md", "---\ntitle: Other\n---\n");
  page("parent.md", "Target");
  page("log_tags.md", "Target");
  page("nested/sibling.md", "Target", "url: /custom/sibling/\n");
  write("content/manuals/other/sibling.md", "---\ntitle: Other sibling\n---\n");
  page("bundle/index.md", "![Bundle](asset.png)");
  page(
    "stub.md",
    "",
    "params:\n  sidebar:\n    goto: https://example.com/releases\n",
  );
  const code =
    '```json\n{"text":"[Example](missing.md)","src":"relative.png"}\n```';
  page(
    "nested/ordinary.md",
    `[Sibling](sibling.md?view=raw#target)
[Source](/manuals/guide/nested/sibling.md)
[Parent](../parent.md)
[Lookup](log_tags.md)
[Reference][ref]
[Same page](#part)
[Query](?view=raw#part)
![Image](images/figure.png?width=10#preview)
{{< button text="Button" url="../parent.md" >}}
{{< inline-image src="images/inline.png" alt="Inline" >}}
<img src="../images/raw.svg?x=1&amp;y=2">
<img src='../images/single.svg'>
<script>const example = 'src="untouched.png"';</script>
[External](https://example.com/path.md)
${code}

[ref]: sibling.md#target
[unused]: unused.md
`,
  );
  execFileSync(process.env.HUGO_BIN || "hugo", ["--source", root], {
    stdio: "pipe",
  });
  const publicDir = path.join(root, "public");
  for (let run = 0; run < 2; run++) {
    execFileSync(process.execPath, [script, publicDir], { stdio: "pipe" });
    const read = (file) => fs.readFileSync(path.join(publicDir, file), "utf8");
    const body = read("guide/nested/ordinary.md");
    assert.ok(
      body.includes(
        "[Sibling](https://docs.docker.com/custom/sibling/?view=raw#target)",
      ),
    );
    assert.ok(
      body.includes("[Source](https://docs.docker.com/custom/sibling/)"),
    );
    assert.ok(body.includes("[Parent](https://docs.docker.com/guide/parent/)"));
    assert.ok(
      body.includes("[Lookup](https://docs.docker.com/guide/log_tags/)"),
    );
    assert.ok(
      body.includes("[ref]: https://docs.docker.com/custom/sibling/#target"),
    );
    assert.ok(body.includes("[Same page](#part)"));
    assert.ok(body.includes("[Query](?view=raw#part)"));
    assert.ok(
      body.includes(
        "https://docs.docker.com/guide/nested/images/figure.png?width=10#preview",
      ),
    );
    assert.ok(body.includes("[Button](https://docs.docker.com/guide/parent/)"));
    assert.ok(
      body.includes("https://docs.docker.com/guide/nested/images/inline.png"),
    );
    assert.ok(
      body.includes('<img src="/guide/nested/images/raw.svg?x=1&amp;y=2">'),
    );
    assert.ok(body.includes("<img src='/guide/nested/images/single.svg'>"));
    assert.ok(body.includes(`const example = 'src="untouched.png"'`));
    assert.ok(body.includes(code));
    assert.ok(body.includes("[External](https://example.com/path.md)"));
    assert.ok(body.includes("[unused]: unused.md"));
    assert.match(
      read("guide.md"),
      /https:\/\/docs.docker.com\/guide\/images\/section.png/,
    );
    assert.match(
      read("guide/bundle.md"),
      /https:\/\/docs.docker.com\/guide\/bundle\/asset.png/,
    );
    assert.match(
      read("guide/stub.md"),
      /See \[Example\]\(https:\/\/example.com\/releases\)/,
    );
    assert.ok(read("index.md").includes('<img src="/images/home.png">'));
    const cli = read("reference/command.md");
    assert.ok(
      cli.includes(
        '{"text": "<example>", "path": "/manuals/example/index.md"}',
      ),
    );
    assert.ok(
      cli.includes("[the parent](https://docs.docker.com/guide/parent/)"),
    );
    assert.ok(
      cli.includes("[the target](https://docs.docker.com/custom/sibling/)"),
    );
    assert.ok(!fs.existsSync(path.join(publicDir, "__markdown-links")));
  }
});

test("missing Hugo metadata fails before changing or flattening any file", (t) => {
  const { root, write } = fixture(t);
  write("a/index.md", "[Root](/known/)\n");
  write("b/index.md", "[Missing](missing.md)\n");
  const result = spawnSync(process.execPath, [script, root], {
    encoding: "utf8",
  });
  assert.notEqual(result.status, 0);
  assert.match(
    result.stderr,
    /No Hugo resolution for b\/index.md -> missing.md/,
  );
  assert.ok(fs.existsSync(path.join(root, "a/index.md")));
  assert.ok(!fs.existsSync(path.join(root, "a.md")));
});

test("rewrites destination spans without changing labels, titles, entities, or code examples", (t) => {
  const { root, write } = fixture(t);
  const cases = [
    ['[Title](<file.md> "A title")', '[Title](</target/> "A title")'],
    [
      "[Nested **label**](file(test).md?q=1#part)",
      "[Nested **label**](/target/?q=1#part)",
    ],
    ["[`odd ] bracket`](file.md)", "[`odd ] bracket`](/target/)"],
    ["![Image [label]](picture.png)", "![Image [label]](/assets/picture.png)"],
    [
      '[ref]: /manuals/engine/index.md "Title"\n\n[Reference][ref]',
      '[ref]: /engine/ "Title"\n\n[Reference][ref]',
    ],
    ["> - [Guide](file.md)", "> - [Guide](/target/)"],
    ["[Escaped](file\\(name\\).md?x=1&amp;y=2)", "[Escaped](/target/?x=1&y=2)"],
    [
      "[CDN](//example.com/manuals/index.md)",
      "[CDN](//example.com/manuals/index.md)",
    ],
    ["[Mail](mailto:docs@example.com)", "[Mail](mailto:docs@example.com)"],
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
  write("guide/index.md", cases.map(([before]) => before).join("\n\n"));
  write(
    "__markdown-links/guide/index.md.json",
    JSON.stringify({
      links: {
        "file.md": "/target/",
        "file(test).md?q=1#part": "/target/?q=1#part",
        "/manuals/engine/index.md": "/engine/",
        "file(name).md?x=1&y=2": "/target/?x=1&y=2",
      },
      images: { "picture.png": "/assets/picture.png" },
    }),
  );
  execFileSync(process.execPath, [script, root]);
  assert.equal(
    fs.readFileSync(path.join(root, "guide.md"), "utf8"),
    cases.map(([, after]) => after).join("\n\n"),
  );
});
