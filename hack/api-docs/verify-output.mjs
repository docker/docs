import fs from "node:fs";
import path from "node:path";
import { isDeepStrictEqual } from "node:util";
const root = path.resolve(import.meta.dirname, "../..");
const base = path.resolve(process.argv[2] || "public");
const data = JSON.parse(
  fs.readFileSync(path.join(root, "tmp/api-reference/data/api-reference.json")),
);
const decode = (s) =>
  s
    .replace(/&#(\d+);/g, (_, n) => String.fromCodePoint(Number(n)))
    .replaceAll("&amp;", "&")
    .replaceAll("&#34;", '"')
    .replaceAll("&#39;", "'")
    .replaceAll("&lt;", "<")
    .replaceAll("&gt;", ">");
// Hugo's production minifier removes optional attribute quotes.
function attributes(html, name) {
  const pattern = new RegExp(
    `\\s${name}=(?:"([^"]*)"|'([^']*)'|([^\\s>]+))`,
    "g",
  );
  return [...html.matchAll(pattern)].map((m) => decode(m[1] ?? m[2] ?? m[3]));
}
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
// Check API reference links in every generated preview page, including copied specification prose.
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
  for (const href of attributes(html, "href")) {
    let url;
    try {
      url = new URL(href, "http://localhost:1314/" + path.relative(base, file));
    } catch {
      continue;
    }
    if (
      url.hostname !== "localhost" ||
      !url.pathname.startsWith("/reference/api/")
    )
      continue;
    const target = path.join(
      base,
      decodeURIComponent(url.pathname),
      url.pathname.endsWith("/") ? "index.html" : "",
    );
    if (!fs.existsSync(target))
      problems.push(`Broken API link: ${path.relative(base, file)} -> ${href}`);
  }
}
const governance = fs.readFileSync(
  path.join(base, "reference/api/ai-governance/index.html"),
  "utf8",
);
if (governance.includes("data-api-view="))
  problems.push("Governance renderer changed");
if (fs.existsSync(path.join(base, "api-prototype")))
  problems.push("Unexpected prototype routes");
// Historical Engine pages must keep the original renderer at their existing URLs.
for (let minor = 40; minor <= 56; minor++) {
  const url = `/reference/api/engine/version/v1.${minor}/`;
  const html = fs.readFileSync(path.join(base, url, "index.html"), "utf8");
  if (!html.includes("<redoc") || html.includes("data-api-view="))
    problems.push(`Legacy Engine renderer changed: ${url}`);
}
if (problems.length)
  throw Error(
    problems.slice(0, 25).join("\n") + `\n${problems.length} output failures`,
  );
const bytes = files(path.join(base, "reference/api")).reduce(
  (n, p) => n + fs.statSync(p).size,
  0,
);
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
