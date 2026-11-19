import fs from "node:fs";
import path from "node:path";
import { isDeepStrictEqual } from "node:util";
import {
  decode,
  attributes,
  localLinkProblems,
  markdownLinks,
  engineOutputProblems,
} from "./output-checks.mjs";
const root = path.resolve(import.meta.dirname, "../..");
const base = path.resolve(process.argv[2] || "public");
const data = JSON.parse(
  fs.readFileSync(path.join(root, "data/api-reference.json")),
);
const redirects = JSON.parse(
  fs.readFileSync(path.join(base, "redirects.json")),
);
let count = 0;
const problems = [];
function check(url) {
  const htmlPath = path.join(base, url, "index.html");
  const mdPath = path.join(base, url.replace(/\/$/, "")) + ".md";
  if (!fs.existsSync(htmlPath) || !fs.existsSync(mdPath)) {
    problems.push(`Missing HTML/Markdown: ${url}`);
    return ["", ""];
  }
  const html = fs.readFileSync(htmlPath, "utf8"),
    md = fs.readFileSync(mdPath, "utf8");
  count++;
  return [decode(html), md];
}
check("/reference/api/");
for (const api of data.apis) {
  check(api.url);
  const publishedSource = fs.readFileSync(path.join(base, api.sourceURL));
  if (!publishedSource.equals(fs.readFileSync(path.join(root, api.source))))
    problems.push(`Published specification differs from source: ${api.id}`);

  for (const op of api.operations) {
    const [html, md] = check(op.url);
    for (const variant of op.variants) {
      if (
        !attributes(html, "data-api-variant").includes(variant.pointer) ||
        !md.includes(`### ${variant.direction} ${variant.status}`)
      )
        problems.push(`Missing variant: ${op.url} ${variant.pointer}`);
      for (const example of variant.examples) {
        if (
          !html.includes(example.text.trim()) ||
          !md.includes("```" + example.language + "\n" + example.text)
        )
          problems.push(`Example mismatch: ${op.url} ${variant.pointer}`);
      }
    }
    for (const p of op.parameters)
      if (
        !attributes(html, "data-api-parameter").includes(p.name) ||
        !md.includes(`### ${p.name}`)
      )
        problems.push(`Missing parameter: ${op.url} ${p.name}`);
    if (!html.includes(op.path) || !md.includes(op.method + " " + op.path))
      problems.push(`Missing operation signature: ${op.url}`);
    if (!md.includes(op.curl))
      problems.push(`Request example mismatch: ${op.url}`);
  }
  for (const schema of api.schemas) {
    const [html, md] = check(schema.url);
    if (
      ![...md.matchAll(/```json\n([\s\S]*?)\n```/g)].some((m) => {
        try {
          return isDeepStrictEqual(JSON.parse(m[1]), schema.schema);
        } catch {
          return false;
        }
      })
    )
      problems.push(`Incomplete Markdown schema: ${schema.url}`);
    if (!html.includes(schema.name))
      problems.push(`Missing schema heading: ${schema.url}`);
  }
}
// Check local links and fragments in generated API pages, including specification prose.
function files(dir) {
  return fs
    .readdirSync(dir, { withFileTypes: true })
    .flatMap((e) =>
      e.isDirectory()
        ? files(path.join(dir, e.name))
        : [path.join(dir, e.name)],
    );
}
for (const file of [...files(path.join(base, "reference/api"))].filter((p) =>
  p.endsWith(".html"),
)) {
  const html = fs.readFileSync(file, "utf8");
  if (!html.includes("data-api-view=")) continue;
  const article = html.match(
    /<article\b[^>]*class=["']?api-reference[\s\S]*?<\/article>/,
  )?.[0];
  if (!article) {
    problems.push(`Missing API article: ${path.relative(base, file)}`);
    continue;
  }
  problems.push(
    ...localLinkProblems(base, file, attributes(article, "href"), redirects),
  );
  const mdFile = path.dirname(file) + ".md";
  if (fs.existsSync(mdFile)) {
    problems.push(
      ...localLinkProblems(
        base,
        file,
        markdownLinks(fs.readFileSync(mdFile, "utf8")),
        redirects,
      ).map((problem) => `Markdown: ${problem}`),
    );
  }
}

if (fs.existsSync(path.join(base, "api-prototype")))
  problems.push("Unexpected prototype routes");
problems.push(
  ...engineOutputProblems(
    base,
    path.join(root, "content/reference/api/engine/version"),
    data.apis,
  ),
);

if (problems.length)
  throw Error(
    problems.slice(0, 25).join("\n") + `\n${problems.length} output failures`,
  );
const bytes = files(path.join(base, "reference/api")).reduce(
  (n, p) => n + fs.statSync(p).size,
  0,
);
fs.mkdirSync(path.join(root, "tmp/api-reference/reports"), { recursive: true });
fs.writeFileSync(
  path.join(root, "tmp/api-reference/reports/output.json"),
  JSON.stringify(
    { htmlMarkdownPagePairs: count, referenceBytes: bytes, result: "pass" },
    null,
    2,
  ) + "\n",
);
console.log(
  `Verified ${count} HTML/Markdown page pairs and API reference links.`,
);
